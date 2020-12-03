package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"ravblk/exercises/service/midlware/instrumenting"
	service "ravblk/exercises/service/services"
	"ravblk/exercises/service/services/brackets"
	"syscall"

	"ravblk/exercises/service/transports"

	"ravblk/exercises/service/endpoints"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	metricsPort = ":8080"
	grpcPort    = ":3003"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		return
	}

	defer logger.Sync()

	var bracketsService service.Brackets

	bracketsService = brackets.NewService(logger)
	bracketsService = instrumenting.Middleware(bracketsService)(bracketsService)

	bracketsEndpoint := endpoints.MakeEndpoints(bracketsService)
	grpcServer := transports.NewGRPCServer(bracketsEndpoint)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Error("listen err", zap.Error(err))
		return
	}

	go func() {
		baseServer := grpc.NewServer()
		brackets.RegisterBracketsServiceServer(baseServer, grpcServer)
		logger.Info("server started")
		baseServer.Serve(grpcListener)
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Info("server metrics started")
		errs <- http.ListenAndServe(metricsPort, nil)
	}()

	err = <-errs
	logger.Error("server err", zap.Error(err))
}

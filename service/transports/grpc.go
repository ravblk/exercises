package transports

import (
	"context"
	"ravblk/exercises/service/endpoints"

	"ravblk/exercises/service/services/brackets"

	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	fix      gt.Handler
	validate gt.Handler
}

func NewGRPCServer(endpoints endpoints.Endpoints) *gRPCServer {
	return &gRPCServer{
		fix: gt.NewServer(
			endpoints.Fix,
			decodeBracketsRequest,
			encodeFixResponse,
		),
		validate: gt.NewServer(
			endpoints.Validate,
			decodeBracketsRequest,
			encodeValidateResponse,
		),
	}
}

func (s *gRPCServer) Fix(ctx context.Context, req *brackets.Brackets) (*brackets.ResultFixBrackets, error) {
	_, resp, err := s.fix.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*brackets.ResultFixBrackets), nil
}

func (s *gRPCServer) Validate(ctx context.Context, req *brackets.Brackets) (*brackets.ResultValidateBrackets, error) {
	_, resp, err := s.validate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*brackets.ResultValidateBrackets), nil
}

func decodeBracketsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*brackets.Brackets)
	return req, nil
}

func encodeFixResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*brackets.ResultFixBrackets)
	return resp, nil
}

func encodeValidateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*brackets.ResultValidateBrackets)
	return resp, nil
}

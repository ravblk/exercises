package instrumenting

import (
	"context"
	service "ravblk/exercises/service/services"
	"ravblk/exercises/service/services/brackets"
	"time"

	"github.com/go-kit/kit/metrics"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
)

func Middleware(
	svc service.Brackets,
) service.BracketsMiddleware {
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "metrics",
		Subsystem: "brackets_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: "metrics",
		Subsystem: "brackets_service",
		Name:      "request_latency_microseconds",
		Help:      "Duration of requests in microseconds.",
		Buckets:   stdprometheus.LinearBuckets(1, 10, 10),
	}, fieldKeys)

	return func(next service.Brackets) service.Brackets {
		return &instrumenting{requestCount, requestLatency, next}
	}
}

type instrumenting struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	svc            service.Brackets
}

func (i *instrumenting) Fix(ctx context.Context, in *brackets.Brackets) (*brackets.ResultFixBrackets, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "fix"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds() * 1_000_000)
	}(time.Now())

	return i.svc.Fix(ctx, in)
}

func (i *instrumenting) Validate(ctx context.Context, in *brackets.Brackets) (*brackets.ResultValidateBrackets, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "validate"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds() * 1_000_000)
	}(time.Now())

	return i.svc.Validate(ctx, in)
}

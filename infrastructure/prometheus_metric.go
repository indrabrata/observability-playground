package infrastructure

import (
	"context"

	"github.com/indrabrata/observability-playground/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func NewPrometheusMetric(ctx context.Context) *prometheus.Registry {
	promReg := prometheus.NewRegistry()

	TotalRequest := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "total_requests",
		Help: "Total number of requests",
	}, []string{"method", "endpoint", "status"})
	Latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_latency",
		Help:    "Request latency distribution",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint"})

	middleware.TotalRequest = TotalRequest
	middleware.Latency = Latency

	promReg.MustRegister(TotalRequest, Latency)
	return promReg
}

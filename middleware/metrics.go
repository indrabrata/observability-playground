package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/indrabrata/observability-playground/common"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	totalRequest *prometheus.CounterVec
	latency      *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		totalRequest: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "total_requests",
			Help: "Total number of requests",
		}, []string{"method", "endpoint", "status"}),
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "request_latency",
			Help:    "Request latency distribution",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "endpoint"}),
	}

	reg.MustRegister(m.totalRequest, m.latency)
	return m
}

func (m *Metrics) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := common.NewInterceptor(w)
		next.ServeHTTP(crw, r)

		elapsed := time.Since(start).Milliseconds()
		m.totalRequest.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(crw.StatusCode)).Inc()
		m.latency.WithLabelValues(r.Method, r.URL.Path).Observe(float64(elapsed))
	})
}

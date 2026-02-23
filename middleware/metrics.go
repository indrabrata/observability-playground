package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/indrabrata/observability-playground/utility"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalRequest *prometheus.CounterVec
	Latency      *prometheus.HistogramVec
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := utility.NewInterceptor(w)
		next.ServeHTTP(crw, r)

		elapsed := time.Since(start).Milliseconds()
		TotalRequest.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(crw.StatusCode)).Inc()
		Latency.WithLabelValues(r.Method, r.URL.Path).Observe(float64(elapsed))
	})
}

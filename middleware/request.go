package middleware

import (
	"net/http"

	"github.com/indrabrata/observability-playground/common"
	"go.uber.org/zap"
)

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("Request received", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("Host", r.Host), zap.String("User-Agent", r.UserAgent()), zap.String("IP-Address", r.RemoteAddr), zap.String("requestId", r.Context().Value("requestId").(string)))

		next.ServeHTTP(w, r)

		crw := w.(*common.Interceptor)
		zap.L().Info("Request completed", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("Host", r.Host), zap.String("User-Agent", r.UserAgent()), zap.String("IP-Address", r.RemoteAddr), zap.Int("status", crw.StatusCode), zap.String("requestId", r.Context().Value("requestId").(string)))
	})
}

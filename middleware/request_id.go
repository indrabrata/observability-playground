package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := w.Header().Get("X-Request-Id")
		if requestId == "" {
			requestId = uuid.NewString()
		}
		w.Header().Set("X-Request-Id", requestId)
		r = r.WithContext(context.WithValue(r.Context(), "requestId", requestId))
		next.ServeHTTP(w, r)
	})
}

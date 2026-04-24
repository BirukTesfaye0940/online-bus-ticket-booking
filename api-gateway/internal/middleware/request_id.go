package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

// RequestID injects a unique request ID into each request's context and response headers.
// This ID is propagated to downstream gRPC calls for distributed tracing.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		w.Header().Set(RequestIDKey, requestID)
		r.Header.Set(RequestIDKey, requestID)
		next.ServeHTTP(w, r)
	})
}

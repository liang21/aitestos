// Package middleware provides HTTP middleware implementations
package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Logging creates a logging middleware
func Logging(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate or get trace ID
			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = uuid.New().String()
			}

			// Set trace ID in response header
			w.Header().Set("X-Trace-ID", traceID)

			// Create wrapped response writer to capture status
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Log request
			logger.Info().
				Str("trace_id", traceID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Msg("request started")

			// Call next handler
			next.ServeHTTP(ww, r)

			// Log response
			logger.Info().
				Str("trace_id", traceID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Dur("duration", time.Since(start)).
				Msg("request completed")
		})
	}
}

// Package middleware provides HTTP middleware implementations
package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

// Recovery creates a panic recovery middleware
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					// Get trace ID from header if exists
					traceID := r.Header.Get("X-Trace-ID")

					// Log the panic using zerolog
					stack := debug.Stack()
					log.Error().
						Str("trace_id", traceID).
						Str("service_name", "aitestos").
						Interface("panic", rvr).
						Str("stack", string(stack)).
						Msg("panic recovered")

					// Send error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					response := map[string]string{"error": "internal server error"}
					_ = json.NewEncoder(w).Encode(response)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

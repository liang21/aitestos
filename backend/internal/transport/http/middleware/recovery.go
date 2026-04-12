// Package middleware provides HTTP middleware implementations
package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/liang21/aitestos/internal/ierrors"
)

// Recovery creates a panic recovery middleware
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					// Get trace ID from header if exists
					traceID := r.Header.Get("X-Trace-ID")

					// Log the panic
					stack := debug.Stack()
					fmt.Printf("PANIC: %v\n%s\n", rvr, string(stack))

					// Send error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					response := map[string]interface{}{
						"code":    ierrors.CodeInternalError,
						"message": "Internal server error",
					}
					if traceID != "" {
						response["traceId"] = traceID
					}

					_ = json.NewEncoder(w).Encode(response)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

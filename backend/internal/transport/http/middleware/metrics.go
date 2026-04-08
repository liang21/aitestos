// Package middleware provides HTTP middleware implementations
package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds Prometheus metrics for HTTP requests
type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

// NewMetrics creates a new Metrics instance with the given service name
func NewMetrics(serviceName string) *Metrics {
	return &Metrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_requests_total",
				Help:        "Total number of HTTP requests",
				ConstLabels: prometheus.Labels{"service": serviceName},
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_request_duration_seconds",
				Help:        "Duration of HTTP requests in seconds",
				ConstLabels: prometheus.Labels{"service": serviceName},
				Buckets:     prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}

// MetricsMiddleware creates a metrics middleware
func MetricsMiddleware(metrics *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create wrapped response writer to capture status
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Call next handler
			next.ServeHTTP(ww, r)

			// Record metrics
			duration := time.Since(start).Seconds()
			path := r.URL.Path
			method := r.Method
			status := strconv.Itoa(ww.Status())

			metrics.RequestsTotal.WithLabelValues(method, path, status).Inc()
			metrics.RequestDuration.WithLabelValues(method, path).Observe(duration)
		})
	}
}

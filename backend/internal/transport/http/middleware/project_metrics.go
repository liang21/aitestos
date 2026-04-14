// Package middleware provides HTTP middleware implementations
package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

// ProjectMetrics holds Prometheus metrics for project operations
type ProjectMetrics struct {
	QueryDuration *prometheus.HistogramVec
}

// NewProjectMetrics creates a new ProjectMetrics instance
func NewProjectMetrics(serviceName string) *ProjectMetrics {
	return &ProjectMetrics{
		QueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "project_statistics_query_duration_seconds",
				Help:        "Duration of project statistics queries in seconds",
				ConstLabels: prometheus.Labels{"service": serviceName},
				Buckets:     []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"status"}, // success, error
		),
	}
}

// Package project provides project statistics domain models
package project

import "time"

// ProjectStatistics represents aggregated project statistical data
type ProjectStatistics struct {
	ModuleCount      int64         `json:"module_count"`
	CaseCount        int64         `json:"case_count"`
	DocumentCount    int64         `json:"document_count"`
	PassRate         float64       `json:"pass_rate"`
	CoverageRate     float64       `json:"coverage_rate"`
	AIGeneratedCount int64         `json:"ai_generated_count"`
	RecentTasks      []TaskSummary `json:"recent_tasks,omitempty"`
	PassRateTrend    []TrendData   `json:"pass_rate_trend,omitempty"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// TaskSummary represents a summary of a recent generation task
type TaskSummary struct {
	ID            string             `json:"id"`
	Status        string             `json:"status"`
	ResultSummary *TaskResultSummary `json:"result_summary,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
}

// TaskResultSummary summarizes the results of a generation task
type TaskResultSummary struct {
	TotalDrafts    int64 `json:"total_drafts"`
	ConfirmedCount int64 `json:"confirmed_count"`
	RejectedCount  int64 `json:"rejected_count"`
}

// TrendData represents a data point in a trend series
type TrendData struct {
	Date string  `json:"date"`
	Rate float64 `json:"rate"`
}

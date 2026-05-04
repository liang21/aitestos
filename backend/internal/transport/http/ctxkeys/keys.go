// Package ctxkeys provides centralized context key definitions for HTTP layer
package ctxkeys

// ContextKey is the type for context keys to prevent collisions
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// ProjectIDKey is the context key for project ID
	ProjectIDKey ContextKey = "project_id"
	// ModuleIDKey is the context key for module ID
	ModuleIDKey ContextKey = "module_id"
	// CaseIDKey is the context key for test case ID
	CaseIDKey ContextKey = "case_id"
	// PlanIDKey is the context key for test plan ID
	PlanIDKey ContextKey = "plan_id"
	// TaskIDKey is the context key for generation task ID
	TaskIDKey ContextKey = "task_id"
	// DraftIDKey is the context key for generation draft ID
	DraftIDKey ContextKey = "draft_id"
	// DocumentIDKey is the context key for document ID
	DocumentIDKey ContextKey = "document_id"
)

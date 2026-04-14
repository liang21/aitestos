// Package ierrors defines HTTP error response structures
package ierrors

import "encoding/json"

// ErrorResponse is the unified API error response (for reference only)
// Note: Currently we use the simple format {"error": "message"} across all handlers
// This struct is kept for potential future use if we switch to structured error responses
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId,omitempty"`
}

// ToJSON converts ErrorResponse to JSON bytes
func (e *ErrorResponse) ToJSON() []byte {
	data, _ := json.Marshal(e)
	return data
}

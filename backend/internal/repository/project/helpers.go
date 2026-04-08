// Package project provides repository helper functions
package project

import (
	"encoding/json"
	"time"
)

// parseTime parses a time string into time.Time
func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// toJSON serializes a value to JSON string
func toJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// fromJSON deserializes a JSON string into a value
func fromJSON(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

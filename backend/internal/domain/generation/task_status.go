// Package generation defines TaskStatus value object
package generation

import "errors"

// TaskStatus is a value object representing the status of a generation task
type TaskStatus string

const (
	// TaskPending means the task is pending
	TaskPending TaskStatus = "pending"
	// TaskProcessing means the task is being processed
	TaskProcessing TaskStatus = "processing"
	// TaskCompleted means the task has completed
	TaskCompleted TaskStatus = "completed"
	// TaskFailed means the task has failed
	TaskFailed TaskStatus = "failed"
)

// ParseTaskStatus parses a string into TaskStatus
func ParseTaskStatus(s string) (TaskStatus, error) {
	switch s {
	case string(TaskPending):
		return TaskPending, nil
	case string(TaskProcessing):
		return TaskProcessing, nil
	case string(TaskCompleted):
		return TaskCompleted, nil
	case string(TaskFailed):
		return TaskFailed, nil
	default:
		return "", errors.New("invalid task status")
	}
}

// String returns the string representation of TaskStatus
func (s TaskStatus) String() string {
	return string(s)
}

// CanTransitionTo checks if the status can transition to the target status
func (s TaskStatus) CanTransitionTo(target TaskStatus) bool {
	transitions := map[TaskStatus][]TaskStatus{
		TaskPending:    {TaskProcessing},
		TaskProcessing: {TaskCompleted, TaskFailed},
		TaskCompleted:  {},
		TaskFailed:     {},
	}

	allowed, exists := transitions[s]
	if !exists {
		return false
	}

	for _, t := range allowed {
		if t == target {
			return true
		}
	}
	return false
}

// IsFinal returns true if the status is a final state
func (s TaskStatus) IsFinal() bool {
	return s == TaskCompleted || s == TaskFailed
}

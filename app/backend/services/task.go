package services

import (
	"strings"
	"unicode/utf8"

	"github.com/oklog/ulid/v2"
)

// INFO: タスクを表す構造体
type Task struct {
	ID          string
	Title       string
	Description *string
	Priority    int
	Status      string
}

// FUNCTION: 新しいタスクを作成する関数
func NewTask(id string, title string, description string, priority int, status string) (*Task, error) {

	if id != "" {
		if _, err := ulid.ParseStrict(id); err != nil {
			return nil, &ModelValidationError{
				Detail: "id is not a valid ULID",
			}
		}
	}

	if strings.TrimSpace(title) == "" {
		return nil, &ModelValidationError{
			Detail: "title is required",
		}
	}

	if utf8.RuneCountInString(title) > 50 {
		return nil, &ModelValidationError{
			Detail: "title exceeds max length 50",
		}
	}

	var normalizedDescription *string

	if strings.TrimSpace(description) != "" {
		normalizedDescription = &description
	}

	if priority < 1 || priority > 5 {
		return nil, &ModelValidationError{
			Detail: "priority must be between 1 and 5",
		}
	}

	if status != "TODO" && status != "DOING" && status != "DONE" {
		return nil, &ModelValidationError{
			Detail: "status must be TODO, DOING, or DONE",
		}
	}

	if priority == 5 && normalizedDescription == nil {
		return nil, &ModelValidationError{
			Detail: "description is required when priority is 5",
		}
	}

	return &Task{
		ID:          id,
		Title:       title,
		Description: normalizedDescription,
		Priority:    priority,
		Status:      status,
	}, nil
}

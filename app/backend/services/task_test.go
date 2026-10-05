package services

import (
	"ariadne/backend/tests/testutil"
	"errors"
	"strings"
	"testing"
)

// TEST:
func Test_NewTask_IDValidation(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "empty id is valid",
			id:      "",
			wantErr: false,
		},
		{
			name:    "valid ulid is valid",
			id:      "01M4655XT94C1RCKEN92HWW6BM",
			wantErr: false,
		},
		{
			name:    "invalid ulid",
			id:      "invalid-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTask(
				tt.id,
				"Task A",
				"",
				3,
				"TODO",
			)

			if tt.wantErr {
				var validationErr *ModelValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("got %v, want ModelValidationError", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TEST:
func Test_NewTask_TitleValidation(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "empty title",
			title:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only title",
			title:   "   ",
			wantErr: true,
		},
		{
			name:    "50 characters",
			title:   strings.Repeat("あ", 50),
			wantErr: false,
		},
		{
			name:    "51 characters",
			title:   strings.Repeat("あ", 51),
			wantErr: true,
		},
		{
			name:    "title with surrounding spaces",
			title:   "  Task A  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTask(
				"",
				tt.title,
				"",
				3,
				"TODO",
			)

			if tt.wantErr {
				var validationErr *ModelValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("got %v, want ModelValidationError", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Title != tt.title {
				t.Errorf("Title = %q, want %q", got.Title, tt.title)
			}
		})
	}
}

// TEST:
func Test_NewTask_DescriptionNormalization(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	tests := []struct {
		name            string
		description     string
		wantNil         bool
		wantDescription string
	}{
		{
			name:        "empty description",
			description: "",
			wantNil:     true,
		},
		{
			name:        "whitespace only description",
			description: "   ",
			wantNil:     true,
		},
		{
			name:            "description with value",
			description:     "Description A",
			wantNil:         false,
			wantDescription: "Description A",
		},
		{
			name:            "description with surrounding spaces",
			description:     "  Description A  ",
			wantNil:         false,
			wantDescription: "  Description A  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTask(
				"",
				"Task A",
				tt.description,
				3,
				"TODO",
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantNil {
				if got.Description != nil {
					t.Errorf("Description = %q, want nil", *got.Description)
				}
				return
			}

			if got.Description == nil {
				t.Fatalf("Description = nil, want %q", tt.wantDescription)
			}

			if *got.Description != tt.wantDescription {
				t.Errorf(
					"Description = %q, want %q",
					*got.Description,
					tt.wantDescription,
				)
			}
		})
	}
}

// TEST:
func Test_NewTask_PriorityValidation(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	tests := []struct {
		name     string
		priority int
		wantErr  bool
	}{
		{
			name:     "below minimum",
			priority: 0,
			wantErr:  true,
		},
		{
			name:     "minimum",
			priority: 1,
			wantErr:  false,
		},
		{
			name:     "maximum",
			priority: 5,
			wantErr:  false,
		},
		{
			name:     "above maximum",
			priority: 6,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTask(
				"",
				"Task A",
				"Description A",
				tt.priority,
				"TODO",
			)

			if tt.wantErr {
				var validationErr *ModelValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("got %v, want ModelValidationError", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TEST:
func Test_NewTask_StatusValidation(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{
			name:    "TODO",
			status:  "TODO",
			wantErr: false,
		},
		{
			name:    "DOING",
			status:  "DOING",
			wantErr: false,
		},
		{
			name:    "DONE",
			status:  "DONE",
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  "INVALID",
			wantErr: true,
		},
		{
			name:    "",
			status:  "EMPTY",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTask(
				"",
				"Task A",
				"Description A",
				3,
				tt.status,
			)

			if tt.wantErr {
				var validationErr *ModelValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("got %v, want ModelValidationError", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TEST:
func Test_NewTask_CorrelationValidation(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	tests := []struct {
		name        string
		description string
		priority    int
		wantErr     bool
	}{
		{
			name:        "priority 5 with description",
			description: "Description A",
			priority:    5,
			wantErr:     false,
		},
		{
			name:        "priority 5 without description",
			description: "",
			priority:    5,
			wantErr:     true,
		},
		{
			name:        "priority 4 without description",
			description: "",
			priority:    4,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTask(
				"",
				"Task A",
				tt.description,
				tt.priority,
				"TODO",
			)

			if tt.wantErr {
				var validationErr *ModelValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("got %v, want ModelValidationError", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

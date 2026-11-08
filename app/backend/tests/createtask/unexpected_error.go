package createtask

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/services/public"
	"ariadne/backend/tests/testutil"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TEST:
func UnexpectedError(t *testing.T) {
	// INFO: GIVEN
	yamlData := []byte(`
format: ariadne-tasks
version: "2.1"

tasks:
  - id: "01M465B5JK59B98Q8AXTJ6HSDR"
    title: "Existing Task"
    description: ""
    priority: 2
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"
`)

	dir := t.TempDir()
	testutil.WriteTaskFile(t, dir, yamlData)

	var repository services.TaskRepository = file.NewTaskRepository(dir)
	service := services.NewTaskService(repository)

	request := public.CreateTaskRequest{
		Title:       "",
		Description: "New Description",
		Priority:    4,
		Status:      "TODO",
	}

	// INFO: ACTION
	_, err := service.CreateTask(request)

	// INFO: VERIFY
	if !errors.Is(err, services.ErrUnexpectedError) {
		t.Errorf("CreateTask() error = %v, want ErrUnexpectedError", err)
	}

	// Business ModelのValidation Errorが保持されていること
	var validationErr *services.ModelValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("CreateTask() error = %v, want *ModelValidationError in error chain", err)
	}

	saved, err := os.ReadFile(filepath.Join(dir, "tasks.yaml"))
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	if !bytes.Equal(saved, yamlData) {
		t.Errorf("tasks.yaml was modified unexpectedly")
	}
}

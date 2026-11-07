package modifytask

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
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
    title: "Task A"
    description: "Existing A"
    priority: 2
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"

  - id: "01M465BRVCYNR0H09A30JE03PW"
    title: "Task B"
    description: "Existing B"
    priority: 3
    status: DOING
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"
`)

	dir := t.TempDir()
	testutil.WriteTaskFile(t, dir, yamlData)

	var repository services.TaskRepository = file.NewTaskRepository(dir)
	service := services.NewTaskService(repository)

	title, err := services.NewTitleChange("")
	if err != nil {
		t.Fatalf("NewTitleChange() returned error: %v", err)
	}

	changes := []services.TaskChange{title}

	// INFO: ACTION
	err = service.ModifyTask("01M465BRVCYNR0H09A30JE03PW", changes)

	// INFO: VERIFY
	if !errors.Is(err, services.ErrUnexpectedError) {
		t.Errorf("ModifyTask() error = %v, want ErrUnexpectedError", err)
	}

	// Business ModelのValidation Errorが保持されていること
	var validationErr *services.ModelValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("ModifyTask() error = %v, want *ModelValidationError in error chain", err)
	}

	saved, err := os.ReadFile(filepath.Join(dir, "tasks.yaml"))
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	if !bytes.Equal(saved, yamlData) {
		t.Errorf("tasks.yaml was modified unexpectedly")
	}
}

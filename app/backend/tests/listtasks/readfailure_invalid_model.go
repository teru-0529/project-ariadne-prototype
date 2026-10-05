package listtasks

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/tests/testutil"
	"errors"
	"testing"
)

// TEST:
func ReadFailureInvalidModel(t *testing.T) {
	// INFO: GIVEN
	yaml := []byte(`
format: ariadne-tasks
version: "2.1"

tasks:
  - id: "＜ValidなULID＞"
    title: "Task A"
    description: "Description A"
    priority: 6
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"
`)

	dir := t.TempDir()
	testutil.WriteTaskFile(t, dir, yaml)

	repo := file.NewTaskRepository(dir)
	service := services.NewTaskService(repo)

	// INFO: ACTION
	got, err := service.ListTasks()

	// INFO: VERIFY
	if !errors.Is(err, services.ErrReadFailure) {
		t.Errorf("got %v, want ErrReadFailure", err)
	}

	var validationErr *services.ModelValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("got %v, want ModelValidationError", err)
	}

	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

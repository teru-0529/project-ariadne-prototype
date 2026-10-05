package listtasks

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/tests/testutil"
	"errors"
	"testing"
)

// TEST:
func ReadFailureFormatHeader(t *testing.T) {
	// INFO: GIVEN
	yaml := []byte(`
format: ariadne-tasks
tasks: []
`)

	dir := t.TempDir()
	testutil.WriteTaskFile(t, dir, yaml)

	var repository services.TaskRepository = file.NewTaskRepository(dir)
	service := services.NewTaskService(repository)

	// INFO: ACTION
	got, err := service.ListTasks()

	// INFO: VERIFY
	if !errors.Is(err, services.ErrReadFailure) {
		t.Fatalf("ListTasks() error = %v, want ErrReadFailure", err)
	}

	// 元のFormat Errorが失われていないこと
	var formatErr *services.FormatError
	if !errors.As(err, &formatErr) {
		t.Fatalf("ListTasks() error = %v, want *services.FormatError in error chain", err)
	}

	if got != nil {
		t.Errorf("ListTasks() = %v, want nil", got)
	}
}

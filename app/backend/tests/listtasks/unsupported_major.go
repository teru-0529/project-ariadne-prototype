package listtasks

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/tests/testutil"
	"errors"
	"testing"
)

// TEST:
func UnsupportedMajor(t *testing.T) {
	// INFO: GIVEN
	yaml := []byte(`
format: ariadne-tasks
version: "3.0"
tasks:
  current-format-does-not-understand-this:
    whatever: true
`)

	dir := t.TempDir()
	testutil.WriteTaskFile(t, dir, yaml)

	var repository services.TaskRepository = file.NewTaskRepository(dir)
	service := services.NewTaskService(repository)

	// INFO: ACTION
	got, err := service.ListTasks()

	// INFO: VERIFY
	var versionErr *services.UnsupportedMajorVersionError
	if !errors.As(err, &versionErr) {
		t.Fatalf("ListTasks() error = %v, want *services.UnsupportedMajorVersionError", err)
	}

	if versionErr.Major != 3 {
		t.Errorf("Major = %d, want 3", versionErr.Major)
	}

	if errors.Is(err, services.ErrReadFailure) {
		t.Errorf("ListTasks() error = %v, must not be ErrReadFailure", err)
	}

	if got != nil {
		t.Errorf("ListTasks() = %v, want nil", got)
	}
}

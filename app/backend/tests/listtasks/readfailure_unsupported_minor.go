package listtasks

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/tests/testutil"
	"errors"
	"testing"
)

// TEST:
func ReadFailureUnsupportedMinor(t *testing.T) {
	// INFO: GIVEN
	yaml := []byte(`
format: ariadne-tasks
version: "2.2"
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
	if !errors.Is(err, services.ErrReadFailure) {
		t.Fatalf("ListTasks() error = %v, want ErrReadFailure", err)
	}

	var versionErr *services.UnsupportedMinorVersionError
	if !errors.As(err, &versionErr) {
		t.Fatalf("ListTasks() error = %v, want *services.UnsupportedMinorVersionError", err)
	}

	if versionErr.Major != 2 || versionErr.Minor != 2 {
		t.Errorf("Major = %d, want 2; Minor = %d, want 2", versionErr.Major, versionErr.Minor)
	}

	if got != nil {
		t.Errorf("ListTasks() = %v, want nil", got)
	}
}

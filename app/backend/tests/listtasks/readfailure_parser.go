package listtasks

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/tests/testutil"
	"errors"
	"testing"
)

// TEST:
func ReadFailureParser(t *testing.T) {
	// INFO: GIVEN
	yaml := []byte(`
format: ariadne-tasks
version: "2.1"
tasks:
  - title: "unterminated
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

	// WARN: 元のYAMLパースエラーは固定での判定が不可能の為チェック対象外

	if got != nil {
		t.Errorf("ListTasks() = %v, want nil", got)
	}
}

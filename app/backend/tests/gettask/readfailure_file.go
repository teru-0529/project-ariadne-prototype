package gettask

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"errors"
	"os"
	"testing"
)

// TEST:
func ReadFailure(t *testing.T) {
	// INFO: GIVEN
	dir := t.TempDir()
	// tasks.yaml は作成しない

	var repository services.TaskRepository = file.NewTaskRepository(dir)
	service := services.NewTaskService(repository)

	// INFO: ACTION
	got, err := service.GetTask("01M465BRVCYNR0H09A30JE03PW")

	// INFO: VERIFY
	if !errors.Is(err, services.ErrReadFailure) {
		t.Fatalf("GetTask() error = %v, want ErrReadFailure", err)
	}

	// 元のFile I/O Errorが失われていないこと
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("ListTasks() error = %v, want *os.PathError in error chain", err)
	}

	if got != nil {
		t.Errorf("GetTask() = %+v, want nil", got)
	}
}

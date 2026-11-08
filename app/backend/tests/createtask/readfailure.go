package createtask

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/services/public"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TEST:
func ReadFailure(t *testing.T) {
	// INFO: GIVEN
	dir := t.TempDir()

	var repository services.TaskRepository = file.NewTaskRepository(dir)
	service := services.NewTaskService(repository)

	request := public.CreateTaskRequest{
		Title:       "New Task",
		Description: "New Description",
		Priority:    4,
		Status:      "TODO",
	}

	// INFO: ACTION
	_, err := service.CreateTask(request)

	// INFO: VERIFY
	if !errors.Is(err, services.ErrReadFailure) {
		t.Errorf("CreateTask() error = %v, want ErrReadFailure", err)
	}

	// 元のFile I/O Errorが失われていないこと
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Errorf("CreateTask() error = %v, want *os.PathError in error chain", err)
	}

	_, err = os.Stat(filepath.Join(dir, "tasks.yaml"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("tasks.yaml exists or Stat failed unexpectedly: %v", err)
	}
}

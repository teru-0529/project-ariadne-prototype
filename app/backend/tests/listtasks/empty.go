package listtasks

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/tests/testutil"
	"testing"
)

// TEST:
func Empty(t *testing.T) {
	// INFO: GIVEN
	yaml := []byte(`
format: ariadne-tasks
version: "2.1"

tasks: []
`)

	dir := t.TempDir()
	testutil.WriteTaskFile(t, dir, yaml)

	var repository services.TaskRepository = file.NewTaskRepository(dir)
	service := services.NewTaskService(repository)

	// INFO: ACTION
	got, err := service.ListTasks()

	// INFO: VERIFY
	if err != nil {
		t.Fatalf("ListTasks() returned error: %v", err)
	}

	if got == nil {
		t.Fatal("ListTasks() returned nil, want empty slice")
	}

	if len(got) != 0 {
		t.Errorf("ListTasks() returned %d tasks, want 0", len(got))
	}
}

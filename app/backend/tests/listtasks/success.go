package listtasks

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/tests/testutil"
	"testing"
)

// TEST:
func Success(t *testing.T) {
	// INFO: GIVEN
	yaml := []byte(`
format: ariadne-tasks
version: "2.1"

tasks:
  - id: "01J9C0000000000000000000000"
    title: "Task C"
    description: ""
    priority: 5
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"

  - id: "01J9C0000000000000000000001"
    title: "Task A"
    description: ""
    priority: 1
    status: DOING
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"

  - id: "01J9C0000000000000000000002"
    title: "Task B"
    description: ""
    priority: 3
    status: DONE
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"
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

	if len(got) != 3 {
		t.Fatalf("ListTasks() returned %d tasks, want 3", len(got))
	}

	want := []services.TaskSummary{
		{
			ID:       "01J9C0000000000000000000000",
			Title:    "Task C",
			Priority: 5,
			Status:   "TODO",
		},
		{
			ID:       "01J9C0000000000000000000001",
			Title:    "Task A",
			Priority: 1,
			Status:   "DOING",
		},
		{
			ID:       "01J9C0000000000000000000002",
			Title:    "Task B",
			Priority: 3,
			Status:   "DONE",
		},
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("TaskSummary[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

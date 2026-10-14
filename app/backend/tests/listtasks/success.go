package listtasks

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/services/publish"
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
  - id: "01M465B5JK59B98Q8AXTJ6HSDR"
    title: "Task C"
    description: |
      "重要度5の場合は説明必須"
    priority: 5
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"

  - id: "01M465BRVCYNR0H09A30JE03PW"
    title: "Task A"
    description: ""
    priority: 1
    status: DOING
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"

  - id: "01M465BRVCDDWKHCTF84YK2QQD"
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

	want := []publish.TaskSummary{
		{
			ID:       "01M465B5JK59B98Q8AXTJ6HSDR",
			Title:    "Task C",
			Priority: 5,
			Status:   "TODO",
		},
		{
			ID:       "01M465BRVCYNR0H09A30JE03PW",
			Title:    "Task A",
			Priority: 1,
			Status:   "DOING",
		},
		{
			ID:       "01M465BRVCDDWKHCTF84YK2QQD",
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

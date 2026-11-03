package createtask

import (
	"ariadne/backend/adapters/file"
	"ariadne/backend/services"
	"ariadne/backend/services/public"
	"ariadne/backend/tests/testutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/oklog/ulid/v2"
	"gopkg.in/yaml.v3"
)

// TEST:
func Success(t *testing.T) {
	// INFO: GIVEN
	yamlData := []byte(`
format: ariadne-tasks
version: "2.1"

tasks:
  - id: "01M465B5JK59B98Q8AXTJ6HSDR"
    title: "Task A"
    description: |
      Existing A
    priority: 2
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"

  - id: "01M465BRVCYNR0H09A30JE03PW"
    title: "Task B"
    description: |
      Existing B
    priority: 3
    status: DOING
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"
`)

	dir := t.TempDir()
	testutil.WriteTaskFile(t, dir, yamlData)

	var repository services.TaskRepository = file.NewTaskRepository(dir)
	service := services.NewTaskService(repository)

	request := public.CreateTaskRequest{
		Title:       "New Task",
		Description: "New Description",
		Priority:    4,
		Status:      "TODO",
	}

	// INFO: ACTION
	id, err := service.CreateTask(request)

	// INFO: VERIFY
	if err != nil {
		t.Fatalf("CreateTask() returned error: %v", err)
	}

	if _, err := ulid.ParseStrict(id); err != nil {
		t.Errorf("CreateTask() returned invalid ULID: %q", id)
	}

	saved, err := os.ReadFile(filepath.Join(dir, "tasks.yaml"))
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	var document struct {
		Tasks []struct {
			ID          string `yaml:"id"`
			Title       string `yaml:"title"`
			Description string `yaml:"description"`
			Priority    int    `yaml:"priority"`
			Status      string `yaml:"status"`
			CreatedAt   string `yaml:"createdAt"`
			UpdatedAt   string `yaml:"updatedAt"`
		} `yaml:"tasks"`
	}

	if err := yaml.Unmarshal(saved, &document); err != nil {
		t.Fatalf("YAML Unmarshal failed: %v", err)
	}

	if len(document.Tasks) != 3 {
		t.Fatalf("Task count = %d, want 3", len(document.Tasks))
	}

	// 既存Taskが維持されていること
	if document.Tasks[0].ID != "01M465B5JK59B98Q8AXTJ6HSDR" ||
		document.Tasks[0].Title != "Task A" ||
		document.Tasks[0].Description != "Existing A\n" ||
		document.Tasks[0].Priority != 2 ||
		document.Tasks[0].Status != "TODO" ||
		document.Tasks[0].CreatedAt != "2026-10-04T09:00:00+09:00" ||
		document.Tasks[0].UpdatedAt != "2026-10-04T09:00:00+09:00" {
		t.Errorf("Existing Task A was modified: %+v", document.Tasks[0])
	}

	if document.Tasks[1].ID != "01M465BRVCYNR0H09A30JE03PW" ||
		document.Tasks[1].Title != "Task B" ||
		document.Tasks[1].Description != "Existing B\n" ||
		document.Tasks[1].Priority != 3 ||
		document.Tasks[1].Status != "DOING" ||
		document.Tasks[1].CreatedAt != "2026-10-04T09:00:00+09:00" ||
		document.Tasks[1].UpdatedAt != "2026-10-04T09:00:00+09:00" {
		t.Errorf("Existing Task B was modified: %+v", document.Tasks[1])
	}

	// 新規Taskが末尾に追加されていること
	created := document.Tasks[2]

	if created.ID != id ||
		created.Title != request.Title ||
		created.Description != request.Description ||
		created.Priority != request.Priority ||
		created.Status != request.Status {
		t.Errorf("Created Task = %+v, request = %+v, ID = %q",
			created, request, id)
	}
}

package file

import (
	"ariadne/backend/tests/testutil"
	"testing"
)

// TEST:
func Test_ParseYAML_Success(t *testing.T) {

	if !*testutil.Logic {
		t.SkipNow()
	}

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
    description: |
      Description A
    priority: 1
    status: DOING
    createdAt: "2026-10-04T10:00:00+09:00"
    updatedAt: "2026-10-04T11:00:00+09:00"
`)

	// INFO: ACTION
	got, err := parseYAML(yaml)

	// INFO: VERIFY
	if err != nil {
		t.Fatalf("parseYAML() returned error: %v", err)
	}

	if got.Format != "ariadne-tasks" {
		t.Errorf("Format = %q, want %q", got.Format, "ariadne-tasks")
	}

	if got.Version != "2.1" {
		t.Errorf("Version = %q, want %q", got.Version, "2.1")
	}

	if len(got.Tasks) != 2 {
		t.Fatalf("Tasks = %d, want 2", len(got.Tasks))
	}

	if got.Tasks[0].Title != "Task C" {
		t.Errorf("Tasks[0].Title = %q, want %q", got.Tasks[0].Title, "Task C")
	}

	if got.Tasks[1].Title != "Task A" {
		t.Errorf("Tasks[1].Title = %q, want %q", got.Tasks[1].Title, "Task A")
	}
}

// TEST:
func Test_ReadYAML_Success(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	// INFO: GIVEN
	yaml := []byte(`
format: ariadne-tasks
version: "2.1"

tasks:
  - id: "01J9C0000000000000000000000"
    title: "Task A"
    description: ""
    priority: 1
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T09:00:00+09:00"
`)

	dir := t.TempDir()
	testutil.WriteTaskFile(t, dir, yaml)

	repo := NewTaskRepository(dir)

	// INFO: ACTION
	got, err := repo.readYAML()

	// INFO: VERIFY
	if err != nil {
		t.Fatalf("readYAML() returned error: %v", err)
	}

	if got.Format != "ariadne-tasks" {
		t.Errorf("Format = %q, want %q", got.Format, "ariadne-tasks")
	}

	if got.Version != "2.1" {
		t.Errorf("Version = %q, want %q", got.Version, "2.1")
	}

	if len(got.Tasks) != 1 {
		t.Fatalf("Tasks = %d, want 1", len(got.Tasks))
	}

	if got.Tasks[0].Title != "Task A" {
		t.Errorf("Tasks[0].Title = %q, want %q", got.Tasks[0].Title, "Task A")
	}
}

package file

import (
	"ariadne/backend/services"
	"ariadne/backend/tests/testutil"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

const storeTestYAML = `format: ariadne-tasks
version: "2.1"
tasks:
  - id: 01M46599N791CM7QP2YN9KK94K
    title: Existing Task
    description: Existing Description
    priority: 2
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T11:00:00+09:00"
`

// TEST: DepositでTaskを末尾に追加できること
func Test_TaskRepository_Store_Deposit_Success(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	// INFO: GIVEN
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.yaml")

	if err := os.WriteFile(path, []byte(storeTestYAML), 0600); err != nil {
		t.Fatal(err)
	}

	repo := NewTaskRepository(dir)

	task, err := services.NewTask(
		"01M46599N791CM7QP2YN9KK95K",
		"New Task",
		"New Description",
		3,
		"DOING",
	)
	if err != nil {
		t.Fatal(err)
	}

	// INFO: ACTION
	err = repo.Store(task, services.StoreDeposit)

	// INFO: VERIFY
	if err != nil {
		t.Fatalf("Store() returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var dataset taskDataset
	if err := yaml.Unmarshal(data, &dataset); err != nil {
		t.Fatalf("YAML parse failed: %v", err)
	}

	if dataset.Format != formatName {
		t.Errorf("Format = %q, want %q", dataset.Format, formatName)
	}
	if dataset.Version != "2.1" {
		t.Errorf("Version = %q, want %q", dataset.Version, "2.1")
	}

	if len(dataset.Tasks) != 2 {
		t.Fatalf("Task count = %d, want 2", len(dataset.Tasks))
	}

	existing := dataset.Tasks[0]
	if existing.ID == nil || *existing.ID != "01M46599N791CM7QP2YN9KK94K" {
		t.Errorf("Existing Task ID was changed")
	}
	if existing.Title == nil || *existing.Title != "Existing Task" {
		t.Errorf("Existing Task Title was changed")
	}
	if existing.Description == nil || *existing.Description != "Existing Description" {
		t.Errorf("Existing Task Description was changed")
	}
	if existing.Priority == nil || *existing.Priority != 2 {
		t.Errorf("Existing Task Priority was changed")
	}
	if existing.Status == nil || *existing.Status != "TODO" {
		t.Errorf("Existing Task Status was changed")
	}
	if existing.CreatedAt == nil || *existing.CreatedAt != "2026-10-04T09:00:00+09:00" {
		t.Errorf("Existing Task CreatedAt was changed")
	}
	if existing.UpdatedAt == nil || *existing.UpdatedAt != "2026-10-04T11:00:00+09:00" {
		t.Errorf("Existing Task UpdatedAt was changed")
	}

	added := dataset.Tasks[1]

	if added.ID == nil || *added.ID != task.ID {
		t.Errorf("ID mismatch")
	}
	if added.Title == nil || *added.Title != task.Title {
		t.Errorf("Title mismatch")
	}
	if added.Description == nil || *added.Description != *task.Description {
		t.Errorf("Description mismatch")
	}
	if added.Priority == nil || *added.Priority != task.Priority {
		t.Errorf("Priority mismatch")
	}
	if added.Status == nil || *added.Status != task.Status {
		t.Errorf("Status mismatch")
	}

	if added.CreatedAt == nil || added.UpdatedAt == nil {
		t.Fatal("Timestamp is missing")
	}

	if _, err := time.Parse(time.RFC3339, *added.CreatedAt); err != nil {
		t.Errorf("CreatedAt is not RFC3339: %v", err)
	}
	if _, err := time.Parse(time.RFC3339, *added.UpdatedAt); err != nil {
		t.Errorf("UpdatedAt is not RFC3339: %v", err)
	}
	if *added.CreatedAt != *added.UpdatedAt {
		t.Errorf("CreatedAt and UpdatedAt differ")
	}
}

// TEST: DepositでIDが重複した場合にエラーとなること
func Test_TaskRepository_Store_Deposit_KeyDuplicate(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	// INFO: GIVEN
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.yaml")

	before := []byte(storeTestYAML)
	if err := os.WriteFile(path, before, 0600); err != nil {
		t.Fatal(err)
	}

	repo := NewTaskRepository(dir)

	task, err := services.NewTask(
		"01M46599N791CM7QP2YN9KK94K",
		"Duplicate Task",
		"Duplicate Description",
		3,
		"DOING",
	)
	if err != nil {
		t.Fatal(err)
	}

	// INFO: ACTION
	err = repo.Store(task, services.StoreDeposit)

	// INFO: VERIFY
	var duplicateErr *services.KeyDuplicateError
	if !errors.As(err, &duplicateErr) {
		t.Errorf("Store() error = %v, want KeyDuplicateError", err)
	} else if duplicateErr.Key != task.ID {
		t.Errorf("Duplicate key = %q, want %q", duplicateErr.Key, task.ID)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(before, after) {
		t.Errorf("YAML was modified on Key Duplicate")
	}
}

// TEST: Description未設定時に空文字列へ変換されること
func Test_NewPersistenceTask_EmptyDescription(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	// INFO: GIVEN
	task, err := services.NewTask(
		"01M4655XT94C1RCKEN92HWW6BM",
		"Task A",
		"",
		3,
		"TODO",
	)
	if err != nil {
		t.Fatal(err)
	}

	createdAt := "2026-11-05T09:00:00+09:00"
	updatedAt := "2026-11-05T09:00:00+09:00"

	// INFO: ACTION
	got := newPersistenceTask(task, createdAt, updatedAt)

	// INFO: VERIFY
	if task.Description != nil {
		t.Errorf("Business Task Description = %v, want nil", task.Description)
	}

	if got.Description == nil {
		t.Fatal("Persistence Task Description = nil, want empty string")
	}

	if *got.Description != "" {
		t.Errorf("Persistence Task Description = %q, want empty string", *got.Description)
	}
}

// TEST: Replaceで指定Taskを丸ごと置換できること
func Test_TaskRepository_Store_Replace_Success(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	// INFO: GIVEN
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.yaml")

	before := `format: ariadne-tasks
version: "2.1"
tasks:
  - id: 01M46599N791CM7QP2YN9KK94K
    title: Task A
    description: Description A
    priority: 2
    status: TODO
    createdAt: "2026-10-04T09:00:00+09:00"
    updatedAt: "2026-10-04T11:00:00+09:00"
  - id: 01M46599N791CM7QP2YN9KK95K
    title: Task B
    description: Description B
    priority: 3
    status: DOING
    createdAt: "2026-10-05T09:00:00+09:00"
    updatedAt: "2026-10-05T11:00:00+09:00"
`
	if err := os.WriteFile(path, []byte(before), 0600); err != nil {
		t.Fatal(err)
	}

	repo := NewTaskRepository(dir)

	task, err := services.NewTask(
		"01M46599N791CM7QP2YN9KK95K",
		"Updated Task B",
		"Updated Description B",
		5,
		"DONE",
	)
	if err != nil {
		t.Fatal(err)
	}

	// INFO: ACTION
	err = repo.Store(task, services.StoreReplace)

	// INFO: VERIFY
	if err != nil {
		t.Fatalf("Store() returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var dataset taskDataset
	if err := yaml.Unmarshal(data, &dataset); err != nil {
		t.Fatalf("YAML parse failed: %v", err)
	}

	if dataset.Format != formatName || dataset.Version != "2.1" {
		t.Errorf("YAML header was changed")
	}

	if len(dataset.Tasks) != 2 {
		t.Fatalf("Task count = %d, want 2", len(dataset.Tasks))
	}

	// 置換対象外のTaskは変更されない
	original := dataset.Tasks[0]
	if original.ID == nil || *original.ID != "01M46599N791CM7QP2YN9KK94K" ||
		original.Title == nil || *original.Title != "Task A" ||
		original.Description == nil || *original.Description != "Description A" ||
		original.Priority == nil || *original.Priority != 2 ||
		original.Status == nil || *original.Status != "TODO" ||
		original.CreatedAt == nil || *original.CreatedAt != "2026-10-04T09:00:00+09:00" ||
		original.UpdatedAt == nil || *original.UpdatedAt != "2026-10-04T11:00:00+09:00" {
		t.Errorf("Unrelated Task was changed")
	}

	// 2件目が同じ位置で置換される
	replaced := dataset.Tasks[1]

	if replaced.ID == nil || *replaced.ID != task.ID {
		t.Errorf("ID mismatch")
	}
	if replaced.Title == nil || *replaced.Title != task.Title {
		t.Errorf("Title mismatch")
	}
	if replaced.Description == nil || *replaced.Description != *task.Description {
		t.Errorf("Description mismatch")
	}
	if replaced.Priority == nil || *replaced.Priority != task.Priority {
		t.Errorf("Priority mismatch")
	}
	if replaced.Status == nil || *replaced.Status != task.Status {
		t.Errorf("Status mismatch")
	}

	if replaced.CreatedAt == nil || replaced.UpdatedAt == nil {
		t.Fatal("Timestamp is missing")
	}

	if *replaced.CreatedAt != "2026-10-05T09:00:00+09:00" {
		t.Errorf("CreatedAt was changed: %q", *replaced.CreatedAt)
	}

	if _, err := time.Parse(time.RFC3339, *replaced.UpdatedAt); err != nil {
		t.Errorf("UpdatedAt is not RFC3339: %v", err)
	}

	if *replaced.UpdatedAt == "2026-10-05T11:00:00+09:00" {
		t.Errorf("UpdatedAt was not updated")
	}
}

// TEST: Replaceで対象Taskが存在しない場合にエラーとなること
func Test_TaskRepository_Store_Replace_NotFound(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	// INFO: GIVEN
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.yaml")

	before := []byte(storeTestYAML)
	if err := os.WriteFile(path, before, 0600); err != nil {
		t.Fatal(err)
	}

	repo := NewTaskRepository(dir)

	task, err := services.NewTask(
		"01M46599N791CM7QP2YN9KK95K",
		"Unknown Task",
		"Unknown Description",
		3,
		"DOING",
	)
	if err != nil {
		t.Fatal(err)
	}

	// INFO: ACTION
	err = repo.Store(task, services.StoreReplace)

	// INFO: VERIFY
	var notFoundErr *services.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("Store() error = %v, want NotFoundError", err)
	} else if notFoundErr.Key != task.ID {
		t.Errorf("NotFound key = %q, want %q", notFoundErr.Key, task.ID)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(before, after) {
		t.Errorf("YAML was modified on Not Found")
	}
}

package file

import (
	"ariadne/backend/tests/testutil"
	"testing"
)

// TEST:
func Test_PersistenceTask_To_Task_Success(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	// INFO: GIVEN
	source := persistenceTask{
		ID:          "01J9C0000000000000000000000",
		Title:       "Task A",
		Description: "Description A",
		Priority:    3,
		Status:      "DOING",
		CreatedAt:   "2026-10-04T09:00:00+09:00",
		UpdatedAt:   "2026-10-04T11:00:00+09:00",
	}

	// INFO: ACTION
	got, err := source.toTask()

	// INFO: VERIFY
	if err != nil {
		t.Fatalf("toTask() returned error: %v", err)
	}

	if got.ID != source.ID {
		t.Errorf("ID = %q, want %q", got.ID, source.ID)
	}

	if got.Title != source.Title {
		t.Errorf("Title = %q, want %q", got.Title, source.Title)
	}

	if got.Description != source.Description {
		t.Errorf("Description = %q, want %q", got.Description, source.Description)
	}

	if got.Priority != source.Priority {
		t.Errorf("Priority = %d, want %d", got.Priority, source.Priority)
	}

	if got.Status != source.Status {
		t.Errorf("Status = %q, want %q", got.Status, source.Status)
	}
}

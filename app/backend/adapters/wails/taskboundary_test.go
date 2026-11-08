package wails

import (
	"errors"
	"testing"

	"ariadne/backend/services"
	"ariadne/backend/services/public"
	"ariadne/backend/tests/testutil"
)

// TEST:
func Test_ToPublicError_ReadFailure(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	err := toPublicError(services.ErrReadFailure)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "タスクを読み込めませんでした。" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TEST:
func Test_ToPublicError_UnsupportedMajorVersion(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	err := toPublicError(&services.UnsupportedMajorVersionError{
		ActualMajor:    3,
		SupportedMajor: 2,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	want := "データのバージョン 3.x には対応していません。対応バージョンは 2.x です。"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

// TEST:
func Test_ToPublicError_Unknown(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	err := toPublicError(errors.New("something happened"))

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	want := "予期しないエラーが発生しました。"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

// TEST:
func Test_ToPublicError_NotFound(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	err := toPublicError(&services.NotFoundError{Key: "task-001"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	want := "指定されたタスクが見つかりませんでした。"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

// TEST:
func Test_ToPublicError_KeyDuplicate(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	err := toPublicError(&services.KeyDuplicateError{Key: "task-001"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	want := "同じIDのタスクが既に存在します。"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

// TEST:
func Test_ToPublicError_WriteFailure(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	err := toPublicError(services.ErrWriteFailure)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	want := "タスクを保存できませんでした。"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

// TEST: TaskChangeの正常変換
func Test_ConvertTaskChanges_Success(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	request := public.EditTaskRequest{
		Changes: []public.TaskPropertyChange{
			{Type: "title", Value: "Updated Title"},
			{Type: "description", Value: "Updated Description"},
			{Type: "priority", Value: 2},
			{Type: "status", Value: "done"},
		},
	}

	changes, err := convertTaskChanges(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 4 {
		t.Fatalf("len(changes) = %d, want 4", len(changes))
	}

	if _, ok := changes[0].(*services.TitleChange); !ok {
		t.Errorf("changes[0] = %T, want *services.TitleChange", changes[0])
	}
	if _, ok := changes[1].(*services.DescriptionChange); !ok {
		t.Errorf("changes[1] = %T, want *services.DescriptionChange", changes[1])
	}
	if _, ok := changes[2].(*services.PriorityChange); !ok {
		t.Errorf("changes[2] = %T, want *services.PriorityChange", changes[2])
	}
	if _, ok := changes[3].(*services.StatusChange); !ok {
		t.Errorf("changes[3] = %T, want *services.StatusChange", changes[3])
	}
}

// TEST: 未知のProperty Type
func Test_ConvertTaskChanges_UnknownType(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	request := public.EditTaskRequest{
		Changes: []public.TaskPropertyChange{
			{Type: "title", Value: "Updated Title"},
			{Type: "unknown", Value: "value"},
		},
	}

	changes, err := convertTaskChanges(request)

	if !errors.Is(err, services.ErrUnexpectedError) {
		t.Fatalf("err = %v, want ErrUnexpectedError", err)
	}
	if changes != nil {
		t.Fatalf("changes = %v, want nil", changes)
	}
}

// TEST: Property Valueの型不一致
func Test_ConvertTaskChanges_InvalidValueType(t *testing.T) {
	if !*testutil.Logic {
		t.SkipNow()
	}

	request := public.EditTaskRequest{
		Changes: []public.TaskPropertyChange{
			{Type: "title", Value: "Updated Title"},
			{Type: "priority", Value: "high"},
		},
	}

	changes, err := convertTaskChanges(request)

	if !errors.Is(err, services.ErrUnexpectedError) {
		t.Fatalf("err = %v, want ErrUnexpectedError", err)
	}
	if changes != nil {
		t.Fatalf("changes = %v, want nil", changes)
	}
}

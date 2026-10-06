package wails

import (
	"errors"
	"testing"

	"ariadne/backend/services"
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

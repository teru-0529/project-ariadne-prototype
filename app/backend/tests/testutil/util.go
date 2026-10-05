package testutil

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// D-Test実施フラグ
var Logic = flag.Bool("logic", false, "run logic tests")

// FUNCTION: テスト用のtasks.yamlを作成する関数
func WriteTaskFile(t *testing.T, dir string, data []byte) {
	t.Helper()

	if err := os.WriteFile(
		filepath.Join(dir, "tasks.yaml"),
		data,
		0644,
	); err != nil {
		t.Fatalf("failed to write tasks.yaml: %v", err)
	}
}

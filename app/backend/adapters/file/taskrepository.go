package file

import (
	services "ariadne/backend/services"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v3"
)

// INFO: タスクリポジトリの構造体
type TaskRepository struct {
	dir string // tasks.yamlが存在するディレクトリのパス
}

// FUNCTION: TaskRepositoryのコンストラクタ
func NewTaskRepository(dir string) *TaskRepository {
	return &TaskRepository{
		dir: dir,
	}
}

// FUNCTION: ディレクトリからタスクデータを取得する関数
func (r *TaskRepository) Retrieve() ([]*services.Task, error) {
	dataset, err := r.readYAML()
	if err != nil {
		return nil, err
	}

	var tasks []*services.Task
	for i, pt := range dataset.Tasks {
		task, err := pt.toTask()
		if err != nil {
			return nil, fmt.Errorf("%w: tasks[%d]: %w", services.ErrReadFailure, i, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// FUNCTION: ディレクトリからYAML形式のタスクデータを読み込む関数
func (r *TaskRepository) readYAML() (*taskDataset, error) {
	data, err := os.ReadFile(filepath.Join(r.dir, "tasks.yaml"))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", services.ErrReadFailure, err)
	}

	major, minor, err := parseHeader(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", services.ErrReadFailure, err)
	}

	if major != supportedMajorVersion {
		return nil, &services.UnsupportedMajorVersionError{
			ActualMajor:    major,
			SupportedMajor: supportedMajorVersion,
		}
	}

	if minor > supportedMinorVersion {
		return nil, fmt.Errorf("%w: %w", services.ErrReadFailure, &services.UnsupportedMinorVersionError{
			Major: major,
			Minor: minor,
		})
	}

	dataset, err := parseYAML(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", services.ErrReadFailure, err)
	}

	return dataset, nil
}

var headerPattern = regexp.MustCompile(
	fmt.Sprintf(
		`(?m)^format:\s*%s\s*\r?\nversion:\s*"(\d+)\.(\d+)"\s*$`,
		formatName,
	),
)

// FUNCTION: YAML形式のタスクデータのヘッダをパースする関数
func parseHeader(data []byte) (int, int, error) {
	matches := headerPattern.FindSubmatch(data)
	if matches == nil {
		return 0, 0, &services.FormatError{
			Detail: "invalid header",
		}
	}

	major, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		return 0, 0, &services.FormatError{
			Detail: "invalid major version",
		}
	}

	minor, err := strconv.Atoi(string(matches[2]))
	if err != nil {
		return 0, 0, &services.FormatError{
			Detail: "invalid minor version",
		}
	}

	return major, minor, nil
}

// FUNCTION: YAML形式のタスクデータをパースする関数
func parseYAML(data []byte) (*taskDataset, error) {
	var dataset taskDataset

	if err := yaml.Unmarshal(data, &dataset); err != nil {
		return nil, err //TODO: エラー処理wrap候補
	}
	return &dataset, nil
}

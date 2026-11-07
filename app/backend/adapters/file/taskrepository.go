package file

import (
	services "ariadne/backend/services"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

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

// FUNCTION: IDを指定してTaskを取得する関数
func (r *TaskRepository) Pick(id string) (*services.Task, error) {
	tasks, err := r.Retrieve()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if task.ID() == id {
			return task, nil
		}
	}

	return nil, &services.NotFoundError{Key: id}
}

// FUNCTION: Taskを永続化する関数 TODO:（A15-05以降で本実装）
func (r *TaskRepository) Store(task *services.Task, intent services.StoreIntent) error {
	dataset, err := r.readYAML()
	if err != nil {
		return err
	}

	switch intent {
	case services.StoreDeposit:
		for _, existing := range dataset.Tasks {
			if existing.ID != nil && *existing.ID == task.ID() {
				return &services.KeyDuplicateError{Key: task.ID()}
			}
		}

		now := time.Now().Format(time.RFC3339)
		dataset.Tasks = append(
			dataset.Tasks,
			newPersistenceTask(task, now, now),
		)

	case services.StoreReplace:
		found := false

		for i, existing := range dataset.Tasks {
			if existing.ID != nil && *existing.ID == task.ID() {
				now := time.Now().Format(time.RFC3339)

				dataset.Tasks[i] = newPersistenceTask(
					task,
					*existing.CreatedAt,
					now,
				)

				found = true
				break
			}
		}

		if !found {
			return &services.NotFoundError{Key: task.ID()}
		}

	default:
		return fmt.Errorf("%w: unsupported store intent: %d", services.ErrUnexpectedError, intent)
	}

	return r.writeYAML(dataset)
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
		return nil, err
	}
	return &dataset, nil
}

// FUNCTION: Dataset全体をYAML形式で安全に保存する関数
func (r *TaskRepository) writeYAML(dataset *taskDataset) error {
	// YAML形式に変換
	data, err := yaml.Marshal(dataset)
	if err != nil {
		return fmt.Errorf("%w: %w", services.ErrWriteFailure, err)
	}

	// 正本と同じディレクトリに一時ファイルを作成
	temp, err := os.CreateTemp(r.dir, ".tasks-*.tmp")
	if err != nil {
		return fmt.Errorf("%w: %w", services.ErrWriteFailure, err)
	}

	// 処理失敗時の一時ファイルを後始末
	defer os.Remove(temp.Name())

	// 一時ファイルに全量書き込み
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return fmt.Errorf("%w: %w", services.ErrWriteFailure, err)
	}

	// 書き込み完了後にファイルを閉じる
	if err := temp.Close(); err != nil {
		return fmt.Errorf("%w: %w", services.ErrWriteFailure, err)
	}

	// 一時ファイルで正本を置換
	if err := os.Rename(
		temp.Name(),
		filepath.Join(r.dir, "tasks.yaml"),
	); err != nil {
		return fmt.Errorf("%w: %w", services.ErrWriteFailure, err)
	}

	return nil
}

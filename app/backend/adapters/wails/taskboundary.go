package wails

import (
	"errors"
	"fmt"

	"ariadne/backend/services"
	"ariadne/backend/services/publish"
)

// INFO: タスクバウンダリの構造体
type TaskBoundary struct {
	service *services.TaskService
	logger  Logger
}

// FUNCTION: TaskBoundaryのコンストラクタ関数
func NewTaskBoundary(
	service *services.TaskService,
	logger Logger,
) *TaskBoundary {
	return &TaskBoundary{
		service: service,
		logger:  logger,
	}
}

// FUNCTION: タスクの一覧を取得する関数
func (b *TaskBoundary) ListTasks() ([]publish.TaskSummary, error) {
	tasks, err := b.service.ListTasks()
	if err != nil {
		b.logger.Error("ListTasks failed", "error", err)
		return nil, toPublicError(err)
	}

	return tasks, nil
}

// FUNCTION: IDを指定してタスクを取得する関数
func (b *TaskBoundary) GetTask(id string) (*publish.Task, error) {
	task, err := b.service.GetTask(id)
	if err != nil {
		b.logger.Error("GetTask failed", "error", err)
		return nil, toPublicError(err)
	}

	return task, nil
}

// FUNCTION: 公開エラーに変換する
func toPublicError(err error) error {
	var message string

	var versionErr *services.UnsupportedMajorVersionError
	var notFound *services.NotFoundError

	if errors.As(err, &versionErr) {
		message = fmt.Sprintf(
			"データのバージョン %d.x には対応していません。対応バージョンは %d.x です。",
			versionErr.ActualMajor,
			versionErr.SupportedMajor,
		)
	} else if errors.As(err, &notFound) {
		message = "指定されたタスクが見つかりませんでした。"
	} else if errors.Is(err, services.ErrReadFailure) {
		message = "タスクを読み込めませんでした。"
	} else {
		message = "予期しないエラーが発生しました。"
	}

	return &PublicError{
		Processed: true,
		Message:   message,
	}
}

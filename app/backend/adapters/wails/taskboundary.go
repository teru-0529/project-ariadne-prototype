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

// FUNCTION: 公開エラーに変換する
func toPublicError(err error) error {
	var message string

	var versionErr *services.UnsupportedMajorVersionError
	if errors.As(err, &versionErr) {
		message = fmt.Sprintf(
			"データのバージョン %d.x には対応していません。対応バージョンは %d.x です。",
			versionErr.ActualMajor,
			versionErr.SupportedMajor,
		)
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

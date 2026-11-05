package wails

import (
	"errors"
	"fmt"

	"ariadne/backend/services"
	"ariadne/backend/services/public"
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
func (b *TaskBoundary) ListTasks() ([]public.TaskSummary, error) {
	tasks, err := b.service.ListTasks()
	if err != nil {
		b.logger.Error("ListTasks failed", "error", err)
		return nil, toPublicError(err)
	}

	return tasks, nil
}

// FUNCTION: IDを指定してタスクを取得する関数
func (b *TaskBoundary) GetTask(id string) (*public.Task, error) {
	task, err := b.service.GetTask(id)
	if err != nil {
		b.logger.Error("GetTask failed", "error", err)
		return nil, toPublicError(err)
	}

	return task, nil
}

// FUNCTION: 新しいタスクを登録する関数
func (b *TaskBoundary) CreateTask(request public.CreateTaskRequest) (string, error) {
	id, err := b.service.CreateTask(request)
	if err != nil {
		b.logger.Error("CreateTask failed", "error", err)
		return "", toPublicError(err)
	}

	return id, nil
}

// FUNCTION: 公開エラーに変換する
func toPublicError(err error) error {
	var message string

	var versionErr *services.UnsupportedMajorVersionError
	var notFound *services.NotFoundError
	var keyDuplicate *services.KeyDuplicateError

	if errors.As(err, &versionErr) {
		message = fmt.Sprintf(
			"データのバージョン %d.x には対応していません。対応バージョンは %d.x です。",
			versionErr.ActualMajor,
			versionErr.SupportedMajor,
		)
	} else if errors.As(err, &notFound) {
		message = "指定されたタスクが見つかりませんでした。"
	} else if errors.As(err, &keyDuplicate) {
		message = "同じIDのタスクが既に存在します。"
	} else if errors.Is(err, services.ErrReadFailure) {
		message = "タスクを読み込めませんでした。"
	} else if errors.Is(err, services.ErrWriteFailure) {
		message = "タスクを保存できませんでした。"
	} else {
		message = "予期しないエラーが発生しました。"
	}
	return &PublicError{
		Processed: true,
		Message:   message,
	}
}

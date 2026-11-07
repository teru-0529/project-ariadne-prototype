package services

import (
	"ariadne/backend/services/public"
	"fmt"

	"github.com/oklog/ulid/v2"
)

// INFO: タスクサービスの構造体
type TaskService struct {
	repo TaskRepository
}

// FUNCTION: TaskServiceのコンストラクタ関数
func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// FUNCTION: タスクの一覧を取得する関数
func (s *TaskService) ListTasks() ([]public.TaskSummary, error) {
	tasks, err := s.repo.Retrieve()
	if err != nil {
		return nil, err
	}

	summaries := make([]public.TaskSummary, 0, len(tasks))
	for _, task := range tasks {
		summaries = append(summaries, public.TaskSummary{
			ID:       task.ID,
			Title:    task.Title,
			Priority: task.Priority,
			Status:   task.Status,
		})
	}

	return summaries, nil
}

// FUNCTION: IDを指定してTaskを取得する関数
func (s *TaskService) GetTask(id string) (*public.Task, error) {
	task, err := s.repo.Pick(id)
	if err != nil {
		return nil, err
	}

	description := ""
	if task.Description != nil {
		description = *task.Description
	}

	return &public.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: description,
		Priority:    task.Priority,
		Status:      task.Status,
	}, nil
}

// FUNCTION: 新しいTaskを登録する関数
func (s *TaskService) CreateTask(request public.CreateTaskRequest) (string, error) {
	id := ulid.Make().String()

	task, err := NewTask(
		id,
		request.Title,
		request.Description,
		request.Priority,
		request.Status,
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrUnexpectedError, err)
	}

	if err := s.repo.Store(task, StoreDeposit); err != nil {
		return "", err
	}

	return id, nil
}

// FUNCTION: Taskを更新する関数 TODO:（A04-02 RED確認用スタブ）
func (s *TaskService) ModifyTask(id string, changes []TaskChange) error {
	return ErrUnexpectedError
}

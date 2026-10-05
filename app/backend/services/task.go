package services

// INFO: タスクを表す構造体
type Task struct {
	ID          string
	Title       string
	Description string
	Priority    int
	Status      string
}

// FUNCTION: 新しいタスクを作成する関数
func NewTask(id string, title string, description string, priority int, status string) (*Task, error) {
	// TODO: バリデーション実装
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      status,
	}, nil
}

// タスクの概要を表す構造体
type TaskSummary struct {
	ID       string
	Title    string
	Priority int
	Status   string
}

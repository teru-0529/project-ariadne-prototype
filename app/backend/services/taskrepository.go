package services

// INFO: タスクリポジトリのインターフェース
type TaskRepository interface {
	Retrieve() ([]*Task, error)
	Pick(id string) (*Task, error)
	Store(task *Task, intent StoreIntent) error
}

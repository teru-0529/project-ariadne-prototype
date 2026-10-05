package services

// INFO: タスクリポジトリのインターフェース
// TODO:
type TaskRepository interface {
	Retrieve() ([]*Task, error)
}

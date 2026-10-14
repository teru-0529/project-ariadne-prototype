package publish

// INFO: Task詳細取得用のPublic Model
type Task struct {
	ID          string
	Title       string
	Description string
	Priority    int
	Status      string
}

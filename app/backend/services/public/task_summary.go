package public

// INFO: Task一覧取得用のPublic Model
type TaskSummary struct {
	ID       string
	Title    string
	Priority int
	Status   string
}

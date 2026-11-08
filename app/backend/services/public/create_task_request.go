package public

// INFO: Task登録要求用のPublic Model
type CreateTaskRequest struct {
	Title       string
	Description string
	Priority    int
	Status      string
}

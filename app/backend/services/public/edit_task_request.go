package public

// INFO: Taskの単一Propertyに対する変更要求
type TaskPropertyChange struct {
	Type  string
	Value any
}

// INFO: Taskの更新要求
type EditTaskRequest struct {
	Changes []TaskPropertyChange
}

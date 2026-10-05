package file

import (
	services "ariadne/backend/services"
)

// INFO: YAML形式のタスクデータセットを表す構造体
type taskDataset struct {
	Format  string            `yaml:"format"`
	Version string            `yaml:"version"`
	Tasks   []persistenceTask `yaml:"tasks"`
}

// INFO: YAML形式のタスクデータを表す構造体
type persistenceTask struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Priority    int    `yaml:"priority"`
	Status      string `yaml:"status"`
	CreatedAt   string `yaml:"createdAt"`
	UpdatedAt   string `yaml:"updatedAt"`
}

// FUNCTION: persistenceTaskをservices.Taskに変換する関数
func (p *persistenceTask) toTask() (*services.Task, error) {
	return services.NewTask(
		p.ID,
		p.Title,
		p.Description,
		p.Priority,
		p.Status,
	)
}

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
	ID          *string `yaml:"id"`
	Title       *string `yaml:"title"`
	Description *string `yaml:"description"`
	Priority    *int    `yaml:"priority"`
	Status      *string `yaml:"status"`
	CreatedAt   *string `yaml:"createdAt"`
	UpdatedAt   *string `yaml:"updatedAt"`
}

// FUNCTION: persistenceTaskをservices.Taskに変換する関数
func (p *persistenceTask) toTask() (*services.Task, error) {
	if err := p.validateFormat(); err != nil {
		return nil, err
	}

	return services.NewTask(
		*p.ID,
		*p.Title,
		*p.Description,
		*p.Priority,
		*p.Status,
	)
}

// FUNCTION: persistenceTaskのフォーマットを検証する関数
func (p *persistenceTask) validateFormat() error {
	if p.ID == nil {
		return &services.FormatError{Detail: "id is missing"}
	}
	if p.Title == nil {
		return &services.FormatError{Detail: "title is missing"}
	}
	if p.Description == nil {
		return &services.FormatError{Detail: "description is missing"}
	}
	if p.Priority == nil {
		return &services.FormatError{Detail: "priority is missing"}
	}
	if p.Status == nil {
		return &services.FormatError{Detail: "status is missing"}
	}
	if p.CreatedAt == nil {
		return &services.FormatError{Detail: "createdAt is missing"}
	}
	if p.UpdatedAt == nil {
		return &services.FormatError{Detail: "updatedAt is missing"}
	}

	return nil
}

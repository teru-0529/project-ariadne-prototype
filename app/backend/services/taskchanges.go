package services

import "fmt"

// INFO: Titleの変更指示
type TitleChange struct {
	value string
}

// FUNCTION: TitleChangeのコンストラクタ関数
func NewTitleChange(value any) (TaskChange, error) {
	v, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: title must be string", ErrUnexpectedError)
	}
	return &TitleChange{value: v}, nil
}

// FUNCTION: TitleChangeのapplyメソッド
func (c *TitleChange) apply(edit *TaskEdit) {
}

// INFO: Descriptionの変更指示
type DescriptionChange struct {
	value string
}

// FUNCTION: DescriptionChangeのコンストラクタ関数
func NewDescriptionChange(value any) (TaskChange, error) {
	v, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: description must be string", ErrUnexpectedError)
	}
	return &DescriptionChange{value: v}, nil
}

// FUNCTION: DescriptionChangeのapplyメソッド
func (c *DescriptionChange) apply(edit *TaskEdit) {
}

// INFO: Priorityの変更指示
type PriorityChange struct {
	value int
}

// FUNCTION: PriorityChangeのコンストラクタ関数
func NewPriorityChange(value any) (TaskChange, error) {
	v, ok := value.(int)
	if !ok {
		return nil, fmt.Errorf("%w: priority must be int", ErrUnexpectedError)
	}
	return &PriorityChange{value: v}, nil
}

// FUNCTION: PriorityChangeのapplyメソッド
func (c *PriorityChange) apply(edit *TaskEdit) {
}

// INFO: Statusの変更指示
type StatusChange struct {
	value string
}

// FUNCTION: StatusChangeのコンストラクタ関数
func NewStatusChange(value any) (TaskChange, error) {
	v, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: status must be string", ErrUnexpectedError)
	}
	return &StatusChange{value: v}, nil
}

// FUNCTION: StatusChangeのapplyメソッド
func (c *StatusChange) apply(edit *TaskEdit) {
}

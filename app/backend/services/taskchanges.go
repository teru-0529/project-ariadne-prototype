package services

import "fmt"

// INFO: Titleの変更指示
type TitleChange struct {
	value string
}

func NewTitleChange(value any) (TaskChange, error) {
	v, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: title must be string", ErrUnexpectedError)
	}
	return &TitleChange{value: v}, nil
}

// INFO: Descriptionの変更指示
type DescriptionChange struct {
	value string
}

func NewDescriptionChange(value any) (TaskChange, error) {
	v, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: description must be string", ErrUnexpectedError)
	}
	return &DescriptionChange{value: v}, nil
}

// INFO: Priorityの変更指示
type PriorityChange struct {
	value int
}

func NewPriorityChange(value any) (TaskChange, error) {
	v, ok := value.(int)
	if !ok {
		return nil, fmt.Errorf("%w: priority must be int", ErrUnexpectedError)
	}
	return &PriorityChange{value: v}, nil
}

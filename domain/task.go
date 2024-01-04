package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskID int

type (
	// タスクが持つ抽象的な実装
	TaskRepository interface {
		Create(context.Context, Task) (Task, error)
		// UpdateTitle() error
		// FindAll() ([]Task, error)
	}

	// タスク
	Task struct {
		ID        TaskID
		Title     string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

package domain

import (
	"context"
	"time"
)

type (
	// タスクが持つ抽象的な実装
	TaskRepository interface {
		Create(context.Context, Task) (Task, error)
		// Update(context.Context, Task) (Task, error)
		// FindAll() ([]Task, error)
	}

	// タスク
	Task struct {
		ID        int
		Title     string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

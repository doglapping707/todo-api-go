package domain

import (
	"context"
	"time"
)

type TaskID uint64

type (
	// タスクが持つ抽象的な実装
	TaskRepository interface {
		Create(context.Context, Task) (Task, error)
		Update(context.Context, Task, TaskID) error
		FindAll(context.Context) ([]Task, error)
	}

	// タスク
	Task struct {
		ID        TaskID
		Title     string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

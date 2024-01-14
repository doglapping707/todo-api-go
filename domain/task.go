package domain

import (
	"context"
	"time"
)

type TaskID uint64

type (
	TaskRepository interface {
		Create(context.Context, Task) (Task, error)
		Update(context.Context, Task, TaskID) error
		FindAll(context.Context) ([]Task, error)
	}

	Task struct {
		ID        TaskID
		Title     string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

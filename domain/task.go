package domain

import (
	"context"
	"strconv"
	"time"
)

type TaskID uint64

func Uint64(taskID string) (uint64, error) {
	return strconv.ParseUint(taskID, 10, 64)
}

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

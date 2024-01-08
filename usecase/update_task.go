package usecase

import (
	"context"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type (
	UpdateTaskUseCase interface {
		Execute(context.Context, UpdateTaskInput) error
	}

	UpdateTaskInput struct {
		ID    uint64 `json:"id" validate:"required"`
		Title string `json:"title" validate:"required,gte=1,lte=15"`
	}

	UpdateTaskInteractor struct {
		repo       domain.TaskRepository
		ctxTimeout time.Duration
	}
)

func NewUpdateTaskInteractor(
	taskRepo domain.TaskRepository,
	t time.Duration,
) UpdateTaskUseCase {
	return UpdateTaskInteractor{
		repo: taskRepo,
		ctxTimeout: t,
	}
}

func (t UpdateTaskInteractor) Execute(ctx context.Context, input UpdateTaskInput) error {
	// タイムアウトを設定したコンテキストを取得する
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	// タスクを成形する
	var task = domain.Task{
		ID:        input.ID,
		Title:     input.Title,
		UpdatedAt: time.Now(),
	}

	// タスクを更新する
	err := t.repo.Update(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

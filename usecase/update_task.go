package usecase

import (
	"context"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type (
	UpdateTaskUseCase interface {
		Execute(context.Context, UpdateTaskInput) (UpdateTaskOutput, error)
	}

	UpdateTaskInput struct {
		ID    uint64 `json:"id" validate:"required"`
		Title string `json:"title" validate:"required,gte=1,lte=15"`
	}

	UpdateTaskPresenter interface {
		Output(domain.Task) UpdateTaskOutput
	}

	UpdateTaskOutput struct {
		ID        uint64 `json:"id"`
		Title     string `json:"title"`
		UpdatedAt string `json:"updated_at"`
	}

	UpdateTaskInteractor struct {
		repo       domain.TaskRepository
		presenter  UpdateTaskPresenter
		ctxTimeout time.Duration
	}
)

func NewUpdateTaskInteractor(
	taskRepo domain.TaskRepository,
	presenter  UpdateTaskPresenter,
	t time.Duration,
) UpdateTaskUseCase {
	return UpdateTaskInteractor{
		repo: taskRepo,
		presenter: presenter,
		ctxTimeout: t,
	}
}

func (t UpdateTaskInteractor) Execute(ctx context.Context, input UpdateTaskInput) (UpdateTaskOutput, error) {
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
	task, err := t.repo.Update(ctx, task)
	if err != nil {
		return t.presenter.Output(domain.Task{}), err
	}

	// 出力用のタスクを返却する
	return t.presenter.Output(task), nil
}

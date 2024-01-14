package usecase

import (
	"context"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type (
	CreateTaskUseCase interface {
		Execute(context.Context, CreateTaskInput) (CreateTaskOutput, error)
	}

	CreateTaskInput struct {
		Title string `json:"title" validate:"required,gte=1,lte=15"`
	}

	CreateTaskPresenter interface {
		Output(domain.Task) CreateTaskOutput
	}

	CreateTaskOutput struct {
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	createTaskInteractor struct {
		repo       domain.TaskRepository
		presenter  CreateTaskPresenter
		ctxTimeout time.Duration
	}
)

func NewCreateTaskInteractor(
	repo domain.TaskRepository,
	presenter  CreateTaskPresenter,
	t time.Duration,
) CreateTaskUseCase {
	return createTaskInteractor{
		repo: repo,
		presenter: presenter,
		ctxTimeout: t,
	}
}

func (t createTaskInteractor) Execute(ctx context.Context, input CreateTaskInput) (CreateTaskOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	var task = domain.Task{
		Title: input.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	task, err := t.repo.Create(ctx, task)
	if err != nil {
		return t.presenter.Output(domain.Task{}), err
	}

	return t.presenter.Output(task), nil
}

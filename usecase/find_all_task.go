package usecase

import (
	"context"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type (
	FindAllTaskUseCase interface {
		Execute(context.Context) ([]FindAllTaskOutput, error)
	}

	FindAllTaskPresenter interface {
		Output([]domain.Task) []FindAllTaskOutput
	}

	FindAllTaskOutput struct {
		ID        uint64  `json:"id"`
		Title     string  `json:"title"`
	}

	findAllTaskInteractor struct {
		repo       domain.TaskRepository
		presenter  FindAllTaskPresenter
		ctxTimeout time.Duration
	}
)

func NewFindAllTaskInteractor(
	repo domain.TaskRepository,
	presenter FindAllTaskPresenter,
	t time.Duration,
) FindAllTaskUseCase {
	return findAllTaskInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

func (t findAllTaskInteractor) Execute(ctx context.Context) ([]FindAllTaskOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	tasks, err := t.repo.FindAll(ctx)
	if err != nil {
		return t.presenter.Output([]domain.Task{}), err
	}

	return t.presenter.Output(tasks), nil
}

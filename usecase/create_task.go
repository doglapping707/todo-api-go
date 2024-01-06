package usecase

import (
	"context"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type (
	// タスク作成が持つ抽象的な実装
	CreateTaskUseCase interface {
		Execute(context.Context, CreateTaskInput) (CreateTaskOutput, error)
	}

	// タスク作成の入力時
	CreateTaskInput struct {
		Title string `json:"title" validate:"required,gte=1,lte=15"`
	}

	// タスク作成の成形を行う抽象的な実装
	CreateTaskPresenter interface {
		Output(domain.Task) CreateTaskOutput
	}

	// タスク作成の出力時
	CreateTaskOutput struct {
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	// タスク作成
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
	// タイムアウトを設定したコンテキストを取得する
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	// タスクを成形する
	var task = domain.Task{
		Title: input.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// タスクを作成する
	task, err := t.repo.Create(ctx, task)
	if err != nil {
		return t.presenter.Output(domain.Task{}), err
	}

	// 出力用のタスクを返却する
	return t.presenter.Output(task), nil
}

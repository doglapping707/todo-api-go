package presenter

import (
	"time"

	"github.com/doglapping707/todo-api-go/domain"
	"github.com/doglapping707/todo-api-go/usecase"
)

type createTaskPresenter struct{}

func NewCreateTaskPresenter() usecase.CreateTaskPresenter {
	return createTaskPresenter{}
}

func (t createTaskPresenter) Output(task domain.Task) usecase.CreateTaskOutput {
	return usecase.CreateTaskOutput{
		Title:     task.Title,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
		UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
	}
}

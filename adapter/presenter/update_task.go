package presenter

import (
	"time"

	"github.com/doglapping707/todo-api-go/domain"
	"github.com/doglapping707/todo-api-go/usecase"
)

type updateTaskPresenter struct{}

func NewUpdateTaskPresenter() usecase.UpdateTaskPresenter {
	return updateTaskPresenter{}
}

func (t updateTaskPresenter) Output(task domain.Task) usecase.UpdateTaskOutput {
	return usecase.UpdateTaskOutput{
		ID:        task.ID,
		Title:     task.Title,
		UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
	}
}
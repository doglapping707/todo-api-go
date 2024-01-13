package presenter

import (
	"github.com/doglapping707/todo-api-go/domain"
	"github.com/doglapping707/todo-api-go/usecase"
)

type findAllTaskPresenter struct{}

func NewFindAllTaskPresenter() usecase.FindAllTaskPresenter {
	return findAllTaskPresenter{}
}

func (a findAllTaskPresenter) Output(tasks []domain.Task) []usecase.FindAllTaskOutput {
	var o = make([]usecase.FindAllTaskOutput, 0)

	for _, task := range tasks {
		o = append(o, usecase.FindAllTaskOutput{
			ID:        task.ID,
			Title:     task.Title,
		})
	}

	return o
}

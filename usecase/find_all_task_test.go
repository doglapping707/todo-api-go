package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type mockTaskRepoFindAll struct {
	domain.TaskRepository

	result []domain.Task
	err    error
}

func (m mockTaskRepoFindAll) FindAll(_ context.Context) ([]domain.Task, error) {
	return m.result, m.err
}

type mockFindAllTaskPresenter struct {
	result []FindAllTaskOutput
}

func (m mockFindAllTaskPresenter) Output(_ []domain.Task) []FindAllTaskOutput {
	return m.result
}

func TestFindAllTaskInteractor_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		repository    domain.TaskRepository
		presenter     FindAllTaskPresenter
		expected      []FindAllTaskOutput
		expectedError interface{}
	}{
		{
			name: "Success when returning the task list",
			repository: mockTaskRepoFindAll{
				result: []domain.Task{
					{
						ID:    1,
						Title: "Task_1",
					},
					{
						ID:    2,
						Title: "Task_2",
					},
				},
				err: nil,
			},
			presenter: mockFindAllTaskPresenter{
				result: []FindAllTaskOutput{
					{
						ID:    1,
						Title: "Task_1",
					},
					{
						ID:    2,
						Title: "Task_2",
					},
				},
			},
			expected: []FindAllTaskOutput{
				{
					ID:    1,
					Title: "Task_1",
				},
				{
					ID:    2,
					Title: "Task_2",
				},
			},
		},
		{
			name: "Success when returning the empty task list",
			repository: mockTaskRepoFindAll{
				result: []domain.Task{},
				err:    nil,
			},
			presenter: mockFindAllTaskPresenter{
				result: []FindAllTaskOutput{},
			},
			expected: []FindAllTaskOutput{},
		},
		{
			name: "Error when returning the list of tasks",
			repository: mockTaskRepoFindAll{
				result: []domain.Task{},
				err:    errors.New("error"),
			},
			presenter: mockFindAllTaskPresenter{
				result: []FindAllTaskOutput{},
			},
			expectedError: "error",
			expected:      []FindAllTaskOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewFindAllTaskInteractor(tt.repository, tt.presenter, time.Second)

			result, err := uc.Execute(context.Background())
			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}

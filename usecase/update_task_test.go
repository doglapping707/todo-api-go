package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type mockTaskRepo struct {
	domain.TaskRepository

	result domain.Task
	err    error
}

func (t mockTaskRepo) Update(_ context.Context, _ domain.Task) (domain.Task, error) {
	return t.result, t.err
}

type mockUpdateTaskPresenter struct {
	result UpdateTaskOutput
}

func (m mockUpdateTaskPresenter) Output(_ domain.Task) UpdateTaskOutput {
	return m.result
}

func TestUpdateTaskInteractor_Execute(t *testing.T) {
	// 並列処理を可能にする
	t.Parallel()

	type args struct {
		input UpdateTaskInput
	}

	tests := []struct {
		name          string
		args          args
		repository    domain.TaskRepository
		presenter     UpdateTaskPresenter
		expected      UpdateTaskOutput
		expectedError interface{}
	}{
		// 正常値
		{
			name: "Update account successful",
			args: args{
				input: UpdateTaskInput{
					ID:    1,
					Title: "Test Task2",
				},
			},
			repository: mockTaskRepo{
				result: domain.Task{
					ID:        1,
					Title:     "Test Task2",
					UpdatedAt: time.Time{},
				},
				err: nil,
			},
			presenter: mockUpdateTaskPresenter{
				result: UpdateTaskOutput{
					ID:        1,
					Title:     "Test Task2",
					UpdatedAt: time.Time{}.String(),
				},
			},
			expected: UpdateTaskOutput{
				ID:        1,
				Title:     "Test Task2",
				UpdatedAt: time.Time{}.String(),
			},
		},

		// 異常値
		{
			name: "Update task generic error",
			args: args{
				input: UpdateTaskInput{
					ID:    0,
					Title: "",
				},
			},
			repository: mockTaskRepo{
				result: domain.Task{},
				err: nil,
			},
			presenter: mockUpdateTaskPresenter{
				result: UpdateTaskOutput{},
			},
			expected: UpdateTaskOutput{},
			expectedError: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewUpdateTaskInteractor(tt.repository, tt.presenter, time.Second)

			result, err := uc.Execute(context.TODO(), tt.args.input)
			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}

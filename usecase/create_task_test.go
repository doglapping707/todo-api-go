package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type mockTaskRepoStore struct {
	domain.TaskRepository

	result domain.Task
	err    error
}

func (m mockTaskRepoStore) Create(_ context.Context, _ domain.Task) (domain.Task, error) {
	return m.result, m.err
}

type mockCreateTaskPresenter struct {
	result CreateTaskOutput
}

func (m mockCreateTaskPresenter) Output(_ domain.Task) CreateTaskOutput {
	return m.result
}

func TestCreateTaskInteractor_Execute(t *testing.T) {
	// 並列処理を可能にする
	t.Parallel()

	type args struct {
		input CreateTaskInput
	}

	// テスト対象
	tests := []struct {
		name          string
		args          args
		repository    domain.TaskRepository
		presenter     CreateTaskPresenter
		expected      CreateTaskOutput
		expectedError interface{}
	}{
		// 正常値
		{
			name: "Create task successful",
			args: args{
				input: CreateTaskInput{
					Title: "Test",
				},
			},
			repository: mockTaskRepoStore{
				result: domain.Task{
					Title:     "Test",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				err: nil,
			},
			presenter: mockCreateTaskPresenter{
				result: CreateTaskOutput{
					Title:     "Test",
					CreatedAt: time.Time{}.String(),
					UpdatedAt: time.Time{}.String(),
				},
			},
			expected: CreateTaskOutput{
				Title:     "Test",
				CreatedAt: time.Time{}.String(),
				UpdatedAt: time.Time{}.String(),
			},
		},

		// 異常値
		{
			name: "Create task generic error",
			args: args{
				input: CreateTaskInput{
					Title: "",
				},
			},
			repository: mockTaskRepoStore{
				result: domain.Task{},
				err:    errors.New("error"),
			},
			presenter: mockCreateTaskPresenter{
				result: CreateTaskOutput{},
			},
			expectedError: "error",
			expected:      CreateTaskOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewCreateTaskInteractor(tt.repository, tt.presenter, time.Second)

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
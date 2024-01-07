package presenter

import (
	"reflect"
	"testing"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
	"github.com/doglapping707/todo-api-go/usecase"
)

func Test_updateTaskPresenter_Output(t *testing.T) {
	type args struct {
		task domain.Task
	}

	tests := []struct {
		name string
		args args
		want usecase.UpdateTaskOutput
	}{
		{
			name: "Update task output",

			// 入力値
			args: args{
				task: domain.Task{
					ID:        1,
					Title:     "Test Task",
					UpdatedAt: time.Time{},
				},
			},

			// 期待値
			want: usecase.UpdateTaskOutput{
				ID:        1,
				Title:     "Test Task",
				UpdatedAt: "0001-01-01T00:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewCreateTaskPresenter()
			if got := pre.Output(tt.args.task); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

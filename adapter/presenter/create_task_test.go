package presenter

import (
	"reflect"
	"testing"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
	"github.com/doglapping707/todo-api-go/usecase"
)

func Test_createTaskPresenter_Output(t *testing.T) {
	type args struct {
		task domain.Task
	}

	tests := []struct {
		name string
		args args
		want usecase.CreateTaskOutput
	}{
		{
			name: "Create task output",
	
			// 入力値
			args: args{
				task: domain.Task{
					Title:     "Testing",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},

			// 期待値
			want: usecase.CreateTaskOutput{
				Title:     "Testing",
				CreatedAt: "0001-01-01T00:00:00Z",
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
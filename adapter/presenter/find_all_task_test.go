package presenter

import (
	"reflect"
	"testing"

	"github.com/doglapping707/todo-api-go/domain"
	"github.com/doglapping707/todo-api-go/usecase"
)

func Test_findAllTaskPresenter_Output(t *testing.T) {
	type args struct {
		tasks []domain.Task
	}
	tests := []struct {
		name string
		args args
		want []usecase.FindAllTaskOutput
	}{
		{
			name: "Find all task output",
			args: args{
				tasks: []domain.Task{
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
			want: []usecase.FindAllTaskOutput{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewFindAllTaskPresenter()
			if got := pre.Output(tt.args.tasks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

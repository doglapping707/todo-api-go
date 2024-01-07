package action

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/doglapping707/todo-api-go/infrastructure/log"
	"github.com/doglapping707/todo-api-go/infrastructure/validation"
	"github.com/doglapping707/todo-api-go/usecase"
)

type mockUpdateTask struct {
	result usecase.UpdateTaskOutput
	err    error
}

func (m mockUpdateTask) Execute(_ context.Context, _ usecase.UpdateTaskInput) (usecase.UpdateTaskOutput, error) {
	return m.result, m.err
}

func TestUpdateTaskAction_Execute(t *testing.T) {
	t.Parallel()

	validator, _ := validation.NewValidatorFactory(validation.InstanceGoPlayground)

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.UpdateTaskUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		// 正常値
		{
			// input
			name: "UpdateTaskAction success",
			args: args{
				rawPayload: []byte(
					`{
						"id": 1,
						"title": "Test Task"
					}`,
				),
			},

			// output
			ucMock: mockUpdateTask{
				result: usecase.UpdateTaskOutput{
					ID:        1,
					Title:     "Test Task",
					UpdatedAt: time.Time{}.String(),
				},
				err: nil,
			},

			// 期待値
			expectedBody:       `{"id":1,"title":"Test Task","updated_at":"0001-01-01 00:00:00 +0000 UTC"}`,
			expectedStatusCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPut,
				"/tasks",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewUpdateTaskAction(tt.ucMock, log.LoggerMock{}, validator)
			)

			action.Execute(w, req)

			// レスポンスされたステータスコードが期待値と異なる場合
			if w.Code != tt.expectedStatusCode {
				// エラーを記録する
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					w.Code,
					tt.expectedStatusCode,
				)
			}

			// レスポンスされたボディが期待値と異なる場合
			var result = strings.TrimSpace(w.Body.String())
			if !strings.EqualFold(result, tt.expectedBody) {
				// エラーを記録する
				t.Errorf(
					"[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					result,
					tt.expectedBody,
				)
			}
		})
	}
}
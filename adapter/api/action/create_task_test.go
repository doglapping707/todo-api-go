package action

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/doglapping707/todo-api-go/infrastructure/log"
	"github.com/doglapping707/todo-api-go/infrastructure/validation"
	"github.com/doglapping707/todo-api-go/usecase"
)

type mockCreateTask struct {
	result usecase.CreateTaskOutput
	err    error
}

func (m mockCreateTask) Execute(_ context.Context, _ usecase.CreateTaskInput) (usecase.CreateTaskOutput, error) {
	return m.result, m.err
}

func TestCreateTaskAction_Execute(t *testing.T) {
	// 並列処理を可能にする
	t.Parallel()

	validator, _ := validation.NewValidatorFactory(validation.InstanceGoPlayground)

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.CreateTaskUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		// 正常値
		{
			// input
			name: "CreateTaskAction success",
			args: args{
				rawPayload: []byte(
					`{
						"title": "Test Task"
					}`,
				),
			},

			// output
			ucMock: mockCreateTask{
				result: usecase.CreateTaskOutput{
					Title:     "Test Task",
					CreatedAt: time.Time{}.String(),
					UpdatedAt: time.Time{}.String(),
				},
				err: nil,
			},

			// 期待値
			expectedBody:       `{"title":"Test Task","created_at":"0001-01-01 00:00:00 +0000 UTC","updated_at":"0001-01-01 00:00:00 +0000 UTC"}`,
			expectedStatusCode: http.StatusCreated,
		},

		// 異常値
		{
			// input
			name: "CreateTaskAction generic error",
			args: args{
				rawPayload: []byte(
					`{
						"title": "Test Task"
					}`,
				),
			},

			// output
			ucMock: mockCreateTask{
				result: usecase.CreateTaskOutput{},
				err:    errors.New("error"),
			},

			// 期待値
			expectedBody:       `{"errors":["error"]}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			// input
			name: "CreateTaskAction error invalid title",
			args: args{
				rawPayload: []byte(
					`{
						"title": ""
					}`,
				),
			},

			// output
			ucMock: mockCreateTask{
				result: usecase.CreateTaskOutput{},
				err:    errors.New("error"),
			},

			// 期待値
			expectedBody:       `{"errors":["Title is a required field"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			// input
			name: "CreateTaskAction error invalid title",
			args: args{
				rawPayload: []byte(
					`{
						"title": "aaaabbbbccccdddd"
					}`,
				),
			},

			// output
			ucMock: mockCreateTask{
				result: usecase.CreateTaskOutput{},
				err:    errors.New("error"),
			},

			// 期待値
			expectedBody:       `{"errors":["Title must be at maximum 15 characters in length"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			// input
			name: "CreateTaskAction error invalid field",
			args: args{
				rawPayload: []byte(
					`{
						"title1234": "Test Task"
					}`,
				),
			},

			// output
			ucMock: mockCreateTask{
				result: usecase.CreateTaskOutput{},
				err:    nil,
			},

			// 期待値
			expectedBody:       `{"errors":["Title is a required field"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			// input
			name: "CreateTaskAction error invalid JSON",
			args: args{
				rawPayload: []byte(
					`{
						"title":
					}`,
				),
			},

			// output
			ucMock: mockCreateTask{
				result: usecase.CreateTaskOutput{},
				err:    nil,
			},

			// 期待値
			expectedBody:       `{"errors":["invalid character '}' looking for beginning of value"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewCreateTaskAction(tt.ucMock, log.LoggerMock{}, validator)
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

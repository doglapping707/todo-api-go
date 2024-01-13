package action

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/doglapping707/todo-api-go/infrastructure/log"
	"github.com/doglapping707/todo-api-go/usecase"
)

type mockFindAllTask struct {
	result []usecase.FindAllTaskOutput
	err    error
}

func (m mockFindAllTask) Execute(_ context.Context) ([]usecase.FindAllTaskOutput, error) {
	return m.result, m.err
}

func TestFindAllTaskAction_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		ucMock             usecase.FindAllTaskUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "FindAllTaskAction success one task",
			ucMock: mockFindAllTask{
				result: []usecase.FindAllTaskOutput{
					{
						ID:    1,
						Title: "Task_1",
					},
				},
				err: nil,
			},
			expectedBody:       `[{"id":1,"title":"Task_1"}]`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAllTaskAction success empty",
			ucMock: mockFindAllTask{
				result: []usecase.FindAllTaskOutput{},
				err:    nil,
			},
			expectedBody:       `[]`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAllTaskAction generic error",
			ucMock: mockFindAllTask{
				err: errors.New("error"),
			},
			expectedBody:       `{"errors":["error"]}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

			var (
				w      = httptest.NewRecorder()
				action = NewFindAllTaskAction(tt.ucMock, log.LoggerMock{})
			)

			action.Execute(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					w.Code,
					tt.expectedStatusCode,
				)
			}

			var result = strings.TrimSpace(w.Body.String())
			if !strings.EqualFold(result, tt.expectedBody) {
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

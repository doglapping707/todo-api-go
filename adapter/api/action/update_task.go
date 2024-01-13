package action

import (
	"encoding/json"
	"net/http"

	"github.com/doglapping707/todo-api-go/adapter/api/logging"
	"github.com/doglapping707/todo-api-go/adapter/api/response"
	"github.com/doglapping707/todo-api-go/adapter/logger"
	"github.com/doglapping707/todo-api-go/adapter/validator"
	"github.com/doglapping707/todo-api-go/domain"
	"github.com/doglapping707/todo-api-go/usecase"
)

type UpdateTaskAction struct {
	uc        usecase.UpdateTaskUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewUpdateTaskAction(uc usecase.UpdateTaskUseCase, log logger.Logger, v validator.Validator) UpdateTaskAction {
	return UpdateTaskAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (t UpdateTaskAction) Execute(w http.ResponseWriter, r *http.Request) {
	var logKey = "update_task"

	var taskID, err = domain.Uint64(r.URL.Query().Get("task_id"))
	if err != nil {
		var err = response.ErrParameterInvalid
		logging.NewError(
			t.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log("invalid parameter")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	var input usecase.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.NewError(
			t.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log("error when decoding json")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := t.validateInput(input); len(errs) > 0 {
		logging.NewError(
			t.log,
			response.ErrInvalidInput,
			logKey,
			http.StatusBadRequest,
		).Log("invalid input")

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	if err := t.uc.Execute(r.Context(), input, domain.TaskID(taskID)); err != nil {
		logging.NewError(
			t.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when updating a new task")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(t.log, logKey, http.StatusNoContent).Log("success updating task")

	response.NewSuccess(nil, http.StatusNoContent).Send(w)
}

func (t UpdateTaskAction) validateInput(input usecase.UpdateTaskInput) []string {
	var msgs []string

	if err := t.validator.Validate(input); err != nil {
		for _, msg := range t.validator.Messages() {
			msg := msg
			msgs = append(msgs, msg)
		}
	}

	return msgs
}

package action

import (
	"encoding/json"
	"net/http"

	"github.com/doglapping707/todo-api-go/adapter/api/logging"
	"github.com/doglapping707/todo-api-go/adapter/api/response"
	"github.com/doglapping707/todo-api-go/adapter/logger"
	"github.com/doglapping707/todo-api-go/adapter/validator"
	"github.com/doglapping707/todo-api-go/usecase"
)

type CreateTaskAction struct {
	uc        usecase.CreateTaskUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewCreateTaskAction(uc usecase.CreateTaskUseCase, log logger.Logger, v validator.Validator) CreateTaskAction {
	return CreateTaskAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (t CreateTaskAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_task"

	var input usecase.CreateTaskInput
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

	output, err := t.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			t.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when creating a new task")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(t.log, logKey, http.StatusCreated).Log("success creating task")

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

func (t CreateTaskAction) validateInput(input usecase.CreateTaskInput) []string {
	var msgs []string

	err := t.validator.Validate(input)
	if err != nil {
		for _, msg := range t.validator.Messages() {
			msg := msg
			msgs = append(msgs, msg)
		}
	}

	return msgs
}

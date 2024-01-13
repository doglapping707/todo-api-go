package action

import (
	"net/http"

	"github.com/doglapping707/todo-api-go/adapter/api/logging"
	"github.com/doglapping707/todo-api-go/adapter/api/response"
	"github.com/doglapping707/todo-api-go/adapter/logger"
	"github.com/doglapping707/todo-api-go/usecase"
)

type FindAllTaskAction struct {
	uc  usecase.FindAllTaskUseCase
	log logger.Logger
}

func NewFindAllTaskAction(uc usecase.FindAllTaskUseCase, log logger.Logger) FindAllTaskAction {
	return FindAllTaskAction{
		uc:  uc,
		log: log,
	}
}

func (a FindAllTaskAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_all_task"

	output, err := a.uc.Execute(r.Context())
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when returning task list")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success when returning task list")

	response.NewSuccess(output, http.StatusOK).Send(w)
}

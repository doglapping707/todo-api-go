package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/doglapping707/todo-api-go/adapter/api/action"
	"github.com/doglapping707/todo-api-go/adapter/api/middleware"
	"github.com/doglapping707/todo-api-go/adapter/logger"
	"github.com/doglapping707/todo-api-go/adapter/presenter"
	"github.com/doglapping707/todo-api-go/adapter/repository"
	"github.com/doglapping707/todo-api-go/adapter/validator"
	"github.com/doglapping707/todo-api-go/usecase"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// マルチプレクサ
type gorillaMux struct {
	router     *mux.Router
	middleware *negroni.Negroni
	log        logger.Logger
	db         repository.SQL
	validator  validator.Validator
	port       Port
	ctxTimeout time.Duration
}

// マルチプレクサを返却する
func newGorillaMux(
	log logger.Logger,
	db repository.SQL,
	validator validator.Validator,
	port Port,
	t time.Duration,
) *gorillaMux {
	return &gorillaMux{
		router:     mux.NewRouter(),
		middleware: negroni.New(),
		log:        log,
		db:         db,
		validator:  validator,
		port:       port,
		ctxTimeout: t,
	}
}

func (g gorillaMux) Listen() {
	// HTTPハンドラーをセットする
	g.setAppHandlers(g.router)
	// HTTPハンドラーを登録する
	g.middleware.UseHandler(g.router)

	// HTTPサーバーを起動するためのパラメータを成形する
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      g.middleware,
	}

	// シグナルの受付を開始する
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// HTTPサーバーを起動する
	go func() {
		g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
		if err := server.ListenAndServe(); err != nil {
			g.log.WithError(err).Fatalln("Error starting HTTP server")
		}
	}()

	// シグナルを検知する
	<-stop

	// 指定時間でタイムアウトするコンテキストを作成する
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	// サーバーを停止する
	if err := server.Shutdown(ctx); err != nil {
		g.log.WithError(err).Fatalln("Server Shutdown Failed")
	}

	// ログを出力する
	g.log.Infof("Service down")
}

func (g gorillaMux) setAppHandlers(router *mux.Router) {
	// プレフィックスを設定する
	api := router.PathPrefix("/v1").Subrouter()

	// 送金
	api.Handle("/transfers", g.buildCreateTransferAction()).Methods(http.MethodPost)
	api.Handle("/transfers", g.buildFindAllTransferAction()).Methods(http.MethodGet)

	// アカウント
	api.Handle("/accounts/{account_id}/balance", g.buildFindBalanceAccountAction()).Methods(http.MethodGet)
	api.Handle("/accounts", g.buildCreateAccountAction()).Methods(http.MethodPost)
	api.Handle("/accounts", g.buildFindAllAccountAction()).Methods(http.MethodGet)

	// タスク
	api.Handle("/tasks", g.buildCreateTaskAction()).Methods(http.MethodPost)
	api.Handle("/tasks/{task_id}", g.buildUpdateTaskAction()).Methods(http.MethodPut)
	api.Handle("/tasks", g.buildFindAllTaskAction()).Methods(http.MethodGet)

	// ヘルスチェック
	api.HandleFunc("/health", action.HealthCheck).Methods(http.MethodGet)
}

func (g gorillaMux) buildCreateTransferAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewCreateTransferInteractor(
				repository.NewTransferSQL(g.db),
				repository.NewAccountSQL(g.db),
				presenter.NewCreateTransferPresenter(),
				g.ctxTimeout,
			)
			act = action.NewCreateTransferAction(uc, g.log, g.validator)
		)

		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildFindAllTransferAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewFindAllTransferInteractor(
				repository.NewTransferSQL(g.db),
				presenter.NewFindAllTransferPresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAllTransferAction(uc, g.log)
		)

		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildCreateAccountAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewCreateAccountInteractor(
				repository.NewAccountSQL(g.db),
				presenter.NewCreateAccountPresenter(),
				g.ctxTimeout,
			)
			act = action.NewCreateAccountAction(uc, g.log, g.validator)
		)

		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildFindAllAccountAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewFindAllAccountInteractor(
				repository.NewAccountSQL(g.db),
				presenter.NewFindAllAccountPresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAllAccountAction(uc, g.log)
		)

		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildFindBalanceAccountAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewFindBalanceAccountInteractor(
				repository.NewAccountSQL(g.db),
				presenter.NewFindAccountBalancePresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAccountBalanceAction(uc, g.log)
		)

		var (
			vars = mux.Vars(req)
			q    = req.URL.Query()
		)

		q.Add("account_id", vars["account_id"])
		req.URL.RawQuery = q.Encode()

		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildCreateTaskAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewCreateTaskInteractor(
				repository.NewTaskSQL(g.db),
				presenter.NewCreateTaskPresenter(),
				g.ctxTimeout,
			)
			act = action.NewCreateTaskAction(uc, g.log, g.validator)
		)
		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildUpdateTaskAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewUpdateTaskInteractor(
				repository.NewTaskSQL(g.db),
				g.ctxTimeout,
			)
			act = action.NewUpdateTaskAction(uc, g.log, g.validator)
		)

		var (
			vars = mux.Vars(req)   // Get path params
			q    = req.URL.Query() // Get query param
		)

		q.Add("task_id", vars["task_id"])
		req.URL.RawQuery = q.Encode()

		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildFindAllTaskAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewFindAllTaskInteractor(
				repository.NewTaskSQL(g.db),
				presenter.NewFindAllTaskPresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAllTaskAction(uc, g.log)
		)
		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

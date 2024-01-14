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

type gorillaMux struct {
	router     *mux.Router
	middleware *negroni.Negroni
	log        logger.Logger
	db         repository.SQL
	validator  validator.Validator
	port       Port
	ctxTimeout time.Duration
}

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

	// タスク
	api.Handle("/tasks", g.buildCreateTaskAction()).Methods(http.MethodPost)
	api.Handle("/tasks/{task_id}", g.buildUpdateTaskAction()).Methods(http.MethodPut)
	api.Handle("/tasks", g.buildFindAllTaskAction()).Methods(http.MethodGet)

	// ヘルスチェック
	api.HandleFunc("/health", action.HealthCheck).Methods(http.MethodGet)
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

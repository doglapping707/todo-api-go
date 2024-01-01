package infrastructure

import (
	"strconv"
	"time"

	"github.com/doglapping707/todo-api-go/adapter/logger"
	"github.com/doglapping707/todo-api-go/adapter/repository"
	"github.com/doglapping707/todo-api-go/adapter/validator"
	"github.com/doglapping707/todo-api-go/infrastructure/database"
	"github.com/doglapping707/todo-api-go/infrastructure/log"
	"github.com/doglapping707/todo-api-go/infrastructure/router"
	"github.com/doglapping707/todo-api-go/infrastructure/validation"
)

// サーバー接続設定
type config struct {
	appName       string
	logger        logger.Logger
	validator     validator.Validator
	dbSQL         repository.SQL
	ctxTimeout    time.Duration
	webServerPort router.Port
	webServer     router.Server
}

// サーバー接続設定を返す
func NewConfig() *config {
	return &config{}
}

// サーバー接続設定に "コンテキストがタイムアウトする時間" をセットし返却する
func (c *config) ContextTimeout(t time.Duration) *config {
	c.ctxTimeout = t
	return c
}

// サーバー接続設定に "アプリケーション名" をセットし返却する
func (c *config) Name(name string) *config {
	c.appName = name
	return c
}

// サーバー接続設定に "ロガー" をセットし返却する
func (c *config) Logger(instance int) *config {
	log, err := log.NewLoggerFactory(instance)
	if err != nil {
		log.Fatalln(err)
	}

	c.logger = log
	c.logger.Infof("Successfully configured log")
	return c
}

// サーバー接続設定に "DBハンドラー" をセットし返却する
func (c *config) DbSQL(instance int) *config {
	db, err := database.NewDatabaseSQLFactory(instance)
	if err != nil {
		c.logger.Fatalln(err, "Could not make a connection to the database")
	}

	c.logger.Infof("Successfully connected to the SQL database")

	c.dbSQL = db
	return c
}

// サーバー接続設定に "バリデーター" をセットし返却する
func (c *config) Validator(instance int) *config {
	v, err := validation.NewValidatorFactory(instance)
	if err != nil {
		c.logger.Fatalln(err)
	}

	c.logger.Infof("Successfully configured validator")

	c.validator = v
	return c
}

// サーバー接続設定に "マルチプレクサー" をセットし返却する
func (c *config) WebServer(instance int) *config {
	s, err := router.NewWebServerFactory(
		instance,
		c.logger,
		c.dbSQL,
		c.validator,
		c.webServerPort,
		c.ctxTimeout,
	)

	if err != nil {
		c.logger.Fatalln(err)
	}

	c.logger.Infof("Successfully configured router server")

	c.webServer = s
	return c
}

// サーバー接続設定に "ポート" をセットし返却する
func (c *config) WebServerPort(port string) *config {
	p, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		c.logger.Fatalln(err)
	}

	c.webServerPort = router.Port(p)
	return c
}

// サーバーを起動する
func (c *config) Start() {
	c.webServer.Listen()
}

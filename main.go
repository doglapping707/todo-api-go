package main

import (
	"os"
	"time"

	"github.com/doglapping707/todo-api-go/infrastructure"
	"github.com/doglapping707/todo-api-go/infrastructure/database"
	"github.com/doglapping707/todo-api-go/infrastructure/log"
	"github.com/doglapping707/todo-api-go/infrastructure/router"
	"github.com/doglapping707/todo-api-go/infrastructure/validation"
)

func main() {
	var app = infrastructure.NewConfig().
		Name(os.Getenv("APP_NAME")).
		ContextTimeout(10 * time.Second).
		Logger(log.InstanceLogrusLogger).
		Validator(validation.InstanceGoPlayground).
		DbSQL(database.InstancePostgres).
		DbNoSQL(database.InstanceMongoDB)

	app.WebServerPort(os.Getenv("APP_PORT")).
		WebServer(router.InstanceGorillaMux).
		Start()
}

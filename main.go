package main

import (
	"fmt"
	"os"
	"time"

	"github.com/doglapping707/todo-api-go/infrastructure"
	"github.com/doglapping707/todo-api-go/infrastructure/log"
	"github.com/doglapping707/todo-api-go/infrastructure/validation"
)

func main() {
	app := infrastructure.NewConfig().
	Name(os.Getenv("APP_NAME")).
	ContextTimeout(10 * time.Second).
	Logger(log.InstanceLogrusLogger).
	Validator(validation.InstanceGoPlayground)
	fmt.Println(app)
}

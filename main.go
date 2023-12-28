package main

import (
	"fmt"
	"os"
	"time"

	"github.com/doglapping707/todo-api-go/infrastructure"
	"github.com/doglapping707/todo-api-go/infrastructure/log"
)

func main() {
	app := infrastructure.NewConfig().
	Name(os.Getenv("APP_NAME")).
	ContextTimeout(10 * time.Second).
	Logger(log.InstanceLogrusLogger)
	fmt.Println(app)
}

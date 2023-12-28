package infrastructure

import (
	"time"

	"github.com/doglapping707/todo-api-go/adapter/logger"
	"github.com/doglapping707/todo-api-go/infrastructure/log"
)

// server connection config
type config struct {
	appName    string        // application
	logger     logger.Logger // logger
	ctxTimeout time.Duration //
}

// init server connection config
func NewConfig() *config {
	return &config{}
}

// 
func (c *config) ContextTimeout(t time.Duration) *config {
	c.ctxTimeout = t
	return c
}

// set application name
func (c *config) Name(name string) *config {
	c.appName = name
	return c
}

// set logger
func (c *config) Logger(instance int) *config {
	log, err := log.NewLoggerFactory(instance)
	if err != nil {
		log.Fatalln(err)
	}

	c.logger = log
	c.logger.Infof("Successfully configured log")
	return c
}

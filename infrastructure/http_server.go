package infrastructure

import (
	"time"

	"github.com/doglapping707/todo-api-go/adapter/logger"
	"github.com/doglapping707/todo-api-go/adapter/validator"
	"github.com/doglapping707/todo-api-go/infrastructure/log"
	"github.com/doglapping707/todo-api-go/infrastructure/validation"
)

type config struct {
	appName    string
	ctxTimeout time.Duration
	logger     logger.Logger
	validator  validator.Validator
}

func NewConfig() *config {
	return &config{}
}

func (c *config) Name(name string) *config {
	c.appName = name
	return c
}

func (c *config) ContextTimeout(t time.Duration) *config {
	c.ctxTimeout = t
	return c
}

func (c *config) Logger(instance int) *config {
	log, err := log.NewLoggerFactory(instance)
	if err != nil {
		log.Fatalln(err)
	}

	c.logger = log
	c.logger.Infof("Successfully configured log")
	return c
}

func (c *config) Validator(instance int) *config {
	v, err := validation.NewValidatorFactory(instance)
	if err != nil {
		c.logger.Fatalln(err)
	}

	c.logger.Infof("Successfully configured validator")

	c.validator = v
	return c
}

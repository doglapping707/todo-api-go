package router

import (
	"errors"
	"time"

	"github.com/doglapping707/todo-api-go/adapter/logger"
	"github.com/doglapping707/todo-api-go/adapter/repository"
	"github.com/doglapping707/todo-api-go/adapter/validator"
)

type Server interface {
	Listen()
}

type Port int64

var (
	errInvalidWebServerInstance = errors.New("invalid router server instance")
)

const (
	InstanceGorillaMux int = iota
	InstanceGin
)

func NewWebServerFactory(
	instance int,
	log logger.Logger,
	dbSQL repository.SQL,
	dbNoSQL repository.NoSQL,
	validator validator.Validator,
	port Port,
	ctxTimeout time.Duration,
) (Server, error) {
	switch instance {
	case InstanceGorillaMux:
		return newGorillaMux(log, dbSQL, validator, port, ctxTimeout), nil
	case InstanceGin:
		return newGinServer(log, dbNoSQL, validator, port, ctxTimeout), nil
	default:
		return nil, errInvalidWebServerInstance
	}
}

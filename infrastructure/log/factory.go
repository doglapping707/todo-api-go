package log

import (
	"errors"

	"github.com/doglapping707/todo-api-go/adapter/logger"
)

const (
	InstanceLogrusLogger int = iota
)

var (
	errInvalidLoggerInstance = errors.New("invalid log instance")
)

// 生成されたロガーを返却する
func NewLoggerFactory(instance int) (logger.Logger, error) {
	switch instance {
	case InstanceLogrusLogger:
		return NewLogrusLogger(), nil
	default:
		return nil, errInvalidLoggerInstance
	}
}

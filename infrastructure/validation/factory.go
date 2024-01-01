package validation

import (
	"errors"

	"github.com/doglapping707/todo-api-go/adapter/validator"
)

var (
	errInvalidValidatorInstance = errors.New("invalid validator instance")
)

const (
	InstanceGoPlayground int = iota
)

// 生成されたバリデーターを返却する
func NewValidatorFactory(instance int) (validator.Validator, error) {
	switch instance {
	case InstanceGoPlayground:
		return NewGoPlayground()
	default:
		return nil, errInvalidValidatorInstance
	}
}

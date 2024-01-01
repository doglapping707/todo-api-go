package database

import (
	"errors"

	"github.com/doglapping707/todo-api-go/adapter/repository"
)

var (
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

const (
	InstancePostgres int = iota
)

// 生成されたDBハンドラーを返却する
func NewDatabaseSQLFactory(instance int) (repository.SQL, error) {
	switch instance {
	case InstancePostgres:
		return NewPostgresHandler(newConfigPostgres())
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
}

package repository

import (
	"context"

	"github.com/doglapping707/todo-api-go/domain"
	"github.com/pkg/errors"
)

type TaskSQL struct {
	db SQL
}

func NewTaskSQL(db SQL) TaskSQL {
	return TaskSQL{
		db: db,
	}	
}

func (t TaskSQL) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	var query = "INSERT INTO tasks (title) VALUES ($1)"

	// sqlを実行する
	if err := t.db.ExecuteContext(
		ctx,
		query,
		task.Title,
	); err != nil {
		return domain.Task{}, errors.Wrap(err, "error creating task")
	}

	return task, nil
}

func (t TaskSQL) Update(ctx context.Context, task domain.Task) (domain.Task, error) {
	query := "UPDATE tasks SET title = $1 WHERE id = $2"

	// sqlを実行する
	if err := t.db.ExecuteContext(
		ctx,
		query,
		task.Title,
		task.ID,
	); err != nil {
		return domain.Task{}, errors.Wrap(err, "error updating task")
	}

	return task, nil
}

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

func (t TaskSQL) Update(ctx context.Context, task domain.Task) error {
	var query = "UPDATE tasks SET title = $1 WHERE id = $2"

	// sqlを実行する
	if err := t.db.ExecuteContext(
		ctx,
		query,
		task.Title,
		task.ID,
	); err != nil {
		return errors.Wrap(err, "error updating task")
	}

	return nil
}

func (t TaskSQL) FindAll(ctx context.Context) ([]domain.Task, error) {
	var query = "SELECT id, title FROM tasks"

	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return []domain.Task{}, errors.Wrap(err, "error listing tasks")
	}
	
	var tasks = make([]domain.Task, 0)
	for rows.Next() {
		var (
			ID        uint64
			title     string
		)

		if err = rows.Scan(&ID, &title); err != nil {
			return []domain.Task{}, errors.Wrap(err, "error listing tasks")
		}

		tasks = append(tasks, domain.Task{
			ID:    ID,
			Title: title,
		})
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return []domain.Task{}, err
	}

	return tasks, nil
}

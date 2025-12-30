package task

import (
	"context"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s Store) CreateTask(ctx context.Context, input CreateTaskInput) (Task, error) {
	const insertStatement = `
		INSERT INTO tasks (user_id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, false, current_timestamp, current_timestamp)
		RETURNING id, completed, created_at, updated_at
	`

	task := Task{
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
	}

	err := s.db.
		QueryRowContext(
			ctx,
			insertStatement,
			input.UserID,
			input.Title,
			input.Description,
		).
		Scan(&task.ID, &task.Completed, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

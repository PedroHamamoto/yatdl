package task

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func (s Store) FindByID(ctx context.Context, id uint64) (*Task, error) {
	const selectStatement = `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var task Task
	err := s.db.
		QueryRowContext(ctx, selectStatement, id).
		Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if task.ID == 0 {
		return nil, sql.ErrNoRows
	}

	return &task, nil
}

func (s Store) UpdateTaskByID(ctx context.Context, input *UpdateTaskInput) error {
	var setClauses []string
	var args []any
	argPos := 1

	if input.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argPos))
		args = append(args, *input.Title)
		argPos++
	}

	if input.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *input.Description)
		argPos++
	}

	if input.Completed != nil {
		setClauses = append(setClauses, fmt.Sprintf("completed = $%d", argPos))
		args = append(args, *input.Completed)
		argPos++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	updateStatement := fmt.Sprintf(`
		UPDATE tasks
		SET %s,
		    updated_at = current_timestamp
		WHERE id = $%d
	`, strings.Join(setClauses, ", "), argPos)

	args = append(args, input.ID)

	_, err := s.db.ExecContext(ctx, updateStatement, args...)
	return err
}

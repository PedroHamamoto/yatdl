package task

import (
	"context"
	"errors"
)

var (
	ErrCannotUpdateTaskFromAnotherUser = errors.New("cannot update task from another user")
)

type Service struct {
	store Store
}

type CreateTaskInput struct {
	UserID      uint64
	Title       string
	Description string
}

type UpdateTaskInput struct {
	ID          uint64
	UserID      uint64
	Title       *string
	Description *string
	Completed   *bool
}

func NewService(store *Store) *Service {
	return &Service{store: *store}
}

func (s *Service) CreateTask(ctx context.Context, input CreateTaskInput) (Task, error) {
	task, err := s.store.CreateTask(ctx, input)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, input UpdateTaskInput) error {
	task, err := s.store.FindByID(ctx, input.ID)

	if err != nil {
		return err
	}

	// TODO include additional validations
	if task.UserID != input.UserID {
		return ErrCannotUpdateTaskFromAnotherUser
	}

	err = s.store.UpdateTaskByID(ctx, &input)

	return err
}

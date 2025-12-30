package task

import (
	"context"
)

type Service struct {
	store Store
}

type CreateTaskInput struct {
	UserID      uint64
	Title       string
	Description string
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

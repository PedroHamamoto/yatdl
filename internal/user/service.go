package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, email string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.store.CreateUser(ctx, email, hash)
	if err != nil {
		return err
	}

	// TODO Handle errors

	return nil
}

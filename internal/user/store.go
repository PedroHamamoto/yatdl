package user

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

func (s *Store) CreateUser(ctx context.Context, email string, passwordHash []byte) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO users (email, password_hash) VALUES ($1, $2)",
		email, passwordHash)

	return err
}

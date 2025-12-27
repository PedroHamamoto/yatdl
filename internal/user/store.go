package user

import (
	"context"
	"database/sql"
	"errors"
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

func (s *Store) FindByEmail(ctx context.Context, email string) (*User, error) {
	const query = "SELECT id, email, password_hash FROM users WHERE email = $1"

	var user User
	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil

}

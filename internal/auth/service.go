package auth

import (
	"context"
	"errors"
	"yatdl/internal/user"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	jwtSecret []byte
	userStore *user.Store
}

func NewService(jwtSecret string, userStore *user.Store) *Service {
	return &Service{
		jwtSecret: []byte(jwtSecret),
		userStore: userStore,
	}
}

func (s *Service) Login(ctx context.Context, email string, password string) (*LoginResponse, error) {
	user, err := s.userStore.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, expiresIn, err := GenerateJWT(int64(user.ID), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   "Bearer",
	}, nil

}

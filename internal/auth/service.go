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
	jwt       *Jwt
	userStore *user.Store
}

func NewService(jwt *Jwt, userStore *user.Store) *Service {
	return &Service{
		jwt:       jwt,
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

	token, expiresIn, err := s.jwt.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   "Bearer",
	}, nil

}

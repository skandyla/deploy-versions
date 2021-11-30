package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/skandyla/deploy-versions/internal/domain"
)

// PasswordHasher provides hashing logic to securely store passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type Users struct {
	repo   UsersRepository
	hasher PasswordHasher
}

func NewUsers(repo UsersRepository, hasher PasswordHasher) *Users {
	return &Users{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *Users) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *Users) SignIn(ctx context.Context, inp domain.SignInInput) (err error) {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return
	}

	_, err = s.repo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return
}

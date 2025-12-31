package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Uranury/WorkoutTracker/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strings"
)

type Service interface {
	Create(ctx context.Context, request SignUpRequest) (*User, error)
	ValidateCredentials(ctx context.Context, username, password string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, id int64, updates UpdateUserInput) (*User, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) Create(ctx context.Context, request SignUpRequest) (*User, error) {
	existing, _ := s.repo.GetByEmail(ctx, request.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	existing, _ = s.repo.GetByUsername(ctx, request.Username)
	if existing != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Age:      request.Age,
		Gender:   request.Gender,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *service) ValidateCredentials(ctx context.Context, username, password string) (*User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id int64, updates UpdateUserInput) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if updates.Username != nil {
		user.Username = *updates.Username
	}
	if updates.Email != nil {
		user.Email = *updates.Email
	}
	if updates.Age != nil {
		user.Age = *updates.Age
	}
	if updates.Gender != nil {
		user.Gender = strings.ToLower(*updates.Gender)
	}
	if updates.Weight != nil {
		user.Weight = *updates.Weight
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

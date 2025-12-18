package user

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type Service interface {
	Create(ctx context.Context, request SignUpRequest) (*User, error)
	ValidateCredentials(ctx context.Context, username, password string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, id int64, updates UpdateUserInput) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request SignUpRequest) (*User, error) {
	// Validation
	if err := validateEmail(request.Email); err != nil {
		return nil, err
	}
	if err := validatePassword(request.Password); err != nil {
		return nil, err
	}
	if err := validateUsername(request.Username); err != nil {
		return nil, err
	}

	// Check if user already exists
	existing, _ := s.repo.GetByEmail(ctx, request.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	existing, _ = s.repo.GetByUsername(ctx, request.Username)
	if existing != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
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
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

type UpdateUserInput struct {
	Username *string
	Email    *string
	Age      *int
	Gender   *string
}

func (s *service) Update(ctx context.Context, id int64, updates UpdateUserInput) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if updates.Username != nil {
		if err := validateUsername(*updates.Username); err != nil {
			return nil, err
		}
		user.Username = *updates.Username
	}
	if updates.Email != nil {
		if err := validateEmail(*updates.Email); err != nil {
			return nil, err
		}
		user.Email = *updates.Email
	}
	if updates.Age != nil {
		user.Age = *updates.Age
	}
	if updates.Gender != nil {
		user.Gender = *updates.Gender
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Validation helpers
func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	return nil
}

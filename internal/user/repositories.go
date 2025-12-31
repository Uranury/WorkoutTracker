package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewRepository(db *sqlx.DB, logger *slog.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	query := `
        INSERT INTO users (username, email, age, gender, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		user.Username, user.Email, user.Age, user.Gender, user.Password,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE email = $1"
	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE username = $1"
	err := r.db.GetContext(ctx, &user, query, username)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	query := `
        UPDATE users 
        SET username = $1, email = $2, age = $3, gender = $4, updated_at = NOW(), weight = $5
        WHERE id = $6`

	_, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.Age, user.Gender, user.Weight, user.ID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

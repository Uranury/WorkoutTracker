package auth

import (
	"context"
	"github.com/Uranury/WorkoutTracker/pkg/database"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	FindByHashForUpdate(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
	DeleteExpired(ctx context.Context) error
}

type repository struct {
	executor database.Executor
}

func NewRepository(db *sqlx.DB) RefreshTokenRepository {
	return &repository{
		executor: db,
	}
}

func NewRepositoryFromTx(tx *sqlx.Tx) RefreshTokenRepository {
	return &repository{
		executor: tx,
	}
}

func (repo *repository) Save(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, user_agent, ip)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return repo.executor.QueryRowxContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt, token.UserAgent, token.IP).Scan(&token.ID)
}

func (repo *repository) FindByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	query := `
		SELECT * FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()
	`
	token := &RefreshToken{}
	err := repo.executor.QueryRowxContext(ctx, query, tokenHash).StructScan(token)
	return token, err
}

func (repo *repository) FindByHashForUpdate(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	query := `
		SELECT * FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()
		FOR UPDATE 
    `
	token := &RefreshToken{}
	err := repo.executor.QueryRowxContext(ctx, query, tokenHash).StructScan(token)
	return token, err
}

func (repo *repository) RevokeByHash(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1`
	_, err := repo.executor.ExecContext(ctx, query, tokenHash)
	return err
}

func (repo *repository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := repo.executor.ExecContext(ctx, query)
	return err
}

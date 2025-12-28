package auth

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
	DeleteExpired(ctx context.Context) error
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
}

type repository struct {
	Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) RefreshTokenRepository {
	return &repository{
		Db: db,
	}
}

func (repo *repository) Save(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, user_agent, ip)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return repo.Db.QueryRowxContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt, token.UserAgent, token.IP).Scan(&token.ID)
}

func (repo *repository) FindByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at, user_agent, ip
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()
	`
	token := &RefreshToken{}
	err := repo.Db.QueryRowxContext(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
		&token.UserAgent,
		&token.IP,
	)
	return token, err
}

func (repo *repository) FindByHashTx(tx *sqlx.Tx, tokenHash string) (*RefreshToken, error) {
	return nil, nil
}

func (repo *repository) RevokeByHash(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1`
	_, err := repo.Db.ExecContext(ctx, query, tokenHash)
	return err
}

func (repo *repository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := repo.Db.ExecContext(ctx, query)
	return err
}

func (repo *repository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := repo.Db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

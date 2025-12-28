package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log/slog"
	"time"
)

type Service interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(tokenString string) (*Claims, error)

	GenerateRefreshToken(ctx context.Context, userID int64, userAgent string, ip string) (string, error)
	ValidateRefreshToken(ctx context.Context, tokenString string) (*RefreshToken, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error)
}

type auth struct {
	jwtKey []byte
	logger *slog.Logger
	repo   RefreshTokenRepository
}

func NewAuth(secret string, repo RefreshTokenRepository, logger *slog.Logger) Service {
	return &auth{jwtKey: []byte(secret), repo: repo, logger: logger}
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *auth) GenerateToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

func (s *auth) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtKey, nil
	})
	if err != nil || !token.Valid {
		if err != nil {
			return nil, fmt.Errorf("failed to parse token: %w", err)
		}
		return nil, fmt.Errorf("token is not valid")
	}
	return claims, nil
}

func (s *auth) GenerateRefreshToken(ctx context.Context, userID int64, userAgent string, ip string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(b)
	tokenHash := hashToken(token)

	refreshToken := RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		UserAgent: userAgent,
		IP:        ip,
	}
	if err := s.repo.Save(ctx, &refreshToken); err != nil {
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}
	return token, nil
}

func (s *auth) ValidateRefreshToken(ctx context.Context, tokenString string) (*RefreshToken, error) {
	tokenHash := hashToken(tokenString)

	token, err := s.repo.FindByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}
	if token.RevokedAt != nil {
		return nil, fmt.Errorf("refresh token is revoked")
	}
	if time.Now().After(token.ExpiresAt) {
		return nil, fmt.Errorf("refresh token is expired")
	}

	return token, nil
}

func (s *auth) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction: %w", err)
		return "", "", fmt.Errorf("failed to start transaction")
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			s.logger.Error("failed to rollback transaction: %w", err)
		}
	}()

	rt, err := s.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}
	accessToken, err := s.GenerateToken(rt.UserID)
	if err != nil {
		return "", "", err
	}

	if err := s.RevokeRefreshToken(ctx, refreshToken); err != nil {
		return "", "", fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	newRefresh, err := s.GenerateRefreshToken(ctx, rt.UserID, rt.UserAgent, rt.IP)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", "", fmt.Errorf("failed to commit transaction: %w", err)
	}
	return accessToken, newRefresh, nil
}

func (s *auth) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	if err := s.repo.RevokeByHash(ctx, tokenHash); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

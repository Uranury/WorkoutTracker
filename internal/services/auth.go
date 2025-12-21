package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Auth interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type auth struct {
	jwtKey []byte
}

func NewAuth(secret string) Auth {
	return &auth{jwtKey: []byte(secret)}
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

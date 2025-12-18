package middleware

import (
	"errors"
	"github.com/Uranury/WorkoutTracker/internal/services"
	"github.com/Uranury/WorkoutTracker/pkg/apperrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Auth struct {
	authService *services.Auth
}

func NewAuth(authService *services.Auth) *Auth {
	return &Auth{authService: authService}
}

type contextKey string

const UserIDKey contextKey = "user_id"

func (m *Auth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			apperrors.GenHTTPError(c, http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := m.authService.VerifyJWT(tokenString)
		if err != nil {
			apperrors.GenHTTPError(c, http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (int64, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, errors.New("user id not found in context")
	}

	id, ok := userID.(int64)
	if !ok {
		return 0, errors.New("user id has invalid type")
	}

	return id, nil
}

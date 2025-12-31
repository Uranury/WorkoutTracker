package http_server

import (
	"github.com/gin-gonic/gin"
)

func (h *HTTPServer) setupRoutes() {
	// Public routes
	authMiddleware := h.app.AuthMiddleware()
	userHandler := h.app.UserHandler()

	h.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
		return
	})
	auth := h.router.Group("/auth")
	auth.POST("/signup", userHandler.SignUp)
	auth.POST("/login", userHandler.Login)
	auth.POST("/refresh", userHandler.RefreshToken)

	// Protected routes
	api := h.router.Group("/api")
	api.Use(authMiddleware.JWTAuth())
	{
		api.GET("/users/me", userHandler.GetProfile)
		api.PUT("/users/me", userHandler.UpdateProfile)
	}
}

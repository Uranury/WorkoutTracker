package http_server

import "github.com/gin-gonic/gin"

func (h *HTTPServer) setupRoutes() {
	// Public routes
	h.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
		return
	})
	auth := h.router.Group("/auth")
	auth.POST("/signup", h.userHandler.SignUp)
	auth.POST("/login", h.userHandler.Login)

	// Protected routes
	api := h.router.Group("/api")
	api.Use(h.authMiddleware.JWTAuth())
	{
		api.GET("/users/me", h.userHandler.GetProfile)
		api.PUT("/users/me", h.userHandler.UpdateProfile)
	}
}

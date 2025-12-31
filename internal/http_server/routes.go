package http_server

import (
	_ "github.com/Uranury/WorkoutTracker/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (h *HTTPServer) setupRoutes() {
	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Public routes
	authMiddleware := h.app.AuthMiddleware()

	h.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
		return
	})
	auth := h.router.Group("/auth")
	auth.POST("/signup", h.app.UserHandler().SignUp)
	auth.POST("/login", h.app.UserHandler().Login)
	auth.POST("/logout", h.app.UserHandler().Logout)
	auth.POST("/refresh", h.app.UserHandler().RefreshToken)

	// Protected routes
	api := h.router.Group("/api")
	api.Use(authMiddleware.JWTAuth())
	{
		users := api.Group("/users")
		users.GET("/me", h.app.UserHandler().GetProfile)
		users.PATCH("/me", h.app.UserHandler().UpdateProfile)
		users.GET("/:id", h.app.UserHandler().GetUserByID)
	}
}

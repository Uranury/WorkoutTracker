package user

import (
	"github.com/Uranury/WorkoutTracker/internal/middleware"
	"github.com/Uranury/WorkoutTracker/pkg/apperrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}

// SignUp registers a new user
func (h *Handler) SignUp(c *gin.Context) {
	request := SignUpRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	user, err := h.service.Create(c.Request.Context(), request)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetProfile returns the current user's profile
func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the current user's profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	var req UpdateUserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user, err := h.service.Update(c.Request.Context(), userID, req)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByID is for admin use cases (optional)
func (h *Handler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to get user", nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

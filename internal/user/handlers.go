package user

import (
	"errors"
	"github.com/Uranury/WorkoutTracker/internal/auth"
	"github.com/Uranury/WorkoutTracker/internal/middleware"
	"github.com/Uranury/WorkoutTracker/pkg/apperrors"
	"github.com/Uranury/WorkoutTracker/pkg/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service     Service
	authService auth.Service
}

func NewHandler(service Service, authService auth.Service) *Handler {
	return &Handler{service: service, authService: authService}
}

type SignUpRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=32"`
	Password string `json:"password" binding:"required" validate:"required,min=8"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Age      int    `json:"age" binding:"required" validate:"required,gt=0"`
	Gender   string `json:"gender" binding:"required" validate:"required,oneof=male female"`
}

// SignUp registers a new user
func (h *Handler) SignUp(c *gin.Context) {
	req, ok := validation.BindAndValidate[SignUpRequest](c)
	if !ok {
		return
	}
	user, err := h.service.Create(c.Request.Context(), *req)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, user)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=32"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (h *Handler) Login(c *gin.Context) {
	req, ok := validation.BindAndValidate[LoginRequest](c)
	if !ok {
		return
	}

	user, err := h.service.ValidateCredentials(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.GenHTTPError(c, http.StatusNotFound, err.Error(), nil)
		} else {
			apperrors.GenHTTPError(c, http.StatusUnauthorized, err.Error(), nil)
		}
		return
	}

	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: token, User: *user})
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

type UpdateUserInput struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=32"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Age      *int    `json:"age" validate:"omitempty,gt=0"`
	Gender   *string `json:"gender" validate:"omitempty,oneof=male female"`
}

// UpdateProfile updates the current user's profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	req, ok := validation.BindAndValidate[UpdateUserInput](c)
	if !ok {
		return
	}

	user, err := h.service.Update(c.Request.Context(), userID, *req)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

type IntIDPathParam struct {
	ID int64 `uri:"id" binding:"required" validate:"required,gt=0"`
}

// GetUserByID is for admin use cases (optional)
func (h *Handler) GetUserByID(c *gin.Context) {
	idParam, ok := validation.BindAndValidateURI[IntIDPathParam](c)
	if !ok {
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), idParam.ID)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to get user", nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

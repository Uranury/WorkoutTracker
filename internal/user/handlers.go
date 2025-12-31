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
// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "Sign up payload"
// @Success 201 {object} User
// @Failure 400 {object} apperrors.HTTPError
// @Failure 500 {object} apperrors.HTTPError
// @Router /auth/signup [post]
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
	c.JSON(http.StatusCreated, user)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=32"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

// Login authenticates a user
// @Summary Login user
// @Description Validates credentials and returns access token + sets refresh token cookie
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login payload"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} apperrors.HTTPError
// @Failure 404 {object} apperrors.HTTPError
// @Failure 500 {object} apperrors.HTTPError
// @Router /auth/login [post]
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

	accessToken, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(c.Request.Context(), user.ID, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(auth.RefreshTokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	c.JSON(http.StatusOK, LoginResponse{AccessToken: accessToken, User: *user})
}

// Logout revokes the user's refresh token
// @Summary Logout user
// @Description Revokes the refresh token and clears the cookie
// @Tags auth
// @Produce json
// @Success 200 {string} string "Logged out successfully"
// @Failure 500 {object} apperrors.HTTPError
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.Status(http.StatusOK)
		return
	}

	if err := h.authService.RevokeRefreshToken(c.Request.Context(), refreshToken); err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to log out", nil)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	c.Status(http.StatusOK)
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// RefreshToken refreshes access token
// @Summary Refresh access token
// @Description Rotates refresh token and returns new access token
// @Tags auth
// @Produce json
// @Success 200 {object} AccessTokenResponse
// @Failure 401 {object} apperrors.HTTPError
// @Failure 500 {object} apperrors.HTTPError
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	accessToken, newRefresh, err := h.authService.RefreshAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to refresh access token", nil)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		Path:     "/",
		MaxAge:   int(auth.RefreshTokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	c.JSON(http.StatusOK, AccessTokenResponse{AccessToken: accessToken})
}

// GetProfile returns current user's profile
// @Summary Get current user profile
// @Description Returns the authenticated user's profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} User
// @Failure 401 {object} apperrors.HTTPError
// @Failure 500 {object} apperrors.HTTPError
// @Router /api/users/me [get]
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
	Username *string  `json:"username" validate:"omitempty,min=3,max=32"`
	Email    *string  `json:"email" validate:"omitempty,email"`
	Age      *int     `json:"age" validate:"omitempty,gt=0"`
	Gender   *string  `json:"gender" validate:"omitempty,oneof=male female"`
	Weight   *float64 `json:"weight" validate:"omitempty,gt=0"`
}

// UpdateProfile updates current user's profile
// @Summary Update user profile
// @Description Updates fields of the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateUserInput true "Update user payload"
// @Success 200 {object} User
// @Failure 400 {object} apperrors.HTTPError
// @Failure 401 {object} apperrors.HTTPError
// @Failure 500 {object} apperrors.HTTPError
// @Router /api/users/me [patch]
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

// GetUserByID returns user by ID
// @Summary Get user by ID
// @Description Admin endpoint to fetch user by ID
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} apperrors.HTTPError
// @Failure 401 {object} apperrors.HTTPError
// @Failure 500 {object} apperrors.HTTPError
// @Router /api/users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) { // todo: roles not implemented yet
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

package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"mylittleprice/internal/app"
	"mylittleprice/internal/middleware"
	"mylittleprice/internal/models"
	"mylittleprice/internal/services"
)

type AuthHandler struct {
	container *app.Container
}

func NewAuthHandler(c *app.Container) *AuthHandler {
	return &AuthHandler{container: c}
}

// Signup handles user registration
// POST /api/auth/signup
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req models.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Email and password are required",
		})
	}

	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Password must be at least 8 characters",
		})
	}

	// Create user
	authResp, err := h.container.AuthService.Signup(&req)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
				Error:   "user_exists",
				Message: "User with this email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(authResp)
}

// GoogleLogin handles Google OAuth authentication
// POST /api/auth/google
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	var req models.GoogleAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.IDToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID token is required",
		})
	}

	// Authenticate user with Google
	authResp, err := h.container.AuthService.GoogleLogin(req.IDToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "invalid_token",
			Message: "Failed to verify Google token",
		})
	}

	return c.JSON(authResp)
}

// Login handles user authentication
// POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Email and password are required",
		})
	}

	// Authenticate user
	authResp, err := h.container.AuthService.Login(&req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) || errors.Is(err, services.ErrInvalidPassword) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to authenticate user",
		})
	}

	return c.JSON(authResp)
}

// RefreshToken handles access token refresh
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Refresh token is required",
		})
	}

	authResp, err := h.container.AuthService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "invalid_token",
				Message: "Invalid or expired refresh token",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to refresh token",
		})
	}

	return c.JSON(authResp)
}

// Logout handles user logout (revokes refresh token)
// POST /api/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Refresh token is required",
		})
	}

	if err := h.container.AuthService.Logout(req.RefreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to logout",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// GetMe returns current user info (requires authentication)
// GET /api/auth/me
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
	}

	user, err := h.container.AuthService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
		})
	}

	return c.JSON(models.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Picture:   user.Picture,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
	})
}

// ClaimSessions links anonymous chat sessions to the authenticated user
// POST /api/auth/claim-sessions
func (h *AuthHandler) ClaimSessions(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
	}

	var req models.ClaimSessionsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if len(req.SessionIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "At least one session ID is required",
		})
	}

	if err := h.container.AuthService.ClaimSessions(userID, req.SessionIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to claim sessions",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Sessions claimed successfully",
		"claimed": len(req.SessionIDs),
	})
}

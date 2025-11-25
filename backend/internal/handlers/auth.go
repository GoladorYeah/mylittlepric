package handlers

import (
	"errors"
	"fmt"

	"mylittleprice/internal/container"
	"mylittleprice/internal/middleware"
	"mylittleprice/internal/models"
	"mylittleprice/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	container *container.Container
}

func NewAuthHandler(c *container.Container) *AuthHandler {
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

// ChangePassword handles password change for authenticated users
// POST /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
	}

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Current password and new password are required",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "New password must be at least 8 characters",
		})
	}

	// Change password
	if err := h.container.AuthService.ChangePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrPasswordNotSet) {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "password_not_set",
				Message: "Cannot change password for OAuth users",
			})
		}
		if errors.Is(err, services.ErrInvalidPassword) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "invalid_password",
				Message: "Current password is incorrect",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to change password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

// RequestPasswordReset generates a password reset token
// POST /api/auth/request-password-reset
func (h *AuthHandler) RequestPasswordReset(c *fiber.Ctx) error {
	var req models.RequestPasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Email is required",
		})
	}

	// Generate reset token
	resetToken, err := h.container.AuthService.RequestPasswordReset(req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate reset token",
		})
	}

	// If token was generated (not empty), try to send email
	if resetToken != "" {
		// Send reset email
		if err := h.container.EmailService.SendPasswordResetEmail(req.Email, resetToken); err != nil {
			// Log error but don't fail the request - token is still valid
			fmt.Printf("⚠️ Failed to send password reset email: %v\n", err)
			// For testing, return the token if email fails
			return c.JSON(fiber.Map{
				"message": "Password reset token generated (email failed, using fallback)",
				"token":   resetToken, // Fallback for testing
			})
		}
	}

	// Success response - don't reveal if user exists
	return c.JSON(fiber.Map{
		"message": "If an account exists with this email, a password reset link has been sent",
	})
}

// ResetPassword resets the user's password using a reset token
// POST /api/auth/reset-password
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req models.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate input
	if req.Token == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Token and new password are required",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "New password must be at least 8 characters",
		})
	}

	// Reset password
	if err := h.container.AuthService.ResetPassword(req.Token, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrResetTokenNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error:   "token_not_found",
				Message: "Reset token not found",
			})
		}
		if errors.Is(err, services.ErrTokenExpired) {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "token_expired",
				Message: "Reset token has expired",
			})
		}
		if errors.Is(err, services.ErrTokenAlreadyUsed) {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "token_used",
				Message: "Reset token has already been used",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to reset password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

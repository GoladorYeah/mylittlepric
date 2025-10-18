package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

var (
	ErrUserExists       = errors.New("user with this email already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidToken     = errors.New("invalid or expired token")
)

type AuthService struct {
	redis      *redis.Client
	jwtService *utils.JWTService
	ctx        context.Context
}

func NewAuthService(redis *redis.Client, jwtService *utils.JWTService) *AuthService {
	return &AuthService{
		redis:      redis,
		jwtService: jwtService,
		ctx:        context.Background(),
	}
}

// Signup creates a new user account
func (s *AuthService) Signup(req *models.SignupRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	exists, err := s.userExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.saveUser(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Generate tokens
	return s.generateAuthResponse(user)
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.getUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.saveUser(user); err != nil {
		// Don't fail login if we can't update last login time
		fmt.Printf("Warning: failed to update last login time: %v\n", err)
	}

	// Generate tokens
	return s.generateAuthResponse(user)
}

// RefreshAccessToken generates a new access token using a refresh token
func (s *AuthService) RefreshAccessToken(refreshToken string) (*models.AuthResponse, error) {
	// Validate refresh token
	tokenHash := s.hashToken(refreshToken)
	refreshTokenData, err := s.getRefreshToken(tokenHash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Check if token is revoked or expired
	if refreshTokenData.RevokedAt != nil {
		return nil, ErrInvalidToken
	}
	if time.Now().After(refreshTokenData.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.getUserByID(refreshTokenData.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate new access token (keep same refresh token)
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         s.toUserInfo(user),
		ExpiresIn:    int64(s.jwtService.GetAccessTTL().Seconds()),
	}, nil
}

// Logout revokes a refresh token
func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := s.hashToken(refreshToken)
	return s.revokeRefreshToken(tokenHash)
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.getUserByID(userID)
}

// ClaimSessions links anonymous sessions to a user account
func (s *AuthService) ClaimSessions(userID uuid.UUID, sessionIDs []string) error {
	for _, sessionID := range sessionIDs {
		key := fmt.Sprintf("session:%s", sessionID)

		// Add user_id to session data
		if err := s.redis.HSet(s.ctx, key, "user_id", userID.String()).Err(); err != nil {
			return fmt.Errorf("failed to claim session %s: %w", sessionID, err)
		}
	}
	return nil
}

// ==================== Private Helper Methods ====================

func (s *AuthService) generateAuthResponse(user *models.User) (*models.AuthResponse, error) {
	// Generate access token
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token
	tokenHash := s.hashToken(refreshToken)
	refreshTokenData := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTTL()),
		CreatedAt: time.Now(),
	}

	if err := s.saveRefreshToken(refreshTokenData); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         s.toUserInfo(user),
		ExpiresIn:    int64(s.jwtService.GetAccessTTL().Seconds()),
	}, nil
}

func (s *AuthService) userExists(email string) (bool, error) {
	key := fmt.Sprintf("user:email:%s", email)
	exists, err := s.redis.Exists(s.ctx, key).Result()
	return exists > 0, err
}

func (s *AuthService) saveUser(user *models.User) error {
	// Save user data by ID
	userKey := fmt.Sprintf("user:id:%s", user.ID.String())
	userData := map[string]interface{}{
		"id":            user.ID.String(),
		"email":         user.Email,
		"password_hash": user.PasswordHash,
		"full_name":     user.FullName,
		"created_at":    user.CreatedAt.Format(time.RFC3339),
		"updated_at":    user.UpdatedAt.Format(time.RFC3339),
	}

	if user.LastLoginAt != nil {
		userData["last_login_at"] = user.LastLoginAt.Format(time.RFC3339)
	}

	if err := s.redis.HSet(s.ctx, userKey, userData).Err(); err != nil {
		return err
	}

	// Create email -> ID mapping
	emailKey := fmt.Sprintf("user:email:%s", user.Email)
	return s.redis.Set(s.ctx, emailKey, user.ID.String(), 0).Err()
}

func (s *AuthService) getUserByEmail(email string) (*models.User, error) {
	emailKey := fmt.Sprintf("user:email:%s", email)
	userIDStr, err := s.redis.Get(s.ctx, emailKey).Result()
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return s.getUserByID(userID)
}

func (s *AuthService) getUserByID(userID uuid.UUID) (*models.User, error) {
	userKey := fmt.Sprintf("user:id:%s", userID.String())
	userData, err := s.redis.HGetAll(s.ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	if len(userData) == 0 {
		return nil, redis.Nil
	}

	user := &models.User{}
	user.ID = userID
	user.Email = userData["email"]
	user.PasswordHash = userData["password_hash"]
	user.FullName = userData["full_name"]

	if createdAt, err := time.Parse(time.RFC3339, userData["created_at"]); err == nil {
		user.CreatedAt = createdAt
	}
	if updatedAt, err := time.Parse(time.RFC3339, userData["updated_at"]); err == nil {
		user.UpdatedAt = updatedAt
	}
	if lastLoginStr, ok := userData["last_login_at"]; ok && lastLoginStr != "" {
		if lastLogin, err := time.Parse(time.RFC3339, lastLoginStr); err == nil {
			user.LastLoginAt = &lastLogin
		}
	}

	return user, nil
}

func (s *AuthService) saveRefreshToken(token *models.RefreshToken) error {
	key := fmt.Sprintf("refresh_token:%s", token.TokenHash)
	tokenData := map[string]interface{}{
		"id":         token.ID.String(),
		"user_id":    token.UserID.String(),
		"token_hash": token.TokenHash,
		"expires_at": token.ExpiresAt.Format(time.RFC3339),
		"created_at": token.CreatedAt.Format(time.RFC3339),
	}

	ttl := time.Until(token.ExpiresAt)
	if err := s.redis.HSet(s.ctx, key, tokenData).Err(); err != nil {
		return err
	}
	return s.redis.Expire(s.ctx, key, ttl).Err()
}

func (s *AuthService) getRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	key := fmt.Sprintf("refresh_token:%s", tokenHash)
	tokenData, err := s.redis.HGetAll(s.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(tokenData) == 0 {
		return nil, redis.Nil
	}

	token := &models.RefreshToken{}
	if id, err := uuid.Parse(tokenData["id"]); err == nil {
		token.ID = id
	}
	if userID, err := uuid.Parse(tokenData["user_id"]); err == nil {
		token.UserID = userID
	}
	token.TokenHash = tokenData["token_hash"]

	if expiresAt, err := time.Parse(time.RFC3339, tokenData["expires_at"]); err == nil {
		token.ExpiresAt = expiresAt
	}
	if createdAt, err := time.Parse(time.RFC3339, tokenData["created_at"]); err == nil {
		token.CreatedAt = createdAt
	}
	if revokedAtStr, ok := tokenData["revoked_at"]; ok && revokedAtStr != "" {
		if revokedAt, err := time.Parse(time.RFC3339, revokedAtStr); err == nil {
			token.RevokedAt = &revokedAt
		}
	}

	return token, nil
}

func (s *AuthService) revokeRefreshToken(tokenHash string) error {
	key := fmt.Sprintf("refresh_token:%s", tokenHash)
	now := time.Now().Format(time.RFC3339)
	return s.redis.HSet(s.ctx, key, "revoked_at", now).Err()
}

func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) toUserInfo(user *models.User) *models.UserInfo {
	return &models.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
	}
}

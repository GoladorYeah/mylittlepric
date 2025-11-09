package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

var (
	ErrUserExists      = errors.New("user with this email already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken    = errors.New("invalid or expired token")
)

type AuthService struct {
	db          *sqlx.DB
	redis       *redis.Client
	jwtService  *utils.JWTService
	googleOAuth *GoogleOAuthService
	ctx         context.Context
}

func NewAuthService(db *sqlx.DB, redis *redis.Client, jwtService *utils.JWTService, googleOAuth *GoogleOAuthService) *AuthService {
	return &AuthService{
		db:          db,
		redis:       redis,
		jwtService:  jwtService,
		googleOAuth: googleOAuth,
		ctx:         context.Background(),
	}
}

// Signup creates a new user account
func (s *AuthService) Signup(req *models.SignupRequest) (*models.AuthResponse, error) {
	// Validate input
	if err := validateEmail(req.Email); err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

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
		Provider:     "email",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.saveUser(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Generate tokens
	return s.generateAuthResponse(user)
}

// GoogleLogin authenticates a user via Google OAuth and returns tokens
func (s *AuthService) GoogleLogin(idToken string) (*models.AuthResponse, error) {
	// Verify the Google ID token
	googleUser, err := s.googleOAuth.VerifyIDToken(s.ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Google token: %w", err)
	}

	// Check if user exists by provider ID
	user, err := s.getUserByProviderID("google", googleUser.Sub)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// If user doesn't exist, create new account
	if errors.Is(err, redis.Nil) {
		user = &models.User{
			ID:         uuid.New(),
			Email:      googleUser.Email,
			FullName:   googleUser.Name,
			Picture:    googleUser.Picture,
			Provider:   "google",
			ProviderID: googleUser.Sub,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := s.saveUser(user); err != nil {
			return nil, fmt.Errorf("failed to save user: %w", err)
		}
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	user.Picture = googleUser.Picture // Update picture in case it changed
	user.FullName = googleUser.Name   // Update name in case it changed
	if err := s.saveUser(user); err != nil {
		fmt.Printf("Warning: failed to update user info: %v\n", err)
	}

	// Generate tokens
	return s.generateAuthResponse(user)
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Validate input
	if err := validateEmail(req.Email); err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

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
// Updates chat sessions, search history in PostgreSQL, and Redis cache
// This ensures cross-device synchronization after login
func (s *AuthService) ClaimSessions(userID uuid.UUID, sessionIDs []string) error {
	sessionCount := 0
	historyCount := 0

	for _, sessionID := range sessionIDs {
		// 1. Update Redis cache (for fast access)
		key := fmt.Sprintf("session:%s", sessionID)
		if err := s.redis.HSet(s.ctx, key, "user_id", userID.String()).Err(); err != nil {
			fmt.Printf("⚠️ Failed to claim session %s in Redis: %v\n", sessionID, err)
			// Continue to PostgreSQL even if Redis fails
		}

		// 2. Update chat_sessions in PostgreSQL (primary source of truth)
		// This enables cross-device session synchronization
		sessionQuery := `UPDATE chat_sessions SET user_id = $1, updated_at = NOW() WHERE session_id = $2`
		sessionResult, err := s.db.ExecContext(s.ctx, sessionQuery, userID, sessionID)
		if err != nil {
			fmt.Printf("⚠️ Failed to claim session %s in PostgreSQL: %v\n", sessionID, err)
			continue
		}

		// Check if session was updated
		if rows, _ := sessionResult.RowsAffected(); rows > 0 {
			sessionCount++
			fmt.Printf("✅ Claimed session %s for user %s\n", sessionID, userID.String())
		}

		// 3. Update search_history in PostgreSQL
		// Link all anonymous search history from this session to the user
		historyQuery := `UPDATE search_history SET user_id = $1 WHERE session_id = $2 AND user_id IS NULL`
		historyResult, err := s.db.ExecContext(s.ctx, historyQuery, userID, sessionID)
		if err != nil {
			fmt.Printf("⚠️ Failed to claim search history for session %s: %v\n", sessionID, err)
			continue
		}

		// Check how many history records were updated
		if rows, _ := historyResult.RowsAffected(); rows > 0 {
			historyCount += int(rows)
			fmt.Printf("✅ Claimed %d search history records for session %s\n", rows, sessionID)
		}
	}

	if sessionCount == 0 && len(sessionIDs) > 0 {
		return fmt.Errorf("failed to claim any sessions")
	}

	fmt.Printf("✅ Successfully claimed %d sessions and %d search history records for user %s\n",
		sessionCount, historyCount, userID.String())
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
		"picture":       user.Picture,
		"provider":      user.Provider,
		"provider_id":   user.ProviderID,
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
	if err := s.redis.Set(s.ctx, emailKey, user.ID.String(), 0).Err(); err != nil {
		return err
	}

	// Create provider -> ID mapping for OAuth users
	if user.Provider != "email" && user.ProviderID != "" {
		providerKey := fmt.Sprintf("user:provider:%s:%s", user.Provider, user.ProviderID)
		if err := s.redis.Set(s.ctx, providerKey, user.ID.String(), 0).Err(); err != nil {
			return err
		}
	}

	return nil
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
	user.Picture = userData["picture"]
	user.Provider = userData["provider"]
	user.ProviderID = userData["provider_id"]

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

func (s *AuthService) getUserByProviderID(provider, providerID string) (*models.User, error) {
	providerKey := fmt.Sprintf("user:provider:%s:%s", provider, providerID)
	userIDStr, err := s.redis.Get(s.ctx, providerKey).Result()
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return s.getUserByID(userID)
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
		Picture:   user.Picture,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
	}
}

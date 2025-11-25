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

	"mylittleprice/ent"
	"mylittleprice/ent/chatsession"
	"mylittleprice/ent/searchhistory"
	"mylittleprice/ent/user"
	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

var (
	ErrUserExists           = errors.New("user with this email already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrPasswordNotSet       = errors.New("password not set for OAuth users")
	ErrTokenExpired         = errors.New("reset token has expired")
	ErrTokenAlreadyUsed     = errors.New("reset token has already been used")
	ErrResetTokenNotFound   = errors.New("reset token not found")
)

type AuthService struct {
	client      *ent.Client
	redis       *redis.Client
	jwtService  *utils.JWTService
	googleOAuth *GoogleOAuthService
	ctx         context.Context
}

func NewAuthService(client *ent.Client, redis *redis.Client, jwtService *utils.JWTService, googleOAuth *GoogleOAuthService) *AuthService {
	return &AuthService{
		client:      client,
		redis:       redis,
		jwtService:  jwtService,
		googleOAuth: googleOAuth,
		ctx:         context.Background(),
	}
}

// UserLookup defines criteria for looking up a user
type UserLookup struct {
	ByID         *uuid.UUID
	ByEmail      *string
	ByProviderID *struct {
		Provider string
		ID       string
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

		// 2. Update chat_sessions in PostgreSQL using Ent
		// This enables cross-device session synchronization
		sessionRows, err := s.client.ChatSession.Update().
			Where(chatsession.SessionIDEQ(sessionID)).
			SetUserID(userID).
			SetUpdatedAt(time.Now()).
			Save(s.ctx)

		if err != nil {
			fmt.Printf("⚠️ Failed to claim session %s in PostgreSQL: %v\n", sessionID, err)
			continue
		}

		// Check if session was updated
		if sessionRows > 0 {
			sessionCount++
			fmt.Printf("✅ Claimed session %s for user %s\n", sessionID, userID.String())
		}

		// 3. Update search_history in PostgreSQL using Ent
		// Link all anonymous search history from this session to the user
		historyRows, err := s.client.SearchHistory.Update().
			Where(
				searchhistory.And(
					searchhistory.SessionIDEQ(sessionID),
					searchhistory.UserIDIsNil(),
				),
			).
			SetUserID(userID).
			Save(s.ctx)

		if err != nil {
			fmt.Printf("⚠️ Failed to claim search history for session %s: %v\n", sessionID, err)
			continue
		}

		// Check how many history records were updated
		if historyRows > 0 {
			historyCount += historyRows
			fmt.Printf("✅ Claimed %d search history records for session %s\n", historyRows, sessionID)
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
	// 1. Save to PostgreSQL (source of truth for foreign key constraints)
	if err := s.SaveUserToPostgres(user); err != nil {
		return fmt.Errorf("failed to save user to postgres: %w", err)
	}

	// 2. Save to Redis (for fast access)
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

// SaveUserToPostgres saves a user to PostgreSQL using Ent (exported for use by other services)
func (s *AuthService) SaveUserToPostgres(userModel *models.User) error {
	// Check if user exists
	exists, err := s.client.User.Query().
		Where(user.IDEQ(userModel.ID)).
		Exist(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		// Update existing user
		updateBuilder := s.client.User.UpdateOneID(userModel.ID).
			SetEmail(userModel.Email).
			SetPasswordHash(userModel.PasswordHash).
			SetName(userModel.FullName).
			SetAvatarURL(userModel.Picture).
			SetProvider(userModel.Provider).
			SetUpdatedAt(userModel.UpdatedAt)

		// Only set google_id if it's not empty (to avoid unique constraint violation)
		if userModel.ProviderID != "" {
			updateBuilder.SetGoogleID(userModel.ProviderID)
		} else {
			updateBuilder.ClearGoogleID()
		}

		// Set optional last_login
		if userModel.LastLoginAt != nil {
			updateBuilder.SetLastLogin(*userModel.LastLoginAt)
		}

		_, err = updateBuilder.Save(s.ctx)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		// Create new user
		createBuilder := s.client.User.Create().
			SetID(userModel.ID).
			SetEmail(userModel.Email).
			SetPasswordHash(userModel.PasswordHash).
			SetName(userModel.FullName).
			SetAvatarURL(userModel.Picture).
			SetProvider(userModel.Provider).
			SetCreatedAt(userModel.CreatedAt).
			SetUpdatedAt(userModel.UpdatedAt)

		// Only set google_id if it's not empty (to avoid unique constraint violation)
		if userModel.ProviderID != "" {
			createBuilder.SetGoogleID(userModel.ProviderID)
		}

		// Set optional last_login
		if userModel.LastLoginAt != nil {
			createBuilder.SetLastLogin(*userModel.LastLoginAt)
		}

		_, err = createBuilder.Save(s.ctx)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	return nil
}

func (s *AuthService) getUserByEmail(email string) (*models.User, error) {
	return s.getUser(UserLookup{ByEmail: &email})
}

func (s *AuthService) getUserByID(userID uuid.UUID) (*models.User, error) {
	return s.getUser(UserLookup{ByID: &userID})
}

func (s *AuthService) getUserByProviderID(provider, providerID string) (*models.User, error) {
	return s.getUser(UserLookup{
		ByProviderID: &struct {
			Provider string
			ID       string
		}{
			Provider: provider,
			ID:       providerID,
		},
	})
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

// entUserToModel converts an Ent User entity to a models.User
func (s *AuthService) entUserToModel(entUser *ent.User) *models.User {
	userModel := &models.User{
		ID:           entUser.ID,
		Email:        entUser.Email,
		PasswordHash: entUser.PasswordHash,
		FullName:     entUser.Name,
		Picture:      entUser.AvatarURL,
		Provider:     entUser.Provider,
		ProviderID:   entUser.GoogleID,
		CreatedAt:    entUser.CreatedAt,
		UpdatedAt:    entUser.UpdatedAt,
	}

	// Handle optional last_login
	if !entUser.LastLogin.IsZero() {
		userModel.LastLoginAt = &entUser.LastLogin
	}

	return userModel
}

// syncUserToRedis syncs a user to Redis cache without touching PostgreSQL
func (s *AuthService) syncUserToRedis(user *models.User) error {
	// 1. Save user data
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
		return fmt.Errorf("failed to sync user data to Redis: %w", err)
	}

	// 2. Create email -> ID mapping
	emailKey := fmt.Sprintf("user:email:%s", user.Email)
	if err := s.redis.Set(s.ctx, emailKey, user.ID.String(), 0).Err(); err != nil {
		return fmt.Errorf("failed to sync email mapping to Redis: %w", err)
	}

	// 3. Create provider -> ID mapping for OAuth users
	if user.Provider != "email" && user.ProviderID != "" {
		providerKey := fmt.Sprintf("user:provider:%s:%s", user.Provider, user.ProviderID)
		if err := s.redis.Set(s.ctx, providerKey, user.ID.String(), 0).Err(); err != nil {
			return fmt.Errorf("failed to sync provider mapping to Redis: %w", err)
		}
	}

	return nil
}

// getUserWithFallback retrieves a user by ID with Redis -> PostgreSQL fallback
func (s *AuthService) getUserWithFallback(userID uuid.UUID) (*models.User, error) {
	// 1. Try Redis first
	userKey := fmt.Sprintf("user:id:%s", userID.String())
	userData, err := s.redis.HGetAll(s.ctx, userKey).Result()

	if err == nil && len(userData) > 0 {
		// Redis hit - parse and return
		user := &models.User{}
		user.ID = userID
		user.Email = userData["email"]
		user.PasswordHash = userData["password_hash"]
		user.FullName = userData["full_name"]
		user.Picture = userData["picture"]
		user.Provider = userData["provider"]
		user.ProviderID = userData["provider_id"]

		if createdAt, parseErr := time.Parse(time.RFC3339, userData["created_at"]); parseErr == nil {
			user.CreatedAt = createdAt
		}
		if updatedAt, parseErr := time.Parse(time.RFC3339, userData["updated_at"]); parseErr == nil {
			user.UpdatedAt = updatedAt
		}
		if lastLoginStr, ok := userData["last_login_at"]; ok && lastLoginStr != "" {
			if lastLogin, parseErr := time.Parse(time.RFC3339, lastLoginStr); parseErr == nil {
				user.LastLoginAt = &lastLogin
			}
		}

		return user, nil
	}

	// 2. Fallback to PostgreSQL through Ent
	fmt.Printf("⚠️ Redis miss for user %s, falling back to PostgreSQL\n", userID)
	entUser, err := s.client.User.Get(s.ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, redis.Nil // Maintain consistent error type
		}
		return nil, fmt.Errorf("user not found in PostgreSQL: %w", err)
	}

	// 3. Convert Ent user to model
	userModel := s.entUserToModel(entUser)

	// 4. Sync back to Redis for next requests
	if syncErr := s.syncUserToRedis(userModel); syncErr != nil {
		fmt.Printf("⚠️ Failed to sync user to Redis: %v\n", syncErr)
		// Don't return error - user retrieved from DB successfully
	}

	return userModel, nil
}

// lookupUserIDByEmail gets user ID from Redis by email, returns empty UUID if not found
func (s *AuthService) lookupUserIDByEmail(email string) (uuid.UUID, error) {
	emailKey := fmt.Sprintf("user:email:%s", email)
	userIDStr, err := s.redis.Get(s.ctx, emailKey).Result()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in Redis: %w", err)
	}

	return userID, nil
}

// lookupUserIDByProvider gets user ID from Redis by provider, returns empty UUID if not found
func (s *AuthService) lookupUserIDByProvider(provider, providerID string) (uuid.UUID, error) {
	providerKey := fmt.Sprintf("user:provider:%s:%s", provider, providerID)
	userIDStr, err := s.redis.Get(s.ctx, providerKey).Result()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in Redis: %w", err)
	}

	return userID, nil
}

// getUser is a unified method for retrieving users by different criteria
func (s *AuthService) getUser(lookup UserLookup) (*models.User, error) {
	var userID uuid.UUID
	var err error

	// 1. Determine userID from lookup criteria
	switch {
	case lookup.ByID != nil:
		userID = *lookup.ByID
	case lookup.ByEmail != nil:
		// Try Redis first for email -> ID mapping
		userID, err = s.lookupUserIDByEmail(*lookup.ByEmail)
		if err == nil && userID != uuid.Nil {
			// Found in Redis, now get the user
			user, getUserErr := s.getUserWithFallback(userID)
			if getUserErr == nil {
				return user, nil
			}
		}

		// Fallback to PostgreSQL for email lookup
		fmt.Printf("⚠️ Redis miss for user email %s, falling back to PostgreSQL\n", *lookup.ByEmail)
		entUser, dbErr := s.client.User.Query().
			Where(user.EmailEQ(*lookup.ByEmail)).
			Only(s.ctx)

		if dbErr != nil {
			if ent.IsNotFound(dbErr) {
				return nil, redis.Nil // Maintain consistent error type
			}
			return nil, fmt.Errorf("user not found in PostgreSQL: %w", dbErr)
		}

		// Convert and sync to Redis
		userModel := s.entUserToModel(entUser)
		if syncErr := s.syncUserToRedis(userModel); syncErr != nil {
			fmt.Printf("⚠️ Failed to sync user to Redis: %v\n", syncErr)
		}
		return userModel, nil

	case lookup.ByProviderID != nil:
		// Try Redis first for provider -> ID mapping
		userID, err = s.lookupUserIDByProvider(lookup.ByProviderID.Provider, lookup.ByProviderID.ID)
		if err == nil && userID != uuid.Nil {
			// Found in Redis, now get the user
			user, getUserErr := s.getUserWithFallback(userID)
			if getUserErr == nil {
				return user, nil
			}
		}

		// Fallback to PostgreSQL for provider lookup
		fmt.Printf("⚠️ Redis miss for provider %s:%s, falling back to PostgreSQL\n",
			lookup.ByProviderID.Provider, lookup.ByProviderID.ID)

		var entUser *ent.User
		if lookup.ByProviderID.Provider == "google" {
			entUser, err = s.client.User.Query().
				Where(user.GoogleIDEQ(lookup.ByProviderID.ID)).
				Only(s.ctx)
		} else {
			return nil, fmt.Errorf("unsupported provider: %s", lookup.ByProviderID.Provider)
		}

		if err != nil {
			if ent.IsNotFound(err) {
				return nil, redis.Nil // Maintain consistent error type
			}
			return nil, fmt.Errorf("user not found in PostgreSQL: %w", err)
		}

		// Convert and sync to Redis
		userModel := s.entUserToModel(entUser)
		if syncErr := s.syncUserToRedis(userModel); syncErr != nil {
			fmt.Printf("⚠️ Failed to sync user to Redis: %v\n", syncErr)
		}
		return userModel, nil

	default:
		return nil, fmt.Errorf("no lookup criteria provided")
	}

	// For ByID case, use the unified fallback method
	return s.getUserWithFallback(userID)
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

// ==================== Password Management Methods ====================

// ChangePassword updates the user's password after verifying the current password
func (s *AuthService) ChangePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	// Get user
	user, err := s.getUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is using OAuth (Google)
	if user.Provider != "email" {
		return ErrPasswordNotSet
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrInvalidPassword
	}

	// Validate new password
	if err := validatePassword(newPassword); err != nil {
		return fmt.Errorf("invalid new password: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password in database
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	// Save to database
	if err := s.saveUser(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all existing refresh tokens to force re-login on all devices
	if err := s.revokeAllUserRefreshTokens(userID); err != nil {
		fmt.Printf("Warning: failed to revoke refresh tokens: %v\n", err)
		// Don't fail the password change if token revocation fails
	}

	return nil
}

// RequestPasswordReset generates a password reset token and returns it
// In production, this token should be sent via email, not returned directly
func (s *AuthService) RequestPasswordReset(email string) (string, error) {
	// Get user by email
	user, err := s.getUserByEmail(email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Don't reveal if user exists or not for security
			return "", nil
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is using OAuth
	if user.Provider != "email" {
		// Don't reveal OAuth users for security
		return "", nil
	}

	// Generate reset token (random secure token)
	resetToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Hash token for storage
	tokenHash := s.hashToken(resetToken)

	// Save reset token to database (expires in 1 hour)
	resetTokenData := &models.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.savePasswordResetToken(resetTokenData); err != nil {
		return "", fmt.Errorf("failed to save reset token: %w", err)
	}

	return resetToken, nil
}

// ResetPassword resets the user's password using a reset token
func (s *AuthService) ResetPassword(resetToken, newPassword string) error {
	// Validate new password
	if err := validatePassword(newPassword); err != nil {
		return fmt.Errorf("invalid new password: %w", err)
	}

	// Hash token
	tokenHash := s.hashToken(resetToken)

	// Get reset token from database
	resetTokenData, err := s.getPasswordResetToken(tokenHash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrResetTokenNotFound
		}
		return fmt.Errorf("failed to get reset token: %w", err)
	}

	// Check if token is already used
	if resetTokenData.UsedAt != nil {
		return ErrTokenAlreadyUsed
	}

	// Check if token is expired
	if time.Now().After(resetTokenData.ExpiresAt) {
		return ErrTokenExpired
	}

	// Get user
	user, err := s.getUserByID(resetTokenData.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	// Save to database
	if err := s.saveUser(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	now := time.Now()
	resetTokenData.UsedAt = &now
	if err := s.markPasswordResetTokenAsUsed(tokenHash); err != nil {
		fmt.Printf("Warning: failed to mark token as used: %v\n", err)
	}

	// Revoke all existing refresh tokens to force re-login
	if err := s.revokeAllUserRefreshTokens(user.ID); err != nil {
		fmt.Printf("Warning: failed to revoke refresh tokens: %v\n", err)
	}

	return nil
}

// ==================== Password Reset Token Helpers ====================

func (s *AuthService) savePasswordResetToken(token *models.PasswordResetToken) error {
	key := fmt.Sprintf("password_reset:%s", token.TokenHash)
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

func (s *AuthService) getPasswordResetToken(tokenHash string) (*models.PasswordResetToken, error) {
	key := fmt.Sprintf("password_reset:%s", tokenHash)
	tokenData, err := s.redis.HGetAll(s.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(tokenData) == 0 {
		return nil, redis.Nil
	}

	token := &models.PasswordResetToken{}
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
	if usedAtStr, ok := tokenData["used_at"]; ok && usedAtStr != "" {
		if usedAt, err := time.Parse(time.RFC3339, usedAtStr); err == nil {
			token.UsedAt = &usedAt
		}
	}

	return token, nil
}

func (s *AuthService) markPasswordResetTokenAsUsed(tokenHash string) error {
	key := fmt.Sprintf("password_reset:%s", tokenHash)
	now := time.Now().Format(time.RFC3339)
	return s.redis.HSet(s.ctx, key, "used_at", now).Err()
}

func (s *AuthService) revokeAllUserRefreshTokens(userID uuid.UUID) error {
	// In Redis, we need to find all refresh tokens for this user
	// This is a simplification - in production you might want to maintain a user->tokens index
	pattern := "refresh_token:*"
	var cursor uint64
	for {
		keys, nextCursor, err := s.redis.Scan(s.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			tokenData, err := s.redis.HGetAll(s.ctx, key).Result()
			if err != nil {
				continue
			}

			if tokenData["user_id"] == userID.String() {
				now := time.Now().Format(time.RFC3339)
				s.redis.HSet(s.ctx, key, "revoked_at", now)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

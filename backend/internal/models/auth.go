package models

import (
	"time"

	"github.com/google/uuid"
)

// ═══════════════════════════════════════════════════════════
// AUTHENTICATION MODELS
// ═══════════════════════════════════════════════════════════

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"` // Never expose password
	FullName     string     `json:"full_name,omitempty" db:"full_name"`
	Picture      string     `json:"picture,omitempty" db:"picture"`         // Profile picture URL
	Provider     string     `json:"provider" db:"provider"`                 // "google", "email"
	ProviderID   string     `json:"provider_id,omitempty" db:"provider_id"` // Google user ID
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash string     `json:"-" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

// ═══════════════════════════════════════════════════════════
// AUTH REQUEST/RESPONSE MODELS
// ═══════════════════════════════════════════════════════════

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user"`
	ExpiresIn    int64     `json:"expires_in"` // seconds
}

type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name,omitempty"`
	Picture   string    `json:"picture,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ClaimSessionsRequest struct {
	SessionIDs []string `json:"session_ids" validate:"required"`
}

// ═══════════════════════════════════════════════════════════
// GOOGLE OAUTH MODELS
// ═══════════════════════════════════════════════════════════

type GoogleAuthRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"` // Google user ID
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

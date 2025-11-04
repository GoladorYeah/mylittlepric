package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mylittleprice/internal/config"
	"mylittleprice/internal/models"
)

const (
	googleTokenInfoURL = "https://oauth2.googleapis.com/tokeninfo?id_token="
)

type GoogleOAuthService struct {
	config *config.Config
	client *http.Client
}

func NewGoogleOAuthService(cfg *config.Config) *GoogleOAuthService {
	return &GoogleOAuthService{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// tokenInfoResponse represents Google's tokeninfo response
type tokenInfoResponse struct {
	Aud           string `json:"aud"`             // Audience (client ID)
	Sub           string `json:"sub"`             // User ID
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"` // "true" or "false" as string
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Exp           string `json:"exp"` // Expiration time
}

// VerifyIDToken verifies a Google ID token and returns user information
func (s *GoogleOAuthService) VerifyIDToken(ctx context.Context, idToken string) (*models.GoogleUserInfo, error) {
	// Call Google's tokeninfo endpoint to verify the token
	url := googleTokenInfoURL + idToken
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Google token verification failed. Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("invalid token (status %d)", resp.StatusCode)
	}

	// Parse the response
	var tokenInfo tokenInfoResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("🔍 Google tokeninfo response: %s\n", string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Verify the audience (client ID) matches our application
	if tokenInfo.Aud != s.config.GoogleClientID {
		fmt.Printf("❌ Client ID mismatch. Expected: %s, Got: %s\n", s.config.GoogleClientID, tokenInfo.Aud)
		return nil, fmt.Errorf("token is for different client ID")
	}

	// Verify email is verified
	if tokenInfo.EmailVerified != "true" {
		return nil, fmt.Errorf("email not verified")
	}

	// Convert to our model
	userInfo := &models.GoogleUserInfo{
		Sub:           tokenInfo.Sub,
		Email:         tokenInfo.Email,
		EmailVerified: true,
		Name:          tokenInfo.Name,
		Picture:       tokenInfo.Picture,
		GivenName:     tokenInfo.GivenName,
		FamilyName:    tokenInfo.FamilyName,
		Locale:        tokenInfo.Locale,
	}

	fmt.Printf("✅ Token verified successfully for user: %s (%s)\n", userInfo.Name, userInfo.Email)
	return userInfo, nil
}

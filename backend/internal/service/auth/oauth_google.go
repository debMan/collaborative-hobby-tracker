package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/auth"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// GoogleTokenResponse represents the response from Google's token endpoint
type GoogleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// GoogleUserInfo represents user information from Google
type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// NewGoogleOAuthService creates a new OAuth service configured for Google
func NewGoogleOAuthService(config *OAuthConfig) *OAuthService {
	// Set defaults for Google production URLs if not specified (for testing)
	if config.TokenURL == "" {
		config.TokenURL = googleTokenURL
	}
	if config.UserInfoURL == "" {
		config.UserInfoURL = googleUserInfoURL
	}

	return NewOAuthService(config)
}

// GetGoogleAuthURL generates the Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", s.config.ClientID)
	params.Add("redirect_uri", s.config.RedirectURL)
	params.Add("response_type", "code")
	params.Add("scope", "email profile")
	params.Add("state", state)

	return googleAuthURL + "?" + params.Encode()
}

// HandleGoogleCallback handles the OAuth callback from Google
func (s *OAuthService) HandleGoogleCallback(ctx context.Context, code string, authService *Service) (*AuthResult, error) {
	if code == "" {
		return nil, ErrOAuthCodeRequired
	}

	// Exchange authorization code for access token
	accessToken, err := s.exchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOAuthTokenExchange, err)
	}

	// Get user info from Google
	userInfo, err := s.getUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOAuthUserInfo, err)
	}

	// Find or create user in database
	user, err := authService.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil && err != repository.ErrUserNotFound {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		// Create new user
		user = &models.User{
			Email:             userInfo.Email,
			Name:              userInfo.Name,
			OAuthProvider:     userInfo.ID,
			OAuthProviderName: "google",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err = authService.userRepo.Create(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Update existing user with OAuth info if not set
		if user.OAuthProvider == "" {
			user.OAuthProvider = userInfo.ID
			user.OAuthProviderName = "google"
			user.UpdatedAt = time.Now()

			err = authService.userRepo.Update(ctx, user)
			if err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
		}
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, authService.jwtSecret, authService.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResult{
		Token:  token,
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Name:   user.Name,
	}, nil
}

// exchangeCodeForToken exchanges the authorization code for an access token
func (s *OAuthService) exchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", s.config.ClientID)
	data.Set("client_secret", s.config.ClientSecret)
	data.Set("redirect_uri", s.config.RedirectURL)

	req, err := http.NewRequest(http.MethodPost, s.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// getUserInfo fetches user information from Google
func (s *OAuthService) getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, s.config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user info request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

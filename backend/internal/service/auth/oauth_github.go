package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/auth"
)

const (
	githubAuthURL      = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
)

// GitHubTokenResponse represents the response from GitHub's token endpoint
type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GitHubUserInfo represents user information from GitHub
type GitHubUserInfo struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// GitHubEmail represents an email from GitHub
type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// NewGitHubOAuthService creates a new OAuth service configured for GitHub
func NewGitHubOAuthService(config *OAuthConfig) *OAuthService {
	// Set defaults for GitHub production URLs if not specified (for testing)
	if config.TokenURL == "" {
		config.TokenURL = githubTokenURL
	}
	if config.UserInfoURL == "" {
		config.UserInfoURL = githubUserURL
	}
	if config.UserEmailsURL == "" {
		config.UserEmailsURL = githubUserEmailURL
	}

	return NewOAuthService(config)
}

// GetGitHubAuthURL generates the GitHub OAuth authorization URL
func (s *OAuthService) GetGitHubAuthURL(state string) string {

	params := url.Values{}
	params.Add("client_id", s.config.ClientID)
	params.Add("redirect_uri", s.config.RedirectURL)
	params.Add("scope", "user:email")
	params.Add("state", state)

	return githubAuthURL + "?" + params.Encode()
}

// HandleGitHubCallback handles the OAuth callback from GitHub
func (s *OAuthService) HandleGitHubCallback(ctx context.Context, code string, authService *Service) (*AuthResult, error) {
	if code == "" {
		return nil, ErrOAuthCodeRequired
	}

	// Exchange authorization code for access token
	accessToken, err := s.exchangeGitHubCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOAuthTokenExchange, err)
	}

	// Get user info from GitHub
	userInfo, err := s.getGitHubUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOAuthUserInfo, err)
	}

	// Get user emails from GitHub
	email, err := s.getGitHubUserEmail(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user email: %w", err)
	}

	// Find or create user in database
	user, err := authService.userRepo.FindByEmail(ctx, email)
	if err != nil && err != repository.ErrUserNotFound {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		// Create new user
		user = &models.User{
			Email:             email,
			Name:              userInfo.Name,
			OAuthProvider:     strconv.Itoa(userInfo.ID),
			OAuthProviderName: "github",
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
			user.OAuthProvider = strconv.Itoa(userInfo.ID)
			user.OAuthProviderName = "github"
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

// exchangeGitHubCodeForToken exchanges the authorization code for an access token
func (s *OAuthService) exchangeGitHubCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", s.config.ClientID)
	data.Set("client_secret", s.config.ClientSecret)
	data.Set("redirect_uri", s.config.RedirectURL)

	req, err := http.NewRequest(http.MethodPost, s.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

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

	var tokenResp GitHubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// getGitHubUserInfo fetches user information from GitHub
func (s *OAuthService) getGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, s.config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

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

	var userInfo GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// getGitHubUserEmail fetches the primary verified email from GitHub
func (s *OAuthService) getGitHubUserEmail(accessToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, s.config.UserEmailsURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("emails request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var emails []GitHubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode emails: %w", err)
	}

	// Find primary verified email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// If no primary verified email, use first verified email
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	// If no verified email, use first email
	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", fmt.Errorf("no email found for user")
}

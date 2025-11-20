package auth

import (
	"errors"
)

// OAuth Architecture:
//
// This package provides a flexible OAuth 2.0 implementation that supports multiple providers.
// The architecture is designed to be extensible and maintainable:
//
// 1. Common Components (this file):
//    - OAuthConfig: Shared configuration structure
//    - OAuthService: Core service that can be used with any provider
//    - Common error definitions
//    - Optional OAuthProvider interface for future strategy pattern implementation
//
// 2. Provider-Specific Implementations:
//    - oauth_google.go: Google OAuth implementation
//    - oauth_github.go: GitHub OAuth implementation
//    - Future providers (Apple, Twitter/X, etc.) can follow the same pattern
//
// 3. Usage Pattern:
//    // For Google:
//    config := &OAuthConfig{ClientID: "...", ClientSecret: "...", RedirectURL: "..."}
//    service := NewGoogleOAuthService(config)  // Sets Google-specific defaults
//    url := service.GetGoogleAuthURL(state)
//    result, err := service.HandleGoogleCallback(ctx, code, authService)
//
//    // For GitHub:
//    service := NewGitHubOAuthService(config)  // Sets GitHub-specific defaults
//    url := service.GetGitHubAuthURL(state)
//    result, err := service.HandleGitHubCallback(ctx, code, authService)
//
// Adding a New Provider:
// 1. Create oauth_<provider>.go file
// 2. Define provider-specific constants (auth URL, token URL, etc.)
// 3. Create New<Provider>OAuthService() helper that sets defaults
// 4. Implement Get<Provider>AuthURL() method
// 5. Implement Handle<Provider>Callback() method
// 6. Add comprehensive tests in oauth_<provider>_test.go

// Common OAuth errors shared across all providers
var (
	// ErrOAuthCodeRequired is returned when authorization code is empty
	ErrOAuthCodeRequired = errors.New("authorization code is required")
	// ErrOAuthTokenExchange is returned when token exchange fails
	ErrOAuthTokenExchange = errors.New("failed to exchange authorization code for token")
	// ErrOAuthUserInfo is returned when fetching user info fails
	ErrOAuthUserInfo = errors.New("failed to fetch user information")
)

// OAuthConfig holds OAuth configuration for any provider
// Provider-specific URLs are set by helper functions like NewGoogleOAuthService()
type OAuthConfig struct {
	ClientID      string // OAuth client ID from provider
	ClientSecret  string // OAuth client secret from provider
	RedirectURL   string // Callback URL for OAuth flow
	TokenURL      string // Provider token endpoint (auto-set in production, customizable for testing)
	UserInfoURL   string // Provider user info endpoint (auto-set in production, customizable for testing)
	UserEmailsURL string // GitHub-specific: endpoint to fetch user emails (auto-set for GitHub)
}

// OAuthService handles OAuth operations for multiple providers
// Use provider-specific constructors (NewGoogleOAuthService, NewGitHubOAuthService)
// instead of creating this directly
type OAuthService struct {
	config *OAuthConfig
}

// NewOAuthService creates a new OAuth service with the given configuration
// Note: Prefer using provider-specific helpers like NewGoogleOAuthService()
// or NewGitHubOAuthService() which set appropriate defaults
func NewOAuthService(config *OAuthConfig) *OAuthService {
	return &OAuthService{
		config: config,
	}
}

// OAuthProvider defines the interface for OAuth providers
// This interface can be implemented by specific providers for a strategy pattern approach
// Currently, we use a simpler approach with provider-specific methods on OAuthService
// Future refactoring could use this interface for more formal provider abstraction
type OAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(code string) (string, error)
	GetUserInfo(accessToken string) (OAuthUserInfo, error)
}

// OAuthUserInfo represents common user information from OAuth providers
// Provider-specific implementations convert their API responses to this common format
type OAuthUserInfo struct {
	ID    string // Provider-specific user ID
	Email string // User email address
	Name  string // User display name
}

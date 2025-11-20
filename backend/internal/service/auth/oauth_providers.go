package auth

import (
	"context"
)

// This file demonstrates an alternative architecture using provider-specific types.
// Currently NOT in use - kept as documentation for potential future refactoring.
//
// The current approach (provider-specific methods on OAuthService) is simpler
// and sufficient for our needs. Use this pattern if:
// - You have many providers (5+)
// - You need runtime provider selection
// - You want strict interface enforcement

// Provider defines the contract that all OAuth providers must implement
// This is more specific than OAuthProvider and includes callback handling
type Provider interface {
	// GetAuthorizationURL generates the OAuth authorization URL
	GetAuthorizationURL(state string) string

	// HandleCallback processes the OAuth callback and returns user info
	HandleCallback(ctx context.Context, code string, authService *Service) (*AuthResult, error)

	// Name returns the provider name (e.g., "google", "github")
	Name() string
}

// Example implementation for Google:
//
// type GoogleProvider struct {
//     config *OAuthConfig
// }
//
// func NewGoogleProvider(config *OAuthConfig) Provider {
//     return &GoogleProvider{config: config}
// }
//
// func (p *GoogleProvider) GetAuthorizationURL(state string) string {
//     // Google-specific implementation
// }
//
// func (p *GoogleProvider) HandleCallback(ctx context.Context, code string, authService *Service) (*AuthResult, error) {
//     // Google-specific implementation
// }
//
// func (p *GoogleProvider) Name() string {
//     return "google"
// }

// ProviderRegistry could manage multiple providers
// type ProviderRegistry struct {
//     providers map[string]Provider
// }
//
// func (r *ProviderRegistry) Register(provider Provider) {
//     r.providers[provider.Name()] = provider
// }
//
// func (r *ProviderRegistry) Get(name string) (Provider, error) {
//     p, ok := r.providers[name]
//     if !ok {
//         return nil, fmt.Errorf("provider %s not found", name)
//     }
//     return p, nil
// }

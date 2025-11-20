package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository/testutil"
	pkgauth "github.com/debMan/collaborative-hobby-tracker/backend/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testGitHubClientID = "test-github-client-id"
const testGitHubClientSecret = "test-github-client-secret"
const testGitHubRedirectURL = "http://localhost:8080/auth/github/callback"

func TestGitHubOAuthURL(t *testing.T) {
	config := &OAuthConfig{
		ClientID:     testGitHubClientID,
		ClientSecret: testGitHubClientSecret,
		RedirectURL:  testGitHubRedirectURL,
	}

	oauthService := NewOAuthService(config)

	t.Run("Generate GitHub OAuth URL", func(t *testing.T) {
		state := "random-state-string"
		url := oauthService.GetGitHubAuthURL(state)

		// Verify URL contains required parameters
		assert.Contains(t, url, "https://github.com/login/oauth/authorize")
		assert.Contains(t, url, "client_id="+testGitHubClientID)
		assert.Contains(t, url, "redirect_uri=")
		assert.Contains(t, url, "state="+state)
		assert.Contains(t, url, "scope=")
		// user:email will be URL-encoded
		assert.Contains(t, url, "email")
	})

	t.Run("Generate different URLs with different states", func(t *testing.T) {
		state1 := "state-1"
		state2 := "state-2"

		url1 := oauthService.GetGitHubAuthURL(state1)
		url2 := oauthService.GetGitHubAuthURL(state2)

		assert.NotEqual(t, url1, url2)
		assert.Contains(t, url1, "state="+state1)
		assert.Contains(t, url2, "state="+state2)
	})
}

func TestGitHubOAuthCallback(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create indexes
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(mongo.DB)
	authService := NewService(userRepo, testJWTSecret, testTokenDuration)

	// Mock GitHub OAuth token endpoint
	mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/login/oauth/access_token", r.URL.Path)

		// Verify request contains required parameters
		err := r.ParseForm()
		require.NoError(t, err)
		assert.NotEmpty(t, r.FormValue("code")) // Accept any code
		assert.Equal(t, testGitHubClientID, r.FormValue("client_id"))
		assert.Equal(t, testGitHubClientSecret, r.FormValue("client_secret"))
		assert.Equal(t, testGitHubRedirectURL, r.FormValue("redirect_uri"))

		// Return mock access token
		response := map[string]any{
			"access_token": "mock-github-token",
			"token_type":   "Bearer",
			"scope":        "user:email",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockTokenServer.Close()

	// Mock GitHub user endpoint
	mockUserServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/user", r.URL.Path)

		// Verify Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer mock-github-token", authHeader)

		// Return mock user info
		userInfo := map[string]any{
			"id":    12345,
			"login": "githubuser",
			"name":  "GitHub User",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer mockUserServer.Close()

	// Mock GitHub emails endpoint
	mockEmailsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/user/emails", r.URL.Path)

		// Verify Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer mock-github-token", authHeader)

		// Return mock emails
		emails := []map[string]any{
			{
				"email":    "githubuser@example.com",
				"primary":  true,
				"verified": true,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(emails)
	}))
	defer mockEmailsServer.Close()

	config := &OAuthConfig{
		ClientID:      testGitHubClientID,
		ClientSecret:  testGitHubClientSecret,
		RedirectURL:   testGitHubRedirectURL,
		TokenURL:      mockTokenServer.URL + "/login/oauth/access_token",
		UserInfoURL:   mockUserServer.URL + "/user",
		UserEmailsURL: mockEmailsServer.URL + "/user/emails",
	}

	oauthService := NewOAuthService(config)

	t.Run("Handle GitHub OAuth callback for new user", func(t *testing.T) {
		code := "test-github-code"

		result, err := oauthService.HandleGitHubCallback(ctx, code, authService)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify auth result
		assert.NotEmpty(t, result.Token)
		assert.NotEmpty(t, result.UserID)
		assert.Equal(t, "githubuser@example.com", result.Email)
		assert.Equal(t, "GitHub User", result.Name)

		// Verify user was created in database
		user, err := userRepo.FindByEmail(ctx, "githubuser@example.com")
		require.NoError(t, err)
		assert.Equal(t, "githubuser@example.com", user.Email)
		assert.Equal(t, "GitHub User", user.Name)
		assert.Empty(t, user.PasswordHash) // OAuth users don't have password
		assert.Equal(t, "12345", user.OAuthProvider)
		assert.Equal(t, "github", user.OAuthProviderName)
	})

	t.Run("Handle GitHub OAuth callback for existing user", func(t *testing.T) {
		// First login creates the user
		code1 := "test-github-code-1"
		result1, err := oauthService.HandleGitHubCallback(ctx, code1, authService)
		require.NoError(t, err)

		// Second login should find existing user
		code2 := "test-github-code-2"
		result2, err := oauthService.HandleGitHubCallback(ctx, code2, authService)
		require.NoError(t, err)

		// Should return same user
		assert.Equal(t, result1.UserID, result2.UserID)
		assert.Equal(t, result1.Email, result2.Email)

		// Both tokens should be valid
		_, err = pkgauth.ValidateToken(result1.Token, testJWTSecret)
		assert.NoError(t, err)
		_, err = pkgauth.ValidateToken(result2.Token, testJWTSecret)
		assert.NoError(t, err)
	})
}

func TestGitHubOAuthErrors(t *testing.T) {
	config := &OAuthConfig{
		ClientID:      testGitHubClientID,
		ClientSecret:  testGitHubClientSecret,
		RedirectURL:   testGitHubRedirectURL,
		TokenURL:      "http://invalid-url/token",
		UserInfoURL:   "http://invalid-url/user",
		UserEmailsURL: "http://invalid-url/emails",
	}

	oauthService := NewOAuthService(config)

	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(mongo.DB)
	authService := NewService(userRepo, testJWTSecret, testTokenDuration)

	t.Run("Reject empty authorization code", func(t *testing.T) {
		result, err := oauthService.HandleGitHubCallback(ctx, "", authService)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "authorization code is required")
	})

	t.Run("Handle token exchange failure", func(t *testing.T) {
		result, err := oauthService.HandleGitHubCallback(ctx, "invalid-code", authService)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

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

const testGoogleClientID = "test-google-client-id"
const testGoogleClientSecret = "test-google-client-secret"
const testGoogleRedirectURL = "http://localhost:8080/auth/google/callback"

func TestGoogleOAuthURL(t *testing.T) {
	config := &OAuthConfig{
		ClientID:     testGoogleClientID,
		ClientSecret: testGoogleClientSecret,
		RedirectURL:  testGoogleRedirectURL,
	}

	oauthService := NewOAuthService(config)

	t.Run("Generate Google OAuth URL", func(t *testing.T) {
		state := "random-state-string"
		url := oauthService.GetGoogleAuthURL(state)

		// Verify URL contains required parameters
		assert.Contains(t, url, "https://accounts.google.com/o/oauth2/v2/auth")
		assert.Contains(t, url, "client_id="+testGoogleClientID)
		// redirect_uri will be URL-encoded in the query string
		assert.Contains(t, url, "redirect_uri=")
		assert.Contains(t, url, "response_type=code")
		assert.Contains(t, url, "scope=")
		assert.Contains(t, url, "state="+state)
		assert.Contains(t, url, "email")
		assert.Contains(t, url, "profile")
	})

	t.Run("Generate different URLs with different states", func(t *testing.T) {
		state1 := "state-1"
		state2 := "state-2"

		url1 := oauthService.GetGoogleAuthURL(state1)
		url2 := oauthService.GetGoogleAuthURL(state2)

		assert.NotEqual(t, url1, url2)
		assert.Contains(t, url1, "state="+state1)
		assert.Contains(t, url2, "state="+state2)
	})
}

func TestGoogleOAuthCallback(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create indexes
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(mongo.DB)
	authService := NewService(userRepo, testJWTSecret, testTokenDuration)

	// Mock Google OAuth token endpoint
	mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/token", r.URL.Path)

		// Verify request contains required parameters
		err := r.ParseForm()
		require.NoError(t, err)
		assert.Equal(t, "authorization_code", r.FormValue("grant_type"))
		assert.NotEmpty(t, r.FormValue("code")) // Accept any code
		assert.Equal(t, testGoogleClientID, r.FormValue("client_id"))
		assert.Equal(t, testGoogleClientSecret, r.FormValue("client_secret"))
		assert.Equal(t, testGoogleRedirectURL, r.FormValue("redirect_uri"))

		// Return mock access token
		response := map[string]interface{}{
			"access_token": "mock-access-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockTokenServer.Close()

	// Mock Google user info endpoint
	mockUserInfoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/userinfo", r.URL.Path)

		// Verify Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer mock-access-token", authHeader)

		// Return mock user info
		userInfo := map[string]interface{}{
			"id":    "google-user-123",
			"email": "googleuser@example.com",
			"name":  "Google User",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer mockUserInfoServer.Close()

	config := &OAuthConfig{
		ClientID:         testGoogleClientID,
		ClientSecret:     testGoogleClientSecret,
		RedirectURL:      testGoogleRedirectURL,
		TokenURL:         mockTokenServer.URL + "/token",
		UserInfoURL:      mockUserInfoServer.URL + "/userinfo",
	}

	oauthService := NewOAuthService(config)

	t.Run("Handle Google OAuth callback for new user", func(t *testing.T) {
		code := "test-auth-code"

		result, err := oauthService.HandleGoogleCallback(ctx, code, authService)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify auth result
		assert.NotEmpty(t, result.Token)
		assert.NotEmpty(t, result.UserID)
		assert.Equal(t, "googleuser@example.com", result.Email)
		assert.Equal(t, "Google User", result.Name)

		// Verify user was created in database
		user, err := userRepo.FindByEmail(ctx, "googleuser@example.com")
		require.NoError(t, err)
		assert.Equal(t, "googleuser@example.com", user.Email)
		assert.Equal(t, "Google User", user.Name)
		assert.Empty(t, user.PasswordHash) // OAuth users don't have password
		assert.Equal(t, "google-user-123", user.OAuthProvider)
		assert.Equal(t, "google", user.OAuthProviderName)
	})

	t.Run("Handle Google OAuth callback for existing user", func(t *testing.T) {
		// First login creates the user
		code1 := "test-auth-code"
		result1, err := oauthService.HandleGoogleCallback(ctx, code1, authService)
		require.NoError(t, err)

		// Second login should find existing user
		code2 := "test-auth-code-2"
		result2, err := oauthService.HandleGoogleCallback(ctx, code2, authService)
		require.NoError(t, err)

		// Should return same user
		assert.Equal(t, result1.UserID, result2.UserID)
		assert.Equal(t, result1.Email, result2.Email)

		// Both tokens should be valid (they might be identical if generated in same second)
		_, err = pkgauth.ValidateToken(result1.Token, testJWTSecret)
		assert.NoError(t, err)
		_, err = pkgauth.ValidateToken(result2.Token, testJWTSecret)
		assert.NoError(t, err)
	})
}

func TestGoogleOAuthErrors(t *testing.T) {
	config := &OAuthConfig{
		ClientID:     testGoogleClientID,
		ClientSecret: testGoogleClientSecret,
		RedirectURL:  testGoogleRedirectURL,
		TokenURL:     "http://invalid-url/token",
		UserInfoURL:  "http://invalid-url/userinfo",
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
		result, err := oauthService.HandleGoogleCallback(ctx, "", authService)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "authorization code is required")
	})

	t.Run("Handle token exchange failure", func(t *testing.T) {
		result, err := oauthService.HandleGoogleCallback(ctx, "invalid-code", authService)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGoogleOAuthUserInfoParsing(t *testing.T) {
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(mongo.DB)
	authService := NewService(userRepo, testJWTSecret, testTokenDuration)

	t.Run("Handle user info with minimal fields", func(t *testing.T) {
		mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"access_token": "mock-token",
				"token_type":   "Bearer",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer mockTokenServer.Close()

		mockUserInfoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userInfo := map[string]interface{}{
				"id":    "minimal-user-123",
				"email": "minimal@example.com",
				"name":  "Minimal User",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(userInfo)
		}))
		defer mockUserInfoServer.Close()

		config := &OAuthConfig{
			ClientID:     testGoogleClientID,
			ClientSecret: testGoogleClientSecret,
			RedirectURL:  testGoogleRedirectURL,
			TokenURL:     mockTokenServer.URL,
			UserInfoURL:  mockUserInfoServer.URL,
		}

		oauthService := NewOAuthService(config)
		result, err := oauthService.HandleGoogleCallback(ctx, "test-code", authService)
		require.NoError(t, err)

		assert.Equal(t, "minimal@example.com", result.Email)
		assert.Equal(t, "Minimal User", result.Name)
	})
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository/testutil"
	authservice "github.com/debMan/collaborative-hobby-tracker/backend/internal/service/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAuthJWTSecret     = "test-auth-secret-key"
	testAuthTokenDuration = 24 * time.Hour
)

// setupAuthTestHandler creates a test handler for auth endpoints
func setupAuthTestHandler(t *testing.T) (*AuthHandler, *testutil.MongoDBContainer, *gin.Engine) {
	mongo := testutil.SetupMongoDB(t)
	ctx := context.Background()

	// Create indexes
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	// Create repositories
	userRepo := repository.NewUserRepository(mongo.DB)

	// Create auth service
	authService := authservice.NewService(userRepo, testAuthJWTSecret, testAuthTokenDuration)

	// Create mock OAuth services (with invalid URLs for testing)
	googleOAuth := authservice.NewGoogleOAuthService(&authservice.OAuthConfig{
		ClientID:     "test-google-client-id",
		ClientSecret: "test-google-client-secret",
		RedirectURL:  "http://localhost:8080/api/v1/auth/google/callback",
	})

	githubOAuth := authservice.NewGitHubOAuthService(&authservice.OAuthConfig{
		ClientID:     "test-github-client-id",
		ClientSecret: "test-github-client-secret",
		RedirectURL:  "http://localhost:8080/api/v1/auth/github/callback",
	})

	// Create handler
	handler := NewAuthHandler(authService, googleOAuth, githubOAuth, testAuthJWTSecret)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Auth routes
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.GET("/google", handler.GoogleAuth)
		auth.GET("/google/callback", handler.GoogleCallback)
		auth.GET("/github", handler.GitHubAuth)
		auth.GET("/github/callback", handler.GitHubCallback)
	}

	return handler, mongo, router
}

func TestRegister(t *testing.T) {
	_, mongo, router := setupAuthTestHandler(t)
	defer mongo.Teardown(t)

	t.Run("Register successfully with valid data", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "newuser@example.com",
			"password": "password123",
			"name":     "New User",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response["token"])
		assert.NotEmpty(t, response["userId"])
		assert.Equal(t, "newuser@example.com", response["email"])
		assert.Equal(t, "New User", response["name"])
	})

	t.Run("Reject registration with existing email", func(t *testing.T) {
		// First registration
		reqBody1 := map[string]interface{}{
			"email":    "duplicate@example.com",
			"password": "password123",
			"name":     "First User",
		}
		body1, _ := json.Marshal(reqBody1)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body1))
		req1.Header.Set("Content-Type", "application/json")
		resp1 := httptest.NewRecorder()
		router.ServeHTTP(resp1, req1)
		require.Equal(t, http.StatusCreated, resp1.Code)

		// Second registration with same email
		reqBody2 := map[string]interface{}{
			"email":    "duplicate@example.com",
			"password": "password456",
			"name":     "Second User",
		}
		body2, _ := json.Marshal(reqBody2)

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		resp2 := httptest.NewRecorder()
		router.ServeHTTP(resp2, req2)

		assert.Equal(t, http.StatusConflict, resp2.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp2.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "email already exists")
	})

	t.Run("Reject registration with invalid email", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "not-an-email",
			"password": "password123",
			"name":     "Test User",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid email")
	})

	t.Run("Reject registration with short password", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "user@example.com",
			"password": "short",
			"name":     "Test User",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "at least 8 characters")
	})

	t.Run("Reject registration with missing email", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"password": "password123",
			"name":     "Test User",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject registration with missing password", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email": "user@example.com",
			"name":  "Test User",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject registration with missing name", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "user@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject registration with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestLogin(t *testing.T) {
	_, mongo, router := setupAuthTestHandler(t)
	defer mongo.Teardown(t)

	// Register a user first
	registerBody := map[string]interface{}{
		"email":    "loginuser@example.com",
		"password": "password123",
		"name":     "Login User",
	}
	body, _ := json.Marshal(registerBody)
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp := httptest.NewRecorder()
	router.ServeHTTP(registerResp, registerReq)
	require.Equal(t, http.StatusCreated, registerResp.Code)

	t.Run("Login successfully with correct credentials", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "loginuser@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response["token"])
		assert.NotEmpty(t, response["userId"])
		assert.Equal(t, "loginuser@example.com", response["email"])
		assert.Equal(t, "Login User", response["name"])
	})

	t.Run("Reject login with incorrect password", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "loginuser@example.com",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid credentials")
	})

	t.Run("Reject login with non-existent email", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid credentials")
	})

	t.Run("Reject login with missing email", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject login with missing password", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email": "loginuser@example.com",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject login with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestGoogleAuth(t *testing.T) {
	_, mongo, router := setupAuthTestHandler(t)
	defer mongo.Teardown(t)

	t.Run("Redirect to Google OAuth URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// Should redirect to Google
		assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)

		// Check Location header contains Google OAuth URL
		location := resp.Header().Get("Location")
		assert.Contains(t, location, "https://accounts.google.com/o/oauth2/v2/auth")
		assert.Contains(t, location, "client_id=")
		assert.Contains(t, location, "redirect_uri=")
		assert.Contains(t, location, "state=")
		assert.Contains(t, location, "scope=")

		// Check state cookie is set
		cookies := resp.Result().Cookies()
		var stateCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "oauth_state" {
				stateCookie = cookie
				break
			}
		}
		assert.NotNil(t, stateCookie, "oauth_state cookie should be set")
		assert.NotEmpty(t, stateCookie.Value, "oauth_state cookie should have a value")
	})
}

func TestGitHubAuth(t *testing.T) {
	_, mongo, router := setupAuthTestHandler(t)
	defer mongo.Teardown(t)

	t.Run("Redirect to GitHub OAuth URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/github", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// Should redirect to GitHub
		assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)

		// Check Location header contains GitHub OAuth URL
		location := resp.Header().Get("Location")
		assert.Contains(t, location, "https://github.com/login/oauth/authorize")
		assert.Contains(t, location, "client_id=")
		assert.Contains(t, location, "redirect_uri=")
		assert.Contains(t, location, "state=")
		assert.Contains(t, location, "scope=")

		// Check state cookie is set
		cookies := resp.Result().Cookies()
		var stateCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "oauth_state" {
				stateCookie = cookie
				break
			}
		}
		assert.NotNil(t, stateCookie, "oauth_state cookie should be set")
		assert.NotEmpty(t, stateCookie.Value, "oauth_state cookie should have a value")
	})
}

func TestGoogleCallback(t *testing.T) {
	_, mongo, router := setupAuthTestHandler(t)
	defer mongo.Teardown(t)

	t.Run("Reject callback with missing state", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?code=test-code", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid state")
	})

	t.Run("Reject callback with invalid state", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?code=test-code&state=invalid-state", nil)
		req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "different-state"})
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid state")
	})

	t.Run("Reject callback with missing code", func(t *testing.T) {
		state := "test-state-123"
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?state="+state, nil)
		req.AddCookie(&http.Cookie{Name: "oauth_state", Value: state})
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "authorization code is required")
	})
}

func TestGitHubCallback(t *testing.T) {
	_, mongo, router := setupAuthTestHandler(t)
	defer mongo.Teardown(t)

	t.Run("Reject callback with missing state", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/github/callback?code=test-code", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid state")
	})

	t.Run("Reject callback with invalid state", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/github/callback?code=test-code&state=invalid-state", nil)
		req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "different-state"})
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid state")
	})

	t.Run("Reject callback with missing code", func(t *testing.T) {
		state := "test-state-456"
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/github/callback?state="+state, nil)
		req.AddCookie(&http.Cookie{Name: "oauth_state", Value: state})
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "authorization code is required")
	})
}

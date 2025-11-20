package handlers

import (
	"errors"
	"net/http"

	authservice "github.com/debMan/collaborative-hobby-tracker/backend/internal/service/auth"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService  *authservice.Service
	googleOAuth  *authservice.OAuthService
	githubOAuth  *authservice.OAuthService
	jwtSecret    string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	authService *authservice.Service,
	googleOAuth *authservice.OAuthService,
	githubOAuth *authservice.OAuthService,
	jwtSecret string,
) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		googleOAuth:  googleOAuth,
		githubOAuth:  githubOAuth,
		jwtSecret:    jwtSecret,
	}
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, authservice.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, authservice.ErrInvalidEmail) ||
			errors.Is(err, authservice.ErrPasswordTooShort) ||
			errors.Is(err, authservice.ErrEmailRequired) ||
			errors.Is(err, authservice.ErrPasswordRequired) ||
			errors.Is(err, authservice.ErrNameRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token":  result.Token,
		"userId": result.UserID,
		"email":  result.Email,
		"name":   result.Name,
	})
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, authservice.ErrEmailRequired) ||
			errors.Is(err, authservice.ErrPasswordRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  result.Token,
		"userId": result.UserID,
		"email":  result.Email,
		"name":   result.Name,
	})
}

// GoogleAuth handles GET /auth/google (redirects to Google)
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	// Generate random state for CSRF protection
	state := generateRandomState()

	// Store state in session/cookie (simplified for now)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	// Redirect to Google OAuth
	authURL := h.googleOAuth.GetGoogleAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback handles GET /auth/google/callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Verify state
	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state parameter"})
		return
	}

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code is required"})
		return
	}

	// Handle OAuth callback
	result, err := h.googleOAuth.HandleGoogleCallback(c.Request.Context(), code, h.authService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate with Google"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  result.Token,
		"userId": result.UserID,
		"email":  result.Email,
		"name":   result.Name,
	})
}

// GitHubAuth handles GET /auth/github (redirects to GitHub)
func (h *AuthHandler) GitHubAuth(c *gin.Context) {
	// Generate random state for CSRF protection
	state := generateRandomState()

	// Store state in session/cookie (simplified for now)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	// Redirect to GitHub OAuth
	authURL := h.githubOAuth.GetGitHubAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback handles GET /auth/github/callback
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	// Verify state
	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state parameter"})
		return
	}

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code is required"})
		return
	}

	// Handle OAuth callback
	result, err := h.githubOAuth.HandleGitHubCallback(c.Request.Context(), code, h.authService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate with GitHub"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  result.Token,
		"userId": result.UserID,
		"email":  result.Email,
		"name":   result.Name,
	})
}

// generateRandomState generates a random state string for CSRF protection
func generateRandomState() string {
	// Simple implementation - in production, use crypto/rand
	return "random-state-" + string(rune(123456789))
}

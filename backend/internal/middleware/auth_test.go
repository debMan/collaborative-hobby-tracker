package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const testJWTSecret = "test-jwt-secret-key-for-middleware"

func TestAuthMiddleware(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	email := "test@example.com"

	// Generate a valid token
	validToken, err := auth.GenerateToken(userID, email, testJWTSecret, 1*time.Hour)
	assert.NoError(t, err)

	// Generate an expired token
	expiredToken, err := auth.GenerateToken(userID, email, testJWTSecret, -1*time.Hour)
	assert.NoError(t, err)

	t.Run("Allow request with valid token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/protected", func(c *gin.Context) {
			// Get user ID from context
			userIDFromContext, exists := c.Get("userID")
			assert.True(t, exists)
			assert.Equal(t, userID.Hex(), userIDFromContext)

			// Get email from context
			emailFromContext, exists := c.Get("email")
			assert.True(t, exists)
			assert.Equal(t, email, emailFromContext)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("Reject request without Authorization header", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "authorization header required")
	})

	t.Run("Reject request with malformed Authorization header", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid authorization header format")
	})

	t.Run("Reject request with missing Bearer prefix", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", validToken) // Missing "Bearer " prefix
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid authorization header format")
	})

	t.Run("Reject request with invalid token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid token")
	})

	t.Run("Reject request with expired token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid token")
	})

	t.Run("Reject request with tampered token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Tamper with the token
		tamperedToken := validToken[:len(validToken)-5] + "XXXXX"

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tamperedToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid token")
	})

	t.Run("Reject request with token signed with different secret", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Generate token with different secret
		differentToken, err := auth.GenerateToken(userID, email, "different-secret", 1*time.Hour)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+differentToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid token")
	})
}

func TestAuthMiddlewareContextInjection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	email := "context@example.com"

	token, err := auth.GenerateToken(userID, email, testJWTSecret, 1*time.Hour)
	assert.NoError(t, err)

	t.Run("Inject user information into context", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(testJWTSecret))
		router.GET("/me", func(c *gin.Context) {
			userIDFromContext, exists := c.Get("userID")
			assert.True(t, exists, "userID should exist in context")
			assert.Equal(t, userID.Hex(), userIDFromContext, "userID should match")

			emailFromContext, exists := c.Get("email")
			assert.True(t, exists, "email should exist in context")
			assert.Equal(t, email, emailFromContext, "email should match")

			c.JSON(http.StatusOK, gin.H{
				"userID": userIDFromContext,
				"email":  emailFromContext,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), userID.Hex())
		assert.Contains(t, w.Body.String(), email)
	})
}

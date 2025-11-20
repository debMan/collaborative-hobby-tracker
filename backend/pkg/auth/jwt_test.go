package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const testSecret = "test-secret-key-for-jwt-signing"

func TestGenerateToken(t *testing.T) {
	userID := primitive.NewObjectID()
	email := "test@example.com"

	t.Run("Generate valid token", func(t *testing.T) {
		token, err := GenerateToken(userID, email, testSecret, 1*time.Hour)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		// JWT format: header.payload.signature
		assert.Regexp(t, `^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`, token)
	})

	t.Run("Generate token with different users produces different tokens", func(t *testing.T) {
		user1ID := primitive.NewObjectID()
		user2ID := primitive.NewObjectID()

		token1, err1 := GenerateToken(user1ID, "user1@example.com", testSecret, 1*time.Hour)
		token2, err2 := GenerateToken(user2ID, "user2@example.com", testSecret, 1*time.Hour)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})

	t.Run("Generate token with zero duration creates expired token", func(t *testing.T) {
		token, err := GenerateToken(userID, email, testSecret, 0)

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Token should be immediately expired
		claims, err := ValidateToken(token, testSecret)
		assert.Error(t, err) // Should fail validation due to expiration
		assert.Nil(t, claims)
	})
}

func TestValidateToken(t *testing.T) {
	userID := primitive.NewObjectID()
	email := "test@example.com"

	t.Run("Validate valid token", func(t *testing.T) {
		token, err := GenerateToken(userID, email, testSecret, 1*time.Hour)
		require.NoError(t, err)

		claims, err := ValidateToken(token, testSecret)
		require.NoError(t, err)
		require.NotNil(t, claims)

		assert.Equal(t, userID.Hex(), claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.True(t, time.Unix(claims.ExpiresAt, 0).After(time.Now()))
	})

	t.Run("Reject token with wrong secret", func(t *testing.T) {
		token, err := GenerateToken(userID, email, testSecret, 1*time.Hour)
		require.NoError(t, err)

		claims, err := ValidateToken(token, "wrong-secret")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Reject malformed token", func(t *testing.T) {
		malformedToken := "this.is.not.a.valid.jwt"

		claims, err := ValidateToken(malformedToken, testSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Reject empty token", func(t *testing.T) {
		claims, err := ValidateToken("", testSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Reject expired token", func(t *testing.T) {
		// Generate token that expires immediately
		token, err := GenerateToken(userID, email, testSecret, -1*time.Hour)
		require.NoError(t, err)

		claims, err := ValidateToken(token, testSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Reject token with invalid signature", func(t *testing.T) {
		// Generate a valid token
		token, err := GenerateToken(userID, email, testSecret, 1*time.Hour)
		require.NoError(t, err)

		// Tamper with the token by changing last character
		tamperedToken := token[:len(token)-1] + "X"

		claims, err := ValidateToken(tamperedToken, testSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestTokenRoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		userID   primitive.ObjectID
		email    string
		duration time.Duration
	}{
		{
			name:     "Short-lived token",
			userID:   primitive.NewObjectID(),
			email:    "short@example.com",
			duration: 15 * time.Minute,
		},
		{
			name:     "Long-lived token",
			userID:   primitive.NewObjectID(),
			email:    "long@example.com",
			duration: 24 * time.Hour,
		},
		{
			name:     "Email with special characters",
			email:    "user+tag@example.co.uk",
			userID:   primitive.NewObjectID(),
			duration: 1 * time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate token
			token, err := GenerateToken(tc.userID, tc.email, testSecret, tc.duration)
			require.NoError(t, err)

			// Validate token
			claims, err := ValidateToken(token, testSecret)
			require.NoError(t, err)
			require.NotNil(t, claims)

			// Verify claims
			assert.Equal(t, tc.userID.Hex(), claims.UserID)
			assert.Equal(t, tc.email, claims.Email)

			// Verify expiration is approximately correct (within 1 second)
			expectedExpiry := time.Now().Add(tc.duration)
			actualExpiry := time.Unix(claims.ExpiresAt, 0)
			assert.WithinDuration(t, expectedExpiry, actualExpiry, 1*time.Second)
		})
	}
}

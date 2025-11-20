package auth

import (
	"context"
	"testing"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository/testutil"
	pkgauth "github.com/debMan/collaborative-hobby-tracker/backend/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-jwt-secret-key-for-auth-service"
const testTokenDuration = 1 * time.Hour

func TestRegister(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create indexes
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(mongo.DB)
	service := NewService(userRepo, testJWTSecret, testTokenDuration)

	t.Run("Register new user successfully", func(t *testing.T) {
		email := "newuser@example.com"
		password := "securePassword123"
		name := "New User"

		result, err := service.Register(ctx, email, password, name)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify token is returned
		assert.NotEmpty(t, result.Token)
		assert.NotEmpty(t, result.UserID)
		assert.Equal(t, email, result.Email)
		assert.Equal(t, name, result.Name)

		// Verify token is valid
		claims, err := pkgauth.ValidateToken(result.Token, testJWTSecret)
		require.NoError(t, err)
		assert.Equal(t, result.UserID, claims.UserID)
		assert.Equal(t, email, claims.Email)

		// Verify user was created in database
		user, err := userRepo.FindByEmail(ctx, email)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, name, user.Name)
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("Reject registration with duplicate email", func(t *testing.T) {
		email := "duplicate@example.com"
		password := "password123"
		name := "User 1"

		// Register first user
		_, err := service.Register(ctx, email, password, name)
		require.NoError(t, err)

		// Try to register again with same email - should fail
		result, err := service.Register(ctx, email, "differentPassword", "User 2")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "email already exists")
	})

	t.Run("Reject registration with empty email", func(t *testing.T) {
		result, err := service.Register(ctx, "", "password123", "User")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("Reject registration with empty password", func(t *testing.T) {
		result, err := service.Register(ctx, "user@example.com", "", "User")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("Reject registration with empty name", func(t *testing.T) {
		result, err := service.Register(ctx, "user@example.com", "password123", "")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("Reject registration with invalid email format", func(t *testing.T) {
		result, err := service.Register(ctx, "not-an-email", "password123", "User")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid email")
	})

	t.Run("Reject registration with short password", func(t *testing.T) {
		result, err := service.Register(ctx, "user@example.com", "short", "User")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "password must be at least")
	})
}

func TestLogin(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create indexes
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(mongo.DB)
	service := NewService(userRepo, testJWTSecret, testTokenDuration)

	// Create a test user
	email := "testuser@example.com"
	password := "testPassword123"
	name := "Test User"

	registerResult, err := service.Register(ctx, email, password, name)
	require.NoError(t, err)

	t.Run("Login with correct credentials", func(t *testing.T) {
		result, err := service.Login(ctx, email, password)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify token is returned
		assert.NotEmpty(t, result.Token)
		assert.Equal(t, registerResult.UserID, result.UserID)
		assert.Equal(t, email, result.Email)
		assert.Equal(t, name, result.Name)

		// Verify token is valid
		claims, err := pkgauth.ValidateToken(result.Token, testJWTSecret)
		require.NoError(t, err)
		assert.Equal(t, result.UserID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("Reject login with wrong password", func(t *testing.T) {
		result, err := service.Login(ctx, email, "wrongPassword")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid credentials")
	})

	t.Run("Reject login with non-existent email", func(t *testing.T) {
		result, err := service.Login(ctx, "nonexistent@example.com", password)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid credentials")
	})

	t.Run("Reject login with empty email", func(t *testing.T) {
		result, err := service.Login(ctx, "", password)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("Reject login with empty password", func(t *testing.T) {
		result, err := service.Login(ctx, email, "")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "password is required")
	})
}

func TestAuthResult(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create indexes
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(mongo.DB)
	service := NewService(userRepo, testJWTSecret, testTokenDuration)

	t.Run("Register and Login return consistent user information", func(t *testing.T) {
		email := "consistent@example.com"
		password := "password123"
		name := "Consistent User"

		// Register
		registerResult, err := service.Register(ctx, email, password, name)
		require.NoError(t, err)

		// Login
		loginResult, err := service.Login(ctx, email, password)
		require.NoError(t, err)

		// Both should return same user info
		assert.Equal(t, registerResult.UserID, loginResult.UserID)
		assert.Equal(t, registerResult.Email, loginResult.Email)
		assert.Equal(t, registerResult.Name, loginResult.Name)

		// Both tokens should be valid (they might be identical if generated in same second)
		registerClaims, err := pkgauth.ValidateToken(registerResult.Token, testJWTSecret)
		require.NoError(t, err)
		loginClaims, err := pkgauth.ValidateToken(loginResult.Token, testJWTSecret)
		require.NoError(t, err)

		// Claims should have same user info
		assert.Equal(t, registerClaims.UserID, loginClaims.UserID)
		assert.Equal(t, registerClaims.Email, loginClaims.Email)
	})
}

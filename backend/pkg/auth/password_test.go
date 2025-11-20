package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("Successfully hash a password", func(t *testing.T) {
		password := "mySecurePassword123"
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)

		// Hash should start with bcrypt prefix
		assert.Contains(t, hash, "$2a$")
	})

	t.Run("Different hashes for same password", func(t *testing.T) {
		password := "samePassword"
		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		require.NoError(t, err1)
		require.NoError(t, err2)

		// Same password should produce different hashes (salt)
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Empty password should still hash", func(t *testing.T) {
		password := ""
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Run("Verify correct password", func(t *testing.T) {
		password := "correctPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		err = VerifyPassword(hash, password)
		assert.NoError(t, err)
	})

	t.Run("Reject incorrect password", func(t *testing.T) {
		password := "correctPassword"
		wrongPassword := "wrongPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		err = VerifyPassword(hash, wrongPassword)
		assert.Error(t, err)
	})

	t.Run("Reject empty password when hash exists", func(t *testing.T) {
		password := "correctPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		err = VerifyPassword(hash, "")
		assert.Error(t, err)
	})

	t.Run("Reject invalid hash format", func(t *testing.T) {
		invalidHash := "not-a-valid-bcrypt-hash"
		password := "anyPassword"

		err := VerifyPassword(invalidHash, password)
		assert.Error(t, err)
	})

	t.Run("Verify with empty hash should error", func(t *testing.T) {
		err := VerifyPassword("", "password")
		assert.Error(t, err)
	})
}

func TestPasswordHashingRoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		password string
	}{
		{"Simple password", "password123"},
		{"Complex password", "P@ssw0rd!#$%^&*()"},
		{"Unicode password", "пароль密码🔒"},
		{"Long password", "ThisIsALongPasswordButStillUnder72BytesLimit!@#$123"},
		{"Empty password", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(tc.password)
			require.NoError(t, err)

			// Verify the correct password
			err = VerifyPassword(hash, tc.password)
			assert.NoError(t, err)

			// Verify that a different password fails
			if tc.password != "" {
				err = VerifyPassword(hash, tc.password+"wrong")
				assert.Error(t, err)
			}
		})
	}
}

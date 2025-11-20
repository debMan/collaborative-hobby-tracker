package repository

import (
	"context"
	"testing"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserRepository_Create(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewUserRepository(mongo.DB)
	ctx := context.Background()

	// Test case: Create a new user
	t.Run("Create new user successfully", func(t *testing.T) {
		user := &models.User{
			Email:        "test@example.com",
			PasswordHash: "hashed_password",
			Name:         "Test User",
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, user.ID)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	// Test case: Create user with existing ID
	t.Run("Create user with pre-set ID", func(t *testing.T) {
		id := primitive.NewObjectID()
		user := &models.User{
			ID:           id,
			Email:        "test2@example.com",
			PasswordHash: "hashed_password",
			Name:         "Test User 2",
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.Equal(t, id, user.ID)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewUserRepository(mongo.DB)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		Email:        "find@example.com",
		PasswordHash: "hashed_password",
		Name:         "Find User",
	}
	require.NoError(t, repo.Create(ctx, user))

	// Test case: Find existing user
	t.Run("Find existing user by ID", func(t *testing.T) {
		found, err := repo.FindByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)
		assert.Equal(t, user.Name, found.Name)
		assert.Equal(t, user.PasswordHash, found.PasswordHash)
	})

	// Test case: Find non-existent user
	t.Run("Find non-existent user returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		found, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewUserRepository(mongo.DB)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		Email:        "unique@example.com",
		PasswordHash: "hashed_password",
		Name:         "Email User",
	}
	require.NoError(t, repo.Create(ctx, user))

	// Test case: Find existing user by email
	t.Run("Find existing user by email", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, "unique@example.com")
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, user.Name, found.Name)
	})

	// Test case: Find non-existent user by email
	t.Run("Find non-existent user by email returns error", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, "nonexistent@example.com")
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestUserRepository_Update(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewUserRepository(mongo.DB)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		Email:        "update@example.com",
		PasswordHash: "hashed_password",
		Name:         "Original Name",
	}
	require.NoError(t, repo.Create(ctx, user))

	// Test case: Update existing user
	t.Run("Update existing user successfully", func(t *testing.T) {
		// Save original updated time
		originalUpdatedAt := user.UpdatedAt

		// Wait a bit to ensure timestamp changes
		time.Sleep(10 * time.Millisecond)

		// Update user
		user.Name = "Updated Name"
		user.AvatarURL = "https://example.com/avatar.jpg"

		err := repo.Update(ctx, user)
		require.NoError(t, err)
		assert.True(t, user.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		found, err := repo.FindByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.Equal(t, "https://example.com/avatar.jpg", found.AvatarURL)
	})

	// Test case: Update non-existent user
	t.Run("Update non-existent user returns error", func(t *testing.T) {
		nonExistentUser := &models.User{
			ID:    primitive.NewObjectID(),
			Email: "ghost@example.com",
			Name:  "Ghost User",
		}

		err := repo.Update(ctx, nonExistentUser)
		assert.Error(t, err)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewUserRepository(mongo.DB)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		Email:        "delete@example.com",
		PasswordHash: "hashed_password",
		Name:         "Delete User",
	}
	require.NoError(t, repo.Create(ctx, user))

	// Test case: Delete existing user
	t.Run("Delete existing user successfully", func(t *testing.T) {
		err := repo.Delete(ctx, user.ID)
		require.NoError(t, err)

		// Verify deletion
		found, err := repo.FindByID(ctx, user.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	// Test case: Delete non-existent user
	t.Run("Delete non-existent user returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestUserRepository_Integration(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewUserRepository(mongo.DB)
	ctx := context.Background()

	// Integration test: Full lifecycle
	t.Run("Full user lifecycle", func(t *testing.T) {
		// 1. Create user
		user := &models.User{
			Email:         "lifecycle@example.com",
			PasswordHash:  "hashed_password",
			Name:          "Lifecycle User",
			OAuthProvider: "google",
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// 2. Find by ID
		found, err := repo.FindByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)

		// 3. Find by Email
		foundByEmail, err := repo.FindByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.Equal(t, user.ID, foundByEmail.ID)

		// 4. Update
		user.Name = "Updated Lifecycle User"
		user.AvatarURL = "https://example.com/new-avatar.jpg"
		err = repo.Update(ctx, user)
		require.NoError(t, err)

		// 5. Verify update
		updated, err := repo.FindByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Lifecycle User", updated.Name)
		assert.Equal(t, "https://example.com/new-avatar.jpg", updated.AvatarURL)

		// 6. Delete
		err = repo.Delete(ctx, user.ID)
		require.NoError(t, err)

		// 7. Verify deletion
		deleted, err := repo.FindByID(ctx, user.ID)
		assert.Error(t, err)
		assert.Nil(t, deleted)
	})
}

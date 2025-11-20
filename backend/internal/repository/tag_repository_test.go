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

func TestTagRepository_Create(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewTagRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	userID := primitive.NewObjectID()

	// Test case: Create a new tag
	t.Run("Create new tag successfully", func(t *testing.T) {
		tag := &models.Tag{
			Name:       "Action",
			Color:      "#FF5733",
			UserID:     userID,
			UsageCount: 0,
		}

		err := repo.Create(ctx, tag)
		require.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, tag.ID)
		assert.False(t, tag.CreatedAt.IsZero())
		assert.False(t, tag.UpdatedAt.IsZero())
	})

	// Test case: Create tag without color
	t.Run("Create tag without color", func(t *testing.T) {
		tag := &models.Tag{
			Name:   "Comedy",
			UserID: userID,
		}

		err := repo.Create(ctx, tag)
		require.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, tag.ID)
	})
}

func TestTagRepository_FindByID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewTagRepository(mongo.DB)
	ctx := context.Background()

	// Create test tag
	userID := primitive.NewObjectID()
	tag := &models.Tag{
		Name:   "Drama",
		Color:  "#3498DB",
		UserID: userID,
	}
	require.NoError(t, repo.Create(ctx, tag))

	// Test case: Find existing tag
	t.Run("Find existing tag by ID", func(t *testing.T) {
		found, err := repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, tag.Name, found.Name)
		assert.Equal(t, tag.Color, found.Color)
		assert.Equal(t, tag.UserID, found.UserID)
	})

	// Test case: Find non-existent tag
	t.Run("Find non-existent tag returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		found, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestTagRepository_FindByUserID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewTagRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()

	// Create tags for user1
	tags1 := []*models.Tag{
		{
			Name:   "Action",
			Color:  "#FF5733",
			UserID: user1ID,
		},
		{
			Name:   "Comedy",
			Color:  "#FFC300",
			UserID: user1ID,
		},
		{
			Name:   "Drama",
			Color:  "#3498DB",
			UserID: user1ID,
		},
	}

	for _, tag := range tags1 {
		require.NoError(t, repo.Create(ctx, tag))
	}

	// Create tag for user2
	tag2 := &models.Tag{
		Name:   "User 2 Tag",
		UserID: user2ID,
	}
	require.NoError(t, repo.Create(ctx, tag2))

	// Test case: Find all tags for user1
	t.Run("Find all tags for user", func(t *testing.T) {
		found, err := repo.FindByUserID(ctx, user1ID)
		require.NoError(t, err)
		assert.Len(t, found, 3)

		// Verify all tags belong to user1
		for _, tag := range found {
			assert.Equal(t, user1ID, tag.UserID)
		}
	})

	// Test case: Find tags for user with no tags
	t.Run("Find tags for user with no tags", func(t *testing.T) {
		emptyUserID := primitive.NewObjectID()
		found, err := repo.FindByUserID(ctx, emptyUserID)
		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestTagRepository_FindByName(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewTagRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()

	// Create tag for user1
	tag1 := &models.Tag{
		Name:   "SciFi",
		Color:  "#9B59B6",
		UserID: user1ID,
	}
	require.NoError(t, repo.Create(ctx, tag1))

	// Create tag with same name for user2
	tag2 := &models.Tag{
		Name:   "SciFi",
		Color:  "#E74C3C",
		UserID: user2ID,
	}
	require.NoError(t, repo.Create(ctx, tag2))

	// Test case: Find tag by name for specific user
	t.Run("Find tag by name for user", func(t *testing.T) {
		found, err := repo.FindByName(ctx, user1ID, "SciFi")
		require.NoError(t, err)
		assert.Equal(t, tag1.Name, found.Name)
		assert.Equal(t, tag1.Color, found.Color)
		assert.Equal(t, user1ID, found.UserID)
	})

	// Test case: Find non-existent tag by name
	t.Run("Find non-existent tag by name returns error", func(t *testing.T) {
		found, err := repo.FindByName(ctx, user1ID, "NonExistent")
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	// Test case: Find tag for different user
	t.Run("Find tag by name for different user", func(t *testing.T) {
		found, err := repo.FindByName(ctx, user2ID, "SciFi")
		require.NoError(t, err)
		assert.Equal(t, user2ID, found.UserID)
		assert.Equal(t, tag2.Color, found.Color)
	})
}

func TestTagRepository_Update(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewTagRepository(mongo.DB)
	ctx := context.Background()

	// Create test tag
	userID := primitive.NewObjectID()
	tag := &models.Tag{
		Name:       "Original Name",
		Color:      "#FF0000",
		UserID:     userID,
		UsageCount: 5,
	}
	require.NoError(t, repo.Create(ctx, tag))

	// Test case: Update existing tag
	t.Run("Update existing tag successfully", func(t *testing.T) {
		// Save original updated time
		originalUpdatedAt := tag.UpdatedAt

		// Wait a bit to ensure timestamp changes
		time.Sleep(10 * time.Millisecond)

		// Update tag
		tag.Name = "Updated Name"
		tag.Color = "#00FF00"
		tag.UsageCount = 10

		err := repo.Update(ctx, tag)
		require.NoError(t, err)
		assert.True(t, tag.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		found, err := repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.Equal(t, "#00FF00", found.Color)
		assert.Equal(t, 10, found.UsageCount)
	})

	// Test case: Update non-existent tag
	t.Run("Update non-existent tag returns error", func(t *testing.T) {
		nonExistentTag := &models.Tag{
			ID:     primitive.NewObjectID(),
			Name:   "Ghost Tag",
			UserID: userID,
		}

		err := repo.Update(ctx, nonExistentTag)
		assert.Error(t, err)
	})
}

func TestTagRepository_Delete(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewTagRepository(mongo.DB)
	ctx := context.Background()

	// Create test tag
	userID := primitive.NewObjectID()
	tag := &models.Tag{
		Name:   "Delete Me",
		UserID: userID,
	}
	require.NoError(t, repo.Create(ctx, tag))

	// Test case: Delete existing tag
	t.Run("Delete existing tag successfully", func(t *testing.T) {
		err := repo.Delete(ctx, tag.ID)
		require.NoError(t, err)

		// Verify deletion
		found, err := repo.FindByID(ctx, tag.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	// Test case: Delete non-existent tag
	t.Run("Delete non-existent tag returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestTagRepository_IncrementUsage(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewTagRepository(mongo.DB)
	ctx := context.Background()

	// Create test tag
	userID := primitive.NewObjectID()
	tag := &models.Tag{
		Name:       "Usage Test",
		UserID:     userID,
		UsageCount: 0,
	}
	require.NoError(t, repo.Create(ctx, tag))

	// Test case: Increment usage count
	t.Run("Increment usage count successfully", func(t *testing.T) {
		err := repo.IncrementUsage(ctx, tag.ID)
		require.NoError(t, err)

		// Verify increment
		found, err := repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, found.UsageCount)

		// Increment again
		err = repo.IncrementUsage(ctx, tag.ID)
		require.NoError(t, err)

		found, err = repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, found.UsageCount)

		// Increment multiple times
		for i := 0; i < 5; i++ {
			err = repo.IncrementUsage(ctx, tag.ID)
			require.NoError(t, err)
		}

		found, err = repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, 7, found.UsageCount)
	})

	// Test case: Increment non-existent tag
	t.Run("Increment non-existent tag returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.IncrementUsage(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestTagRepository_Integration(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewTagRepository(mongo.DB)
	ctx := context.Background()

	// Integration test: Full lifecycle
	t.Run("Full tag lifecycle", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// 1. Create tag
		tag := &models.Tag{
			Name:       "Thriller",
			Color:      "#8E44AD",
			UserID:     userID,
			UsageCount: 0,
		}
		err := repo.Create(ctx, tag)
		require.NoError(t, err)

		// 2. Find by ID
		found, err := repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, tag.Name, found.Name)

		// 3. Find by UserID
		userTags, err := repo.FindByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, userTags, 1)

		// 4. Find by Name
		foundByName, err := repo.FindByName(ctx, userID, "Thriller")
		require.NoError(t, err)
		assert.Equal(t, tag.ID, foundByName.ID)

		// 5. Increment usage
		err = repo.IncrementUsage(ctx, tag.ID)
		require.NoError(t, err)

		incremented, err := repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, incremented.UsageCount)

		// 6. Increment multiple times
		for i := 0; i < 5; i++ {
			err = repo.IncrementUsage(ctx, tag.ID)
			require.NoError(t, err)
		}

		multiIncremented, err := repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, 6, multiIncremented.UsageCount)

		// 7. Update tag
		tag.Name = "Psychological Thriller"
		tag.Color = "#2C3E50"
		err = repo.Update(ctx, tag)
		require.NoError(t, err)

		// 8. Verify update
		updated, err := repo.FindByID(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, "Psychological Thriller", updated.Name)
		assert.Equal(t, "#2C3E50", updated.Color)

		// 9. Find by new name
		foundByNewName, err := repo.FindByName(ctx, userID, "Psychological Thriller")
		require.NoError(t, err)
		assert.Equal(t, tag.ID, foundByNewName.ID)

		// 10. Delete
		err = repo.Delete(ctx, tag.ID)
		require.NoError(t, err)

		// 11. Verify deletion
		deleted, err := repo.FindByID(ctx, tag.ID)
		assert.Error(t, err)
		assert.Nil(t, deleted)
	})
}

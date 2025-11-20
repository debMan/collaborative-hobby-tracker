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

func TestHobbyItemRepository_Create(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewHobbyItemRepository(mongo.DB)
	ctx := context.Background()

	// Create test user and category
	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	// Test case: Create a new item
	t.Run("Create new item successfully", func(t *testing.T) {
		item := &models.HobbyItem{
			Title:       "Watch Inception",
			Description: "A mind-bending thriller",
			CategoryID:  categoryID,
			OwnerID:     userID,
			IsCompleted: false,
			Source:      models.SourceManual,
			Tags:        []string{"movie", "thriller"},
		}

		err := repo.Create(ctx, item)
		require.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, item.ID)
		assert.False(t, item.AddedAt.IsZero())
		assert.False(t, item.CreatedAt.IsZero())
		assert.False(t, item.UpdatedAt.IsZero())
	})

	// Test case: Create item with pre-set ID
	t.Run("Create item with pre-set ID", func(t *testing.T) {
		id := primitive.NewObjectID()
		item := &models.HobbyItem{
			ID:          id,
			Title:       "Visit Tokyo",
			CategoryID:  categoryID,
			OwnerID:     userID,
			IsCompleted: false,
			Source:      models.SourceManual,
		}

		err := repo.Create(ctx, item)
		require.NoError(t, err)
		assert.Equal(t, id, item.ID)
	})

	// Test case: Create item with metadata
	t.Run("Create item with metadata", func(t *testing.T) {
		item := &models.HobbyItem{
			Title:       "Try Sushi",
			CategoryID:  categoryID,
			OwnerID:     userID,
			Source:      models.SourceInstagram,
			SourceURL:   "https://instagram.com/post/123",
			ImageURL:    "https://example.com/sushi.jpg",
			Metadata: map[string]interface{}{
				"restaurant": "Nobu",
				"rating":     4.5,
			},
		}

		err := repo.Create(ctx, item)
		require.NoError(t, err)
		assert.NotNil(t, item.Metadata)
		assert.Equal(t, "Nobu", item.Metadata["restaurant"])
	})
}

func TestHobbyItemRepository_FindByID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewHobbyItemRepository(mongo.DB)
	ctx := context.Background()

	// Create test item
	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	item := &models.HobbyItem{
		Title:       "Watch The Matrix",
		Description: "A cyberpunk classic",
		CategoryID:  categoryID,
		OwnerID:     userID,
		Source:      models.SourceManual,
		Tags:        []string{"movie", "sci-fi"},
	}
	require.NoError(t, repo.Create(ctx, item))

	// Test case: Find existing item
	t.Run("Find existing item by ID", func(t *testing.T) {
		found, err := repo.FindByID(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, item.Title, found.Title)
		assert.Equal(t, item.Description, found.Description)
		assert.Equal(t, item.CategoryID, found.CategoryID)
		assert.Equal(t, item.OwnerID, found.OwnerID)
		assert.Equal(t, item.Tags, found.Tags)
	})

	// Test case: Find non-existent item
	t.Run("Find non-existent item returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		found, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestHobbyItemRepository_FindByUserID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewHobbyItemRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	// Create items for user1
	items1 := []*models.HobbyItem{
		{
			Title:      "Item 1",
			CategoryID: categoryID,
			OwnerID:    user1ID,
			Source:     models.SourceManual,
		},
		{
			Title:      "Item 2",
			CategoryID: categoryID,
			OwnerID:    user1ID,
			Source:     models.SourceManual,
		},
		{
			Title:      "Item 3",
			CategoryID: categoryID,
			OwnerID:    user1ID,
			Source:     models.SourceManual,
		},
	}

	for _, item := range items1 {
		require.NoError(t, repo.Create(ctx, item))
	}

	// Create item for user2
	item2 := &models.HobbyItem{
		Title:      "User 2 Item",
		CategoryID: categoryID,
		OwnerID:    user2ID,
		Source:     models.SourceManual,
	}
	require.NoError(t, repo.Create(ctx, item2))

	// Test case: Find all items for user1
	t.Run("Find all items for user", func(t *testing.T) {
		found, err := repo.FindByUserID(ctx, user1ID)
		require.NoError(t, err)
		assert.Len(t, found, 3)

		// Verify all items belong to user1
		for _, item := range found {
			assert.Equal(t, user1ID, item.OwnerID)
		}
	})

	// Test case: Find items for user with no items
	t.Run("Find items for user with no items", func(t *testing.T) {
		emptyUserID := primitive.NewObjectID()
		found, err := repo.FindByUserID(ctx, emptyUserID)
		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestHobbyItemRepository_FindByCategoryID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewHobbyItemRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	userID := primitive.NewObjectID()
	category1ID := primitive.NewObjectID()
	category2ID := primitive.NewObjectID()

	// Create items for category1
	items1 := []*models.HobbyItem{
		{
			Title:      "Movie 1",
			CategoryID: category1ID,
			OwnerID:    userID,
			Source:     models.SourceManual,
		},
		{
			Title:      "Movie 2",
			CategoryID: category1ID,
			OwnerID:    userID,
			Source:     models.SourceManual,
		},
	}

	for _, item := range items1 {
		require.NoError(t, repo.Create(ctx, item))
	}

	// Create item for category2
	item2 := &models.HobbyItem{
		Title:      "Restaurant 1",
		CategoryID: category2ID,
		OwnerID:    userID,
		Source:     models.SourceManual,
	}
	require.NoError(t, repo.Create(ctx, item2))

	// Test case: Find all items in category1
	t.Run("Find all items in category", func(t *testing.T) {
		found, err := repo.FindByCategoryID(ctx, category1ID)
		require.NoError(t, err)
		assert.Len(t, found, 2)

		// Verify all items belong to category1
		for _, item := range found {
			assert.Equal(t, category1ID, item.CategoryID)
		}
	})

	// Test case: Find items in empty category
	t.Run("Find items in empty category", func(t *testing.T) {
		emptyCategoryID := primitive.NewObjectID()
		found, err := repo.FindByCategoryID(ctx, emptyCategoryID)
		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestHobbyItemRepository_Update(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewHobbyItemRepository(mongo.DB)
	ctx := context.Background()

	// Create test item
	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	item := &models.HobbyItem{
		Title:       "Original Title",
		Description: "Original Description",
		CategoryID:  categoryID,
		OwnerID:     userID,
		Source:      models.SourceManual,
		Tags:        []string{"original"},
	}
	require.NoError(t, repo.Create(ctx, item))

	// Test case: Update existing item
	t.Run("Update existing item successfully", func(t *testing.T) {
		// Save original updated time
		originalUpdatedAt := item.UpdatedAt

		// Wait a bit to ensure timestamp changes
		time.Sleep(10 * time.Millisecond)

		// Update item
		item.Title = "Updated Title"
		item.Description = "Updated Description"
		item.Tags = []string{"updated", "new"}
		item.ImageURL = "https://example.com/image.jpg"

		err := repo.Update(ctx, item)
		require.NoError(t, err)
		assert.True(t, item.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		found, err := repo.FindByID(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", found.Title)
		assert.Equal(t, "Updated Description", found.Description)
		assert.Equal(t, []string{"updated", "new"}, found.Tags)
		assert.Equal(t, "https://example.com/image.jpg", found.ImageURL)
	})

	// Test case: Update non-existent item
	t.Run("Update non-existent item returns error", func(t *testing.T) {
		nonExistentItem := &models.HobbyItem{
			ID:         primitive.NewObjectID(),
			Title:      "Ghost Item",
			CategoryID: categoryID,
			OwnerID:    userID,
			Source:     models.SourceManual,
		}

		err := repo.Update(ctx, nonExistentItem)
		assert.Error(t, err)
	})
}

func TestHobbyItemRepository_Delete(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewHobbyItemRepository(mongo.DB)
	ctx := context.Background()

	// Create test item
	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	item := &models.HobbyItem{
		Title:      "Delete Me",
		CategoryID: categoryID,
		OwnerID:    userID,
		Source:     models.SourceManual,
	}
	require.NoError(t, repo.Create(ctx, item))

	// Test case: Delete existing item
	t.Run("Delete existing item successfully", func(t *testing.T) {
		err := repo.Delete(ctx, item.ID)
		require.NoError(t, err)

		// Verify deletion
		found, err := repo.FindByID(ctx, item.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	// Test case: Delete non-existent item
	t.Run("Delete non-existent item returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestHobbyItemRepository_ToggleComplete(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewHobbyItemRepository(mongo.DB)
	ctx := context.Background()

	// Create test item
	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	item := &models.HobbyItem{
		Title:       "Toggle Me",
		CategoryID:  categoryID,
		OwnerID:     userID,
		Source:      models.SourceManual,
		IsCompleted: false,
	}
	require.NoError(t, repo.Create(ctx, item))

	// Test case: Toggle to completed
	t.Run("Toggle item to completed", func(t *testing.T) {
		err := repo.ToggleComplete(ctx, item.ID)
		require.NoError(t, err)

		// Verify completion
		found, err := repo.FindByID(ctx, item.ID)
		require.NoError(t, err)
		assert.True(t, found.IsCompleted)
		assert.NotNil(t, found.CompletedAt)
	})

	// Test case: Toggle to incomplete
	t.Run("Toggle item to incomplete", func(t *testing.T) {
		err := repo.ToggleComplete(ctx, item.ID)
		require.NoError(t, err)

		// Verify incompletion
		found, err := repo.FindByID(ctx, item.ID)
		require.NoError(t, err)
		assert.False(t, found.IsCompleted)
		assert.Nil(t, found.CompletedAt)
	})

	// Test case: Toggle non-existent item
	t.Run("Toggle non-existent item returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.ToggleComplete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestHobbyItemRepository_Integration(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewHobbyItemRepository(mongo.DB)
	ctx := context.Background()

	// Integration test: Full lifecycle
	t.Run("Full item lifecycle", func(t *testing.T) {
		userID := primitive.NewObjectID()
		categoryID := primitive.NewObjectID()

		// 1. Create item
		item := &models.HobbyItem{
			Title:       "Watch Interstellar",
			Description: "Epic space movie",
			CategoryID:  categoryID,
			OwnerID:     userID,
			Source:      models.SourceYouTube,
			SourceURL:   "https://youtube.com/watch?v=xyz",
			Tags:        []string{"movie", "space", "sci-fi"},
			Metadata: map[string]interface{}{
				"director": "Christopher Nolan",
				"year":     2014,
			},
		}
		err := repo.Create(ctx, item)
		require.NoError(t, err)

		// 2. Find by ID
		found, err := repo.FindByID(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, item.Title, found.Title)

		// 3. Find by UserID
		userItems, err := repo.FindByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, userItems, 1)

		// 4. Find by CategoryID
		categoryItems, err := repo.FindByCategoryID(ctx, categoryID)
		require.NoError(t, err)
		assert.Len(t, categoryItems, 1)

		// 5. Update
		item.Title = "Watch Interstellar (IMAX)"
		item.Tags = append(item.Tags, "imax")
		err = repo.Update(ctx, item)
		require.NoError(t, err)

		// 6. Verify update
		updated, err := repo.FindByID(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, "Watch Interstellar (IMAX)", updated.Title)
		assert.Contains(t, updated.Tags, "imax")

		// 7. Toggle complete
		err = repo.ToggleComplete(ctx, item.ID)
		require.NoError(t, err)

		completed, err := repo.FindByID(ctx, item.ID)
		require.NoError(t, err)
		assert.True(t, completed.IsCompleted)
		assert.NotNil(t, completed.CompletedAt)

		// 8. Toggle incomplete
		err = repo.ToggleComplete(ctx, item.ID)
		require.NoError(t, err)

		incomplete, err := repo.FindByID(ctx, item.ID)
		require.NoError(t, err)
		assert.False(t, incomplete.IsCompleted)
		assert.Nil(t, incomplete.CompletedAt)

		// 9. Delete
		err = repo.Delete(ctx, item.ID)
		require.NoError(t, err)

		// 10. Verify deletion
		deleted, err := repo.FindByID(ctx, item.ID)
		assert.Error(t, err)
		assert.Nil(t, deleted)
	})
}

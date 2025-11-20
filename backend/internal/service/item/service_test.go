package item

import (
	"context"
	"testing"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// setupTestService creates a test service with MongoDB
func setupTestService(t *testing.T) (*Service, *testutil.MongoDBContainer, primitive.ObjectID, primitive.ObjectID) {
	mongo := testutil.SetupMongoDB(t)
	ctx := context.Background()

	// Create indexes
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	itemRepo := repository.NewHobbyItemRepository(mongo.DB)
	service := NewService(itemRepo)

	// Create test user IDs
	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	return service, mongo, userID, categoryID
}

func TestCreateItem(t *testing.T) {
	service, mongo, userID, categoryID := setupTestService(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	t.Run("Create item successfully", func(t *testing.T) {
		req := &CreateItemRequest{
			Title:       "Watch Inception",
			Description: "Amazing sci-fi movie",
			CategoryID:  categoryID,
			Source:      models.SourceManual,
			Tags:        []string{"movie", "sci-fi"},
		}

		item, err := service.CreateItem(ctx, userID, req)
		require.NoError(t, err)
		require.NotNil(t, item)

		// Verify item properties
		assert.NotEqual(t, primitive.NilObjectID, item.ID)
		assert.Equal(t, "Watch Inception", item.Title)
		assert.Equal(t, "Amazing sci-fi movie", item.Description)
		assert.Equal(t, categoryID, item.CategoryID)
		assert.Equal(t, userID, item.OwnerID)
		assert.False(t, item.IsCompleted)
		assert.Equal(t, models.SourceManual, item.Source)
		assert.Equal(t, []string{"movie", "sci-fi"}, item.Tags)
		assert.NotZero(t, item.AddedAt)
		assert.NotZero(t, item.CreatedAt)
		assert.NotZero(t, item.UpdatedAt)
	})

	t.Run("Create item with minimal fields", func(t *testing.T) {
		req := &CreateItemRequest{
			Title:      "Minimal item",
			CategoryID: categoryID,
		}

		item, err := service.CreateItem(ctx, userID, req)
		require.NoError(t, err)
		assert.Equal(t, "Minimal item", item.Title)
		assert.Empty(t, item.Description)
		assert.Equal(t, models.SourceManual, item.Source) // Default source
	})

	t.Run("Reject empty title", func(t *testing.T) {
		req := &CreateItemRequest{
			Title:      "",
			CategoryID: categoryID,
		}

		item, err := service.CreateItem(ctx, userID, req)
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("Reject empty category ID", func(t *testing.T) {
		req := &CreateItemRequest{
			Title:      "Test item",
			CategoryID: primitive.NilObjectID,
		}

		item, err := service.CreateItem(ctx, userID, req)
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "category ID is required")
	})

	t.Run("Reject empty user ID", func(t *testing.T) {
		req := &CreateItemRequest{
			Title:      "Test item",
			CategoryID: categoryID,
		}

		item, err := service.CreateItem(ctx, primitive.NilObjectID, req)
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "user ID is required")
	})
}

func TestGetItemByID(t *testing.T) {
	service, mongo, userID, categoryID := setupTestService(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create a test item
	req := &CreateItemRequest{
		Title:      "Test Item",
		CategoryID: categoryID,
	}
	createdItem, err := service.CreateItem(ctx, userID, req)
	require.NoError(t, err)

	t.Run("Get item successfully", func(t *testing.T) {
		item, err := service.GetItemByID(ctx, userID, createdItem.ID)
		require.NoError(t, err)
		assert.Equal(t, createdItem.ID, item.ID)
		assert.Equal(t, createdItem.Title, item.Title)
	})

	t.Run("Reject access to item by different user", func(t *testing.T) {
		otherUserID := primitive.NewObjectID()
		item, err := service.GetItemByID(ctx, otherUserID, createdItem.ID)
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("Return error for non-existent item", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		item, err := service.GetItemByID(ctx, userID, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, item)
	})
}

func TestGetUserItems(t *testing.T) {
	service, mongo, userID, categoryID := setupTestService(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create multiple items for the user
	items := []string{"Item 1", "Item 2", "Item 3"}
	for _, title := range items {
		req := &CreateItemRequest{
			Title:      title,
			CategoryID: categoryID,
		}
		_, err := service.CreateItem(ctx, userID, req)
		require.NoError(t, err)
	}

	// Create item for another user
	otherUserID := primitive.NewObjectID()
	req := &CreateItemRequest{
		Title:      "Other user item",
		CategoryID: categoryID,
	}
	_, err := service.CreateItem(ctx, otherUserID, req)
	require.NoError(t, err)

	t.Run("Get all items for user", func(t *testing.T) {
		userItems, err := service.GetUserItems(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, userItems, 3)

		// Verify all items belong to the user
		for _, item := range userItems {
			assert.Equal(t, userID, item.OwnerID)
		}
	})

	t.Run("Return empty list for user with no items", func(t *testing.T) {
		emptyUserID := primitive.NewObjectID()
		userItems, err := service.GetUserItems(ctx, emptyUserID)
		require.NoError(t, err)
		assert.Empty(t, userItems)
	})
}

func TestUpdateItem(t *testing.T) {
	service, mongo, userID, categoryID := setupTestService(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create a test item
	req := &CreateItemRequest{
		Title:      "Original Title",
		CategoryID: categoryID,
	}
	createdItem, err := service.CreateItem(ctx, userID, req)
	require.NoError(t, err)

	t.Run("Update item successfully", func(t *testing.T) {
		updateReq := &UpdateItemRequest{
			Title:       "Updated Title",
			Description: "Updated description",
			Tags:        []string{"updated", "tags"},
		}

		updatedItem, err := service.UpdateItem(ctx, userID, createdItem.ID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updatedItem.Title)
		assert.Equal(t, "Updated description", updatedItem.Description)
		assert.Equal(t, []string{"updated", "tags"}, updatedItem.Tags)
		assert.True(t, updatedItem.UpdatedAt.After(createdItem.UpdatedAt))
	})

	t.Run("Reject update by different user", func(t *testing.T) {
		otherUserID := primitive.NewObjectID()
		updateReq := &UpdateItemRequest{
			Title: "Unauthorized update",
		}

		updatedItem, err := service.UpdateItem(ctx, otherUserID, createdItem.ID, updateReq)
		assert.Error(t, err)
		assert.Nil(t, updatedItem)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("Reject empty title", func(t *testing.T) {
		updateReq := &UpdateItemRequest{
			Title: "",
		}

		updatedItem, err := service.UpdateItem(ctx, userID, createdItem.ID, updateReq)
		assert.Error(t, err)
		assert.Nil(t, updatedItem)
		assert.Contains(t, err.Error(), "title is required")
	})
}

func TestDeleteItem(t *testing.T) {
	service, mongo, userID, categoryID := setupTestService(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create a test item
	req := &CreateItemRequest{
		Title:      "Item to delete",
		CategoryID: categoryID,
	}
	createdItem, err := service.CreateItem(ctx, userID, req)
	require.NoError(t, err)

	t.Run("Delete item successfully", func(t *testing.T) {
		err := service.DeleteItem(ctx, userID, createdItem.ID)
		require.NoError(t, err)

		// Verify item is deleted
		item, err := service.GetItemByID(ctx, userID, createdItem.ID)
		assert.Error(t, err)
		assert.Nil(t, item)
	})

	t.Run("Reject delete by different user", func(t *testing.T) {
		// Create another item
		req := &CreateItemRequest{
			Title:      "Another item",
			CategoryID: categoryID,
		}
		anotherItem, err := service.CreateItem(ctx, userID, req)
		require.NoError(t, err)

		// Try to delete with different user
		otherUserID := primitive.NewObjectID()
		err = service.DeleteItem(ctx, otherUserID, anotherItem.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")

		// Verify item still exists
		item, err := service.GetItemByID(ctx, userID, anotherItem.ID)
		assert.NoError(t, err)
		assert.NotNil(t, item)
	})

	t.Run("Return error for non-existent item", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := service.DeleteItem(ctx, userID, nonExistentID)
		assert.Error(t, err)
	})
}

func TestToggleItemCompletion(t *testing.T) {
	service, mongo, userID, categoryID := setupTestService(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	t.Run("Mark item as completed", func(t *testing.T) {
		// Create a new item for this test
		req := &CreateItemRequest{
			Title:      "Item to complete",
			CategoryID: categoryID,
		}
		createdItem, err := service.CreateItem(ctx, userID, req)
		require.NoError(t, err)
		assert.False(t, createdItem.IsCompleted)

		// Toggle to completed
		updatedItem, err := service.ToggleItemCompletion(ctx, userID, createdItem.ID)
		require.NoError(t, err)
		assert.True(t, updatedItem.IsCompleted)
		assert.NotNil(t, updatedItem.CompletedAt)
		assert.True(t, updatedItem.CompletedAt.After(time.Time{}))
	})

	t.Run("Mark completed item as incomplete", func(t *testing.T) {
		// Create a new item for this test
		req := &CreateItemRequest{
			Title:      "Item to toggle",
			CategoryID: categoryID,
		}
		newItem, err := service.CreateItem(ctx, userID, req)
		require.NoError(t, err)

		// First, mark as completed
		completedItem, err := service.ToggleItemCompletion(ctx, userID, newItem.ID)
		require.NoError(t, err)
		assert.True(t, completedItem.IsCompleted)

		// Then, mark as incomplete
		incompleteItem, err := service.ToggleItemCompletion(ctx, userID, newItem.ID)
		require.NoError(t, err)
		assert.False(t, incompleteItem.IsCompleted)
		assert.Nil(t, incompleteItem.CompletedAt)
	})

	t.Run("Reject toggle by different user", func(t *testing.T) {
		// Create a new item for this test
		req := &CreateItemRequest{
			Title:      "Item for auth test",
			CategoryID: categoryID,
		}
		authTestItem, err := service.CreateItem(ctx, userID, req)
		require.NoError(t, err)

		otherUserID := primitive.NewObjectID()
		updatedItem, err := service.ToggleItemCompletion(ctx, otherUserID, authTestItem.ID)
		assert.Error(t, err)
		assert.Nil(t, updatedItem)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("Return error for non-existent item", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		updatedItem, err := service.ToggleItemCompletion(ctx, userID, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, updatedItem)
	})
}

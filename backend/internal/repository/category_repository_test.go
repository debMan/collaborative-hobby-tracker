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

func TestCategoryRepository_Create(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	ownerID := primitive.NewObjectID()
	circleID := primitive.NewObjectID()

	// Test case: Create a new category
	t.Run("Create new category successfully", func(t *testing.T) {
		category := &models.Category{
			Name:      "Movies",
			Icon:      "film",
			CircleID:  circleID,
			OwnerID:   ownerID,
			ItemCount: 0,
		}

		err := repo.Create(ctx, category)
		require.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, category.ID)
		assert.False(t, category.CreatedAt.IsZero())
		assert.False(t, category.UpdatedAt.IsZero())
	})

	// Test case: Create category with pre-set ID
	t.Run("Create category with pre-set ID", func(t *testing.T) {
		id := primitive.NewObjectID()
		category := &models.Category{
			ID:       id,
			Name:     "Restaurants",
			Icon:     "utensils",
			CircleID: circleID,
			OwnerID:  ownerID,
		}

		err := repo.Create(ctx, category)
		require.NoError(t, err)
		assert.Equal(t, id, category.ID)
	})
}

func TestCategoryRepository_FindByID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Create test category
	ownerID := primitive.NewObjectID()
	circleID := primitive.NewObjectID()
	category := &models.Category{
		Name:     "Travel",
		Icon:     "plane",
		CircleID: circleID,
		OwnerID:  ownerID,
	}
	require.NoError(t, repo.Create(ctx, category))

	// Test case: Find existing category
	t.Run("Find existing category by ID", func(t *testing.T) {
		found, err := repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, category.Name, found.Name)
		assert.Equal(t, category.Icon, found.Icon)
		assert.Equal(t, category.CircleID, found.CircleID)
		assert.Equal(t, category.OwnerID, found.OwnerID)
	})

	// Test case: Find non-existent category
	t.Run("Find non-existent category returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		found, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestCategoryRepository_FindByUserID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()
	circleID := primitive.NewObjectID()

	// Create categories for user1
	categories1 := []*models.Category{
		{
			Name:     "Movies",
			Icon:     "film",
			CircleID: circleID,
			OwnerID:  user1ID,
		},
		{
			Name:     "Books",
			Icon:     "book",
			CircleID: circleID,
			OwnerID:  user1ID,
		},
		{
			Name:     "Music",
			Icon:     "music",
			CircleID: circleID,
			OwnerID:  user1ID,
		},
	}

	for _, category := range categories1 {
		require.NoError(t, repo.Create(ctx, category))
	}

	// Create category for user2
	category2 := &models.Category{
		Name:     "User 2 Category",
		CircleID: circleID,
		OwnerID:  user2ID,
	}
	require.NoError(t, repo.Create(ctx, category2))

	// Test case: Find all categories for user1
	t.Run("Find all categories for user", func(t *testing.T) {
		found, err := repo.FindByUserID(ctx, user1ID)
		require.NoError(t, err)
		assert.Len(t, found, 3)

		// Verify all categories belong to user1
		for _, category := range found {
			assert.Equal(t, user1ID, category.OwnerID)
		}
	})

	// Test case: Find categories for user with no categories
	t.Run("Find categories for user with no categories", func(t *testing.T) {
		emptyUserID := primitive.NewObjectID()
		found, err := repo.FindByUserID(ctx, emptyUserID)
		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestCategoryRepository_FindByCircleID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	ownerID := primitive.NewObjectID()
	circle1ID := primitive.NewObjectID()
	circle2ID := primitive.NewObjectID()

	// Create categories for circle1
	categories1 := []*models.Category{
		{
			Name:     "Movies",
			CircleID: circle1ID,
			OwnerID:  ownerID,
		},
		{
			Name:     "TV Shows",
			CircleID: circle1ID,
			OwnerID:  ownerID,
		},
	}

	for _, category := range categories1 {
		require.NoError(t, repo.Create(ctx, category))
	}

	// Create category for circle2
	category2 := &models.Category{
		Name:     "Restaurants",
		CircleID: circle2ID,
		OwnerID:  ownerID,
	}
	require.NoError(t, repo.Create(ctx, category2))

	// Test case: Find all categories in circle1
	t.Run("Find all categories in circle", func(t *testing.T) {
		found, err := repo.FindByCircleID(ctx, circle1ID)
		require.NoError(t, err)
		assert.Len(t, found, 2)

		// Verify all categories belong to circle1
		for _, category := range found {
			assert.Equal(t, circle1ID, category.CircleID)
		}
	})

	// Test case: Find categories in empty circle
	t.Run("Find categories in empty circle", func(t *testing.T) {
		emptyCircleID := primitive.NewObjectID()
		found, err := repo.FindByCircleID(ctx, emptyCircleID)
		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestCategoryRepository_Update(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Create test category
	ownerID := primitive.NewObjectID()
	circleID := primitive.NewObjectID()
	category := &models.Category{
		Name:      "Original Name",
		Icon:      "star",
		CircleID:  circleID,
		OwnerID:   ownerID,
		ItemCount: 5,
	}
	require.NoError(t, repo.Create(ctx, category))

	// Test case: Update existing category
	t.Run("Update existing category successfully", func(t *testing.T) {
		// Save original updated time
		originalUpdatedAt := category.UpdatedAt

		// Wait a bit to ensure timestamp changes
		time.Sleep(10 * time.Millisecond)

		// Update category
		category.Name = "Updated Name"
		category.Icon = "heart"
		category.ItemCount = 10

		err := repo.Update(ctx, category)
		require.NoError(t, err)
		assert.True(t, category.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		found, err := repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.Equal(t, "heart", found.Icon)
		assert.Equal(t, 10, found.ItemCount)
	})

	// Test case: Update non-existent category
	t.Run("Update non-existent category returns error", func(t *testing.T) {
		nonExistentCategory := &models.Category{
			ID:       primitive.NewObjectID(),
			Name:     "Ghost Category",
			CircleID: circleID,
			OwnerID:  ownerID,
		}

		err := repo.Update(ctx, nonExistentCategory)
		assert.Error(t, err)
	})
}

func TestCategoryRepository_Delete(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Create test category
	ownerID := primitive.NewObjectID()
	circleID := primitive.NewObjectID()
	category := &models.Category{
		Name:     "Delete Me",
		CircleID: circleID,
		OwnerID:  ownerID,
	}
	require.NoError(t, repo.Create(ctx, category))

	// Test case: Delete existing category
	t.Run("Delete existing category successfully", func(t *testing.T) {
		err := repo.Delete(ctx, category.ID)
		require.NoError(t, err)

		// Verify deletion
		found, err := repo.FindByID(ctx, category.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	// Test case: Delete non-existent category
	t.Run("Delete non-existent category returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestCategoryRepository_IncrementItemCount(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Create test category
	ownerID := primitive.NewObjectID()
	circleID := primitive.NewObjectID()
	category := &models.Category{
		Name:      "Counter Test",
		CircleID:  circleID,
		OwnerID:   ownerID,
		ItemCount: 0,
	}
	require.NoError(t, repo.Create(ctx, category))

	// Test case: Increment item count
	t.Run("Increment item count successfully", func(t *testing.T) {
		err := repo.IncrementItemCount(ctx, category.ID)
		require.NoError(t, err)

		// Verify increment
		found, err := repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, found.ItemCount)

		// Increment again
		err = repo.IncrementItemCount(ctx, category.ID)
		require.NoError(t, err)

		found, err = repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, found.ItemCount)
	})

	// Test case: Increment non-existent category
	t.Run("Increment non-existent category returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.IncrementItemCount(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestCategoryRepository_DecrementItemCount(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Create test category with initial count
	ownerID := primitive.NewObjectID()
	circleID := primitive.NewObjectID()
	category := &models.Category{
		Name:      "Counter Test",
		CircleID:  circleID,
		OwnerID:   ownerID,
		ItemCount: 5,
	}
	require.NoError(t, repo.Create(ctx, category))

	// Test case: Decrement item count
	t.Run("Decrement item count successfully", func(t *testing.T) {
		err := repo.DecrementItemCount(ctx, category.ID)
		require.NoError(t, err)

		// Verify decrement
		found, err := repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, 4, found.ItemCount)

		// Decrement again
		err = repo.DecrementItemCount(ctx, category.ID)
		require.NoError(t, err)

		found, err = repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, found.ItemCount)
	})

	// Test case: Decrement non-existent category
	t.Run("Decrement non-existent category returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.DecrementItemCount(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestCategoryRepository_Integration(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCategoryRepository(mongo.DB)
	ctx := context.Background()

	// Integration test: Full lifecycle
	t.Run("Full category lifecycle", func(t *testing.T) {
		ownerID := primitive.NewObjectID()
		circleID := primitive.NewObjectID()

		// 1. Create category
		category := &models.Category{
			Name:      "Movies",
			Icon:      "film",
			CircleID:  circleID,
			OwnerID:   ownerID,
			ItemCount: 0,
		}
		err := repo.Create(ctx, category)
		require.NoError(t, err)

		// 2. Find by ID
		found, err := repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, category.Name, found.Name)

		// 3. Find by UserID
		userCategories, err := repo.FindByUserID(ctx, ownerID)
		require.NoError(t, err)
		assert.Len(t, userCategories, 1)

		// 4. Find by CircleID
		circleCategories, err := repo.FindByCircleID(ctx, circleID)
		require.NoError(t, err)
		assert.Len(t, circleCategories, 1)

		// 5. Increment item count
		err = repo.IncrementItemCount(ctx, category.ID)
		require.NoError(t, err)

		incremented, err := repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, incremented.ItemCount)

		// 6. Increment again
		err = repo.IncrementItemCount(ctx, category.ID)
		require.NoError(t, err)

		incremented, err = repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, incremented.ItemCount)

		// 7. Decrement item count
		err = repo.DecrementItemCount(ctx, category.ID)
		require.NoError(t, err)

		decremented, err := repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, decremented.ItemCount)

		// 8. Update category
		category.Name = "Movies & TV Shows"
		category.Icon = "tv"
		err = repo.Update(ctx, category)
		require.NoError(t, err)

		// 9. Verify update
		updated, err := repo.FindByID(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, "Movies & TV Shows", updated.Name)
		assert.Equal(t, "tv", updated.Icon)

		// 10. Delete
		err = repo.Delete(ctx, category.ID)
		require.NoError(t, err)

		// 11. Verify deletion
		deleted, err := repo.FindByID(ctx, category.ID)
		assert.Error(t, err)
		assert.Nil(t, deleted)
	})
}

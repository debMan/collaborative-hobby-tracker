package repository

import (
	"context"
	"testing"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateIndexes(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Test case: Create all indexes successfully
	t.Run("Create all indexes successfully", func(t *testing.T) {
		err := CreateIndexes(ctx, mongo.DB)
		require.NoError(t, err)

		// Verify user indexes
		userIndexes := getIndexNames(t, ctx, mongo, models.User{}.CollectionName())
		assert.Contains(t, userIndexes, "email_unique")

		// Verify hobby item indexes
		itemIndexes := getIndexNames(t, ctx, mongo, models.HobbyItem{}.CollectionName())
		assert.Contains(t, itemIndexes, "owner_id_idx")
		assert.Contains(t, itemIndexes, "category_id_idx")
		assert.Contains(t, itemIndexes, "owner_category_idx")
		assert.Contains(t, itemIndexes, "is_completed_idx")
		assert.Contains(t, itemIndexes, "owner_completed_idx")

		// Verify category indexes
		categoryIndexes := getIndexNames(t, ctx, mongo, models.Category{}.CollectionName())
		assert.Contains(t, categoryIndexes, "owner_id_idx")
		assert.Contains(t, categoryIndexes, "circle_id_idx")
		assert.Contains(t, categoryIndexes, "circle_owner_idx")

		// Verify circle indexes
		circleIndexes := getIndexNames(t, ctx, mongo, models.Circle{}.CollectionName())
		assert.Contains(t, circleIndexes, "owner_id_idx")
		assert.Contains(t, circleIndexes, "members_user_id_idx")

		// Verify tag indexes
		tagIndexes := getIndexNames(t, ctx, mongo, models.Tag{}.CollectionName())
		assert.Contains(t, tagIndexes, "user_id_idx")
		assert.Contains(t, tagIndexes, "user_name_unique")
		assert.Contains(t, tagIndexes, "usage_count_idx")
	})

	// Test case: Idempotent - calling CreateIndexes multiple times should not error
	t.Run("CreateIndexes is idempotent", func(t *testing.T) {
		err := CreateIndexes(ctx, mongo.DB)
		require.NoError(t, err)

		// Call again - should not error
		err = CreateIndexes(ctx, mongo.DB)
		require.NoError(t, err)
	})
}

func TestUserIndexes(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create indexes
	err := CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	// Test case: Email unique index works
	t.Run("Email unique index prevents duplicates", func(t *testing.T) {
		collection := mongo.DB.Collection(models.User{}.CollectionName())

		// Insert first user
		_, err := collection.InsertOne(ctx, bson.M{
			"email": "test@example.com",
			"name":  "User 1",
		})
		require.NoError(t, err)

		// Try to insert duplicate email - should fail
		_, err = collection.InsertOne(ctx, bson.M{
			"email": "test@example.com",
			"name":  "User 2",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
	})
}

func TestTagIndexes(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	ctx := context.Background()

	// Create indexes
	err := CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	// Test case: User-Name unique index works
	t.Run("User-Name unique index prevents duplicates per user", func(t *testing.T) {
		collection := mongo.DB.Collection(models.Tag{}.CollectionName())

		userID1 := "user1"
		userID2 := "user2"

		// Insert first tag for user1
		_, err := collection.InsertOne(ctx, bson.M{
			"user_id": userID1,
			"name":    "Action",
		})
		require.NoError(t, err)

		// Try to insert duplicate tag name for user1 - should fail
		_, err = collection.InsertOne(ctx, bson.M{
			"user_id": userID1,
			"name":    "Action",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")

		// Insert same tag name for user2 - should succeed
		_, err = collection.InsertOne(ctx, bson.M{
			"user_id": userID2,
			"name":    "Action",
		})
		require.NoError(t, err)
	})
}

// Helper function to get index names for a collection
func getIndexNames(t *testing.T, ctx context.Context, mongo *testutil.MongoDBContainer, collectionName string) []string {
	t.Helper()

	collection := mongo.DB.Collection(collectionName)
	cursor, err := collection.Indexes().List(ctx)
	require.NoError(t, err)
	defer cursor.Close(ctx)

	var indexes []bson.M
	err = cursor.All(ctx, &indexes)
	require.NoError(t, err)

	var names []string
	for _, index := range indexes {
		if name, ok := index["name"].(string); ok {
			names = append(names, name)
		}
	}

	return names
}

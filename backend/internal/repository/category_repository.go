package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCategoryRepository struct {
	db *database.MongoDB
}

// NewCategoryRepository creates a new MongoDB category repository
func NewCategoryRepository(db *database.MongoDB) CategoryRepository {
	return &mongoCategoryRepository{db: db}
}

// Create creates a new category
func (r *mongoCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	// Generate ID if not provided
	if category.ID.IsZero() {
		category.ID = primitive.NewObjectID()
	}

	// Set timestamps
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	// Insert into database
	collection := r.db.Collection(category.CollectionName())
	_, err := collection.InsertOne(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// FindByID finds a category by ID
func (r *mongoCategoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Category, error) {
	collection := r.db.Collection(models.Category{}.CollectionName())

	var category models.Category
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	return &category, nil
}

// FindByUserID finds all categories owned by a user
func (r *mongoCategoryRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Category, error) {
	collection := r.db.Collection(models.Category{}.CollectionName())

	cursor, err := collection.Find(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find categories by user: %w", err)
	}
	defer cursor.Close(ctx)

	var categories []*models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("failed to decode categories: %w", err)
	}

	// Return empty slice instead of nil if no categories found
	if categories == nil {
		categories = []*models.Category{}
	}

	return categories, nil
}

// FindByCircleID finds all categories in a circle
func (r *mongoCategoryRepository) FindByCircleID(ctx context.Context, circleID primitive.ObjectID) ([]*models.Category, error) {
	collection := r.db.Collection(models.Category{}.CollectionName())

	cursor, err := collection.Find(ctx, bson.M{"circle_id": circleID})
	if err != nil {
		return nil, fmt.Errorf("failed to find categories by circle: %w", err)
	}
	defer cursor.Close(ctx)

	var categories []*models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("failed to decode categories: %w", err)
	}

	// Return empty slice instead of nil if no categories found
	if categories == nil {
		categories = []*models.Category{}
	}

	return categories, nil
}

// Update updates a category
func (r *mongoCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	// Update timestamp
	category.UpdatedAt = time.Now()

	collection := r.db.Collection(category.CollectionName())
	result, err := collection.ReplaceOne(ctx, bson.M{"_id": category.ID}, category)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// Delete deletes a category by ID
func (r *mongoCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(models.Category{}.CollectionName())

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// IncrementItemCount increments the item count for a category
func (r *mongoCategoryRepository) IncrementItemCount(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(models.Category{}.CollectionName())

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"item_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to increment item count: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// DecrementItemCount decrements the item count for a category
func (r *mongoCategoryRepository) DecrementItemCount(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(models.Category{}.CollectionName())

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"item_count": -1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to decrement item count: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

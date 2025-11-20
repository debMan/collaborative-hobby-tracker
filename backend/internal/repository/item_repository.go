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

type mongoHobbyItemRepository struct {
	db *database.MongoDB
}

// NewHobbyItemRepository creates a new MongoDB hobby item repository
func NewHobbyItemRepository(db *database.MongoDB) HobbyItemRepository {
	return &mongoHobbyItemRepository{db: db}
}

// Create creates a new hobby item
func (r *mongoHobbyItemRepository) Create(ctx context.Context, item *models.HobbyItem) error {
	// Generate ID if not provided
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}

	// Set timestamps
	now := time.Now()
	item.AddedAt = now
	item.CreatedAt = now
	item.UpdatedAt = now

	// Insert into database
	collection := r.db.Collection(item.CollectionName())
	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to create hobby item: %w", err)
	}

	return nil
}

// FindByID finds a hobby item by ID
func (r *mongoHobbyItemRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.HobbyItem, error) {
	collection := r.db.Collection(models.HobbyItem{}.CollectionName())

	var item models.HobbyItem
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("hobby item not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find hobby item: %w", err)
	}

	return &item, nil
}

// FindByUserID finds all hobby items for a user
func (r *mongoHobbyItemRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.HobbyItem, error) {
	collection := r.db.Collection(models.HobbyItem{}.CollectionName())

	cursor, err := collection.Find(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find hobby items by user: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*models.HobbyItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode hobby items: %w", err)
	}

	// Return empty slice instead of nil if no items found
	if items == nil {
		items = []*models.HobbyItem{}
	}

	return items, nil
}

// FindByCategoryID finds all hobby items in a category
func (r *mongoHobbyItemRepository) FindByCategoryID(ctx context.Context, categoryID primitive.ObjectID) ([]*models.HobbyItem, error) {
	collection := r.db.Collection(models.HobbyItem{}.CollectionName())

	cursor, err := collection.Find(ctx, bson.M{"category_id": categoryID})
	if err != nil {
		return nil, fmt.Errorf("failed to find hobby items by category: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*models.HobbyItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode hobby items: %w", err)
	}

	// Return empty slice instead of nil if no items found
	if items == nil {
		items = []*models.HobbyItem{}
	}

	return items, nil
}

// Update updates a hobby item
func (r *mongoHobbyItemRepository) Update(ctx context.Context, item *models.HobbyItem) error {
	// Update timestamp
	item.UpdatedAt = time.Now()

	collection := r.db.Collection(item.CollectionName())
	result, err := collection.ReplaceOne(ctx, bson.M{"_id": item.ID}, item)
	if err != nil {
		return fmt.Errorf("failed to update hobby item: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("hobby item not found")
	}

	return nil
}

// Delete deletes a hobby item by ID
func (r *mongoHobbyItemRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(models.HobbyItem{}.CollectionName())

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete hobby item: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("hobby item not found")
	}

	return nil
}

// ToggleComplete toggles the completion status of an item
func (r *mongoHobbyItemRepository) ToggleComplete(ctx context.Context, id primitive.ObjectID) error {
	// First, find the item to get its current state
	item, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Toggle the completion status
	if item.IsCompleted {
		item.MarkIncomplete()
	} else {
		item.MarkCompleted()
	}

	// Update the item
	collection := r.db.Collection(item.CollectionName())
	result, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, item)
	if err != nil {
		return fmt.Errorf("failed to toggle hobby item completion: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("hobby item not found")
	}

	return nil
}

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

type mongoTagRepository struct {
	db *database.MongoDB
}

// NewTagRepository creates a new MongoDB tag repository
func NewTagRepository(db *database.MongoDB) TagRepository {
	return &mongoTagRepository{db: db}
}

// Create creates a new tag
func (r *mongoTagRepository) Create(ctx context.Context, tag *models.Tag) error {
	// Generate ID if not provided
	if tag.ID.IsZero() {
		tag.ID = primitive.NewObjectID()
	}

	// Set timestamps
	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now

	// Insert into database
	collection := r.db.Collection(tag.CollectionName())
	_, err := collection.InsertOne(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// FindByID finds a tag by ID
func (r *mongoTagRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Tag, error) {
	collection := r.db.Collection(models.Tag{}.CollectionName())

	var tag models.Tag
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&tag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("tag not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}

	return &tag, nil
}

// FindByUserID finds all tags for a user
func (r *mongoTagRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Tag, error) {
	collection := r.db.Collection(models.Tag{}.CollectionName())

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find tags by user: %w", err)
	}
	defer cursor.Close(ctx)

	var tags []*models.Tag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, fmt.Errorf("failed to decode tags: %w", err)
	}

	// Return empty slice instead of nil if no tags found
	if tags == nil {
		tags = []*models.Tag{}
	}

	return tags, nil
}

// FindByName finds a tag by name for a specific user
func (r *mongoTagRepository) FindByName(ctx context.Context, userID primitive.ObjectID, name string) (*models.Tag, error) {
	collection := r.db.Collection(models.Tag{}.CollectionName())

	var tag models.Tag
	err := collection.FindOne(ctx, bson.M{
		"user_id": userID,
		"name":    name,
	}).Decode(&tag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("tag not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find tag by name: %w", err)
	}

	return &tag, nil
}

// Update updates a tag
func (r *mongoTagRepository) Update(ctx context.Context, tag *models.Tag) error {
	// Update timestamp
	tag.UpdatedAt = time.Now()

	collection := r.db.Collection(tag.CollectionName())
	result, err := collection.ReplaceOne(ctx, bson.M{"_id": tag.ID}, tag)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}

// Delete deletes a tag by ID
func (r *mongoTagRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(models.Tag{}.CollectionName())

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}

// IncrementUsage increments the usage count for a tag
func (r *mongoTagRepository) IncrementUsage(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(models.Tag{}.CollectionName())

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"usage_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to increment tag usage: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}

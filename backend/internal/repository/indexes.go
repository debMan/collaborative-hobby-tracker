package repository

import (
	"context"
	"fmt"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes creates all necessary indexes for the database
func CreateIndexes(ctx context.Context, db *database.MongoDB) error {
	// Create indexes for each collection
	if err := createUserIndexes(ctx, db); err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	if err := createHobbyItemIndexes(ctx, db); err != nil {
		return fmt.Errorf("failed to create hobby item indexes: %w", err)
	}

	if err := createCategoryIndexes(ctx, db); err != nil {
		return fmt.Errorf("failed to create category indexes: %w", err)
	}

	if err := createCircleIndexes(ctx, db); err != nil {
		return fmt.Errorf("failed to create circle indexes: %w", err)
	}

	if err := createTagIndexes(ctx, db); err != nil {
		return fmt.Errorf("failed to create tag indexes: %w", err)
	}

	return nil
}

// createUserIndexes creates indexes for the users collection
func createUserIndexes(ctx context.Context, db *database.MongoDB) error {
	collection := db.Collection(models.User{}.CollectionName())

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createHobbyItemIndexes creates indexes for the hobby_items collection
func createHobbyItemIndexes(ctx context.Context, db *database.MongoDB) error {
	collection := db.Collection(models.HobbyItem{}.CollectionName())

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "owner_id", Value: 1}},
			Options: options.Index().SetName("owner_id_idx"),
		},
		{
			Keys:    bson.D{{Key: "category_id", Value: 1}},
			Options: options.Index().SetName("category_id_idx"),
		},
		{
			Keys: bson.D{
				{Key: "owner_id", Value: 1},
				{Key: "category_id", Value: 1},
			},
			Options: options.Index().SetName("owner_category_idx"),
		},
		{
			Keys:    bson.D{{Key: "is_completed", Value: 1}},
			Options: options.Index().SetName("is_completed_idx"),
		},
		{
			Keys: bson.D{
				{Key: "owner_id", Value: 1},
				{Key: "is_completed", Value: 1},
			},
			Options: options.Index().SetName("owner_completed_idx"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createCategoryIndexes creates indexes for the categories collection
func createCategoryIndexes(ctx context.Context, db *database.MongoDB) error {
	collection := db.Collection(models.Category{}.CollectionName())

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "owner_id", Value: 1}},
			Options: options.Index().SetName("owner_id_idx"),
		},
		{
			Keys:    bson.D{{Key: "circle_id", Value: 1}},
			Options: options.Index().SetName("circle_id_idx"),
		},
		{
			Keys: bson.D{
				{Key: "circle_id", Value: 1},
				{Key: "owner_id", Value: 1},
			},
			Options: options.Index().SetName("circle_owner_idx"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createCircleIndexes creates indexes for the circles collection
func createCircleIndexes(ctx context.Context, db *database.MongoDB) error {
	collection := db.Collection(models.Circle{}.CollectionName())

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "owner_id", Value: 1}},
			Options: options.Index().SetName("owner_id_idx"),
		},
		{
			Keys:    bson.D{{Key: "members.user_id", Value: 1}},
			Options: options.Index().SetName("members_user_id_idx"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createTagIndexes creates indexes for the tags collection
func createTagIndexes(ctx context.Context, db *database.MongoDB) error {
	collection := db.Collection(models.Tag{}.CollectionName())

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("user_id_idx"),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "name", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("user_name_unique"),
		},
		{
			Keys:    bson.D{{Key: "usage_count", Value: -1}},
			Options: options.Index().SetName("usage_count_idx"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

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

// mongoUserRepository implements UserRepository using MongoDB
type mongoUserRepository struct {
	db *database.MongoDB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *database.MongoDB) UserRepository {
	return &mongoUserRepository{
		db: db,
	}
}

// Create creates a new user
func (r *mongoUserRepository) Create(ctx context.Context, user *models.User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := r.db.Collection(user.CollectionName())
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID finds a user by ID
func (r *mongoUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	collection := r.db.Collection((&models.User{}).CollectionName())

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := r.db.Collection((&models.User{}).CollectionName())

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// Update updates a user
func (r *mongoUserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	collection := r.db.Collection(user.CollectionName())
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user by ID
func (r *mongoUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection((&models.User{}).CollectionName())

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

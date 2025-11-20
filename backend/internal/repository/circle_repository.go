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

type mongoCircleRepository struct {
	db *database.MongoDB
}

// NewCircleRepository creates a new MongoDB circle repository
func NewCircleRepository(db *database.MongoDB) CircleRepository {
	return &mongoCircleRepository{db: db}
}

// Create creates a new circle
func (r *mongoCircleRepository) Create(ctx context.Context, circle *models.Circle) error {
	// Generate ID if not provided
	if circle.ID.IsZero() {
		circle.ID = primitive.NewObjectID()
	}

	// Set timestamps
	now := time.Now()
	circle.CreatedAt = now
	circle.UpdatedAt = now

	// Initialize members array if nil
	if circle.Members == nil {
		circle.Members = []models.CircleMember{}
	}

	// Insert into database
	collection := r.db.Collection(circle.CollectionName())
	_, err := collection.InsertOne(ctx, circle)
	if err != nil {
		return fmt.Errorf("failed to create circle: %w", err)
	}

	return nil
}

// FindByID finds a circle by ID
func (r *mongoCircleRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Circle, error) {
	collection := r.db.Collection(models.Circle{}.CollectionName())

	var circle models.Circle
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&circle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("circle not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find circle: %w", err)
	}

	return &circle, nil
}

// FindByUserID finds all circles where the user is owner or member
func (r *mongoCircleRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Circle, error) {
	collection := r.db.Collection(models.Circle{}.CollectionName())

	// Find circles where user is owner OR a member
	filter := bson.M{
		"$or": []bson.M{
			{"owner_id": userID},
			{"members.user_id": userID},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find circles by user: %w", err)
	}
	defer cursor.Close(ctx)

	var circles []*models.Circle
	if err := cursor.All(ctx, &circles); err != nil {
		return nil, fmt.Errorf("failed to decode circles: %w", err)
	}

	// Return empty slice instead of nil if no circles found
	if circles == nil {
		circles = []*models.Circle{}
	}

	return circles, nil
}

// Update updates a circle
func (r *mongoCircleRepository) Update(ctx context.Context, circle *models.Circle) error {
	// Update timestamp
	circle.UpdatedAt = time.Now()

	collection := r.db.Collection(circle.CollectionName())
	result, err := collection.ReplaceOne(ctx, bson.M{"_id": circle.ID}, circle)
	if err != nil {
		return fmt.Errorf("failed to update circle: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("circle not found")
	}

	return nil
}

// Delete deletes a circle by ID
func (r *mongoCircleRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(models.Circle{}.CollectionName())

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete circle: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("circle not found")
	}

	return nil
}

// AddMember adds a member to a circle
func (r *mongoCircleRepository) AddMember(ctx context.Context, circleID primitive.ObjectID, member models.CircleMember) error {
	collection := r.db.Collection(models.Circle{}.CollectionName())

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": circleID},
		bson.M{
			"$push": bson.M{"members": member},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add member to circle: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("circle not found")
	}

	return nil
}

// RemoveMember removes a member from a circle
func (r *mongoCircleRepository) RemoveMember(ctx context.Context, circleID, userID primitive.ObjectID) error {
	collection := r.db.Collection(models.Circle{}.CollectionName())

	// First, check if the circle exists and has the member
	circle, err := r.FindByID(ctx, circleID)
	if err != nil {
		return err
	}

	// Check if member exists
	memberExists := false
	for _, member := range circle.Members {
		if member.UserID == userID {
			memberExists = true
			break
		}
	}

	if !memberExists {
		return fmt.Errorf("member not found in circle")
	}

	// Remove the member
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": circleID},
		bson.M{
			"$pull": bson.M{"members": bson.M{"user_id": userID}},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to remove member from circle: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("circle not found")
	}

	return nil
}

// UpdateMemberAccess updates a member's access level
func (r *mongoCircleRepository) UpdateMemberAccess(ctx context.Context, circleID, userID primitive.ObjectID, accessLevel models.AccessLevel) error {
	collection := r.db.Collection(models.Circle{}.CollectionName())

	// First, check if the circle exists and has the member
	circle, err := r.FindByID(ctx, circleID)
	if err != nil {
		return err
	}

	// Check if member exists
	memberExists := false
	for _, member := range circle.Members {
		if member.UserID == userID {
			memberExists = true
			break
		}
	}

	if !memberExists {
		return fmt.Errorf("member not found in circle")
	}

	// Update the member's access level using positional operator
	result, err := collection.UpdateOne(
		ctx,
		bson.M{
			"_id":              circleID,
			"members.user_id": userID,
		},
		bson.M{
			"$set": bson.M{
				"members.$.access_level": accessLevel,
				"updated_at":             time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update member access: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("circle or member not found")
	}

	return nil
}

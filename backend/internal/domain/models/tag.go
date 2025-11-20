package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tag represents a user-defined tag for categorizing items
// AI learns from user's tag usage patterns
type Tag struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Color      string             `bson:"color,omitempty" json:"color,omitempty"` // Hex color code
	UsageCount int                `bson:"usage_count" json:"usageCount"`          // How many times this tag has been used
	UserID     primitive.ObjectID `bson:"user_id" json:"userId"`                  // Tags are user-specific
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CollectionName returns the MongoDB collection name for Tag
func (Tag) CollectionName() string {
	return "tags"
}

// IncrementUsage increments the usage count of the tag
func (t *Tag) IncrementUsage() {
	t.UsageCount++
	t.UpdatedAt = time.Now()
}

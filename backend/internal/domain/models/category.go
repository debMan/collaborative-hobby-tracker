package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category represents a collection of hobby items
// Examples: Movies, Restaurants, Travel Destinations, Music, Activities
type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Icon      string             `bson:"icon,omitempty" json:"icon,omitempty"` // Icon name (e.g., "film", "utensils")
	CircleID  primitive.ObjectID `bson:"circle_id" json:"circleId"`
	OwnerID   primitive.ObjectID `bson:"owner_id" json:"ownerId"` // For quick "my categories" queries
	ItemCount int                `bson:"item_count" json:"itemCount"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CollectionName returns the MongoDB collection name for Category
func (Category) CollectionName() string {
	return "categories"
}

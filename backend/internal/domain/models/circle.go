package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccessLevel defines the access level for circle members
type AccessLevel string

const (
	AccessLevelPrivate AccessLevel = "private" // Owner only
	AccessLevelView    AccessLevel = "view"    // Can view items
	AccessLevelEdit    AccessLevel = "edit"    // Can add/edit items
	AccessLevelAdmin   AccessLevel = "admin"   // Can manage circle and members
)

// CircleMember represents a member of a circle with their access level
type CircleMember struct {
	UserID     primitive.ObjectID `bson:"user_id" json:"userId"`
	AccessLevel AccessLevel        `bson:"access_level" json:"accessLevel"`
	InvitedAt  time.Time          `bson:"invited_at" json:"invitedAt"`
	AcceptedAt *time.Time         `bson:"accepted_at,omitempty" json:"acceptedAt,omitempty"`
}

// Circle represents a group for sharing categories
// Examples: Partner, Friends, Family, Colleagues
type Circle struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	OwnerID   primitive.ObjectID `bson:"owner_id" json:"ownerId"`
	Members   []CircleMember     `bson:"members" json:"members"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CollectionName returns the MongoDB collection name for Circle
func (Circle) CollectionName() string {
	return "circles"
}

// HasMember checks if a user is a member of the circle
func (c *Circle) HasMember(userID primitive.ObjectID) bool {
	for _, member := range c.Members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}

// GetMemberAccessLevel returns the access level of a user in the circle
func (c *Circle) GetMemberAccessLevel(userID primitive.ObjectID) (AccessLevel, bool) {
	// Owner has admin access
	if c.OwnerID == userID {
		return AccessLevelAdmin, true
	}

	for _, member := range c.Members {
		if member.UserID == userID {
			return member.AccessLevel, true
		}
	}
	return "", false
}

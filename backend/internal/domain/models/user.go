package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email           string             `bson:"email" json:"email"`
	PasswordHash    string             `bson:"password_hash" json:"-"` // Never expose in JSON
	Name            string             `bson:"name" json:"name"`
	AvatarURL       string             `bson:"avatar_url,omitempty" json:"avatarUrl,omitempty"`
	OAuthProvider   string             `bson:"oauth_provider,omitempty" json:"oauthProvider,omitempty"`     // User ID from OAuth provider (e.g., "google-user-123")
	OAuthProviderName string           `bson:"oauth_provider_name,omitempty" json:"oauthProviderName,omitempty"` // Provider name: google, github, apple
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}

// TableName returns the MongoDB collection name for User
func (User) CollectionName() string {
	return "users"
}

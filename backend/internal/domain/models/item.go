package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DataSource represents the source of an imported item
type DataSource string

const (
	SourceManual    DataSource = "manual"
	SourceInstagram DataSource = "instagram"
	SourceYouTube   DataSource = "youtube"
	SourceTwitter   DataSource = "twitter"
	SourceTikTok    DataSource = "tiktok"
	SourceTelegram  DataSource = "telegram"
	SourceWeb       DataSource = "web"
	SourceWikipedia DataSource = "wikipedia"
)

// HobbyItem represents a single item to track
// Examples: a movie to watch, a restaurant to visit, a destination to travel
type HobbyItem struct {
	ID                 primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Title              string                 `bson:"title" json:"title"`
	Description        string                 `bson:"description,omitempty" json:"description,omitempty"`
	CategoryID         primitive.ObjectID     `bson:"category_id" json:"categoryId"`
	OwnerID            primitive.ObjectID     `bson:"owner_id" json:"ownerId"` // User who added this item
	IsCompleted        bool                   `bson:"is_completed" json:"isCompleted"`
	AddedAt            time.Time              `bson:"added_at" json:"addedAt"`
	CompletedAt        *time.Time             `bson:"completed_at,omitempty" json:"completedAt,omitempty"`
	DueDate            *time.Time             `bson:"due_date,omitempty" json:"dueDate,omitempty"`
	Source             DataSource             `bson:"source" json:"source"`
	SourceURL          string                 `bson:"source_url,omitempty" json:"sourceUrl,omitempty"`
	ImageURL           string                 `bson:"image_url,omitempty" json:"imageUrl,omitempty"`
	CategoryConfidence float64                `bson:"category_confidence,omitempty" json:"categoryConfidence,omitempty"` // AI confidence (0-1)
	Tags               []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Metadata           map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // Flexible metadata per category
	CreatedAt          time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt          time.Time              `bson:"updated_at" json:"updatedAt"`
}

// CollectionName returns the MongoDB collection name for HobbyItem
func (HobbyItem) CollectionName() string {
	return "hobby_items"
}

// MarkCompleted marks the item as completed with the current timestamp
func (h *HobbyItem) MarkCompleted() {
	h.IsCompleted = true
	now := time.Now()
	h.CompletedAt = &now
	h.UpdatedAt = now
}

// MarkIncomplete marks the item as incomplete
func (h *HobbyItem) MarkIncomplete() {
	h.IsCompleted = false
	h.CompletedAt = nil
	h.UpdatedAt = time.Now()
}

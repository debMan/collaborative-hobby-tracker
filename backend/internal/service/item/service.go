package item

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrUnauthorized is returned when user doesn't have permission
	ErrUnauthorized = errors.New("unauthorized: you don't have permission to access this item")
	// ErrTitleRequired is returned when title is empty
	ErrTitleRequired = errors.New("title is required")
	// ErrCategoryRequired is returned when category ID is empty
	ErrCategoryRequired = errors.New("category ID is required")
	// ErrUserIDRequired is returned when user ID is empty
	ErrUserIDRequired = errors.New("user ID is required")
)

// Service handles business logic for hobby items
type Service struct {
	itemRepo repository.HobbyItemRepository
}

// NewService creates a new item service
func NewService(itemRepo repository.HobbyItemRepository) *Service {
	return &Service{
		itemRepo: itemRepo,
	}
}

// CreateItemRequest represents a request to create an item
type CreateItemRequest struct {
	Title              string                 `json:"title"`
	Description        string                 `json:"description,omitempty"`
	CategoryID         primitive.ObjectID     `json:"categoryId"`
	Source             models.DataSource      `json:"source,omitempty"`
	SourceURL          string                 `json:"sourceUrl,omitempty"`
	ImageURL           string                 `json:"imageUrl,omitempty"`
	CategoryConfidence float64                `json:"categoryConfidence,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	DueDate            *time.Time             `json:"dueDate,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateItemRequest represents a request to update an item
type UpdateItemRequest struct {
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	DueDate     *time.Time             `json:"dueDate,omitempty"`
	ImageURL    string                 `json:"imageUrl,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CreateItem creates a new hobby item for a user
func (s *Service) CreateItem(ctx context.Context, userID primitive.ObjectID, req *CreateItemRequest) (*models.HobbyItem, error) {
	// Validate user ID
	if userID.IsZero() {
		return nil, ErrUserIDRequired
	}

	// Validate request
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if req.CategoryID.IsZero() {
		return nil, ErrCategoryRequired
	}

	// Set default source if not provided
	source := req.Source
	if source == "" {
		source = models.SourceManual
	}

	// Create item
	now := time.Now()
	item := &models.HobbyItem{
		Title:              req.Title,
		Description:        req.Description,
		CategoryID:         req.CategoryID,
		OwnerID:            userID,
		IsCompleted:        false,
		Source:             source,
		SourceURL:          req.SourceURL,
		ImageURL:           req.ImageURL,
		CategoryConfidence: req.CategoryConfidence,
		Tags:               req.Tags,
		DueDate:            req.DueDate,
		Metadata:           req.Metadata,
		AddedAt:            now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err := s.itemRepo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return item, nil
}

// GetItemByID retrieves an item by ID with authorization check
func (s *Service) GetItemByID(ctx context.Context, userID primitive.ObjectID, itemID primitive.ObjectID) (*models.HobbyItem, error) {
	item, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Authorization check: user must own the item
	if item.OwnerID != userID {
		return nil, ErrUnauthorized
	}

	return item, nil
}

// GetUserItems retrieves all items for a user
func (s *Service) GetUserItems(ctx context.Context, userID primitive.ObjectID) ([]*models.HobbyItem, error) {
	items, err := s.itemRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user items: %w", err)
	}

	return items, nil
}

// UpdateItem updates an existing item with authorization check
func (s *Service) UpdateItem(ctx context.Context, userID primitive.ObjectID, itemID primitive.ObjectID, req *UpdateItemRequest) (*models.HobbyItem, error) {
	// Validate title if provided
	if req.Title == "" {
		return nil, ErrTitleRequired
	}

	// Get existing item and check authorization
	item, err := s.GetItemByID(ctx, userID, itemID)
	if err != nil {
		return nil, err
	}

	// Update fields
	item.Title = req.Title
	item.Description = req.Description
	item.Tags = req.Tags
	item.DueDate = req.DueDate
	item.ImageURL = req.ImageURL
	item.Metadata = req.Metadata
	item.UpdatedAt = time.Now()

	err = s.itemRepo.Update(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return item, nil
}

// DeleteItem deletes an item with authorization check
func (s *Service) DeleteItem(ctx context.Context, userID primitive.ObjectID, itemID primitive.ObjectID) error {
	// Check authorization
	_, err := s.GetItemByID(ctx, userID, itemID)
	if err != nil {
		return err
	}

	err = s.itemRepo.Delete(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// ToggleItemCompletion toggles the completion status of an item with authorization check
func (s *Service) ToggleItemCompletion(ctx context.Context, userID primitive.ObjectID, itemID primitive.ObjectID) (*models.HobbyItem, error) {
	// Check authorization
	_, err := s.GetItemByID(ctx, userID, itemID)
	if err != nil {
		return nil, err
	}

	err = s.itemRepo.ToggleComplete(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to toggle item completion: %w", err)
	}

	// Return updated item
	updatedItem, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated item: %w", err)
	}

	return updatedItem, nil
}

package repository

import (
	"context"
	"errors"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrItemNotFound is returned when a hobby item is not found
	ErrItemNotFound = errors.New("hobby item not found")
	// ErrCategoryNotFound is returned when a category is not found
	ErrCategoryNotFound = errors.New("category not found")
	// ErrCircleNotFound is returned when a circle is not found
	ErrCircleNotFound = errors.New("circle not found")
	// ErrTagNotFound is returned when a tag is not found
	ErrTagNotFound = errors.New("tag not found")
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	// Update updates a user
	Update(ctx context.Context, user *models.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// HobbyItemRepository defines the interface for hobby item data access
type HobbyItemRepository interface {
	// Create creates a new hobby item
	Create(ctx context.Context, item *models.HobbyItem) error

	// FindByID finds a hobby item by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.HobbyItem, error)

	// FindByUserID finds all hobby items for a user
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.HobbyItem, error)

	// FindByCategoryID finds all hobby items in a category
	FindByCategoryID(ctx context.Context, categoryID primitive.ObjectID) ([]*models.HobbyItem, error)

	// Update updates a hobby item
	Update(ctx context.Context, item *models.HobbyItem) error

	// Delete deletes a hobby item by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// ToggleComplete toggles the completion status of an item
	ToggleComplete(ctx context.Context, id primitive.ObjectID) error
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	// Create creates a new category
	Create(ctx context.Context, category *models.Category) error

	// FindByID finds a category by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Category, error)

	// FindByUserID finds all categories owned by a user
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Category, error)

	// FindByCircleID finds all categories in a circle
	FindByCircleID(ctx context.Context, circleID primitive.ObjectID) ([]*models.Category, error)

	// Update updates a category
	Update(ctx context.Context, category *models.Category) error

	// Delete deletes a category by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// IncrementItemCount increments the item count for a category
	IncrementItemCount(ctx context.Context, id primitive.ObjectID) error

	// DecrementItemCount decrements the item count for a category
	DecrementItemCount(ctx context.Context, id primitive.ObjectID) error
}

// CircleRepository defines the interface for circle data access
type CircleRepository interface {
	// Create creates a new circle
	Create(ctx context.Context, circle *models.Circle) error

	// FindByID finds a circle by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Circle, error)

	// FindByUserID finds all circles where the user is owner or member
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Circle, error)

	// Update updates a circle
	Update(ctx context.Context, circle *models.Circle) error

	// Delete deletes a circle by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// AddMember adds a member to a circle
	AddMember(ctx context.Context, circleID primitive.ObjectID, member models.CircleMember) error

	// RemoveMember removes a member from a circle
	RemoveMember(ctx context.Context, circleID, userID primitive.ObjectID) error

	// UpdateMemberAccess updates a member's access level
	UpdateMemberAccess(ctx context.Context, circleID, userID primitive.ObjectID, accessLevel models.AccessLevel) error
}

// TagRepository defines the interface for tag data access
type TagRepository interface {
	// Create creates a new tag
	Create(ctx context.Context, tag *models.Tag) error

	// FindByID finds a tag by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Tag, error)

	// FindByUserID finds all tags for a user
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Tag, error)

	// FindByName finds a tag by name for a specific user
	FindByName(ctx context.Context, userID primitive.ObjectID, name string) (*models.Tag, error)

	// Update updates a tag
	Update(ctx context.Context, tag *models.Tag) error

	// Delete deletes a tag by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// IncrementUsage increments the usage count for a tag
	IncrementUsage(ctx context.Context, id primitive.ObjectID) error
}

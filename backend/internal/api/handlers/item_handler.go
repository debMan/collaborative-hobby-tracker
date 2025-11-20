package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/service/item"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ItemHandler handles HTTP requests for hobby items
type ItemHandler struct {
	itemService *item.Service
}

// NewItemHandler creates a new item handler
func NewItemHandler(itemService *item.Service) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// CreateItemRequest represents the HTTP request body for creating an item
type CreateItemRequest struct {
	Title              string                 `json:"title" binding:"required"`
	Description        string                 `json:"description"`
	CategoryID         string                 `json:"categoryId" binding:"required"`
	Source             models.DataSource      `json:"source"`
	SourceURL          string                 `json:"sourceUrl"`
	ImageURL           string                 `json:"imageUrl"`
	CategoryConfidence float64                `json:"categoryConfidence"`
	Tags               []string               `json:"tags"`
	DueDate            *string                `json:"dueDate"` // ISO 8601 date string
	Metadata           map[string]interface{} `json:"metadata"`
}

// UpdateItemRequest represents the HTTP request body for updating an item
type UpdateItemRequest struct {
	Title       string                 `json:"title" binding:"required"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	DueDate     *string                `json:"dueDate"` // ISO 8601 date string
	ImageURL    string                 `json:"imageUrl"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// CreateItem handles POST /items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse request body
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse category ID
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	// Create item service request
	serviceReq := &item.CreateItemRequest{
		Title:              req.Title,
		Description:        req.Description,
		CategoryID:         categoryID,
		Source:             req.Source,
		SourceURL:          req.SourceURL,
		ImageURL:           req.ImageURL,
		CategoryConfidence: req.CategoryConfidence,
		Tags:               req.Tags,
		Metadata:           req.Metadata,
	}

	// TODO: Parse DueDate if provided

	// Create item
	createdItem, err := h.itemService.CreateItem(c.Request.Context(), userID, serviceReq)
	if err != nil {
		if errors.Is(err, item.ErrTitleRequired) || errors.Is(err, item.ErrCategoryRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, createdItem)
}

// GetUserItems handles GET /items
func (h *ItemHandler) GetUserItems(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get items
	items, err := h.itemService.GetUserItems(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItemByID handles GET /items/:id
func (h *ItemHandler) GetItemByID(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse item ID from URL
	itemIDStr := c.Param("id")
	itemID, err := primitive.ObjectIDFromHex(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	// Get item
	foundItem, err := h.itemService.GetItemByID(c.Request.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, item.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get item"})
		return
	}

	c.JSON(http.StatusOK, foundItem)
}

// UpdateItem handles PUT /items/:id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse item ID from URL
	itemIDStr := c.Param("id")
	itemID, err := primitive.ObjectIDFromHex(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	// Parse request body
	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create item service request
	serviceReq := &item.UpdateItemRequest{
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		ImageURL:    req.ImageURL,
		Metadata:    req.Metadata,
	}

	// TODO: Parse DueDate if provided

	// Update item
	updatedItem, err := h.itemService.UpdateItem(c.Request.Context(), userID, itemID, serviceReq)
	if err != nil {
		if errors.Is(err, item.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, item.ErrTitleRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item"})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

// DeleteItem handles DELETE /items/:id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse item ID from URL
	itemIDStr := c.Param("id")
	itemID, err := primitive.ObjectIDFromHex(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	// Delete item
	err = h.itemService.DeleteItem(c.Request.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, item.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ToggleItemCompletion handles PATCH /items/:id/toggle
func (h *ItemHandler) ToggleItemCompletion(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse item ID from URL
	itemIDStr := c.Param("id")
	itemID, err := primitive.ObjectIDFromHex(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	// Toggle completion
	updatedItem, err := h.itemService.ToggleItemCompletion(c.Request.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, item.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to toggle item completion"})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

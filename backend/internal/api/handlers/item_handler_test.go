package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/middleware"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository/testutil"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/service/item"
	pkgauth "github.com/debMan/collaborative-hobby-tracker/backend/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	testJWTSecret     = "test-secret-key-for-jwt-tokens"
	testTokenDuration = 24 * time.Hour
)

// setupTestHandler creates a test handler with all dependencies
func setupTestHandler(t *testing.T) (*ItemHandler, *testutil.MongoDBContainer, *gin.Engine, primitive.ObjectID, primitive.ObjectID, string) {
	mongo := testutil.SetupMongoDB(t)
	ctx := context.Background()

	// Create indexes
	err := repository.CreateIndexes(ctx, mongo.DB)
	require.NoError(t, err)

	// Create repositories
	itemRepo := repository.NewHobbyItemRepository(mongo.DB)

	// Create services
	itemService := item.NewService(itemRepo)

	// Create handler
	handler := NewItemHandler(itemService)

	// Create test user and get JWT token
	userID := primitive.NewObjectID()
	userEmail := "test@example.com"
	token, err := pkgauth.GenerateToken(userID, userEmail, testJWTSecret, testTokenDuration)
	require.NoError(t, err)

	// Setup Gin router with auth middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(testJWTSecret))
	{
		protected.POST("/items", handler.CreateItem)
		protected.GET("/items", handler.GetUserItems)
		protected.GET("/items/:id", handler.GetItemByID)
		protected.PUT("/items/:id", handler.UpdateItem)
		protected.DELETE("/items/:id", handler.DeleteItem)
		protected.PATCH("/items/:id/toggle", handler.ToggleItemCompletion)
	}

	// Create test category ID
	categoryID := primitive.NewObjectID()

	return handler, mongo, router, userID, categoryID, token
}

func TestCreateItem(t *testing.T) {
	_, mongo, router, userID, categoryID, token := setupTestHandler(t)
	defer mongo.Teardown(t)

	t.Run("Create item successfully", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":       "Watch Inception",
			"description": "Amazing sci-fi movie",
			"categoryId":  categoryID.Hex(),
			"source":      "manual",
			"tags":        []string{"movie", "sci-fi"},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Watch Inception", response["title"])
		assert.Equal(t, "Amazing sci-fi movie", response["description"])
		assert.Equal(t, userID.Hex(), response["ownerId"])
		assert.NotEmpty(t, response["id"])
	})

	t.Run("Create item with minimal fields", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":      "Minimal item",
			"categoryId": categoryID.Hex(),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Reject request without auth token", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":      "Test item",
			"categoryId": categoryID.Hex(),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Reject empty title", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":      "",
			"categoryId": categoryID.Hex(),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestGetUserItems(t *testing.T) {
	handler, mongo, router, userID, categoryID, token := setupTestHandler(t)
	defer mongo.Teardown(t)

	// Create test items
	ctx := context.Background()
	items := []string{"Item 1", "Item 2", "Item 3"}
	for _, title := range items {
		req := &item.CreateItemRequest{
			Title:      title,
			CategoryID: categoryID,
		}
		_, err := handler.itemService.CreateItem(ctx, userID, req)
		require.NoError(t, err)
	}

	t.Run("Get all user items", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response, 3)
	})

	t.Run("Reject request without auth token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
		// No Authorization header
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestGetItemByID(t *testing.T) {
	handler, mongo, router, userID, categoryID, token := setupTestHandler(t)
	defer mongo.Teardown(t)

	// Create test item
	ctx := context.Background()
	createdItem, err := handler.itemService.CreateItem(ctx, userID, &item.CreateItemRequest{
		Title:      "Test Item",
		CategoryID: categoryID,
	})
	require.NoError(t, err)

	t.Run("Get item successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items/"+createdItem.ID.Hex(), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, createdItem.ID.Hex(), response["id"])
		assert.Equal(t, "Test Item", response["title"])
	})

	t.Run("Return 404 for non-existent item", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items/"+nonExistentID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("Return 400 for invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items/invalid-id", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject request without auth token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items/"+createdItem.ID.Hex(), nil)
		// No Authorization header
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestUpdateItem(t *testing.T) {
	handler, mongo, router, userID, categoryID, token := setupTestHandler(t)
	defer mongo.Teardown(t)

	// Create test item
	ctx := context.Background()
	createdItem, err := handler.itemService.CreateItem(ctx, userID, &item.CreateItemRequest{
		Title:      "Original Title",
		CategoryID: categoryID,
	})
	require.NoError(t, err)

	t.Run("Update item successfully", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":       "Updated Title",
			"description": "Updated description",
			"tags":        []string{"updated", "tags"},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/"+createdItem.ID.Hex(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Updated Title", response["title"])
		assert.Equal(t, "Updated description", response["description"])
	})

	t.Run("Reject empty title", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/"+createdItem.ID.Hex(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Return 404 for non-existent item", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		reqBody := map[string]interface{}{
			"title": "Updated",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/"+nonExistentID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("Reject request without auth token", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "Updated",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/"+createdItem.ID.Hex(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestDeleteItem(t *testing.T) {
	handler, mongo, router, userID, categoryID, token := setupTestHandler(t)
	defer mongo.Teardown(t)

	// Create test item
	ctx := context.Background()
	createdItem, err := handler.itemService.CreateItem(ctx, userID, &item.CreateItemRequest{
		Title:      "Item to delete",
		CategoryID: categoryID,
	})
	require.NoError(t, err)

	t.Run("Delete item successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/"+createdItem.ID.Hex(), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("Return 404 for non-existent item", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/"+nonExistentID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("Return 400 for invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/invalid-id", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject request without auth token", func(t *testing.T) {
		// Create another item for this test
		anotherItem, err := handler.itemService.CreateItem(ctx, userID, &item.CreateItemRequest{
			Title:      "Another item",
			CategoryID: categoryID,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/"+anotherItem.ID.Hex(), nil)
		// No Authorization header
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestToggleItemCompletion(t *testing.T) {
	handler, mongo, router, userID, categoryID, token := setupTestHandler(t)
	defer mongo.Teardown(t)

	// Create test item
	ctx := context.Background()
	createdItem, err := handler.itemService.CreateItem(ctx, userID, &item.CreateItemRequest{
		Title:      "Item to toggle",
		CategoryID: categoryID,
	})
	require.NoError(t, err)

	t.Run("Toggle item to completed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/items/"+createdItem.ID.Hex()+"/toggle", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, true, response["isCompleted"])
		assert.NotNil(t, response["completedAt"])
	})

	t.Run("Toggle completed item to incomplete", func(t *testing.T) {
		// Create a new item for this test
		newItem, err := handler.itemService.CreateItem(ctx, userID, &item.CreateItemRequest{
			Title:      "Another item to toggle",
			CategoryID: categoryID,
		})
		require.NoError(t, err)

		// First toggle to complete
		req1 := httptest.NewRequest(http.MethodPatch, "/api/v1/items/"+newItem.ID.Hex()+"/toggle", nil)
		req1.Header.Set("Authorization", "Bearer "+token)
		resp1 := httptest.NewRecorder()
		router.ServeHTTP(resp1, req1)
		require.Equal(t, http.StatusOK, resp1.Code)

		// Then toggle back to incomplete
		req2 := httptest.NewRequest(http.MethodPatch, "/api/v1/items/"+newItem.ID.Hex()+"/toggle", nil)
		req2.Header.Set("Authorization", "Bearer "+token)
		resp2 := httptest.NewRecorder()
		router.ServeHTTP(resp2, req2)

		assert.Equal(t, http.StatusOK, resp2.Code)

		var response map[string]interface{}
		err = json.Unmarshal(resp2.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, false, response["isCompleted"])
		assert.Nil(t, response["completedAt"])
	})

	t.Run("Return 404 for non-existent item", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/items/%s/toggle", nonExistentID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("Reject request without auth token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/items/"+createdItem.ID.Hex()+"/toggle", nil)
		// No Authorization header
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

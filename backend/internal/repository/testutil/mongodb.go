package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/database"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

// MongoDBContainer wraps a testcontainer MongoDB instance
type MongoDBContainer struct {
	Container *mongodb.MongoDBContainer
	DB        *database.MongoDB
	URI       string
}

// SetupMongoDB creates a test MongoDB container and returns a connected database
func SetupMongoDB(t *testing.T) *MongoDBContainer {
	t.Helper()

	ctx := context.Background()

	// Start MongoDB container
	mongoContainer, err := mongodb.Run(ctx,
		"mongo:7.0",
	)
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}

	// Get connection string
	uri, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get MongoDB connection string: %v", err)
	}

	// Connect to MongoDB
	db, err := database.NewMongoDB(database.MongoDBConfig{
		URI:      uri,
		Name:     "test_db",
		Timeout:  10,
		PoolSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to connect to test MongoDB: %v", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.HealthCheck(ctx); err != nil {
		t.Fatalf("MongoDB health check failed: %v", err)
	}

	return &MongoDBContainer{
		Container: mongoContainer,
		DB:        db,
		URI:       uri,
	}
}

// Teardown cleans up the MongoDB container and closes connections
func (m *MongoDBContainer) Teardown(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	// Close database connection
	if err := m.DB.Close(ctx); err != nil {
		t.Errorf("Failed to close MongoDB connection: %v", err)
	}

	// Terminate container
	if err := m.Container.Terminate(ctx); err != nil {
		t.Errorf("Failed to terminate MongoDB container: %v", err)
	}
}

// CleanupCollections removes all documents from all collections
func (m *MongoDBContainer) CleanupCollections(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	// Get all collection names
	collections, err := m.DB.DB.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to list collections: %v", err)
	}

	// Drop each collection
	for _, collectionName := range collections {
		if err := m.DB.DB.Collection(collectionName).Drop(ctx); err != nil {
			t.Fatalf("Failed to drop collection %s: %v", collectionName, err)
		}
	}
}

// CreateIndexes creates the necessary indexes for testing
func (m *MongoDBContainer) CreateIndexes(t *testing.T, collectionName string, indexes []interface{}) {
	t.Helper()

	_ = context.Background()
	_ = m.DB.Collection(collectionName)
	_ = indexes

	// Note: Index creation can be added here when needed
	// For basic tests, indexes are not required
}

// PrintDatabaseInfo prints database information for debugging
func (m *MongoDBContainer) PrintDatabaseInfo(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	collections, err := m.DB.DB.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		t.Logf("Failed to list collections: %v", err)
		return
	}

	t.Logf("MongoDB URI: %s", m.URI)
	t.Logf("Database: %s", m.DB.DB.Name())
	t.Logf("Collections: %v", collections)

	for _, collectionName := range collections {
		count, err := m.DB.Collection(collectionName).CountDocuments(ctx, map[string]interface{}{})
		if err != nil {
			t.Logf("Failed to count documents in %s: %v", collectionName, err)
			continue
		}
		t.Logf("  %s: %d documents", collectionName, count)
	}
}

// AssertNoError is a helper to assert no error occurred
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertError is a helper to assert an error occurred
func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", message)
	}
}

// AssertEqual is a helper to assert two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if fmt.Sprintf("%v", expected) != fmt.Sprintf("%v", actual) {
		t.Fatalf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertNotNil is a helper to assert a value is not nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Fatalf("%s: expected non-nil value", message)
	}
}

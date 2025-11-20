package repository

import (
	"context"
	"testing"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCircleRepository_Create(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	ownerID := primitive.NewObjectID()

	// Test case: Create a new circle
	t.Run("Create new circle successfully", func(t *testing.T) {
		circle := &models.Circle{
			Name:    "Partner",
			OwnerID: ownerID,
			Members: []models.CircleMember{},
		}

		err := repo.Create(ctx, circle)
		require.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, circle.ID)
		assert.False(t, circle.CreatedAt.IsZero())
		assert.False(t, circle.UpdatedAt.IsZero())
	})

	// Test case: Create circle with members
	t.Run("Create circle with members", func(t *testing.T) {
		memberID := primitive.NewObjectID()
		circle := &models.Circle{
			Name:    "Friends",
			OwnerID: ownerID,
			Members: []models.CircleMember{
				{
					UserID:      memberID,
					AccessLevel: models.AccessLevelView,
					InvitedAt:   time.Now(),
				},
			},
		}

		err := repo.Create(ctx, circle)
		require.NoError(t, err)
		assert.Len(t, circle.Members, 1)
	})
}

func TestCircleRepository_FindByID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Create test circle
	ownerID := primitive.NewObjectID()
	circle := &models.Circle{
		Name:    "Family",
		OwnerID: ownerID,
		Members: []models.CircleMember{},
	}
	require.NoError(t, repo.Create(ctx, circle))

	// Test case: Find existing circle
	t.Run("Find existing circle by ID", func(t *testing.T) {
		found, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Equal(t, circle.Name, found.Name)
		assert.Equal(t, circle.OwnerID, found.OwnerID)
	})

	// Test case: Find non-existent circle
	t.Run("Find non-existent circle returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		found, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestCircleRepository_FindByUserID(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Create test data
	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()
	user3ID := primitive.NewObjectID()

	// Create circles owned by user1
	ownedCircle1 := &models.Circle{
		Name:    "User 1 Circle 1",
		OwnerID: user1ID,
		Members: []models.CircleMember{},
	}
	require.NoError(t, repo.Create(ctx, ownedCircle1))

	ownedCircle2 := &models.Circle{
		Name:    "User 1 Circle 2",
		OwnerID: user1ID,
		Members: []models.CircleMember{},
	}
	require.NoError(t, repo.Create(ctx, ownedCircle2))

	// Create circle owned by user2 with user1 as member
	memberCircle := &models.Circle{
		Name:    "User 2 Circle",
		OwnerID: user2ID,
		Members: []models.CircleMember{
			{
				UserID:      user1ID,
				AccessLevel: models.AccessLevelView,
				InvitedAt:   time.Now(),
			},
		},
	}
	require.NoError(t, repo.Create(ctx, memberCircle))

	// Create circle for user3 (not related to user1)
	unrelatedCircle := &models.Circle{
		Name:    "User 3 Circle",
		OwnerID: user3ID,
		Members: []models.CircleMember{},
	}
	require.NoError(t, repo.Create(ctx, unrelatedCircle))

	// Test case: Find all circles for user1 (owner + member)
	t.Run("Find all circles for user (owner + member)", func(t *testing.T) {
		found, err := repo.FindByUserID(ctx, user1ID)
		require.NoError(t, err)
		assert.Len(t, found, 3) // 2 owned + 1 member

		// Verify user1 is either owner or member in all circles
		for _, circle := range found {
			isOwner := circle.OwnerID == user1ID
			isMember := circle.HasMember(user1ID)
			assert.True(t, isOwner || isMember)
		}
	})

	// Test case: Find circles for user with no circles
	t.Run("Find circles for user with no circles", func(t *testing.T) {
		emptyUserID := primitive.NewObjectID()
		found, err := repo.FindByUserID(ctx, emptyUserID)
		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestCircleRepository_Update(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Create test circle
	ownerID := primitive.NewObjectID()
	circle := &models.Circle{
		Name:    "Original Name",
		OwnerID: ownerID,
		Members: []models.CircleMember{},
	}
	require.NoError(t, repo.Create(ctx, circle))

	// Test case: Update existing circle
	t.Run("Update existing circle successfully", func(t *testing.T) {
		// Save original updated time
		originalUpdatedAt := circle.UpdatedAt

		// Wait a bit to ensure timestamp changes
		time.Sleep(10 * time.Millisecond)

		// Update circle
		circle.Name = "Updated Name"

		err := repo.Update(ctx, circle)
		require.NoError(t, err)
		assert.True(t, circle.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		found, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
	})

	// Test case: Update non-existent circle
	t.Run("Update non-existent circle returns error", func(t *testing.T) {
		nonExistentCircle := &models.Circle{
			ID:      primitive.NewObjectID(),
			Name:    "Ghost Circle",
			OwnerID: ownerID,
		}

		err := repo.Update(ctx, nonExistentCircle)
		assert.Error(t, err)
	})
}

func TestCircleRepository_Delete(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Create test circle
	ownerID := primitive.NewObjectID()
	circle := &models.Circle{
		Name:    "Delete Me",
		OwnerID: ownerID,
		Members: []models.CircleMember{},
	}
	require.NoError(t, repo.Create(ctx, circle))

	// Test case: Delete existing circle
	t.Run("Delete existing circle successfully", func(t *testing.T) {
		err := repo.Delete(ctx, circle.ID)
		require.NoError(t, err)

		// Verify deletion
		found, err := repo.FindByID(ctx, circle.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	// Test case: Delete non-existent circle
	t.Run("Delete non-existent circle returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestCircleRepository_AddMember(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Create test circle
	ownerID := primitive.NewObjectID()
	circle := &models.Circle{
		Name:    "Test Circle",
		OwnerID: ownerID,
		Members: []models.CircleMember{},
	}
	require.NoError(t, repo.Create(ctx, circle))

	// Test case: Add member to circle
	t.Run("Add member to circle successfully", func(t *testing.T) {
		memberID := primitive.NewObjectID()
		member := models.CircleMember{
			UserID:      memberID,
			AccessLevel: models.AccessLevelView,
			InvitedAt:   time.Now(),
		}

		err := repo.AddMember(ctx, circle.ID, member)
		require.NoError(t, err)

		// Verify member was added
		found, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Len(t, found.Members, 1)
		assert.Equal(t, memberID, found.Members[0].UserID)
		assert.Equal(t, models.AccessLevelView, found.Members[0].AccessLevel)
	})

	// Test case: Add multiple members
	t.Run("Add multiple members to circle", func(t *testing.T) {
		member2ID := primitive.NewObjectID()
		member2 := models.CircleMember{
			UserID:      member2ID,
			AccessLevel: models.AccessLevelEdit,
			InvitedAt:   time.Now(),
		}

		err := repo.AddMember(ctx, circle.ID, member2)
		require.NoError(t, err)

		// Verify both members exist
		found, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Len(t, found.Members, 2)
	})

	// Test case: Add member to non-existent circle
	t.Run("Add member to non-existent circle returns error", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		member := models.CircleMember{
			UserID:      primitive.NewObjectID(),
			AccessLevel: models.AccessLevelView,
			InvitedAt:   time.Now(),
		}

		err := repo.AddMember(ctx, nonExistentID, member)
		assert.Error(t, err)
	})
}

func TestCircleRepository_RemoveMember(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Create test circle with members
	ownerID := primitive.NewObjectID()
	member1ID := primitive.NewObjectID()
	member2ID := primitive.NewObjectID()

	circle := &models.Circle{
		Name:    "Test Circle",
		OwnerID: ownerID,
		Members: []models.CircleMember{
			{
				UserID:      member1ID,
				AccessLevel: models.AccessLevelView,
				InvitedAt:   time.Now(),
			},
			{
				UserID:      member2ID,
				AccessLevel: models.AccessLevelEdit,
				InvitedAt:   time.Now(),
			},
		},
	}
	require.NoError(t, repo.Create(ctx, circle))

	// Test case: Remove member from circle
	t.Run("Remove member from circle successfully", func(t *testing.T) {
		err := repo.RemoveMember(ctx, circle.ID, member1ID)
		require.NoError(t, err)

		// Verify member was removed
		found, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Len(t, found.Members, 1)
		assert.Equal(t, member2ID, found.Members[0].UserID)
	})

	// Test case: Remove non-existent member
	t.Run("Remove non-existent member returns error", func(t *testing.T) {
		nonExistentMemberID := primitive.NewObjectID()
		err := repo.RemoveMember(ctx, circle.ID, nonExistentMemberID)
		assert.Error(t, err)
	})

	// Test case: Remove member from non-existent circle
	t.Run("Remove member from non-existent circle returns error", func(t *testing.T) {
		nonExistentCircleID := primitive.NewObjectID()
		err := repo.RemoveMember(ctx, nonExistentCircleID, member2ID)
		assert.Error(t, err)
	})
}

func TestCircleRepository_UpdateMemberAccess(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Create test circle with a member
	ownerID := primitive.NewObjectID()
	memberID := primitive.NewObjectID()

	circle := &models.Circle{
		Name:    "Test Circle",
		OwnerID: ownerID,
		Members: []models.CircleMember{
			{
				UserID:      memberID,
				AccessLevel: models.AccessLevelView,
				InvitedAt:   time.Now(),
			},
		},
	}
	require.NoError(t, repo.Create(ctx, circle))

	// Test case: Update member access level
	t.Run("Update member access level successfully", func(t *testing.T) {
		err := repo.UpdateMemberAccess(ctx, circle.ID, memberID, models.AccessLevelEdit)
		require.NoError(t, err)

		// Verify access level was updated
		found, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Len(t, found.Members, 1)
		assert.Equal(t, models.AccessLevelEdit, found.Members[0].AccessLevel)
	})

	// Test case: Update to admin access
	t.Run("Update member to admin access", func(t *testing.T) {
		err := repo.UpdateMemberAccess(ctx, circle.ID, memberID, models.AccessLevelAdmin)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Equal(t, models.AccessLevelAdmin, found.Members[0].AccessLevel)
	})

	// Test case: Update non-existent member
	t.Run("Update non-existent member returns error", func(t *testing.T) {
		nonExistentMemberID := primitive.NewObjectID()
		err := repo.UpdateMemberAccess(ctx, circle.ID, nonExistentMemberID, models.AccessLevelEdit)
		assert.Error(t, err)
	})

	// Test case: Update member in non-existent circle
	t.Run("Update member in non-existent circle returns error", func(t *testing.T) {
		nonExistentCircleID := primitive.NewObjectID()
		err := repo.UpdateMemberAccess(ctx, nonExistentCircleID, memberID, models.AccessLevelEdit)
		assert.Error(t, err)
	})
}

func TestCircleRepository_Integration(t *testing.T) {
	// Setup
	mongo := testutil.SetupMongoDB(t)
	defer mongo.Teardown(t)

	repo := NewCircleRepository(mongo.DB)
	ctx := context.Background()

	// Integration test: Full lifecycle
	t.Run("Full circle lifecycle", func(t *testing.T) {
		ownerID := primitive.NewObjectID()
		member1ID := primitive.NewObjectID()
		member2ID := primitive.NewObjectID()

		// 1. Create circle
		circle := &models.Circle{
			Name:    "Partner Circle",
			OwnerID: ownerID,
			Members: []models.CircleMember{},
		}
		err := repo.Create(ctx, circle)
		require.NoError(t, err)

		// 2. Find by ID
		found, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Equal(t, circle.Name, found.Name)

		// 3. Find by UserID (as owner)
		userCircles, err := repo.FindByUserID(ctx, ownerID)
		require.NoError(t, err)
		assert.Len(t, userCircles, 1)

		// 4. Add first member
		member1 := models.CircleMember{
			UserID:      member1ID,
			AccessLevel: models.AccessLevelView,
			InvitedAt:   time.Now(),
		}
		err = repo.AddMember(ctx, circle.ID, member1)
		require.NoError(t, err)

		// 5. Add second member
		member2 := models.CircleMember{
			UserID:      member2ID,
			AccessLevel: models.AccessLevelEdit,
			InvitedAt:   time.Now(),
		}
		err = repo.AddMember(ctx, circle.ID, member2)
		require.NoError(t, err)

		// 6. Verify both members
		withMembers, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Len(t, withMembers.Members, 2)

		// 7. Update member1 access level
		err = repo.UpdateMemberAccess(ctx, circle.ID, member1ID, models.AccessLevelEdit)
		require.NoError(t, err)

		updated, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		accessLevel, memberFound := updated.GetMemberAccessLevel(member1ID)
		assert.True(t, memberFound)
		assert.Equal(t, models.AccessLevelEdit, accessLevel)

		// 8. Find by member's UserID
		memberCircles, err := repo.FindByUserID(ctx, member1ID)
		require.NoError(t, err)
		assert.Len(t, memberCircles, 1)

		// 9. Remove member1
		err = repo.RemoveMember(ctx, circle.ID, member1ID)
		require.NoError(t, err)

		afterRemoval, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Len(t, afterRemoval.Members, 1)
		assert.Equal(t, member2ID, afterRemoval.Members[0].UserID)

		// 10. Update circle name
		circle.Name = "Close Friends"
		err = repo.Update(ctx, circle)
		require.NoError(t, err)

		renamed, err := repo.FindByID(ctx, circle.ID)
		require.NoError(t, err)
		assert.Equal(t, "Close Friends", renamed.Name)

		// 11. Delete circle
		err = repo.Delete(ctx, circle.ID)
		require.NoError(t, err)

		// 12. Verify deletion
		deleted, err := repo.FindByID(ctx, circle.ID)
		assert.Error(t, err)
		assert.Nil(t, deleted)
	})
}

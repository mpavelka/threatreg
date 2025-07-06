package service

import (
	"testing"
	"threatreg/internal/database"
	"threatreg/internal/models"
	"threatreg/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTagService_CreateTag(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Test data
		name := "Test Tag"
		description := "A test tag description"
		color := "#FF0000"

		// Create tag
		tag, err := CreateTag(name, description, color)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, tag)
		assert.NotEqual(t, uuid.Nil, tag.ID)
		assert.Equal(t, name, tag.Name)
		assert.Equal(t, description, tag.Description)
		assert.Equal(t, color, tag.Color)

		// Verify tag was actually saved to database
		db := database.GetDB()
		var dbTag models.Tag
		err = db.First(&dbTag, "id = ?", tag.ID).Error
		require.NoError(t, err)
		assert.Equal(t, tag.ID, dbTag.ID)
		assert.Equal(t, name, dbTag.Name)
		assert.Equal(t, description, dbTag.Description)
		assert.Equal(t, color, dbTag.Color)
	})

	t.Run("DuplicateName", func(t *testing.T) {
		// Create a tag
		name := "Duplicate Name Tag"
		_, err := CreateTag(name, "First tag", "#111111")
		require.NoError(t, err)

		// Try to create another tag with the same name
		_, err = CreateTag(name, "Second tag", "#222222")
		assert.Error(t, err) // Should fail due to unique constraint on name
	})
}

func TestTagService_GetTag(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create a tag first
		name := "Get Test Tag"
		description := "Tag for get test"
		color := "#00FF00"
		createdTag, err := CreateTag(name, description, color)
		require.NoError(t, err)

		// Get the tag
		retrievedTag, err := GetTag(createdTag.ID)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedTag)
		assert.Equal(t, createdTag.ID, retrievedTag.ID)
		assert.Equal(t, name, retrievedTag.Name)
		assert.Equal(t, description, retrievedTag.Description)
		assert.Equal(t, color, retrievedTag.Color)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Try to get a non-existent tag
		nonExistentID := uuid.New()
		tag, err := GetTag(nonExistentID)

		// Should return error and nil tag
		assert.Error(t, err)
		assert.Nil(t, tag)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestTagService_GetTagByName(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create a tag first
		name := "GetByName Test Tag"
		description := "Tag for get by name test"
		color := "#0000FF"
		createdTag, err := CreateTag(name, description, color)
		require.NoError(t, err)

		// Get the tag by name
		retrievedTag, err := GetTagByName(name)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, retrievedTag)
		assert.Equal(t, createdTag.ID, retrievedTag.ID)
		assert.Equal(t, name, retrievedTag.Name)
		assert.Equal(t, description, retrievedTag.Description)
		assert.Equal(t, color, retrievedTag.Color)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Try to get a non-existent tag by name
		tag, err := GetTagByName("NonExistentTag")

		// Should return error and nil tag
		assert.Error(t, err)
		assert.Nil(t, tag)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestTagService_UpdateTag(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("FullUpdate", func(t *testing.T) {
		// Create a tag first
		originalName := "Original Tag"
		originalDescription := "Original description"
		originalColor := "#FFFF00"
		createdTag, err := CreateTag(originalName, originalDescription, originalColor)
		require.NoError(t, err)

		// Update the tag
		newName := "Updated Tag"
		newDescription := "Updated description"
		newColor := "#FF00FF"
		updatedTag, err := UpdateTag(createdTag.ID, &newName, &newDescription, &newColor)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedTag)
		assert.Equal(t, createdTag.ID, updatedTag.ID)
		assert.Equal(t, newName, updatedTag.Name)
		assert.Equal(t, newDescription, updatedTag.Description)
		assert.Equal(t, newColor, updatedTag.Color)

		// Verify the update was persisted to database
		db := database.GetDB()
		var dbTag models.Tag
		err = db.First(&dbTag, "id = ?", createdTag.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbTag.Name)
		assert.Equal(t, newDescription, dbTag.Description)
		assert.Equal(t, newColor, dbTag.Color)
	})

	t.Run("PartialUpdate", func(t *testing.T) {
		// Create a tag first
		originalName := "Partial Update Tag"
		originalDescription := "Original description"
		originalColor := "#CCCCCC"
		createdTag, err := CreateTag(originalName, originalDescription, originalColor)
		require.NoError(t, err)

		// Update only the name
		newName := "New Name Only"
		updatedTag, err := UpdateTag(createdTag.ID, &newName, nil, nil)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, updatedTag)
		assert.Equal(t, createdTag.ID, updatedTag.ID)
		assert.Equal(t, newName, updatedTag.Name)
		assert.Equal(t, originalDescription, updatedTag.Description) // Should remain unchanged
		assert.Equal(t, originalColor, updatedTag.Color)             // Should remain unchanged

		// Verify the partial update was persisted
		db := database.GetDB()
		var dbTag models.Tag
		err = db.First(&dbTag, "id = ?", createdTag.ID).Error
		require.NoError(t, err)
		assert.Equal(t, newName, dbTag.Name)
		assert.Equal(t, originalDescription, dbTag.Description)
		assert.Equal(t, originalColor, dbTag.Color)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Try to update a non-existent tag
		nonExistentID := uuid.New()
		newName := "New Name"
		tag, err := UpdateTag(nonExistentID, &newName, nil, nil)

		// Should return error and nil tag
		assert.Error(t, err)
		assert.Nil(t, tag)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestTagService_DeleteTag(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create a tag first
		name := "Delete Test Tag"
		description := "Tag to be deleted"
		color := "#888888"
		createdTag, err := CreateTag(name, description, color)
		require.NoError(t, err)

		// Delete the tag
		err = DeleteTag(createdTag.ID)

		// Assertions
		require.NoError(t, err)

		// Verify the tag was actually deleted from database
		db := database.GetDB()
		var dbTag models.Tag
		err = db.First(&dbTag, "id = ?", createdTag.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Try to delete a non-existent tag
		nonExistentID := uuid.New()
		err := DeleteTag(nonExistentID)

		// Delete should succeed even if tag doesn't exist (GORM behavior)
		assert.NoError(t, err)
	})
}

func TestTagService_ListTags(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("WithTags", func(t *testing.T) {
		// Create multiple tags
		tags := []struct {
			name        string
			description string
			color       string
		}{
			{"Tag 1", "Description 1", "#111111"},
			{"Tag 2", "Description 2", "#222222"},
			{"Tag 3", "Description 3", "#333333"},
		}

		var createdTags []*models.Tag
		for _, tagData := range tags {
			tag, err := CreateTag(tagData.name, tagData.description, tagData.color)
			require.NoError(t, err)
			createdTags = append(createdTags, tag)
		}

		// List all tags
		retrievedTags, err := ListTags()

		// Assertions
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(retrievedTags), len(tags))

		// Verify all created tags are in the list
		tagMap := make(map[uuid.UUID]models.Tag)
		for _, retrievedTag := range retrievedTags {
			tagMap[retrievedTag.ID] = retrievedTag
		}

		for _, created := range createdTags {
			retrieved, exists := tagMap[created.ID]
			assert.True(t, exists, "Created tag should exist in list")
			assert.Equal(t, created.Name, retrieved.Name)
			assert.Equal(t, created.Description, retrieved.Description)
			assert.Equal(t, created.Color, retrieved.Color)
		}
	})

}

func TestTagService_ListTags_Empty(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("NoTags", func(t *testing.T) {
		// List tags when no tags exist in this test database
		tags, err := ListTags()

		// Should return empty slice, not error
		require.NoError(t, err)
		assert.Len(t, tags, 0)
	})
}

func TestTagService_AssignTagToProduct(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data
		product, err := CreateProduct("Test Product", "A test product for tagging")
		require.NoError(t, err)

		tag, err := CreateTag("Product Tag", "A tag for products", "#AAAAAA")
		require.NoError(t, err)

		// Assign tag to product
		err = AssignTagToProduct(tag.ID, product.ID)
		require.NoError(t, err)

		// Verify assignment
		tags, err := ListTagsByProductID(product.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, tag.ID, tags[0].ID)
		assert.Equal(t, tag.Name, tags[0].Name)

		// Verify reverse lookup
		products, err := ListProductsByTagID(tag.ID)
		require.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Equal(t, product.ID, products[0].ID)
		assert.Equal(t, product.Name, products[0].Name)
	})

	t.Run("Duplicate", func(t *testing.T) {
		// Create test data
		product, err := CreateProduct("Duplicate Test Product", "A test product for duplicate tagging")
		require.NoError(t, err)

		tag, err := CreateTag("Duplicate Tag", "A tag for duplicate test", "#BBBBBB")
		require.NoError(t, err)

		// Assign tag to product first time
		err = AssignTagToProduct(tag.ID, product.ID)
		require.NoError(t, err)

		// Try to assign the same tag to the same product again
		err = AssignTagToProduct(tag.ID, product.ID)
		require.NoError(t, err) // Should not error

		// Verify only one assignment exists
		tags, err := ListTagsByProductID(product.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, tag.ID, tags[0].ID)
	})

	t.Run("InvalidTagID", func(t *testing.T) {
		// Create test product
		product, err := CreateProduct("Invalid Tag Test Product", "A test product")
		require.NoError(t, err)

		// Try to assign non-existent tag
		nonExistentTagID := uuid.New()
		err = AssignTagToProduct(nonExistentTagID, product.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tag not found")
	})

	t.Run("InvalidProductID", func(t *testing.T) {
		// Create test tag
		tag, err := CreateTag("Invalid Product Test Tag", "A test tag", "#DDDDDD")
		require.NoError(t, err)

		// Try to assign to non-existent product
		nonExistentProductID := uuid.New()
		err = AssignTagToProduct(tag.ID, nonExistentProductID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product not found")
	})
}

func TestTagService_UnassignTagFromProduct(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data and assign tag
		product, err := CreateProduct("Unassign Test Product", "A test product for unassigning tags")
		require.NoError(t, err)

		tag, err := CreateTag("Unassign Tag", "A tag for unassign test", "#EEEEEE")
		require.NoError(t, err)

		err = AssignTagToProduct(tag.ID, product.ID)
		require.NoError(t, err)

		// Verify assignment exists
		tags, err := ListTagsByProductID(product.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)

		// Unassign tag from product
		err = UnassignTagFromProduct(tag.ID, product.ID)
		require.NoError(t, err)

		// Verify assignment was removed
		tags, err = ListTagsByProductID(product.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 0)
	})
}

func TestTagService_AssignTagToInstance(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data
		product, err := CreateProduct("Instance Tag Product", "A product for instance tagging")
		require.NoError(t, err)

		instance, err := CreateInstance("Test Instance", product.ID)
		require.NoError(t, err)

		tag, err := CreateTag("Instance Tag", "A tag for instances", "#FFFFFF")
		require.NoError(t, err)

		// Assign tag to instance
		err = AssignTagToInstance(tag.ID, instance.ID)
		require.NoError(t, err)

		// Verify assignment
		tags, err := ListTagsByInstanceID(instance.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, tag.ID, tags[0].ID)
		assert.Equal(t, tag.Name, tags[0].Name)

		// Verify reverse lookup
		instances, err := ListInstancesByTagID(tag.ID)
		require.NoError(t, err)
		assert.Len(t, instances, 1)
		assert.Equal(t, instance.ID, instances[0].ID)
		assert.Equal(t, instance.Name, instances[0].Name)
	})

	t.Run("Duplicate", func(t *testing.T) {
		// Create test data
		product, err := CreateProduct("Duplicate Instance Product", "A product for duplicate instance tagging")
		require.NoError(t, err)

		instance, err := CreateInstance("Duplicate Test Instance", product.ID)
		require.NoError(t, err)

		tag, err := CreateTag("Duplicate Instance Tag", "A tag for duplicate instance test", "#000000")
		require.NoError(t, err)

		// Assign tag to instance first time
		err = AssignTagToInstance(tag.ID, instance.ID)
		require.NoError(t, err)

		// Try to assign the same tag to the same instance again
		err = AssignTagToInstance(tag.ID, instance.ID)
		require.NoError(t, err) // Should not error

		// Verify only one assignment exists
		tags, err := ListTagsByInstanceID(instance.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, tag.ID, tags[0].ID)
	})
}

func TestTagService_UnassignTagFromInstance(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		// Create test data and assign tag
		product, err := CreateProduct("Unassign Instance Product", "A product for unassigning instance tags")
		require.NoError(t, err)

		instance, err := CreateInstance("Unassign Test Instance", product.ID)
		require.NoError(t, err)

		tag, err := CreateTag("Unassign Instance Tag", "A tag for instance unassign test", "#999999")
		require.NoError(t, err)

		err = AssignTagToInstance(tag.ID, instance.ID)
		require.NoError(t, err)

		// Verify assignment exists
		tags, err := ListTagsByInstanceID(instance.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)

		// Unassign tag from instance
		err = UnassignTagFromInstance(tag.ID, instance.ID)
		require.NoError(t, err)

		// Verify assignment was removed
		tags, err = ListTagsByInstanceID(instance.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 0)
	})
}

func TestTagService_AssignTagByName(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Product{},
		&models.Instance{},
		&models.Tag{},
	)
	defer cleanup()

	t.Run("AssignToProduct_ExistingTag", func(t *testing.T) {
		// Create test data
		product, err := CreateProduct("ByName Test Product", "A test product for tag assignment by name")
		require.NoError(t, err)

		tag, err := CreateTag("ByName Tag", "A tag for by name test", "#ABCDEF")
		require.NoError(t, err)

		// Assign tag to product by name
		err = AssignTagToProductByName(tag.Name, product.ID)
		require.NoError(t, err)

		// Verify assignment
		tags, err := ListTagsByProductID(product.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, tag.ID, tags[0].ID)
		assert.Equal(t, tag.Name, tags[0].Name)
	})

	t.Run("AssignToProduct_NewTag", func(t *testing.T) {
		// Create test data
		product, err := CreateProduct("NewTag Test Product", "A test product for new tag creation")
		require.NoError(t, err)

		tagName := "New Auto-Created Tag"

		// Assign non-existent tag to product by name (should create the tag)
		err = AssignTagToProductByName(tagName, product.ID)
		require.NoError(t, err)

		// Verify tag was created and assigned
		tags, err := ListTagsByProductID(product.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, tagName, tags[0].Name)
		assert.Empty(t, tags[0].Description) // Should be empty for auto-created tag
		assert.Empty(t, tags[0].Color)       // Should be empty for auto-created tag
	})

	t.Run("AssignToInstance_ExistingTag", func(t *testing.T) {
		// Create test data
		product, err := CreateProduct("Instance ByName Product", "A product for instance tag assignment by name")
		require.NoError(t, err)

		instance, err := CreateInstance("ByName Test Instance", product.ID)
		require.NoError(t, err)

		tag, err := CreateTag("Instance ByName Tag", "A tag for instance by name test", "#FEDCBA")
		require.NoError(t, err)

		// Assign tag to instance by name
		err = AssignTagToInstanceByName(tag.Name, instance.ID)
		require.NoError(t, err)

		// Verify assignment
		tags, err := ListTagsByInstanceID(instance.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, tag.ID, tags[0].ID)
		assert.Equal(t, tag.Name, tags[0].Name)
	})

	t.Run("AssignToInstance_NewTag", func(t *testing.T) {
		// Create test data
		product, err := CreateProduct("Instance NewTag Product", "A product for new instance tag creation")
		require.NoError(t, err)

		instance, err := CreateInstance("NewTag Test Instance", product.ID)
		require.NoError(t, err)

		tagName := "New Instance Auto-Created Tag"

		// Assign non-existent tag to instance by name (should create the tag)
		err = AssignTagToInstanceByName(tagName, instance.ID)
		require.NoError(t, err)

		// Verify tag was created and assigned
		tags, err := ListTagsByInstanceID(instance.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, tagName, tags[0].Name)
		assert.Empty(t, tags[0].Description) // Should be empty for auto-created tag
		assert.Empty(t, tags[0].Color)       // Should be empty for auto-created tag
	})
}

package service

import (
	"testing"
	"threatreg/internal/models"
	"threatreg/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestComponentAttributeService_CreateComponentAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	// Create test component
	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("CreateStringAttribute", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, "environment", models.ComponentAttributeTypeString, "production")
		
		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.NotEqual(t, uuid.Nil, attribute.ID)
		assert.Equal(t, component.ID, attribute.ComponentID)
		assert.Equal(t, "environment", attribute.Name)
		assert.Equal(t, models.ComponentAttributeTypeString, attribute.Type)
		assert.Equal(t, "production", attribute.Value)
	})

	t.Run("CreateTextAttribute", func(t *testing.T) {
		longText := "This is a long description that would be stored as text type"
		attribute, err := CreateComponentAttribute(component.ID, "description", models.ComponentAttributeTypeText, longText)
		
		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, models.ComponentAttributeTypeText, attribute.Type)
		assert.Equal(t, longText, attribute.Value)
	})

	t.Run("CreateNumberAttribute_Integer", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, "port", models.ComponentAttributeTypeNumber, "8080")
		
		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, models.ComponentAttributeTypeNumber, attribute.Type)
		assert.Equal(t, "8080", attribute.Value)
	})

	t.Run("CreateNumberAttribute_Float", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, "version", models.ComponentAttributeTypeNumber, "1.5")
		
		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, models.ComponentAttributeTypeNumber, attribute.Type)
		assert.Equal(t, "1.5", attribute.Value)
	})

	t.Run("CreateComponentAttribute_ValidReference", func(t *testing.T) {
		// Create another component to reference
		referencedComponent, err := CreateComponent("Referenced Component", "Referenced Description", models.ComponentTypeProduct)
		require.NoError(t, err)

		attribute, err := CreateComponentAttribute(component.ID, "depends_on", models.ComponentAttributeTypeComponent, referencedComponent.ID.String())
		
		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, models.ComponentAttributeTypeComponent, attribute.Type)
		assert.Equal(t, referencedComponent.ID.String(), attribute.Value)
	})

	t.Run("CreateComponentAttribute_InvalidReference", func(t *testing.T) {
		nonExistentID := uuid.New()
		attribute, err := CreateComponentAttribute(component.ID, "depends_on", models.ComponentAttributeTypeComponent, nonExistentID.String())
		
		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "referenced component does not exist")
	})

	t.Run("CreateNumberAttribute_InvalidValue", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, "invalid_port", models.ComponentAttributeTypeNumber, "not-a-number")
		
		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "number attribute value must be a valid integer or float")
	})

	t.Run("CreateComponentAttribute_InvalidUUID", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, "depends_on", models.ComponentAttributeTypeComponent, "not-a-uuid")
		
		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "invalid UUID")
	})

	t.Run("CreateAttribute_EmptyName", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, "", models.ComponentAttributeTypeString, "value")
		
		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "attribute name is required")
	})

	t.Run("CreateAttribute_EmptyValue", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, "empty_value", models.ComponentAttributeTypeString, "")
		
		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "attribute value is required")
	})
}

func TestComponentAttributeService_GetComponentAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("GetExistingAttribute", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, "environment", models.ComponentAttributeTypeString, "staging")
		require.NoError(t, err)

		retrievedAttribute, err := GetComponentAttribute(createdAttribute.ID)
		
		require.NoError(t, err)
		assert.NotNil(t, retrievedAttribute)
		assert.Equal(t, createdAttribute.ID, retrievedAttribute.ID)
		assert.Equal(t, createdAttribute.Name, retrievedAttribute.Name)
		assert.Equal(t, createdAttribute.Type, retrievedAttribute.Type)
		assert.Equal(t, createdAttribute.Value, retrievedAttribute.Value)
	})

	t.Run("GetNonExistentAttribute", func(t *testing.T) {
		nonExistentID := uuid.New()
		attribute, err := GetComponentAttribute(nonExistentID)
		
		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestComponentAttributeService_UpdateComponentAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("UpdateAttributeValue", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, "environment", models.ComponentAttributeTypeString, "development")
		require.NoError(t, err)

		newValue := "production"
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, nil, &newValue)
		
		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, createdAttribute.ID, updatedAttribute.ID)
		assert.Equal(t, "production", updatedAttribute.Value)
		assert.Equal(t, createdAttribute.Name, updatedAttribute.Name)
		assert.Equal(t, createdAttribute.Type, updatedAttribute.Type)
	})

	t.Run("UpdateAttributeType", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, "config", models.ComponentAttributeTypeString, "8080")
		require.NoError(t, err)

		newType := models.ComponentAttributeTypeNumber
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, &newType, nil)
		
		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, models.ComponentAttributeTypeNumber, updatedAttribute.Type)
		assert.Equal(t, "8080", updatedAttribute.Value) // Should still be valid as number
	})

	t.Run("UpdateAttributeName", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, "old_name", models.ComponentAttributeTypeString, "value")
		require.NoError(t, err)

		newName := "new_name"
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, &newName, nil, nil)
		
		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, "new_name", updatedAttribute.Name)
	})

	t.Run("UpdateToComponentType_ValidReference", func(t *testing.T) {
		// Create referenced component
		referencedComponent, err := CreateComponent("Referenced Component", "Referenced Description", models.ComponentTypeProduct)
		require.NoError(t, err)

		createdAttribute, err := CreateComponentAttribute(component.ID, "reference", models.ComponentAttributeTypeString, "old_value")
		require.NoError(t, err)

		newType := models.ComponentAttributeTypeComponent
		newValue := referencedComponent.ID.String()
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, &newType, &newValue)
		
		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, models.ComponentAttributeTypeComponent, updatedAttribute.Type)
		assert.Equal(t, referencedComponent.ID.String(), updatedAttribute.Value)
	})

	t.Run("UpdateToComponentType_InvalidReference", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, "reference", models.ComponentAttributeTypeString, "old_value")
		require.NoError(t, err)

		newType := models.ComponentAttributeTypeComponent
		nonExistentID := uuid.New().String()
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, &newType, &nonExistentID)
		
		assert.Error(t, err)
		assert.Nil(t, updatedAttribute)
		assert.Contains(t, err.Error(), "referenced component does not exist")
	})

	t.Run("UpdateToNumberType_InvalidValue", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, "numeric", models.ComponentAttributeTypeString, "old_value")
		require.NoError(t, err)

		newType := models.ComponentAttributeTypeNumber
		newValue := "not-a-number"
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, &newType, &newValue)
		
		assert.Error(t, err)
		assert.Nil(t, updatedAttribute)
		assert.Contains(t, err.Error(), "number attribute value must be a valid integer or float")
	})

	t.Run("UpdateNonExistentAttribute", func(t *testing.T) {
		nonExistentID := uuid.New()
		newValue := "new_value"
		updatedAttribute, err := UpdateComponentAttribute(nonExistentID, nil, nil, &newValue)
		
		assert.Error(t, err)
		assert.Nil(t, updatedAttribute)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestComponentAttributeService_DeleteComponentAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("DeleteExistingAttribute", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, "temp_attr", models.ComponentAttributeTypeString, "temp_value")
		require.NoError(t, err)

		err = DeleteComponentAttribute(createdAttribute.ID)
		require.NoError(t, err)

		// Verify it's deleted
		retrievedAttribute, err := GetComponentAttribute(createdAttribute.ID)
		assert.Error(t, err)
		assert.Nil(t, retrievedAttribute)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteNonExistentAttribute", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := DeleteComponentAttribute(nonExistentID)
		
		// Should not return error for non-existent records (GORM behavior)
		assert.NoError(t, err)
	})
}

func TestComponentAttributeService_GetComponentAttributes(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("GetAttributesForComponentWithAttributes", func(t *testing.T) {
		// Create multiple attributes
		_, err := CreateComponentAttribute(component.ID, "environment", models.ComponentAttributeTypeString, "production")
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component.ID, "port", models.ComponentAttributeTypeNumber, "8080")
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component.ID, "description", models.ComponentAttributeTypeText, "Long description")
		require.NoError(t, err)

		attributes, err := GetComponentAttributes(component.ID)
		
		require.NoError(t, err)
		assert.Len(t, attributes, 3)

		// Create a map for easier assertion
		attributeMap := make(map[string]models.ComponentAttribute)
		for _, attr := range attributes {
			attributeMap[attr.Name] = attr
		}

		assert.Contains(t, attributeMap, "environment")
		assert.Contains(t, attributeMap, "port")
		assert.Contains(t, attributeMap, "description")
		assert.Equal(t, "production", attributeMap["environment"].Value)
		assert.Equal(t, "8080", attributeMap["port"].Value)
		assert.Equal(t, models.ComponentAttributeTypeNumber, attributeMap["port"].Type)
	})

	t.Run("GetAttributesForComponentWithNoAttributes", func(t *testing.T) {
		emptyComponent, err := CreateComponent("Empty Component", "Component with no attributes", models.ComponentTypeInstance)
		require.NoError(t, err)

		attributes, err := GetComponentAttributes(emptyComponent.ID)
		
		require.NoError(t, err)
		assert.Len(t, attributes, 0)
	})
}

func TestComponentAttributeService_GetComponentAttributeByName(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("GetExistingAttributeByName", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, "environment", models.ComponentAttributeTypeString, "staging")
		require.NoError(t, err)

		retrievedAttribute, err := GetComponentAttributeByName(component.ID, "environment")
		
		require.NoError(t, err)
		assert.NotNil(t, retrievedAttribute)
		assert.Equal(t, createdAttribute.ID, retrievedAttribute.ID)
		assert.Equal(t, "environment", retrievedAttribute.Name)
		assert.Equal(t, "staging", retrievedAttribute.Value)
	})

	t.Run("GetNonExistentAttributeByName", func(t *testing.T) {
		attribute, err := GetComponentAttributeByName(component.ID, "non_existent")
		
		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestComponentAttributeService_FindComponentsByAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	// Create test components
	component1, err := CreateComponent("Component 1", "First component", models.ComponentTypeInstance)
	require.NoError(t, err)
	component2, err := CreateComponent("Component 2", "Second component", models.ComponentTypeInstance)
	require.NoError(t, err)
	component3, err := CreateComponent("Component 3", "Third component", models.ComponentTypeProduct)
	require.NoError(t, err)

	t.Run("FindComponentsByAttributeValue", func(t *testing.T) {
		// Add same attribute to multiple components
		_, err := CreateComponentAttribute(component1.ID, "environment", models.ComponentAttributeTypeString, "production")
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component2.ID, "environment", models.ComponentAttributeTypeString, "production")
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component3.ID, "environment", models.ComponentAttributeTypeString, "development")
		require.NoError(t, err)

		components, err := FindComponentsByAttribute("environment", "production")
		
		require.NoError(t, err)
		assert.Len(t, components, 2)

		componentIDs := []uuid.UUID{components[0].ID, components[1].ID}
		assert.Contains(t, componentIDs, component1.ID)
		assert.Contains(t, componentIDs, component2.ID)
		assert.NotContains(t, componentIDs, component3.ID)
	})

	t.Run("FindComponentsByNonExistentAttribute", func(t *testing.T) {
		components, err := FindComponentsByAttribute("non_existent", "value")
		
		require.NoError(t, err)
		assert.Len(t, components, 0)
	})
}

func TestComponentAttributeService_FindComponentsByAttributeAndType(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component1, err := CreateComponent("Component 1", "First component", models.ComponentTypeInstance)
	require.NoError(t, err)
	component2, err := CreateComponent("Component 2", "Second component", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("FindComponentsByAttributeAndType", func(t *testing.T) {
		// Add same value but different types
		_, err := CreateComponentAttribute(component1.ID, "port", models.ComponentAttributeTypeString, "8080")
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component2.ID, "port", models.ComponentAttributeTypeNumber, "8080")
		require.NoError(t, err)

		// Find only number type
		components, err := FindComponentsByAttributeAndType("port", "8080", models.ComponentAttributeTypeNumber)
		
		require.NoError(t, err)
		assert.Len(t, components, 1)
		assert.Equal(t, component2.ID, components[0].ID)
	})
}

func TestComponentAttributeService_ComponentHasAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("ComponentHasExistingAttribute", func(t *testing.T) {
		_, err := CreateComponentAttribute(component.ID, "environment", models.ComponentAttributeTypeString, "production")
		require.NoError(t, err)

		hasAttribute, err := ComponentHasAttribute(component.ID, "environment")
		
		require.NoError(t, err)
		assert.True(t, hasAttribute)
	})

	t.Run("ComponentDoesNotHaveAttribute", func(t *testing.T) {
		hasAttribute, err := ComponentHasAttribute(component.ID, "non_existent")
		
		require.NoError(t, err)
		assert.False(t, hasAttribute)
	})
}

func TestComponentAttributeService_ComponentHasAttributeWithValue(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("ComponentHasAttributeWithCorrectValue", func(t *testing.T) {
		_, err := CreateComponentAttribute(component.ID, "environment", models.ComponentAttributeTypeString, "production")
		require.NoError(t, err)

		hasAttribute, err := ComponentHasAttributeWithValue(component.ID, "environment", "production")
		
		require.NoError(t, err)
		assert.True(t, hasAttribute)
	})

	t.Run("ComponentHasAttributeWithDifferentValue", func(t *testing.T) {
		hasAttribute, err := ComponentHasAttributeWithValue(component.ID, "environment", "staging")
		
		require.NoError(t, err)
		assert.False(t, hasAttribute)
	})
}

func TestComponentAttributeService_SetComponentAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("SetAttributeOnNewComponent", func(t *testing.T) {
		attribute, err := SetComponentAttribute(component.ID, "environment", models.ComponentAttributeTypeString, "production")
		
		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, "environment", attribute.Name)
		assert.Equal(t, "production", attribute.Value)
	})

	t.Run("SetAttributeOnExistingAttribute", func(t *testing.T) {
		// First call creates the attribute
		_, err := SetComponentAttribute(component.ID, "version", models.ComponentAttributeTypeString, "1.0")
		require.NoError(t, err)

		// Second call should update it
		updatedAttribute, err := SetComponentAttribute(component.ID, "version", models.ComponentAttributeTypeString, "2.0")
		
		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, "version", updatedAttribute.Name)
		assert.Equal(t, "2.0", updatedAttribute.Value)

		// Verify only one attribute exists
		attributes, err := GetComponentAttributes(component.ID)
		require.NoError(t, err)
		
		versionCount := 0
		for _, attr := range attributes {
			if attr.Name == "version" {
				versionCount++
				assert.Equal(t, "2.0", attr.Value)
			}
		}
		assert.Equal(t, 1, versionCount, "Should have exactly one version attribute")
	})
}

func TestComponentAttributeService_DeleteComponentAttributeByName(t *testing.T) {
	cleanup := testutil.SetupTestDatabaseWithCustomModels(t,
		&models.Component{},
		&models.ComponentAttribute{},
	)
	defer cleanup()

	component, err := CreateComponent("Test Component", "Test Description", models.ComponentTypeInstance)
	require.NoError(t, err)

	t.Run("DeleteExistingAttributeByName", func(t *testing.T) {
		_, err := CreateComponentAttribute(component.ID, "temp_attr", models.ComponentAttributeTypeString, "temp_value")
		require.NoError(t, err)

		err = DeleteComponentAttributeByName(component.ID, "temp_attr")
		require.NoError(t, err)

		// Verify it's deleted
		hasAttribute, err := ComponentHasAttribute(component.ID, "temp_attr")
		require.NoError(t, err)
		assert.False(t, hasAttribute)
	})

	t.Run("DeleteNonExistentAttributeByName", func(t *testing.T) {
		err := DeleteComponentAttributeByName(component.ID, "non_existent")
		
		// Should not return error for non-existent records
		assert.NoError(t, err)
	})
}
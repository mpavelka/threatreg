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
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	// Create test component
	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("CreateStringAttribute", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("environment"), models.ComponentAttributeTypeString, testutil.AddRandSuffix("production"))

		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.NotEqual(t, uuid.Nil, attribute.ID)
		assert.Equal(t, component.ID, attribute.ComponentID)
		assert.Contains(t, attribute.Name, "environment")
		assert.Equal(t, models.ComponentAttributeTypeString, attribute.Type)
		assert.Contains(t, attribute.Value, "production")
	})

	t.Run("CreateTextAttribute", func(t *testing.T) {
		longText := testutil.AddRandSuffix("This is a long description that would be stored as text type")
		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("description"), models.ComponentAttributeTypeText, longText)

		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, models.ComponentAttributeTypeText, attribute.Type)
		assert.Equal(t, longText, attribute.Value)
	})

	t.Run("CreateNumberAttribute_Integer", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("port"), models.ComponentAttributeTypeNumber, "8080")

		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, models.ComponentAttributeTypeNumber, attribute.Type)
		assert.Equal(t, "8080", attribute.Value)
	})

	t.Run("CreateNumberAttribute_Float", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("version"), models.ComponentAttributeTypeNumber, "1.5")

		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, models.ComponentAttributeTypeNumber, attribute.Type)
		assert.Equal(t, "1.5", attribute.Value)
	})

	t.Run("CreateComponentAttribute_ValidReference", func(t *testing.T) {
		// Create another component to reference
		referencedComponent, err := CreateComponent(
			testutil.AddRandSuffix("Referenced Component"),
			testutil.AddRandSuffix("Referenced Description"),
			models.ComponentTypeProduct,
		)
		require.NoError(t, err)

		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("depends_on"), models.ComponentAttributeTypeComponent, referencedComponent.ID.String())

		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, models.ComponentAttributeTypeComponent, attribute.Type)
		assert.Equal(t, referencedComponent.ID.String(), attribute.Value)
	})

	t.Run("CreateComponentAttribute_InvalidReference", func(t *testing.T) {
		nonExistentID := uuid.New()
		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("depends_on"), models.ComponentAttributeTypeComponent, nonExistentID.String())

		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "referenced component does not exist")
	})

	t.Run("CreateNumberAttribute_InvalidValue", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("invalid_port"), models.ComponentAttributeTypeNumber, testutil.AddRandSuffix("not-a-number"))

		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "number attribute value must be a valid integer or float")
	})

	t.Run("CreateComponentAttribute_InvalidUUID", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("depends_on"), models.ComponentAttributeTypeComponent, testutil.AddRandSuffix("not-a-uuid"))

		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "invalid UUID")
	})

	t.Run("CreateAttribute_EmptyName", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, "", models.ComponentAttributeTypeString, testutil.AddRandSuffix("value"))

		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "attribute name is required")
	})

	t.Run("CreateAttribute_EmptyValue", func(t *testing.T) {
		attribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("empty_value"), models.ComponentAttributeTypeString, "")

		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Contains(t, err.Error(), "attribute value is required")
	})
}

func TestComponentAttributeService_GetComponentAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("GetExistingAttribute", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("environment"), models.ComponentAttributeTypeString, testutil.AddRandSuffix("staging"))
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
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("UpdateAttributeValue", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("environment"), models.ComponentAttributeTypeString, testutil.AddRandSuffix("development"))
		require.NoError(t, err)

		newValue := testutil.AddRandSuffix("production")
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, nil, &newValue)

		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, createdAttribute.ID, updatedAttribute.ID)
		assert.Equal(t, newValue, updatedAttribute.Value)
		assert.Equal(t, createdAttribute.Name, updatedAttribute.Name)
		assert.Equal(t, createdAttribute.Type, updatedAttribute.Type)
	})

	t.Run("UpdateAttributeType", func(t *testing.T) {
		attrConfig := testutil.AddRandSuffix("config")
		createdAttribute, err := CreateComponentAttribute(component.ID, attrConfig, models.ComponentAttributeTypeString, "8080")
		require.NoError(t, err)

		newType := models.ComponentAttributeTypeNumber
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, &newType, nil)

		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, models.ComponentAttributeTypeNumber, updatedAttribute.Type)
		assert.Contains(t, updatedAttribute.Value, "8080") // Should still be valid as number
	})

	t.Run("UpdateAttributeName", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("old_name"), models.ComponentAttributeTypeString, testutil.AddRandSuffix("value"))
		require.NoError(t, err)

		newName := testutil.AddRandSuffix("new_name")
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, &newName, nil, nil)

		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, newName, updatedAttribute.Name)
	})

	t.Run("UpdateToComponentType_ValidReference", func(t *testing.T) {
		// Create referenced component
		referencedComponent, err := CreateComponent(
			testutil.AddRandSuffix("Referenced Component"),
			testutil.AddRandSuffix("Referenced Description"),
			models.ComponentTypeProduct,
		)
		require.NoError(t, err)

		createdAttribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("reference"), models.ComponentAttributeTypeString, testutil.AddRandSuffix("old_value"))
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
		createdAttribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("reference"), models.ComponentAttributeTypeString, testutil.AddRandSuffix("old_value"))
		require.NoError(t, err)

		newType := models.ComponentAttributeTypeComponent
		nonExistentID := uuid.New().String()
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, &newType, &nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, updatedAttribute)
		assert.Contains(t, err.Error(), "referenced component does not exist")
	})

	t.Run("UpdateToNumberType_InvalidValue", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("numeric"), models.ComponentAttributeTypeString, testutil.AddRandSuffix("old_value"))
		require.NoError(t, err)

		newType := models.ComponentAttributeTypeNumber
		newValue := testutil.AddRandSuffix("not-a-number")
		updatedAttribute, err := UpdateComponentAttribute(createdAttribute.ID, nil, &newType, &newValue)

		assert.Error(t, err)
		assert.Nil(t, updatedAttribute)
		assert.Contains(t, err.Error(), "number attribute value must be a valid integer or float")
	})

	t.Run("UpdateNonExistentAttribute", func(t *testing.T) {
		nonExistentID := uuid.New()
		newValue := testutil.AddRandSuffix("new_value")
		updatedAttribute, err := UpdateComponentAttribute(nonExistentID, nil, nil, &newValue)

		assert.Error(t, err)
		assert.Nil(t, updatedAttribute)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestComponentAttributeService_DeleteComponentAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("DeleteExistingAttribute", func(t *testing.T) {
		createdAttribute, err := CreateComponentAttribute(component.ID, testutil.AddRandSuffix("temp_attr"), models.ComponentAttributeTypeString, testutil.AddRandSuffix("temp_value"))
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
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("GetAttributesForComponentWithAttributes", func(t *testing.T) {
		attrEnvironment := testutil.AddRandSuffix("environment")
		attrPort := testutil.AddRandSuffix("port")
		attrDescription := testutil.AddRandSuffix("description")

		// Create multiple attributes
		_, err := CreateComponentAttribute(component.ID, attrEnvironment, models.ComponentAttributeTypeString, "production")
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component.ID, attrPort, models.ComponentAttributeTypeNumber, "8080")
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component.ID, attrDescription, models.ComponentAttributeTypeText, "Long description")
		require.NoError(t, err)

		attributes, err := GetComponentAttributes(component.ID)

		require.NoError(t, err)
		assert.Len(t, attributes, 3)

		// Create a map for easier assertion
		attributeMap := make(map[string]models.ComponentAttribute)
		for _, attr := range attributes {
			attributeMap[attr.Name] = attr
		}

		assert.Contains(t, attributeMap, attrEnvironment)
		assert.Contains(t, attributeMap, attrPort)
		assert.Contains(t, attributeMap, attrDescription)
		assert.Equal(t, "production", attributeMap[attrEnvironment].Value)
		assert.Equal(t, "8080", attributeMap[attrPort].Value)
		assert.Equal(t, models.ComponentAttributeTypeNumber, attributeMap[attrPort].Type)

	})

	t.Run("GetAttributesForComponentWithNoAttributes", func(t *testing.T) {
		emptyComponent, err := CreateComponent(
			testutil.AddRandSuffix("Empty Component"),
			testutil.AddRandSuffix("Component with no attributes"),
			models.ComponentTypeInstance,
		)
		require.NoError(t, err)

		attributes, err := GetComponentAttributes(emptyComponent.ID)

		require.NoError(t, err)
		assert.Len(t, attributes, 0)
	})
}

func TestComponentAttributeService_GetComponentAttributeByName(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("GetExistingAttributeByName", func(t *testing.T) {
		attrEnvironment := testutil.AddRandSuffix("environment")
		createdAttribute, err := CreateComponentAttribute(component.ID, attrEnvironment, models.ComponentAttributeTypeString, "staging")
		require.NoError(t, err)

		retrievedAttribute, err := GetComponentAttributeByName(component.ID, createdAttribute.Name)

		require.NoError(t, err)
		assert.NotNil(t, retrievedAttribute)
		assert.Equal(t, createdAttribute.ID, retrievedAttribute.ID)
		assert.Equal(t, createdAttribute.Name, retrievedAttribute.Name)
		assert.Equal(t, "staging", retrievedAttribute.Value)
	})

	t.Run("GetNonExistentAttributeByName", func(t *testing.T) {
		attribute, err := GetComponentAttributeByName(component.ID, testutil.AddRandSuffix("non_existent"))

		assert.Error(t, err)
		assert.Nil(t, attribute)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestComponentAttributeService_FindComponentsByAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	// Create test components
	component1, err := CreateComponent(
		testutil.AddRandSuffix("Component 1"),
		testutil.AddRandSuffix("First component"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)
	component2, err := CreateComponent(
		testutil.AddRandSuffix("Component 2"),
		testutil.AddRandSuffix("Second component"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)
	component3, err := CreateComponent(
		testutil.AddRandSuffix("Component 3"),
		testutil.AddRandSuffix("Third component"),
		models.ComponentTypeProduct,
	)
	require.NoError(t, err)

	t.Run("FindComponentsByAttributeValue", func(t *testing.T) {
		// Add same attribute to multiple components
		valProduction := testutil.AddRandSuffix("production")
		valDevelopment := testutil.AddRandSuffix("development")
		attrEnvironment := testutil.AddRandSuffix("environment")

		_, err := CreateComponentAttribute(component1.ID, attrEnvironment, models.ComponentAttributeTypeString, valProduction)
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component2.ID, attrEnvironment, models.ComponentAttributeTypeString, valProduction)
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component3.ID, attrEnvironment, models.ComponentAttributeTypeString, valDevelopment)
		require.NoError(t, err)

		components, err := FindComponentsByAttribute(attrEnvironment, valProduction)

		require.NoError(t, err)
		assert.Len(t, components, 2)

		componentIDs := []uuid.UUID{components[0].ID, components[1].ID}
		assert.Contains(t, componentIDs, component1.ID)
		assert.Contains(t, componentIDs, component2.ID)
		assert.NotContains(t, componentIDs, component3.ID)
	})

	t.Run("FindComponentsByNonExistentAttribute", func(t *testing.T) {
		components, err := FindComponentsByAttribute(
			testutil.AddRandSuffix("non_existent"),
			testutil.AddRandSuffix("value"),
		)

		require.NoError(t, err)
		assert.Len(t, components, 0)
	})
}

func TestComponentAttributeService_FindComponentsByAttributeAndType(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component1, err := CreateComponent(
		testutil.AddRandSuffix("Component 1"),
		testutil.AddRandSuffix("First component"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)
	component2, err := CreateComponent(
		testutil.AddRandSuffix("Component 2"),
		testutil.AddRandSuffix("Second component"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("FindComponentsByAttributeAndType", func(t *testing.T) {
		attrPort := testutil.AddRandSuffix("port")
		valuePort := "8080"
		// Add same value but different types
		_, err := CreateComponentAttribute(component1.ID, attrPort, models.ComponentAttributeTypeString, valuePort)
		require.NoError(t, err)
		_, err = CreateComponentAttribute(component2.ID, attrPort, models.ComponentAttributeTypeNumber, valuePort)
		require.NoError(t, err)

		// Find only number type
		components, err := FindComponentsByAttributeAndType(
			attrPort,
			valuePort,
			models.ComponentAttributeTypeNumber,
		)

		require.NoError(t, err)
		assert.Len(t, components, 1)
		assert.Equal(t, component2.ID, components[0].ID)
	})
}

func TestComponentAttributeService_ComponentHasAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("ComponentHasExistingAttribute", func(t *testing.T) {
		attrEnvironment := testutil.AddRandSuffix("environment")
		_, err := CreateComponentAttribute(component.ID, attrEnvironment, models.ComponentAttributeTypeString, "production")
		require.NoError(t, err)

		hasAttribute, err := ComponentHasAttribute(component.ID, attrEnvironment)

		require.NoError(t, err)
		assert.True(t, hasAttribute)
	})

	t.Run("ComponentDoesNotHaveAttribute", func(t *testing.T) {
		attrNonExistent := testutil.AddRandSuffix("non_existent")
		hasAttribute, err := ComponentHasAttribute(component.ID, attrNonExistent)

		require.NoError(t, err)
		assert.False(t, hasAttribute)
	})
}

func TestComponentAttributeService_ComponentHasAttributeWithValue(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("ComponentHasAttributeWithCorrectValue", func(t *testing.T) {
		attrEnvironment := testutil.AddRandSuffix("environment")
		_, err := CreateComponentAttribute(component.ID, attrEnvironment, models.ComponentAttributeTypeString, "production")
		require.NoError(t, err)

		hasAttribute, err := ComponentHasAttributeWithValue(component.ID, attrEnvironment, "production")

		require.NoError(t, err)
		assert.True(t, hasAttribute)
	})

	t.Run("ComponentHasAttributeWithDifferentValue", func(t *testing.T) {
		attrEnvironment := testutil.AddRandSuffix("environment")
		hasAttribute, err := ComponentHasAttributeWithValue(component.ID, attrEnvironment, "staging")

		require.NoError(t, err)
		assert.False(t, hasAttribute)
	})
}

func TestComponentAttributeService_SetComponentAttribute(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("SetAttributeOnNewComponent", func(t *testing.T) {
		attrEnvironment := testutil.AddRandSuffix("environment")
		attribute, err := SetComponentAttribute(component.ID, attrEnvironment, models.ComponentAttributeTypeString, "production")

		require.NoError(t, err)
		assert.NotNil(t, attribute)
		assert.Equal(t, attribute.Name, attrEnvironment)
		assert.Equal(t, attribute.Value, "production")
	})

	t.Run("SetAttributeOnExistingAttribute", func(t *testing.T) {
		attrVersion := testutil.AddRandSuffix("version")
		// First call creates the attribute
		_, err := SetComponentAttribute(component.ID, attrVersion, models.ComponentAttributeTypeString, "1.0")
		require.NoError(t, err)

		// Second call should update it
		updatedAttribute, err := SetComponentAttribute(component.ID, attrVersion, models.ComponentAttributeTypeString, "2.0")

		require.NoError(t, err)
		assert.NotNil(t, updatedAttribute)
		assert.Equal(t, updatedAttribute.Name, attrVersion)
		assert.Equal(t, updatedAttribute.Value, "2.0")

		// Verify only one attribute exists
		attributes, err := GetComponentAttributes(component.ID)
		require.NoError(t, err)

		versionCount := 0
		for _, attr := range attributes {
			if attr.Name == attrVersion {
				versionCount++
				assert.Equal(t, "2.0", attr.Value)
			}
		}
		assert.Equal(t, 1, versionCount, "Should have exactly one version attribute")
	})
}

func TestComponentAttributeService_DeleteComponentAttributeByName(t *testing.T) {
	cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	component, err := CreateComponent(
		testutil.AddRandSuffix("Test Component"),
		testutil.AddRandSuffix("Test Description"),
		models.ComponentTypeInstance,
	)
	require.NoError(t, err)

	t.Run("DeleteExistingAttributeByName", func(t *testing.T) {
		tempAttr := testutil.AddRandSuffix("temp_attr")
		tempValue := testutil.AddRandSuffix("temp_value")
		_, err := CreateComponentAttribute(component.ID, tempAttr, models.ComponentAttributeTypeString, tempValue)
		require.NoError(t, err)

		err = DeleteComponentAttributeByName(component.ID, tempAttr)
		require.NoError(t, err)

		// Verify it's deleted
		hasAttribute, err := ComponentHasAttribute(component.ID, tempAttr)
		require.NoError(t, err)
		assert.False(t, hasAttribute)
	})

	t.Run("DeleteNonExistentAttributeByName", func(t *testing.T) {
		err := DeleteComponentAttributeByName(component.ID, testutil.AddRandSuffix("non_existent"))

		// Should not return error for non-existent records
		assert.NoError(t, err)
	})
}

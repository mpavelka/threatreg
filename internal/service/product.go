package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getProductRepository() (*models.ProductRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewProductRepository(db), nil
}

// GetProduct retrieves a product by its unique identifier.
// Returns the product if found, or an error if the product does not exist or database access fails.
func GetProduct(id uuid.UUID) (*models.Product, error) {
	productRepository, err := getProductRepository()
	if err != nil {
		return nil, err
	}

	return productRepository.GetByID(nil, id)
}

// CreateProduct creates a new product with the specified name and description.
// Returns the created product with its assigned ID, or an error if creation fails.
func CreateProduct(
	Name string,
	Description string,
) (*models.Product, error) {

	product := &models.Product{
		Name:        Name,
		Description: Description,
	}
	productRepository, err := getProductRepository()
	if err != nil {
		return nil, err
	}

	err = productRepository.Create(nil, product)
	if err != nil {
		fmt.Println("Error creating product:", err)
		return nil, err
	}

	return product, nil
}

// UpdateProduct updates an existing product's name and/or description within a transaction.
// Only non-nil fields are updated. Returns the updated product or an error if the update fails.
func UpdateProduct(
	id uuid.UUID,
	Name *string,
	Description *string,
) (*models.Product, error) {
	var updatedProduct *models.Product
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		productRepository, err := getProductRepository()
		if err != nil {
			return err
		}
		product, err := productRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// New values
		if Name != nil {
			product.Name = *Name
		}
		if Description != nil {
			product.Description = *Description
		}

		err = productRepository.Update(tx, product)
		if err != nil {
			return err
		}
		updatedProduct = product
		return nil
	})

	return updatedProduct, err
}

// DeleteProduct removes a product from the system by its unique identifier.
// Returns an error if the product does not exist or if deletion fails.
func DeleteProduct(id uuid.UUID) error {
	productRepository, err := getProductRepository()
	if err != nil {
		return err
	}

	return productRepository.Delete(nil, id)
}

// ListProducts retrieves all products in the system.
// Returns a slice of products or an error if database access fails.
func ListProducts() ([]models.Product, error) {
	productRepository, err := getProductRepository()
	if err != nil {
		return nil, err
	}

	return productRepository.List(nil)
}

func getThreatAssignmentRepository() (*models.ThreatAssignmentRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatAssignmentRepository(db), nil
}

// AssignThreatToProduct creates a threat assignment linking a threat to a specific product.
// Returns the created threat assignment or an error if the assignment already exists or creation fails.
func AssignThreatToProduct(productID, threatID uuid.UUID) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.AssignThreatToProduct(nil, threatID, productID)
}

// ListThreatAssignmentsByProductID retrieves all threat assignments for a specific product.
// Returns a slice of threat assignments with threat details or an error if database access fails.
func ListThreatAssignmentsByProductID(productID uuid.UUID) ([]models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.ListByProductID(nil, productID)
}

// ListThreatAssignmentsByProductIDWithResolutionByInstanceID retrieves all threat assignments for a product with resolution information filtered to a specific instance.
// Returns threat assignments for the product, but only shows resolution status if it matches the given resolutionInstanceID.
func ListThreatAssignmentsByProductIDWithResolutionByInstanceID(productID, resolutionInstanceID uuid.UUID) ([]models.ThreatAssignmentWithResolution, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.ListWithResolutionByProductID(nil, productID, resolutionInstanceID)
}

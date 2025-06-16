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

func GetProduct(id uuid.UUID) (*models.Product, error) {
	productRepository, err := getProductRepository()
	if err != nil {
		return nil, err
	}

	return productRepository.GetByID(nil, id)
}

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

func DeleteProduct(id uuid.UUID) error {
	productRepository, err := getProductRepository()
	if err != nil {
		return err
	}

	return productRepository.Delete(nil, id)
}

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

func AssignThreatToProduct(productID, threatID uuid.UUID) (*models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.AssignThreatToProduct(nil, threatID, productID)
}

func ListThreatAssignmentsByProductID(productID uuid.UUID) ([]models.ThreatAssignment, error) {
	threatAssignmentRepository, err := getThreatAssignmentRepository()
	if err != nil {
		return nil, err
	}

	return threatAssignmentRepository.ListByProductID(nil, productID)
}

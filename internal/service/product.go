package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"gorm.io/gorm"
)

func GetProduct(id string) (*models.Product, error) {
	product := &models.Product{}
	result := database.GetDB().First(product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found with id: %s", id)
		}
		return nil, fmt.Errorf("error retrieving product: %w", result.Error)
	}
	return product, nil
}

func CreateProduct(
	Name string,
	Description string,
) (*models.Product, error) {

	product := &models.Product{
		Name:        Name,
		Description: Description,
	}
	fmt.Printf("Creating product with name: '%s' and description: '%s'\n", product.Name, product.Description)
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, err
	}
	result := db.Create(product)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating product: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("unexpected number of rows affected: %d", result.RowsAffected)
	}
	return product, nil
}

func UpdateProduct(
	id string,
	Name *string,
	Description *string,
) (*models.Product, error) {
	var updatedProduct *models.Product
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		product := &models.Product{}
		result := tx.First(product, id)
		if result.Error != nil {
			return fmt.Errorf("error finding product: %w", result.Error)
		}
		if Name != nil {
			product.Name = *Name
		}
		if Description != nil {
			product.Description = *Description
		}
		saveResult := tx.Save(product)
		if saveResult.Error != nil {
			return fmt.Errorf("error updating product: %w", saveResult.Error)
		}
		updatedProduct = product
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedProduct, nil
}

func DeleteProduct(id string) error {
	result := database.GetDB().Delete(&models.Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting product: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected: %d", result.RowsAffected)
	}
	return nil
}

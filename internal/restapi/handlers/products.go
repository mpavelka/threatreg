package handlers

import (
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateProductRequest represents the request payload for creating a product
type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateProductRequest represents the request payload for updating a product
type UpdateProductRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListProducts handles GET /api/v1/products
func ListProducts(c *gin.Context) {
	products, err := service.ListProducts()
	if err != nil {
		InternalError(c, err, "Failed to retrieve products")
		return
	}

	ListResponse(c, products)
}

// CreateProduct handles POST /api/v1/products
func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	product, err := service.CreateProduct(req.Name, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to create product")
		return
	}

	CreatedResponse(c, product, "Product")
}

// GetProduct handles GET /api/v1/products/:id
func GetProduct(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	product, err := service.GetProduct(id)
	if err != nil {
		NotFoundError(c, err, "Product")
		return
	}

	GetResponse(c, product)
}

// UpdateProduct handles PUT /api/v1/products/:id
func UpdateProduct(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	product, err := service.UpdateProduct(id, req.Name, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to update product")
		return
	}

	UpdatedResponse(c, product, "Product")
}

// DeleteProduct handles DELETE /api/v1/products/:id
func DeleteProduct(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	if err := service.DeleteProduct(id); err != nil {
		InternalError(c, err, "Failed to delete product")
		return
	}

	DeletedResponse(c, "Product")
}
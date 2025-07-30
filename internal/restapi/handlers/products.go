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
// @Summary List all products
// @Description Get a list of all products in the system
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Product}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /products [get]
func ListProducts(c *gin.Context) {
	products, err := service.ListProducts()
	if err != nil {
		InternalError(c, err, "Failed to retrieve products")
		return
	}

	ListResponse(c, products)
}

// CreateProduct handles POST /api/v1/products
// @Summary Create a new product
// @Description Create a new product with the provided name and description
// @Tags Products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.Product}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /products [post]
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
// @Summary Get a product by ID
// @Description Get a specific product by its UUID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=models.Product}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /products/{id} [get]
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
// @Summary Update a product
// @Description Update a product's name and/or description
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param product body UpdateProductRequest true "Product update request"
// @Success 200 {object} handlers.SuccessResponse{data=models.Product}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /products/{id} [put]
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
// @Summary Delete a product
// @Description Delete a product by its UUID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /products/{id} [delete]
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
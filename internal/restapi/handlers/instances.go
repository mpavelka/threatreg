package handlers

import (
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateInstanceRequest represents the request payload for creating an instance
type CreateInstanceRequest struct {
	Name       string    `json:"name" binding:"required"`
	InstanceOf uuid.UUID `json:"instance_of" binding:"required"`
}

// UpdateInstanceRequest represents the request payload for updating an instance
type UpdateInstanceRequest struct {
	Name       *string    `json:"name,omitempty"`
	InstanceOf *uuid.UUID `json:"instance_of,omitempty"`
}

// ListInstances handles GET /api/v1/instances
// @Summary List all instances
// @Description Get a list of all product instances in the system
// @Tags Instances
// @Accept json
// @Produce json
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Instance}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances [get]
func ListInstances(c *gin.Context) {
	instances, err := service.ListInstances()
	if err != nil {
		InternalError(c, err, "Failed to retrieve instances")
		return
	}

	ListResponse(c, instances)
}

// CreateInstance handles POST /api/v1/instances
// @Summary Create a new instance
// @Description Create a new product instance with the provided name and product reference
// @Tags Instances
// @Accept json
// @Produce json
// @Param instance body CreateInstanceRequest true "Instance creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.Instance}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances [post]
func CreateInstance(c *gin.Context) {
	var req CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	instance, err := service.CreateInstance(req.Name, req.InstanceOf)
	if err != nil {
		InternalError(c, err, "Failed to create instance")
		return
	}

	CreatedResponse(c, instance, "Instance")
}

// GetInstance handles GET /api/v1/instances/:id
// @Summary Get an instance by ID
// @Description Get a specific product instance by its UUID
// @Tags Instances
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=models.Instance}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /instances/{id} [get]
func GetInstance(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	instance, err := service.GetInstance(id)
	if err != nil {
		NotFoundError(c, err, "Instance")
		return
	}

	GetResponse(c, instance)
}

// UpdateInstance handles PUT /api/v1/instances/:id
// @Summary Update an instance
// @Description Update an instance's name and/or product reference
// @Tags Instances
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Param instance body UpdateInstanceRequest true "Instance update request"
// @Success 200 {object} handlers.SuccessResponse{data=models.Instance}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances/{id} [put]
func UpdateInstance(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req UpdateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	instance, err := service.UpdateInstance(id, req.Name, req.InstanceOf)
	if err != nil {
		InternalError(c, err, "Failed to update instance")
		return
	}

	UpdatedResponse(c, instance, "Instance")
}

// DeleteInstance handles DELETE /api/v1/instances/:id
// @Summary Delete an instance
// @Description Delete a product instance by its UUID
// @Tags Instances
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances/{id} [delete]
func DeleteInstance(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	if err := service.DeleteInstance(id); err != nil {
		InternalError(c, err, "Failed to delete instance")
		return
	}

	DeletedResponse(c, "Instance")
}
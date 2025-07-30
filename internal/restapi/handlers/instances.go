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
func ListInstances(c *gin.Context) {
	instances, err := service.ListInstances()
	if err != nil {
		InternalError(c, err, "Failed to retrieve instances")
		return
	}

	ListResponse(c, instances)
}

// CreateInstance handles POST /api/v1/instances
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
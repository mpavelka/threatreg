package handlers

import (
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateControlRequest represents the request payload for creating a control
type CreateControlRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateControlRequest represents the request payload for updating a control
type UpdateControlRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListControls handles GET /api/v1/controls
func ListControls(c *gin.Context) {
	controls, err := service.ListControls()
	if err != nil {
		InternalError(c, err, "Failed to retrieve controls")
		return
	}

	ListResponse(c, controls)
}

// CreateControl handles POST /api/v1/controls
func CreateControl(c *gin.Context) {
	var req CreateControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	control, err := service.CreateControl(req.Title, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to create control")
		return
	}

	CreatedResponse(c, control, "Control")
}

// GetControl handles GET /api/v1/controls/:id
func GetControl(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	control, err := service.GetControl(id)
	if err != nil {
		NotFoundError(c, err, "Control")
		return
	}

	GetResponse(c, control)
}

// UpdateControl handles PUT /api/v1/controls/:id
func UpdateControl(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req UpdateControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	control, err := service.UpdateControl(id, req.Title, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to update control")
		return
	}

	UpdatedResponse(c, control, "Control")
}

// DeleteControl handles DELETE /api/v1/controls/:id
func DeleteControl(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	if err := service.DeleteControl(id); err != nil {
		InternalError(c, err, "Failed to delete control")
		return
	}

	DeletedResponse(c, "Control")
}
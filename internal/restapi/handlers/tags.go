package handlers

import (
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateTagRequest represents the request payload for creating a tag
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Color       string `json:"color" binding:"required"`
}

// UpdateTagRequest represents the request payload for updating a tag
type UpdateTagRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// ListTags handles GET /api/v1/tags
func ListTags(c *gin.Context) {
	tags, err := service.ListTags()
	if err != nil {
		InternalError(c, err, "Failed to retrieve tags")
		return
	}

	ListResponse(c, tags)
}

// CreateTag handles POST /api/v1/tags
func CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	tag, err := service.CreateTag(req.Name, req.Description, req.Color)
	if err != nil {
		InternalError(c, err, "Failed to create tag")
		return
	}

	CreatedResponse(c, tag, "Tag")
}

// GetTag handles GET /api/v1/tags/:id
func GetTag(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	tag, err := service.GetTag(id)
	if err != nil {
		NotFoundError(c, err, "Tag")
		return
	}

	GetResponse(c, tag)
}

// UpdateTag handles PUT /api/v1/tags/:id
func UpdateTag(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	tag, err := service.UpdateTag(id, req.Name, req.Description, req.Color)
	if err != nil {
		InternalError(c, err, "Failed to update tag")
		return
	}

	UpdatedResponse(c, tag, "Tag")
}

// DeleteTag handles DELETE /api/v1/tags/:id
func DeleteTag(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	if err := service.DeleteTag(id); err != nil {
		InternalError(c, err, "Failed to delete tag")
		return
	}

	DeletedResponse(c, "Tag")
}
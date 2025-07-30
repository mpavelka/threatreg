package handlers

import (
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateThreatRequest represents the request payload for creating a threat
type CreateThreatRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateThreatRequest represents the request payload for updating a threat
type UpdateThreatRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListThreats handles GET /api/v1/threats
func ListThreats(c *gin.Context) {
	threats, err := service.ListThreats()
	if err != nil {
		InternalError(c, err, "Failed to retrieve threats")
		return
	}

	ListResponse(c, threats)
}

// CreateThreat handles POST /api/v1/threats
func CreateThreat(c *gin.Context) {
	var req CreateThreatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	threat, err := service.CreateThreat(req.Title, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to create threat")
		return
	}

	CreatedResponse(c, threat, "Threat")
}

// GetThreat handles GET /api/v1/threats/:id
func GetThreat(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	threat, err := service.GetThreat(id)
	if err != nil {
		NotFoundError(c, err, "Threat")
		return
	}

	GetResponse(c, threat)
}

// UpdateThreat handles PUT /api/v1/threats/:id
func UpdateThreat(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req UpdateThreatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	threat, err := service.UpdateThreat(id, req.Title, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to update threat")
		return
	}

	UpdatedResponse(c, threat, "Threat")
}

// DeleteThreat handles DELETE /api/v1/threats/:id
func DeleteThreat(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	if err := service.DeleteThreat(id); err != nil {
		InternalError(c, err, "Failed to delete threat")
		return
	}

	DeletedResponse(c, "Threat")
}
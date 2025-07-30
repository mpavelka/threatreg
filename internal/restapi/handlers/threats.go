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
// @Summary List all threats
// @Description Get a list of all security threats in the system
// @Tags Threats
// @Accept json
// @Produce json
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Threat}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threats [get]
func ListThreats(c *gin.Context) {
	threats, err := service.ListThreats()
	if err != nil {
		InternalError(c, err, "Failed to retrieve threats")
		return
	}

	ListResponse(c, threats)
}

// CreateThreat handles POST /api/v1/threats
// @Summary Create a new threat
// @Description Create a new security threat with the provided title and description
// @Tags Threats
// @Accept json
// @Produce json
// @Param threat body CreateThreatRequest true "Threat creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.Threat}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threats [post]
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
// @Summary Get a threat by ID
// @Description Get a specific security threat by its UUID
// @Tags Threats
// @Accept json
// @Produce json
// @Param id path string true "Threat ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=models.Threat}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /threats/{id} [get]
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
// @Summary Update a threat
// @Description Update a threat's title and/or description
// @Tags Threats
// @Accept json
// @Produce json
// @Param id path string true "Threat ID (UUID)"
// @Param threat body UpdateThreatRequest true "Threat update request"
// @Success 200 {object} handlers.SuccessResponse{data=models.Threat}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threats/{id} [put]
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
// @Summary Delete a threat
// @Description Delete a security threat by its UUID
// @Tags Threats
// @Accept json
// @Produce json
// @Param id path string true "Threat ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threats/{id} [delete]
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
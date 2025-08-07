package handlers

import (
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetThreatAssignment handles GET /api/v1/threat-assignments/:id
// @Summary Get a threat assignment
// @Description Retrieve a threat assignment by its ID
// @Tags Threat Assignments
// @Accept json
// @Produce json
// @Param id path string true "Threat Assignment ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=models.ThreatAssignment}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-assignments/{id} [get]
func GetThreatAssignment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ValidationError(c, err)
		return
	}

	assignment, err := service.GetThreatAssignmentById(id)
	if err != nil {
		NotFoundError(c, err, "Threat assignment")
		return
	}

	GetResponse(c, assignment)
}

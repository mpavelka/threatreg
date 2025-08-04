package handlers

import (
	"strconv"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateThreatResolutionRequest represents the request payload for creating a threat resolution
type CreateThreatResolutionRequest struct {
	ThreatAssignmentID int                                     `json:"threatAssignmentId" binding:"required"`
	ComponentID        uuid.UUID                               `json:"instanceId"`
	Status             models.ThreatAssignmentResolutionStatus `json:"status" binding:"required"`
	Description        string                                  `json:"description" binding:"required"`
}

// UpdateThreatResolutionRequest represents the request payload for updating a threat resolution
type UpdateThreatResolutionRequest struct {
	Status      *models.ThreatAssignmentResolutionStatus `json:"status,omitempty"`
	Description *string                                  `json:"description,omitempty"`
}

// DelegateResolutionRequest represents the request payload for delegating a resolution
type DelegateResolutionRequest struct {
	TargetResolutionID uuid.UUID `json:"targetResolutionId" binding:"required"`
}

// CreateThreatResolution handles POST /api/v1/threat-resolutions
// @Summary Create a new threat resolution
// @Description This endpoint creates a new threat resolution.
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param resolution body CreateThreatResolutionRequest true "Threat resolution creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolution} "Returns created threat resolution"
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions [post]
func CreateThreatResolution(c *gin.Context) {
	var req CreateThreatResolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	resolution, err := service.CreateThreatResolution(
		req.ThreatAssignmentID,
		req.ComponentID,
		req.Status,
		req.Description,
	)
	if err != nil {
		InternalError(c, err, "Failed to create threat resolution")
		return
	}

	CreatedResponse(c, resolution, "Threat resolution (DUMMY - endpoint disabled)")
	return
}

// UpdateThreatResolution handles PUT /api/v1/threat-resolutions/:id
// @Summary Update a threat resolution
// @Description Update an existing threat resolution's status and/or description
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param id path string true "Resolution ID"
// @Param resolution body UpdateThreatResolutionRequest true "Threat resolution update request"
// @Success 200 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/{id} [put]
func UpdateThreatResolution(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req UpdateThreatResolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	resolution, err := service.UpdateThreatResolution(id, req.Status, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to update threat resolution")
		return
	}

	UpdatedResponse(c, resolution, "Threat resolution")
}

// GetThreatResolution handles GET /api/v1/threat-resolutions/:id
// @Summary Get a threat resolution
// @Description Retrieve a threat resolution by its ID
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param id path string true "Resolution ID"
// @Success 200 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/{id} [get]
func GetThreatResolution(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolution, err := service.GetThreatResolution(id)
	if err != nil {
		NotFoundError(c, err, "Threat resolution")
		return
	}

	GetResponse(c, resolution)
}

// GetThreatResolutionByThreatAssignmentID handles GET /api/v1/threat-resolutions/by-assignment/:assignmentId
// @Summary Get threat resolution by assignment ID
// @Description Retrieve a threat resolution by its threat assignment ID
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param assignmentId path int true "Threat Assignment ID"
// @Success 200 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/by-assignment/{assignmentId} [get]
func GetThreatResolutionByThreatAssignmentID(c *gin.Context) {
	assignmentIDStr := c.Param("assignmentId")
	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolution, err := service.GetThreatResolutionByThreatAssignmentID(assignmentID)
	if err != nil {
		NotFoundError(c, err, "Threat resolution")
		return
	}

	GetResponse(c, resolution)
}

// DeleteThreatResolution handles DELETE /api/v1/threat-resolutions/:id
// @Summary Delete a threat resolution
// @Description Delete a threat resolution and update upstream delegation statuses
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param id path string true "Resolution ID"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/{id} [delete]
func DeleteThreatResolution(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	err = service.DeleteThreatResolution(id)
	if err != nil {
		InternalError(c, err, "Failed to delete threat resolution")
		return
	}

	DeletedResponse(c, "Threat resolution")
}

// DelegateResolution handles POST /api/v1/threat-resolutions/:id/delegate
// @Summary Delegate a threat resolution
// @Description Create a delegation from one threat resolution to another
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param id path string true "Source Resolution ID"
// @Param delegation body DelegateResolutionRequest true "Delegation request"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/{id}/delegate [post]
func DelegateResolution(c *gin.Context) {
	sourceID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req DelegateResolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	// Get both resolutions
	sourceResolution, err := service.GetThreatResolution(sourceID)
	if err != nil {
		NotFoundError(c, err, "Source threat resolution")
		return
	}

	targetResolution, err := service.GetThreatResolution(req.TargetResolutionID)
	if err != nil {
		NotFoundError(c, err, "Target threat resolution")
		return
	}

	err = service.DelegateResolution(*sourceResolution, *targetResolution)
	if err != nil {
		InternalError(c, err, "Failed to delegate threat resolution")
		return
	}

	HandleSuccess(c, 200, nil, "Threat resolution delegated successfully")
}

// GetDelegatedToResolutionByDelegatedByID handles GET /api/v1/threat-resolutions/:id/delegated-to
// @Summary Get delegated-to resolution
// @Description Retrieve the target resolution of a delegation by the source resolution ID
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param id path string true "Source Resolution ID"
// @Success 200 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/{id}/delegated-to [get]
func GetDelegatedToResolutionByDelegatedByID(c *gin.Context) {
	delegatedByID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolution, err := service.GetDelegatedToResolutionByDelegatedByID(delegatedByID)
	if err != nil {
		NotFoundError(c, err, "Delegated threat resolution")
		return
	}

	if resolution == nil {
		NotFoundError(c, err, "No delegation found for this resolution")
		return
	}

	GetResponse(c, resolution)
}

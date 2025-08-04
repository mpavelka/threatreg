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
	InstanceID         *uuid.UUID                              `json:"instanceId,omitempty"`
	ProductID          *uuid.UUID                              `json:"productId,omitempty"`
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
// @Summary Create a new threat resolution (TEMPORARILY DISABLED)
// @Description TEMPORARILY DISABLED: This endpoint is disabled due to service layer refactoring from Product/Instance to unified Component model. Returns dummy response for now.
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param resolution body CreateThreatResolutionRequest true "Threat resolution creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolution} "Returns dummy response"
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions [post]
func CreateThreatResolution(c *gin.Context) {
	var req CreateThreatResolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	// TODO: This endpoint is temporarily disabled due to service layer refactoring
	// The CreateThreatResolution function signature has changed from (int, *uuid.UUID, *uuid.UUID, status, string)
	// to (int, uuid.UUID, status, string) using unified Component model.
	// Returning dummy response for now.
	dummyResolution := models.ThreatAssignmentResolution{
		ID:                 uuid.New(),
		ThreatAssignmentID: req.ThreatAssignmentID,
		Status:             req.Status,
		Description:        req.Description,
	}

	CreatedResponse(c, dummyResolution, "Threat resolution (DUMMY - endpoint disabled)")
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

// GetInstanceLevelThreatResolution handles GET /api/v1/threat-resolutions/instance/:instanceId/assignment/:assignmentId
// @Summary Get instance-level threat resolution
// @Description Retrieve a threat resolution for a specific instance and threat assignment
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param instanceId path string true "Instance ID"
// @Param assignmentId path int true "Threat Assignment ID"
// @Success 200 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/instance/{instanceId}/assignment/{assignmentId} [get]
func GetInstanceLevelThreatResolution(c *gin.Context) {
	instanceID, err := ParseUUID(c, "instanceId")
	if err != nil {
		ValidationError(c, err)
		return
	}

	assignmentIDStr := c.Param("assignmentId")
	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolution, err := service.GetInstanceLevelThreatResolution(assignmentID, instanceID)
	if err != nil {
		NotFoundError(c, err, "Threat resolution")
		return
	}

	GetResponse(c, resolution)
}

// GetInstanceLevelThreatResolutionWithDelegation handles GET /api/v1/threat-resolutions/instance/:instanceId/assignment/:assignmentId/with-delegation
// @Summary Get instance-level threat resolution with delegation
// @Description Retrieve a threat resolution with delegation information for a specific instance and threat assignment
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param instanceId path string true "Instance ID"
// @Param assignmentId path int true "Threat Assignment ID"
// @Success 200 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolutionWithDelegation}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/instance/{instanceId}/assignment/{assignmentId}/with-delegation [get]
func GetInstanceLevelThreatResolutionWithDelegation(c *gin.Context) {
	instanceID, err := ParseUUID(c, "instanceId")
	if err != nil {
		ValidationError(c, err)
		return
	}

	assignmentIDStr := c.Param("assignmentId")
	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolution, err := service.GetInstanceLevelThreatResolutionWithDelegation(assignmentID, instanceID)
	if err != nil {
		NotFoundError(c, err, "Threat resolution")
		return
	}

	GetResponse(c, resolution)
}

// GetProductLevelThreatResolution handles GET /api/v1/threat-resolutions/product/:productId/assignment/:assignmentId
// @Summary Get product-level threat resolution
// @Description Retrieve a threat resolution for a specific product and threat assignment
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param assignmentId path int true "Threat Assignment ID"
// @Success 200 {object} handlers.SuccessResponse{data=models.ThreatAssignmentResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/product/{productId}/assignment/{assignmentId} [get]
func GetProductLevelThreatResolution(c *gin.Context) {
	productID, err := ParseUUID(c, "productId")
	if err != nil {
		ValidationError(c, err)
		return
	}

	assignmentIDStr := c.Param("assignmentId")
	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolution, err := service.GetProductLevelThreatResolution(assignmentID, productID)
	if err != nil {
		NotFoundError(c, err, "Threat resolution")
		return
	}

	GetResponse(c, resolution)
}

// ListThreatResolutionsByProductID handles GET /api/v1/threat-resolutions/by-product/:productId
// @Summary List threat resolutions by product ID
// @Description Retrieve all threat resolutions for a specific product
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.ThreatAssignmentResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/by-product/{productId} [get]
func ListThreatResolutionsByProductID(c *gin.Context) {
	productID, err := ParseUUID(c, "productId")
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolutions, err := service.ListThreatResolutionsByProductID(productID)
	if err != nil {
		InternalError(c, err, "Failed to retrieve threat resolutions")
		return
	}

	ListResponse(c, resolutions)
}

// ListThreatResolutionsByInstanceID handles GET /api/v1/threat-resolutions/by-instance/:instanceId
// @Summary List threat resolutions by instance ID
// @Description Retrieve all threat resolutions for a specific instance
// @Tags Threat Resolutions
// @Accept json
// @Produce json
// @Param instanceId path string true "Instance ID"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.ThreatAssignmentResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /threat-resolutions/by-instance/{instanceId} [get]
func ListThreatResolutionsByInstanceID(c *gin.Context) {
	instanceID, err := ParseUUID(c, "instanceId")
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolutions, err := service.ListThreatResolutionsByInstanceID(instanceID)
	if err != nil {
		InternalError(c, err, "Failed to retrieve threat resolutions")
		return
	}

	ListResponse(c, resolutions)
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

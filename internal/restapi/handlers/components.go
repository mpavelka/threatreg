package handlers

import (
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateComponentRequest represents the request payload for creating a component
type CreateComponentRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Description string                   `json:"description" binding:"required"`
	Type        models.ComponentType     `json:"type" binding:"required"`
}

// UpdateComponentRequest represents the request payload for updating a component
type UpdateComponentRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListComponents handles GET /api/v1/components
// @Summary List all components
// @Description Get a list of all components in the system
// @Tags Components
// @Accept json
// @Produce json
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Component}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components [get]
func ListComponents(c *gin.Context) {
	components, err := service.ListComponents()
	if err != nil {
		InternalError(c, err, "Failed to retrieve components")
		return
	}

	ListResponse(c, components)
}

// CreateComponent handles POST /api/v1/components
// @Summary Create a new component
// @Description Create a new component with the provided name, description, and type
// @Tags Components
// @Accept json
// @Produce json
// @Param component body CreateComponentRequest true "Component creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.Component}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components [post]
func CreateComponent(c *gin.Context) {
	var req CreateComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	component, err := service.CreateComponent(req.Name, req.Description, req.Type)
	if err != nil {
		InternalError(c, err, "Failed to create component")
		return
	}

	CreatedResponse(c, component, "Component")
}

// GetComponent handles GET /api/v1/components/:id
// @Summary Get a component by ID
// @Description Get a specific component by its UUID
// @Tags Components
// @Accept json
// @Produce json
// @Param id path string true "Component ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=models.Component}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /components/{id} [get]
func GetComponent(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	component, err := service.GetComponent(id)
	if err != nil {
		NotFoundError(c, err, "Component")
		return
	}

	GetResponse(c, component)
}

// UpdateComponent handles PUT /api/v1/components/:id
// @Summary Update a component
// @Description Update a component's name and/or description
// @Tags Components
// @Accept json
// @Produce json
// @Param id path string true "Component ID (UUID)"
// @Param component body UpdateComponentRequest true "Component update request"
// @Success 200 {object} handlers.SuccessResponse{data=models.Component}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components/{id} [put]
func UpdateComponent(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req UpdateComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	component, err := service.UpdateComponent(id, req.Name, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to update component")
		return
	}

	UpdatedResponse(c, component, "Component")
}

// DeleteComponent handles DELETE /api/v1/components/:id
// @Summary Delete a component
// @Description Delete a component by its UUID
// @Tags Components
// @Accept json
// @Produce json
// @Param id path string true "Component ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components/{id} [delete]
func DeleteComponent(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	if err := service.DeleteComponent(id); err != nil {
		InternalError(c, err, "Failed to delete component")
		return
	}

	DeletedResponse(c, "Component")
}

// ListComponentsByType handles GET /api/v1/components/by-type/:type
// @Summary List components by type
// @Description Get all components of a specific type
// @Tags Components
// @Accept json
// @Produce json
// @Param type path string true "Component Type (product or instance)"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Component}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components/by-type/{type} [get]
func ListComponentsByType(c *gin.Context) {
	typeStr := c.Param("type")
	componentType := models.ComponentType(typeStr)

	// Validate component type
	if componentType != models.ComponentTypeProduct && componentType != models.ComponentTypeInstance {
		ValidationError(c, gin.Error{})
		return
	}

	components, err := service.ListComponentsByType(componentType)
	if err != nil {
		InternalError(c, err, "Failed to retrieve components by type")
		return
	}

	ListResponse(c, components)
}

// FilterComponents handles GET /api/v1/components/filter
// @Summary Filter components by name
// @Description Search for components by name using case-insensitive partial matching
// @Tags Components
// @Accept json
// @Produce json
// @Param component_name query string false "Component name filter (partial match)"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Component}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components/filter [get]
func FilterComponents(c *gin.Context) {
	componentName := c.Query("component_name")

	components, err := service.FilterComponents(componentName)
	if err != nil {
		InternalError(c, err, "Failed to filter components")
		return
	}

	ListResponse(c, components)
}

// AssignThreatToComponentRequest represents the request payload for assigning a threat to a component
type AssignThreatToComponentRequest struct {
	ThreatID uuid.UUID `json:"threat_id" binding:"required"`
}

// AssignThreatToComponent handles POST /api/v1/components/:id/threats
// @Summary Assign a threat to a component
// @Description Create a threat assignment linking a threat to a specific component
// @Tags Components
// @Accept json
// @Produce json
// @Param id path string true "Component ID (UUID)"
// @Param threat body AssignThreatToComponentRequest true "Threat assignment request"
// @Success 201 {object} handlers.SuccessResponse{data=models.ThreatAssignment}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components/{id}/threats [post]
func AssignThreatToComponent(c *gin.Context) {
	componentID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req AssignThreatToComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	threatAssignment, err := service.AssignThreatToComponent(componentID, req.ThreatID)
	if err != nil {
		InternalError(c, err, "Failed to assign threat to component")
		return
	}

	CreatedResponse(c, threatAssignment, "Threat assignment")
}

// ListThreatAssignmentsByComponent handles GET /api/v1/components/:id/threats
// @Summary List threat assignments for a component
// @Description Get all threat assignments for a specific component
// @Tags Components
// @Accept json
// @Produce json
// @Param id path string true "Component ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.ThreatAssignment}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components/{id}/threats [get]
func ListThreatAssignmentsByComponent(c *gin.Context) {
	componentID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	threatAssignments, err := service.ListThreatAssignmentsByComponentID(componentID)
	if err != nil {
		InternalError(c, err, "Failed to retrieve threat assignments for component")
		return
	}

	ListResponse(c, threatAssignments)
}

// ListThreatAssignmentsByComponentWithResolutionByComponent handles GET /api/v1/components/:id/threats/with-resolution/:resolutionComponentId
// @Summary List threat assignments with resolution status for a component filtered by resolution component
// @Description Get all threat assignments for a specific component, showing resolution status only for the specified resolution component
// @Tags Components
// @Accept json
// @Produce json
// @Param id path string true "Component ID (UUID)"
// @Param resolutionComponentId path string true "Resolution Component ID (UUID) to filter resolutions"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.ThreatAssignmentWithResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /components/{id}/threats/with-resolution/{resolutionComponentId} [get]
func ListThreatAssignmentsByComponentWithResolutionByComponent(c *gin.Context) {
	componentID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolutionComponentID, err := ParseUUID(c, "resolutionComponentId")
	if err != nil {
		ValidationError(c, err)
		return
	}

	threatAssignments, err := service.ListThreatAssignmentsByComponentIDWithResolutionByComponentID(componentID, resolutionComponentID)
	if err != nil {
		InternalError(c, err, "Failed to retrieve threat assignments with resolution for component")
		return
	}

	ListResponse(c, threatAssignments)
}
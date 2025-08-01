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

// GetDomainsByInstance handles GET /api/v1/instances/:id/domains
// @Summary Get domains for an instance
// @Description Get all domains that contain a specific instance
// @Tags Instances
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Domain}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances/{id}/domains [get]
func GetDomainsByInstance(c *gin.Context) {
	instanceID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	domains, err := service.GetDomainsByInstance(instanceID)
	if err != nil {
		InternalError(c, err, "Failed to retrieve domains for instance")
		return
	}

	ListResponse(c, domains)
}

// ListInstancesByProduct handles GET /api/v1/instances/by-product/:productId
// @Summary List instances by product
// @Description Get all instances that belong to a specific product
// @Tags Instances
// @Accept json
// @Produce json
// @Param productId path string true "Product ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Instance}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances/by-product/{productId} [get]
func ListInstancesByProduct(c *gin.Context) {
	productID, err := ParseUUID(c, "productId")
	if err != nil {
		ValidationError(c, err)
		return
	}

	instances, err := service.ListInstancesByProductID(productID)
	if err != nil {
		InternalError(c, err, "Failed to retrieve instances for product")
		return
	}

	ListResponse(c, instances)
}

// FilterInstances handles GET /api/v1/instances/filter
// @Summary Filter instances by name and/or product name
// @Description Search for instances by name and/or product name using case-insensitive partial matching
// @Tags Instances
// @Accept json
// @Produce json
// @Param instance_name query string false "Instance name filter (partial match)"
// @Param product_name query string false "Product name filter (partial match)"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Instance}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances/filter [get]
func FilterInstances(c *gin.Context) {
	instanceName := c.Query("instance_name")
	productName := c.Query("product_name")

	instances, err := service.FilterInstances(instanceName, productName)
	if err != nil {
		InternalError(c, err, "Failed to filter instances")
		return
	}

	ListResponse(c, instances)
}

// AssignThreatToInstanceRequest represents the request payload for assigning a threat to an instance
type AssignThreatToInstanceRequest struct {
	ThreatID uuid.UUID `json:"threat_id" binding:"required"`
}

// AssignThreatToInstance handles POST /api/v1/instances/:id/threats
// @Summary Assign a threat to an instance
// @Description Create a threat assignment linking a threat to a specific instance
// @Tags Instances
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Param threat body AssignThreatToInstanceRequest true "Threat assignment request"
// @Success 201 {object} handlers.SuccessResponse{data=models.ThreatAssignment}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances/{id}/threats [post]
func AssignThreatToInstance(c *gin.Context) {
	instanceID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req AssignThreatToInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	threatAssignment, err := service.AssignThreatToInstance(instanceID, req.ThreatID)
	if err != nil {
		InternalError(c, err, "Failed to assign threat to instance")
		return
	}

	CreatedResponse(c, threatAssignment, "Threat assignment")
}

// ListThreatAssignmentsByInstance handles GET /api/v1/instances/:id/threats
// @Summary List threat assignments for an instance
// @Description Get all threat assignments for a specific instance
// @Tags Instances
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.ThreatAssignment}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances/{id}/threats [get]
func ListThreatAssignmentsByInstance(c *gin.Context) {
	instanceID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	threatAssignments, err := service.ListThreatAssignmentsByInstanceID(instanceID)
	if err != nil {
		InternalError(c, err, "Failed to retrieve threat assignments for instance")
		return
	}

	ListResponse(c, threatAssignments)
}

// ListThreatAssignmentsByInstanceWithResolutionByInstance handles GET /api/v1/instances/:id/threats/with-resolution/:resolutionInstanceId
// @Summary List threat assignments with resolution status for an instance filtered by resolution instance
// @Description Get all threat assignments for a specific instance, showing resolution status only for the specified resolution instance
// @Tags Instances
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Param resolutionInstanceId path string true "Resolution Instance ID (UUID) to filter resolutions"
// @Success 200 {object} handlers.SuccessResponse{data=[]models.ThreatAssignmentWithResolution}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /instances/{id}/threats/with-resolution/{resolutionInstanceId} [get]
func ListThreatAssignmentsByInstanceWithResolutionByInstance(c *gin.Context) {
	instanceID, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	resolutionInstanceID, err := ParseUUID(c, "resolutionInstanceId")
	if err != nil {
		ValidationError(c, err)
		return
	}

	threatAssignments, err := service.ListThreatAssignmentsByInstanceIDWithResolutionByInstanceID(instanceID, resolutionInstanceID)
	if err != nil {
		InternalError(c, err, "Failed to retrieve threat assignments with resolution for instance")
		return
	}

	ListResponse(c, threatAssignments)
}

package handlers

import (
	"threatreg/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateDomainRequest represents the request payload for creating a domain
type CreateDomainRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateDomainRequest represents the request payload for updating a domain
type UpdateDomainRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListDomains handles GET /api/v1/domains
// @Summary List all domains
// @Description Get a list of all security domains in the system
// @Tags Domains
// @Accept json
// @Produce json
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Domain}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /domains [get]
func ListDomains(c *gin.Context) {
	domains, err := service.ListDomains()
	if err != nil {
		InternalError(c, err, "Failed to retrieve domains")
		return
	}

	ListResponse(c, domains)
}

// CreateDomain handles POST /api/v1/domains
// @Summary Create a new domain
// @Description Create a new security domain with the provided name and description
// @Tags Domains
// @Accept json
// @Produce json
// @Param domain body CreateDomainRequest true "Domain creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.Domain}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /domains [post]
func CreateDomain(c *gin.Context) {
	var req CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	domain, err := service.CreateDomain(req.Name, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to create domain")
		return
	}

	CreatedResponse(c, domain, "Domain")
}

// GetDomain handles GET /api/v1/domains/:id
// @Summary Get a domain by ID
// @Description Get a specific security domain by its UUID
// @Tags Domains
// @Accept json
// @Produce json
// @Param id path string true "Domain ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=models.Domain}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /domains/{id} [get]
func GetDomain(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	domain, err := service.GetDomain(id)
	if err != nil {
		NotFoundError(c, err, "Domain")
		return
	}

	GetResponse(c, domain)
}

// UpdateDomain handles PUT /api/v1/domains/:id
// @Summary Update a domain
// @Description Update a domain's name and/or description
// @Tags Domains
// @Accept json
// @Produce json
// @Param id path string true "Domain ID (UUID)"
// @Param domain body UpdateDomainRequest true "Domain update request"
// @Success 200 {object} handlers.SuccessResponse{data=models.Domain}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /domains/{id} [put]
func UpdateDomain(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	var req UpdateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	domain, err := service.UpdateDomain(id, req.Name, req.Description)
	if err != nil {
		InternalError(c, err, "Failed to update domain")
		return
	}

	UpdatedResponse(c, domain, "Domain")
}

// DeleteDomain handles DELETE /api/v1/domains/:id
// @Summary Delete a domain
// @Description Delete a security domain by its UUID
// @Tags Domains
// @Accept json
// @Produce json
// @Param id path string true "Domain ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /domains/{id} [delete]
func DeleteDomain(c *gin.Context) {
	id, err := ParseUUID(c, "id")
	if err != nil {
		ValidationError(c, err)
		return
	}

	if err := service.DeleteDomain(id); err != nil {
		InternalError(c, err, "Failed to delete domain")
		return
	}

	DeletedResponse(c, "Domain")
}
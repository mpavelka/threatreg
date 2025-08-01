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
// @Summary List all controls
// @Description Get a list of all security controls in the system
// @Tags Controls
// @Accept json
// @Produce json
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Control}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /controls [get]
func ListControls(c *gin.Context) {
	controls, err := service.ListControls()
	if err != nil {
		InternalError(c, err, "Failed to retrieve controls")
		return
	}

	ListResponse(c, controls)
}

// CreateControl handles POST /api/v1/controls
// @Summary Create a new control
// @Description Create a new security control with the provided title and description
// @Tags Controls
// @Accept json
// @Produce json
// @Param control body CreateControlRequest true "Control creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.Control}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /controls [post]
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
// @Summary Get a control by ID
// @Description Get a specific security control by its UUID
// @Tags Controls
// @Accept json
// @Produce json
// @Param id path string true "Control ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse{data=models.Control}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /controls/{id} [get]
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
// @Summary Update a control
// @Description Update a control's title and/or description
// @Tags Controls
// @Accept json
// @Produce json
// @Param id path string true "Control ID (UUID)"
// @Param control body UpdateControlRequest true "Control update request"
// @Success 200 {object} handlers.SuccessResponse{data=models.Control}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /controls/{id} [put]
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
// @Summary Delete a control
// @Description Delete a security control by its UUID
// @Tags Controls
// @Accept json
// @Produce json
// @Param id path string true "Control ID (UUID)"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /controls/{id} [delete]
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

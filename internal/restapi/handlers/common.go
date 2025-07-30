package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// HandleError sends a standardized error response
func HandleError(c *gin.Context, statusCode int, err error, message string) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: message,
		Code:    statusCode,
	}
	c.JSON(statusCode, response)
}

// HandleSuccess sends a standardized success response
func HandleSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	response := SuccessResponse{
		Data:    data,
		Message: message,
	}
	c.JSON(statusCode, response)
}

// ParseUUID parses a UUID from a URL parameter
func ParseUUID(c *gin.Context, paramName string) (uuid.UUID, error) {
	idStr := c.Param(paramName)
	return uuid.Parse(idStr)
}

// ParseLimit parses the limit query parameter with a default and maximum value
func ParseLimit(c *gin.Context, defaultLimit, maxLimit int) int {
	limitStr := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

// ParseOffset parses the offset query parameter
func ParseOffset(c *gin.Context) int {
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return 0
	}
	return offset
}

// ValidationError creates a bad request error response for validation failures
func ValidationError(c *gin.Context, err error) {
	HandleError(c, http.StatusBadRequest, err, "Invalid request payload")
}

// NotFoundError creates a not found error response
func NotFoundError(c *gin.Context, err error, resource string) {
	HandleError(c, http.StatusNotFound, err, resource+" not found")
}

// InternalError creates an internal server error response
func InternalError(c *gin.Context, err error, message string) {
	HandleError(c, http.StatusInternalServerError, err, message)
}

// CreatedResponse creates a successful creation response
func CreatedResponse(c *gin.Context, data interface{}, resource string) {
	HandleSuccess(c, http.StatusCreated, data, resource+" created successfully")
}

// UpdatedResponse creates a successful update response
func UpdatedResponse(c *gin.Context, data interface{}, resource string) {
	HandleSuccess(c, http.StatusOK, data, resource+" updated successfully")
}

// DeletedResponse creates a successful deletion response
func DeletedResponse(c *gin.Context, resource string) {
	HandleSuccess(c, http.StatusOK, nil, resource+" deleted successfully")
}

// ListResponse creates a successful list response
func ListResponse(c *gin.Context, data interface{}) {
	HandleSuccess(c, http.StatusOK, data, "")
}

// GetResponse creates a successful get response
func GetResponse(c *gin.Context, data interface{}) {
	HandleSuccess(c, http.StatusOK, data, "")
}
package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Generic API response
type APIResponse struct {
	Status    string      `json:"status"`
	Timestamp string      `json:"timestamp"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Status:    "success",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Data:      data,
	})
}

func SuccessWithCode(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Status:    "success",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Data:      data,
	})
}

// Error sends an error response with custom code
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Status:    "error",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Data:      nil,
	})
}

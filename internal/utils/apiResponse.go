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
	Metadata  interface{} `json:"metadata"`
}

/*
API standard response to created
*/
func SuccessCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Status:    "success",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Data:      data,
	})
}

/*
API standard response to success operation
*/
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Status:    "success",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Data:      data,
	})
}

/*
API standard response to success operation
*/
func SuccessList(c *gin.Context, message string, data interface{}, metadata interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Status:    "success",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Data:      data,
		Metadata:  metadata,
	})
}

/*
API standard response to success operation, and specify the http code
*/
func SuccessWithCode(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Status:    "success",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Data:      data,
	})
}

/*
API standard response to error/failed operation, and specify the http code
*/
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Status:    "error",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Data:      nil,
	})
}

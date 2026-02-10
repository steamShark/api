package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"steamshark-api/internal/models"
)

// Success sends a successful JSON response.
func Success[T any](ctx *gin.Context, status int, data T) {
	ctx.JSON(status, models.APIResponse[T]{
		Success: true,
		Data:    data,
	})
}

// Error sends an error JSON response.
func Error(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, models.APIResponse[any]{
		Success: false,
		Error:   message,
	})
}

// NoContent sends a 204 No Content response.
func NoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

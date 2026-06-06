package middlewares

import (
	"net/http"
	"steamshark-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// AdminKeyAuth rejects requests that don't carry the correct X-Admin-Key header.
// If the server was started without ADMIN_SECRET_KEY set, all requests are rejected
// to avoid accidentally leaving the endpoint open.
func AdminKeyAuth(adminKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if adminKey == "" {
			utils.Error(ctx, http.StatusServiceUnavailable, "admin key not configured on server")
			ctx.Abort()
			return
		}
		if ctx.GetHeader("X-Admin-Key") != adminKey {
			utils.Error(ctx, http.StatusUnauthorized, "invalid or missing authentication")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

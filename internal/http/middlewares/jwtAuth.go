package middlewares

import (
	"net/http"
	"strings"
	"steamshark-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth validates the Bearer token and injects user_id, steam_id and role into the context.
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if jwtSecret == "" {
			utils.Error(ctx, http.StatusServiceUnavailable, "JWT not configured on server")
			ctx.Abort()
			return
		}

		header := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			utils.Error(ctx, http.StatusUnauthorized, "missing or invalid Authorization header")
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.Error(ctx, http.StatusUnauthorized, "invalid or expired token")
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Error(ctx, http.StatusUnauthorized, "invalid token claims")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Set("steam_id", claims["steam_id"])
		ctx.Set("role", claims["role"])
		ctx.Next()
	}
}

package middlewares

import (
	"github.com/gin-gonic/gin"
)

/*
@brief: CORS policy middleware to the response
*/
func CORSPolicy() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* Get the origin */
		origin := c.GetHeader("Origin")

		/* Map with allowed origin */
		allowedOrigins := map[string]bool{
			"http://localhost:8090":      true,
			"https://steamshark.app":     true,
			"https://www.steamshark.app": true,
		}

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			// fallback: no Origin or not allowed
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		//For now, only GET will be supported
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET") //, POST, PUT, PATCH, DELETE, OPTIONS
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-Id")
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(requestIDKey, rid)
		c.Writer.Header().Set("X-Request-Id", rid)
		c.Next()
	}
}

// Middleware function to log with the zap packages
//
// Responsible for all the loggin within the API
func ZapLogger(log *zap.Logger) gin.HandlerFunc {
	l := log.Named("http")
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		rid, _ := c.Get(requestIDKey)

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.Any("request_id", rid),
		}

		// log errors if present
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
			l.Warn("request", fields...)
			return
		}

		if status >= 500 {
			l.Error("request", fields...)
		} else if status >= 400 {
			l.Warn("request", fields...)
		} else {
			l.Info("request", fields...)
		}
	}
}

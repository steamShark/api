package controllers

import (
	"time"

	"steamshark-api/utils"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (ctrl *HealthController) Ping(c *gin.Context) {
	utils.Success(c, "API is healthy", gin.H{
		"uptime": time.Now().Format(time.RFC3339),
	})
}

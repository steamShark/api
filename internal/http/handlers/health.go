package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Health struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewHealth(log *zap.Logger, db *gorm.DB) *Health {
	return &Health{log: log, db: db}
}

func (h *Health) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Health) Readyz(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false, "error": "db handle"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false, "error": "db ping failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ready": true})
}

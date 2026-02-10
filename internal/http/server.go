package http

import (
	"net/http"
	"steamshark-api/internal/config"
	"steamshark-api/internal/http/routes"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
@brief: start the HTTP server
*/
func New(cfg config.Config, log *zap.Logger, db *gorm.DB) *http.Server {
	router := routes.Build(cfg, log, db)

	return &http.Server{
		Addr:              cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * 1e9,  // 5s
		ReadTimeout:       15 * 1e9, // 15s
		WriteTimeout:      15 * 1e9, // 15s
		IdleTimeout:       60 * 1e9, // 60s
	}
}

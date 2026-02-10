package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"steamshark-api/internal/config"
	"steamshark-api/internal/db"
	httpserver "steamshark-api/internal/http"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	/* Get configuration */
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	/* Start Logger */
	logger := zap.Must(zap.NewProduction())
	defer func() { _ = logger.Sync() }()

	logger.Info("starting api", zap.String("env", cfg.Env), zap.String("addr", cfg.Port))

	/* Connect to DB */
	db, err := db.InitDB(cfg.DBPath)
	if err != nil { //If cannot conenct to db, just exit the program with error
		logger.Fatal("Database init error", zap.Error(err))
		os.Exit(1)
	} else {
		logger.Info("Database connected")
	}

	srv := httpserver.New(*cfg, logger, db)

	/* IMPLEMENT ROUTER */
	// Start http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()
	logger.Info("http server listening", zap.String("addr", cfg.Port))

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown requested")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}

	logger.Info("shutdown complete")
}

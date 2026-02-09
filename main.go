package main

import (
	"log"
	"os"
	controllers "steamshark-api/controllers/v1/website"
	"steamshark-api/db"
	"steamshark-api/helpers"
	"steamshark-api/middlewares"
	"steamshark-api/routes"
	"steamshark-api/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	/* Get configuration */
	config, err := helpers.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}
	/* Start gin instance */
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.SecurityHeaders()) //Use better security headers
	r.Use(middlewares.CORSPolicy())      //Use cors

	/* Start Logger */
	logger := zap.Must(zap.NewProduction())

	/* Cehck which env is going to be used */
	switch config.Env {
	case "production", "prod":
		gin.SetMode(gin.ReleaseMode)
		logger = zap.Must(zap.NewProduction())
	default:
		gin.SetMode(gin.DebugMode)
		logger = zap.Must(zap.NewDevelopment())
	}

	/* Start DB */
	db, err := db.InitDB(config.DBPath)
	if err != nil { //If cannot conenct to db, just exit the program with error
		logger.Fatal("Database init error", zap.Error(err))
		os.Exit(1)
	} else {
		logger.Info("Database connected")
	}

	/* START CONTROLLERS AND SERVICES */
	websiteService := services.NewWebsiteService(db)
	websiteController := controllers.NewWebsiteController(websiteService)

	occurrenceWebsiteService := services.NewOccurenceWebsiteService(db)
	occurrenceWebsiteController := controllers.NewOccurrenceWebsiteController(occurrenceWebsiteService)

	/* r.Use(metrics.MetricsMiddleware()) */
	/* r.GET("/metrics", gin.WrapH(promhttp.Handler())) */

	rateLimiter := middlewares.NewRateLimiter(1, 5) // Limit to 30 requests per second per IP
	r.Use(rateLimiter.Limit())

	/* IMPLEMENT ROUTER */
	router := routes.SetupRouter(websiteController, occurrenceWebsiteController)

	/* Start Server */
	if err := router.Run(config.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		logger.Fatal("Failed to start the server!", zap.Error(err))
		os.Exit(1)
	} else {
		logger.Info("API started!")
	}
}

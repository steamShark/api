package main

import (
	"log"
	controllers "steamshark-api/controllers/v1/website"
	"steamshark-api/db"
	"steamshark-api/helpers"
	"steamshark-api/middlewares"
	"steamshark-api/routes"
	"steamshark-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	envConfig, err := helpers.LoadEnvConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.SecurityHeaders()) //Use better security headers
	r.Use(middlewares.CORSPolicy())      //Use cors

	switch envConfig.Env {
	case "development":
		gin.SetMode(gin.DebugMode)
	case "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	host := envConfig.Host
	if host == "" {
		host = "localhost"
	}

	/* Start DB */
	db := db.InitUsersDB()

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

	port := envConfig.Port
	if port == "" {
		port = ":8800"
	}

	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

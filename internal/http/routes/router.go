package routes

import (
	"net/http"
	"steamshark-api/internal/http/handlers"
	"steamshark-api/internal/http/middlewares"
	"steamshark-api/internal/models"
	"steamshark-api/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
@brief: Function responsible to build the gin server, add the middlewares and store the routes to the app
*/
func Build(cfg models.Config, logger *zap.Logger, db *gorm.DB) http.Handler {
	/* Cehck which env is going to be used */
	switch cfg.Env {
	case "production", "prod":
		gin.SetMode(gin.ReleaseMode)
		logger = zap.Must(zap.NewProduction())
	default:
		gin.SetMode(gin.DebugMode)
		logger = zap.Must(zap.NewDevelopment())
	}

	/* Start gin instance */
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//Custom middleware
	router.Use(middlewares.SecurityHeaders()) //Use better security headers
	router.Use(middlewares.CORSPolicy())      //Use cors
	router.Use(middlewares.ZapLogger(logger))
	router.Use(middlewares.PrometheusMetrics())

	rateLimiter := middlewares.NewRateLimiter(1, 5) // Limit to 30 requests per second per IP
	router.Use(rateLimiter.Limit())

	/* Start Handlers */
	logger.Info("Starting http Handlers!")
	healthHandler := handlers.NewHealth(logger, db)
	websiteHandler := handlers.NewWebisteHandler(logger, db)
	occurrenceHandler := handlers.NewOccurrenceHandler(logger, db)

	/* v1 API Group */
	logger.Info("Starting http v1 API group!")
	v1 := router.Group("/api/v1")
	{
		// Health chec
		router.GET("/healthz", healthHandler.Healthz)
		router.GET("/readyz", healthHandler.Readyz)

		//GETS
		v1.GET("/websites", websiteHandler.ListWebsites)
		v1.GET("/websites/:identification", websiteHandler.GetByIdOrDomainOrURL)
		//Extension
		//v1.GET("/websites/extension", websiteController.GetExtensions)

		//POSTS
		//v1.POST("/websites", websiteHandler.Create)
		//VERIFY RECORDS
		//v1.POST("/websites/:id/verify", websiteHandler.VerifyWebsiteById) /* ADMIN ONLY */

		//PUT
		//v1.PUT("/websites/:id", websiteHandler.Update)

		//DELET
		//v1.DELETE("/websites/:id", websiteHandler.Delete)

		// Occurrences — standalone CRUD
		v1.GET("/occurrences", occurrenceHandler.ListOccurrences)
		v1.GET("/occurrences/:id", occurrenceHandler.GetOccurrence)
		v1.POST("/occurrences", occurrenceHandler.CreateOccurrence)
		v1.PUT("/occurrences/:id", occurrenceHandler.UpdateOccurrence)
		v1.DELETE("/occurrences/:id", occurrenceHandler.DeleteOccurrence)

		// Occurrences — scoped to a website
		v1.GET("/websites/:id/occurrences", occurrenceHandler.ListOccurrencesByWebsite)

	}

	router.NoRoute(func(ctx *gin.Context) {
		utils.Error(ctx, http.StatusNotFound, "Route not found!")
	})

	return router
}

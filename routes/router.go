package routes

import (
	healthController "steamshark-api/controllers/v1/health"
	websiteControllers "steamshark-api/controllers/v1/website"

	"github.com/gin-gonic/gin"
)

func SetupRouter(websiteController *websiteControllers.WebsiteController, occurrenceController *websiteControllers.OccurrenceWebsiteController) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", healthController.NewHealthController().Ping)

		//GETS
		v1.GET("/websites", websiteController.ListWebsites)
		v1.GET("/websites/:identification", websiteController.GetWebsite)
		//Extension
		v1.GET("/websites/extension", websiteController.GetExtensions)

		//POSTS
		v1.POST("/websites", websiteController.CreateWebsite)
		//VERIFY RECORDS
		v1.POST("/websites/:id/verify", websiteController.VerifyWebsiteById) /* ADMIN ONLY */

		//PUT
		v1.PUT("/websites/:id", websiteController.UpdateWebsite)

		//DELET
		v1.DELETE("/websites/:id", websiteController.DeleteWebsite)

	}
	return router
}

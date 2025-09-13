package controllers

import "steamshark-api/services"

// Wire this to your DI/container as you do for other controllers.
type OccurrenceWebsiteController struct {
	Service *services.OccurrenceWebsiteService
}

func NewOccurrenceWebsiteController(s *services.OccurrenceWebsiteService) *OccurrenceWebsiteController {
	return &OccurrenceWebsiteController{Service: s}
}

/*  */

func (ctrl *OccurrenceWebsiteController) GetOccurences() {

}

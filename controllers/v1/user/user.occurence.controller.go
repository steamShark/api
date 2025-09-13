package controllers

/*
import (
	"fmt"
	"net/http"
	"steamshark-api/models"
	"steamshark-api/services"
	"steamshark-api/utils"

	"github.com/gin-gonic/gin"
)

type OccurrenceController struct {
	Service *services.OccurrenceService
}

func NewOccurrenceController(service *services.OccurrenceService) *OccurrenceController {
	return &OccurrenceController{Service: service}
}

func (ctrl *OccurrenceController) GetBySteamID(c *gin.Context) {
	steamID := c.Param("steamid64")
	occurrences, err := ctrl.Service.GetOccurrencesBySteamID(steamID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "User or occurrences not found")
		return
	}
	utils.Success(c, "Occurrences fetched", occurrences)
}

func (ctrl *OccurrenceController) GetAll(c *gin.Context) {
	occurrences, err := ctrl.Service.GetAll()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch occurrences")
		return
	}
	utils.Success(c, "All occurrences retrieved", occurrences)
}

func (ctrl *OccurrenceController) GetByID(c *gin.Context) {
	varID := c.Param("id")
	var occID uint
	if _, err := fmt.Sscanf(varID, "%d", &occID); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid occurrence ID")
		return
	}

	occurrence, err := ctrl.Service.GetByID(occID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "Occurrence not found")
		return
	}

	utils.Success(c, "Occurrence found", occurrence)
}

func (ctrl *OccurrenceController) Create(c *gin.Context) {
	var occ models.Occurrence
	if err := c.ShouldBindJSON(&occ); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	fmt.Println(occ)

	if err := ctrl.Service.Create(&occ); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create occurrence")
		return
	}

	utils.Success(c, "Occurrence created", occ)
}
*/

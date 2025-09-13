package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"steamshark-api/dtos"
	helpers "steamshark-api/helpers/convertDTO"
	"steamshark-api/models"
	"steamshark-api/services"
	"steamshark-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Wire this to your DI/container as you do for other controllers.
type WebsiteController struct {
	Service *services.WebsiteService
}

func NewWebsiteController(s *services.WebsiteService) *WebsiteController {
	return &WebsiteController{Service: s}
}

// POST /websites
func (ctrl *WebsiteController) CreateWebsite(c *gin.Context) {
	var in dtos.WebsiteCreationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	//If both is scam and is Official are true
	/* if *in.IsNotTrusted && *in.IsOfficial {
		utils.Error(c, http.StatusBadRequest, "is_scam cannot be true and is_officiial be true at the same time")
		return
	} */

	created, err := ctrl.Service.CreateWebsite(c.Request.Context(), in)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	/* Convert into return DTO */
	createWebsiteReturnDTO, err := helpers.ConvertWebsiteModelDTOReturn(*created)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, "Website created", createWebsiteReturnDTO)
}

// GET /websites/:identification
func (ctrl *WebsiteController) GetWebsite(c *gin.Context) {
	identification := strings.TrimSpace(c.Param("identification"))
	if identification == "" {
		utils.Error(c, http.StatusBadRequest, "missing identification, it must be either id or website url")
		return
	}

	var (
		website *models.Website
		err     error
	)

	//Check if it's an uuid
	// If it's a UUID, fetch by ID; otherwise by domain.
	if uuid.Validate(identification) == nil {
		website, err = ctrl.Service.GetWebsiteByID(c.Request.Context(), identification)
	} else {
		website, err = ctrl.Service.GetWebsiteByDomain(c.Request.Context(), identification)
	}

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if website == nil {
		utils.Error(c, http.StatusNotFound, "website not found")
		return
	}

	//Append Occurences, to not appear null
	if website.Occurrences == nil {
		website.Occurrences = []models.Occurrence{}
	}

	utils.Success(c, "Website found", website)
}

// GET /websites
// Supports ?domain=&status=&risk_level=&limit=&offset=
func (ctrl *WebsiteController) ListWebsites(c *gin.Context) {
	var isNotTrusted bool
	var IsNotTrustedEnabled bool = false
	domain := strings.TrimSpace(c.Query("domain"))
	status := strings.TrimSpace(c.Query("status"))
	risk := strings.TrimSpace(c.Query("risk_level"))
	limit := utils.ParseIntDefault(c.Query("limit"), 50)
	offset := utils.ParseIntDefault(c.Query("offset"), 0)
	/* Convert the string into bool */
	if s := strings.TrimSpace(c.Query("is_not_trusted")); s != "" {
		parsed, err := strconv.ParseBool(s) // accepts: 1/0, t/f, true/false (any case)
		if err == nil {
			isNotTrusted = parsed
			IsNotTrustedEnabled = true
		}
	}

	res, err := ctrl.Service.ListWebsites(c.Request.Context(), services.ListWebsitesFilter{
		IsNotTrustedEnabled: &IsNotTrustedEnabled,
		IsNotTrusted:        &isNotTrusted,
		Domain:              domain,
		Status:              status,
		RiskLevel:           risk,
		Limit:               utils.Clamp(limit, 1, 200),
		Offset:              utils.Max(0, offset),
	})
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to list websites")
		return
	}

	// shape: { data, count, limit, offset }
	utils.Success(c, "Websites listed", gin.H{
		"data":   res.Items,
		"count":  res.Count,
		"limit":  res.Limit,
		"offset": res.Offset,
	})
}

// GET /websites/extension
// Supports ?domain=&status=&risk_level=&limit=&offset=
func (ctrl *WebsiteController) GetWebsitesExtension(c *gin.Context) {
	var isNotTrusted bool
	var IsNotTrustedEnabled bool = false
	if s := strings.TrimSpace(c.Query("is_not_trusted")); s != "" {
		parsed, err := strconv.ParseBool(s) // accepts: 1/0, t/f, true/false (any case)
		if err == nil {
			isNotTrusted = parsed
			IsNotTrustedEnabled = true
		}
	}

	res, err := ctrl.Service.GetWebsitesExtension(c.Request.Context(), services.ListWebsitesExtensionFilter{
		IsNotTrustedEnabled: &IsNotTrustedEnabled,
		IsNotTrusted:        &isNotTrusted,
	})
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to list websites")
		return
	}

	websites, err := helpers.ConvertListWebsiteModelDTOReturnExtension(*res)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to list websites")
		return
	}

	// shape: { data, count, limit, offset }
	utils.Success(c, "Websites listed", websites)
}

// PUT /websites/:id
func (ctrl *WebsiteController) UpdateWebsite(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		utils.Error(c, http.StatusBadRequest, "missing id")
		return
	}

	var in dtos.WebsiteUpdateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := ctrl.Service.UpdateWebsite(c.Request.Context(), id, in)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	//convert model to dto return
	createWebsiteReturnDTO, err := helpers.ConvertWebsiteModelDTOReturn(*updated)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, "Website updated", createWebsiteReturnDTO)
}

// DELETE /websites/:id
func (ctrl *WebsiteController) DeleteWebsite(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		utils.Error(c, http.StatusBadRequest, "missing id")
		return
	}
	if err := ctrl.Service.DeleteWebsite(c.Request.Context(), id); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	// 204 with no payload, or 200 with message — choose one. Keeping consistent with utils:
	utils.SuccessWithCode(c, http.StatusNoContent, "Website deleted", nil)
}

package controllers

import (
	"fmt"
	"net/http"
	"net/url"
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

type WebsiteController struct {
	Service *services.WebsiteService
}

func NewWebsiteController(s *services.WebsiteService) *WebsiteController {
	return &WebsiteController{Service: s}
}

/*
Create website, POST method

Needs:
  - in: Website dto for creations

Returns:
  - 200: already exists
  - 201: created
  - 400: bad request
  - 500: internal server error
*/
func (ctrl *WebsiteController) CreateWebsite(c *gin.Context) {
	var in dtos.WebsiteCreationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic input validation
	if strings.TrimSpace(in.Domain) == "" {
		utils.Error(c, http.StatusBadRequest, "domain is required")
		return
	}
	//add the rest of input validation for required fields

	/* VERIFICATIONS */
	//Verify if everything relate to the url are correct
	if in.URL != nil {
		parsed, err := url.Parse(*in.URL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			utils.Error(c, http.StatusBadRequest, "Invalid URL")
			return
		}

		// Domain consistency
		if !strings.Contains(parsed.Host, in.Domain) {
			utils.Error(c, http.StatusBadRequest, "Domain does not match URL")
			return
		}

		// Tld consistency
		if !strings.Contains(parsed.Host, in.TLD) {
			utils.Error(c, http.StatusBadRequest, "TLD does not match URL")
			return
		}
	}

	//If urls contains https assign SSLCertificate to true
	if in.URL != nil && strings.HasPrefix(strings.ToLower(*in.URL), "https") {
		ssl := true
		in.SSLCertificate = &ssl
	} else {
		ssl := false
		in.SSLCertificate = &ssl
	}

	//Verification of the percentage and risk level

	/* END VERIFICATIONS */

	existing, created, err := ctrl.Service.CreateWebsite(c.Request.Context(), in)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	/* If value already exists, return a 200 */
	if *existing {
		utils.Success(c, "Website already exists", created)
		return
	}
	//if not, return a 201 with the dto
	utils.SuccessCreated(c, "Webiste created", created)

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
	if website == nil {
		utils.Error(c, http.StatusNotFound, "website not found")
		return
	}
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
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
	page := utils.ParseIntDefault(c.Query("page"), 0)
	pageSize := utils.ParseIntDefault(c.Query("page_size"), 50)
	/* Convert the string into bool */
	if s := strings.TrimSpace(c.Query("is_not_trusted")); s != "" {
		parsed, err := strconv.ParseBool(s) // accepts: 1/0, t/f, true/false (any case)
		if err == nil {
			isNotTrusted = parsed
			IsNotTrustedEnabled = true
		}
	}

	res, err := ctrl.Service.ListWebsites(c.Request.Context(), models.Pagination{
		Page:     utils.Clamp(page, 0, 200),
		PageSize: utils.Max(0, pageSize),
	}, models.ListWebsitesFilter{
		IsNotTrustedEnabled: &IsNotTrustedEnabled,
		IsNotTrusted:        &isNotTrusted,
		Domain:              domain,
		Status:              status,
		RiskLevel:           risk,
	})
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to list websites")
		return
	}

	// shape: { data, count, limit, offset }
	utils.SuccessList(c, "Websites listed", gin.H{
		"data": res.Items,
	}, gin.H{
		"total":     res.Total,
		"page":      res.Page,
		"page_size": res.PageSize,
	})
}

// GET /websites/extension
// Supports ?domain=&status=&risk_level=&limit=&offset=
func (ctrl *WebsiteController) GetExtensions(c *gin.Context) {
	var isNotTrusted bool
	var IsNotTrustedEnabled bool = false
	if s := strings.TrimSpace(c.Query("is_not_trusted")); s != "" {
		parsed, err := strconv.ParseBool(s) // accepts: 1/0, t/f, true/false (any case)
		if err == nil {
			isNotTrusted = parsed
			IsNotTrustedEnabled = true
		}
	}

	res, err := ctrl.Service.GetWebsitesExtension(c.Request.Context(), models.ListWebsitesExtensionFilter{
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

	fmt.Println("in ", in)
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
	fmt.Println("website dto return ", createWebsiteReturnDTO)
	utils.Success(c, "Website updated", createWebsiteReturnDTO)
}

/*
Delete a website by id

Needs:
  - id of the webiste

Returns:
  - 204: no content, but it was successful
  - 400: Bad request
  - 401: unauthorized (made through middleware)
  - 403: no access (made through middleware)
  - 500: internal server error
*/
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
	// 204 with no payload
	utils.SuccessWithCode(c, http.StatusNoContent, "Website deleted", nil)
}

/*
Function to verify a website as an ADMIN

Needs:

  - id of the website

Returns:

  - 200: website verified
  - 400: bad request
  - 401: unauthorized (made through middleware)
  - 403: no access (made through middleware)
  - 500: internal server error
*/
func (ctrl *WebsiteController) VerifyWebsiteById(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		utils.Error(c, http.StatusBadRequest, "missing id")
		return
	}

	//verify the website
	website, err := ctrl.Service.VerifyWebsiteById(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 200 with message
	utils.Success(c, "Website verified", website)
}

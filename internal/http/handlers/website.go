package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"steamshark-api/internal/dtos"
	helpers "steamshark-api/internal/helpers/convertDTO"
	"steamshark-api/internal/models"
	"steamshark-api/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WebisteHandler struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewWebisteHandler(logger *zap.Logger, db *gorm.DB) *WebisteHandler {
	return &WebisteHandler{logger: logger, db: db}
}

/*
Gets all websites

# GET /websites

# Needs:
  - id of the webiste

# Query params:
  - domain
  - status
  - risk_level
  - page
  - page_size

# Returns:
  - 200: success
  - 400: Bad request
  - 401: unauthorized (made through middleware)
  - 403: no access (made through middleware)
  - 500: internal server error
*/
func (handler *WebisteHandler) ListWebsites(ctx *gin.Context) {
	var isNotTrusted bool
	var IsNotTrustedEnabled bool = false
	domain := strings.TrimSpace(ctx.Query("domain"))
	status := strings.TrimSpace(ctx.Query("status"))
	risk := strings.TrimSpace(ctx.Query("risk_level"))
	page := utils.ParseIntDefault(ctx.Query("page"), 0)
	pageSize := utils.ParseIntDefault(ctx.Query("page_size"), 50)
	/* Convert the string into bool */
	if s := strings.TrimSpace(ctx.Query("is_not_trusted")); s != "" {
		parsed, err := strconv.ParseBool(s) // accepts: 1/0, t/f, true/false (any case)
		if err == nil {
			isNotTrusted = parsed
			IsNotTrustedEnabled = true
		}
	}

	filters := models.ListWebsitesFilter{
		IsNotTrustedEnabled: &IsNotTrustedEnabled,
		IsNotTrusted:        &isNotTrusted,
		Domain:              domain,
		Status:              status,
		RiskLevel:           risk,
	}

	pagination := models.Pagination{
		Page:     utils.Clamp(page, 0, 200),
		PageSize: utils.Max(0, pageSize),
	}

	/* res, err := ctrl.Service.ListWebsites(ctx.Request.Context(), models.Pagination{
		Page:     utils.Clamp(page, 0, 200),
		PageSize: utils.Max(0, pageSize),
	}, )
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to list websites")
		return
	} */

	query := handler.db.WithContext(ctx).Model(&models.Website{})

	fmt.Println("query ", query)

	/* VERIFY IF THERE IS IN PARAMS */
	if *filters.IsNotTrustedEnabled {
		if filters.IsNotTrusted != nil {
			query = query.Where("is_not_trusted = ?", filters.IsNotTrusted)
		}
	}
	if filters.Domain != "" {
		query = query.Where("domain LIKE ?", "%"+filters.Domain+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.RiskLevel != "" {
		query = query.Where("risk_level = ?", filters.RiskLevel)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		handler.logger.Error("error while getting total of items " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "Error while trying to found the websites")
		return
	}

	fmt.Println("total ", total)

	var items []models.Website
	if err := query. /* Preload("Website"). */
				Order("updated_at DESC").
				Limit(pagination.PageSize).
				Offset((pagination.Page - 1) * pagination.PageSize).
				Find(&items).Error; err != nil {
		handler.logger.Error("error while processing the list websites " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "Error while trying to found the websites")
		return
	}

	fmt.Println(items)

	// shape: { data, count, limit, offset }
	utils.SuccessList(ctx, "Websites listed", gin.H{
		"data": items,
	}, gin.H{
		"total":     total,
		"page":      pagination.Page,
		"page_size": pagination.PageSize,
	})
}

/*
Gets a specific website, by identification, id or the domain of the website (eg.: steamcommunity.com)

# GET /websites/:identification

# Needs:
  - identification of the webiste (id, or domain (eg.: steamcommunity.com))

# Returns:
  - 200: success
  - 400: Bad request
  - 401: unauthorized (made through middleware)
  - 403: no access (made through middleware)
  - 500: internal server error
*/
func (handler *WebisteHandler) GetByIdorDomain(ctx *gin.Context) {
	identification := strings.TrimSpace(ctx.Param("identification"))
	if identification == "" {
		handler.logger.Error("missing identification, it must be either id or website url")
		utils.Error(ctx, http.StatusBadRequest, "missing identification, it must be either id or website url")
		return
	}

	var website *models.Website

	//Check if it's an uuid
	// If it's a UUID, fetch by ID; otherwise by domain.
	if uuid.Validate(identification) == nil {
		db := handler.db.WithContext(ctx).
			Preload("Occurrences", func(tx *gorm.DB) *gorm.DB {
				return tx.Order("created_at DESC")
			})

		// First try lookup by ID
		err := db.Where("id = ?", identification).First(&website).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				handler.logger.Error("no website found with provided id")
				utils.Error(ctx, http.StatusNotFound, "no website found with provided id")
				return
			} else { // DB error on ID lookup
				handler.logger.Error("Error while searching for the id " + err.Error())
				utils.Error(ctx, http.StatusNotFound, "no website found with provided id")
				return
			}
		}
	} else { //if it's to search by domain/name
		db := handler.db.WithContext(ctx).
			Preload("Occurrences", func(tx *gorm.DB) *gorm.DB {
				return tx.Order("created_at DESC")
			})

		err := db.Where("domain = ?", identification).First(&website).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handler.logger.Error("Error while searching for the domain " + err.Error())
			utils.Error(ctx, http.StatusNotFound, "no website found with provided domain")
			return
		}
	}

	handler.logger.Error("Website found!")
	utils.Success(ctx, "Website found", website)
}

/*
# Create website, POST method

# POST /websites

# Needs:
  - in: Website dto for creations

# Returns:
  - 200: already exists
  - 201: created
  - 400: bad request
  - 500: internal server error
*/
func (handler *WebisteHandler) Create(ctx *gin.Context) {
	var in dtos.WebsiteCreationInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		handler.logger.Error("Invalid request body. Error: " + err.Error())
		utils.Error(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic input validation
	if strings.TrimSpace(in.Domain) == "" {
		handler.logger.Error("A domain is required to create a new website")
		utils.Error(ctx, http.StatusBadRequest, "domain is required")
		return
	}
	//add the rest of input validation for required fields

	/* VERIFICATIONS */
	//Verify if everything relate to the url are correct
	if in.URL != nil {
		parsed, err := url.Parse(*in.URL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			utils.Error(ctx, http.StatusBadRequest, "Invalid URL")
			return
		}

		// Domain consistency
		if !strings.Contains(parsed.Host, in.Domain) {
			utils.Error(ctx, http.StatusBadRequest, "Domain does not match URL")
			return
		}

		// Tld consistency
		if !strings.Contains(parsed.Host, in.TLD) {
			utils.Error(ctx, http.StatusBadRequest, "TLD does not match URL")
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

	/* CREATE THE MODLE THROUGH THE DTO */
	modelWebsite, err := helpers.ConvertWebsiteDTOModelCreation(in)
	if err != nil {
		handler.logger.Error("An error occured while converting the dto into object. Error: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "an error occurred while creating the website")
		return
	}

	// Idempotent create: if domain already exists, return the existing record
	var existing models.Website
	if err := handler.db.WithContext(ctx).Where("domain = ?", modelWebsite.Domain).First(&existing).Error; err == nil {
		/* Create the DTO inside the service and return DTO to controller */
		existingWebsiteReturnDTO, err := helpers.ConvertWebsiteModelDTOReturn(*modelWebsite)
		if err != nil {
			handler.logger.Error("An error occured while converting the dto into object. Error: " + err.Error())
			utils.Error(ctx, http.StatusInternalServerError, "an error occurred while creating the website")
			return
		}
		handler.logger.Error("Website was already created, send the current version")
		utils.Success(ctx, "Website already exists!", existingWebsiteReturnDTO)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		handler.logger.Error("An error occured searching for a possible duplicaiton in db. Error: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "an error occurred while creating the website")
		return
	}

	//verify and determine the is_not_trusted, otherwise is already
	if modelWebsite.RiskScore <= 20 {
		modelWebsite.IsNotTrusted = utils.ReturnPointerBool(false)
	}

	//set the risk level (text)
	if modelWebsite.RiskScore == 0 && modelWebsite.IsOfficial { //
		modelWebsite.RiskLevel = "none"
	} else if modelWebsite.RiskScore > 0 && modelWebsite.RiskScore <= 10 { /*  */
		modelWebsite.RiskLevel = "low"
	} else if modelWebsite.RiskScore > 10 && modelWebsite.RiskScore <= 50 { /*  */
		modelWebsite.RiskLevel = "medium"
	} else if modelWebsite.RiskScore > 50 && modelWebsite.RiskScore < 90 { /*  */
		modelWebsite.RiskLevel = "high"
	} else if modelWebsite.RiskScore >= 90 && modelWebsite.RiskScore <= 100 { /*  */
		modelWebsite.RiskLevel = "critical"
	} else { /* just to be sure */
		modelWebsite.RiskLevel = "unknown"
	}
	// Create
	if err := handler.db.Create(&modelWebsite).Error; err != nil {
		handler.logger.Error("An error occured while trrying to create the record on the database. Error: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "an error occurred while creating the website")
		return
	}

	/* Convert into return DTO */
	createWebsiteReturnDTO, err := helpers.ConvertWebsiteModelDTOReturn(*modelWebsite)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "an error occured while creating the website")

	}
	//if not, return a 201 with the dto
	utils.SuccessCreated(ctx, "Webiste created", createWebsiteReturnDTO)
}

/*
Update a website by id

# PUT /websites/:id

# Needs:
  - id of the webiste

# Returns:
  - 200: updated
  - 400: Bad request
  - 401: unauthorized (made through middleware)
  - 403: no access (made through middleware)
  - 500: internal server error
*/
func (handler *WebisteHandler) Update(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		handler.logger.Error("missing an id in request for update.")
		utils.Error(ctx, http.StatusBadRequest, "missing id")
		return
	}

	var in dtos.WebsiteUpdateInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		handler.logger.Error("Invalid request body. " + err.Error())
		utils.Error(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	fmt.Println("in ", in)
	/* convert dto to model */
	updateWebsite, err := helpers.ConvertWebsiteDTOModelUpdate(in)
	if err != nil {
		handler.logger.Error("Error converting dto to object. " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error while updating the website")
		return
	}

	//Check if there is with th id given
	var website models.Website
	if err := handler.db.WithContext(ctx).First(&website, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handler.logger.Error("website with id " + id + " not found")
			utils.Error(ctx, http.StatusInternalServerError, "website with id "+id+" not found")
			return
		}
		handler.logger.Error("Error while trying to check if website exists. " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error while updating the website")
		return
	}

	// Apply updates
	if err := handler.db.WithContext(ctx).Model(&website).Updates(updateWebsite).Error; err != nil {
		handler.logger.Error("Error while updating website. " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error while updating the website")
		return
	}

	//convert model to dto return
	createWebsiteReturnDTO, err := helpers.ConvertWebsiteModelDTOReturn(website)
	if err != nil {
		handler.logger.Error("Error while converting the model into DTO. " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error while updating the website")
		return
	}
	fmt.Println("website dto return ", createWebsiteReturnDTO)
	utils.Success(ctx, "Website updated", createWebsiteReturnDTO)
}

/*
Delete a website by id

# DELETE /websites/:id

# Needs:
  - id of the webiste

# Returns:
  - 204: no content, but it was successful
  - 400: Bad request
  - 401: unauthorized (made through middleware)
  - 403: no access (made through middleware)
  - 500: internal server error
*/
func (handler *WebisteHandler) Delete(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		handler.logger.Error("missing an id in request for update.")
		utils.Error(ctx, http.StatusBadRequest, "missing id")
		return
	}
	//Delete the
	result := handler.db.WithContext(ctx).Delete(&models.Website{}, "id = ?", id)
	if result.Error != nil {
		handler.logger.Error("Error while trying to delete from database. " + result.Error.Error())
		utils.Error(ctx, http.StatusInternalServerError, "Error while trying to delete from database")
		return
	}
	// 204 with no payload
	utils.SuccessWithCode(ctx, http.StatusNoContent, "Website deleted", nil)
}

/*
Function to verify a website as an ADMIN

# POST /webistes/:id

# Needs:
  - id of the website

# Returns:
  - 200: website verified
  - 400: bad request
  - 401: unauthorized (made through middleware)
  - 403: no access (made through middleware)
  - 500: internal server error
*/
func (handler *WebisteHandler) VerifyWebsiteById(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		handler.logger.Error("missing an id in request for update.")
		utils.Error(ctx, http.StatusBadRequest, "missing id")
		return
	}

	var website models.Website

	//verifiedWebsite := s.DB.WithContext(ctx).Model(&models.Website{}).Where("id = ?", id).Update("verified", true)

	// Load the website
	if err := handler.db.WithContext(ctx).
		First(&website, "id = ?", id).Error; err != nil {
		handler.logger.Error("an error occured while verifying the website. " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "an error occured while verifying the website")
		return
	}

	// Update the field in Go
	website.Verified = true

	// Persist changes
	if err := handler.db.WithContext(ctx).Save(&website).Error; err != nil {
		handler.logger.Error("an error occured while verifying the website on update. " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "an error occured while verifying the website")
		return
	}

	// 200 with message
	utils.Success(ctx, "Website verified", website)
}

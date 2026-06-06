package handlers

import (
	"errors"
	"net/http"
	"steamshark-api/internal/dtos"
	"steamshark-api/internal/models"
	"steamshark-api/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OccurrenceHandler struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewOccurrenceHandler(logger *zap.Logger, db *gorm.DB) *OccurrenceHandler {
	return &OccurrenceHandler{logger: logger, db: db}
}

/*
List all occurrences across all websites

GET /occurrences

Query params:
  - website_id
  - severity
  - status
  - page
  - page_size
*/
func (h *OccurrenceHandler) ListOccurrences(ctx *gin.Context) {
	page := utils.ParseIntDefault(ctx.Query("page"), 1)
	pageSize := utils.ParseIntDefault(ctx.Query("page_size"), 50)
	pagination := models.Pagination{
		Page:     utils.Clamp(page, 1, 200),
		PageSize: utils.Max(1, pageSize),
	}

	query := h.db.WithContext(ctx).Model(&models.Occurrence{})

	if websiteID := strings.TrimSpace(ctx.Query("website_id")); websiteID != "" {
		query = query.Where("website_id = ?", websiteID)
	}
	if severity := strings.TrimSpace(ctx.Query("severity")); severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if status := strings.TrimSpace(ctx.Query("status")); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		h.logger.Error("error counting occurrences: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error listing occurrences")
		return
	}

	var items []models.Occurrence
	if err := query.Order("created_at DESC").
		Limit(pagination.PageSize).
		Offset((pagination.Page - 1) * pagination.PageSize).
		Find(&items).Error; err != nil {
		h.logger.Error("error fetching occurrences: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error listing occurrences")
		return
	}

	utils.SuccessList(ctx, "Occurrences listed", gin.H{"data": items}, gin.H{
		"total":     total,
		"page":      pagination.Page,
		"page_size": pagination.PageSize,
	})
}

/*
Get a single occurrence by its own ID

GET /occurrences/:id
*/
func (h *OccurrenceHandler) GetOccurrence(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		utils.Error(ctx, http.StatusBadRequest, "missing occurrence id")
		return
	}

	var occurrence models.Occurrence
	if err := h.db.WithContext(ctx).First(&occurrence, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(ctx, http.StatusNotFound, "occurrence not found")
			return
		}
		h.logger.Error("error fetching occurrence: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error fetching occurrence")
		return
	}

	utils.Success(ctx, "Occurrence found", toOccurrenceDTO(occurrence))
}

/*
Create an occurrence

POST /occurrences

Body must include website_id.
*/
func (h *OccurrenceHandler) CreateOccurrence(ctx *gin.Context) {
	var in dtos.OccurrenceCreationInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		h.logger.Error("invalid request body: " + err.Error())
		utils.Error(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.db.WithContext(ctx).First(&models.Website{}, "id = ?", in.WebsiteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(ctx, http.StatusNotFound, "website not found")
			return
		}
		h.logger.Error("error checking website: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error creating occurrence")
		return
	}

	severity := in.Severity
	if severity == "" {
		severity = "medium"
	}

	occurrence := models.Occurrence{
		WebsiteID:   in.WebsiteID,
		Description: in.Description,
		URLReported: in.URLReported,
		CountryCode: in.CountryCode,
		Severity:    severity,
	}

	if err := h.db.WithContext(ctx).Create(&occurrence).Error; err != nil {
		h.logger.Error("error creating occurrence: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error creating occurrence")
		return
	}

	utils.SuccessCreated(ctx, "Occurrence created", toOccurrenceDTO(occurrence))
}

/*
Update an occurrence by its own ID

PUT /occurrences/:id
*/
func (h *OccurrenceHandler) UpdateOccurrence(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		utils.Error(ctx, http.StatusBadRequest, "missing occurrence id")
		return
	}

	var in dtos.OccurrenceUpdateInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		h.logger.Error("invalid request body: " + err.Error())
		utils.Error(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	var occurrence models.Occurrence
	if err := h.db.WithContext(ctx).First(&occurrence, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(ctx, http.StatusNotFound, "occurrence not found")
			return
		}
		h.logger.Error("error fetching occurrence: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error updating occurrence")
		return
	}

	updates := map[string]any{}
	if in.Description != nil {
		updates["description"] = in.Description
	}
	if in.URLReported != nil {
		updates["url_reported"] = *in.URLReported
	}
	if in.CountryCode != nil {
		updates["country_code"] = in.CountryCode
	}
	if in.Severity != nil {
		updates["severity"] = *in.Severity
	}
	if in.Status != nil {
		updates["status"] = *in.Status
	}

	if err := h.db.WithContext(ctx).Model(&occurrence).Updates(updates).Error; err != nil {
		h.logger.Error("error updating occurrence: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error updating occurrence")
		return
	}

	utils.Success(ctx, "Occurrence updated", toOccurrenceDTO(occurrence))
}

/*
Delete an occurrence by its own ID

DELETE /occurrences/:id
*/
func (h *OccurrenceHandler) DeleteOccurrence(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		utils.Error(ctx, http.StatusBadRequest, "missing occurrence id")
		return
	}

	result := h.db.WithContext(ctx).Delete(&models.Occurrence{}, "id = ?", id)
	if result.Error != nil {
		h.logger.Error("error deleting occurrence: " + result.Error.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error deleting occurrence")
		return
	}
	if result.RowsAffected == 0 {
		utils.Error(ctx, http.StatusNotFound, "occurrence not found")
		return
	}

	utils.SuccessWithCode(ctx, http.StatusNoContent, "Occurrence deleted", nil)
}

/*
List occurrences scoped to a specific website

GET /websites/:id/occurrences
*/
func (h *OccurrenceHandler) ListOccurrencesByWebsite(ctx *gin.Context) {
	websiteID := strings.TrimSpace(ctx.Param("id"))
	if websiteID == "" {
		utils.Error(ctx, http.StatusBadRequest, "missing website id")
		return
	}

	if err := h.db.WithContext(ctx).First(&models.Website{}, "id = ?", websiteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(ctx, http.StatusNotFound, "website not found")
			return
		}
		h.logger.Error("error checking website: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error listing occurrences")
		return
	}

	page := utils.ParseIntDefault(ctx.Query("page"), 1)
	pageSize := utils.ParseIntDefault(ctx.Query("page_size"), 50)
	pagination := models.Pagination{
		Page:     utils.Clamp(page, 1, 200),
		PageSize: utils.Max(1, pageSize),
	}

	query := h.db.WithContext(ctx).Model(&models.Occurrence{}).Where("website_id = ?", websiteID)

	if severity := strings.TrimSpace(ctx.Query("severity")); severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if status := strings.TrimSpace(ctx.Query("status")); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		h.logger.Error("error counting occurrences: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error listing occurrences")
		return
	}

	var items []models.Occurrence
	if err := query.Order("created_at DESC").
		Limit(pagination.PageSize).
		Offset((pagination.Page - 1) * pagination.PageSize).
		Find(&items).Error; err != nil {
		h.logger.Error("error fetching occurrences: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "error listing occurrences")
		return
	}

	utils.SuccessList(ctx, "Occurrences listed", gin.H{"data": items}, gin.H{
		"total":     total,
		"page":      pagination.Page,
		"page_size": pagination.PageSize,
	})
}

func toOccurrenceDTO(o models.Occurrence) dtos.OccurrenceReturnDTO {
	return dtos.OccurrenceReturnDTO{
		ID:          o.ID,
		WebsiteID:   o.WebsiteID,
		Description: o.Description,
		URLReported: o.URLReported,
		CountryCode: o.CountryCode,
		Severity:    o.Severity,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   o.UpdatedAt.Format(time.RFC3339),
	}
}
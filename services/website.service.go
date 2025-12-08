package services

import (
	"context"
	"errors"
	"fmt"
	"steamshark-api/dtos"
	helpers "steamshark-api/helpers/convertDTO"
	"steamshark-api/models"
	"steamshark-api/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebsiteService struct {
	DB *gorm.DB
}

func NewWebsiteService(db *gorm.DB) *WebsiteService {
	return &WebsiteService{DB: db}
}

/*
Create an website through an DTO provided

@params:

	ctx
	website dto

@returns:

	existing: to now if already exists or no
	dto website model: actual object in dto
	error: if has an error
*/
func (s *WebsiteService) CreateWebsite(ctx context.Context, in dtos.WebsiteCreationInput) (*bool, *dtos.WebsiteReturnDTO, error) { /* dtos.WebsiteDTO */
	/* CREATE THE MODLE THROUGH THE DTO */
	modelWebsite, err := helpers.ConvertWebsiteDTOModelCreation(in)
	if err != nil {
		return nil, nil, errors.New("an error occurred while creating the website")
	}

	// Idempotent create: if domain already exists, return the existing record
	var existing models.Website
	if err := s.DB.WithContext(ctx).Where("domain = ?", modelWebsite.Domain).First(&existing).Error; err == nil {
		/* Create the DTO inside the service and return DTO to controller */
		existingWebsiteReturnDTO, err := helpers.ConvertWebsiteModelDTOReturn(*modelWebsite)
		if err != nil {
			return nil, nil, errors.New("an error occurred while creating the website")
		}
		return utils.ReturnPointerBool(true), existingWebsiteReturnDTO, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
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
	if err := s.DB.Create(&modelWebsite).Error; err != nil {
		return nil, nil, err
	}

	/* Convert into return DTO */
	createWebsiteReturnDTO, err := helpers.ConvertWebsiteModelDTOReturn(*modelWebsite)
	if err != nil {
		return nil, nil, errors.New("an error occurred while creating the website")
	}

	return utils.ReturnPointerBool(false), createWebsiteReturnDTO, nil
}

/*
List all websites with pagination and custom params

@pagination:

	Limit     int
	Offset    int

@returns:

	PaginatedListResult: List with pagination of the Websites with params
*/
func (s *WebsiteService) ListWebsites(ctx context.Context, pagination models.Pagination, filters models.ListWebsitesFilter) (*models.PaginatedListResult[models.Website], error) {
	q := s.DB.WithContext(ctx).Model(&models.Website{})

	fmt.Println("ListWebsitesFilter ", filters)

	/* VERIFY IF THERE IS IN PARAMS */
	if *filters.IsNotTrustedEnabled {
		if filters.IsNotTrusted != nil {
			q = q.Where("is_not_trusted = ?", filters.IsNotTrusted)
		}
	}
	if filters.Domain != "" {
		q = q.Where("domain LIKE ?", "%"+filters.Domain+"%")
	}
	if filters.Status != "" {
		q = q.Where("status = ?", filters.Status)
	}
	if filters.RiskLevel != "" {
		q = q.Where("risk_level = ?", filters.RiskLevel)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []models.Website
	if err := q. /* Preload("Website"). */
			Order("updated_at DESC").
			Limit(pagination.PageSize).
			Offset((pagination.Page - 1) * pagination.PageSize).
			Find(&items).Error; err != nil {
		return nil, err
	}

	fmt.Println(items)

	return &models.PaginatedListResult[models.Website]{Items: items, Total: total, Page: pagination.Page, PageSize: pagination.PageSize}, nil
}

/*
List a unique website weither with id or url
*/
func (s *WebsiteService) GetWebsiteByID(ctx context.Context, identification string) (*models.Website, error) {
	if identification == "" {
		return nil, errors.New("identification is required")
	}

	var w models.Website
	db := s.DB.WithContext(ctx).
		Preload("Occurrences", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("created_at DESC")
		})

	//Check if it's an uuid
	if err := uuid.Validate(identification); err == nil {

	}

	// First try lookup by ID
	err := db.Where("id = ?", identification).First(&w).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no website found with provided id")
		} else {
			return nil, err // DB error on ID lookup
		}
	}

	if w.Occurrences == nil {
		w.Occurrences = []models.Occurrence{}
	}

	return &w, nil
}

func (s *WebsiteService) GetWebsitesExtension(ctx context.Context, f models.ListWebsitesExtensionFilter) (*[]models.Website, error) {
	q := s.DB.WithContext(ctx).Model(&models.Website{})

	fmt.Println("ListWebsitesFilter ", f)

	/* VERIFY IF THERE IS IN PARAMS */
	if *f.IsNotTrustedEnabled {
		if f.IsNotTrusted != nil {
			q = q.Where("is_not_trusted = ?", f.IsNotTrusted)
		}
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []models.Website
	if err := q.Preload("Occurrences").
		Order("updated_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return &items, nil
}

/*
Get website by domain
*/
func (s *WebsiteService) GetWebsiteByDomain(ctx context.Context, identification string) (*models.Website, error) {
	if identification == "" {
		return nil, errors.New("identification is required")
	}

	var w models.Website
	db := s.DB.WithContext(ctx).
		Preload("Occurrences", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("created_at DESC")
		})

	err := db.Where("domain = ?", identification).First(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("no website found with provided domain")
	}
	if err != nil {
		return nil, err // DB error on domain lookup
	}

	if w.Occurrences == nil {
		w.Occurrences = []models.Occurrence{}
	}

	return &w, nil
}

/*
Delete Website deletes a website (and cascades its occurrences).
*/
func (s *WebsiteService) UpdateWebsite(ctx context.Context, id string, in dtos.WebsiteUpdateInput) (*models.Website, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	/* convert dto to model */
	updateWebsite, err := helpers.ConvertWebsiteDTOModelUpdate(in)
	if err != nil {
		return nil, err
	}

	//Check if there is with th id given
	var website models.Website
	if err := s.DB.WithContext(ctx).First(&website, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("website with id %s not found", id)
		}
		return nil, err
	}

	// Apply updates
	if err := s.DB.WithContext(ctx).Model(&website).Updates(updateWebsite).Error; err != nil {
		return nil, err
	}

	return &website, nil
}

/*
Delete Website deletes a website (and cascades its occurrences).
*/
func (s *WebsiteService) DeleteWebsite(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	result := s.DB.WithContext(ctx).Delete(&models.Website{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

/*
Verify the an website by id
*/
func (s *WebsiteService) VerifyWebsiteById(ctx context.Context, id string) (*models.Website, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	var website models.Website

	//verifiedWebsite := s.DB.WithContext(ctx).Model(&models.Website{}).Where("id = ?", id).Update("verified", true)

	// Load the website
	if err := s.DB.WithContext(ctx).
		First(&website, "id = ?", id).Error; err != nil {
		return nil, err // handle not found / db error
	}

	// Update the field in Go
	website.Verified = true

	// Persist changes
	if err := s.DB.WithContext(ctx).Save(&website).Error; err != nil {
		return nil, err
	}

	// website now has the updated data
	return &website, nil
}

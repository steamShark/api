package services

import (
	"context"
	"errors"
	"fmt"
	"steamshark-api/db"
	"steamshark-api/dtos"
	helpers "steamshark-api/helpers/convertDTO"
	"steamshark-api/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Filter the /websites list
type ListWebsitesFilter struct {
	IsNotTrustedEnabled *bool
	IsNotTrusted        *bool
	Domain              string
	Status              string
	RiskLevel           string
	Limit               int
	Offset              int
}

type ListWebsitesExtensionFilter struct {
	IsNotTrustedEnabled *bool
	IsNotTrusted        *bool
}

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

	website model
	error
*/
func (s *WebsiteService) CreateWebsite(ctx context.Context, in dtos.WebsiteCreationInput) (*models.Website, error) { /* dtos.WebsiteDTO */
	// Basic validation
	if strings.TrimSpace(in.Domain) == "" {
		return nil, errors.New("domain is required")
	}

	/* CREATE THE MODLE THROUGH THE DTO */
	modelWebsite, err := helpers.ConvertWebsiteDTOModelCreation(in)
	if err != nil {
		return nil, errors.New("an error occurred while creating the website")
	}

	db := s.DB.WithContext(ctx)

	// Idempotent create: if domain already exists, return the existing record
	var existing models.Website
	if err := db.Where("domain = ?", modelWebsite.Domain).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if modelWebsite.RiskLevel == "low" {
		b := false
		modelWebsite.IsNotTrusted = &b
	} else {
		b := true
		modelWebsite.IsNotTrusted = &b
	}

	// Create
	if err := db.Create(&modelWebsite).Error; err != nil {
		return nil, err
	}

	return modelWebsite, nil
}

/*
List all websites with pagination and custom params

@pagination:

	Limit     int
	Offset    int

@returns:

	PaginatedListResult: List with pagination of the Websites with params
*/
func (s *WebsiteService) ListWebsites(ctx context.Context, f ListWebsitesFilter) (*db.PaginatedListResult[models.Website], error) {
	q := s.DB.WithContext(ctx).Model(&models.Website{})

	fmt.Println("ListWebsitesFilter ", f)

	/* VERIFY IF THERE IS IN PARAMS */
	if *f.IsNotTrustedEnabled {
		if f.IsNotTrusted != nil {
			q = q.Where("is_not_trusted = ?", f.IsNotTrusted)
		}
	}
	if f.Domain != "" {
		q = q.Where("domain LIKE ?", "%"+f.Domain+"%")
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.RiskLevel != "" {
		q = q.Where("risk_level = ?", f.RiskLevel)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []models.Website
	if err := q.Preload("Occurrences").
		Order("updated_at DESC").
		Limit(f.Limit).
		Offset(f.Offset).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return &db.PaginatedListResult[models.Website]{Items: items, Count: total, Limit: f.Limit, Offset: f.Offset}, nil
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

func (s *WebsiteService) GetWebsitesExtension(ctx context.Context, f ListWebsitesExtensionFilter) (*[]models.Website, error) {
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

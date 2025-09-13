package services

import "gorm.io/gorm"

type OccurrenceWebsiteService struct {
	DB *gorm.DB
}

func NewOccurenceWebsiteService(db *gorm.DB) *OccurrenceWebsiteService {
	return &OccurrenceWebsiteService{DB: db}
}

/* Get all occurrences */
func (s *OccurrenceWebsiteService) GetOccurrences() {

}

package services

/*
import (
	"fmt"
	"steamshark-api/models"

	"gorm.io/gorm"
)

type OccurrenceService struct {
	DB *gorm.DB
}

func NewOccurrenceService(db *gorm.DB) *OccurrenceService {
	return &OccurrenceService{DB: db}
}

func (s *OccurrenceService) GetOccurrencesBySteamID(steamid64 string) ([]models.Occurrence, error) {
	var user models.User
	if err := s.DB.Where("steamid64 = ?", steamid64).First(&user).Error; err != nil {
		return nil, err
	}

	var occurrences []models.Occurrence
	if err := s.DB.Where("steamid64 = ?", user.SteamID64).Order("created_at desc").Find(&occurrences).Error; err != nil {
		return nil, err
	}

	return occurrences, nil
}

func (s *OccurrenceService) GetAll() ([]models.Occurrence, error) {
	var occurrences []models.Occurrence
	err := s.DB.Order("created_at desc").Find(&occurrences).Error
	return occurrences, err
}

func (s *OccurrenceService) GetByID(id uint) (*models.Occurrence, error) {
	var occ models.Occurrence
	err := s.DB.First(&occ, id).Error
	if err != nil {
		return nil, err
	}
	return &occ, nil
}

func (s *OccurrenceService) Create(occ *models.Occurrence) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		userSvc := &UserService{db: tx}

		fmt.Println(userSvc)

		err := userSvc.CreateFromOccurrence(occ.SteamID64, !occ.IsTrusted)
		if err != nil {
			return err
		}

		if err := tx.Create(occ).Error; err != nil {
			return err
		}

		return nil
	})
}
*/

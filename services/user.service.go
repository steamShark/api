package services

/*
import (
	"errors"
	"fmt"

	"steamshark-api/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserBySteamID(steamID string) (*models.User, error) {
	fmt.Println(steamID)
	var user models.User
	result := s.db.First(&user, "steamid64 = ?", steamID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Function to create the user from an occurrence
func (s *UserService) CreateFromOccurrence(steamID string, suspicious bool) error {
	var user models.User
	err := s.db.First(&user, "steamid64 = ?", steamID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = models.User{
				SteamID64: steamID,
			}
		} else {
			return err
		}
	}

	if suspicious {
		user.Status = "suspicious"
	} else {
		user.Status = "clean"
	}

	return s.db.Save(&user).Error
}
*/

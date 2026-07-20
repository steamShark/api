package model_user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

type User struct {
	ID          string    `json:"id" gorm:"primaryKey;size:36"`
	SteamID     string    `json:"steam_id" gorm:"uniqueIndex;not null;size:32"`
	DisplayName string    `json:"display_name" gorm:"size:255"`
	AvatarURL   string    `json:"avatar_url" gorm:"size:512"`
	Role        string    `json:"role" gorm:"size:16;default:user"` // user, admin
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

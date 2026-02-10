package model_website

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (o *Occurrence) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return nil
}

type Occurrence struct {
	ID          string    `json:"id" gorm:"primaryKey;size:36"`
	WebsiteID   string    `json:"website_id" gorm:"size:36;index;not null"`
	Description *string   `json:"description"`
	URLReported string    `json:"url_reported" gorm:"size:2048"`
	CountryCode *string   `json:"country_code" gorm:"size:2"`
	Severity    string    `json:"severity" gorm:"size:16;default:medium"` // info, low, medium, high, critical
	Status      string    `json:"status" gorm:"size:16;default:pending"`  // pending, verified, rejected
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

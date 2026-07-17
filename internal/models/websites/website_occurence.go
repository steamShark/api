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
	ID              string     `json:"id" gorm:"primaryKey;size:36"`
	WebsiteID       string     `json:"website_id" gorm:"size:36;index;not null"`
	Type            string     `json:"type" gorm:"size:32;default:other"`          // phishing, scam, malware, impersonation, fake_giveaway, credential_harvesting, other
	Source          string     `json:"source" gorm:"size:16;default:user_report"`  // user_report, automated_scan, admin
	Description     *string    `json:"description"`
	URLReported     string     `json:"url_reported" gorm:"size:2048"`
	CountryCode     *string    `json:"country_code" gorm:"size:2"`
	Severity        string     `json:"severity" gorm:"size:16;default:medium"`     // info, low, medium, high, critical
	Status          string     `json:"status" gorm:"size:16;default:pending"`      // pending, verified, rejected
	ResolutionNotes *string    `json:"resolution_notes" gorm:"type:text"`
	ResolvedAt      *time.Time `json:"resolved_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
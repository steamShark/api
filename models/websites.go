package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
Function to be executed before the creation
*/
func (w *Website) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	return nil
}

type Website struct {
	ID                string       `json:"id" gorm:"primaryKey;size:36"`
	URL               string       `json:"url" gorm:"not null;size:255"`
	Domain            string       `json:"domain" gorm:"uniqueIndex;not null;size:255"`
	SSLCertificate    bool         `json:"ssl_certificate" gorm:"default:false"` // if its website, false -> not valid, true -> valid
	DisplayName       *string      `json:"display_name" gorm:"size:255"`         //Display name for website and api
	TLD               string       `json:"tld" gorm:"size:63"`
	Description       *string      `json:"description" gorm:"size:255"`                  //Description for website and api
	Type              string       `json:"type" gorm:"not null;size:16;default:website"` // website, tool, extension
	IsNotTrusted      *bool        `json:"is_not_trusted" gorm:"not null;default:true"`
	IsOfficial        bool         `json:"is_official" gorm:"default:false"`
	SteamLoginPresent bool         `json:"steam_login_present" gorm:"default:false"`
	Verified          bool         `json:"verified" gorm:"default:false;not null"`    //if it was verified by an admin
	RiskScore         float64      `json:"risk_score" gorm:"default:0.0"`             //0 - 100 of the sum of everything, if official this is 100, if not check every var we register about website + occurences
	RiskLevel         string       `json:"risk_level" gorm:"size:16;default:unknown"` // unknown, none (only fo official steam websites), low, medium, high, critical
	Status            string       `json:"status" gorm:"size:16;default:active"`      // active, inactive, blocked, archived
	Notes             *string      `json:"notes" gorm:"type:text"`                    //Internal notes
	Occurrences       []Occurrence `json:"occurrences" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

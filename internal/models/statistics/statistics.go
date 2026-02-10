package model_statistics

import "time"

type Statistics struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	Name      string    `json:"name" gorm:"size:255"`                        //names like total_downloads, total_types, total_inteartions, uptime
	Value     string    `json:"value" gorm:"size:255"`                       //a number, or a percentage
	Type      string    `json:"type" gorm:"not null;size:16;default:number"` //number, percentage
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

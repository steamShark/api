package dtos

type OccurrenceCreationInput struct {
	WebsiteID   string  `json:"website_id" binding:"required"`
	Description *string `json:"description,omitempty"`
	URLReported string  `json:"url_reported" binding:"required"`
	CountryCode *string `json:"country_code,omitempty" binding:"omitempty,len=2"`
	Severity    string  `json:"severity,omitempty" binding:"omitempty,oneof=info low medium high critical"`
}

type OccurrenceUpdateInput struct {
	Description *string `json:"description,omitempty"`
	URLReported *string `json:"url_reported,omitempty"`
	CountryCode *string `json:"country_code,omitempty" binding:"omitempty,len=2"`
	Severity    *string `json:"severity,omitempty" binding:"omitempty,oneof=info low medium high critical"`
	Status      *string `json:"status,omitempty" binding:"omitempty,oneof=pending verified rejected"`
}

type OccurrenceReturnDTO struct {
	ID          string  `json:"id"`
	WebsiteID   string  `json:"website_id"`
	Description *string `json:"description,omitempty"`
	URLReported string  `json:"url_reported"`
	CountryCode *string `json:"country_code,omitempty"`
	Severity    string  `json:"severity"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

package dtos

type OccurrenceCreationInput struct {
	WebsiteID   string  `json:"website_id" binding:"required"`
	Type        string  `json:"type,omitempty" binding:"omitempty,oneof=phishing scam malware impersonation fake_giveaway credential_harvesting other"`
	Source      string  `json:"source,omitempty" binding:"omitempty,oneof=user_report automated_scan admin"`
	Description *string `json:"description,omitempty"`
	URLReported string  `json:"url_reported" binding:"required"`
	CountryCode *string `json:"country_code,omitempty" binding:"omitempty,len=2"`
	Severity    string  `json:"severity,omitempty" binding:"omitempty,oneof=info low medium high critical"`
}

type OccurrenceUpdateInput struct {
	Type            *string `json:"type,omitempty" binding:"omitempty,oneof=phishing scam malware impersonation fake_giveaway credential_harvesting other"`
	Description     *string `json:"description,omitempty"`
	URLReported     *string `json:"url_reported,omitempty"`
	CountryCode     *string `json:"country_code,omitempty" binding:"omitempty,len=2"`
	Severity        *string `json:"severity,omitempty" binding:"omitempty,oneof=info low medium high critical"`
	Status          *string `json:"status,omitempty" binding:"omitempty,oneof=pending verified rejected"`
	ResolutionNotes *string `json:"resolution_notes,omitempty"`
}

type OccurrenceReturnDTO struct {
	ID              string  `json:"id"`
	WebsiteID       string  `json:"website_id"`
	Type            string  `json:"type"`
	Source          string  `json:"source"`
	Description     *string `json:"description,omitempty"`
	URLReported     string  `json:"url_reported"`
	CountryCode     *string `json:"country_code,omitempty"`
	Severity        string  `json:"severity"`
	Status          string  `json:"status"`
	ResolutionNotes *string `json:"resolution_notes,omitempty"`
	ResolvedAt      *string `json:"resolved_at,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

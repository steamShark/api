package dtos

type WebsiteCreationInput struct {
	URL               *string  `json:"url" binding:"required"`
	Domain            string   `json:"domain" binding:"required"`
	DisplayName       *string  `json:"display_name,omitempty"`
	TLD               string   `json:"tld,omitempty"`
	Description       *string  `json:"description,omitempty"`
	Type              *string  `json:"type" binding:"required,oneof=website tool extension"`
	SSLCertificate    *bool    `json:"ssl_certificate,omitempty" binding:"required"`
	IsOfficial        *bool    `json:"is_official,omitempty"`
	SteamLoginPresent *bool    `json:"steam_login_present,omitempty"`
	Status            string   `json:"status,omitempty" binding:"omitempty,oneof=active blocked archived"`
	RiskScore         *float64 `json:"risk_score,omitempty"` /* its going to evaluate if is going to be trusted or not, in the first iteration */
	Notes             *string  `json:"notes,omitempty"`
}

type WebsiteUpdateInput struct {
	DisplayName       *string `json:"display_name,omitempty"`
	Description       *string `json:"description,omitempty"`
	Type              *string `json:"type" binding:"required,oneof=website tool extension"`
	SSLCertificate    *bool   `json:"ssl_certificate,omitempty" binding:"required"`
	IsOfficial        *bool   `json:"is_official,omitempty"`
	SteamLoginPresent *bool   `json:"steam_login_present,omitempty"`
	Status            string  `json:"status,omitempty" binding:"omitempty,oneof=active blocked archived"`
	Notes             *string `json:"notes,omitempty"`
}

type WebsiteReturnDTO struct {
	URL               *string  `json:"url" binding:"required"`
	Domain            string   `json:"domain" binding:"required"`
	DisplayName       *string  `json:"display_name,omitempty"`
	TLD               string   `json:"tld,omitempty"`
	Description       *string  `json:"description,omitempty"`
	Type              *string  `json:"type" binding:"required,oneof=website tool extension"`
	SSLCertificate    *bool    `json:"ssl_certificate,omitempty" binding:"required"`
	IsNotTrusted      *bool    `json:"is_not_trusted,omitempty"`
	IsOfficial        *bool    `json:"is_official,omitempty"`
	SteamLoginPresent *bool    `json:"steam_login_present,omitempty"`
	Verified          *bool    `json:"verified" binding:"required"`
	RiskScore         *float64 `json:"risk_score,omitempty"`
	RiskLevel         string   `json:"risk_level,omitempty" binding:"omitempty,oneof=unknown none low medium high critical"`
	Status            string   `json:"status,omitempty" binding:"omitempty,oneof=active blocked archived"`
}

type ExtensionReturnDTO struct {
	URL         *string `json:"url" binding:"required"`
	Description *string `json:"description,omitempty"`
}

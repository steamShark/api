package helpers

import (
	"steamshark-api/dtos"
	"steamshark-api/models"
)

/*
Convert an Creation Website DTO into model
*/
func ConvertWebsiteDTOModelCreation(in dtos.WebsiteCreationInput) (*models.Website, error) {
	if in.Status == "" {
		in.Status = "active"
	}

	websiteModel := models.Website{
		URL:         *in.URL,
		Domain:      in.Domain,
		TLD:         in.TLD,
		Type:        *in.Type,
		DisplayName: in.DisplayName,
		Description: in.Description,
		RiskScore:   *in.RiskScore,
		Status:      in.Status,
		Notes:       in.Notes,
	}
	if in.SSLCertificate != nil {
		websiteModel.SSLCertificate = *in.SSLCertificate
	}
	if in.IsOfficial != nil {
		websiteModel.IsOfficial = *in.IsOfficial
	}
	if in.SteamLoginPresent != nil {
		websiteModel.SteamLoginPresent = *in.SteamLoginPresent
	}

	return &websiteModel, nil
}

/*
Convert an website update dto into model
*/
func ConvertWebsiteDTOModelUpdate(in dtos.WebsiteUpdateInput) (*models.Website, error) {
	websiteModel := models.Website{
		Type:        *in.Type,
		DisplayName: in.DisplayName,
		Description: in.Description,
		Status:      in.Status,
		Notes:       in.Notes,
	}
	if in.SSLCertificate != nil {
		websiteModel.SSLCertificate = *in.SSLCertificate
	}
	if in.IsOfficial != nil {
		websiteModel.IsOfficial = *in.IsOfficial
	}
	if in.SteamLoginPresent != nil {
		websiteModel.SteamLoginPresent = *in.SteamLoginPresent
	}

	return &websiteModel, nil
}

/*
Convert an website model into return dto
*/
func ConvertWebsiteModelDTOReturn(website models.Website) (*dtos.WebsiteReturnDTO, error) {
	websiteReturnDto := dtos.WebsiteReturnDTO{
		URL:               &website.URL,
		Domain:            website.Domain,
		TLD:               website.TLD,
		Type:              &website.Type,
		SSLCertificate:    &website.SSLCertificate,
		IsNotTrusted:      website.IsNotTrusted,
		IsOfficial:        &website.IsOfficial,
		SteamLoginPresent: &website.SteamLoginPresent,
		RiskScore:         &website.RiskScore,
		DisplayName:       website.DisplayName,
		Description:       website.Description,
		RiskLevel:         website.RiskLevel,
		Status:            website.Status,
	}

	return &websiteReturnDto, nil
}

/*
Convert list of model.websites into get extension dto
*/
func ConvertListWebsiteModelDTOReturnExtension(websites []models.Website) (*[]dtos.ExtensionReturnDTO, error) {
	if websites == nil {
		empty := []dtos.ExtensionReturnDTO{}
		return &empty, nil
	}

	result := make([]dtos.ExtensionReturnDTO, 0, len(websites))

	for _, w := range websites {
		url := w.URL
		description := w.Description

		dto := dtos.ExtensionReturnDTO{
			URL:         &url,
			Description: description,
		}

		result = append(result, dto)
	}

	return &result, nil
}

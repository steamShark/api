package models

// Filter the /websites list
type ListWebsitesFilter struct {
	IsNotTrustedEnabled *bool
	IsNotTrusted        *bool
	Domain              string
	Status              string
	RiskLevel           string
}

type ListWebsitesExtensionFilter struct {
	IsNotTrustedEnabled *bool
	IsNotTrusted        *bool
}

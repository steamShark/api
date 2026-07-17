package utils

import (
	"steamshark-api/internal/models"

	"gorm.io/gorm"
)

var severityWeights = map[string]float64{
	"info":     2,
	"low":      5,
	"medium":   15,
	"high":     25,
	"critical": 40,
}

// RecalculateWebsiteRisk recomputes and persists the risk_score, risk_level and
// is_not_trusted for the given website based on its properties and verified occurrences.
func RecalculateWebsiteRisk(db *gorm.DB, websiteID string) error {
	var website models.Website
	if err := db.First(&website, "id = ?", websiteID).Error; err != nil {
		return err
	}

	// Official sites are always safe
	if website.IsOfficial {
		return db.Model(&website).Updates(map[string]any{
			"risk_score":    0.0,
			"risk_level":    "none",
			"is_not_trusted": false,
		}).Error
	}

	var score float64

	// Base website factors
	if !website.SSLCertificate {
		score += 10
	}
	if website.SteamLoginPresent {
		score += 5
	}

	// Sum weights of all verified occurrences
	var occurrences []models.Occurrence
	db.Where("website_id = ? AND status = ?", websiteID, "verified").Find(&occurrences)
	for _, o := range occurrences {
		score += severityWeights[o.Severity]
	}

	if score > 100 {
		score = 100
	}

	level := riskLevel(score)
	isNotTrusted := score > 20

	return db.Model(&website).Updates(map[string]any{
		"risk_score":     score,
		"risk_level":     level,
		"is_not_trusted": isNotTrusted,
	}).Error
}

func riskLevel(score float64) string {
	switch {
	case score == 0:
		return "unknown"
	case score <= 20:
		return "low"
	case score <= 50:
		return "medium"
	case score <= 80:
		return "high"
	default:
		return "critical"
	}
}

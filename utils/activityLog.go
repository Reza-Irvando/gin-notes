package utils

import (
	"encoding/json"
	"gin-notes/models"

	"github.com/jinzhu/gorm"
)

// Fungsi untuk log activity
func LogActivity(db *gorm.DB, userID uint, action string, entity string, entityID uint, details interface{}) error {
	detailsJSON, _ := json.Marshal(details)

	activity := models.ActivityLog{
		UserID:   userID,
		Action:   action,
		Entity:   entity,
		EntityID: entityID,
		Details:  string(detailsJSON),
	}

	return db.Create(&activity).Error
}

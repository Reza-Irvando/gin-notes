package handlers

import (
	"gin-notes/models"
	"gin-notes/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Mendapatkan activity log user
func GetActivityLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		// Dapatkan parameter pagination
		page := 1
		if p := c.Query("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
				page = parsed
			}
		}

		// Build query
		query := db.Where("user_id = ?", userID.(uint))

		// Dapatkan total count
		var total int64
		if err := query.Model(&models.ActivityLog{}).Count(&total).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to get activity log count")
			return
		}

		// Dapatkan activity logs dengan pagination
		var logs []models.ActivityLog
		offset := (page - 1) * DefaultPageSize
		if err := query.Offset(offset).Limit(DefaultPageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to get activity logs")
			return
		}

		utils.SuccessResponse(c, 200, "Activity logs retrieved successfully", gin.H{
			"data":        logs,
			"total":       total,
			"page":        page,
			"page_size":   DefaultPageSize,
			"total_pages": (total + int64(DefaultPageSize) - 1) / int64(DefaultPageSize),
		})
	}
}

// Mendapatkan activity log untuk entity tertentu
func GetEntityActivityLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		entity := c.Param("entity")
		entityID := c.Param("entityId")

		var logs []models.ActivityLog
		if err := db.Where("user_id = ? AND entity = ? AND entity_id = ?", userID.(uint), entity, entityID).
			Order("created_at DESC").
			Find(&logs).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to get activity logs")
			return
		}

		utils.SuccessResponse(c, 200, "Activity logs retrieved successfully", logs)
	}
}

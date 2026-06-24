package handlers

import (
	"gin-notes/models"
	"gin-notes/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Menambahkan note ke favorit
func AddToFavorite(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		noteID := c.Param("id")
		id, err := strconv.ParseUint(noteID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid note ID")
			return
		}

		// Verify note belongs to user
		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).
			First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		// Check if already favorited
		var existingFavorite models.Favorite
		if !db.Where("note_id = ? AND user_id = ?", uint(id), userID.(uint)).First(&existingFavorite).RecordNotFound() {
			utils.ErrorResponse(c, 400, "Note already in favorites")
			return
		}

		favorite := models.Favorite{
			UserID: userID.(uint),
			NoteID: uint(id),
		}

		if err := db.Create(&favorite).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to add to favorite")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "ADD_FAVORITE", "Note", uint(id), gin.H{})

		utils.SuccessResponse(c, 200, "Note added to favorites successfully", nil)
	}
}

// Menghapus note dari favorit
func RemoveFromFavorite(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		noteID := c.Param("id")
		id, err := strconv.ParseUint(noteID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid note ID")
			return
		}

		// Verify note belongs to user
		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).
			First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		if err := db.Where("note_id = ? AND user_id = ?", uint(id), userID.(uint)).Delete(&models.Favorite{}).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to remove from favorite")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "REMOVE_FAVORITE", "Note", uint(id), gin.H{})

		utils.SuccessResponse(c, 200, "Note removed from favorites successfully", nil)
	}
}

// Mendapatkan semua favorite notes
func GetFavoriteNotes(db *gorm.DB) gin.HandlerFunc {
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
		query := db.Joins("INNER JOIN favorites ON notes.id = favorites.note_id").
			Where("favorites.user_id = ? AND notes.deleted_at IS NULL", userID.(uint))

		// Dapatkan total count
		var total int64
		if err := query.Model(&models.Note{}).Count(&total).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to get favorite notes count")
			return
		}

		// Dapatkan favorite notes dengan pagination
		var notes []models.Note
		offset := (page - 1) * DefaultPageSize
		if err := query.Offset(offset).Limit(DefaultPageSize).Order("notes.created_at DESC").Find(&notes).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to get favorite notes")
			return
		}

		utils.SuccessResponse(c, 200, "Favorite notes retrieved successfully", gin.H{
			"data":        notes,
			"total":       total,
			"page":        page,
			"page_size":   DefaultPageSize,
			"total_pages": (total + int64(DefaultPageSize) - 1) / int64(DefaultPageSize),
		})
	}
}
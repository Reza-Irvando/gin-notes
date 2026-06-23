package handlers

import (
	"gin-notes/models"
	"gin-notes/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Struct untuk request membuat tag
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// Membuat tag baru
func CreateTag(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		var input CreateTagRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, 400, "Invalid input")
			return
		}

		tag := models.Tag{
			UserID: userID.(uint),
			Name:   input.Name,
		}

		if err := db.Create(&tag).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to create tag")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "CREATE", "Tag", tag.ID, input)

		utils.SuccessResponse(c, 201, "Tag created successfully", tag)
	}
}

// Mendapatkan semua tag untuk user
func GetAllTags(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		var tags []models.Tag
		if err := db.Where("user_id = ? AND deleted_at IS NULL", userID.(uint)).
			Order("name ASC").
			Find(&tags).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to get tags")
			return
		}

		utils.SuccessResponse(c, 200, "Tags retrieved successfully", tags)
	}
}

// Struct untuk request update tag
type UpdateTagRequest struct {
	Name string `json:"name"`
}

// Update tag
func UpdateTag(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		tagID := c.Param("id")
		id, err := strconv.ParseUint(tagID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid tag ID")
			return
		}

		var tag models.Tag
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).
			First(&tag).Error; err != nil {
			utils.ErrorResponse(c, 404, "Tag not found")
			return
		}

		var input UpdateTagRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, 400, "Invalid input")
			return
		}

		if input.Name != "" {
			tag.Name = input.Name
		}

		if err := db.Save(&tag).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to update tag")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "UPDATE", "Tag", tag.ID, input)

		utils.SuccessResponse(c, 200, "Tag updated successfully", tag)
	}
}

// Delete tag
func DeleteTag(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		tagID := c.Param("id")
		id, err := strconv.ParseUint(tagID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid tag ID")
			return
		}

		var tag models.Tag
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).
			First(&tag).Error; err != nil {
			utils.ErrorResponse(c, 404, "Tag not found")
			return
		}

		// Soft delete
		if err := db.Model(&tag).Update("deleted_at", time.Now()).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to delete tag")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "DELETE", "Tag", tag.ID, gin.H{})

		utils.SuccessResponse(c, 200, "Tag deleted successfully", nil)
	}
}

// Menambahkan tag ke note
func AddTagToNote(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		noteID := c.Param("noteId")
		nID, err := strconv.ParseUint(noteID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid note ID")
			return
		}

		tagID := c.Param("tagId")
		tID, err := strconv.ParseUint(tagID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid tag ID")
			return
		}

		// Verify note belongs to user
		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(nID), userID.(uint)).
			First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		// Verify tag belongs to user
		var tag models.Tag
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(tID), userID.(uint)).
			First(&tag).Error; err != nil {
			utils.ErrorResponse(c, 404, "Tag not found")
			return
		}

		// Check if tag already added
		var existingNoteTag models.NoteTag
		if db.Where("note_id = ? AND tag_id = ?", uint(nID), uint(tID)).First(&existingNoteTag).RecordNotFound() {
			noteTag := models.NoteTag{
				NoteID: uint(nID),
				TagID:  uint(tID),
			}

			if err := db.Create(&noteTag).Error; err != nil {
				utils.ErrorResponse(c, 500, "Failed to add tag to note")
				return
			}

			// Log activity
			utils.LogActivity(db, userID.(uint), "ADD_TAG", "Note", uint(nID), gin.H{"tag_id": tID})
		}

		utils.SuccessResponse(c, 200, "Tag added to note successfully", nil)
	}
}

// Menghapus tag dari note
func RemoveTagFromNote(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		noteID := c.Param("noteId")
		nID, err := strconv.ParseUint(noteID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid note ID")
			return
		}

		tagID := c.Param("tagId")
		tID, err := strconv.ParseUint(tagID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid tag ID")
			return
		}

		// Verify note belongs to user
		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(nID), userID.(uint)).
			First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		if err := db.Where("note_id = ? AND tag_id = ?", uint(nID), uint(tID)).Delete(&models.NoteTag{}).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to remove tag from note")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "REMOVE_TAG", "Note", uint(nID), gin.H{"tag_id": tID})

		utils.SuccessResponse(c, 200, "Tag removed from note successfully", nil)
	}
}

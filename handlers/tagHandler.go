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

// Struct untuk menangkap Query Parameter (?note_id=X&tag_id=Y)
type TagQuery struct {
	NoteID uint `form:"note_id" binding:"required"`
	TagID  uint `form:"tag_id" binding:"required"`
}

// Menambahkan tag ke note
func AddTagToNote(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		// Menangkap parameter dari URL Query (?note_id=...&tag_id=...)
		var query TagQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			utils.ErrorResponse(c, 400, "Query parameter note_id dan tag_id wajib diisi dengan angka")
			return
		}

		// Verify note belongs to user
		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", query.NoteID, userID.(uint)).
			First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		// Verify tag belongs to user
		var tag models.Tag
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", query.TagID, userID.(uint)).
			First(&tag).Error; err != nil {
			utils.ErrorResponse(c, 404, "Tag not found")
			return
		}

		// Check if tag already added
		var existingNoteTag models.NoteTag
		err := db.Where("note_id = ? AND tag_id = ?", query.NoteID, query.TagID).First(&existingNoteTag).Error
		
		if err != nil { // Jika data belum ada (error record not found), lakukan insert
			noteTag := models.NoteTag{
				NoteID: query.NoteID,
				TagID:  query.TagID,
			}

			if errCreate := db.Create(&noteTag).Error; errCreate != nil {
				utils.ErrorResponse(c, 500, "Failed to add tag to note")
				return
			}

			// Log activity
			utils.LogActivity(db, userID.(uint), "ADD_TAG", "Note", query.NoteID, gin.H{"tag_id": query.TagID})
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

		// Menangkap parameter dari URL Query
		var query TagQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			utils.ErrorResponse(c, 400, "Query parameter note_id dan tag_id wajib diisi dengan angka")
			return
		}

		// Verify note belongs to user
		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", query.NoteID, userID.(uint)).
			First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		// Menghapus relasi tag pada note
		if err := db.Where("note_id = ? AND tag_id = ?", query.NoteID, query.TagID).Delete(&models.NoteTag{}).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to remove tag from note")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "REMOVE_TAG", "Note", query.NoteID, gin.H{"tag_id": query.TagID})

		utils.SuccessResponse(c, 200, "Tag removed from note successfully", nil)
	}
}
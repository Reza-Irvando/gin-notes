package handlers

import (
	"gin-notes/models"
	"gin-notes/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Konstanta untuk pagination
const DefaultPageSize int = 5

// Struct untuk request membuat catatan
type CreateNoteRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Status     string `json:"status"`
	CategoryID *uint  `json:"category_id"`
}

// Membuat catatan baru
func CreateNote(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		var input CreateNoteRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, 400, "Invalid input")
			return
		}

		// Validasi judul dan konten
		if !utils.IsValidTitle(input.Title) {
			utils.ErrorResponse(c, 400, "Title is required")
			return
		}

		if !utils.IsValidContent(input.Content) {
			utils.ErrorResponse(c, 400, "Content is required")
			return
		}

		// Set status default jika tidak diberikan
		if input.Status == "" {
			input.Status = "Active"
		}

		note := models.Note{
			UserID:     userID.(uint),
			Title:      input.Title,
			Content:    input.Content,
			Status:     input.Status,
			CategoryID: input.CategoryID,
		}

		if err := db.Create(&note).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to create note")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "CREATE", "Note", note.ID, input)

		utils.SuccessResponse(c, 201, "Note created successfully", note)
	}
}

// Struct untuk request mengubah catatan
type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

// Mengubah catatan yang sudah ada
func UpdateNote(db *gorm.DB) gin.HandlerFunc {
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

		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		var input UpdateNoteRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, 400, "Invalid input")
			return
		}

		// Update field jika disediakan
		if input.Title != "" {
			note.Title = input.Title
		}
		if input.Content != "" {
			note.Content = input.Content
		}
		if input.Status != "" {
			note.Status = input.Status
		}

		if err := db.Save(&note).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to update note")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "UPDATE", "Note", note.ID, input)

		utils.SuccessResponse(c, 200, "Note updated successfully", note)
	}
}

// Menghapus catatan
func DeleteNote(db *gorm.DB) gin.HandlerFunc {
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

		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		// Soft delete dengan set deleted_at
		if err := db.Model(&note).Update("deleted_at", time.Now()).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to delete note")
			return
		}

		// Log activity
		utils.LogActivity(db, userID.(uint), "DELETE", "Note", note.ID, gin.H{})

		utils.SuccessResponse(c, 200, "Note deleted successfully", nil)
	}
}

// Mendapatkan detail catatan
func GetNoteDetail(db *gorm.DB) gin.HandlerFunc {
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

		var note models.Note
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).First(&note).Error; err != nil {
			utils.ErrorResponse(c, 404, "Note not found")
			return
		}

		utils.SuccessResponse(c, 200, "Note retrieved successfully", note)
	}
}

// Struct untuk parameter pagination
type PaginationParams struct {
	Page     int
	PageSize int
	Search   string
	Status   string
}

// Mendapatkan semua catatan dengan pagination, pencarian, dan filter
func GetAllNotes(db *gorm.DB) gin.HandlerFunc {
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

	// Dapatkan parameter pencarian dan filter
	search := c.Query("search")
	status := c.Query("status")

	// Build query dengan soft delete filter
	query := db.Where("user_id = ? AND deleted_at IS NULL", userID.(uint))

	// Tambahkan filter pencarian
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Tambahkan filter status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Dapatkan total count
	var total int64
	if err := query.Model(&models.Note{}).Count(&total).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to get notes count")
		return
	}

	// Dapatkan catatan dengan pagination menggunakan DefaultPageSize
	var notes []models.Note
	offset := (page - 1) * DefaultPageSize
	if err := query.Offset(offset).Limit(DefaultPageSize).Order("created_at DESC").Find(&notes).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to get notes")
		return
	}

	utils.SuccessResponse(c, 200, "Notes retrieved successfully", gin.H{
		"data":        notes,
		"total":       total,
		"page":        page,
		"page_size":   DefaultPageSize,
		"total_pages": (total + int64(DefaultPageSize) - 1) / int64(DefaultPageSize),
	})
	}
}

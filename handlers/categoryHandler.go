package handlers

import (
	"gin-notes/models"
	"gin-notes/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Struct untuk request membuat category
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

// Membuat category baru
func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		var input CreateCategoryRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, 400, "Invalid input")
			return
		}

		// Set default color jika tidak ada
		if input.Color == "" {
			input.Color = "#3498db"
		}

		category := models.Category{
			UserID: userID.(uint),
			Name:   input.Name,
			Color:  input.Color,
		}

		if err := db.Create(&category).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to create category")
			return
		}
	}
}

// Mendapatkan semua category untuk user
func GetAllCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		var categories []models.Category
		if err := db.Where("user_id = ? AND deleted_at IS NULL", userID.(uint)).
			Order("name ASC").
			Find(&categories).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to get categories")
			return
		}

		utils.SuccessResponse(c, 200, "Categories retrieved successfully", categories)
	}
}

// Struct untuk request update category
type UpdateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Update category
func UpdateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		categoryID := c.Param("id")
		id, err := strconv.ParseUint(categoryID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid category ID")
			return
		}

		var category models.Category
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).
			First(&category).Error; err != nil {
			utils.ErrorResponse(c, 404, "Category not found")
			return
		}

		var input UpdateCategoryRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, 400, "Invalid input")
			return
		}

		if input.Name != "" {
			category.Name = input.Name
		}
		if input.Color != "" {
			category.Color = input.Color
		}

		if err := db.Save(&category).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to update category")
			return
		}
	}
}

// Delete category
func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Unauthorized")
			return
		}

		categoryID := c.Param("id")
		id, err := strconv.ParseUint(categoryID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid category ID")
			return
		}

		var category models.Category
		if err := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", uint(id), userID.(uint)).
			First(&category).Error; err != nil {
			utils.ErrorResponse(c, 404, "Category not found")
			return
		}

		// Soft delete
		if err := db.Model(&category).Update("deleted_at", time.Now()).Error; err != nil {
			utils.ErrorResponse(c, 500, "Failed to delete category")
			return
		}
	}
}

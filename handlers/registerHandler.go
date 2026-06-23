package handlers

import (
	"gin-notes/models"
	"gin-notes/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Struct untuk request pendaftaran pengguna
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Fungsi untuk mendaftarkan pengguna baru
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, 400, "Invalid input")
			return
		}

		// Validasi format email
		if !utils.IsValidEmail(input.Email) {
			utils.ErrorResponse(c, 400, "Invalid email format")
			return
		}

		// Validasi panjang password
		if !utils.IsValidPassword(input.Password) {
			utils.ErrorResponse(c, 400, "Password must be at least 8 characters")
			return
		}

		// Cek apakah email sudah terdaftar
		var existingUser models.User
		if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			utils.ErrorResponse(c, 400, "Email already exists")
			return
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			utils.ErrorResponse(c, 500, "Internal Server Error")
			return
		}

		// Buat pengguna baru
		newUser := models.User{
			Email:    input.Email,
			Password: hashedPassword,
		}

		if err := db.Create(&newUser).Error; err != nil {
			utils.ErrorResponse(c, 500, "Internal Server Error")
			return
		}

		// Generate token untuk pengguna baru setelah pendaftaran
		token, err := CreateToken(newUser.ID)
		if err != nil {
			utils.ErrorResponse(c, 500, "Internal Server Error")
			return
		}

		utils.SuccessResponse(c, 201, "User registered successfully", gin.H{
			"token": token,
			"user": gin.H{
				"id":    newUser.ID,
				"email": newUser.Email,
			},
		})
	}
}
package handlers

import (
	"gin-notes/models"
	"gin-notes/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Struct untuk request login
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Fungsi untuk login pengguna
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, 400, "Invalid input")
			return
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			utils.ErrorResponse(c, 401, "Invalid credentials")
			return
		}

		// Verifikasi password
		if !utils.VerifyPassword(user.Password, input.Password) {
			utils.ErrorResponse(c, 401, "Invalid credentials")
			return
		}

		// Generate token
		token, err := CreateToken(user.ID)
		if err != nil {
			utils.ErrorResponse(c, 500, "Internal Server Error")
			return
		}

		utils.SuccessResponse(c, 200, "Login successful", gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
		})
	}
}
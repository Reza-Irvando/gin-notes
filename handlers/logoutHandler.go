package handlers

import (
	"gin-notes/utils"

	"github.com/gin-gonic/gin"
)

// Fungsi untuk logout pengguna
// Memerlukan autentikasi - pengguna harus memiliki token JWT yang valid
// Dalam sistem JWT yang stateless, logout dilakukan di sisi klien dengan menghapus token
// Handler ini mengonfirmasi bahwa pengguna yang terautentikasi telah logout
func Logout(c *gin.Context) {
	// Mengambil user_id dari context yang di-set oleh AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, 401, "User tidak terautentikasi")
		return
	}

	utils.SuccessResponse(c, 200, "Logout berhasil", gin.H{
		"user_id": userID,
	})
}
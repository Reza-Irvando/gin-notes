package handlers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Kunci rahasia untuk JWT - ubah dengan kunci yang aman
var jwtKey = []byte("your-secret-key")

// Fungsi untuk membuat token JWT
func CreateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
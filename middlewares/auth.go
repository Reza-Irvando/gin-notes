package middlewares

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gin-notes/utils"
)

// Kunci rahasia untuk JWT - ubah dengan kunci yang aman
var jwtKey = []byte("your-secret-key")

// Middleware untuk autentikasi JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utils.ErrorResponse(c, 401, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			utils.ErrorResponse(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)

		c.Next()
	}
}

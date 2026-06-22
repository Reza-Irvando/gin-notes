package utils

import (
	"regexp"
)

// Fungsi untuk validasi format email
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// Fungsi untuk validasi panjang password
func IsValidPassword(password string) bool {
	return len(password) >= 8
}

// Fungsi untuk validasi judul tidak kosong
func IsValidTitle(title string) bool {
	return len(title) > 0
}

// Fungsi untuk validasi konten tidak kosong
func IsValidContent(content string) bool {
	return len(content) > 0
}

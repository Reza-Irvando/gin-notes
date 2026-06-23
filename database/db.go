package database

import (
	"gin-notes/models"
	"gin-notes/utils"
	"time"

	"github.com/jinzhu/gorm"
)

// Fungsi untuk melakukan migrasi database
func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Note{},
		&models.Category{},
		&models.Tag{},
		&models.NoteTag{},
		&models.Favorite{},
		&models.ActivityLog{},
	)
}

// Fungsi untuk seed data user
func SeedUsers(db *gorm.DB) error {
	users := []models.User{
		{
			Email:    "user1@example.com",
			Password: "password123",
		},
		{
			Email:    "user2@example.com",
			Password: "password123",
		},
		{
			Email:    "user3@example.com",
			Password: "password123",
		},
	}

	for i := range users {
		hashedPassword, err := utils.HashPassword(users[i].Password)
		if err != nil {
			return err
		}
		users[i].Password = hashedPassword

		// Cek apakah user sudah ada
		var existingUser models.User
		if db.Where("email = ?", users[i].Email).First(&existingUser).RecordNotFound() {
			if err := db.Create(&users[i]).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// Fungsi untuk seed data notes
func SeedNotes(db *gorm.DB) error {
	// Ambil semua user
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	noteTitles := []string{
		"Rapat Pagi Tim",
		"Review Code",
		"Update Database",
		"Fixing Bug Login",
		"Dokumentasi API",
		"Testing Fitur Baru",
		"Diskusi dengan PM",
		"Deployment Production",
		"Optimasi Query",
		"Backup Data Server",
	}

	noteContents := []string{
		"Membahas progress sprint mingguan. Tim backend selesaikan API authentication. Frontend siap integrate dengan endpoint baru.",
		"Review pull request dari 3 developer. Menemukan issue pada error handling. Suggest refactor code untuk readability lebih baik.",
		"Update schema database untuk support fitur baru. Tambah tabel relation, optimize index untuk query performa.",
		"Temukan bug di login handler. User tidak bisa login dengan password benar. Fixed issue pada password verification logic.",
		"Tulis dokumentasi lengkap untuk API endpoints. Include request/response examples, error codes, dan authentication requirements.",
		"Test semua fitur aplikasi di development environment. Cek validasi input, edge cases, dan response time setiap endpoint.",
		"Meeting dengan Product Manager bahas roadmap quarter depan. Prioritas fitur: export notes, collaboration sharing notes.",
		"Deploy aplikasi ke production server. Update environment variables, run migration database, verify health check semua service.",
		"Analisis slow query di database. Optimize dengan index strategy, rewrite query menggunakan JOIN lebih efisien.",
		"Jalankan backup full database dan verify integritas data. Prepare disaster recovery plan untuk business continuity.",
	}

	// Buat 10 notes untuk setiap user
	for _, user := range users {
		for i := 0; i < 10; i++ {
			note := models.Note{
				UserID:    user.ID,
				Title:     noteTitles[i],
				Content:   noteContents[i],
				Status:    "Active",
				CreatedAt: time.Now().Add(-time.Hour * time.Duration(i)),
				UpdatedAt: time.Now().Add(-time.Hour * time.Duration(i)),
				DeletedAt: nil,
			}

			// Cek apakah note sudah ada
			var existingNote models.Note
			if db.Where("user_id = ? AND title = ?", user.ID, note.Title).First(&existingNote).RecordNotFound() {
				if err := db.Create(&note).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Fungsi untuk menjalankan semua seeding
func Seed(db *gorm.DB) error {
	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedNotes(db); err != nil {
		return err
	}

	return nil
}
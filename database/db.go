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

// Fungsi untuk seed data categories
func SeedCategories(db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	categoryColors := map[string]string{
		"Work":       "#FF5733",
		"Personal":   "#3498DB",
		"Important":  "#E74C3C",
		"Ideas":      "#F39C12",
	}

	for _, user := range users {
		for name, color := range categoryColors {
			category := models.Category{
				UserID: user.ID,
				Name:   name,
				Color:  color,
			}

			var existingCategory models.Category
			if db.Where("user_id = ? AND name = ?", user.ID, name).First(&existingCategory).RecordNotFound() {
				if err := db.Create(&category).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Fungsi untuk seed data tags
func SeedTags(db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	tagNames := []string{
		"Urgent",
		"Completed",
		"In Progress",
		"Blocked",
		"Review Needed",
		"High Priority",
		"Low Priority",
		"Bug Fix",
		"Feature Request",
		"Documentation",
	}

	for _, user := range users {
		for _, tagName := range tagNames {
			tag := models.Tag{
				UserID: user.ID,
				Name:   tagName,
			}

			var existingTag models.Tag
			if db.Where("user_id = ? AND name = ?", user.ID, tagName).First(&existingTag).RecordNotFound() {
				if err := db.Create(&tag).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Fungsi untuk seed data note-tag relationships
func SeedNoteTags(db *gorm.DB) error {
	var notes []models.Note
	if err := db.Where("deleted_at IS NULL").Find(&notes).Error; err != nil {
		return err
	}

	for _, note := range notes {
		// Ambil beberapa random tags untuk user ini
		var tags []models.Tag
		if err := db.Where("user_id = ? AND deleted_at IS NULL", note.UserID).
			Limit(3).
			Find(&tags).Error; err != nil {
			return err
		}

		for _, tag := range tags {
			noteTag := models.NoteTag{
				NoteID: note.ID,
				TagID:  tag.ID,
			}

			var existingNoteTag models.NoteTag
			if db.Where("note_id = ? AND tag_id = ?", note.ID, tag.ID).First(&existingNoteTag).RecordNotFound() {
				if err := db.Create(&noteTag).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Fungsi untuk seed data favorites
func SeedFavorites(db *gorm.DB) error {
	var notes []models.Note
	if err := db.Where("deleted_at IS NULL").Find(&notes).Error; err != nil {
		return err
	}

	// Tandai 30% dari notes sebagai favorit (cukup random)
	count := 0
	for _, note := range notes {
		if count%3 == 0 {
			favorite := models.Favorite{
				UserID: note.UserID,
				NoteID: note.ID,
			}

			var existingFavorite models.Favorite
			if db.Where("user_id = ? AND note_id = ?", note.UserID, note.ID).First(&existingFavorite).RecordNotFound() {
				if err := db.Create(&favorite).Error; err != nil {
					return err
				}
			}
		}
		count++
	}

	return nil
}

// Fungsi untuk update notes dengan category
func UpdateNotesWithCategories(db *gorm.DB) error {
	var notes []models.Note
	if err := db.Where("deleted_at IS NULL").Find(&notes).Error; err != nil {
		return err
	}

	categoryIndex := 0
	for _, note := range notes {
		// Ambil categories user ini
		var categories []models.Category
		if err := db.Where("user_id = ? AND deleted_at IS NULL", note.UserID).
			Find(&categories).Error; err != nil {
			return err
		}

		if len(categories) > 0 {
			categoryID := categories[categoryIndex%len(categories)].ID
			if err := db.Model(&note).Update("category_id", categoryID).Error; err != nil {
				return err
			}
			categoryIndex++
		}
	}

	return nil
}

// Fungsi untuk seed data activity logs
func SeedActivityLogs(db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		// Log untuk setiap note yang dibuat
		var notes []models.Note
		if err := db.Where("user_id = ? AND deleted_at IS NULL", user.ID).Find(&notes).Error; err != nil {
			return err
		}

		for _, note := range notes {
			// CREATE activity
			activityCreate := models.ActivityLog{
				UserID:   user.ID,
				Action:   "CREATE",
				Entity:   "Note",
				EntityID: note.ID,
				Details:  `{"title":"` + note.Title + `"}`,
			}

			var existingActivity models.ActivityLog
			if db.Where("user_id = ? AND action = ? AND entity = ? AND entity_id = ?", 
				user.ID, "CREATE", "Note", note.ID).First(&existingActivity).RecordNotFound() {
				if err := db.Create(&activityCreate).Error; err != nil {
					return err
				}
			}
		}

		// Log untuk setiap category yang dibuat
		var categories []models.Category
		if err := db.Where("user_id = ? AND deleted_at IS NULL", user.ID).Find(&categories).Error; err != nil {
			return err
		}

		for _, category := range categories {
			activityCategory := models.ActivityLog{
				UserID:   user.ID,
				Action:   "CREATE",
				Entity:   "Category",
				EntityID: category.ID,
				Details:  `{"name":"` + category.Name + `"}`,
			}

			var existingActivity models.ActivityLog
			if db.Where("user_id = ? AND action = ? AND entity = ? AND entity_id = ?", 
				user.ID, "CREATE", "Category", category.ID).First(&existingActivity).RecordNotFound() {
				if err := db.Create(&activityCategory).Error; err != nil {
					return err
				}
			}
		}

		// Log untuk setiap tag yang dibuat
		var tags []models.Tag
		if err := db.Where("user_id = ? AND deleted_at IS NULL", user.ID).Find(&tags).Error; err != nil {
			return err
		}

		for _, tag := range tags {
			activityTag := models.ActivityLog{
				UserID:   user.ID,
				Action:   "CREATE",
				Entity:   "Tag",
				EntityID: tag.ID,
				Details:  `{"name":"` + tag.Name + `"}`,
			}

			var existingActivity models.ActivityLog
			if db.Where("user_id = ? AND action = ? AND entity = ? AND entity_id = ?", 
				user.ID, "CREATE", "Tag", tag.ID).First(&existingActivity).RecordNotFound() {
				if err := db.Create(&activityTag).Error; err != nil {
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

	if err := SeedCategories(db); err != nil {
		return err
	}

	if err := SeedTags(db); err != nil {
		return err
	}

	if err := SeedNoteTags(db); err != nil {
		return err
	}

	if err := SeedFavorites(db); err != nil {
		return err
	}

	if err := UpdateNotesWithCategories(db); err != nil {
		return err
	}

	if err := SeedActivityLogs(db); err != nil {
		return err
	}

	return nil
}
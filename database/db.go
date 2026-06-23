package database

import (
	"gin-notes/models"
	"github.com/jinzhu/gorm"
)

// Fungsi untuk melakukan migrasi database
func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Note{},
	)
}
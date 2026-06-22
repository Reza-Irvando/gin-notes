package database

import (
	"gin-notes/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Note{},
	)
}
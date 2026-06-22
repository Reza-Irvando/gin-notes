package database

import (
	"gin-notes/models"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Note{},
	)
}
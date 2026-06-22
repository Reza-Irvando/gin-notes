package configs

import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/go-notes"))
	if err != nil {
		return nil, err
	}
	DB = db
	return DB, nil
}
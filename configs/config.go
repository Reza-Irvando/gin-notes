package configs

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Variabel global untuk database
var DB *gorm.DB

// Fungsi untuk inisialisasi koneksi database
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/go-notes?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	DB = db
	return DB, nil
}
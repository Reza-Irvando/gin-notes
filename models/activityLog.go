package models

import "time"

// Model ActivityLog untuk mencatat semua aktivitas pengguna
type ActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Action    string    `gorm:"not null" json:"action"`
	Entity    string    `gorm:"not null" json:"entity"` // Note, Category, Tag, etc
	EntityID  uint      `gorm:"not null" json:"entity_id"`
	Details   string    `gorm:"type:text" json:"details"` // JSON details
	IPAddress string    `gorm:"default:null" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

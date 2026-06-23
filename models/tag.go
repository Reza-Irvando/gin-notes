package models

import "time"

// Model Tag
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

// Model NoteTag - junction table untuk relasi many-to-many antara Note dan Tag
type NoteTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NoteID    uint      `gorm:"not null;index" json:"note_id"`
	TagID     uint      `gorm:"not null;index" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

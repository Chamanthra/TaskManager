package models

import "time"

type File struct {
	ID         uint      `gorm:"primaryKey"`
	FilePath   string    `gorm:"not null"`
	FileName   string    `gorm:"not null"`
	FileType   string    `gorm:"size:100"`
	FileSize   int64     `gorm:"default:0"`
	TaskID     uint      `gorm:"index"`
	UserID     uint      `gorm:"index"`
	UploadedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	Task Task `gorm:"foreignKey:TaskID"`
	User User `gorm:"foreignKey:UserID"`
}

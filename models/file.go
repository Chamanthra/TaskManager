package models

import "time"

type File struct {
	ID         uint   `gorm:"primaryKey"`
	FilePath   string `gorm:"not null"`
	FileType   string
	UploadedAt time.Time
}

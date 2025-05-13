// In models/comment.go
package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"not null;size:1000"`
	TaskID    uint      `gorm:"index"`
	UserID    uint      `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	Task Task `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID"`
}

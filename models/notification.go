// In models/notification.go
package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey"`
	Message   string    `gorm:"not null;size:1000"`
	Status    string    `gorm:"type:enum('unread','read');default:'unread'"`
	UserID    uint      `gorm:"index"`
	TaskID    uint      `gorm:"index"`
	Type      string    `gorm:"type:enum('comment','status_change','file_upload','due_date');not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
	Task Task `gorm:"foreignKey:TaskID"`
}

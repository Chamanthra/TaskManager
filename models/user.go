package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserName  string    `json:"user_name" gorm:"unique;not null;size:50"`
	Email     string    `json:"email" gorm:"unique;not null;size:100"`
	Password  string    `json:"-" gorm:"not null"` // - means don't include in JSON
	FirstName string    `json:"first_name" gorm:"size:50"`
	LastName  string    `json:"last_name" gorm:"size:50"`
	RoleID    uint      `json:"-" gorm:"index"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Role          Role           `gorm:"foreignKey:RoleID"`
	Tasks         []Task         `gorm:"foreignKey:UserID"`
	Comments      []Comment      `gorm:"foreignKey:UserID"`
	Files         []File         `gorm:"foreignKey:UserID"`
	Notifications []Notification `gorm:"foreignKey:UserID"`
}

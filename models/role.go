package models

type Role struct {
	ID         uint   `gorm:"primaryKey"`
	Role       string `gorm:"unique"`
	Permission string
}

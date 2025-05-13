package models

type Role struct {
	ID                uint   `gorm:"primaryKey"`
	Role              string `gorm:"unique;not null"`
	Permission        string
	Name              string
	Description       string
	CanManageTasks    bool `gorm:"default:false"`
	CanManageUsers    bool `gorm:"default:false"`
	CanManageComments bool `gorm:"default:false"`
	CanManageFiles    bool `gorm:"default:false"`
	IsAdmin           bool `gorm:"default:false"`
}

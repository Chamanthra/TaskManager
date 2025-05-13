package migrations

import (
	"github.com/Chamanthra/TaskManager/models"
	"gorm.io/gorm"
)

func InitRoles(db *gorm.DB) error {
	defaultRoles := []models.Role{
		{
			Role:              "user",
			Permission:        "basic",
			Name:              "Regular User",
			Description:       "Can manage own tasks",
			CanManageTasks:    true,
			CanManageUsers:    false,
			CanManageComments: true,
			CanManageFiles:    true,
			IsAdmin:           false,
		},
		{
			Role:              "admin",
			Permission:        "all",
			Name:              "Administrator",
			Description:       "Can manage all system resources",
			CanManageTasks:    true,
			CanManageUsers:    true,
			CanManageComments: true,
			CanManageFiles:    true,
			IsAdmin:           true,
		},
	}

	for _, role := range defaultRoles {
		if err := db.FirstOrCreate(&role, models.Role{Role: role.Role}).Error; err != nil {
			return err
		}
	}
	return nil
}

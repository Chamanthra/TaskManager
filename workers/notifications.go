package workers

import (
	"fmt"
	"time"

	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/models"
)

func StartNotificationWorker() {
	ticker := time.NewTicker(24 * time.Hour) // Check daily
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkDueTasks()
		}
	}
}

func checkDueTasks() {
	var tasks []models.Task
	oneWeekFromNow := time.Now().Add(7 * 24 * time.Hour)

	config.DB.Where("due_date BETWEEN ? AND ? AND notified = ?",
		time.Now(), oneWeekFromNow, false).Find(&tasks)

	for _, task := range tasks {
		notification := models.Notification{
			Message: fmt.Sprintf("Task '%s' is due in less than a week!", task.Title),
			UserID:  task.UserID,
			TaskID:  task.ID,
			Type:    "due_date",
		}

		if err := config.DB.Create(&notification).Error; err == nil {
			// Mark as notified if creation succeeded
			config.DB.Model(&task).Update("notified", true)
		}
	}
}

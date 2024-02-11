package database

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/models"
	"fmt"
	"gorm.io/gorm"
	"slices"
)

// Добавляет задание в БД
func AddTask(ex, req_id, user_id string, time int) {
	task := models.Task{
		Model:      gorm.Model{},
		Expression: ex,
		Req_id:     req_id,
		User_id:    user_id,
		Status:     false,
		ToDoTime:   time,
		Res:        "",
		Err:        "",
	}

	res := db.Create(&task)
	if res.Error != nil {
		config.Log.Warn(res.Error)
	}
}

// Обновляют данные в задание
func UpdateTask(reqId, user_id, res string, status bool, err string) {
	var task models.Task
	if err := db.First(&task, "req_id = ?", reqId).Error; err != nil {
		config.Log.Error(err)
		return
	}

	if status != false {
		task.Status = status
	}

	if user_id != "" {
		tmp := task.GetUserIDs()
		if !slices.Contains(tmp, user_id) {
			tmp = append(tmp, user_id)
		}
		task.SetUserIDs(tmp)
	}

	if res != "" {
		task.Res = res
	}

	if err != "" {
		task.Err = err
	}

	if err := db.Save(&task).Error; err != nil {
		config.Log.Error(err)
		return
	}

}

// Выдает задание по ID
func GetTask(reqId string) (models.Task, bool) {
	var task models.Task
	if err := db.First(&task, "req_id = ?", reqId).Error; err != nil {
		config.Log.WithField("DB", "Не удалось найти значение").Warn(err)
		return task, false
	}
	return task, true
}

// Выдает все задания пользователя
func GetAllUserTask(user_id string) ([]models.Task, bool) {
	var tasks []models.Task
	err := db.Model(&models.Task{}).Where("user_id LIKE ?", fmt.Sprintf("%%%s%%", user_id)).Find(&tasks).Error
	if err != nil {
		config.Log.Error(err)
		return []models.Task{}, false
	}
	return tasks, true
}

// Выдает все задания
func GetAllTask() ([]models.Task, bool) {
	var task []models.Task
	if err := db.Find(&task).Error; err != nil {
		config.Log.WithField("DB", "Не удалось найти значение").Warn(err)
		return task, false
	}
	return task, true
}

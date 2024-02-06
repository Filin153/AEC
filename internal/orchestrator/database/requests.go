package database

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/models"
	"gorm.io/gorm"
)

func AddTask(ex, req_id, user_id string) {
	task := models.Task{
		Model:      gorm.Model{},
		Expression: ex,
		Req_id:     req_id,
		User_id:    user_id,
		Status:     false,
		Res:        "",
		Err:        "",
	}

	res := db.Create(&task)
	if res.Error != nil {
		config.Log.Error(res.Error)
	}
}

func UpdateTask(id string, res string, status bool, err string) {
	model := db.Model(&models.Task{}).Where("req_id = ?", id)

	if status != false {
		model.Update("status", status)

	}

	if res != "" {
		model.Update("res", res)
	}

	if err != "" {
		model.Update("err", err)
	}

}

func GetTask(req_id string) models.Task {
	var task models.Task
	db.Model(&models.Task{}).Where("req_id = ?", req_id).First(&task)
	return task
}

func GetAllTask(user_id string) []models.Task {
	var tasks []models.Task
	db.Model(&models.Task{}).Where("user_id = ?", user_id).Find(&tasks)
	return tasks
}
func FindUser(user_id string) bool {
	var task models.Task
	db.Model(&models.Task{}).Where("user_id = ?", user_id).First(&task)
	if task.User_id != "" {
		return true
	}

	return false
}

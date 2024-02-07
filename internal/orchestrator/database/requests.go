package database

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/models"
	"gorm.io/gorm"
	"slices"
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
		config.Log.Warn(res.Error)
		UpdateTask(req_id, user_id, "", false, "")
	}
}

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

//func FindUser(user_id string) bool {
//	var task models.Task
//	db.Model(&models.Task{}).Where("user_id = ?", user_id).First(&task)
//	if task.User_id != "" {
//		return true
//	}
//
//	return false
//}

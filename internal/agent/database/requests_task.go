package database

import (
	"AEC/internal/agent/config"
	"AEC/internal/agent/models"
	"slices"
)

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

package database

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/models"
	"gorm.io/gorm"
)

func AddCalRes(id, ex string) {
	res := models.CalRes{
		Model:      gorm.Model{},
		RId:        id,
		Expression: ex,
		Res:        "",
		Err:        "",
	}

	info := db.Create(&res)
	if info.Error != nil {
		config.Log.Warn(info.Error)
	}
}

func GetCalRes(id string) (models.CalRes, bool) {
	res := models.CalRes{}

	if err := db.First(&res, "r_id = ?", id).Error; err != nil {
		config.Log.WithField("DB", "Не удалось найти значение").Warn(err)
		return res, false
	}

	return res, true
}

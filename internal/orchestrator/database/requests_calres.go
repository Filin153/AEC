package database

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/models"
)

// Выдает выражение по ID
func GetCalRes(id string) (models.CalRes, bool) {
	res := models.CalRes{}

	if err := db.First(&res, "r_id = ?", id).Error; err != nil {
		config.Log.WithField("DB", "Не удалось найти значение").Warn(err)
		return res, false
	}

	return res, true
}

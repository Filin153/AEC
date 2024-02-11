package database

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Подключается к БД
// Делает миграции
func init() {

	dbUrl := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", config.Conf.DbHost, config.Conf.DbUsername, config.Conf.DbName, config.Conf.DbPassword)

	conn, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		config.Log.Error(err)
	}

	db = conn
	config.Log.Info("connected to DB successful")

	err = db.AutoMigrate(&models.Task{}, &models.CalRes{})
	if err != nil {
		config.Log.Error(err)
	}

	config.Log.Info("AutoMigrate successful")

}

package database

import (
	conf "AEC/internal/orchestrator/config"
	mod "AEC/internal/orchestrator/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func init() {

	err := godotenv.Load()
	if err != nil {
		conf.Log.Error(err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUrl := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		conf.Log.Error(err)
	}

	db = conn
	conf.Log.Info("connected to DB successful")

	err = db.AutoMigrate(&mod.Task{})
	if err != nil {
		conf.Log.Error(err)
	}

	conf.Log.Info("AutoMigrate successful")

}

package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

type config struct {
	Redis_host string
	Redis_port string
	Port       string
}

var Conf config

func init() {
	relativePath := "internal/orchestrator/config/.env"
	currentDir, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("Ошибка при получении текущей директории:", err)
		return
	}
	absolutePath := filepath.Join(currentDir, relativePath)

	err = godotenv.Load(absolutePath)
	if err != nil {
		Log.Error(err)
		panic(err)
	}

	Conf.Redis_host = os.Getenv("redis_host")
	Conf.Redis_port = os.Getenv("redis_port")
	Conf.Port = os.Getenv("port")

	fmt.Println(Conf)
}

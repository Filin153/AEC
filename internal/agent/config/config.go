package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type config struct {
	worker     int
	redis_host string
	redis_port string
}

var Conf config

func init() {

	err := godotenv.Load()
	if err != nil {
		Log.Error(err)
		panic(err)
	}

	Conf.worker, _ = strconv.Atoi(os.Getenv("worker"))
	Conf.redis_host = os.Getenv("redis_host")
	Conf.redis_port = os.Getenv("redis_port")

}

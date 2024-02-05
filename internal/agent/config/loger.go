package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func init() {
	Log.SetLevel(logrus.DebugLevel)

	Log.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Info("Не удалось открыть файл логов, используется стандартный stderr")
	}
}

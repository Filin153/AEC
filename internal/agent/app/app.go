package app

import (
	"AEC/internal/agent/config"
	"AEC/internal/agent/services"
	"AEC/internal/agent/transport"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Run() {
	defer config.RedisClientQ.Close()

	taskChan := make(chan []byte)

	go services.StartWorkers(config.Conf.Worker, taskChan)
	go services.AddTask(taskChan)
	go services.PING(config.Conf.Connect_to, config.Conf.Сonnect_path, config.Conf.I_host)

	router := mux.NewRouter()

	router.HandleFunc("/", transport.AddCal).Methods("POST")

	err := http.ListenAndServe(":"+config.Conf.Port, router)
	fmt.Println("Server start - " + config.Conf.Port)

	if err != nil {
		config.Log.Error(err)
		panic(err)
	}
}

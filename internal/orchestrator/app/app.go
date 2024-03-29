package app

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/docs"
	"AEC/internal/orchestrator/services"
	"AEC/internal/orchestrator/transport"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Run() {
	defer config.RedisClient.Close()

	go docs.Swag()
	go services.CheckServer()
	go services.CheckNoReadyEx()

	router := mux.NewRouter()

	router.HandleFunc("/server/newcon", transport.Connect).Methods("POST")
	router.HandleFunc("/server/all", transport.AllServ).Methods("GET")
	router.HandleFunc("/server/add/{id}/{add}", transport.AddWorkerFor).Methods("POST")
	router.HandleFunc("/server/del/{id}", transport.DeleteServer).Methods("DELETE")
	router.HandleFunc("/", transport.Calc).Methods("POST")
	router.HandleFunc("/task/{id}", transport.GetOneTask).Methods("GET")
	router.HandleFunc("/user/{id}", transport.GetAllTaskFromUser).Methods("GET")

	err := http.ListenAndServe(":"+config.Conf.Port, router)
	fmt.Println("Server start - " + config.Conf.Port)

	if err != nil {
		panic(err)
	}
}

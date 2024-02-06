package app

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/services"
	"AEC/internal/orchestrator/transport"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Run() {
	go services.CheckServer()

	router := mux.NewRouter()

	router.HandleFunc("/server/newcon", transport.Connect).Methods("POST")
	router.HandleFunc("/server/all", transport.AllServ).Methods("GET")
	router.HandleFunc("/server/del", transport.DeleteServer).Methods("DELETE")

	err := http.ListenAndServe(":"+config.Conf.Port, router)
	fmt.Println("Server start - " + config.Conf.Port)

	if err != nil {
		panic(err)
	}
}

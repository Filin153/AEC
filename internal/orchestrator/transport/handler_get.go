package transport

import (
	"AEC/internal/orchestrator/database"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func GetOneTask(w http.ResponseWriter, r *http.Request) {
	var a answer

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		a.Err = errors.New("id не обнаружен")
		a.Info = "/task/{id}, то как должен выглядеть путь"
		w.WriteHeader(400)
		jsonResp, _ := json.Marshal(a)
		w.Write(jsonResp)
	}

	if task, ok := database.GetTask(id); ok {
		a.Data = task
		w.WriteHeader(200)
		jsonResp, _ := json.Marshal(a)
		w.Write(jsonResp)
		return
	}

	a.Info = "Не удалось найти запись"
	w.WriteHeader(400)
	jsonResp, _ := json.Marshal(a)
	w.Write(jsonResp)

}

func GetAllTaskFromUser(w http.ResponseWriter, r *http.Request) {
	var a answer

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		a.Err = errors.New("id не обнаружен")
		a.Info = "/user/{id}, то как должен выглядеть путь"
		w.WriteHeader(400)
		jsonResp, _ := json.Marshal(a)
		w.Write(jsonResp)
	}

	if task, ok := database.GetAllUserTask(id); ok {
		a.Data = task
		w.WriteHeader(200)
		jsonResp, _ := json.Marshal(a)
		w.Write(jsonResp)
		return
	}

	a.Info = "Не удалось найти записи"
	w.WriteHeader(400)
	jsonResp, _ := json.Marshal(a)
	w.Write(jsonResp)

}

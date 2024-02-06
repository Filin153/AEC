package transport

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/services"
	"encoding/json"
	"net/http"
)

type answer struct {
	Err  error       `json:"err"`
	Data interface{} `json:"data"`
	Info string      `json:"info"`
}

func Connect(w http.ResponseWriter, r *http.Request) {
	fromURL := services.GetClientIP(r)
	a := answer{}
	err := services.AddServer(services.HashSome(fromURL), fromURL)
	if err != nil {
		w.WriteHeader(400)
		a.Err = err
		data, _ := json.Marshal(a)
		w.Write(data)
		return
	}

	data, _ := json.Marshal(a)
	w.Write(data)
}

func AllServ(w http.ResponseWriter, r *http.Request) {
	res := services.AllServer()

	a := answer{
		Err:  nil,
		Data: res,
	}

	data, _ := json.Marshal(a)
	w.Write(data)
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	a := answer{}
	var serv services.Server

	err := json.NewDecoder(r.Body).Decode(&serv)
	if err != nil {
		config.Log.WithField("err", "Не удалось декодировать JSON").Error(err)
		w.WriteHeader(400)
		a.Err = err
		a.Info = "Не удалось декадировать JSON"
		data, _ := json.Marshal(a)
		w.Write(data)
		return
	}

	a.Data = serv.URL

	err = services.RemoveServerFromRedis(services.HashSome(serv.URL))
	if err != nil {
		config.Log.WithField("err", "Не удалось удалить").Error(err)
		w.WriteHeader(400)
		a.Err = err
		a.Info = "Не удалось удалить"
		data, _ := json.Marshal(a)
		w.Write(data)
		return
	}

	a.Info = "Successful delete"

	data, _ := json.Marshal(a)
	w.Write(data)
}

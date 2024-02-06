package transport

import (
	"AEC/internal/orchestrator/database"
	"AEC/internal/orchestrator/services"
	"encoding/json"
	"net/http"
	"time"
)

func Calc(w http.ResponseWriter, r *http.Request) {
	a := answer{}
	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		a.Err = err
		a.Info = "Не удалось декодировать JSON"
		w.WriteHeader(400)
		jsonResp, _ := json.Marshal(a)
		w.Write(jsonResp)
		return
	}

	var userId string
	reqId := services.HashSome(data["task"])

	if data["user_id"] == "" {
		userId = services.HashSome(time.Now().String())
	} else {
		userId = data["user_id"]
	}

	go services.Direct(data["task"], reqId)
	go database.AddTask(data["task"], reqId, userId)

	w.Write([]byte(userId))
}

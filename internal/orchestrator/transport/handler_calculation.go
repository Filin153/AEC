package transport

import (
	"AEC/internal/orchestrator/database"
	"AEC/internal/orchestrator/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	add, _ := strconv.Atoi(fmt.Sprintf("%v", data["add_time"]))
	sub, _ := strconv.Atoi(fmt.Sprintf("%v", data["sub_time"]))
	mult, _ := strconv.Atoi(fmt.Sprintf("%v", data["mult_time"]))
	dev, _ := strconv.Atoi(fmt.Sprintf("%v", data["dev_time"]))

	if _, ok := database.GetTask(reqId); !ok {
		go services.Direct(data["task"], reqId, add, sub, mult, dev)
		go database.AddTask(data["task"], reqId, userId)
	} else {
		go database.UpdateTask(reqId, userId, "", false, "", -1)
	}

	a.Data = map[string]string{
		"reqID":  reqId,
		"userID": userId,
	}

	jsonResp, _ := json.Marshal(a)
	w.Write(jsonResp)
}

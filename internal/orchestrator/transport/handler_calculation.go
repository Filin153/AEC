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

type TaskType struct {
	Task     string `json:"task"`
	UserId   string `json:"user_id"`
	AddTime  string `json:"add_time"`
	SubTime  string `json:"sub_time"`
	MultTime string `json:"mult_time"`
	DevTime  string `json:"dev_time"`
}

func Calc(w http.ResponseWriter, r *http.Request) {
	a := Answer{}
	var data TaskType

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
	reqId := services.HashSome(data.Task)

	if data.UserId == "" {
		userId = services.HashSome(time.Now().String())
	} else {
		userId = data.UserId
	}

	add, _ := strconv.Atoi(fmt.Sprintf("%v", data.AddTime))
	sub, _ := strconv.Atoi(fmt.Sprintf("%v", data.SubTime))
	mult, _ := strconv.Atoi(fmt.Sprintf("%v", data.MultTime))
	dev, _ := strconv.Atoi(fmt.Sprintf("%v", data.DevTime))
	waitTime := services.GetWaitTime(data.Task, add, sub, mult, dev)

	if _, ok := database.GetTask(reqId); !ok {
		go services.Direct(data.Task, reqId, add, sub, mult, dev)
		go database.AddTask(data.Task, reqId, userId, int(waitTime.Seconds()))
	} else {
		go database.UpdateTask(reqId, userId, "", false, "")
	}

	a.Data = map[string]string{
		"reqID":  reqId,
		"userID": userId,
	}

	jsonResp, _ := json.Marshal(a)
	w.Write(jsonResp)
}

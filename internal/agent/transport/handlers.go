package transport

import (
	"AEC/internal/agent/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONdata struct {
	Id   string `json:"id"`
	Task string `json:"task"`
}

func addCal(w http.ResponseWriter, r *http.Request) {
	data := &JSONdata{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		config.Log.Error(err)
		w.WriteHeader(400)
		resp, _ := json.Marshal(`{"err":"Не удалось декодировать JSON"}`)
		w.Write(resp)
		return
	}

	ctx := context.Background()
	config.RedisClientQ.Set(ctx, data.Id, data.Task, 0)

	go func(id string, ww http.ResponseWriter) {
		for {
			answer := config.RedisClientA.Get(ctx, id)
			if answer != nil {
				ww.WriteHeader(200)
				ww.Write([]byte(fmt.Sprintf("%s", answer)))
				return
			}
		}
	}(data.Id, w)
}

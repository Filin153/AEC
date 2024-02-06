package transport

import (
	"AEC/internal/agent/config"
	"context"
	"encoding/json"
	"net/http"
)

type JSONdata struct {
	Id   string `json:"id"`
	Task string `json:"task"`
}

func AddCal(w http.ResponseWriter, r *http.Request) {
	data := &JSONdata{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		config.Log.Error(err)
		w.WriteHeader(400)
		resp, _ := json.Marshal(`{"err":"Не удалось декодировать JSON"}`)
		w.Write(resp)
		return
	}

	config.RedisClientQ.Set(context.Background(), data.Id, data.Task, 0)

	res := make(chan []byte)
	go func(id string) {
		defer close(res)
		for {
			answer := config.RedisClientA.Get(context.Background(), id)
			if answer.Err() == nil && answer.Val() != "" {
				dataA, _ := answer.Bytes()
				res <- dataA
				config.RedisClientA.Del(context.Background(), id)
				return
			}
		}
	}(data.Id)

	w.Write(<-res)
}

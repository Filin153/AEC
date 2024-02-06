package app

import (
	"AEC/internal/agent/config"
	"AEC/internal/agent/services"
	"AEC/internal/agent/transport"
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

func Run() {
	defer config.RedisClientQ.Close()
	defer config.RedisClientA.Close()

	taskChan := make(chan []interface{})

	go services.StartWorkers(config.Conf.Worker, taskChan)

	go func() {
		data := make([]interface{}, 2)
		for {
			keys, err := config.RedisClientQ.Keys(context.Background(), "*").Result()
			if err != nil {
				config.Log.Error(err)
				continue
			}

			for _, key := range keys {
				data[0] = key

				val := config.RedisClientQ.Get(context.Background(), key)

				data[1] = val.Val()

				taskChan <- data

				config.RedisClientQ.Del(context.Background(), key)

			}
		}
	}()

	router := mux.NewRouter()

	router.HandleFunc("/do", transport.AddCal).Methods("POST")

	err := http.ListenAndServe(":"+config.Conf.Port, router)

	if err != nil {
		config.Log.Error(err)
		panic(err)
	}
}

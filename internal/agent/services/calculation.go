package services

import (
	"AEC/internal/agent/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
)

type AnswerData struct {
	Ex     string `json:"ex"`
	Answer string `json:"answer"`
	Err    string `json:"err"`
}

func StartWorkers(max int, task chan []interface{}) {
	for i := 0; i < max; i++ {
		go func() {
			for v := range task {
				calRes, err := calculation(fmt.Sprintf("%s", v[1]))
				if err != nil {
					config.Log.Error(err)
				}

				data, _ := json.Marshal(calRes)
				info := config.RedisClientA.Set(context.Background(), fmt.Sprintf("%s", v[0]), data, 0)
				if info.Err() != nil {
					config.Log.Error(info.Err())
				}
			}
		}()
	}
}

func AddTask(task chan []interface{}) {
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

			task <- data

			config.RedisClientQ.Del(context.Background(), key)

		}
	}
}

func calculation(data string) (AnswerData, error) {
	a := AnswerData{}
	expression, err := govaluate.NewEvaluableExpression(data)
	if err != nil {
		config.Log.WithField("err", "Ошибка при создании выражения").Error(err)
		a.Err = fmt.Sprintf("%v", err)
		return a, err
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		config.Log.WithField("err", "Ошибка при вычислении выражения").Error(err)
		a.Err = fmt.Sprintf("%v", err)
		return a, err
	}

	a.Ex = data
	a.Answer = fmt.Sprintf("%v", result)

	return a, nil
}

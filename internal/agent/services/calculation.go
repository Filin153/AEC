package services

import (
	"AEC/internal/agent/config"
	"AEC/internal/agent/database"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"strings"
	"time"
)

type AnswerData struct {
	Ex     string `json:"ex"`
	Answer string `json:"answer"`
	Err    string `json:"err"`
}

type JSONdata struct {
	Id       string        `json:"id"`
	Task     string        `json:"task"`
	WaitTime time.Duration `json:"wait_time"`
}

func StartWorkers(max int, task chan []byte) {
	for i := 0; i < max; i++ {
		go func() {
			config.Log.Info("Worker start")
			var data = &JSONdata{}
			for v := range task {
				json.Unmarshal(v, data)
				config.Log.Info("Start do - " + data.Task)

				calRes, err := calculation(fmt.Sprintf("%s", data.Task))
				if err != nil {
					config.Log.Error(err)
				}

				time.Sleep(data.WaitTime)
				go database.UpdateCalRes(fmt.Sprintf("%v", data.Id), calRes.Ex, calRes.Answer, calRes.Err)

			}
		}()
	}
}

func GetWaitTime(val string, add, sub, mult, dev int) time.Duration {
	res := 0
	res += strings.Count(val, "+") * add
	res += strings.Count(val, "-") * sub
	res += strings.Count(val, "*") * mult
	res += strings.Count(val, "/") * dev

	return time.Duration(res) * time.Second
}

func AddTask(task chan []byte) {
	for {
		keys, err := config.RedisClientQ.Keys(context.Background(), "*").Result()
		if err != nil {
			config.Log.Error(err)
			continue
		}

		for _, key := range keys {
			val := config.RedisClientQ.Get(context.Background(), key)
			if val.Err() != nil {
				config.Log.Error(val.Err())
				continue
			}
			jsonByte, err1 := val.Bytes()
			if err1 != nil {
				config.Log.Error(err)
				continue
			}

			task <- jsonByte
			err = config.RedisClientQ.Del(context.Background(), key).Err()
			if err != nil {
				config.Log.Error(err)
				continue
			}
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

func CheckNoReadyEx(task chan []byte) {
	var jsonData = &JSONdata{}
	if data, ok := database.GetAllCalRes(); ok {
		for _, v := range data {
			if v.Res == "" && v.Err == "" {
				jsonData.Id = v.RId
				jsonData.Task = v.Expression
				jsonData.WaitTime = time.Second * time.Duration(v.ToDoTime)

				jsonByte, err := json.Marshal(jsonData)
				if err != nil {
					config.Log.Error(err)
					continue
				}

				task <- jsonByte
			}
		}
	}
}

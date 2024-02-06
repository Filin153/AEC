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
					fmt.Println(info.Err())
				}
			}
		}()
	}
}

func calculation(data string) (AnswerData, error) {
	expression, err := govaluate.NewEvaluableExpression(data)
	if err != nil {
		fmt.Println("Ошибка при создании выражения:", err)
		return AnswerData{}, err
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		fmt.Println("Ошибка при вычислении выражения:", err)
		return AnswerData{}, err
	}

	answer := AnswerData{
		Ex:     data,
		Answer: fmt.Sprintf("%v", result),
	}

	return answer, nil
}

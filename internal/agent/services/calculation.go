package services

import (
	"AEC/internal/agent/config"
	"context"
	"fmt"
	"github.com/Knetic/govaluate"
)

func startWorkers(max int, task chan []interface{}) {
	for i := 0; i < max; i++ {
		go func() {
			for v := range task {
				calRes, err := calculation(fmt.Sprintf("%s", v[0]))
				if err != nil {
					config.Log.Error(err)
				}

				ctx := context.Background()
				config.RedisClientA.Set(ctx, fmt.Sprintf("%s", v[1]), calRes, 0)
			}
		}()
	}
}

func calculation(data string) ([]string, error) {
	fmt.Println(data)
	res := make([]string, 2)
	expression, err := govaluate.NewEvaluableExpression(data)
	if err != nil {
		fmt.Println("Ошибка при создании выражения:", err)
		return nil, err
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		fmt.Println("Ошибка при вычислении выражения:", err)
		return nil, err
	}

	res[0] = fmt.Sprintf("%v", result)
	res[1] = data

	return res, nil
}

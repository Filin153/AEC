package services

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"strings"
	"sync"
)

func start(val string) string {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for _, subexpression := range findSubexpressions(val) {
			calRes, err := calculation(subexpression)
			if err != nil {
				fmt.Println(err)
			}
			val = strings.Replace(val, calRes[1], calRes[0], 1)
		}
	}()

	wg.Wait()

	return val
}

func findSubexpressions(expression string) []string {
	var subexpressions []string
	stack := 0
	start := -1

	for i, char := range expression {
		if char == '(' {
			if stack == 0 {
				start = i
			}
			stack++
		} else if char == ')' {
			stack--
			if stack == 0 && start != -1 {
				subexpression := expression[start+1 : i]
				subexpressions = append(subexpressions, subexpression)
				start = -1
			}
		}
	}

	return subexpressions
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

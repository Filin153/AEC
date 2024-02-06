package services

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/database"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type dataForReq struct {
	Id   string `json:"id"`
	Task string `json:"task"`
}

func shuffleSlice(input []string) {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел с текущим временем

	n := len(input)
	for i := n - 1; i > 0; i-- {
		// Генерация случайного индекса от 0 до i (включительно)
		j := rand.Intn(i + 1)

		// Обмен значениями между i-м и случайно выбранным индексом
		input[i], input[j] = input[j], input[i]
	}
}

func containsLetters(input string) bool {
	re := regexp.MustCompile("[a-zA-Z]")
	return re.MatchString(input)
}

func getRandomServer() (Server, error) {
	keys, err := config.RedisClient.Keys(context.Background(), "*").Result()
	if err != nil {
		config.Log.Error(err)
		return Server{}, err
	}

	shuffleSlice(keys)

	for _, key := range keys {
		info := config.RedisClient.Get(context.Background(), key)

		dataByte, _ := info.Bytes()

		var Serv Server
		json.Unmarshal(dataByte, &Serv)

		if Serv.Status == 1 {
			return Serv, nil
		}
	}

	return Server{}, errors.New("Нету доступных серверов")
}

func requestToCalculation(serv Server, subst string, idPrefix int) (map[string]interface{}, error) {
	var data dataForReq

	data.Id = HashSome(subst) + "_" + strconv.Itoa(idPrefix)
	data.Task = subst

	jsonData, err := json.Marshal(data)
	if err != nil {
		config.Log.Error(err)
		return map[string]interface{}{}, err
	}

	req, err := http.NewRequest("POST", serv.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		config.Log.Error(err)
		return map[string]interface{}{}, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		config.Log.Error(err)
		return map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	var answer map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&answer)
	if err != nil {
		config.Log.Error(err)
		return map[string]interface{}{}, err
	}

	return answer, nil
}

func Direct(val, id string) (string, error) {
	if containsLetters(val) {
		config.Log.WithField("ex", val).Error("Выражение содержит буквы")
		go database.UpdateTask(id, "", false, "Выражение содержит буквы")
		return "", errors.New("Выражение содержит буквы")
	}

	trueServ, err := getRandomServer()
	if err != nil {
		config.Log.Error(err)
		return "", err
	}

	for i, subexpression := range findSubexpressions(val) {
		calRes, err1 := requestToCalculation(trueServ, subexpression, i)
		if err1 != nil {
			config.Log.Error(err1)
			go database.UpdateTask(id, "", false, fmt.Sprintf("%v", err))
			return "", err1
		} else if calRes["err"] != "" {
			config.Log.Error(calRes["err"])
			go database.UpdateTask(id, "", false, fmt.Sprintf("%v", calRes["err"]))
			return "", err
		}
		val = strings.Replace(val, fmt.Sprintf("%v", calRes["ex"]), fmt.Sprintf("%v", calRes["answer"]), 1)
	}

	calRes, err := requestToCalculation(trueServ, val, 0)
	if err != nil {
		config.Log.Error(err)
		go database.UpdateTask(id, "", false, fmt.Sprintf("%v", err))
		return "", err
	} else if calRes["err"] != "" {
		config.Log.Error(calRes["err"])
		go database.UpdateTask(id, "", false, fmt.Sprintf("%v", calRes["err"]))
		return "", err
	}
	val = strings.Replace(val, fmt.Sprintf("%v", calRes["ex"]), fmt.Sprintf("%v", calRes["answer"]), 1)

	go database.UpdateTask(id, val, true, "")

	return val, nil
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

func AddCompToBD() {
	for {

	}
}

package services

import (
	"AEC/internal/orchestrator/config"
	"AEC/internal/orchestrator/database"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type dataForReq struct {
	Id       string        `json:"id"`
	Task     string        `json:"task"`
	WaitTime time.Duration `json:"wait_time"`
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

func removeFromSlice(sl []string, id int) []string {
	if id >= 0 && id < len(sl) {
		sl = append(sl[:id], sl[id+1:]...)
	} else {
		fmt.Println("Index out of bounds")
	}
	return sl
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

func requestToCalculation(serv Server, subst string, add, sub, mult, div int) (map[string]interface{}, error) {
	var data dataForReq

	data.Id = HashSome(subst)
	data.Task = subst
	data.WaitTime = GetWaitTime(subst, add, sub, mult, div)

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

func takeCalRes(ids []string, val string) ([]string, string, int, string) {
	needTime := 0
	for {
		if len(ids) == 0 {
			break
		}
		for i, vid := range ids {
			if calRes, ok := database.GetCalRes(vid); ok && calRes.Res != "" {
				if calRes.Err != "" {
					return []string{}, "", 0, calRes.Err
				}
				val = strings.Replace(val, calRes.Expression, calRes.Res, 1)
				ids = removeFromSlice(ids, i)
				needTime += calRes.ToDoTime
				config.Log.WithFields(logrus.Fields{
					"val":        val,
					"Expression": calRes.Expression,
					"Res":        calRes.Res,
					"AllId":      ids,
				}).Info("OK")
			}
		}
		time.Sleep(time.Second)
	}

	return ids, val, needTime, ""
}

func Direct(val, id string, add, sub, mult, div int) (string, error) {
	tempVal := val
	resTime := 0

	if containsLetters(val) {
		config.Log.WithField("ex", val).Error("Выражение содержит буквы")
		go database.UpdateTask(id, "", "", false, "Выражение содержит буквы", 0)
		return "", errors.New("Выражение содержит буквы")
	}

	trueServ, err := getRandomServer()
	if err != nil {
		config.Log.Error(err)
		return "", err
	}

	allId := []string{}
	for _, subexpression := range findSubexpressions(val) {
		go requestToCalculation(trueServ, subexpression, add, sub, mult, div)
		allId = append(allId, HashSome(subexpression))
	}

	allId, val, resTime, err1 := takeCalRes(allId, val)
	if err1 != "" {
		go database.UpdateTask(id, "", "", false, err1, 0)
		config.Log.Error(err1)
		return "", errors.New(err1)
	}
	go requestToCalculation(trueServ, val, add, sub, mult, div)
	allId = append(allId, HashSome(val))
	_, val, tempTime, err1 := takeCalRes(allId, val)
	if err1 != "" {
		go database.UpdateTask(id, "", "", false, err1, 0)
		config.Log.Error(err1)
		return "", errors.New(err1)
	}

	resTime += tempTime
	go database.UpdateTask(id, "", val, true, "", resTime)

	config.Log.Info("Готово - ", tempVal)
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

func GetWaitTime(val string, add, sub, mult, dev int) time.Duration {
	res := 0
	res += strings.Count(val, "+") * add
	res += strings.Count(val, "-") * sub
	res += strings.Count(val, "*") * mult
	res += strings.Count(val, "/") * dev

	return time.Duration(res) * time.Second
}

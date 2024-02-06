package services

import (
	"AEC/internal/orchestrator/config"
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	Id       string    `json:"id"`
	URL      string    `json:"url"`
	Status   int       `json:"status"`
	LastPing time.Time `json:"last_ping"`
}

func serverToRedis(s Server) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	info1 := config.RedisClient.Set(context.Background(), s.Id, data, 0)
	if info1.Err() != nil {
		return info1.Err()
	}

	return nil
}

func RemoveServerFromRedis(id string) error {
	info1 := config.RedisClient.Del(context.Background(), id)
	if info1.Err() != nil {
		config.Log.Error(info1.Err())
		return info1.Err()
	}

	return nil
}

func HashSome(val string) string {
	hasher := sha256.New()
	hasher.Write([]byte(val))
	hashInBytes := hasher.Sum(nil)
	return string(hashInBytes)
}

func GetClientIP(r *http.Request) string {
	// В реальном приложении, возможно, вам захочется учесть возможность, что X-Forwarded-For может содержать несколько IP-адресов.
	// Также следует учесть, что значение X-Forwarded-For может быть легко поддельно.
	// Ваша реализация может варьироваться в зависимости от конкретных требований вашего приложения.

	// Попытка получить IP-адрес из заголовка X-Forwarded-For
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// Если X-Forwarded-For отсутствует, используем RemoteAddr
	return r.RemoteAddr
}

func AddServer(id, URL string) error {

	info := config.RedisClient.Get(context.Background(), URL)

	dataByte, _ := info.Bytes()
	if len(dataByte) > 0 {
		var oldServ Server
		err := json.Unmarshal(dataByte, &oldServ)
		if err != nil {
			return err
		}

		oldServ.LastPing = time.Now()
		oldServ.Status = 1

		err = serverToRedis(oldServ)
		if err != nil {
			return err
		}

		return nil
	}

	server := Server{
		Id:       id,
		URL:      URL,
		Status:   1,
		LastPing: time.Now(),
	}

	err := serverToRedis(server)
	if err != nil {
		return err
	}

	return nil
}

func AllServer() []Server {
	keys, err := config.RedisClient.Keys(context.Background(), "*").Result()
	if err != nil {
		config.Log.Error(err)
		return nil
	}

	allS := []Server{}
	for _, key := range keys {
		info := config.RedisClient.Get(context.Background(), key)

		dataByte, _ := info.Bytes()

		var Serv Server
		json.Unmarshal(dataByte, &Serv)
		allS = append(allS, Serv)
	}

	return allS
}

func CheckServer() {
	for {
		keys, err := config.RedisClient.Keys(context.Background(), "*").Result()
		if err != nil {
			config.Log.Error(err)
			continue
		}

		for _, key := range keys {
			info := config.RedisClient.Get(context.Background(), key)

			dataByte, _ := info.Bytes()

			var Serv Server
			json.Unmarshal(dataByte, &Serv)

			if time.Now().Sub(Serv.LastPing) >= (time.Minute * 5) {
				Serv.Status = 2
				serverToRedis(Serv)
			}
		}
	}
}

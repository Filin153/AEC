package models

import (
	"gorm.io/gorm"
	"strings"
)

type Task struct {
	gorm.Model
	Expression string `gorm:"type:varchar(500)"`
	Req_id     string `gorm:"type:varchar(65);unique"` // Хэш выражения
	User_id    string `gorm:"type:string"`             // Хэш времени, хранит все ID через , если данную задачу запросило несколько пользователей
	Status     bool   `gorm:"default:false"`
	ToDoTime   int    `gorm:"type:integer"`
	Res        string `gorm:"type:string"`
	Err        string `gorm:"type:string"`
}

// Переводит строку в слайс
func (t *Task) GetUserIDs() []string {
	return strings.Split(t.User_id, ",")
}

// Переводит слайс в строку
func (t *Task) SetUserIDs(userIDs []string) {
	t.User_id = strings.Join(userIDs, ",")
}

type CalRes struct {
	gorm.Model
	RId        string `gorm:"type:varchar(65);unique"` // Хэш выражения
	Expression string `gorm:"type:varchar(500)"`
	Res        string `gorm:"type:string"`
	Err        string `gorm:"type:string"`
	ToDoTime   int    `gorm:"type:integer"`
}

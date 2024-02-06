package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Expression string `gorm:"type:varchar(500)"`
	Req_id     string `gorm:"type:varchar(65);unique"` // Хэш выражения
	User_id    string `gorm:"type:varchar(65)"`        // Хэш времени
	Status     bool   `gorm:"default:false"`
	Res        string `gorm:"type:string"`
	Err        string `gorm:"type:string"`
}

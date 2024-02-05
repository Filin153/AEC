package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Expression string `gorm:"type:varchar(500)"`
	Req_id     string `gorm:"type:varchar(65)"`
	Status     int    `gorm:"type:integer"`
	Res        string `gorm:"type:string"`
}

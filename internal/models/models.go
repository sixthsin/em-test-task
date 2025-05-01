package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null"`
	Surname     string `gorm:"type:varchar(100);not null"`
	Patronymic  string `gorm:"type:varchar(100)"`
	Age         uint   `gorm:"type:uint"`
	Gender      string `gorm:"type:string"`
	Nationality string `gorm:"type:string"`
}

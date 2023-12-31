package model

import "gorm.io/gorm"

type Car struct {
	gorm.Model
	CarModel     string  `json:"model"`
	Registration string  `json:"registration" gorm:"unique;not null"`
	Mileage      float64 `json:"mileage"`
	Available    bool    `json:"available"`
}

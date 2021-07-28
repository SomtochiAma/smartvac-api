package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	MinUnit int `json:"min_unit"`
	Address string `json:"address"`
	CurrentReading []CurrentReading
}

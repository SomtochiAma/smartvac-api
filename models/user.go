package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	MinUnit        int    `json:"min_unit"`
	UsedUnit   int    `json:"used_unit"`
	TotalUnit int    `json:"total_unit"`
	CurrentReading []Reading
	Payment        []Payment
}

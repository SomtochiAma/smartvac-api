package models

import "time"

type Reading struct {
	UserID     int       `json:"user_id" binding:"required"`
	Current    int       `json:"current" binding:"required"`
	Power      int       `json:"power" binding:"required"`
	TotalPower int       `json:"total_power" binding:"required"`
	Time       time.Time `json:"time" binding:"required"`
}

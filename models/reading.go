package models

import "time"

type Reading struct {
	UserID     int       `json:"user_id"`
	Current    float32   `json:"current"`
	Power      float32   `json:"power"`
	TotalPower float32   `json:"total_power" binding:"required"`
	Time       time.Time `json:"time"`
}

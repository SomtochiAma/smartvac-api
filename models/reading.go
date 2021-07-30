package models

import "time"

type Reading struct {
	UserID     int       `json:"user_id"`
	Current    int       `json:"current"`
	Power      int       `json:"power"`
	TotalPower int       `json:"total_power" binding:"required"`
	Time       time.Time `json:"time" binding:"required"`
}

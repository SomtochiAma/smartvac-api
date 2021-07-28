package models

import "time"

type CurrentReading struct {
	UserID uint    `json:"user_id" binding:"required"`
	Value  uint    `json:"value" binding:"required"`
	Time   time.Time `json:"time" binding:"required"`
}

package models

import (
	"time"
)

type Payment struct {
	UserID        int       `json:"user_id"`
	Amount        int       `json:"amount" binding:"required"`
	Units         int       `json:"units" binding:"required"`
	Day           time.Time `json:"day" binding:"required"`
	AvailableUnit int       `json:"available_unit"`
}

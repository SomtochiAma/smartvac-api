package models

import (
	"time"
)

type Payment struct {
	UserID int       `json:"user_id"`
	Amount int       `json:"amount"`
	Units  int       `json:"units"`
	Day    time.Time `json:"day"`
}

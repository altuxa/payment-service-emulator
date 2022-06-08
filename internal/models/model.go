package models

import "time"

const (
	StatusNew     = "NEW"
	StatusSuccess = "SUCCESS"
	StatusFail    = "FAIL"
	StatusError   = "ERROR"
)

type Transaction struct {
	ID           int
	UserID       int    `json:"userID"`
	UserEmail    string `json:"email"`
	Sum          int    `json:"sum"`
	Valute       string `json:"valuta"`
	CreationDate time.Time
	ChangeDate   time.Time
	Status       string
}

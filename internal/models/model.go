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
	UserID       int
	UserEmail    string
	Sum          int
	Valute       string
	CreationDate time.Time
	ChangeDate   time.Time
	Status       string
}

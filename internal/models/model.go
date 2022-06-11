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
	UserID       int    `json:"UserID"`
	UserEmail    string `json:"Email"`
	Sum          int    `json:"Sum"`
	Currency     string `json:"Currency"`
	CreationDate time.Time
	ChangeDate   time.Time
	Status       string
}

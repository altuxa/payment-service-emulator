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
	UserID       int     `json:"UserID"`
	UserEmail    string  `json:"Email"`
	Sum          float64 `json:"Sum"`
	Currency     string  `json:"Currency"`
	CreationDate time.Time
	ChangeDate   time.Time
	Status       string
}

type PaymentProcessingInput struct {
	Email string `json:"Email"`
}

type InputByUserEmail struct {
	Email string `json:"email"`
}

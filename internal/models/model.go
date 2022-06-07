package models

import "time"

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

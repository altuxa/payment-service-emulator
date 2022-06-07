package repository

import "database/sql"

type User interface{}

// type Transaction interface{}

type Payment interface {
	NewPayment()
}

type Repository struct {
	User
	Payment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:    NewUserRepo(db),
		Payment: NewPaymentRepo(db),
	}
}

package repository

import "database/sql"

type User interface{}

// type Transaction interface{}

type Payment interface {
	NewPayment(id int, email string, sum int, val string)
	PaymentStatus(paymentId int)
	GetAllPaymentsByID(userId int)
	GetAllPaymentsByEmail(email string)
	CancelPayment(paymentId int)
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

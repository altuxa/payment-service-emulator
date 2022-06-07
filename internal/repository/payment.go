package repository

import "database/sql"

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{
		db: db,
	}
}

func (p *PaymentRepo) NewPayment(id int, email string, sum int, val string) {
}

func (p *PaymentRepo) PaymentStatus(paymentId int) {
}

func (p *PaymentRepo) GetAllPaymentsByID(userId int) {
}

func (p *PaymentRepo) GetAllPaymentsByEmail(email string) {
}

func (p *PaymentRepo) CancelPayment(paymentId int) {
}

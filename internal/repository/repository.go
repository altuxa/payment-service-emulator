package repository

import (
	"database/sql"

	"github.com/altuxa/payment-service-emulator/internal/models"
)

type User interface {
	NewUser()
}

// type Transaction interface{}

type Payment interface {
	NewPayment(id int, email string, sum int, val string) error
	PaymentStatus(paymentId int) (string, error)
	GetAllPaymentsByUserID(userId int) ([]models.Transaction, error)
	GetAllPaymentsByEmail(email string) ([]models.Transaction, error)
	// CancelPayment(paymentId int)
	DeletePayment(paymentId int) error
}

type Repositories struct {
	User
	Payment
}

func NewRepository(db *sql.DB) *Repositories {
	return &Repositories{
		User:    NewUserRepo(db),
		Payment: NewPaymentRepo(db),
	}
}

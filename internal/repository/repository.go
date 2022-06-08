package repository

import (
	"database/sql"

	"github.com/altuxa/payment-service-emulator/internal/models"
)

type User interface {
	NewUser()
}

type Payment interface {
	NewPayment(id int, email string, sum int, val string, status string) (int, error)
	PaymentStatus(paymentId int) (string, error)
	GetAllPaymentsByUserID(userId int) ([]models.Transaction, error)
	GetAllPaymentsByEmail(email string) ([]models.Transaction, error)
	DeletePayment(paymentId int) error
	SetStatusSuccess(paymentId int) error
	SetStatusFail(paymentId int) error
	// SetStatusError(paymentId int) error
	PaymentErrorImitation() bool
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

package service

import (
	"github.com/altuxa/payment-service-emulator/internal/models"
	"github.com/altuxa/payment-service-emulator/internal/repository"
)

type User interface{}

type Payment interface {
	CancelPayment(paymentId int) error
	CreatePayment(id int, email string, sum int, val string) (int, string, error)
	PaymentProcessing(id int) error
	PaymentStatus(paymentId int) (string, error)
	ByUserID(userID int) ([]models.Transaction, error)
	ByUserEmail(email string) ([]models.Transaction, error)
}

type Services struct {
	User
	Payment
}

type ServiceDeps struct {
	Repos *repository.Repositories
}

func NewService(repo *repository.Repositories) *Services {
	return &Services{
		User:    NewUserService(repo.User),
		Payment: NewPaymentService(repo.Payment),
	}
}

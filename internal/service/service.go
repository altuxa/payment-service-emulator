package service

import "github.com/altuxa/payment-service-emulator/internal/repository"

type User interface{}

type Payment interface {
	CancelPayment(paymentId int) error
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

package service

import "github.com/altuxa/payment-service-emulator/internal/repository"

type User interface{}

type Payment interface{}

type Services struct {
	User
	Payment
}

type ServiceDeps struct {
	Repos *repository.Repositories
}

func NewService(repo *repository.Repositories) *Services {
	return &Services{
		Payment: NewPaymentService(repo.Payment),
	}
}

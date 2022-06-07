package service

import "github.com/altuxa/payment-service-emulator/internal/repository"

type PaymentService struct {
	repo repository.Payment
}

func NewPaymentService(repo repository.Payment) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

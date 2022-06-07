package service

import (
	"fmt"

	"github.com/altuxa/payment-service-emulator/internal/repository"
)

type PaymentService struct {
	repo repository.Payment
}

func NewPaymentService(repo repository.Payment) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

func (p *PaymentService) CancelPayment(paymentId int) error {
	status, err := p.repo.PaymentStatus(paymentId)
	if err != nil {
		return err
	}
	if status != "NEW" {
		return fmt.Errorf("error payment status = %v", status)
	}
	err = p.repo.DeletePayment(paymentId)
	if err != nil {
		return err
	}
	return nil
}

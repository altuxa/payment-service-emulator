package service

import (
	"fmt"

	"github.com/altuxa/payment-service-emulator/internal/models"
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

func (p *PaymentService) CreatePayment(id int, email string, sum int, val string) (int, error) {
	status := models.StatusNew
	random := p.repo.PaymentErrorImitation()
	if !random {
		status = models.StatusError
	}
	id, err := p.repo.NewPayment(id, email, sum, val, status)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PaymentService) PaymentProcessing(id int) error {
	err := p.repo.SetStatusSuccess(id)
	if err != nil {
		return err
	}
	return nil
}

// func (p *PaymentService) PaymentErrorImitation(id int) bool {
// 	rand.Seed(time.Now().UnixNano())
// 	a := rand.Intn(60)
// 	b := rand.Intn(45)
// 	if a > b {
// 		return true
// 	}
// 	err := p.repo.SetStatusError(id)
// 	return err != nil
// }

package service

import (
	"fmt"
	"time"

	"github.com/altuxa/payment-service-emulator/internal/helpers"
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
	if status != models.StatusNew {
		return fmt.Errorf("error payment status = %v", status)
	}
	err = p.repo.DeletePayment(paymentId)
	if err != nil {
		return err
	}
	return nil
}

func (p *PaymentService) CreatePayment(id int, email string, sum int, val string) (int, string, error) {
	status := models.StatusNew
	random := helpers.PaymentErrorImitation()
	if !random {
		status = models.StatusError
	}
	id, err := p.repo.NewPayment(id, email, sum, val, status)
	if err != nil {
		return 0, status, err
	}
	return id, status, nil
}

func (p *PaymentService) PaymentProcessing(id int) error {
	time.Sleep(2 * time.Second)
	status, err := p.repo.PaymentStatus(id)
	if err != nil {
		return fmt.Errorf("payment not found %w", err)
	}
	if status != models.StatusNew {
		return fmt.Errorf("incorrect payment status %s", status)
	}
	succes := helpers.FailStatusImitation()
	if succes {
		err = p.repo.SetStatusSuccess(id)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	} else if !succes {
		err = p.repo.SetStatusFail(id)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

func (p *PaymentService) PaymentStatus(paymentId int) (string, error) {
	status, err := p.repo.PaymentStatus(paymentId)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (p *PaymentService) ByUserID(userID int) ([]models.Transaction, error) {
	transactions, err := p.repo.GetAllPaymentsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return transactions, nil
}

func (p *PaymentService) ByUserEmail(email string) ([]models.Transaction, error) {
	transactions, err := p.repo.GetAllPaymentsByEmail(email)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

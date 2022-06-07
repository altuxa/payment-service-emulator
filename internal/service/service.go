package service

import "github.com/altuxa/payment-service-emulator/internal/repository"

type Service struct{}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}

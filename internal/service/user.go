package service

import "github.com/altuxa/payment-service-emulator/internal/repository"

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

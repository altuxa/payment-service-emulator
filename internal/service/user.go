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

func (u *UserService) Verification(payId int, email string) (bool, error) {
	checkEmail, err := u.repo.UserVerification(payId, email)
	if err != nil {
		return false, err
	}
	if checkEmail != email {
		return false, nil
	}
	return true, nil
}

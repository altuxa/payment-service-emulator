package handlers

import "github.com/altuxa/payment-service-emulator/internal/service"

type Handler struct {
	userService    service.User
	paymentService service.Payment
}

func NewHandler(service *service.Services) *Handler {
	return &Handler{
		userService:    service.User,
		paymentService: service.Payment,
	}
}

func (h *Handler) Server() {
}

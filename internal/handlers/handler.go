package handlers

import "github.com/altuxa/payment-service-emulator/internal/service"

type Handler struct{}

func NewHandler(service *service.Services) *Handler {
	return &Handler{}
}

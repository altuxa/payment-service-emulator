package handlers

import "github.com/altuxa/payment-service-emulator/internal/service"

type Handler struct{}

func NewHandler(service *service.Service) *Handler {
	return &Handler{}
}
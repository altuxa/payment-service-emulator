package handlers

import (
	"log"
	"net/http"

	"github.com/altuxa/payment-service-emulator/internal/service"
)

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
	mux := http.NewServeMux()
	mux.HandleFunc("/new", h.NewTransaction)
	mux.HandleFunc("/status/", h.StatusByID)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln(err)
	}
}

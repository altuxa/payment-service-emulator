package handlers

import (
	"log"
	"net/http"

	"github.com/altuxa/payment-service-emulator/internal/service"
)

const addr = ":8080"

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
	mux.HandleFunc("/processing/", h.PaymentStatusChange)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalln(err)
	}
}

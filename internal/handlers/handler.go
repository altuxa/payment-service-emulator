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
	mux.HandleFunc("/payments/new", h.NewTransaction)
	mux.HandleFunc("/payments/status/", h.StatusByID)
	mux.HandleFunc("/payments/processing/", h.PaymentStatusChange)
	mux.HandleFunc("/payments/byid/", h.ByUserID)
	mux.HandleFunc("/payments/byemail", h.ByUserEmail)
	mux.HandleFunc("/payments/cancel/", h.CancelPayment)
	log.Println("Server started localhost:8080")
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalln(err)
	}
}

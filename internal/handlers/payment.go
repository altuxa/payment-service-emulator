package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/altuxa/payment-service-emulator/internal/models"
)

func (h *Handler) NewTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	newPayment := models.Transaction{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &newPayment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, status, err := h.paymentService.CreatePayment(newPayment.UserID, newPayment.UserEmail, newPayment.Sum, newPayment.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	paymentID := strconv.Itoa(id)
	jsonData, _ := json.Marshal("paymentID: " + strconv.Itoa(id) + " status: " + status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
	go http.Post("http://localhost:8080/payments/processing/"+paymentID, "application/json", nil)
}

func (h *Handler) StatusByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	strId := strings.TrimPrefix(r.URL.Path, "/payments/status/")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status, err := h.paymentService.PaymentStatus(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, _ := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) PaymentStatusChange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	strId := strings.TrimPrefix(r.URL.Path, "/payments/processing/")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.paymentService.PaymentProcessing(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	strId := strings.TrimPrefix(r.URL.Path, "/payments/byid/")
	userID, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactions, err := h.paymentService.ByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type input struct {
	Email string `json:"email"`
}

func (h *Handler) ByUserEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	in := input{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactions, err := h.paymentService.ByUserEmail(in.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) CancelPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	strID := strings.TrimPrefix(r.URL.Path, "/payments/cancel/")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.paymentService.CancelPayment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

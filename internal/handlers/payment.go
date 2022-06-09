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
	id, status, err := h.paymentService.CreatePayment(newPayment.UserID, newPayment.UserEmail, newPayment.Sum, newPayment.Valute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	a := strconv.Itoa(id)
	jsonData, _ := json.Marshal("paymentID: " + strconv.Itoa(id) + " status: " + status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	http.Post("http://localhost:8080/processing/"+a, "", nil)
	// log.Println(id)
	// g := strings.NewReader(a)
	// buf := make([]byte, 1)
	// http.NewRequest(http.MethodPost, "http://localhost:8080/processing/"+strconv.Itoa(id), nil)
	w.Write(jsonData)
}

func (h *Handler) StatusByID(w http.ResponseWriter, r *http.Request) {
	strId := strings.TrimPrefix(r.URL.Path, "/status/")
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
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) PaymentStatusChange(w http.ResponseWriter, r *http.Request) {
	strId := strings.TrimPrefix(r.URL.Path, "/processing/")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// var id string
	// reqBody, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// err = json.Unmarshal(reqBody, &id)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	err = h.paymentService.PaymentProcessing(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

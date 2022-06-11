package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/altuxa/payment-service-emulator/internal/models"
	"github.com/altuxa/payment-service-emulator/internal/service"
	mock_service "github.com/altuxa/payment-service-emulator/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	type Mock func(s *mock_service.MockPayment, tr models.Transaction)
	tData := map[string]struct {
		Input               models.Transaction
		InputBody           string
		Method              string
		mock                Mock
		ExpectedRequestBody string
		ExpectedStatusCode  int
	}{
		"OK": {
			Input: models.Transaction{
				UserID:    1,
				UserEmail: "ann@mail.ru",
				Sum:       502.3,
				Currency:  "USD",
			},
			InputBody: `{"UserID":1,"Email":"ann@mail.ru","Sum":502.3,"Currency":"USD"}`,
			Method:    "POST",
			mock: func(s *mock_service.MockPayment, tr models.Transaction) {
				s.EXPECT().CreatePayment(tr.UserID, tr.UserEmail, tr.Sum, tr.Currency).Return(1, models.StatusNew, nil)
			},
			ExpectedRequestBody: "\"paymentID: 1 status: NEW\"",
			ExpectedStatusCode:  201,
		},
	}
	for tName, tCase := range tData {
		v := tCase
		t.Run(tName, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()
			pay := mock_service.NewMockPayment(c)
			v.mock(pay, v.Input)
			services := &service.Services{
				Payment: pay,
			}
			handler := NewHandler(services)
			// test server
			r := http.HandlerFunc(handler.NewTransaction)
			// test req
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/payments/new", bytes.NewBufferString(v.InputBody))

			// perform request
			r.ServeHTTP(w, req)
			assert.Equal(t, v.ExpectedRequestBody, w.Body.String())
			assert.Equal(t, v.ExpectedStatusCode, w.Code)
			assert.Equal(t, v.Method, req.Method)
		})
	}
}

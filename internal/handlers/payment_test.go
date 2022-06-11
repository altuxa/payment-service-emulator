package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
			ExpectedStatusCode:  200,
		},
		"bad req": {
			Input: models.Transaction{
				UserID:    1,
				UserEmail: "ann@mail.ru",
				Sum:       502.3,
				Currency:  "USD",
			},
			InputBody: `{"UserID":1,"Email":"ann@mail.ru","Sum":502.3,"Currency":"USD"}`,
			Method:    "POST",
			mock: func(s *mock_service.MockPayment, tr models.Transaction) {
				s.EXPECT().CreatePayment(tr.UserID, tr.UserEmail, tr.Sum, tr.Currency).Return(0, "", errors.New("bad req"))
			},
			ExpectedRequestBody: "bad req\n",
			ExpectedStatusCode:  400,
		},
		"method not allowed": {
			Method:              "GET",
			mock:                func(s *mock_service.MockPayment, tr models.Transaction) {},
			ExpectedRequestBody: "method not allowed\n",
			ExpectedStatusCode:  405,
		},
		"unmarshal error": {
			Input: models.Transaction{
				UserID:    1,
				UserEmail: "ann@mail.ru",
				Sum:       502.3,
				Currency:  "USD",
			},
			Method:              "POST",
			mock:                func(s *mock_service.MockPayment, tr models.Transaction) {},
			ExpectedRequestBody: "unexpected end of JSON input\n",
			ExpectedStatusCode:  400,
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
			req := httptest.NewRequest(v.Method, "/payments/new", bytes.NewBufferString(v.InputBody))

			// perform request
			r.ServeHTTP(w, req)
			assert.Equal(t, v.ExpectedRequestBody, w.Body.String())
			assert.Equal(t, v.ExpectedStatusCode, w.Code)
		})
	}
}

func TestStatusByID(t *testing.T) {
	type mock func(s *mock_service.MockPayment, id int)
	tData := map[string]struct {
		URL                 string
		Input               int
		Method              string
		ExpectedStatusCode  int
		ExpectedRequestBody string
		Mock                mock
	}{
		"success": {
			URL:                 "/payments/status/1",
			Input:               1,
			Method:              "GET",
			ExpectedStatusCode:  200,
			ExpectedRequestBody: `"SUCCESS"`,
			Mock: func(s *mock_service.MockPayment, id int) {
				s.EXPECT().PaymentStatus(id).Return(models.StatusSuccess, nil)
			},
		},
		"invalid input": {
			URL:                 "/payments/status/",
			Input:               1,
			Method:              "GET",
			ExpectedStatusCode:  400,
			ExpectedRequestBody: "invalid input\n",
			Mock:                func(s *mock_service.MockPayment, id int) {},
		},
		"payment not found": {
			URL:                 "/payments/status/999",
			Input:               999,
			Method:              "GET",
			ExpectedStatusCode:  400,
			ExpectedRequestBody: "payment not found\n",
			Mock: func(s *mock_service.MockPayment, id int) {
				s.EXPECT().PaymentStatus(id).Return("", errors.New("payment not found"))
			},
		},
		"invalid method": {
			URL:                 "/payments/status/1",
			Input:               1,
			Method:              "POST",
			ExpectedStatusCode:  405,
			ExpectedRequestBody: "method not allowed\n",
			Mock: func(s *mock_service.MockPayment, id int) {
			},
		},
	}
	for tName, tCase := range tData {
		v := tCase
		t.Run(tName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			pay := mock_service.NewMockPayment(c)
			v.Mock(pay, v.Input)
			services := service.Services{
				Payment: pay,
			}
			handler := NewHandler(&services)
			r := http.HandlerFunc(handler.StatusByID)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(v.Method, v.URL, nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, v.ExpectedRequestBody, w.Body.String())
			assert.Equal(t, v.ExpectedStatusCode, w.Code)
		})
	}
}

func TestByUserID(t *testing.T) {
	type mock func(s *mock_service.MockPayment, id int)
	tData := map[string]struct {
		URL                 string
		ID                  int
		Method              string
		ExpectedStatusCode  int
		ExpectedRequestBody string
		Mock                mock
	}{
		"succes": {
			URL:                 "/payments/byid/1",
			ID:                  1,
			Method:              "GET",
			ExpectedStatusCode:  200,
			ExpectedRequestBody: `[{"ID":114,"UserID":1,"Email":"ann@mail.ru","Sum":1000,"Currency":"KZ","CreationDate":"2022-06-11T18:45:47.72474801+06:00","ChangeDate":"2022-06-11T18:47:22.683292944+06:00","Status":"SUCCESS"}]`,
			Mock: func(s *mock_service.MockPayment, id int) {
				s.EXPECT().ByUserID(id).Return([]models.Transaction{
					models.Transaction{
						ID:           114,
						UserID:       1,
						UserEmail:    "ann@mail.ru",
						Sum:          1000,
						CreationDate: time.Date(2022, 06, 11, 18, 45, 47, 724748010, time.Local),
						ChangeDate:   time.Date(2022, 06, 11, 18, 47, 22, 683292944, time.Local),
						Currency:     "KZ",
						Status:       "SUCCESS",
					},
				}, nil)
			},
		},
		"invalid method": {
			URL:                 "/payments/byid/1",
			ID:                  1,
			Method:              "POST",
			ExpectedStatusCode:  405,
			ExpectedRequestBody: "method not allowed\n",
			Mock:                func(s *mock_service.MockPayment, id int) {},
		},
		// "bad req": {
		// 	URL:                 "/payments/byid/1",
		// 	ID:                  1,
		// 	Method:              "POST",
		// 	ExpectedStatusCode:  405,
		// 	ExpectedRequestBody: "method not allowed\n",
		// 	Mock:                func(s *mock_service.MockPayment, id int) {},
		// },
	}
	for tName, tCase := range tData {
		v := tCase
		t.Run(tName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			pay := mock_service.NewMockPayment(c)
			v.Mock(pay, v.ID)
			services := service.Services{
				Payment: pay,
			}
			handler := NewHandler(&services)
			r := http.HandlerFunc(handler.ByUserID)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(v.Method, v.URL, nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, v.ExpectedRequestBody, w.Body.String())
			assert.Equal(t, v.ExpectedStatusCode, w.Code)
		})
	}
}

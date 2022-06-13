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
		"success": {
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

func TestPaymentProcessing(t *testing.T) {
	type mockPay func(s *mock_service.MockPayment, payId int)
	type mockUser func(s *mock_service.MockUser, payId int, in models.PaymentProcessingInput)
	tData := map[string]struct {
		URL                 string
		InputBody           string
		Input               models.PaymentProcessingInput
		Method              string
		ExpectedRequestBody string
		ExpectedStatusCode  int
		MockUser            mockUser
		MockPay             mockPay
	}{
		"Success": {
			URL:       "/payments/processing/1",
			InputBody: `{"Email":"ann@mail.ru"}`,
			Input: models.PaymentProcessingInput{
				Email: "ann@mail.ru",
			},
			Method:              "POST",
			ExpectedRequestBody: "\"SUCCESS\"",
			ExpectedStatusCode:  200,
			MockUser: func(s *mock_service.MockUser, payId int, in models.PaymentProcessingInput) {
				s.EXPECT().Verification(payId, in.Email).Return(true, nil)
			},
			MockPay: func(s *mock_service.MockPayment, payId int) {
				s.EXPECT().PaymentProcessing(payId).Return(models.StatusSuccess, nil)
			},
		},
		"Invalid method": {
			URL:                 "/payments/processing/1",
			Method:              "GET",
			ExpectedRequestBody: "method not allowed\n",
			ExpectedStatusCode:  405,
			MockUser:            func(s *mock_service.MockUser, payId int, in models.PaymentProcessingInput) {},
			MockPay:             func(s *mock_service.MockPayment, payId int) {},
		},
		"Invalid payID input": {
			URL:                 "/payments/processing/a",
			Method:              "POST",
			ExpectedRequestBody: "Invalid input\n",
			ExpectedStatusCode:  400,
			MockUser:            func(s *mock_service.MockUser, payId int, in models.PaymentProcessingInput) {},
			MockPay:             func(s *mock_service.MockPayment, payId int) {},
		},
		"User not found": {
			URL:       "/payments/processing/1",
			InputBody: `{"Email":"aboba@mail.ru"}`,
			Input: models.PaymentProcessingInput{
				Email: "aboba@mail.ru",
			},
			Method:              "POST",
			ExpectedRequestBody: "not enough rights not found\n",
			ExpectedStatusCode:  400,
			MockUser: func(s *mock_service.MockUser, payId int, in models.PaymentProcessingInput) {
				s.EXPECT().Verification(payId, in.Email).Return(false, errors.New("not found"))
			},
			MockPay: func(s *mock_service.MockPayment, payId int) {},
		},
		"Unauthorized": {
			URL:       "/payments/processing/1",
			InputBody: `{"Email":"alex@mail.ru"}`,
			Input: models.PaymentProcessingInput{
				Email: "alex@mail.ru",
			},
			Method:              "POST",
			ExpectedRequestBody: "not enough rights\n",
			ExpectedStatusCode:  401,
			MockUser: func(s *mock_service.MockUser, payId int, in models.PaymentProcessingInput) {
				s.EXPECT().Verification(payId, in.Email).Return(false, nil)
			},
			MockPay: func(s *mock_service.MockPayment, payId int) {},
		},
		"Broken json": {
			URL:       "/payments/processing/1",
			InputBody: `{"Email":"alex@mail.ru}`,
			Input: models.PaymentProcessingInput{
				Email: "alex@mail.ru",
			},
			Method:              "POST",
			ExpectedRequestBody: "unexpected end of JSON input\n",
			ExpectedStatusCode:  500,
			MockUser:            func(s *mock_service.MockUser, payId int, in models.PaymentProcessingInput) {},
			MockPay:             func(s *mock_service.MockPayment, payId int) {},
		},
		"Invalid Payments status": {
			URL:       "/payments/processing/1",
			InputBody: `{"Email":"alex@mail.ru"}`,
			Input: models.PaymentProcessingInput{
				Email: "alex@mail.ru",
			},
			Method:              "POST",
			ExpectedRequestBody: "invalid payment status\n",
			ExpectedStatusCode:  400,
			MockUser: func(s *mock_service.MockUser, payId int, in models.PaymentProcessingInput) {
				s.EXPECT().Verification(payId, in.Email).Return(true, nil)
			},
			MockPay: func(s *mock_service.MockPayment, payId int) {
				s.EXPECT().PaymentProcessing(payId).Return("", errors.New("invalid payment status"))
			},
		},
	}
	for tName, tCase := range tData {
		v := tCase
		t.Run(tName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			user := mock_service.NewMockUser(c)
			v.MockUser(user, 1, v.Input)
			pay := mock_service.NewMockPayment(c)
			v.MockPay(pay, 1)
			services := service.Services{
				Payment: pay,
				User:    user,
			}
			handler := NewHandler(&services)
			r := http.HandlerFunc(handler.PaymentProcessing)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(v.Method, v.URL, bytes.NewBufferString(v.InputBody))
			r.ServeHTTP(w, req)
			assert.Equal(t, v.ExpectedRequestBody, w.Body.String())
			assert.Equal(t, v.ExpectedStatusCode, w.Code)
		})
	}
}

func TestByUserEmail(t *testing.T) {
	type mockPay func(s *mock_service.MockPayment, in models.InputByUserEmail)
	tData := map[string]struct {
		InputBody           string
		Input               models.InputByUserEmail
		Method              string
		ExpectedRequestBody string
		ExpectedStatusCode  int
		MockPay             mockPay
	}{
		"Success": {
			InputBody: `{"Email":"ann@mail.ru"}`,
			Input: models.InputByUserEmail{
				Email: "ann@mail.ru",
			},
			Method:              "GET",
			ExpectedRequestBody: `[{"ID":114,"UserID":1,"Email":"ann@mail.ru","Sum":1000,"Currency":"KZ","CreationDate":"2022-06-11T18:45:47.72474801+06:00","ChangeDate":"2022-06-11T18:47:22.683292944+06:00","Status":"SUCCESS"}]`,
			ExpectedStatusCode:  200,
			MockPay: func(s *mock_service.MockPayment, in models.InputByUserEmail) {
				s.EXPECT().ByUserEmail(in.Email).Return([]models.Transaction{
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
		"Invalid method": {
			Method:              "POST",
			ExpectedRequestBody: "method not allowed\n",
			ExpectedStatusCode:  405,
			MockPay:             func(s *mock_service.MockPayment, in models.InputByUserEmail) {},
		},
		"Payments not found": {
			InputBody: `{"Email":"ann@mail.ru"}`,
			Input: models.InputByUserEmail{
				Email: "ann@mail.ru",
			},
			Method:              "GET",
			ExpectedRequestBody: "not found payments\n",
			ExpectedStatusCode:  400,
			MockPay: func(s *mock_service.MockPayment, in models.InputByUserEmail) {
				s.EXPECT().ByUserEmail(in.Email).Return(nil, errors.New("not found payments"))
			},
		},
	}
	for tName, tCase := range tData {
		v := tCase
		t.Run(tName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			pay := mock_service.NewMockPayment(c)
			v.MockPay(pay, v.Input)
			services := service.Services{
				Payment: pay,
			}
			handler := NewHandler(&services)
			r := http.HandlerFunc(handler.ByUserEmail)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(v.Method, "/payments/byemail", bytes.NewBufferString(v.InputBody))
			r.ServeHTTP(w, req)
			assert.Equal(t, v.ExpectedRequestBody, w.Body.String())
			assert.Equal(t, v.ExpectedStatusCode, w.Code)
		})
	}
}

func TestCancelPayment(t *testing.T) {
	type mockPay func(s *mock_service.MockPayment, payId int)
	tData := map[string]struct {
		URL                 string
		Method              string
		ExpectedRequestBody string
		ExpectedStatusCode  int
		MockPay             mockPay
	}{
		"Success": {
			URL:                 "/payments/cancel/1",
			Method:              "POST",
			ExpectedRequestBody: "Done",
			ExpectedStatusCode:  200,
			MockPay: func(s *mock_service.MockPayment, payId int) {
				s.EXPECT().CancelPayment(payId).Return(nil)
			},
		},
		"Invalid payment status": {
			URL:                 "/payments/cancel/1",
			Method:              "POST",
			ExpectedRequestBody: "invalid status\n",
			ExpectedStatusCode:  400,
			MockPay: func(s *mock_service.MockPayment, payId int) {
				s.EXPECT().CancelPayment(payId).Return(errors.New("invalid status"))
			},
		},
		"method not allowed": {
			URL:                 "/payments/cancel/1",
			Method:              "GET",
			ExpectedRequestBody: "method not allowed\n",
			ExpectedStatusCode:  405,
			MockPay:             func(s *mock_service.MockPayment, payId int) {},
		},
	}
	for tName, tCase := range tData {
		v := tCase
		t.Run(tName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			pay := mock_service.NewMockPayment(c)
			v.MockPay(pay, 1)
			services := service.Services{
				Payment: pay,
			}
			handler := NewHandler(&services)
			r := http.HandlerFunc(handler.CancelPayment)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(v.Method, v.URL, nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, v.ExpectedRequestBody, w.Body.String())
			assert.Equal(t, v.ExpectedStatusCode, w.Code)
		})
	}
}

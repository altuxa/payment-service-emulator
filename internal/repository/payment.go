package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/altuxa/payment-service-emulator/internal/models"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{
		db: db,
	}
}

func (p *PaymentRepo) NewPayment(id int, email string, sum int, val string, status string) (int, error) {
	stmt, err := p.db.Prepare("INSERT INTO Transactions(UserID, UserEmail,Sum,Currency,CreationDate,ChangeDate,Status)VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	date := time.Now()
	res, err := stmt.Exec(id, email, sum, val, date, date, status)
	if err != nil {
		return 0, err
	}
	paymentID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if paymentID == 0 {
		return 0, errors.New("payment not create")
	}
	defer stmt.Close()
	return int(paymentID), nil
}

func (p *PaymentRepo) PaymentStatus(paymentId int) (string, error) {
	status := ""
	stmt, err := p.db.Prepare("SELECT Status FROM Transactions WHERE ID = ?")
	if err != nil {
		return "", err
	}
	row := stmt.QueryRow(paymentId)
	row.Scan(&status)
	if len(status) == 0 {
		return "", errors.New("payment not found")
	}
	defer stmt.Close()
	return status, nil
}

func (p *PaymentRepo) GetAllPaymentsByUserID(userId int) ([]models.Transaction, error) {
	payments := []models.Transaction{}
	row, err := p.db.Query("SELECT ID,UserID, UserEmail,Sum,Currency,CreationDate,ChangeDate,Status FROM Transactions WHERE UserID = ?", userId)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		payment := models.Transaction{}
		err := row.Scan(&payment.ID, &payment.UserID, &payment.UserEmail, &payment.Sum, &payment.Valute, &payment.CreationDate, &payment.ChangeDate, &payment.Status)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	defer row.Close()
	return payments, nil
}

func (p *PaymentRepo) GetAllPaymentsByEmail(email string) ([]models.Transaction, error) {
	payments := []models.Transaction{}
	row, err := p.db.Query("SELECT ID,UserID, UserEmail,Sum,Currency,CreationDate,ChangeDate,Status FROM Transactions WHERE UserEmail = ?", email)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		payment := models.Transaction{}
		err := row.Scan(&payment.ID, &payment.UserID, &payment.UserEmail, &payment.Sum, &payment.Valute, &payment.CreationDate, &payment.ChangeDate, &payment.Status)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	defer row.Close()
	return payments, nil
}

func (p *PaymentRepo) DeletePayment(paymentId int) error {
	_, err := p.db.Exec("DELETE FROM Transactions WHERE ID = ?", paymentId)
	if err != nil {
		return err
	}
	return nil
}

func (p *PaymentRepo) SetStatusSuccess(paymentId int) error {
	_, err := p.db.Exec("UPDATE Transactions Set Status  = ? WHERE ID = ?", models.StatusSuccess, paymentId)
	if err != nil {
		return err
	}
	return nil
}

func (p *PaymentRepo) SetStatusFail(paymentId int) error {
	_, err := p.db.Exec("UPDATE Transactions Set Status  = ? WHERE ID = ?", models.StatusFail, paymentId)
	if err != nil {
		return err
	}
	return nil
}

// func (p *PaymentRepo) SetStatusError(paymentId int) error {
// 	_, err := p.db.Exec("UPDATE Transactions Set Status  = ?", models.StatusError)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (p *PaymentRepo) PaymentErrorImitation() bool {
// 	rand.Seed(time.Now().UnixNano())
// 	a := rand.Intn(60)
// 	b := rand.Intn(45)
// 	return a > b
// }

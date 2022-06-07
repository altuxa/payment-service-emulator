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

func (p *PaymentRepo) NewPayment(id int, email string, sum int, val string) error {
	stmt, err := p.db.Prepare("INSERT INTO Transactions(UserID, UserEmail,Sum,Valute,CreationDate,ChangeDate,Status)VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	date := time.Now()
	res, err := stmt.Exec(id, email, sum, val, date, date, "NEW")
	if err != nil {
		return err
	}
	paymentID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	if paymentID == 0 {
		return errors.New("payment not create")
	}
	defer stmt.Close()
	return nil
}

func (p *PaymentRepo) PaymentStatus(paymentId int) (string, error) {
	status := ""
	stmt, err := p.db.Prepare("SELECT Status FROM Transactions WHERE ID = ?")
	if err != nil {
		return "", err
	}
	row := stmt.QueryRow(paymentId)
	row.Scan(&status)
	defer stmt.Close()
	return status, nil
}

func (p *PaymentRepo) GetAllPaymentsByUserID(userId int) ([]models.Transaction, error) {
	payments := []models.Transaction{}
	// stmt, err := p.db.Prepare("SELECT ID,UserID, UserEmail,Sum,Valute,CreationDate,ChangeDate,Status FROM Transactions WHERE UserID = ?")
	// if err != nil {
	// 	return nil, err
	// }
	row, err := p.db.Query("SELECT ID,UserID, UserEmail,Sum,Valute,CreationDate,ChangeDate,Status FROM Transactions WHERE UserID = ?", userId)
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
	return payments, nil
}

func (p *PaymentRepo) GetAllPaymentsByEmail(email string) ([]models.Transaction, error) {
	payments := []models.Transaction{}
	// stmt, err := p.db.Prepare("SELECT ID,UserID, UserEmail,Sum,Valute,CreationDate,ChangeDate,Status FROM Transactions WHERE UserID = ?")
	// if err != nil {
	// 	return nil, err
	// }
	row, err := p.db.Query("SELECT ID,UserID, UserEmail,Sum,Valute,CreationDate,ChangeDate,Status FROM Transactions WHERE UserEmail = ?", email)
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
	return payments, nil
}

func (p *PaymentRepo) CancelPayment(paymentId int) {
}

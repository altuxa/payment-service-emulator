package repository

import (
	"database/sql"
	"errors"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) UserVerification(paymentID int, email string) (string, error) {
	var res string
	stmt, err := u.db.Prepare("SELECT UserEmail FROM Transactions WHERE ID = ? AND UserEmail = ?")
	if err != nil {
		return "", err
	}
	row := stmt.QueryRow(paymentID, email)
	row.Scan(&res)
	if res == "" {
		return "", errors.New("not found")
	}
	defer stmt.Close()
	return res, nil
}

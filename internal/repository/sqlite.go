package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./sqlite3.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE "Transactions" (
		"ID"	INTEGER NOT NULL UNIQUE,
		"UserID"	INTEGER,
		"UserEmail"	TEXT,
		"Sum"	REAL,
		"Currency"	TEXT,
		"CreationDate"	DATETIME NOT NULL,
		"ChangeDate"	DATETIME NOT NULL,
		"Status"	TEXT,
		PRIMARY KEY("ID" AUTOINCREMENT)
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE "Users" (
		"ID"	INTEGER NOT NULL UNIQUE,
		"Email"	TEXT,
		PRIMARY KEY("ID" AUTOINCREMENT)
	)`)
	return err
}

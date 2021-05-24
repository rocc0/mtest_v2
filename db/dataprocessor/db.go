package dataprocessor

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewService(db *sql.DB) (*Service, error) {
	return &Service{db: db}, nil
}

func ConnectToSQL(address string) (*sql.DB, error) {
	if address == "" {
		address = "root:password@tcp(localhost:3306)/mtest"
	}

	db, err := sql.Open("mysql", address)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

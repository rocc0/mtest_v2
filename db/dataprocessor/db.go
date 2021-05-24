package dataprocessor

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewService(db *sql.DB) (*Service, error) {
	return &Service{db: db}, nil
}

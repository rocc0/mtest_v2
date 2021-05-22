package dataprocessor

import (
	"database/sql"
)

func NewService(db *sql.DB) (*Service, error) {
	return &Service{db: db}, nil
}

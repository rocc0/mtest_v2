package main

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/mtest")
	if err != nil {
		log.Error(err)
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
	}
}

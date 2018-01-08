package main

import (

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

var db *sql.DB


func init() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(192.168.99.100:3306)/mtest")
	check(err)

	err = db.Ping()

}
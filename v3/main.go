package main

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"mtest.com.ua/v3/db/dataprocessor"
	hashpkg "mtest.com.ua/v3/db/hasher"
	handlerspkg "mtest.com.ua/v3/handlers"
	routes "mtest.com.ua/v3/router"
	searchpkg "mtest.com.ua/v3/search"
)

func main() {
	hash, err := hashpkg.NewHashHandler("")
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := connectToSQL("")
	if err != nil {
		logrus.Fatal(err)
	}

	data, err := dataprocessor.NewService(db)
	if err != nil {
		logrus.Fatal(err)
	}

	searchService, err := searchpkg.NewService(db)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := searchService.Connect(""); err != nil {
		logrus.Fatal(err)
	}

	router := routes.NewRouter(handlerspkg.NewService(data, hash, searchService), data)
	if err := router.Init(); err != nil {
		logrus.Fatal(err)
	}

	if err := data.InitUsersTable(); err != nil {
		logrus.Fatal(err)
	}

	if err := searchService.Init(); err != nil {
		logrus.Fatal(err)
	}

	// Start serving the application
	router.Run(":80")
}

func connectToSQL(address string) (*sql.DB, error) {
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

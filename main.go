package main

import (
	"github.com/sirupsen/logrus"

	"mtest.com.ua/config"
	datapkg "mtest.com.ua/db/dataprocessor"
	hashpkg "mtest.com.ua/db/hasher"
	handlerspkg "mtest.com.ua/handlers"
	routes "mtest.com.ua/router"
	searchpkg "mtest.com.ua/search"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := datapkg.ConnectToSQL(cfg.DatabaseURL)
	if err != nil {
		logrus.Fatal(err)
	}

	hash, err := hashpkg.NewHashHandler(db)
	if err != nil {
		logrus.Error(err)
	}

	if err := hash.Init(); err != nil {
		logrus.Fatal(err)
	}

	data, err := datapkg.NewService(db)
	if err != nil {
		logrus.Fatal(err)
	}

	searchService, err := searchpkg.NewService(db)
	if err != nil {
		logrus.Error(err)
	}

	if err := searchService.Connect(cfg.ElasticURL); err != nil {
		logrus.Error(err)
	} else {
		if err := searchService.Init(); err != nil {
			logrus.Error(err)
		}
	}

	router := routes.NewRouter(handlerspkg.NewService(data, hash, searchService), data)
	if err := router.Init(); err != nil {
		logrus.Fatal(err)
	}

	if err := data.InitUsersTable(); err != nil {
		logrus.Fatal(err)
	}

	// Start serving the application
	router.Run(":80")
}

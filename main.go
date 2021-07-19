package main

import (
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"mtest.com.ua/mail"

	"github.com/gin-gonic/autotls"
	"mtest.com.ua/config"
	datapkg "mtest.com.ua/db/dataprocessor"
	hashpkg "mtest.com.ua/db/hasher"
	handlerspkg "mtest.com.ua/handlers"
	routes "mtest.com.ua/router"
	searchpkg "mtest.com.ua/search"
)

func main() {
	formatter := runtime.Formatter{ChildFormatter: &logrus.TextFormatter{
		FullTimestamp: true,
	}}
	formatter.Line = true
	formatter.File = true
	logrus.SetFormatter(&formatter)
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
		logrus.Fatal(err)
	}

	if err := hash.Init(); err != nil {
		logrus.Fatal(err)
	}

	data, err := datapkg.NewService(db)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := data.Init(); err != nil {
		logrus.Fatal(err)
	}
	searchService, err := searchpkg.NewService(db)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := searchService.Connect(cfg.ElasticURL); err != nil {
		logrus.Error("Elastic connect:", err)
	} else {
		if err := searchService.Init(); err != nil {
			logrus.Error("Elastic init:", err)
		} else {
			if err := searchService.ElasticIndex(); err != nil {
				logrus.Errorf("Elastic index: %v", err)
			}
		}
	}

	router := routes.NewRouter(handlerspkg.NewService(data, hash, searchService, mail.Auth{
		Email:    cfg.Email,
		Password: cfg.Password,
	}), data)
	if err := router.Init(); err != nil {
		logrus.Fatal(err)
	}

	if err := data.InitUsersTable(); err != nil {
		logrus.Fatal(err)
	}

	// Start serving the application
	//router.Run(":8099")
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("mtest.org.ua"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	autotls.RunWithManager(router, &m)
}

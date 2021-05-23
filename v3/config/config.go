package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"

	validatorpkg "github.com/go-playground/validator"
)

type Config struct {
	ElasticURL  string `env:"ELASTIC_URL" validate:"required,url"`
	PostgresURL string `env:"POSTGRES_URL" validate:"required,url"`
	MongoURL    string `env:"MONGO_URL" validate:"required,url"`
}

func FromEnv() (*Config, error) {
	var conf Config
	err := env.Parse(&conf)
	if err != nil {
		return nil, errors.Wrap(err, "config loading")
	}
	validator := validatorpkg.New()
	err = validator.Struct(conf)
	if err != nil {
		return nil, errors.Wrap(err, "config validation")
	}
	return &conf, nil
}

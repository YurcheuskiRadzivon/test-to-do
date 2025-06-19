package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP       HTTP
		PG         PG
		APP        APP
		JWT        JWT
		ADMIN      ADMIN
		LOCALSTACK LOCALSTACK
	}

	HTTP struct {
		PORT string `env:"HTTP_PORT,required"`
	}

	PG struct {
		URL string `env:"PG_URL,required"`
	}

	APP struct {
		DOMAIN string `env:"APP_DOMAIN,required"`
	}

	JWT struct {
		SECRETKEY string `env:"JWT_SECRET_KEY,required"`
	}

	ADMIN struct {
		ID int `env:"ADMIN_ID,required"`
	}

	LOCALSTACK struct {
		EXTERNAL_ENDPOINT string `env:"LOCALSTACK_ENDPOINT_EXTERNAL,required"`
		INTERNAL_ENDPOINT string `env:"LOCALSTACK_ENDPOINT_INTERNAL,required"`
		ACCESS_KEY        string `env:"LOCALSTACK_ACCESS_KEY,required"`
		SECRET_KEY        string `env:"LOCALSTACK_SECRET_KEY,required"`
		BUCKET            string `env:"LOCALSTACK_BUCKET,required"`
	}
)

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

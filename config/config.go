package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP  HTTP
		PG    PG
		APP   APP
		JWT   JWT
		ADMIN ADMIN
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
		SECRETKEY string `env:"SECRET_KEY,required"`
	}

	ADMIN struct {
		ID int `env:"ADMIN_ID,required"`
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

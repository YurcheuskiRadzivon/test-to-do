package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const (
	StorageMinio      = "MINIO"
	StorageLocalstack = "LOCALSTACK"
	StorageFS         = "FS"
	DefaultStorage    = "FS"
)

type (
	Config struct {
		HTTP            HTTP
		PG              PG
		APP             APP
		JWT             JWT
		ADMIN           ADMIN
		LOCALSTACK      LOCALSTACK
		FSSTORAGE       FSSTORAGE
		MINIO           MINIO
		STORAGESWITCHER STORAGESWITCHER
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

	STORAGESWITCHER struct {
		STORAGE string `env:"STORAGE,required"`
	}

	LOCALSTACK struct {
		EXTERNAL_ENDPOINT string `env:"LOCALSTACK_ENDPOINT_EXTERNAL,required"`
		INTERNAL_ENDPOINT string `env:"LOCALSTACK_ENDPOINT_INTERNAL,required"`
		ACCESS_KEY        string `env:"LOCALSTACK_ACCESS_KEY,required"`
		SECRET_KEY        string `env:"LOCALSTACK_SECRET_KEY,required"`
		BUCKET            string `env:"LOCALSTACK_BUCKET,required"`
	}

	MINIO struct {
		EXTERNAL_ENDPOINT string `env:"MINIO_ENDPOINT_EXTERNAL,required"`
		INTERNAL_ENDPOINT string `env:"MINIO_ENDPOINT_INTERNAL,required"`
		ACCESS_KEY        string `env:"MINIO_ACCESS_KEY,required"`
		SECRET_KEY        string `env:"MINIO_SECRET_KEY,required"`
		BUCKET            string `env:"MINIO_BUCKET,required"`
	}

	FSSTORAGE struct {
		EXTERNAL_ENDPOINT string `env:"FSS_ENDPOINT_EXTERNAL,required"`
		PATH              string `env:"FSS_PATH,required"`
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

	cfg.CheckStorageSwitcher()

	return cfg, nil
}

func (c *Config) CheckStorageSwitcher() {
	switch c.STORAGESWITCHER.STORAGE {
	case StorageFS, StorageMinio, StorageLocalstack:
		log.Printf("Check storage switcher - %v", c.STORAGESWITCHER.STORAGE)
	default:
		c.STORAGESWITCHER.STORAGE = DefaultStorage
		log.Printf("Invalis Storage switcher, default check storage switcher - %v", c.STORAGESWITCHER.STORAGE)
	}
}

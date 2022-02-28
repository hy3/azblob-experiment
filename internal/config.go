package internal

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	LogLevel string `env:"LOGLEVEL"`

	// to indicate specification clearly
	TenantID     string `env:"AZURE_TENANT_ID"`
	ClientID     string `env:"AZURE_CLIENT_ID"`
	ClientSecret string `env:"AZURE_CLIENT_SECRET"`

	AccountName   string `env:"AZBLOB_ACCOUNT_NAME"`
	ContainerName string `env:"AZBLOB_CONTAINER_NAME"`
	BlobName      string `env:"AZBLOB_BLOB_NAME"`
}

func InitConfig() (Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}

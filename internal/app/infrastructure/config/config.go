package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"web-studio-backend/internal/app/infrastructure/logger"
)

var Config struct {
	App struct {
		Jwt struct {
			AccessExpirationMinutes int16  `yaml:"access_expiration_minutes"`
			RefreshExpirationDays   int16  `yaml:"refresh_expiration_days"`
			AccessTokenSecretKey    string `yaml:"access_token_secret_key"`
			RefreshTokenSecretKey   string `yaml:"refresh_token_secret_key"`
		} `yaml:"jwt"`
	} `yaml:"app"`
	Interfaces struct {
		Https struct {
			Port         uint32 `yaml:"port"`
			Host         string `yaml:"host"`
			KeyFilePath  string `yaml:"key_file_path"`
			CertFilePath string `yaml:"cert_file_path"`
		} `yaml:"https"`
		Http struct {
			Port uint32 `yaml:"port"`
			Host string `yaml:"host"`
		} `yaml:"http"`
	} `yaml:"interfaces"`
	Logger             logger.Config `yaml:"logger"`
	DatabaseConnString string        `yaml:"database_conn_string"`
}

func Init(configPath string) error {
	err := cleanenv.ReadConfig(configPath, &Config)
	return err
}

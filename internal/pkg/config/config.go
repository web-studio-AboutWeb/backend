package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"env"`
	App struct {
		Jwt struct {
			AccessTokenExpMinutes int16  `yaml:"access_token_exp_minutes"`
			RefreshTokenExpDays   int16  `yaml:"refresh_token_exp_days"`
			AccessTokenSecretKey  string `yaml:"access_token_secret_key"`
			RefreshTokenSecretKey string `yaml:"refresh_token_secret_key"`
		} `yaml:"jwt"`
	} `yaml:"app"`
	Server struct {
		Http struct {
			Port uint32 `yaml:"port"`
			Host string `yaml:"host"`
		} `yaml:"http"`
		Https struct {
			Port         uint32 `yaml:"port"`
			Host         string `yaml:"host"`
			KeyFilePath  string `yaml:"key_file_path"`
			CertFilePath string `yaml:"cert_file_path"`
		} `yaml:"https"`
	}
	Database struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"database"`
}

var (
	cfg  Config
	once sync.Once
)

func Read(configPath string) {
	once.Do(func() {
		err := cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}
	})
}

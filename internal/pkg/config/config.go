package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Env string `yaml:"env" env-default:"dev"`
		Key string `yaml:"key" env-required:"true"`
		Jwt struct {
			AccessTokenExpMinutes int16  `yaml:"access_token_exp_minutes"`
			RefreshTokenExpDays   int16  `yaml:"refresh_token_exp_days"`
			AccessTokenSecretKey  string `yaml:"access_token_secret_key"`
			RefreshTokenSecretKey string `yaml:"refresh_token_secret_key"`
		} `yaml:"jwt"`
	} `yaml:"app" env-required:"true"`
	Http struct {
		HttpsEnabled bool   `yaml:"https_enabled" env-default:"false"`
		Port         uint16 `yaml:"port" env-default:"8080"`
		Host         string `yaml:"host"`
		KeyFilePath  string `yaml:"key_file_path"`
		CertFilePath string `yaml:"cert_file_path"`
	} `yaml:"http"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Database string `yaml:"database"`
	} `yaml:"database" env-required:"true"`
}

var (
	cfg  Config
	once sync.Once

	// Block is needed to encode/decode sensitive data.
	Block cipher.Block
)

func Read(configPath string) {
	once.Do(func() {
		err := cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		key, err := hex.DecodeString(cfg.App.Key)
		if err != nil {
			log.Fatalf("Failed to decode app key: %v", err)
		}

		Block, err = aes.NewCipher(key)
		if err != nil {
			log.Fatalf("Failed to create cipher block: %v", err)
		}
	})
}

func Get() *Config {
	return &cfg
}

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

	k     = "bd53b9d5f1d53318e322f0a4cbb225972781ef66fb840b309e2a1951edc3abfb"
	Block cipher.Block // Block is needed to encode/decode sensitive data.
)

func Read(configPath string) {
	once.Do(func() {
		err := cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		decodedKey, err := hex.DecodeString(k)
		if err != nil {
			log.Fatalf("Failed to decode app key: %v", err)
		}

		Block, err = aes.NewCipher(decodedKey)
		if err != nil {
			log.Fatalf("Failed to create cipher block: %v", err)
		}
	})
}

func Get() *Config {
	return &cfg
}

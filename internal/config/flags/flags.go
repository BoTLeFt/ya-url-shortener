package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AddressForGin   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() *Config {
	return &Config{
		AddressForGin:   "",
		BaseURL:         "",
		FileStoragePath: "",
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.AddressForGin, "a", "localhost:8080", "port for gin")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base URL")
	flag.StringVar(&c.FileStoragePath, "f", "file_storage.json", "path to store data")
	flag.Parse()

	// cleanenv читает переменные окружения и перезаписывает значения флагов, если они установлены.
	// Иначе используются значения из флагов
	envConfig := &Config{}
	if err := cleanenv.ReadEnv(envConfig); err == nil {
		if envConfig.AddressForGin != "" {
			c.AddressForGin = envConfig.AddressForGin
		}
		if envConfig.BaseURL != "" {
			c.BaseURL = envConfig.BaseURL
		}
		if envConfig.FileStoragePath != "" {
			c.FileStoragePath = envConfig.FileStoragePath
		}
	}
}

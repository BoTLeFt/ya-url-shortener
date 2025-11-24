package config

import (
	"flag"
	"os"
)

type Config struct {
	AddressForGin   string
	BaseURL         string
	FileStoragePath string
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

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		c.AddressForGin = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		c.BaseURL = envBaseURL
	}
	if envFileStorage := os.Getenv("FILE_STORAGE_PATH"); envFileStorage != "" {
		c.FileStoragePath = envFileStorage
	}
}

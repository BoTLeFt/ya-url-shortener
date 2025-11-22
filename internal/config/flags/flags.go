package config

import (
	"flag"
	"os"
)

type Config struct {
	AddressForGin string
	BaseURL       string
}

func NewConfig() *Config {
	return &Config{
		AddressForGin: "",
		BaseURL:       "",
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.AddressForGin, "a", "localhost:8080", "port for gin")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base URL")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		c.AddressForGin = envRunAddr
	}
	if envBaseUrl := os.Getenv("BASE_URL"); envBaseUrl != "" {
		c.BaseURL = envBaseUrl
	}
}

package main

import (
	"log"
	"net/http"

	flags "github.com/BoTLeFt/ya-url-shortener/internal/config/flags"
	"github.com/BoTLeFt/ya-url-shortener/internal/handler"
	"github.com/BoTLeFt/ya-url-shortener/internal/logger"
	"github.com/BoTLeFt/ya-url-shortener/internal/repository/memory"
	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
)

func main() {
	if err := logger.Initialize("info"); err != nil {
		log.Fatal("Failed to initialize logger: ", err)
	}

	config := flags.NewConfig()
	config.ParseFlags()
	repo := memory.New()
	svc := shortener.New(repo)
	router := handler.NewRouter(svc, config)

	log.Println("AddressForGin: ", config.AddressForGin)
	log.Println("BaseURL: ", config.BaseURL)
	if err := http.ListenAndServe(config.AddressForGin, router); err != nil {
		log.Fatal(err)
	}
}

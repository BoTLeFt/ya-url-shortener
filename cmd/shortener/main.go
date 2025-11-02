package main

import (
	"log"
	"net/http"

	flags "github.com/BoTLeFt/ya-url-shortener/internal/config/flags"
	"github.com/BoTLeFt/ya-url-shortener/internal/handler"
	"github.com/BoTLeFt/ya-url-shortener/internal/repository/memory"
	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
)

func main() {
	config := flags.NewConfig()
	config.ParseFlags()
	repo := memory.New()
	svc := shortener.New(repo)
	router := handler.NewRouter(svc, config)

	if err := http.ListenAndServe(config.AddressForGin, router); err != nil {
		log.Fatal(err)
	}
}

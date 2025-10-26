package main

import (
	"log"
	"net/http"

	"github.com/BoTLeFt/ya-url-shortener/internal/handler"
	"github.com/BoTLeFt/ya-url-shortener/internal/repository/memory"
	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
)

func main() {
	repo := memory.New()
	svc := shortener.New(repo)
	router := handler.NewRouter(svc)

	addr := ":8080"
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

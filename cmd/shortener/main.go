package main

import (
	"log"
	"net/http"

	flags "github.com/BoTLeFt/ya-url-shortener/internal/config/flags"
	"github.com/BoTLeFt/ya-url-shortener/internal/handler"
	"github.com/BoTLeFt/ya-url-shortener/internal/logger"
	"github.com/BoTLeFt/ya-url-shortener/internal/repository/file"
	"github.com/BoTLeFt/ya-url-shortener/internal/repository/memory"
	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
)

func uploadFromFile(consumer *file.Consumer, memory *memory.Memory) error {
	for {
		var event *file.ShortenedURL
		event, err := consumer.ReadEvent()
		if err != nil {
			return err
		}
		if event == nil {
			return nil
		}
		memory.Save(event.ShortURL, event.OriginalURL, false)
	}
}

func main() {
	if err := logger.Initialize("info"); err != nil {
		log.Fatal("Failed to initialize logger: ", err)
	}

	config := flags.NewConfig()
	config.ParseFlags()

	consumer, err := file.NewConsumer(config.FileStoragePath)
	if err != nil {
		log.Fatal("No file")
	}
	producer, err := file.NewProducer(config.FileStoragePath)
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	err = uploadFromFile(consumer, repo)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = consumer.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	svc := shortener.New(repo)
	router := handler.NewRouter(svc, config)

	log.Println("AddressForGin: ", config.AddressForGin)
	log.Println("BaseURL: ", config.BaseURL)
	if err := http.ListenAndServe(config.AddressForGin, router); err != nil {
		log.Fatal(err)
	}
}

package handler

import (
	"net/http"
	"strings"

	config "github.com/BoTLeFt/ya-url-shortener/internal/config/flags"
	"github.com/BoTLeFt/ya-url-shortener/internal/encoding"
	"github.com/BoTLeFt/ya-url-shortener/internal/logger"
	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
	"github.com/gin-gonic/gin"
)

func NewRouter(svc *shortener.Service, cfg *config.Config) http.Handler {
	router := gin.Default()

	router.Use(encoding.GzipMiddleware())
	router.Use(logger.RequestLogger())

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "bad request")
	})
	router.NoMethod(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "bad request")
	})

	createHandler := NewCreateHandler(svc, cfg)
	redirectHandler := NewRedirectHandler(svc)
	shortenHandler := NewShortenHandler(svc, cfg)

	router.POST("/", createHandler.Create)
	router.POST("/api/shorten", shortenHandler.Shorten)
	router.GET("/:id", redirectHandler.Redirect)

	return router
}

// isValidIDPath проверяет, что путь соответствует формату "/<id>" с корректным ID
func isValidIDPath(path string) bool {
	if path == "" || path == "/" {
		return false
	}
	id := strings.TrimPrefix(path, "/")
	// Проверка на отсутствие слэшей и соответствие формату ID из конфига
	return !strings.Contains(id, "/") && shortener.IsValidID(id)
}

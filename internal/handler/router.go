package handler

import (
	"net/http"
	"strings"

	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
	"github.com/gin-gonic/gin"
)

func NewRouter(svc *shortener.Service) http.Handler {
	router := gin.New()

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "bad request")
	})
	router.NoMethod(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "bad request")
	})

	createHandler := NewCreateHandler(svc)
	redirectHandler := NewRedirectHandler(svc)

	router.POST("/", gin.WrapH(createHandler))
	router.GET("/:id", gin.WrapH(redirectHandler))

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

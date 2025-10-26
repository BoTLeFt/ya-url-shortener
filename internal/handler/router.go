package handler

import (
	"net/http"
	"strings"

	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
)

type Router struct {
	create   *CreateHandler
	redirect *RedirectHandler
}

func NewRouter(svc *shortener.Service) *Router {
	return &Router{
		create:   NewCreateHandler(svc),
		redirect: NewRedirectHandler(svc),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodPost && req.URL.Path == "/":
		r.create.ServeHTTP(w, req)
	case req.Method == http.MethodGet && isValidIDPath(req.URL.Path):
		r.redirect.ServeHTTP(w, req)
	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
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

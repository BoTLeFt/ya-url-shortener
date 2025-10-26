package handler

import (
	"net/http"
	"strings"

	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
)

type RedirectHandler struct {
	svc *shortener.Service
}

func NewRedirectHandler(svc *shortener.Service) *RedirectHandler {
	return &RedirectHandler{svc: svc}
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" || strings.Contains(id, "/") {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	orig, ok := h.svc.Resolve(id)
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", orig)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

package handler

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
)

const maxBodyBytes = 8 << 10 // 8 KiB

type CreateHandler struct {
	svc *shortener.Service
}

func NewCreateHandler(svc *shortener.Service) *CreateHandler {
	return &CreateHandler{svc: svc}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isTextPlain(r.Header.Get("Content-Type")) {
		fmt.Println(r.Header.Get("Content-Type"))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	raw := strings.TrimSpace(string(body))
	if !isValidURL(raw) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := h.svc.Shorten(raw)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}
	shortURL := fmt.Sprintf("http://%s/%s", host, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func isTextPlain(ct string) bool {
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return false
	}
	return mt == "text/plain"
}

func isValidURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

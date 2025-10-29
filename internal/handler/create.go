package handler

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
	"github.com/gin-gonic/gin"
)

const maxBodyBytes = 8 << 10 // 8 KiB

type CreateHandler struct {
	svc *shortener.Service
}

func NewCreateHandler(svc *shortener.Service) *CreateHandler {
	return &CreateHandler{svc: svc}
}

func (h *CreateHandler) Create(c *gin.Context) {
	if !isTextPlain(c.GetHeader("Content-Type")) {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if c.Request == nil || c.Request.Body == nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	defer c.Request.Body.Close()

	limited := io.LimitReader(c.Request.Body, maxBodyBytes)
	body, err := io.ReadAll(limited)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	raw := strings.TrimSpace(string(body))
	if !isValidURL(raw) {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	id, err := h.svc.Shorten(raw)
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error")
		return
	}

	host := c.Request.Host
	if host == "" {
		host = "localhost:8080"
	}
	shortURL := fmt.Sprintf("http://%s/%s", host, id)

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, shortURL)
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

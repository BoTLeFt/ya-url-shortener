package handler

import (
	"net/http"
	"strings"

	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
	"github.com/gin-gonic/gin"
)

type RedirectHandler struct {
	svc *shortener.Service
}

func NewRedirectHandler(svc *shortener.Service) *RedirectHandler {
	return &RedirectHandler{svc: svc}
}

func (h *RedirectHandler) Redirect(c *gin.Context) {
	id := c.Param("id")
	if id == "" || strings.Contains(id, "/") || !shortener.IsValidID(id) {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	orig, ok := h.svc.Resolve(id)
	if !ok {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Header("Location", orig)
	c.Status(http.StatusTemporaryRedirect)
}

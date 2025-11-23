package handler

import (
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"

	config "github.com/BoTLeFt/ya-url-shortener/internal/config/flags"
	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortenHandler struct {
	svc *shortener.Service
	cfg *config.Config
}

func NewShortenHandler(svc *shortener.Service, cfg *config.Config) *ShortenHandler {
	return &ShortenHandler{svc: svc, cfg: cfg}
}

func (h *ShortenHandler) Shorten(c *gin.Context) {
	if !isApplicationJSON(c.GetHeader("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if c.Request == nil || c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	defer c.Request.Body.Close()

	c.Request.Body = io.NopCloser(io.LimitReader(c.Request.Body, maxBodyBytes))

	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if !isValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	id, err := h.svc.Shorten(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		log.Println(err.Error())
		return
	}

	shortURL, err := url.JoinPath(h.cfg.BaseURL, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		log.Println(err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, ShortenResponse{Result: shortURL})
}

func isApplicationJSON(ct string) bool {
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return false
	}
	return mt == "application/json"
}

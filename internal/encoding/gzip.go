package encoding

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	gzipEncoding = "gzip"
)

var shouldCompressTypes = []string{"application/json", "text/html"}

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Распаковка входящего запроса, если он сжат
		if strings.Contains(c.GetHeader("Content-Encoding"), gzipEncoding) {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
				return
			}
			defer reader.Close()
			c.Request.Body = io.NopCloser(reader)
		}

		// Проверяем, что клиент поддерживает gzip-сжатие и content-type имеет смысл для сжатия
		if !strings.Contains(c.GetHeader("Accept-Encoding"), gzipEncoding) || !slices.Contains(shouldCompressTypes, c.GetHeader("Content-Type")) {
			// если gzip не поддерживается, передаём управление дальше без изменений
			c.Next()
			return
		}

		gzw := &gzipResponseWriter{
			ResponseWriter: c.Writer,
		}
		c.Writer = gzw

		c.Next()

		if gzw.Writer != nil {
			gzw.Writer.Close()
		}
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if w.Writer == nil {
		gz, err := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			w.ResponseWriter.WriteHeader(statusCode)
			return
		}
		w.Writer = gz
		w.Header().Set("Content-Encoding", gzipEncoding)
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Writer == nil {
		gz, err := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			return 0, err
		}
		w.Writer = gz
		w.Header().Set("Content-Encoding", gzipEncoding)
	}
	return w.Writer.Write(b)
}

package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	config "github.com/BoTLeFt/ya-url-shortener/internal/config/flags"
	"github.com/BoTLeFt/ya-url-shortener/internal/repository/file"
	"github.com/BoTLeFt/ya-url-shortener/internal/repository/memory"
	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateHandler_Success(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	req, err := http.NewRequest("POST", "/", strings.NewReader("https://yandex.ru"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain")
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	responseBody := string(body)
	assert.True(t, strings.HasPrefix(responseBody, "http://localhost:8080/"))
	assert.Len(t, strings.TrimPrefix(responseBody, "http://localhost:8080/"), 7)
}

func TestCreateHandler_InvalidContentType(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	req, err := http.NewRequest("POST", "/", strings.NewReader("https://yandex.ru"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateHandler_InvalidURL(t *testing.T) {
	testCases := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "invalid URL format",
			body:    "not a valid url",
			wantErr: true,
		},
		{
			name:    "missing scheme",
			body:    "yandex.ru",
			wantErr: true,
		},
		{
			name:    "empty body",
			body:    "",
			wantErr: true,
		},
		{
			name:    "valid HTTP URL",
			body:    "http://yandex.ru",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL",
			body:    "https://yandex.ru",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			producer, err := file.NewProducer("test_file_storage.json")
			if err != nil {
				log.Fatal("No file")
			}
			repo := memory.New(*producer)
			svc := shortener.New(repo)
			cfg := config.NewConfig()
			cfg.BaseURL = "http://localhost:8080"
			router := NewRouter(svc, cfg)

			req, err := http.NewRequest("POST", "/", strings.NewReader(tc.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "text/plain")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tc.wantErr {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusCreated, resp.StatusCode)
			}
		})
	}
}

func TestCreateHandler_LargeBody(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	largeBody := bytes.NewBuffer(make([]byte, 10240))
	largeBody.WriteString("https://yandex.ru")

	req, err := http.NewRequest("POST", "/", largeBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateHandler_EmptyHost(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	req, err := http.NewRequest("POST", "/", strings.NewReader("https://yandex.ru"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	responseBody := string(body)
	assert.True(t, strings.HasPrefix(responseBody, "http://localhost:8080/"))
}

func TestCreateHandler_MultipleURLs(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
	}

	for i, originalURL := range urls {
		req, err := http.NewRequest("POST", "/", strings.NewReader(originalURL))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")
		req.Host = "localhost:8080"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Failed for URL #%d", i+1)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		shortURL := strings.TrimPrefix(string(body), "http://localhost:8080/")
		org, ok := svc.Resolve(shortURL)
		assert.True(t, ok, "URL #%d should be resolved", i+1)
		assert.Equal(t, originalURL, org, "URL #%d should resolve to original", i+1)
	}
}

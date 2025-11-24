package handler

import (
	"bytes"
	"encoding/json"
	"io"
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

func TestShortenHandler_Success(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		panic("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	reqBody := ShortenRequest{URL: "https://practicum.yandex.ru"}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response ShortenResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(response.Result, "http://localhost:8080/"))
	assert.Len(t, strings.TrimPrefix(response.Result, "http://localhost:8080/"), 7)
}

func TestShortenHandler_InvalidContentType(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		panic("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	reqBody := ShortenRequest{URL: "https://practicum.yandex.ru"}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShortenHandler_InvalidJSON(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		panic("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	req, err := http.NewRequest("POST", "/api/shorten", strings.NewReader("invalid json"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShortenHandler_InvalidURL(t *testing.T) {
	testCases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "invalid URL format",
			url:     "not a valid url",
			wantErr: true,
		},
		{
			name:    "missing scheme",
			url:     "yandex.ru",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "valid HTTP URL",
			url:     "http://yandex.ru",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL",
			url:     "https://practicum.yandex.ru",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			producer, err := file.NewProducer("test_file_storage.json")
			if err != nil {
				panic("No file")
			}
			repo := memory.New(*producer)
			svc := shortener.New(repo)
			cfg := config.NewConfig()
			cfg.BaseURL = "http://localhost:8080"
			router := NewRouter(svc, cfg)

			reqBody := ShortenRequest{URL: tc.url}
			bodyBytes, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/shorten", bytes.NewReader(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tc.wantErr {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusCreated, resp.StatusCode)
				assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var response ShortenResponse
				err = json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.True(t, strings.HasPrefix(response.Result, "http://localhost:8080/"))
			}
		})
	}
}

func TestShortenHandler_LargeBody(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		panic("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	largeBody := bytes.NewBuffer(make([]byte, 10240))
	reqBody := ShortenRequest{URL: "https://practicum.yandex.ru"}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)
	largeBody.Write(bodyBytes)

	req, err := http.NewRequest("POST", "/api/shorten", largeBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShortenHandler_MultipleURLs(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		panic("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	urls := []string{
		"https://practicum.yandex.ru",
		"https://google.com",
		"https://github.com",
	}

	for i, originalURL := range urls {
		reqBody := ShortenRequest{URL: originalURL}
		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/shorten", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Host = "localhost:8080"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Failed for URL #%d", i+1)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var response ShortenResponse
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)

		shortURL := strings.TrimPrefix(response.Result, "http://localhost:8080/")
		org, ok := svc.Resolve(shortURL)
		assert.True(t, ok, "URL #%d should be resolved", i+1)
		assert.Equal(t, originalURL, org, "URL #%d should resolve to original", i+1)
	}
}

func TestShortenHandler_MissingURLField(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		panic("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	reqBody := map[string]string{}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

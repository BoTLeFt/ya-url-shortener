package handler

import (
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

func TestRouter_PostCreate(t *testing.T) {
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
}

func TestRouter_GetRedirect(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	createReq, err := http.NewRequest("POST", "/", strings.NewReader("https://yandex.ru"))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "text/plain")

	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	createResp := createW.Result()
	defer createResp.Body.Close()

	body, err := io.ReadAll(createResp.Body)
	require.NoError(t, err)
	shortURL := strings.TrimPrefix(string(body), "http://localhost:8080/")

	redirectReq, err := http.NewRequest("GET", "/"+shortURL, nil)
	require.NoError(t, err)

	redirectW := httptest.NewRecorder()
	router.ServeHTTP(redirectW, redirectReq)

	redirectResp := redirectW.Result()
	defer redirectResp.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, redirectResp.StatusCode)
	assert.Equal(t, "https://yandex.ru", redirectResp.Header.Get("Location"))
}

func TestRouter_InvalidMethod(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	testCases := []struct {
		method string
		path   string
	}{
		{method: "PUT", path: "/"},
		{method: "DELETE", path: "/"},
		{method: "PATCH", path: "/"},
		{method: "POST", path: "/somepath"},
		{method: "GET", path: "/"},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, strings.NewReader(""))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tc.method == "POST" && tc.path != "/" || (tc.method != "POST" && tc.method != "GET") {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			} else if tc.method == "GET" && tc.path == "/" {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			}
		})
	}
}

func TestRouter_InvalidPathForRedirect(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	testCases := []struct {
		name string
		path string
	}{
		{name: "empty path", path: ""},
		{name: "root path", path: "/"},
		{name: "path with slashes", path: "/id1/id2"},
		{name: "invalid ID format", path: "/short"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestRouter_ValidIDPath(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	router := NewRouter(svc, cfg)

	createReq, err := http.NewRequest("POST", "/", strings.NewReader("https://yandex.ru"))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "text/plain")

	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	createResp := createW.Result()
	defer createResp.Body.Close()

	body, err := io.ReadAll(createResp.Body)
	require.NoError(t, err)
	shortURL := strings.TrimPrefix(string(body), "http://localhost:8080/")

	req, err := http.NewRequest("GET", "/"+shortURL, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
}

func TestIsValidIDPath(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "empty path", path: "", expected: false},
		{name: "root path", path: "/", expected: false},
		{name: "valid ID", path: "/AbCdEfG", expected: true},
		{name: "path with slash", path: "/id1/id2", expected: false},
		{name: "too short ID", path: "/short", expected: false},
		{name: "too long ID", path: "/verylongid", expected: false},
		{name: "ID with special chars", path: "/id@#$!", expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidIDPath(tc.path)
			assert.Equal(t, tc.expected, result, "Path: %s", tc.path)
		})
	}
}

func TestRouter_IntegrationFlow(t *testing.T) {
	producer, err := file.NewProducer("test_file_storage.json")
	if err != nil {
		log.Fatal("No file")
	}
	repo := memory.New(*producer)
	svc := shortener.New(repo)
	cfg := config.NewConfig()
	cfg.BaseURL = "http://yandex.ru"
	router := NewRouter(svc, cfg)

	originalURL := "https://yandex.ru/path?query=value"

	createReq, err := http.NewRequest("POST", "/", strings.NewReader(originalURL))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "text/plain")
	createReq.Host = "yandex.ru"

	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	createResp := createW.Result()
	defer createResp.Body.Close()

	assert.Equal(t, http.StatusCreated, createResp.StatusCode)

	body, err := io.ReadAll(createResp.Body)
	require.NoError(t, err)
	shortURL := strings.TrimPrefix(string(body), "http://yandex.ru/")
	assert.Len(t, shortURL, 7)

	redirectReq, err := http.NewRequest("GET", "/"+shortURL, nil)
	require.NoError(t, err)

	redirectW := httptest.NewRecorder()
	router.ServeHTTP(redirectW, redirectReq)

	redirectResp := redirectW.Result()
	defer redirectResp.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, redirectResp.StatusCode)
	assert.Equal(t, originalURL, redirectResp.Header.Get("Location"))
}

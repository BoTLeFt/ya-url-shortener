package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BoTLeFt/ya-url-shortener/internal/repository/memory"
	"github.com/BoTLeFt/ya-url-shortener/internal/service/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler_Success(t *testing.T) {
	repo := memory.New()
	svc := shortener.New(repo)
	handler := NewRedirectHandler(svc)

	originalURL := "https://yandex.ru"
	id, err := svc.Shorten(originalURL)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/"+id, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, originalURL, resp.Header.Get("Location"))
}

func TestRedirectHandler_NotFound(t *testing.T) {
	repo := memory.New()
	svc := shortener.New(repo)
	handler := NewRedirectHandler(svc)

	req, err := http.NewRequest("GET", "/nonexistent123", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRedirectHandler_InvalidIDFormat(t *testing.T) {
	repo := memory.New()
	svc := shortener.New(repo)
	handler := NewRedirectHandler(svc)

	testCases := []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "empty path",
			path:   "/",
			status: http.StatusBadRequest,
		},
		{
			name:   "path with slash",
			path:   "/id1/id2",
			status: http.StatusBadRequest,
		},
		{
			name:   "too short ID",
			path:   "/short",
			status: http.StatusBadRequest,
		},
		{
			name:   "too long ID",
			path:   "/verylongid123",
			status: http.StatusBadRequest,
		},
		{
			name:   "ID with special characters",
			path:   "/id@#$!",
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.status, resp.StatusCode, "Path: %s", tc.path)
		})
	}
}

func TestRedirectHandler_ValidIDNoRedirect(t *testing.T) {
	repo := memory.New()
	svc := shortener.New(repo)
	handler := NewRedirectHandler(svc)

	req, err := http.NewRequest("GET", "/AbCdEfG", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRedirectHandler_PathWithSlash(t *testing.T) {
	repo := memory.New()
	svc := shortener.New(repo)
	handler := NewRedirectHandler(svc)

	req, err := http.NewRequest("GET", "/some/path", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRedirectHandler_MultipleRedirects(t *testing.T) {
	repo := memory.New()
	svc := shortener.New(repo)
	handler := NewRedirectHandler(svc)

	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
	}

	ids := make([]string, len(urls))

	for i, url := range urls {
		id, err := svc.Shorten(url)
		require.NoError(t, err)
		ids[i] = id
	}

	for i, id := range ids {
		req, err := http.NewRequest("GET", "/"+id, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "Failed for redirect #%d", i+1)
		assert.Equal(t, urls[i], resp.Header.Get("Location"), "Wrong redirect for #%d", i+1)
	}
}

func TestRedirectHandler_WhitespaceInPath(t *testing.T) {
	repo := memory.New()
	svc := shortener.New(repo)
	handler := NewRedirectHandler(svc)

	req, err := http.NewRequest("GET", "/  id  ", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

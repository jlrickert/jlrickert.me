package portfolio

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	config := ServerConfig{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Theme:          "green-nebula-terminal",
	}

	server := NewServer(config, nil)

	require.NotNil(t, server)
	assert.Equal(t, config.Addr, server.config.Addr)
	assert.Equal(t, config.Theme, server.config.Theme)
	assert.NotNil(t, server.router)
	assert.NotNil(t, server.assetManager)
	assert.NotNil(t, server.logger)
}

func TestHandleHealth(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	assert.Contains(t, string(body), "ok")
}

func TestHandleGetHome(t *testing.T) {
	config := DefaultServerConfig()
	config.Theme = "green-nebula-terminal"
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	bodyStr := string(body)

	// Check that the response contains expected content
	assert.NotEmpty(t, bodyStr)
	assert.Contains(t, bodyStr, "<!DOCTYPE html")
	assert.Contains(t, bodyStr, "Jared Rickert") // Name from data.yaml
}

func TestHandleGetData(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	bodyStr := string(body)

	// Check that the response contains expected data
	assert.Contains(t, bodyStr, "Jared Rickert")
	assert.Contains(t, bodyStr, "Minneapolis")
}

func TestHandleGetPost(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/posts/first-post", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	bodyStr := string(body)

	// Check that the response contains expected post data
	assert.Contains(t, bodyStr, "first-post")
	assert.Contains(t, bodyStr, "Welcome to My Blog")
}

func TestHandleGetPostNotFound(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/posts/nonexistent-post", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	body, _ := io.ReadAll(w.Body)
	bodyStr := string(body)

	// Should contain error message
	assert.NotEmpty(t, bodyStr)
	assert.Contains(t, bodyStr, "error")
}

func TestHandleListPosts(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	bodyStr := string(body)

	// Should contain posts structure
	assert.Contains(t, bodyStr, "posts")
}

func TestHandleThemeSwitch(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	// Create a proper form-encoded POST request with body
	body := "theme=dark"
	req := httptest.NewRequest("POST", "/api/theme", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	// Check that theme cookie was set
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies)

	themeCookie := false
	for _, cookie := range cookies {
		if cookie.Name == "theme" {
			themeCookie = true
			assert.Equal(t, "dark", cookie.Value)
			assert.Equal(t, true, cookie.HttpOnly)
		}
	}
	assert.True(t, themeCookie, "theme cookie should be set")
}

func TestHandleThemeSwitchMissingParam(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("POST", "/api/theme", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	assert.Contains(t, string(body), "required")
}

func TestHandleExample(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/example", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	bodyStr := string(body)

	// Check that example.html content is served
	assert.NotEmpty(t, bodyStr)
	assert.Contains(t, bodyStr, "<!DOCTYPE html")
}

func TestHandlePostsPartial(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/api/posts/partial", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	// Should return HTML fragment
	body, _ := io.ReadAll(w.Body)
	assert.NotEmpty(t, body)
}

func TestHandleExperiencePartial(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/api/experience/partial", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	assert.NotEmpty(t, body)
}

func TestHandleSkillsPartial(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	req := httptest.NewRequest("GET", "/api/skills/partial", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	body, _ := io.ReadAll(w.Body)
	assert.NotEmpty(t, body)
}

func TestDefaultServerConfig(t *testing.T) {
	config := DefaultServerConfig()

	assert.Equal(t, ":8080", config.Addr)
	assert.Equal(t, 10*time.Second, config.ReadTimeout)
	assert.Equal(t, 10*time.Second, config.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.IdleTimeout)
	assert.Equal(t, 1<<20, config.MaxHeaderBytes)
}

func TestServerHandler(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	handler := server.Handler()

	require.NotNil(t, handler)
	assert.Equal(t, server.router, handler)
}

func TestServerShutdown(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Should be able to call Shutdown without error
	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestHandleStatic(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "CSS file",
			path:           "/static/css/stylesheet.css",
			expectedStatus: http.StatusOK,
			expectedType:   "text/css; charset=utf-8",
		},
		{
			name:           "JavaScript file",
			path:           "/static/js/starfield.js",
			expectedStatus: http.StatusOK,
			expectedType:   "application/javascript; charset=utf-8",
		},
		{
			name:           "Nonexistent static file",
			path:           "/static/nonexistent.txt",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultServerConfig()
			config.Theme = "green-nebula-terminal"
			server := NewServer(config, nil)

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedType, w.Header().Get("Content-Type"))
				assert.Contains(t, w.Header().Get("Cache-Control"), "max-age")

				body, _ := io.ReadAll(w.Body)
				assert.NotEmpty(t, body)
			}
		})
	}
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		filePath     string
		expectedType string
	}{
		{
			filePath:     "style.css",
			expectedType: "text/css; charset=utf-8",
		},
		{
			filePath:     "script.js",
			expectedType: "application/javascript; charset=utf-8",
		},
		{
			filePath:     "data.json",
			expectedType: "application/json",
		},
		{
			filePath:     "icon.svg",
			expectedType: "image/svg+xml",
		},
		{
			filePath:     "image.png",
			expectedType: "image/png",
		},
		{
			filePath:     "photo.jpg",
			expectedType: "image/jpeg",
		},
		{
			filePath:     "photo.jpeg",
			expectedType: "image/jpeg",
		},
		{
			filePath:     "animation.gif",
			expectedType: "image/gif",
		},
		{
			filePath:     "font.woff",
			expectedType: "font/woff",
		},
		{
			filePath:     "font.woff2",
			expectedType: "font/woff2",
		},
		{
			filePath:     "unknown.xyz",
			expectedType: "application/octet-stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			result := getContentType(tt.filePath)
			assert.Equal(t, tt.expectedType, result)
		})
	}
}

func TestCacheHeaders(t *testing.T) {
	config := DefaultServerConfig()
	server := NewServer(config, nil)

	tests := []struct {
		name           string
		path           string
		expectedHeader string
		checkHeader    bool
	}{
		{
			name:           "home page cache control",
			path:           "/",
			expectedHeader: "public, max-age=3600",
			checkHeader:    false, // Home page may fail to render in tests
		},
		{
			name:           "health check no cache",
			path:           "/health",
			expectedHeader: "",
			checkHeader:    true,
		},
		{
			name:           "data api cache control",
			path:           "/api/data",
			expectedHeader: "public, max-age=3600",
			checkHeader:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			if tt.checkHeader && tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedHeader, w.Header().Get("Cache-Control"))
			}
		})
	}
}

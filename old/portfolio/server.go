package portfolio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ServerConfig holds configuration for the portfolio server
type ServerConfig struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int

	Theme string
}

// DefaultServerConfig returns sensible defaults for ServerConfig
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
		Theme:          DefaultTheme,
	}
}

// Server represents the portfolio HTTP server
type Server struct {
	config       ServerConfig
	router       *chi.Mux
	assetManager *AssetManager
	logger       *slog.Logger
	httpServer   *http.Server
}

// NewServer creates and returns a new Server instance
func NewServer(config ServerConfig, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(
			os.Stderr,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Compress(5))
	r.Use(middleware.RequestID)

	server := &Server{
		config:       config,
		router:       r,
		assetManager: NewAssetManager(config.Theme),
		logger:       logger,
	}

	server.assetManager.Logger = logger

	server.setupRoutes()
	server.setupHTTPServer()

	return server
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	s.router.Get("/", s.handleGetHome)
	s.router.Get("/posts", s.handleListPosts)
	s.router.Get("/posts/{slug}", s.handleGetPost)
	s.router.Get("/example", s.handleExample)

	// Static files
	s.router.Get("/static/*", s.handleStatic)

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/data", s.handleGetData)
		// htmx partial endpoints
		r.Get("/posts/partial", s.handlePostsPartial)
		r.Get("/experience/partial", s.handleExperiencePartial)
		r.Get("/skills/partial", s.handleSkillsPartial)
		r.Post("/theme", s.handleThemeSwitch)
	})
	//
	// Health check
	s.router.Get("/health", s.handleHealth)
}

// setupHTTPServer configures the underlying HTTP server
func (s *Server) setupHTTPServer() {
	s.httpServer = &http.Server{
		Addr:           s.config.Addr,
		Handler:        s.router,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}
}

// Start starts the server and blocks until it stops
func (s *Server) Start() error {
	s.logger.Info(
		"starting portfolio server",
		"addr", s.config.Addr,
	)
	return http.ListenAndServe(s.config.Addr, s.router)
}

// Shutdown gracefully shuts down the server with timeout
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	return s.httpServer.Shutdown(ctx)
}

// Handler returns the HTTP handler for use with other servers
func (s *Server) Handler() http.Handler {
	return s.router
}

// Handlers

// handleHealth handles GET /health requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// handleGetData handles GET /api/data requests
func (s *Server) handleGetData(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := s.assetManager.GetData(ctx)
	if err != nil {
		s.logger.Error("failed to get data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to retrieve data",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(data)
}

// handleListPosts handles GET /api/posts requests
func (s *Server) handleListPosts(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement post listing with pagination
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"posts": []any{},
	})
}

// handleGetPost handles GET /api/posts/{slug} requests
func (s *Server) handleGetPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "post slug is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	post, err := s.assetManager.GetPage(ctx, slug)
	if err != nil {
		s.logger.Error(
			"failed to get post",
			"slug", slug,
			"error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "post not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(map[string]any{
		"slug":    post.Path,
		"title":   post.Title(),
		"content": string(post.Content),
		"date":    post.Date(),
		"tags":    post.Tags(),
	})
}

func (s *Server) renderTemplate(name string, data any) ([]byte, error) {
	// Load base template content
	baseContent, err := s.assetManager.GetTemplateContent(s.config.Theme, "_base")
	if err != nil {
		s.logger.Error("failed to get base template", "error", err)
		return nil, err
	}

	// Load content template content
	contentContent, err := s.assetManager.GetTemplateContent(s.config.Theme, name)
	if err != nil {
		s.logger.Error("failed to get content template", "name", name, "error", err)
		return nil, err
	}

	// Parse base template with custom functions
	baseTmpl, err := template.New("base").Funcs(TemplateFuncs()).Parse(string(baseContent))
	if err != nil {
		s.logger.Error("failed to parse base template", "error", err)
		return nil, err
	}

	// Parse content template and associate with base
	baseTmpl, err = baseTmpl.New("content").Parse(string(contentContent))
	if err != nil {
		s.logger.Error("failed to parse content template", "name", name, "error", err)
		return nil, err
	}

	var buf bytes.Buffer
	err = baseTmpl.Execute(&buf, data)
	if err != nil {
		s.logger.Error("failed to execute template", "error", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Server) handleGetHome(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := s.assetManager.GetData(ctx)
	if err != nil {
		s.logger.Error("failed to get data for home page", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<h1>Error</h1><p>Failed to load home page</p>")
		return
	}

	html, err := s.renderTemplate("index", data)
	if err != nil {
		s.logger.Error("failed to render index template", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<h1>Error</h1><p>Failed to render page</p>")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(html)
}

// handlePostsPartial handles GET /api/posts/partial for htmx requests
// Supports query parameters: tag (filter by tag), search (search in title/description)
func (s *Server) handlePostsPartial(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := s.assetManager.GetData(ctx)
	if err != nil {
		s.logger.Error("failed to get data for posts partial", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<p>Error loading posts</p>")
		return
	}

	// TODO: Implement post filtering by tag and search query
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=600")
	fmt.Fprint(w, "<!-- Posts partial placeholder -->")
}

// handleExperiencePartial handles GET /api/experience/partial for htmx requests
func (s *Server) handleExperiencePartial(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := s.assetManager.GetData(ctx)
	if err != nil {
		s.logger.Error("failed to get data for experience partial", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<p>Error loading experience</p>")
		return
	}

	// TODO: Implement experience rendering
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	fmt.Fprint(w, "<!-- Experience partial placeholder -->")
}

// handleSkillsPartial handles GET /api/skills/partial for htmx requests
func (s *Server) handleSkillsPartial(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := s.assetManager.GetData(ctx)
	if err != nil {
		s.logger.Error("failed to get data for skills partial", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<p>Error loading skills</p>")
		return
	}

	// TODO: Implement skills rendering
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	fmt.Fprint(w, "<!-- Skills partial placeholder -->")
}

// handleThemeSwitch handles POST /api/theme for htmx theme switching
// Expects form data with "theme" parameter
func (s *Server) handleThemeSwitch(w http.ResponseWriter, r *http.Request) {
	theme := r.FormValue("theme")
	if theme == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "theme parameter required")
		return
	}

	// Set theme cookie for persistence
	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    theme,
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<!-- Theme switched -->")
}

// handleExample handles GET /example and serves the example.html reference design
func (s *Server) handleExample(w http.ResponseWriter, r *http.Request) {
	exampleContent, err := fs.ReadFile(s.assetManager.Assets, "example.html")
	if err != nil {
		s.logger.Error("failed to read example.html", "error", err)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<h1>Not Found</h1><p>example.html not found</p>")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(exampleContent)
}

// handleStatic handles GET /static/* for serving static files (CSS, JS, etc.)
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	// Get the static file path from the URL
	// Format: /static/css/stylesheet.css -> themes/green-nebula-terminal/static/css/stylesheet.css
	staticPath := chi.URLParam(r, "*")
	if staticPath == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filePath := filepath.Join("themes", s.config.Theme, "static", staticPath)

	// Read file from embedded assets
	content, err := fs.ReadFile(s.assetManager.Assets, filePath)
	if err != nil {
		s.logger.Error("failed to read static file", "path", filePath, "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Set content type based on file extension
	contentType := getContentType(filePath)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day
	w.Write(content)
}

// getContentType returns the appropriate Content-Type header for a file
func getContentType(filePath string) string {
	switch {
	case strings.HasSuffix(filePath, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(filePath, ".js"):
		return "application/javascript; charset=utf-8"
	case strings.HasSuffix(filePath, ".json"):
		return "application/json"
	case strings.HasSuffix(filePath, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(filePath, ".png"):
		return "image/png"
	case strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(filePath, ".gif"):
		return "image/gif"
	case strings.HasSuffix(filePath, ".woff"):
		return "font/woff"
	case strings.HasSuffix(filePath, ".woff2"):
		return "font/woff2"
	default:
		return "application/octet-stream"
	}
}

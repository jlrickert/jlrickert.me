package portfolio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
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

	DefaultTheme string
}

// DefaultServerConfig returns sensible defaults for ServerConfig
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
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
		assetManager: NewAssetManager(),
		logger:       logger,
	}

	server.setupRoutes()
	server.setupHTTPServer()

	return server
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	s.router.Get("/", s.handleGetHome)
	s.router.Get("/posts", s.handleListPosts)
	s.router.Get("/posts/{slug}", s.handleGetPost)

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/data", s.handleGetData)
		// htmx partial endpoints
		r.Get("/posts/partial", s.handlePostsPartial)
		r.Get("/experience/partial", s.handleExperiencePartial)
		r.Get("/skills/partial", s.handleSkillsPartial)
		r.Post("/theme", s.handleThemeSwitch)
	})

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
func (s *Server) handleGetData(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	post, err := s.assetManager.GetPost(ctx, slug)
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
		"slug":    post.Slug,
		"title":   post.Title(),
		"content": string(post.Content),
		"date":    post.Date(),
		"tags":    post.Tags(),
	})
}

func (s *Server) renderTemplate(name string, data any) ([]byte, error) {
	themePath := fmt.Sprintf("templates/theme/%s/%s.html", s.config.DefaultTheme, name)
	tmpl, err := template.ParseFiles(themePath)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.Bytes(), err
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

	html, err := s.renderTemplate("home", data)
	if err != nil {
		s.logger.Error("failed to render home template", "error", err)
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

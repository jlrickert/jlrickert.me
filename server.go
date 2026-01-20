package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type Server struct {
	startTime time.Time
	router    *chi.Mux
	logger    *slog.Logger
	upgrader  websocket.Upgrader
}

func NewServer(logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(
			os.Stderr,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Compress(5))
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(15 * time.Second))

	server := &Server{
		router: r,
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.Get("/", s.handleHello)
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/ping", s.handlePing)
	s.router.Get("/ws", s.handleWebSocket)
}

type PingResponse struct {
	ServerTime   time.Time     `json:"serverTime"`
	Timestamp    int64         `json:"timestamp"` // Unix ms for easy JS math
	ServerUptime time.Duration `json:"uptime"`
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	uptime := time.Since(s.startTime) // Adjust for actual start

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	resp := PingResponse{
		ServerTime:   start,
		Timestamp:    start.UnixMilli(),
		ServerUptime: uptime,
	}
	json.NewEncoder(w).Encode(resp)
}

// handleWebSocket handles WebSocket connections for real-time ping data
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("websocket upgrade error", "error", err)
		return
	}
	defer conn.Close()

	s.logger.Debug("websocket connection established")

	for {
		// Read message from client
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.Error("websocket error", "error", err)
			}
			break
		}

		// Generate ping response
		start := time.Now()
		uptime := time.Since(s.startTime)

		resp := PingResponse{
			ServerTime:   start,
			Timestamp:    start.UnixMilli(),
			ServerUptime: uptime,
		}

		// Send response back to client
		if err := conn.WriteJSON(resp); err != nil {
			s.logger.Error("websocket write error", "error", err)
			break
		}
	}
}

// handleHello handles GET / requests
func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello, World!",
		"status":  "ok",
	})
}

// handleHealth handles GET /health requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (s *Server) Start(addr string) error {
	s.startTime = time.Now()
	s.logger.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, s.router)
}

func main() {
	// Parse command-line flags
	addr := flag.String("addr", ":8080", "Server address")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Setup logger
	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: level},
	))

	// Create and start server
	server := NewServer(logger)
	if err := server.Start(*addr); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

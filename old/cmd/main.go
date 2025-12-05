package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jlrickert/jlrickert.me/portfolio"
)

func main() {
	// Parse command-line flags
	addr := flag.String("addr", ":8080", "HTTP server address")
	theme := flag.String("theme", "green-nebula-terminal", "Default theme name")
	flag.Parse()

	// Create logger
	logger := slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	))

	// Create server config
	config := portfolio.ServerConfig{
		Addr:           *addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
		Theme:          *theme,
	}

	// Create and start server
	server := portfolio.NewServer(config, logger)

	// Start server in a goroutine
	go func() {
		logger.Info(
			"starting portfolio server",
			"addr", config.Addr,
			"theme", config.Theme,
		)
		if err := server.Start(); err != nil {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	logger.Info("received shutdown signal")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}

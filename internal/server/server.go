package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Config holds HTTP server tuning knobs for the orchestration layer.
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig returns sensible defaults for a high-throughput API server.
func DefaultConfig() Config {
	return Config{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// Server wraps net/http with lifecycle helpers (start + graceful shutdown).
type Server struct {
	cfg    Config
	server *http.Server
}

// New builds a Server around the given root handler (typically your router).
func New(cfg Config, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

// Start blocks and serves incoming HTTP traffic.
func (s *Server) Start() error {
	fmt.Printf("orchestration server listening on %s\n", s.cfg.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server start: %w", err)
	}
	return nil
}

// Shutdown gracefully drains in-flight requests before stopping.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}

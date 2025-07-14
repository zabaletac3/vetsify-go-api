package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server ahora tiene su propia instancia de logger.
type Server struct {
	server *http.Server
	logger *slog.Logger 
}

// NewServer es el constructor que ahora recibe el logger como dependencia.
func NewServer(port int, logger *slog.Logger) *Server { 
	mux := http.NewServeMux()

	server := &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		logger: logger, // 👈 4. Asignamos el logger
	}

	server.registerRoutes(mux)

	return server
}

// registerRoutes ahora puede usar el logger del servidor si lo necesita.
func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.handleHealthCheck())
}

// Start ahora usa el logger estructurado.
func (s *Server) Start() {
	s.logger.Info("🚀 Servidor escuchando", "address", s.server.Addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("No se pudo iniciar el servidor", "error", err)
			os.Exit(1)
		}
	}()

	s.gracefulShutdown()
}

// gracefulShutdown también usa el logger estructurado.
func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Servidor apagándose...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Fallo en el cierre elegante del servidor", "error", err)
		os.Exit(1)
	}

	s.logger.Info("Servidor apagado exitosamente.")
}

// handleHealthCheck ahora usa el logger del struct Server (`s.logger`).
func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.logger.Error("Error escribiendo respuesta de health check", "error", err)
		}
	}
}
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server contiene las dependencias y la configuración de nuestro servidor HTTP.
type Server struct {
	server *http.Server
}

// NewServer es el constructor de nuestro servidor.
func NewServer(port int) *Server {
	mux := http.NewServeMux()

	server := &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
	
	// Registramos las rutas generales.
	server.registerRoutes(mux)

	return server
}

// registerRoutes define los endpoints de la aplicación.
func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.handleHealthCheck())
}

// Start inicia el servidor y maneja el cierre elegante (graceful shutdown).
func (s *Server) Start() {
	log.Printf("🚀 Servidor escuchando en http://localhost%s", s.server.Addr)
	
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("No se pudo iniciar el servidor: %v", err)
		}
	}()

	s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Servidor apagándose...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Fallo en el cierre elegante del servidor: %v", err)
	}
	
	log.Println("Servidor apagado exitosamente.")
}

// handleHealthCheck es nuestro handler para verificar el estado del servicio.
func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
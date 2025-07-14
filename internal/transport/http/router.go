package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/zabaletac3/go-vet-api/internal/transport/http/clinics"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/users"
	"go.mongodb.org/mongo-driver/mongo"

	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupAllRoutes recibe las dependencias globales y las distribuye.
func SetupAllRoutes(mux *http.ServeMux, db *mongo.Database, logger *slog.Logger) {

	
	// Módulo de Usuarios
	users.RegisterRoutes(mux, db, logger) // 👈 Añadimos la llamada

	// Módulo de Clínicas
	clinics.RegisterRoutes(mux, db, logger)

	// @Summary     Obtener información de salud
	// @Description Endpoint para verificar el estado del servidor
	// @Tags        health
	// @Accept      json
	// @Produce     json
	// @Success     200 {object} map[string]string
	// @Router      /health [get]
	mux.HandleFunc("GET /health", handleHealthCheck(logger))
	mux.HandleFunc("GET /test/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("La ruta de prueba funciona!"))
	})
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler) 
}

// handleHealthCheck ahora es una función privada dentro del paquete http.
func handleHealthCheck(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status":    "ok, go!",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Error escribiendo respuesta de health check", "error", err)
		}
	}
}
package clinics

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/zabaletac3/go-vet-api/internal/services"
	"github.com/zabaletac3/go-vet-api/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

// contextMiddleware agrega dependencias al contexto
func contextMiddleware(db *mongo.Database, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			ctx = context.WithValue(ctx, "logger", logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RegisterRoutes registra todas las rutas de clinics
func RegisterRoutes(mux *http.ServeMux, db *mongo.Database, logger *slog.Logger) {
	// Crear el repository específico del módulo (implementa ClinicStorer)
	clinicRepo := storage.NewClinicRepository(db)
	
	// Crear el service específico del módulo
	clinicService := services.NewClinicService(clinicRepo, logger)
	
	// Crear el handler específico del módulo
	handler := NewHandler(clinicService)
	
	// Middleware para agregar dependencias al contexto
	middleware := contextMiddleware(db, logger)
	
	// Rutas de clinics
	mux.Handle("POST /api/v1/clinics", middleware(handler.CreateClinic(db, logger)))
	mux.Handle("GET /api/v1/clinics/{id}", middleware(http.HandlerFunc(handler.GetClinicByID)))
}
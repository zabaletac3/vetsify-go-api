package users

import (
	"log/slog"
	"net/http"

	"github.com/zabaletac3/go-vet-api/internal/services"
	"github.com/zabaletac3/go-vet-api/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRoutes construye toda la pila para el dominio de usuarios y registra sus rutas.
func RegisterRoutes(mux *http.ServeMux, db *mongo.Database, logger *slog.Logger) {
	// 1. Construimos la cadena de dependencias.
	userRepo := storage.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, logger)
	handler := NewHandler(userSvc)

	// 2. Registramos las rutas de este dominio.
	mux.HandleFunc("POST /api/v1/users/register", handler.register)
	// Aquí añadiríamos más rutas como GET /api/v1/users/{id}, etc.

	logger.Info("Rutas de Usuarios registradas.")
}
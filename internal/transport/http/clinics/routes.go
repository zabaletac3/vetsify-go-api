package clinics

import (
	"log/slog"
	"net/http"

	"github.com/zabaletac3/go-vet-api/internal/services"
	"github.com/zabaletac3/go-vet-api/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(mux *http.ServeMux, db *mongo.Database, logger *slog.Logger) {
	repo := storage.NewClinicRepository(db)
	svc := services.NewClinicService(repo, logger)
	handler := NewHandler(svc)

	mux.HandleFunc("POST /api/v1/clinics", handler.createClinic)
	mux.HandleFunc("GET /api/v1/clinics/", handler.getClinicByID)

	logger.Info("Rutas de Clínicas registradas.")
}
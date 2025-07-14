package clinics

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zabaletac3/go-vet-api/internal/services"
)

type Handler struct {
	service services.ClinicService
}

func NewHandler(svc services.ClinicService) *Handler {
	return &Handler{service: svc}
}

// createClinic es el handler para crear una clínica, con documentación Swagger.
// @Summary      Crea una nueva clínica
// @Description  Registra una nueva clínica (tenant) en el sistema.
// @Tags         Clinics
// @Accept       json
// @Produce      json
// @Param        clinic  body      services.CreateClinicParams  true  "Datos para crear la clínica"
// @Success      201      {object}  models.Clinic
// @Failure      400      {object}  map[string]string "Error: Petición inválida"
// @Failure      500      {object}  map[string]string "Error: Error interno del servidor"
// @Router       /api/v1/clinics [post]
func (h *Handler) createClinic(w http.ResponseWriter, r *http.Request) {
	var params services.CreateClinicParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	clinic, err := h.service.Create(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(clinic)
}

// getClinicByID es el handler para obtener una clínica por ID, con documentación Swagger.
// @Summary      Obtiene una clínica por ID
// @Description  Recupera los detalles de una clínica específica usando su ID.
// @Tags         Clinics
// @Produce      json
// @Param        id   path      string  true  "ID de la Clínica"
// @Success      200  {object}  models.Clinic
// @Failure      404  {object}  map[string]string "Error: Clínica no encontrada"
// @Router       /api/v1/clinics/{id} [get]
func (h *Handler) getClinicByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/clinics/")
	clinic, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "clínica no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clinic)
}
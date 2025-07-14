package clinics

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/zabaletac3/go-vet-api/internal/middleware"
	"github.com/zabaletac3/go-vet-api/internal/services"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/response"
	"go.mongodb.org/mongo-driver/mongo"
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
// @Failure      400      {object}  response.ErrorResponse "Error: Petición inválida"
// @Failure      500      {object}  response.ErrorResponse "Error: Error interno del servidor"
// @Router       /api/v1/clinics [post]
func (h *Handler) createClinic(w http.ResponseWriter, r *http.Request, params services.CreateClinicParams, db *mongo.Database, logger *slog.Logger) {
	clinic, err := h.service.Create(r.Context(), params)
	if err != nil {
		logger.Error("Error creando clínica", "error", err, "params", params)
		response.InternalServerError(w, "Error al crear la clínica", logger, err)
		return
	}

	response.Created(w, "Clínica creada exitosamente", clinic)
}

// CreateClinic es el wrapper público que usa el middleware de validación
func (h *Handler) CreateClinic(db *mongo.Database, logger *slog.Logger) http.HandlerFunc {
	return middleware.ValidateRequestWithDeps(h.createClinic, db, logger)
}

// getClinicByID es el handler para obtener una clínica por ID, con documentación Swagger.
// @Summary      Obtiene una clínica por ID
// @Description  Recupera los detalles de una clínica específica usando su ID.
// @Tags         Clinics
// @Produce      json
// @Param        id   path      string  true  "ID de la Clínica"
// @Success      200  {object}  models.Clinic
// @Failure      404  {object}  response.ErrorResponse "Error: Clínica no encontrada"
// @Failure      500  {object}  response.ErrorResponse "Error: Error interno del servidor"
// @Router       /api/v1/clinics/{id} [get]
func (h *Handler) GetClinicByID(w http.ResponseWriter, r *http.Request) {
	// Obtener logger del contexto si está disponible
	logger, ok := r.Context().Value("logger").(*slog.Logger)
	if !ok {
		logger = slog.Default()
	}

	// Obtener ID del path parameter
	id := r.PathValue("id")
	if id == "" {
		response.BadRequest(w, "ID de clínica requerido")
		return
	}

	clinic, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no documents") {
			response.NotFound(w, "Clínica no encontrada")
			return
		}
		
		logger.Error("Error obteniendo clínica", "error", err, "id", id)
		response.InternalServerError(w, "Error al obtener la clínica", logger, err)
		return
	}

	response.Success(w, "Clínica encontrada", clinic)
}

// // GetAllClinics obtiene todas las clínicas con paginación
// // @Summary      Obtiene todas las clínicas
// // @Description  Recupera una lista paginada de todas las clínicas del sistema.
// // @Tags         Clinics
// // @Produce      json
// // @Param        page     query    int  false  "Número de página (por defecto: 1)"
// // @Param        limit    query    int  false  "Elementos por página (por defecto: 10)"
// // @Success      200      {object}  response.PaginatedResponse
// // @Failure      500      {object}  response.ErrorResponse "Error: Error interno del servidor"
// // @Router       /api/v1/clinics [get]
// func (h *Handler) GetAllClinics(w http.ResponseWriter, r *http.Request) {
// 	logger, ok := r.Context().Value("logger").(*slog.Logger)
// 	if !ok {
// 		logger = slog.Default()
// 	}

// 	// Obtener parámetros de paginación de la query string
// 	page := 1
// 	limit := 10
	
// 	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
// 		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
// 			page = p
// 		}
// 	}
	
// 	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
// 		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
// 			limit = l
// 		}
// 	}

// 	clinics, total, err := h.service.GetAll(r.Context(), page, limit)
// 	if err != nil {
// 		logger.Error("Error obteniendo clínicas", "error", err, "page", page, "limit", limit)
// 		response.InternalServerError(w, "Error al obtener las clínicas", logger, err)
// 		return
// 	}

// 	totalPages := int(math.Ceil(float64(total) / float64(limit)))
// 	response.Paginated(w, clinics, page, limit, totalPages, total)
// }

// // UpdateClinic actualiza una clínica existente
// // @Summary      Actualiza una clínica
// // @Description  Actualiza los datos de una clínica existente.
// // @Tags         Clinics
// // @Accept       json
// // @Produce      json
// // @Param        id      path      string                       true  "ID de la Clínica"
// // @Param        clinic  body      services.UpdateClinicParams  true  "Datos para actualizar la clínica"
// // @Success      200      {object}  models.Clinic
// // @Failure      400      {object}  response.ValidationErrorResponse "Error: Datos de entrada inválidos"
// // @Failure      404      {object}  response.ErrorResponse "Error: Clínica no encontrada"
// // @Failure      500      {object}  response.ErrorResponse "Error: Error interno del servidor"
// // @Router       /api/v1/clinics/{id} [put]
// func (h *Handler) updateClinic(w http.ResponseWriter, r *http.Request, params services.UpdateClinicParams, db *mongo.Database, logger *slog.Logger) {
// 	id := strings.TrimPrefix(r.URL.Path, "/api/v1/clinics/")
// 	if id == "" {
// 		response.BadRequest(w, "ID de clínica requerido")
// 		return
// 	}

// 	clinic, err := h.service.Update(r.Context(), id, params)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no documents") {
// 			response.NotFound(w, "Clínica no encontrada")
// 			return
// 		}
		
// 		logger.Error("Error actualizando clínica", "error", err, "id", id, "params", params)
// 		response.InternalServerError(w, "Error al actualizar la clínica", logger, err)
// 		return
// 	}

// 	response.Success(w, "Clínica actualizada exitosamente", clinic)
// }

// // UpdateClinic es el wrapper público que usa el middleware de validación
// func (h *Handler) UpdateClinic(db *mongo.Database, logger *slog.Logger) http.HandlerFunc {
// 	return middleware.ValidateRequestWithDeps(h.updateClinic, db, logger)
// }

// // DeleteClinic elimina una clínica
// // @Summary      Elimina una clínica
// // @Description  Elimina una clínica del sistema.
// // @Tags         Clinics
// // @Produce      json
// // @Param        id   path      string  true  "ID de la Clínica"
// // @Success      204  "Sin contenido"
// // @Failure      404  {object}  response.ErrorResponse "Error: Clínica no encontrada"
// // @Failure      500  {object}  response.ErrorResponse "Error: Error interno del servidor"
// // @Router       /api/v1/clinics/{id} [delete]
// func (h *Handler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
// 	logger, ok := r.Context().Value("logger").(*slog.Logger)
// 	if !ok {
// 		logger = slog.Default()
// 	}

// 	id := r.PathValue("id")
// 	if id == "" {
// 		response.BadRequest(w, "ID de clínica requerido")
// 		return
// 	}

// 	err := h.service.Delete(r.Context(), id)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no documents") {
// 			response.NotFound(w, "Clínica no encontrada")
// 			return
// 		}
		
// 		logger.Error("Error eliminando clínica", "error", err, "id", id)
// 		response.InternalServerError(w, "Error al eliminar la clínica", logger, err)
// 		return
// 	}

// 	response.NoContent(w)
// }
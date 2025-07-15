// internal/transport/http/clinics/handler.go
package clinics

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/zabaletac3/go-vet-api/internal/middleware"
	"github.com/zabaletac3/go-vet-api/internal/models"
	"github.com/zabaletac3/go-vet-api/internal/services"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/response"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
    service services.ClinicService
    logger  *slog.Logger
}

func NewHandler(svc services.ClinicService, logger *slog.Logger) *Handler {
    return &Handler{
        service: svc,
        logger:  logger.With("handler", "clinic"),
    }
}

// createClinic maneja la creación de clínicas
// @Summary      Create a new clinic
// @Description  Register a new clinic (tenant) in the system with color palette
// @Tags         Clinics
// @Accept       json
// @Produce      json
// @Param        clinic  body      CreateClinicRequest  true  "Clinic data"
// @Success      201      {object}  ClinicResponse
// @Failure      400      {object}  response.ValidationErrorResponse "Invalid data"
// @Failure      409      {object}  response.ErrorResponse "Name already exists"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Router       /api/v1/clinics [post]
func (h *Handler) createClinic(w http.ResponseWriter, r *http.Request, req CreateClinicRequest, db *mongo.Database, logger *slog.Logger) {
    // Convertir a parámetros de servicio
    params := services.CreateClinicParams{
        Name:        req.Name,
        DisplayName: req.DisplayName,
        Address:     req.Address,
        Phone:       req.Phone,
        Email:       req.Email,
        Website:     req.Website,
        Description: req.Description,
    }

    // Establecer paleta de colores
    if req.Palette != nil {
        params.Palette = req.Palette.ToModel()
    } else {
        params.Palette = models.GetDefaultPalette()
    }

    clinic, err := h.service.Create(r.Context(), params)
    if err != nil {
        switch err {
        case services.ErrClinicNameRequired:
            response.Error(w, http.StatusBadRequest, "Bad Request", "Clinic name is required")
            return
        case services.ErrClinicNameExists:
            response.Error(w, http.StatusConflict, "Conflict", "A clinic with that name already exists")
            return
        case services.ErrDisplayNameExists:
            response.Error(w, http.StatusConflict, "Conflict", "A clinic with that display name already exists")
            return
        default:
            logger.Error("Error creating clinic", "error", err, "params", params)
            response.Error(w, http.StatusInternalServerError, "Internal Server Error", "Failed to create clinic")
            return
        }
    }

    // Convertir a DTO de respuesta
    clinicResponse := FromModel(clinic)
    response.JSON(w, http.StatusCreated, response.SuccessResponse{
        Success: true,
        Message: "Clinic created successfully",
        Data:    clinicResponse,
    })
}

// CreateClinic es el wrapper público que usa el middleware de validación
func (h *Handler) CreateClinic(db *mongo.Database, logger *slog.Logger) http.HandlerFunc {
    return middleware.ValidateRequestWithDeps(h.createClinic, db, logger)
}

// GetClinicByID obtiene una clínica por ID
// @Summary      Get clinic by ID
// @Description  Retrieve a specific clinic using its ID
// @Tags         Clinics
// @Produce      json
// @Param        id   path      string  true  "Clinic ID"
// @Success      200  {object}  ClinicResponse
// @Failure      400  {object}  response.ErrorResponse "Invalid ID"
// @Failure      404  {object}  response.ErrorResponse "Clinic not found"
// @Failure      500  {object}  response.ErrorResponse "Internal server error"
// @Router       /api/v1/clinics/{id} [get]
func (h *Handler) GetClinicByID(w http.ResponseWriter, r *http.Request) {
    logger, ok := r.Context().Value("logger").(*slog.Logger)
    if !ok {
        logger = slog.Default()
    }

    id := r.PathValue("id")
    if id == "" {
        response.Error(w, http.StatusBadRequest, "Bad Request", "Clinic ID is required")
        return
    }

    clinic, err := h.service.GetByID(r.Context(), id)
    if err != nil {
        switch err {
        case services.ErrClinicNotFound:
            response.Error(w, http.StatusNotFound, "Not Found", "Clinic not found")
            return
        case services.ErrInvalidClinicID:
            response.Error(w, http.StatusBadRequest, "Bad Request", "Invalid clinic ID")
            return
        default:
            logger.Error("Error getting clinic", "error", err, "id", id)
            response.Error(w, http.StatusInternalServerError, "Internal Server Error", "Failed to get clinic")
            return
        }
    }

    // Convertir a DTO de respuesta
    clinicResponse := FromModel(clinic)
    response.JSON(w, http.StatusOK, response.SuccessResponse{
        Success: true,
        Message: "Clinic found",
        Data:    clinicResponse,
    })
}

// updateClinic maneja la actualización parcial de clínicas (PATCH)
// @Summary      Update clinic (partial)
// @Description  Partially update an existing clinic's data (only provided fields)
// @Tags         Clinics
// @Accept       json
// @Produce      json
// @Param        id      path      string                true  "Clinic ID"
// @Param        clinic  body      UpdateClinicRequest   true  "Fields to update (partial)"
// @Success      200      {object}  ClinicResponse
// @Failure      400      {object}  response.ValidationErrorResponse "Invalid data"
// @Failure      404      {object}  response.ErrorResponse "Clinic not found"
// @Failure      409      {object}  response.ErrorResponse "Name already exists"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Router       /api/v1/clinics/{id} [patch]
func (h *Handler) updateClinic(w http.ResponseWriter, r *http.Request, req UpdateClinicRequest, db *mongo.Database, logger *slog.Logger) {
    id := r.PathValue("id")
    if id == "" {
        response.Error(w, http.StatusBadRequest, "Bad Request", "Clinic ID is required")
        return
    }

    // Convertir a parámetros de servicio
    params := services.UpdateClinicParams{
        Name:        req.Name,
        DisplayName: req.DisplayName,
        Address:     req.Address,
        Phone:       req.Phone,
        Email:       req.Email,
        Website:     req.Website,
        Description: req.Description,
        IsActive:    req.IsActive,
    }

    // Convertir paleta si se proporciona
    if req.Palette != nil {
        palette := req.Palette.ToModel()
        params.Palette = &palette
    }

    clinic, err := h.service.Update(r.Context(), id, params)
    if err != nil {
        switch err {
        case services.ErrClinicNotFound:
            response.Error(w, http.StatusNotFound, "Not Found", "Clinic not found")
            return
        case services.ErrClinicNameExists:
            response.Error(w, http.StatusConflict, "Conflict", "A clinic with that name already exists")
            return
        case services.ErrDisplayNameExists:
            response.Error(w, http.StatusConflict, "Conflict", "A clinic with that display name already exists")
            return
        case services.ErrInvalidClinicID:
            response.Error(w, http.StatusBadRequest, "Bad Request", "Invalid clinic ID")
            return
        default:
            logger.Error("Error updating clinic", "error", err, "id", id, "params", params)
            response.Error(w, http.StatusInternalServerError, "Internal Server Error", "Failed to update clinic")
            return
        }
    }

    // Convertir a DTO de respuesta
    clinicResponse := FromModel(clinic)
    response.JSON(w, http.StatusOK, response.SuccessResponse{
        Success: true,
        Message: "Clinic updated successfully",
        Data:    clinicResponse,
    })
}

// UpdateClinic es el wrapper público que usa el middleware de validación
func (h *Handler) UpdateClinic(db *mongo.Database, logger *slog.Logger) http.HandlerFunc {
    return middleware.ValidateRequestWithDeps(h.updateClinic, db, logger)
}

// DeleteClinic elimina una clínica (soft delete)
// @Summary      Delete clinic
// @Description  Delete a clinic from the system (soft delete - marks as inactive)
// @Tags         Clinics
// @Produce      json
// @Param        id   path      string  true  "Clinic ID"
// @Success      200  {object}  response.SuccessResponse "Clinic deleted successfully"
// @Failure      400  {object}  response.ErrorResponse "Invalid ID"
// @Failure      404  {object}  response.ErrorResponse "Clinic not found"
// @Failure      500  {object}  response.ErrorResponse "Internal server error"
// @Router       /api/v1/clinics/{id} [delete]
func (h *Handler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
    logger, ok := r.Context().Value("logger").(*slog.Logger)
    if !ok {
        logger = slog.Default()
    }

    id := r.PathValue("id")
    if id == "" {
        response.Error(w, http.StatusBadRequest, "Bad Request", "Clinic ID is required")
        return
    }

    err := h.service.Delete(r.Context(), id)
    if err != nil {
        switch err {
        case services.ErrClinicNotFound:
            response.Error(w, http.StatusNotFound, "Not Found", "Clinic not found")
            return
        case services.ErrInvalidClinicID:
            response.Error(w, http.StatusBadRequest, "Bad Request", "Invalid clinic ID")
            return
        default:
            logger.Error("Error deleting clinic", "error", err, "id", id)
            response.Error(w, http.StatusInternalServerError, "Internal Server Error", "Failed to delete clinic")
            return
        }
    }

    response.JSON(w, http.StatusOK, response.SuccessResponse{
        Success: true,
        Message: "Clinic deleted successfully",
        Data:    nil,
    })
}

// GetAllClinics obtiene todas las clínicas con paginación
// @Summary      Get all clinics
// @Description  Retrieve a paginated list of all clinics
// @Tags         Clinics
// @Produce      json
// @Param        page       query    int     false  "Page number (default: 1)"
// @Param        limit      query    int     false  "Items per page (default: 10, max: 100)"
// @Param        search     query    string  false  "Search term"
// @Param        is_active  query    bool    false  "Filter by active status"
// @Param        sort_by    query    string  false  "Sort field (name, display_name, created_at, updated_at)"
// @Param        sort_desc  query    bool    false  "Sort descending"
// @Success      200        {object}  ListClinicsResponse
// @Failure      400        {object}  response.ErrorResponse "Invalid parameters"
// @Failure      500        {object}  response.ErrorResponse "Internal server error"
// @Router       /api/v1/clinics [get]
func (h *Handler) GetAllClinics(w http.ResponseWriter, r *http.Request) {
    logger, ok := r.Context().Value("logger").(*slog.Logger)
    if !ok {
        logger = slog.Default()
    }

    // Parsear parámetros de query
    req := ListClinicsRequest{}
    
    // Page
    if pageStr := r.URL.Query().Get("page"); pageStr != "" {
        if page, err := strconv.Atoi(pageStr); err == nil {
            req.Page = page
        }
    }
    
    // Limit
    if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
        if limit, err := strconv.Atoi(limitStr); err == nil {
            req.Limit = limit
        }
    }
    
    // Search
    req.Search = r.URL.Query().Get("search")
    
    // IsActive
    if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
        if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
            req.IsActive = &isActive
        }
    }
    
    // SortBy
    req.SortBy = r.URL.Query().Get("sort_by")
    
    // SortDesc
    if sortDescStr := r.URL.Query().Get("sort_desc"); sortDescStr != "" {
        if sortDesc, err := strconv.ParseBool(sortDescStr); err == nil {
            req.SortDesc = sortDesc
        }
    }

    // Establecer valores por defecto
    req.SetDefaults()

    // Convertir a parámetros de servicio
    params := services.ListClinicsParams{
        Page:     req.Page,
        Limit:    req.Limit,
        Search:   req.Search,
        IsActive: req.IsActive,
        SortBy:   req.SortBy,
        SortDesc: req.SortDesc,
    }

    // Ejecutar servicio
    clinics, pagination, err := h.service.List(r.Context(), params)
    if err != nil {
        logger.Error("Error listing clinics", "error", err, "params", params)
        response.Error(w, http.StatusInternalServerError, "Internal Server Error", "Failed to list clinics")
        return
    }

    // Convertir a DTOs de respuesta
    clinicResponses := FromModels(clinics)

    // Crear respuesta paginada específica (para Swagger)
    listResponse := ListClinicsResponse{
        Data:       clinicResponses,
        Pagination: pagination,
    }

    response.JSON(w, http.StatusOK, listResponse)
}
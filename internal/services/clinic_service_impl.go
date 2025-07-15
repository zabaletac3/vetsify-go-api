// internal/services/clinic_service_impl.go
package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/zabaletac3/go-vet-api/internal/models"
	"github.com/zabaletac3/go-vet-api/internal/storage"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/dto"
)

// Errores específicos del dominio de negocio
var (
    ErrClinicNameRequired = errors.New("clinic name is required")
    ErrClinicNameExists   = errors.New("clinic with that name already exists")
    ErrClinicNotFound     = errors.New("clinic not found")
    ErrInvalidClinicID    = errors.New("invalid clinic ID")
    ErrDisplayNameExists  = errors.New("clinic with that display name already exists")
)

type clinicService struct {
    store  storage.ClinicStorer
    logger *slog.Logger
}

func NewClinicService(store storage.ClinicStorer, logger *slog.Logger) ClinicService {
    return &clinicService{
        store:  store,
        logger: logger.With("service", "clinic"),
    }
}

// Create - Lógica de creación robusta
func (s *clinicService) Create(ctx context.Context, params CreateClinicParams) (*models.Clinic, error) {
    // Validación de parámetros
    if strings.TrimSpace(params.Name) == "" {
        return nil, ErrClinicNameRequired
    }

    // Verificar unicidad del nombre
    existing, err := s.store.GetByName(ctx, params.Name)
    if err != nil {
        s.logger.Error("Error checking unique name", "error", err, "name", params.Name)
        return nil, fmt.Errorf("error checking unique name: %w", err)
    }
    if existing != nil {
        return nil, ErrClinicNameExists
    }

    // Verificar unicidad del display name
    existingDisplay, err := s.store.GetByDisplayName(ctx, params.DisplayName)
    if err != nil {
        s.logger.Error("Error checking unique display name", "error", err, "displayName", params.DisplayName)
        return nil, fmt.Errorf("error checking unique display name: %w", err)
    }
    if existingDisplay != nil {
        return nil, ErrDisplayNameExists
    }

    // Crear modelo
    clinic := &models.Clinic{
        Name:        strings.TrimSpace(params.Name),
        DisplayName: strings.TrimSpace(params.DisplayName),
        Address:     strings.TrimSpace(params.Address),
        Phone:       strings.TrimSpace(params.Phone),
        Email:       strings.TrimSpace(params.Email),
        Website:     strings.TrimSpace(params.Website),
        Description: strings.TrimSpace(params.Description),
        Palette:     params.Palette,
    }

    // Establecer paleta por defecto si está vacía
    if clinic.Palette.Primary == "" {
        clinic.Palette = models.GetDefaultPalette()
    }

    // Persistir
    if err := s.store.Create(ctx, clinic); err != nil {
        s.logger.Error("Error creating clinic", "error", err, "name", clinic.Name)
        return nil, fmt.Errorf("failed to create clinic: %w", err)
    }

    s.logger.Info("Clinic created successfully", 
        "clinic_id", clinic.ID.Hex(), 
        "name", clinic.Name,
        "display_name", clinic.DisplayName)

    return clinic, nil
}

// GetByID - Obtención robusta con logging
func (s *clinicService) GetByID(ctx context.Context, id string) (*models.Clinic, error) {
    if strings.TrimSpace(id) == "" {
        return nil, ErrInvalidClinicID
    }

    clinic, err := s.store.GetByID(ctx, id)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            return nil, ErrClinicNotFound
        }
        if strings.Contains(err.Error(), "invalid") {
            return nil, ErrInvalidClinicID
        }
        
        s.logger.Error("Error getting clinic", "error", err, "id", id)
        return nil, fmt.Errorf("failed to get clinic: %w", err)
    }

    return clinic, nil
}

// Update - Actualización robusta con PATCH real
func (s *clinicService) Update(ctx context.Context, id string, params UpdateClinicParams) (*models.Clinic, error) {
    // Obtener clínica existente
    existing, err := s.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Crear map de campos a actualizar (solo los que vienen en el request)
    updateFields := make(map[string]interface{})

    // Verificar nombre único si se está cambiando
    if params.Name != nil {
        nameExists, err := s.store.GetByName(ctx, *params.Name)
        if err != nil {
            s.logger.Error("Error checking unique name", "error", err, "name", *params.Name)
            return nil, fmt.Errorf("error checking unique name: %w", err)
        }
        if nameExists != nil && nameExists.ID != existing.ID {
            return nil, ErrClinicNameExists
        }
        updateFields["name"] = strings.TrimSpace(*params.Name)
    }

    // Verificar display name único si se está cambiando
    if params.DisplayName != nil {
        displayExists, err := s.store.GetByDisplayName(ctx, *params.DisplayName)
        if err != nil {
            s.logger.Error("Error checking unique display name", "error", err, "displayName", *params.DisplayName)
            return nil, fmt.Errorf("error checking unique display name: %w", err)
        }
        if displayExists != nil && displayExists.ID != existing.ID {
            return nil, ErrDisplayNameExists
        }
        updateFields["displayName"] = strings.TrimSpace(*params.DisplayName)
    }

    // Agregar otros campos solo si están presentes
    if params.Address != nil {
        updateFields["address"] = strings.TrimSpace(*params.Address)
    }
    if params.Phone != nil {
        updateFields["phone"] = strings.TrimSpace(*params.Phone)
    }
    if params.Email != nil {
        updateFields["email"] = strings.TrimSpace(*params.Email)
    }
    if params.Website != nil {
        updateFields["website"] = strings.TrimSpace(*params.Website)
    }
    if params.Description != nil {
        updateFields["description"] = strings.TrimSpace(*params.Description)
    }
    if params.IsActive != nil {
        updateFields["isActive"] = *params.IsActive
    }
    if params.Palette != nil {
        updateFields["palette"] = *params.Palette
    }

    // Si no hay campos para actualizar
    if len(updateFields) == 0 {
        return existing, nil // Retornar sin cambios
    }

    // Persistir cambios (SOLO los campos enviados)
    if err := s.store.Update(ctx, id, updateFields); err != nil {
        s.logger.Error("Error updating clinic", "error", err, "id", id, "fields", updateFields)
        return nil, fmt.Errorf("failed to update clinic: %w", err)
    }

    // Obtener la clínica actualizada para devolver
    updated, err := s.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("error retrieving updated clinic: %w", err)
    }

    s.logger.Info("Clinic updated successfully", 
        "clinic_id", id, 
        "updated_fields", updateFields)

    return updated, nil
}

// Delete - Eliminación segura (soft delete)
func (s *clinicService) Delete(ctx context.Context, id string) error {
    // Verificar que existe
    _, err := s.GetByID(ctx, id)
    if err != nil {
        return err
    }

    // TODO: Verificar dependencias (usuarios, mascotas, etc.)
    // Por ahora solo soft delete

    if err := s.store.Delete(ctx, id); err != nil {
        s.logger.Error("Error deleting clinic", "error", err, "id", id)
        return fmt.Errorf("failed to delete clinic: %w", err)
    }

    s.logger.Info("Clinic deleted successfully", "clinic_id", id)
    return nil
}

// List - Listado robusto con validación de parámetros
func (s *clinicService) List(ctx context.Context, params ListClinicsParams) ([]*models.Clinic, dto.PaginationResponse, error) {
    // Validar y normalizar parámetros
    normalizedParams := s.normalizeListParams(params)

    // Convertir a filtros de storage
    filters := storage.ListFilters{
        Page:     normalizedParams.Page,
        Limit:    normalizedParams.Limit,
        Search:   normalizedParams.Search,
        IsActive: normalizedParams.IsActive,
        SortBy:   normalizedParams.SortBy,
        SortDesc: normalizedParams.SortDesc,
    }

    // Ejecutar consulta
    clinics, total, err := s.store.List(ctx, filters)
    if err != nil {
        s.logger.Error("Error listing clinics", "error", err, "filters", filters)
        return nil, dto.PaginationResponse{}, fmt.Errorf("failed to list clinics: %w", err)
    }

    // Calcular metadatos de paginación (usando DTO reutilizable)
    pagination := storage.CalculatePagination(normalizedParams.Page, normalizedParams.Limit, total)

    return clinics, pagination, nil
}

// GetByName - Búsqueda por nombre
func (s *clinicService) GetByName(ctx context.Context, name string) (*models.Clinic, error) {
    return s.store.GetByName(ctx, name)
}

// GetByDisplayName - Búsqueda por display name
func (s *clinicService) GetByDisplayName(ctx context.Context, displayName string) (*models.Clinic, error) {
    return s.store.GetByDisplayName(ctx, displayName)
}

// Exists - Verificar existencia
func (s *clinicService) Exists(ctx context.Context, id string) (bool, error) {
    return s.store.Exists(ctx, id)
}

// Métodos helper privados

func (s *clinicService) mergeUpdateParams(existing *models.Clinic, params UpdateClinicParams) *models.Clinic {
    updated := *existing // Copia

    if params.Name != nil {
        updated.Name = strings.TrimSpace(*params.Name)
    }
    if params.DisplayName != nil {
        updated.DisplayName = strings.TrimSpace(*params.DisplayName)
    }
    if params.Address != nil {
        updated.Address = strings.TrimSpace(*params.Address)
    }
    if params.Phone != nil {
        updated.Phone = strings.TrimSpace(*params.Phone)
    }
    if params.Email != nil {
        updated.Email = strings.TrimSpace(*params.Email)
    }
    if params.Website != nil {
        updated.Website = strings.TrimSpace(*params.Website)
    }
    if params.Description != nil {
        updated.Description = strings.TrimSpace(*params.Description)
    }
    if params.Palette != nil {
        updated.Palette = *params.Palette
    }
    if params.IsActive != nil {
        updated.IsActive = *params.IsActive
    }

    return &updated
}

func (s *clinicService) normalizeListParams(params ListClinicsParams) ListClinicsParams {
    normalized := params

    if normalized.Page < 1 {
        normalized.Page = 1
    }
    if normalized.Limit < 1 || normalized.Limit > 100 {
        normalized.Limit = 10
    }

    // Validar campo de ordenamiento
    validSortFields := map[string]bool{
        "name":         true,
        "display_name": true,
        "created_at":   true,
        "updated_at":   true,
    }

    if normalized.SortBy == "" || !validSortFields[normalized.SortBy] {
        normalized.SortBy = "created_at"
    }

    return normalized
}
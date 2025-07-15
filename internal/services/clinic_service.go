// internal/services/clinic_service.go
package services

import (
	"context"

	"github.com/zabaletac3/go-vet-api/internal/models"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/dto"
)

// CreateClinicParams - Parámetros para crear clínica (viene de los DTOs)
type CreateClinicParams struct {
    Name        string
    DisplayName string
    Address     string
    Phone       string
    Email       string
    Website     string
    Description string
    Palette     models.ColorPalette
}

// UpdateClinicParams - Parámetros para actualizar clínica
type UpdateClinicParams struct {
    Name        *string
    DisplayName *string
    Address     *string
    Phone       *string
    Email       *string
    Website     *string
    Description *string
    Palette     *models.ColorPalette
    IsActive    *bool
}

// ListClinicsParams - Parámetros para listar clínicas
type ListClinicsParams struct {
    Page     int
    Limit    int
    Search   string
    IsActive *bool
    SortBy   string
    SortDesc bool
}

// ClinicService - Interface principal de servicio
type ClinicService interface {
    // Operaciones CRUD
    Create(ctx context.Context, params CreateClinicParams) (*models.Clinic, error)
    GetByID(ctx context.Context, id string) (*models.Clinic, error)
    Update(ctx context.Context, id string, params UpdateClinicParams) (*models.Clinic, error)
    Delete(ctx context.Context, id string) error
    
    // Operaciones de consulta (USA DTO REUTILIZABLE)
    List(ctx context.Context, params ListClinicsParams) ([]*models.Clinic, dto.PaginationResponse, error)
    GetByName(ctx context.Context, name string) (*models.Clinic, error)
    GetByDisplayName(ctx context.Context, displayName string) (*models.Clinic, error)
    Exists(ctx context.Context, id string) (bool, error)
}
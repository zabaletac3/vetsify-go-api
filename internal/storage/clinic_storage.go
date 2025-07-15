// internal/storage/clinic_storage.go
package storage

import (
	"context"

	"github.com/zabaletac3/go-vet-api/internal/models"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/dto"
)

// ClinicStorer - Interface para operaciones de clínica
type ClinicStorer interface {
    // Operaciones CRUD básicas
    Create(ctx context.Context, clinic *models.Clinic) error
    GetByID(ctx context.Context, id string) (*models.Clinic, error)
    Update(ctx context.Context, id string, updateFields map[string]interface{}) error
    Delete(ctx context.Context, id string) error // Soft delete simple
    
    // Operaciones de consulta
    List(ctx context.Context, filters ListFilters) ([]*models.Clinic, int64, error)
    GetByName(ctx context.Context, name string) (*models.Clinic, error)
    GetByDisplayName(ctx context.Context, displayName string) (*models.Clinic, error)
    Exists(ctx context.Context, id string) (bool, error)
}

// ListFilters - Filtros para listados
type ListFilters struct {
    Page     int
    Limit    int
    Search   string
    IsActive *bool
    SortBy   string
    SortDesc bool
}

// CalculatePagination - Calcula metadatos de paginación (usa DTO reutilizable)
func CalculatePagination(page, limit int, total int64) dto.PaginationResponse {
    totalPages := int((total + int64(limit) - 1) / int64(limit))
    
    return dto.PaginationResponse{
        CurrentPage: page,
        PerPage:     limit,
        TotalPages:  totalPages,
        Total:       total,
        HasNext:     page < totalPages,
        HasPrev:     page > 1,
    }
}
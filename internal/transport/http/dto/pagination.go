// internal/transport/http/dto/pagination.go
package dto

import "errors"

// PaginationRequest - DTO reutilizable para requests paginados
type PaginationRequest struct {
    Page     int    `json:"page" validate:"min=1" form:"page"`
    Limit    int    `json:"limit" validate:"min=1,max=100" form:"limit"`
    Search   string `json:"search" form:"search"`
    SortBy   string `json:"sort_by" form:"sort_by"`
    SortDesc bool   `json:"sort_desc" form:"sort_desc"`
}

// PaginationResponse - DTO reutilizable para respuestas paginadas
type PaginationResponse struct {
    CurrentPage int   `json:"currentPage"`
    PerPage     int   `json:"perPage"`
    TotalPages  int   `json:"totalPages"`
    Total       int64 `json:"total"`
    HasNext     bool  `json:"hasNext"`
    HasPrev     bool  `json:"hasPrev"`
}

// PaginatedResponse - Wrapper genérico para respuestas paginadas
type PaginatedResponse[T any] struct {
    Data       []T                `json:"data"`
    Pagination PaginationResponse `json:"pagination"`
}

// SetDefaults establece valores por defecto para paginación
func (p *PaginationRequest) SetDefaults() {
    if p.Page < 1 {
        p.Page = 1
    }
    if p.Limit < 1 || p.Limit > 100 {
        p.Limit = 10
    }
}

// Validate valida los parámetros de paginación
func (p *PaginationRequest) Validate() error {
    if p.Page < 1 {
        return errors.New("page must be greater than 0")
    }
    if p.Limit < 1 {
        return errors.New("limit must be greater than 0")
    }
    if p.Limit > 100 {
        return errors.New("limit cannot be greater than 100")
    }
    return nil
}
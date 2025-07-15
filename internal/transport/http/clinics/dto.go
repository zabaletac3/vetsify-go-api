// internal/transport/http/clinics/dto.go
package clinics

import (
	"strings"
	"time"

	"github.com/zabaletac3/go-vet-api/internal/models"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/dto"
)

// CreateClinicRequest - DTO para crear clínica
type CreateClinicRequest struct {
    Name        string             `json:"name" validate:"required,min=2,max=100,clinic_name"`
    DisplayName string             `json:"displayName" validate:"required,min=2,max=150,display_name"`
    Address     string             `json:"address" validate:"omitempty,min=5,max=200"`
    Phone       string             `json:"phone" validate:"omitempty,min=7,max=20"`
    Email       string             `json:"email" validate:"omitempty,email"`
    Website     string             `json:"website" validate:"omitempty,url"`
    Description string             `json:"description" validate:"omitempty,max=500"`
    Palette     *ColorPaletteDTO   `json:"palette,omitempty"`
}

// UpdateClinicRequest - DTO para actualizar clínica
type UpdateClinicRequest struct {
    Name        *string            `json:"name" validate:"omitempty,min=2,max=100,clinic_name"`
    DisplayName *string            `json:"displayName" validate:"omitempty,min=2,max=150,display_name"`
    Address     *string            `json:"address" validate:"omitempty,min=5,max=200"`
    Phone       *string            `json:"phone" validate:"omitempty,min=7,max=20"`
    Email       *string            `json:"email" validate:"omitempty,email"`
    Website     *string            `json:"website" validate:"omitempty,url"`
    Description *string            `json:"description" validate:"omitempty,max=500"`
    Palette     *ColorPaletteDTO   `json:"palette,omitempty"`
    IsActive    *bool              `json:"isActive"`
}

// ColorPaletteDTO - DTO para paleta de colores
type ColorPaletteDTO struct {
    Primary    *string `json:"primary" validate:"omitempty,hex_color"`
    Secondary  *string `json:"secondary" validate:"omitempty,hex_color"`
    Tertiary   *string `json:"tertiary" validate:"omitempty,hex_color"`
    Quaternary *string `json:"quaternary" validate:"omitempty,hex_color"`
    Background *string `json:"background" validate:"omitempty,hex_color"`
}

// ClinicResponse - DTO de respuesta
type ClinicResponse struct {
    ID          string                `json:"id"`
    Name        string                `json:"name"`
    DisplayName string                `json:"displayName"`
    Address     string                `json:"address,omitempty"`
    Phone       string                `json:"phone,omitempty"`
    Email       string                `json:"email,omitempty"`
    Website     string                `json:"website,omitempty"`
    Description string                `json:"description,omitempty"`
    Palette     ColorPaletteResponse  `json:"palette"`
    IsActive    bool                  `json:"isActive"`
    CreatedAt   time.Time             `json:"createdAt"`
    UpdatedAt   time.Time             `json:"updatedAt"`
}

// ColorPaletteResponse - Respuesta de paleta de colores
type ColorPaletteResponse struct {
    Primary    string `json:"primary"`
    Secondary  string `json:"secondary"`
    Tertiary   string `json:"tertiary"`
    Quaternary string `json:"quaternary"`
    Background string `json:"background,omitempty"`
}

// ListClinicsRequest - DTO para listar clínicas (extiende paginación)
type ListClinicsRequest struct {
    dto.PaginationRequest            // Embedding de paginación reutilizable
    IsActive              *bool      `json:"is_active" form:"is_active"`
}

// ListClinicsResponse - Respuesta específica para listado de clínicas (para Swagger)
type ListClinicsResponse struct {
    Data       []ClinicResponse       `json:"data"`
    Pagination dto.PaginationResponse `json:"pagination"`
}

// Métodos de conversión

// ToModel convierte DTO de paleta a modelo
func (p *ColorPaletteDTO) ToModel() models.ColorPalette {
    palette := models.GetDefaultPalette()

    if p.Primary != nil {
        palette.Primary = *p.Primary
    }
    if p.Secondary != nil {
        palette.Secondary = *p.Secondary
    }
    if p.Tertiary != nil {
        palette.Tertiary = *p.Tertiary
    }
    if p.Quaternary != nil {
        palette.Quaternary = *p.Quaternary
    }
    if p.Background != nil {
        palette.Background = *p.Background
    }

    return palette
}

// FromModel convierte modelo a DTO de respuesta
func FromModel(clinic *models.Clinic) ClinicResponse {
    return ClinicResponse{
        ID:          clinic.ID.Hex(),
        Name:        clinic.Name,
        DisplayName: clinic.DisplayName,
        Address:     clinic.Address,
        Phone:       clinic.Phone,
        Email:       clinic.Email,
        Website:     clinic.Website,
        Description: clinic.Description,
        Palette: ColorPaletteResponse{
            Primary:    clinic.Palette.Primary,
            Secondary:  clinic.Palette.Secondary,
            Tertiary:   clinic.Palette.Tertiary,
            Quaternary: clinic.Palette.Quaternary,
            Background: clinic.Palette.Background,
        },
        IsActive:  clinic.IsActive,
        CreatedAt: clinic.CreatedAt,
        UpdatedAt: clinic.UpdatedAt,
    }
}

// FromModels convierte slice de modelos a DTOs
func FromModels(clinics []*models.Clinic) []ClinicResponse {
    responses := make([]ClinicResponse, len(clinics))
    for i, clinic := range clinics {
        responses[i] = FromModel(clinic)
    }
    return responses
}

// ToUpdateFields convierte DTO a map para PATCH (solo campos no nulos)
func (r *UpdateClinicRequest) ToUpdateFields() map[string]interface{} {
    fields := make(map[string]interface{})

    if r.Name != nil {
        fields["name"] = strings.TrimSpace(*r.Name)
    }
    if r.DisplayName != nil {
        fields["displayName"] = strings.TrimSpace(*r.DisplayName)
    }
    if r.Address != nil {
        fields["address"] = strings.TrimSpace(*r.Address)
    }
    if r.Phone != nil {
        fields["phone"] = strings.TrimSpace(*r.Phone)
    }
    if r.Email != nil {
        fields["email"] = strings.TrimSpace(*r.Email)
    }
    if r.Website != nil {
        fields["website"] = strings.TrimSpace(*r.Website)
    }
    if r.Description != nil {
        fields["description"] = strings.TrimSpace(*r.Description)
    }
    if r.IsActive != nil {
        fields["isActive"] = *r.IsActive
    }
    if r.Palette != nil {
        fields["palette"] = r.Palette.ToModel()
    }

    return fields
}

// SetDefaults establece valores por defecto específicos de clínicas
func (r *ListClinicsRequest) SetDefaults() {
    r.PaginationRequest.SetDefaults()
    
    // Validar campo de ordenamiento específico de clínicas
    validSortFields := map[string]bool{
        "name":         true,
        "display_name": true,
        "created_at":   true,
        "updated_at":   true,
    }
    
    if r.SortBy == "" || !validSortFields[r.SortBy] {
        r.SortBy = "created_at"
    }
}
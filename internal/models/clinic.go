// internal/models/clinic.go
package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Clinic representa la entidad principal (tenant)
type Clinic struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`                         // Nombre único interno
    DisplayName string             `bson:"displayName" json:"displayName"`           // Nombre para mostrar en UI
    Address     string             `bson:"address,omitempty" json:"address,omitempty"`
    Phone       string             `bson:"phone,omitempty" json:"phone,omitempty"`
    Email       string             `bson:"email,omitempty" json:"email,omitempty"`
    Website     string             `bson:"website,omitempty" json:"website,omitempty"`
    Description string             `bson:"description,omitempty" json:"description,omitempty"`
    Palette     ColorPalette       `bson:"palette" json:"palette"`                   // Colores para UI
    IsActive    bool               `bson:"isActive" json:"isActive"`
    
    // Soft Delete simple
    DeletedAt   *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
    
    CreatedAt   time.Time  `bson:"createdAt" json:"createdAt"`
    UpdatedAt   time.Time  `bson:"updatedAt" json:"updatedAt"`
}

// ColorPalette define los colores de la clínica para el frontend
type ColorPalette struct {
    Primary   string `bson:"primary" json:"primary"`     // Color principal (#hexcode)
    Secondary string `bson:"secondary" json:"secondary"` // Color secundario
    Tertiary  string `bson:"tertiary" json:"tertiary"`   // Color tercero
	Quaternary string `bson:"quaternary" json:"quaternary"` // Color cuaternario
    Background string `bson:"background,omitempty" json:"background,omitempty"` // Color de fondo
}

// Validaciones de negocio
func (c *Clinic) IsValid() error {
    if strings.TrimSpace(c.Name) == "" {
        return ErrInvalidClinicName
    }
    if len(c.Name) < 2 || len(c.Name) > 100 {
        return ErrInvalidClinicNameLength
    }
    if strings.TrimSpace(c.DisplayName) == "" {
        return ErrInvalidDisplayName
    }
    if len(c.DisplayName) < 2 || len(c.DisplayName) > 150 {
        return ErrInvalidDisplayNameLength
    }
    if err := c.Palette.IsValid(); err != nil {
        return err
    }
    return nil
}

// IsValid valida la paleta de colores
func (p *ColorPalette) IsValid() error {
    if !isValidHexColor(p.Primary) {
        return ErrInvalidPrimaryColor
    }
    if !isValidHexColor(p.Secondary) {
        return ErrInvalidSecondaryColor
    }
    if p.Tertiary != "" && !isValidHexColor(p.Tertiary) {
        return ErrInvalidTertiaryColor
    }
    if p.Quaternary != "" && !isValidHexColor(p.Quaternary) {
        return ErrInvalidQuaternaryColor
    }
    if p.Background != "" && !isValidHexColor(p.Background) {
        return ErrInvalidBackgroundColor
    }
    return nil
}

// isValidHexColor valida formato de color hexadecimal
func isValidHexColor(color string) bool {
    if color == "" {
        return true // Campos opcionales pueden estar vacíos
    }
    if len(color) != 7 || color[0] != '#' {
        return false
    }
    
    for i := 1; i < 7; i++ {
        c := color[i]
        if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
            return false
        }
    }
    return true
}

// GetDefaultPalette retorna una paleta de colores por defecto
func GetDefaultPalette() ColorPalette {
    return ColorPalette{
        Primary:    "#3B82F6", // Blue-500
        Secondary:  "#8B5CF6", // Violet-500
        Tertiary:   "#1F2937", // Gray-800
		Quaternary: "#1F2937", // Gray-800
        Background: "#FFFFFF", // White
    }
}

// Métodos de utilidad
func (c *Clinic) IsActiveClinic() bool {
    return c.IsActive && c.DeletedAt == nil
}

func (c *Clinic) GetDisplayName() string {
    if c.DisplayName != "" {
        return c.DisplayName
    }
    return c.Name
}

func (c *Clinic) GetPrimaryColor() string {
    return c.Palette.Primary
}

// Soft Delete methods simples
func (c *Clinic) IsDeleted() bool {
    return c.DeletedAt != nil
}

func (c *Clinic) SoftDelete() {
    now := time.Now().UTC()
    c.DeletedAt = &now
    c.UpdatedAt = now
}

// Errores específicos del dominio
var (
    ErrInvalidClinicName         = errors.New("clinic name is required")
    ErrInvalidClinicNameLength   = errors.New("clinic name must be between 2 and 100 characters")
    ErrInvalidDisplayName        = errors.New("display name is required")
    ErrInvalidDisplayNameLength  = errors.New("display name must be between 2 and 150 characters")
    ErrInvalidPrimaryColor       = errors.New("invalid primary color format")
    ErrInvalidSecondaryColor     = errors.New("invalid secondary color format")
    ErrInvalidTertiaryColor      = errors.New("invalid tertiary color format")
    ErrInvalidQuaternaryColor    = errors.New("invalid quaternary color format")
    ErrInvalidBackgroundColor    = errors.New("invalid background color format")
)
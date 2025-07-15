// internal/validators/clinic_validators.go
package validators

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// RegisterClinicValidators registra validadores específicos para clínicas
func RegisterClinicValidators(validate *validator.Validate) {
    validate.RegisterValidation("hex_color", validateHexColor)
    validate.RegisterValidation("clinic_name", validateClinicName)
    validate.RegisterValidation("display_name", validateDisplayName)
}

// validateHexColor valida que el valor sea un color hexadecimal válido
func validateHexColor(fl validator.FieldLevel) bool {
    color := fl.Field().String()
    if color == "" {
        return true // Permitir vacío para campos opcionales
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

// validateClinicName valida nombres de clínica (más restrictivo para IDs internos)
func validateClinicName(fl validator.FieldLevel) bool {
    name := strings.TrimSpace(fl.Field().String())
    if len(name) < 2 || len(name) > 100 {
        return false
    }
    
    // Para nombres internos: solo alfanuméricos, guiones y puntos
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_.]+`, name)
    return matched
}

// validateDisplayName valida display names (MÁS PERMISIVO para nombres públicos)
func validateDisplayName(fl validator.FieldLevel) bool {
    displayName := strings.TrimSpace(fl.Field().String())
    if len(displayName) < 2 || len(displayName) > 150 {
        return false
    }
    
    // Validación más permisiva que permite:
    // - Letras (incluye acentos y caracteres unicode)
    // - Números
    // - Espacios
    // - Caracteres especiales comunes: - _ . ( ) & ' , 
    for _, r := range displayName {
        if !isValidDisplayNameChar(r) {
            return false
        }
    }
    
    return true
}

// isValidDisplayNameChar valida caracteres individuales para display names
func isValidDisplayNameChar(r rune) bool {
    // Permitir letras Unicode (incluye acentos)
    if unicode.IsLetter(r) {
        return true
    }
    
    // Permitir números
    if unicode.IsNumber(r) {
        return true
    }
    
    // Permitir espacios
    if unicode.IsSpace(r) {
        return true
    }
    
    // Permitir caracteres especiales específicos
    allowedSpecialChars := "-_.()&',:"
    for _, allowed := range allowedSpecialChars {
        if r == allowed {
            return true
        }
    }
    
    return false
}
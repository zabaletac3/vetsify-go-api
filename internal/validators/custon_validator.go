package validators

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Init inicializa el validador con las validaciones personalizadas
func Init() {
	validate = validator.New()
	
	// Registrar validaciones personalizadas
	validate.RegisterValidation("strong_password", validateStrongPassword)
	validate.RegisterValidation("valid_species", validateSpecies)
	validate.RegisterValidation("mongodb_id", validateMongoID)
	validate.RegisterValidation("datetime", validateDateTime)
}

// GetValidator retorna la instancia singleton del validator
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		registerCustomValidators()
	})
	return validate
}

func registerCustomValidators() {
	// Validador para contraseñas seguras
	validate.RegisterValidation("strong_password", validateStrongPassword)
	
	// Validador para especies de animales válidas
	validate.RegisterValidation("valid_species", validateSpecies)
	
	// Validador para ObjectID de MongoDB
	validate.RegisterValidation("mongodb_id", validateMongoID)
	
	// Validador para fechas en formato ISO 8601
	validate.RegisterValidation("datetime", validateDateTime)
}

// validateStrongPassword valida que la contraseña tenga al menos una mayúscula, 
// una minúscula, un número y un carácter especial
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// Verificar longitud mínima
	if len(password) < 8 {
		return false
	}
	
	// Verificar que tenga al menos una mayúscula
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return false
	}
	
	// Verificar que tenga al menos una minúscula
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return false
	}
	
	// Verificar que tenga al menos un número
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return false
	}
	
	// Verificar que tenga al menos un carácter especial
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>?]`).MatchString(password)
	if !hasSpecial {
		return false
	}
	
	return true
}

// validateSpecies valida que la especie sea una de las permitidas
func validateSpecies(fl validator.FieldLevel) bool {
	validSpecies := []string{
		"dog", "cat", "bird", "fish", "rabbit", "hamster", "guinea_pig", 
		"ferret", "reptile", "horse", "cow", "pig", "goat", "sheep",
	}
	
	species := strings.ToLower(fl.Field().String())
	
	for _, valid := range validSpecies {
		if species == valid {
			return true
		}
	}
	
	return false
}

// validateMongoID valida que el string sea un ObjectID válido de MongoDB
func validateMongoID(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	
	// Verificar que sea un ObjectID válido
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

// validateDateTime valida que el string sea una fecha válida en formato ISO 8601
func validateDateTime(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	
	// Intentar parsear en formato ISO 8601
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	
	for _, format := range formats {
		if _, err := time.Parse(format, dateStr); err == nil {
			return true
		}
	}
	
	return false
}

// GetSpeciesOptions retorna las opciones válidas para especies
func GetSpeciesOptions() []string {
	return []string{
		"dog", "cat", "bird", "fish", "rabbit", "hamster", "guinea_pig", 
		"ferret", "reptile", "horse", "cow", "pig", "goat", "sheep",
	}
}

// GetUserRoleOptions retorna las opciones válidas para roles de usuario
func GetUserRoleOptions() []string {
	return []string{"admin", "veterinarian", "assistant", "client"}
}

// GetAppointmentTypeOptions retorna las opciones válidas para tipos de cita
func GetAppointmentTypeOptions() []string {
	return []string{"consultation", "vaccination", "surgery", "checkup", "emergency", "grooming"}
}

// GetAppointmentStatusOptions retorna las opciones válidas para estados de cita
func GetAppointmentStatusOptions() []string {
	return []string{"scheduled", "confirmed", "in_progress", "completed", "cancelled"}
}
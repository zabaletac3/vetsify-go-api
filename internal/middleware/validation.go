package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/response"
	"github.com/zabaletac3/go-vet-api/internal/validators"
	"go.mongodb.org/mongo-driver/mongo"
)

// formatValidationErrors formatea los errores de validación para que sean legibles
func formatValidationErrors(err error) []response.ValidationError {
	var errors []response.ValidationError
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			var message string
			
			switch fieldError.Tag() {
			case "required":
				message = "Este campo es requerido"
			case "email":
				message = "Debe ser un email válido"
			case "min":
				if fieldError.Kind().String() == "string" {
					message = "Debe tener al menos " + fieldError.Param() + " caracteres"
				} else {
					message = "El valor mínimo es " + fieldError.Param()
				}
			case "max":
				if fieldError.Kind().String() == "string" {
					message = "No puede tener más de " + fieldError.Param() + " caracteres"
				} else {
					message = "El valor máximo es " + fieldError.Param()
				}
			case "len":
				message = "Debe tener exactamente " + fieldError.Param() + " caracteres"
			case "oneof":
				options := strings.ReplaceAll(fieldError.Param(), " ", ", ")
				message = "Debe ser uno de: " + options
			case "gt":
				message = "Debe ser mayor que " + fieldError.Param()
			case "gte":
				message = "Debe ser mayor o igual que " + fieldError.Param()
			case "lt":
				message = "Debe ser menor que " + fieldError.Param()
			case "lte":
				message = "Debe ser menor o igual que " + fieldError.Param()
			case "strong_password":
				message = "La contraseña debe tener al menos 8 caracteres, incluyendo mayúsculas, minúsculas, números y caracteres especiales"
			case "valid_species":
				message = "Especie no válida. Especies permitidas: dog, cat, bird, fish, rabbit, hamster, guinea_pig, ferret, reptile, horse, cow, pig, goat, sheep"
			case "mongodb_id":
				message = "Debe ser un ID de MongoDB válido"
			case "datetime":
				message = "Debe ser una fecha válida en formato ISO 8601"
			default:
				message = "Valor inválido"
			}
			
			// Convertir el valor a string de manera segura
			var valueStr string
			if fieldError.Value() != nil {
				if str, ok := fieldError.Value().(string); ok {
					valueStr = str
				} else {
					valueStr = "" // O podrías usar fmt.Sprintf("%v", fieldError.Value())
				}
			}
			
			errors = append(errors, response.ValidationError{
				Field:   field,
				Message: message,
				Value:   valueStr,
			})
		}
	}
	
	return errors
}

// ValidateRequest es un middleware genérico para validar requests JSON
func ValidateRequest[T any](handler func(w http.ResponseWriter, r *http.Request, req T, db *mongo.Database, logger *slog.Logger)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener dependencias del contexto
		db := r.Context().Value("db").(*mongo.Database)
		logger := r.Context().Value("logger").(*slog.Logger)
		
		var req T
		
		// Decodificar JSON
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Error decodificando JSON", "error", err, "path", r.URL.Path)
			response.JSON(w, http.StatusBadRequest, response.ErrorResponse{
				Error:   "JSON malformado",
				Message: "El formato del JSON enviado no es válido",
			})
			return
		}
		
		// Validar estructura
		validate := validators.GetValidator()
		if validate == nil {
			logger.Error("Error: validator no inicializado")
			response.JSON(w, http.StatusInternalServerError, response.ErrorResponse{
				Error:   "Error interno del servidor",
				Message: "Sistema de validación no disponible",
			})
			return
		}
		
		if err := validate.Struct(req); err != nil {
			logger.Warn("Errores de validación", 
				"errors", err.Error(), 
				"path", r.URL.Path,
				"method", r.Method,
			)
			
			validationErrors := formatValidationErrors(err)
			response.JSON(w, http.StatusBadRequest, response.ValidationErrorResponse{
				Error:   "Datos de entrada inválidos",
				Message: "Los datos enviados contienen errores de validación",
				Fields:  validationErrors,
			})
			return
		}
		
		// Llamar al handler principal con datos validados
		handler(w, r, req, db, logger)
	}
}

// ValidateRequestWithDeps es una versión que recibe las dependencias directamente
func ValidateRequestWithDeps[T any](handler func(w http.ResponseWriter, r *http.Request, req T, db *mongo.Database, logger *slog.Logger), db *mongo.Database, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req T
		
		// Decodificar JSON
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Error decodificando JSON", "error", err, "path", r.URL.Path)
			response.JSON(w, http.StatusBadRequest, response.ErrorResponse{
				Error:   "JSON malformado",
				Message: "El formato del JSON enviado no es válido",
			})
			return
		}
		
		// Validar estructura
		validate := validators.GetValidator()
		if validate == nil {
			logger.Error("Error: validator no inicializado")
			response.JSON(w, http.StatusInternalServerError, response.ErrorResponse{
				Error:   "Error interno del servidor",
				Message: "Sistema de validación no disponible",
			})
			return
		}
		
		if err := validate.Struct(req); err != nil {
			logger.Warn("Errores de validación", 
				"errors", err.Error(), 
				"path", r.URL.Path,
				"method", r.Method,
			)
			
			validationErrors := formatValidationErrors(err)
			response.JSON(w, http.StatusBadRequest, response.ValidationErrorResponse{
				Error:   "Datos de entrada inválidos",
				Message: "Los datos enviados contienen errores de validación",
				Fields:  validationErrors,
			})
			return
		}
		
		// Llamar al handler principal con datos validados
		handler(w, r, req, db, logger)
	}
}
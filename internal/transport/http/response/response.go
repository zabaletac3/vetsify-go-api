package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse representa una respuesta de error básica
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ValidationError representa un error de validación específico
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse representa una respuesta con errores de validación
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Fields  []ValidationError `json:"fields"`
}

// SuccessResponse representa una respuesta exitosa genérica
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

// MetaData contiene información de paginación
type MetaData struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
	Total       int64 `json:"total"`
}

// JSON envía una respuesta JSON con el código de estado y datos especificados
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Si falla la codificación JSON, enviar un error básico
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// Error envía una respuesta de error básica
func Error(w http.ResponseWriter, statusCode int, error string, message string) {
	JSON(w, statusCode, ErrorResponse{
		Error:   error,
		Message: message,
	})
}

// ValidationError envía una respuesta con errores de validación
func ValidationErrorRes(w http.ResponseWriter, error string, message string, fields []ValidationError) {
	JSON(w, http.StatusBadRequest, ValidationErrorResponse{
		Error:   error,
		Message: message,
		Fields:  fields,
	})
}

// Success envía una respuesta exitosa
func Success(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created envía una respuesta de recurso creado
func Created(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent envía una respuesta sin contenido
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Paginated envía una respuesta paginada
func Paginated(w http.ResponseWriter, data interface{}, currentPage, perPage, totalPages int, total int64) {
	JSON(w, http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: MetaData{
			CurrentPage: currentPage,
			PerPage:     perPage,
			TotalPages:  totalPages,
			Total:       total,
		},
	})
}

// BadRequest envía una respuesta de solicitud incorrecta
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "Solicitud incorrecta", message)
}

// Unauthorized envía una respuesta de no autorizado
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "No autorizado", message)
}

// Forbidden envía una respuesta de prohibido
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, "Prohibido", message)
}

// NotFound envía una respuesta de no encontrado
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "No encontrado", message)
}

// Conflict envía una respuesta de conflicto
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, "Conflicto", message)
}

// InternalServerError envía una respuesta de error interno del servidor
func InternalServerError(w http.ResponseWriter, message string, logger *slog.Logger, err error) {
	if logger != nil && err != nil {
		logger.Error("Error interno del servidor", "error", err, "message", message)
	}
	Error(w, http.StatusInternalServerError, "Error interno del servidor", message)
}

// UnprocessableEntity envía una respuesta de entidad no procesable
func UnprocessableEntity(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnprocessableEntity, "Entidad no procesable", message)
}

// TooManyRequests envía una respuesta de demasiadas solicitudes
func TooManyRequests(w http.ResponseWriter, message string) {
	Error(w, http.StatusTooManyRequests, "Demasiadas solicitudes", message)
}
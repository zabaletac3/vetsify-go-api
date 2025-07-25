package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zabaletac3/go-vet-api/internal/services"
)

// registerUserRequest define la estructura del cuerpo de la petición para el registro.
// Lo definimos como un tipo para poder referenciarlo en la documentación de Swagger.
type registerUserRequest struct {
	ClinicID string `json:"clinicId" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required"`
}

// Handler contiene las dependencias para los handlers de usuario, en este caso, el servicio.
type Handler struct {
	service services.UserService
}

// NewHandler es el constructor para el User Handler.
func NewHandler(svc services.UserService) *Handler {
	return &Handler{service: svc}
}

// register es el método del handler para registrar un nuevo usuario.
// @Summary      Registra un nuevo usuario
// @Description  Crea un nuevo usuario (empleado) asociado a una clínica.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      registerUserRequest  true  "Datos para el registro del usuario"
// @Success      201   {object}  models.User
// @Failure      400   {object}  map[string]string "Error: Petición inválida"
// @Failure      409   {object}  map[string]string "Error: El email ya existe"
// @Failure      500   {object}  map[string]string "Error: Error interno del servidor"
// @Router       /api/v1/users/register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	// DTO: Este struct representa el cuerpo JSON que esperamos en la petición.
	// Es específico de esta capa HTTP.
	var requestBody struct {
		ClinicID string `json:"clinicId"`
		FullName string `json:"fullName"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}
	// Aquí iría la validación del requestBody con la librería 'validator'.

	// Mapeamos el DTO de la petición a los parámetros que espera el servicio.
	// Esto desacopla la capa de servicio de la estructura de la API.
	params := services.CreateUserParams{
		ClinicID: requestBody.ClinicID,
		FullName: requestBody.FullName,
		Email:    requestBody.Email,
		Password: requestBody.Password,
		Role:     requestBody.Role,
	}

	// Llamamos a la lógica de negocio en el servicio.
	user, err := h.service.Register(r.Context(), params)
	if err != nil {
		// Verificamos si es un error de negocio específico para dar una mejor respuesta HTTP.
		if errors.Is(err, services.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		if errors.Is(err, services.ErrPasswordTooShort) {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400 Bad Request
			return
		}
		// Si es otro tipo de error, devolvemos un error de servidor genérico.
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Respondemos con el usuario creado y un código 201 Created.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
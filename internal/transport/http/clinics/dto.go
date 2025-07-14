package clinics

// CreateClinicRequest es el DTO para la petición de crear una clínica.
// Usamos tags 'json' para el mapeo y 'validate' para las reglas.
type CreateClinicRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=100"`
	Address string `json:"address" validate:"omitempty,min=5"`
	Phone   string `json:"phone" validate:"omitempty,min=7"`
}
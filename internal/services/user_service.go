package services

import (
	"context"

	"github.com/zabaletac3/go-vet-api/internal/models"
)

// CreateUserParams contiene los parámetros para registrar un nuevo usuario.
type CreateUserParams struct {
	ClinicID string
	FullName string
	Email    string
	Password string
	Role     string
}

// UserService define el contrato para la lógica de negocio de los usuarios.
type UserService interface {
	Register(ctx context.Context, params CreateUserParams) (*models.User, error)
	// Aquí irían otros métodos como Authenticate, GetByID, etc.
}
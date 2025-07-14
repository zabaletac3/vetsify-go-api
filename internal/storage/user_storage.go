package storage

import (
	"context"

	"github.com/zabaletac3/go-vet-api/internal/models"
)

// UserStorer define la interfaz para las operaciones de la colección de usuarios.
// Cada método recibe el `clinicId` para asegurar el aislamiento de los datos.
type UserStorer interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, clinicID, email string) (*models.User, error)
	FindByID(ctx context.Context, clinicID, userID string) (*models.User, error)
}
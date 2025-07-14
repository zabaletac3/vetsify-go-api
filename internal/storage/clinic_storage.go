package storage

import (
	"context"

	"github.com/zabaletac3/go-vet-api/internal/models"
)

// ClinicStorer define el contrato para las operaciones de la colección de clínicas.
type ClinicStorer interface {
	Create(ctx context.Context, clinic *models.Clinic) error
	FindByID(ctx context.Context, id string) (*models.Clinic, error)
}
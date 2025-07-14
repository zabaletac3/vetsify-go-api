package services

import (
	"context"

	"github.com/zabaletac3/go-vet-api/internal/models"
)

type CreateClinicParams struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type ClinicService interface {
	Create(ctx context.Context, params CreateClinicParams) (*models.Clinic, error)
	GetByID(ctx context.Context, id string) (*models.Clinic, error)
}
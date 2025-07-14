package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jinzhu/copier"
	"github.com/zabaletac3/go-vet-api/internal/models"
	"github.com/zabaletac3/go-vet-api/internal/storage"
)

var (
	ErrClinicNameRequired = errors.New("el nombre de la clínica es requerido")
)

type clinicService struct {
	store  storage.ClinicStorer
	logger *slog.Logger
}

func NewClinicService(store storage.ClinicStorer, logger *slog.Logger) ClinicService {
	return &clinicService{
		store:  store,
		logger: logger.With("service", "clinic"),
	}
}

func (s *clinicService) Create(ctx context.Context, params CreateClinicParams) (*models.Clinic, error) {
	if params.Name == "" {
		return nil, ErrClinicNameRequired
	}

	var newClinic models.Clinic
	copier.Copy(&newClinic, &params)

	if err := s.store.Create(ctx, &newClinic); err != nil {
		s.logger.Error("No se pudo guardar la clínica", "error", err)
		return nil, fmt.Errorf("error al persistir la clínica: %w", err)
	}

	s.logger.Info("Clínica creada exitosamente", "clinic_id", newClinic.ID.Hex())
	return &newClinic, nil
}

func (s *clinicService) GetByID(ctx context.Context, id string) (*models.Clinic, error) {
	return s.store.FindByID(ctx, id)
}
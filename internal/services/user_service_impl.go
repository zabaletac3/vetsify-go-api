package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jinzhu/copier"
	"github.com/zabaletac3/go-vet-api/internal/auth"
	"github.com/zabaletac3/go-vet-api/internal/models"
	"github.com/zabaletac3/go-vet-api/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Errores de negocio específicos para el dominio de usuarios.
var (
	ErrUserAlreadyExists = errors.New("el usuario con ese email ya existe en esta clínica")
	ErrPasswordTooShort  = errors.New("la contraseña debe tener al menos 8 caracteres")
)

type userService struct {
	userStore storage.UserStorer
	logger    *slog.Logger
}

// NewUserService es el constructor para la implementación del servicio de usuario.
func NewUserService(store storage.UserStorer, logger *slog.Logger) UserService {
	return &userService{
		userStore: store,
		logger:    logger.With("service", "user"),
	}
}

// Register implementa la lógica para registrar un nuevo usuario.
func (s *userService) Register(ctx context.Context, params CreateUserParams) (*models.User, error) {
	// 1. Regla de Negocio: Validar la contraseña.
	if len(params.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	// 2. Regla de Negocio: Verificar que el email no esté ya en uso en esa clínica.
	existingUser, err := s.userStore.FindByEmail(ctx, params.ClinicID, params.Email)
	if err != nil {
		s.logger.Error("Error al verificar el email del usuario", "error", err)
		return nil, fmt.Errorf("error al verificar el email: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 3. Lógica de Aplicación: Hashear la contraseña.
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		s.logger.Error("No se pudo hashear la contraseña", "error", err)
		return nil, fmt.Errorf("error interno al procesar la contraseña")
	}

	// 4. Mapear y Crear el Modelo.
	var newUser models.User
	copier.Copy(&newUser, &params) // Usamos copier para el mapeo limpio.
	
	clinicObjID, _ := primitive.ObjectIDFromHex(params.ClinicID)
	newUser.ClinicID = clinicObjID
	newUser.HashedPassword = hashedPassword

	// 5. Persistir el nuevo usuario.
	if err := s.userStore.Create(ctx, &newUser); err != nil {
		s.logger.Error("No se pudo guardar el usuario en la base de datos", "error", err)
		return nil, fmt.Errorf("error al registrar el usuario: %w", err)
	}

	s.logger.Info("Usuario registrado exitosamente", "email", newUser.Email, "userID", newUser.ID.Hex())
	return &newUser, nil
}
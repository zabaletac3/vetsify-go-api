package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zabaletac3/go-vet-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository implementa la interfaz UserStorer.
// Ya no es necesario volver a definir la interfaz aquí.
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios.
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Create inserta un nuevo usuario en la base de datos.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error al crear el usuario: %w", err)
	}
	return nil
}

// FindByEmail busca un usuario por su email DENTRO de una clínica específica.
func (r *UserRepository) FindByEmail(ctx context.Context, clinicID, email string) (*models.User, error) {
	clinicObjID, err := primitive.ObjectIDFromHex(clinicID)
	if err != nil {
		return nil, fmt.Errorf("ID de clínica inválido: %w", err)
	}

	filter := bson.M{"clinicId": clinicObjID, "email": email}

	var user models.User
	if err := r.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No es un error si no se encuentra.
		}
		return nil, fmt.Errorf("error al buscar usuario por email: %w", err)
	}
	return &user, nil
}

// FindByID busca un usuario por su ID DENTRO de una clínica específica.
func (r *UserRepository) FindByID(ctx context.Context, clinicID, userID string) (*models.User, error) {
	clinicObjID, err := primitive.ObjectIDFromHex(clinicID)
	if err != nil {
		return nil, fmt.Errorf("ID de clínica inválido: %w", err)
	}
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("ID de usuario inválido: %w", err)
	}

	filter := bson.M{"_id": userObjID, "clinicId": clinicObjID}

	var user models.User
	if err := r.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, fmt.Errorf("error al buscar usuario por ID: %w", err)
	}
	return &user, nil
}
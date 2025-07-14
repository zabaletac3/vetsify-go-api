package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/zabaletac3/go-vet-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClinicRepository struct {
	collection *mongo.Collection
}

func NewClinicRepository(db *mongo.Database) *ClinicRepository {
	return &ClinicRepository{
		collection: db.Collection("clinics"),
	}
}

func (r *ClinicRepository) Create(ctx context.Context, clinic *models.Clinic) error {
	clinic.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	clinic.CreatedAt = now
	clinic.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, clinic)
	if err != nil {
		return fmt.Errorf("error al crear la clínica: %w", err)
	}
	return nil
}

func (r *ClinicRepository) FindByID(ctx context.Context, id string) (*models.Clinic, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID de clínica inválido: %w", err)
	}

	var clinic models.Clinic
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&clinic); err != nil {
		return nil, fmt.Errorf("error al buscar la clínica: %w", err)
	}
	return &clinic, nil
}
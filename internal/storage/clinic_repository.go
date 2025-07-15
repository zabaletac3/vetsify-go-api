// internal/storage/clinic_repository.go
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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClinicRepository struct {
    collection *mongo.Collection
}

func NewClinicRepository(db *mongo.Database) *ClinicRepository {
    return &ClinicRepository{
        collection: db.Collection("clinics"),
    }
}

// Create - Crea una nueva clínica con validación
func (r *ClinicRepository) Create(ctx context.Context, clinic *models.Clinic) error {
    // Validar antes de insertar
    if err := clinic.IsValid(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Establecer valores automáticos
    clinic.ID = primitive.NewObjectID()
    now := time.Now().UTC()
    clinic.CreatedAt = now
    clinic.UpdatedAt = now
    clinic.IsActive = true
    clinic.DeletedAt = nil // No eliminado

    // Establecer paleta por defecto si no se proporciona
    if clinic.Palette.Primary == "" {
        clinic.Palette = models.GetDefaultPalette()
    }

    _, err := r.collection.InsertOne(ctx, clinic)
    if err != nil {
        if mongo.IsDuplicateKeyError(err) {
            return fmt.Errorf("clinic with that name already exists: %w", err)
        }
        return fmt.Errorf("failed to create clinic: %w", err)
    }

    return nil
}

// GetByID - Obtiene clínica por ID (EXCLUYE eliminadas)
func (r *ClinicRepository) GetByID(ctx context.Context, id string) (*models.Clinic, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, fmt.Errorf("invalid clinic ID '%s': %w", id, err)
    }

    filter := bson.M{
        "_id":       objID,
        "deletedAt": bson.M{"$exists": false}, // Solo no eliminados
    }

    var clinic models.Clinic
    err = r.collection.FindOne(ctx, filter).Decode(&clinic)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, fmt.Errorf("clinic with ID '%s' not found or deleted", id)
        }
        return nil, fmt.Errorf("failed to find clinic: %w", err)
    }

    return &clinic, nil
}

// Update - Actualiza clínica con validación (VERDADERO PATCH)
func (r *ClinicRepository) Update(ctx context.Context, id string, updateFields map[string]interface{}) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid clinic ID '%s': %w", id, err)
    }

    // Verificar que hay campos para actualizar
    if len(updateFields) == 0 {
        return fmt.Errorf("no fields provided for update")
    }

    // Agregar timestamp de actualización automáticamente
    updateFields["updatedAt"] = time.Now().UTC()

    // Usar $set para actualización parcial (PATCH behavior)
    update := bson.M{
        "$set": updateFields,
    }

    result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
    if err != nil {
        if mongo.IsDuplicateKeyError(err) {
            return fmt.Errorf("clinic with that name already exists: %w", err)
        }
        return fmt.Errorf("failed to update clinic: %w", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("clinic with ID '%s' not found", id)
    }

    return nil
}

// Delete - Soft delete simple (marca deletedAt)
func (r *ClinicRepository) Delete(ctx context.Context, id string) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid clinic ID '%s': %w", id, err)
    }

    now := time.Now().UTC()
    update := bson.M{
        "$set": bson.M{
            "deletedAt": now,
            "updatedAt": now,
        },
    }

    result, err := r.collection.UpdateOne(ctx, bson.M{
        "_id":       objID,
        "deletedAt": bson.M{"$exists": false}, // Solo si no está eliminado
    }, update)
    if err != nil {
        return fmt.Errorf("failed to delete clinic: %w", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("clinic with ID '%s' not found or already deleted", id)
    }

    return nil
}

// List - Lista clínicas (EXCLUYE eliminadas)
func (r *ClinicRepository) List(ctx context.Context, filters ListFilters) ([]*models.Clinic, int64, error) {
    // Construir filtro MongoDB
    filter := r.buildFilter(filters)

    // Configurar opciones de consulta
    findOptions := r.buildFindOptions(filters)

    // Ejecutar consulta principal
    cursor, err := r.collection.Find(ctx, filter, findOptions)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to execute query: %w", err)
    }
    defer cursor.Close(ctx)

    // Decodificar resultados
    var clinics []*models.Clinic
    if err := cursor.All(ctx, &clinics); err != nil {
        return nil, 0, fmt.Errorf("failed to decode results: %w", err)
    }

    // Contar total para paginación
    total, err := r.collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count documents: %w", err)
    }

    return clinics, total, nil
}

// GetByName - Busca clínica por nombre (EXCLUYE eliminadas)
func (r *ClinicRepository) GetByName(ctx context.Context, name string) (*models.Clinic, error) {
    if name == "" {
        return nil, fmt.Errorf("name cannot be empty")
    }

    filter := bson.M{
        "name": bson.M{
            "$regex":   "^" + name + "$",
            "$options": "i",
        },
        "deletedAt": bson.M{"$exists": false}, // Excluir eliminadas
    }

    var clinic models.Clinic
    err := r.collection.FindOne(ctx, filter).Decode(&clinic)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil // No es error si no se encuentra
        }
        return nil, fmt.Errorf("failed to find clinic by name: %w", err)
    }

    return &clinic, nil
}

// GetByDisplayName - Busca clínica por display name (EXCLUYE eliminadas)
func (r *ClinicRepository) GetByDisplayName(ctx context.Context, displayName string) (*models.Clinic, error) {
    if displayName == "" {
        return nil, fmt.Errorf("display name cannot be empty")
    }

    filter := bson.M{
        "displayName": bson.M{
            "$regex":   "^" + displayName + "$",
            "$options": "i",
        },
        "deletedAt": bson.M{"$exists": false}, // Excluir eliminadas
    }

    var clinic models.Clinic
    err := r.collection.FindOne(ctx, filter).Decode(&clinic)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil // No es error si no se encuentra
        }
        return nil, fmt.Errorf("failed to find clinic by display name: %w", err)
    }

    return &clinic, nil
}

// Exists - Verifica si la clínica existe (EXCLUYE eliminadas)
func (r *ClinicRepository) Exists(ctx context.Context, id string) (bool, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return false, fmt.Errorf("invalid clinic ID '%s': %w", id, err)
    }

    count, err := r.collection.CountDocuments(ctx, bson.M{
        "_id":       objID,
        "deletedAt": bson.M{"$exists": false},
    })
    if err != nil {
        return false, fmt.Errorf("failed to check existence: %w", err)
    }

    return count > 0, nil
}

// Método helper para construir filtros
func (r *ClinicRepository) buildFilter(filters ListFilters) bson.M {
    filter := bson.M{
        "deletedAt": bson.M{"$exists": false}, // SIEMPRE excluir eliminados
    }

    // Filtro de búsqueda por texto
    if filters.Search != "" {
        filter["$or"] = []bson.M{
            {"name": bson.M{"$regex": filters.Search, "$options": "i"}},
            {"displayName": bson.M{"$regex": filters.Search, "$options": "i"}},
            {"description": bson.M{"$regex": filters.Search, "$options": "i"}},
            {"address": bson.M{"$regex": filters.Search, "$options": "i"}},
        }
    }

    // Filtro por estado activo
    if filters.IsActive != nil {
        filter["isActive"] = *filters.IsActive
    }

    return filter
}

func (r *ClinicRepository) buildFindOptions(filters ListFilters) *options.FindOptions {
    opts := options.Find()

    // Ordenamiento
    sortField := "createdAt"
    if filters.SortBy != "" {
        switch filters.SortBy {
        case "created_at":
            sortField = "createdAt"
        case "updated_at":
            sortField = "updatedAt"
        case "name":
            sortField = "name"
        case "display_name":
            sortField = "displayName"
        default:
            sortField = "createdAt"
        }
    }

    sortDirection := -1 // Descendente por defecto
    if !filters.SortDesc {
        sortDirection = 1
    }

    opts.SetSort(bson.D{{Key: sortField, Value: sortDirection}})

    // Paginación
    if filters.Limit > 0 {
        opts.SetLimit(int64(filters.Limit))
        if filters.Page > 1 {
            skip := int64((filters.Page - 1) * filters.Limit)
            opts.SetSkip(skip)
        }
    }

    return opts
}
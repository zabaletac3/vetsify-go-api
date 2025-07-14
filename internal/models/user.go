package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User representa a un usuario del sistema, que siempre pertenece a una clínica.
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClinicID       primitive.ObjectID `bson:"clinicId" json:"clinicId"` // ¡El discriminador de Tenant!
	FullName       string             `bson:"fullName" json:"fullName"`
	Email          string             `bson:"email" json:"email"`
	HashedPassword string             `bson:"hashedPassword" json:"-"` // `json:"-"` para nunca exponerlo en las respuestas
	Role           string             `bson:"role" json:"role"`       // ej: "admin", "vet"
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
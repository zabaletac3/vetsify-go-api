package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Clinic representa la entidad principal, nuestro tenant.
type Clinic struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Address   string             `bson:"address,omitempty" json:"address,omitempty"`
	Phone     string             `bson:"phone,omitempty" json:"phone,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect ahora recibe el logger para un registro consistente.
func Connect(uri string, logger *slog.Logger) (*mongo.Client, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, fmt.Errorf("no se pudo preparar la conexión con mongo: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, fmt.Errorf("no se pudo conectar con mongo (ping falló): %w", err)
	}

	cleanup := func() {
		logger.Info("Cerrando conexión a MongoDB...")
		if err := client.Disconnect(context.Background()); err != nil {
			logger.Error("Fallo al desconectar de MongoDB", "error", err)
		}
	}

	return client, cleanup, nil
}
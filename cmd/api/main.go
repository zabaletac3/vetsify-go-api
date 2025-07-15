// @title       Veterinary API Multi-Tenant
// @version     1.0
// @description Interactive documentation for the veterinary API, built in Go with multi-tenant support.
//
// @contact.name   API Support
// @contact.email  support@vetapi.com
//
// @host        localhost:8080
// @BasePath    /
package main

import (
	"log/slog"
	"os"

	_ "github.com/zabaletac3/go-vet-api/docs"

	_ "github.com/zabaletac3/go-vet-api/internal/validators"

	"github.com/zabaletac3/go-vet-api/internal/config"
	"github.com/zabaletac3/go-vet-api/internal/database"
	customhttp "github.com/zabaletac3/go-vet-api/internal/transport/http"
)

func main() {

	// 1. Creamos el logger primero.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)


	// 2. Cargamos la configuración.
	cfg := config.Load()
	logger.Info("Configuración cargada exitosamente", "entorno", cfg.Env)

	// 3. Conectamos a la base de datos, pasando el logger.
	mongoClient, cleanup, err := database.Connect(cfg.MongoURI, logger)
	if err != nil {
		logger.Error("Fallo al conectar a MongoDB", "error", err)
		os.Exit(1) // Salimos si la conexión falla.
	}
	defer cleanup()
	
	logger.Info("✅ Conexión a MongoDB establecida exitosamente.")

	db := mongoClient.Database(cfg.DBName)
	logger.Info("✅ Base de datos seleccionada.", "database", db.Name())


	// 4. Creamos e iniciamos el servidor.
	server := customhttp.NewServer(cfg.Port, logger) // Pasamos el logger al servidor también.

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// <-sigChan
	// logger.Info("Cerrando servidor...")

	customhttp.SetupAllRoutes(server.Mux, db, logger,)

	server.Start()
}
// @title       API de Veterinaria Multi-Tenant
// @version     1.0
// @description Esta es la documentación interactiva para la API de la veterinaria, construida en Go.
//
// @contact.name   Tu Nombre
// @contact.email  tu.email@ejemplo.com
//
// @host        localhost:8080
// @BasePath    /
package main

import (
	"log/slog"
	"os"

	_ "github.com/zabaletac3/go-vet-api/docs"

	"github.com/zabaletac3/go-vet-api/internal/config"
	"github.com/zabaletac3/go-vet-api/internal/database"
	customhttp "github.com/zabaletac3/go-vet-api/internal/transport/http"
)

func main() {

	// 1. Creamos el logger primero.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger) // Establecemos como logger global por defecto.

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

	customhttp.SetupAllRoutes(server.Mux, db, logger)

	server.Start()
}
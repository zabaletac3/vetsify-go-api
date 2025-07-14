package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config contiene toda la configuración de la aplicación.
type Config struct {
	Port     int    `envconfig:"PORT" default:"8080"`
	Env      string `envconfig:"ENV" default:"development"`
	MongoURI string `envconfig:"MONGO_URI" required:"true"`
	DBName   string `envconfig:"DB_NAME" required:"true"`
}

// Load carga la configuración desde el archivo .env y el entorno.
func Load() *Config {
	// Carga el archivo .env. No falla si no lo encuentra.
	if err := godotenv.Load(); err != nil {
		log.Printf("Advertencia: No se pudo cargar el archivo .env: %v", err)
	}

	var cfg Config
	// Puebla el struct 'cfg' desde las variables de entorno.
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Fallo al procesar la configuración: %v", err)
	}

	return &cfg
}
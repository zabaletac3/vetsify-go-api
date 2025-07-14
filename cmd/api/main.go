package main

import (
	customhttp "github.com/zabaletac3/go-vet-api/internal/transport/http"
	"github.com/zabaletac3/go-vet-api/internal/transport/http/config"
)

func main() {
	// 1. Cargamos la configuración al inicio.
	cfg := config.Load()

	// 2. Usamos el puerto de la configuración en lugar de un valor hardcodeado.
	server := customhttp.NewServer(cfg.Port)

	// ¡Arrancamos!
	server.Start()

}
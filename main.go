package main

import (
	"fmt"
	"gognito/lib/config"
	"gognito/lib/service"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.Configuration()
	if err != nil {
		log.Fatalf("Error reading application configuration: %s", err)
	}

	service, err := service.NewService(cfg)
	if err != nil {
		log.Fatalf("Error creating new service: %s", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", cfg.AppPort),
		Handler: service.Muxer,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Error running http server: %s", err)
	}
}

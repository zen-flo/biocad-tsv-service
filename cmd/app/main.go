package main

import (
	"biocad-tsv-service/internal/config"
	"biocad-tsv-service/internal/database"
	"log"
	"os"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("failed to validate config: %v", err)
	}

	dbPool, err := database.NewPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	for _, dir := range []string{cfg.Dirs.Input, cfg.Dirs.Output} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("failed to create directory: %s, error: %v", dir, err)
			}
			log.Printf("created missing directory: %s", dir)
		}
	}

	log.Printf("Loaded config: %s", cfg)
	log.Println("Service started successfully")
}

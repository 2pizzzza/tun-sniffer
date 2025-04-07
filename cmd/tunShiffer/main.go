package main

import (
	"github.com/2pizzzza/tunShiffer/internal/app"
	"github.com/2pizzzza/tunShiffer/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app.New(cfg)
}

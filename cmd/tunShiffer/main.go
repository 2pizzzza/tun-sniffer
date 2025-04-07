package main

import (
	"github.com/2pizzzza/tunShiffer/internal/config"
	"github.com/2pizzzza/tunShiffer/internal/logger"
	"github.com/2pizzzza/tunShiffer/internal/tun"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger, err := logger.NewLogger(cfg.LogPath)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	tunHandler, err := tun.NewTunHandler(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to create TUN interface: %v", err)
	}
	defer tunHandler.Close()

	if err := tunHandler.Start(); err != nil {
		log.Fatalf("Failed to start TUN handler: %v", err)
	}
}

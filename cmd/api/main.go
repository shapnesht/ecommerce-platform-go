package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shapnesht/ecommerce-platform-go/internal/config"
	"github.com/shapnesht/ecommerce-platform-go/internal/database"
	"github.com/shapnesht/ecommerce-platform-go/internal/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed connecting to database")
	}

	mainDb, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database connection")
	}
	defer mainDb.Close()

	gin.SetMode(cfg.Server.GinMode)

	log.Info().Msg("Starting server...")
}

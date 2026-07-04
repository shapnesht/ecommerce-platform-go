package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shapnesht/ecommerce-platform-go/internal/config"
	"github.com/shapnesht/ecommerce-platform-go/internal/database"
	"github.com/shapnesht/ecommerce-platform-go/internal/logger"
	"github.com/shapnesht/ecommerce-platform-go/internal/server"
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

	srv := server.New(cfg, db, log)

	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting http server...")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown http server")
	}

}

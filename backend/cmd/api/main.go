package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/config"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/api"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/database"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()
	defer log.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	log.Info("Starting Hobby Tracker API",
		"version", cfg.App.Version,
		"environment", cfg.App.Environment,
		"port", cfg.Server.Port,
	)

	// Initialize database connection
	log.Info("Connecting to MongoDB", "uri", cfg.Database.URI, "database", cfg.Database.Name)
	db, err := database.NewMongoDB(database.MongoDBConfig{
		URI:      cfg.Database.URI,
		Name:     cfg.Database.Name,
		Timeout:  cfg.Database.Timeout,
		PoolSize: cfg.Database.PoolSize,
	})
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", "error", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Close(ctx); err != nil {
			log.Error("Error closing database connection", "error", err)
		}
	}()
	log.Info("MongoDB connected successfully")

	// TODO: Initialize repositories
	// TODO: Initialize services

	// Create server
	server := api.NewServer(cfg, log, db)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	log.Info("Server started successfully", "address", fmt.Sprintf(":%d", cfg.Server.Port))

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}

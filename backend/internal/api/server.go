package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/config"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/api/middleware"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/database"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	logger *logger.Logger
	db     *database.MongoDB
	router *gin.Engine
	server *http.Server
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, log *logger.Logger, db *database.MongoDB) *Server {
	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.CORS(cfg.Server.AllowedOrigins))

	s := &Server{
		config: cfg,
		logger: log,
		db:     db,
		router: router,
	}

	// Setup routes
	s.setupRoutes()

	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// TODO: Add auth routes
		// TODO: Add items routes
		// TODO: Add categories routes
		// TODO: Add circles routes
		// TODO: Add tags routes
		// TODO: Add import routes

		// Ping endpoint
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
				"version": s.config.App.Version,
			})
		})
	}
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	// Check database connection
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := s.db.HealthCheck(ctx); err != nil {
		s.logger.Errorw("Database health check failed", "error", err)
		dbStatus = "error"
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"app":      s.config.App.Name,
			"version":  s.config.App.Version,
			"database": dbStatus,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"app":      s.config.App.Name,
		"version":  s.config.App.Version,
		"database": dbStatus,
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	s.logger.Infof("Starting server on %s", addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

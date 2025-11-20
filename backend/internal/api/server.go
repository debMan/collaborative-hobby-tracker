package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/config"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/api/handlers"
	apimiddleware "github.com/debMan/collaborative-hobby-tracker/backend/internal/api/middleware"
	authmiddleware "github.com/debMan/collaborative-hobby-tracker/backend/internal/middleware"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	authservice "github.com/debMan/collaborative-hobby-tracker/backend/internal/service/auth"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/service/item"
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
	router.Use(apimiddleware.Logger(log))
	router.Use(apimiddleware.Recovery(log))
	router.Use(apimiddleware.CORS(cfg.Server.AllowedOrigins))

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

	// Create repositories
	userRepo := repository.NewUserRepository(s.db)
	itemRepo := repository.NewHobbyItemRepository(s.db)

	// Create auth services
	jwtExpiration := time.Duration(s.config.Auth.JWTExpiration) * time.Minute
	authService := authservice.NewService(userRepo, s.config.Auth.JWTSecret, jwtExpiration)

	// Create OAuth services
	googleOAuth := authservice.NewGoogleOAuthService(&authservice.OAuthConfig{
		ClientID:     s.config.OAuth.Google.ClientID,
		ClientSecret: s.config.OAuth.Google.ClientSecret,
		RedirectURL:  s.config.OAuth.Google.RedirectURL,
	})

	githubOAuth := authservice.NewGitHubOAuthService(&authservice.OAuthConfig{
		ClientID:     s.config.OAuth.GitHub.ClientID,
		ClientSecret: s.config.OAuth.GitHub.ClientSecret,
		RedirectURL:  s.config.OAuth.GitHub.RedirectURL,
	})

	// Create other services
	itemService := item.NewService(itemRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, googleOAuth, githubOAuth, s.config.Auth.JWTSecret)
	itemHandler := handlers.NewItemHandler(itemService)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Ping endpoint (public)
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
				"version": s.config.App.Version,
			})
		})

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/google", authHandler.GoogleAuth)
			auth.GET("/google/callback", authHandler.GoogleCallback)
			auth.GET("/github", authHandler.GitHubAuth)
			auth.GET("/github/callback", authHandler.GitHubCallback)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(authmiddleware.AuthMiddleware(s.config.Auth.JWTSecret))
		{
			// Items routes
			protected.POST("/items", itemHandler.CreateItem)
			protected.GET("/items", itemHandler.GetUserItems)
			protected.GET("/items/:id", itemHandler.GetItemByID)
			protected.PUT("/items/:id", itemHandler.UpdateItem)
			protected.DELETE("/items/:id", itemHandler.DeleteItem)
			protected.PATCH("/items/:id/toggle", itemHandler.ToggleItemCompletion)

			// TODO: Add categories routes
			// TODO: Add circles routes
			// TODO: Add tags routes
			// TODO: Add import routes
		}
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

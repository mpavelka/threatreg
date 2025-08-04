package restapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"threatreg/internal/config"
	"threatreg/internal/database"
	"threatreg/internal/restapi/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the REST API server
type Server struct {
	router *gin.Engine
	server *http.Server
}

// NewServer creates a new REST API server instance
func NewServer() *Server {
	// Set Gin mode based on environment
	if config.GetEnvironment() == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"} // In production, specify allowed origins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "Content-Length", "X-CSRF-Token"}
	router.Use(cors.New(corsConfig))

	// Add basic middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "threatreg-api",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Swagger documentation endpoints
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Alternative OpenAPI endpoints
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &Server{
		router: router,
	}

	// Setup routes
	server.setupRoutes()

	return server
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// API v1 group
	v1 := s.router.Group("/api/v1")
	{

		// Components endpoints
		components := v1.Group("/components")
		{
			components.GET("", handlers.ListComponents)
			components.POST("", handlers.CreateComponent)
			components.GET("/by-type/:type", handlers.ListComponentsByType)
			components.GET("/filter", handlers.FilterComponents)
			components.GET("/:id", handlers.GetComponent)
			components.PUT("/:id", handlers.UpdateComponent)
			components.DELETE("/:id", handlers.DeleteComponent)
			components.POST("/:id/threats", handlers.AssignThreatToComponent)
			components.GET("/:id/threats", handlers.ListThreatAssignmentsByComponent)
			components.GET("/:id/threats/with-resolution/:resolutionComponentId", handlers.ListThreatAssignmentsByComponentWithResolutionByComponent)
		}

		// Threats endpoints
		threats := v1.Group("/threats")
		{
			threats.GET("", handlers.ListThreats)
			threats.POST("", handlers.CreateThreat)
			threats.GET("/by-domain/:domainId/unresolved", handlers.ListThreatsByDomainWithUnresolvedCount)
			threats.GET("/:id", handlers.GetThreat)
			threats.PUT("/:id", handlers.UpdateThreat)
			threats.DELETE("/:id", handlers.DeleteThreat)
		}

		// Controls endpoints
		controls := v1.Group("/controls")
		{
			controls.GET("", handlers.ListControls)
			controls.POST("", handlers.CreateControl)
			controls.GET("/:id", handlers.GetControl)
			controls.PUT("/:id", handlers.UpdateControl)
			controls.DELETE("/:id", handlers.DeleteControl)
		}

		// Domains endpoints
		domains := v1.Group("/domains")
		{
			domains.GET("", handlers.ListDomains)
			domains.POST("", handlers.CreateDomain)
			domains.GET("/:id", handlers.GetDomain)
			domains.PUT("/:id", handlers.UpdateDomain)
			domains.DELETE("/:id", handlers.DeleteDomain)
		}

		// Tags endpoints
		tags := v1.Group("/tags")
		{
			tags.GET("", handlers.ListTags)
			tags.POST("", handlers.CreateTag)
			tags.GET("/:id", handlers.GetTag)
			tags.PUT("/:id", handlers.UpdateTag)
			tags.DELETE("/:id", handlers.DeleteTag)
		}

		// Threat Assignments endpoints
		threatAssignments := v1.Group("/threat-assignments")
		{
			threatAssignments.GET("/:id", handlers.GetThreatAssignment)
		}

		// Threat Resolutions endpoints
		threatResolutions := v1.Group("/threat-resolutions")
		{
			threatResolutions.POST("", handlers.CreateThreatResolution)
			threatResolutions.GET("/:id", handlers.GetThreatResolution)
			threatResolutions.PUT("/:id", handlers.UpdateThreatResolution)
			threatResolutions.DELETE("/:id", handlers.DeleteThreatResolution)
			threatResolutions.GET("/by-assignment/:assignmentId", handlers.GetThreatResolutionByThreatAssignmentID)
			threatResolutions.POST("/:id/delegate", handlers.DelegateResolution)
			threatResolutions.GET("/:id/delegated-to", handlers.GetDelegatedToResolutionByDelegatedByID)
		}
	}
}

// Run starts the REST API server
func (s *Server) Run() error {
	// Initialize database connection
	if err := database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create HTTP server
	host := config.GetAPIHost()
	port := config.GetAPIPort()
	addr := fmt.Sprintf("%s:%s", host, port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// Create channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 REST API server starting on %s", addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("🛑 Server shutting down...")

	// Shutdown server gracefully
	return s.shutdown()
}

// shutdown gracefully stops the server
func (s *Server) shutdown() error {
	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	// Close database connection
	database.Close()

	log.Println("✅ Server shutdown complete")
	return nil
}

// RunServer is a convenience function to create and run a new server
func RunServer() error {
	server := NewServer()
	return server.Run()
}

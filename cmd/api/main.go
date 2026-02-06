package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"assetManager/internal/auth"
	"assetManager/internal/config"
	"assetManager/internal/database"
	"assetManager/internal/handlers"
	"assetManager/internal/middleware"
	"assetManager/internal/repository"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	assetTypeRepo := repository.NewAssetTypeRepository(db.DB)
	assetRepo := repository.NewAssetRepository(db.DB)
	propertyRepo := repository.NewPropertyRepository(db.DB)
	assetPropertyRepo := repository.NewAssetPropertyRepository(db.DB)
	personRepo := repository.NewPersonRepository(db.DB)
	attributeRepo := repository.NewAttributeRepository(db.DB)
	personAttributeRepo := repository.NewPersonAttributeRepository(db.DB)
	assignmentRepo := repository.NewAssetAssignmentRepository(db.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	userHandler := handlers.NewUserHandler(userRepo)
	assetTypeHandler := handlers.NewAssetTypeHandler(assetTypeRepo)
	assetHandler := handlers.NewAssetHandler(assetRepo, assetPropertyRepo)
	propertyHandler := handlers.NewPropertyHandler(propertyRepo)
	personHandler := handlers.NewPersonHandler(personRepo, personAttributeRepo)
	attributeHandler := handlers.NewAttributeHandler(attributeRepo)
	assignmentHandler := handlers.NewAssignmentHandler(assignmentRepo, personRepo)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Public routes
	router.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtService))
	{
		// Auth
		api.GET("/auth/me", authHandler.Me)
		api.POST("/auth/change-password", authHandler.ChangePassword)

		// Users
		api.GET("/users", userHandler.GetAll)
		api.GET("/users/:id", userHandler.GetByID)
		api.POST("/users", userHandler.Create)
		api.PUT("/users/:id", userHandler.Update)
		api.POST("/users/:id/reset-password", userHandler.ResetPassword)
		api.DELETE("/users/:id", userHandler.Delete)

		// Asset Types
		api.GET("/asset-types", assetTypeHandler.GetAll)
		api.GET("/asset-types/:id", assetTypeHandler.GetByID)
		api.POST("/asset-types", assetTypeHandler.Create)
		api.PUT("/asset-types/:id", assetTypeHandler.Update)
		api.DELETE("/asset-types/:id", assetTypeHandler.Delete)

		// Assets
		api.GET("/assets", assetHandler.GetAll)
		api.GET("/assets/with-assignments", assetHandler.GetWithAssignments)
		api.GET("/assets/search", assetHandler.Search)
		api.GET("/assets/:id", assetHandler.GetByID)
		api.GET("/assets/by-type/:typeId", assetHandler.GetByAssetType)
		api.POST("/assets", assetHandler.Create)
		api.PUT("/assets/:id", assetHandler.Update)
		api.DELETE("/assets/:id", assetHandler.Delete)
		api.GET("/assets/:id/properties", assetHandler.GetProperties)
		api.POST("/assets/:id/properties", assetHandler.SetProperty)
		api.DELETE("/assets/:id/properties/:propId", assetHandler.DeleteProperty)

		// Properties (configuration)
		api.GET("/properties", propertyHandler.GetAll)
		api.GET("/properties/:id", propertyHandler.GetByID)
		api.POST("/properties", propertyHandler.Create)
		api.PUT("/properties/:id", propertyHandler.Update)
		api.DELETE("/properties/:id", propertyHandler.Delete)

		// Persons
		api.GET("/persons", personHandler.GetAll)
		api.GET("/persons/search", personHandler.Search)
		api.GET("/persons/:id", personHandler.GetByID)
		api.POST("/persons", personHandler.Create)
		api.PUT("/persons/:id", personHandler.Update)
		api.DELETE("/persons/:id", personHandler.Delete)
		api.GET("/persons/:id/attributes", personHandler.GetAttributes)
		api.POST("/persons/:id/attributes", personHandler.SetAttribute)
		api.DELETE("/persons/:id/attributes/:attrId", personHandler.DeleteAttribute)

		// Attributes (configuration)
		api.GET("/attributes", attributeHandler.GetAll)
		api.GET("/attributes/:id", attributeHandler.GetByID)
		api.POST("/attributes", attributeHandler.Create)
		api.PUT("/attributes/:id", attributeHandler.Update)
		api.DELETE("/attributes/:id", attributeHandler.Delete)

		// Assignments
		api.GET("/assignments/asset/:assetId", assignmentHandler.GetByAssetID)
		api.GET("/assignments/asset/:assetId/current", assignmentHandler.GetCurrentByAssetID)
		api.GET("/assignments/person/:personId", assignmentHandler.GetByPersonID)
		api.GET("/assignments/person/:personId/current", assignmentHandler.GetCurrentByPersonID)
		api.POST("/assignments", assignmentHandler.Create)
		api.POST("/assignments/assign", assignmentHandler.AssignAsset)
		api.POST("/assignments/unassign/:assetId", assignmentHandler.UnassignAsset)
		api.PUT("/assignments/:id", assignmentHandler.Update)
		api.POST("/assignments/:id/end", assignmentHandler.EndAssignment)
		api.DELETE("/assignments/:id", assignmentHandler.Delete)
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.APIPort)
	log.Printf("Starting API server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

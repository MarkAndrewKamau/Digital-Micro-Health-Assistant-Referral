package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/config"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/database"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/handlers"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/middleware"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/repository"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redis, err := database.NewRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Initialize repositories
	patientRepo := repository.NewPatientRepository(db.Pool)
	facilityRepo := repository.NewFacilityRepository(db.Pool)
	userRepo := repository.NewUserRepository(db.Pool)
	triageRepo := repository.NewTriageRepository(db.Pool)

	// Initialize services
	authService := services.NewAuthService(redis, userRepo)

	// Initialize handlers
	healthHandler := handlers.HealthCheck
	authHandler := handlers.NewAuthHandler(authService)
	patientHandler := handlers.NewPatientHandler(patientRepo)
	facilityHandler := handlers.NewFacilityHandler(facilityRepo)
	triageHandler := handlers.NewTriageHandler(triageRepo)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	// Health check endpoint (public)
	router.GET("/health", healthHandler)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Public routes
		v1.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"version": "1.0.0",
				"service": "digital-health-assistant",
			})
		})

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(authService), authHandler.Me)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Patient routes
			patients := protected.Group("/patients")
			{
				patients.POST("", patientHandler.CreatePatient)
				patients.GET("/:id", patientHandler.GetPatient)
				patients.PUT("/:id", patientHandler.UpdatePatient)
			}

			// Facility routes
			facilities := protected.Group("/facilities")
			{
				facilities.GET("", facilityHandler.ListFacilities)
				facilities.GET("/nearby", facilityHandler.GetNearbyFacilities)
				facilities.GET("/:id", facilityHandler.GetFacility)
			}

			// Triage routes
		    triage := protected.Group("/triage")
			{
				triage.POST("", triageHandler.CreateTriage)
				triage.GET("/:id", triageHandler.GetTriage)
				triage.GET("/patient/:patient_id", triageHandler.GetPatientTriages)
			}
		}
	}

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s in %s mode", port, cfg.Environment)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
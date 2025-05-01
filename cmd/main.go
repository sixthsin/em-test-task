package main

import (
	"em-task/cfg"
	"em-task/internal/clients"
	"em-task/internal/handlers"
	"em-task/internal/repository"
	"em-task/internal/services"
	"em-task/migrations"
	"em-task/pkg/db"
	"em-task/pkg/logging"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.InitLogger("info.log", "error.log")
	defer logger.Close()

	gin.DefaultWriter = logger.InfoLogger.Writer()
	gin.DefaultErrorWriter = logger.ErrorLogger.Writer()

	gin.SetMode(gin.DebugMode)

	config := cfg.LoadConfig()
	db := db.NewDB(config)
	migrations.AutoMigrate()

	repository := repository.NewRepository(db, logger)

	// Clients
	ageAPIClient := clients.NewAgeAPIClient(clients.AgeAPIClientDeps{
		Config: config,
	})
	genderAPIClient := clients.NewGenderAPIClient(clients.GenderAPIClientDeps{
		Config: config,
	})
	nationalityAPIClient := clients.NewNationalityAPIClient(clients.NationalityAPIClientDeps{
		Config: config,
	})

	// Services
	service := services.NewEnrichService(services.EnrichServiceDeps{
		Config:                    config,
		EnrichedPersonsRepository: repository,
		AgeAPIClient:              ageAPIClient,
		GenderAPIClient:           genderAPIClient,
		NationalityAPIClient:      nationalityAPIClient,
		Logger:                    logger,
	})

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorLogger.Errorf("panic recovered: %v", r)
				c.AbortWithStatusJSON(500, gin.H{"status": "error", "code": 500, "message": "Internal Server Error"})
			}
		}()
		c.Next()
	})

	// Handlers
	handlers.NewEnrichHandler(router, handlers.EnrichHandlerDeps{
		Config:        config,
		EnrichService: service,
		Logger:        logger,
	})

	log.Printf("Server started on port %s\n", config.Server.Port)

	if err := router.Run(config.Server.Port); err != nil {
		logger.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}

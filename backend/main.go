package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/qunova/protein-analyzer/internal/handlers"
	"github.com/qunova/protein-analyzer/internal/middleware"
	"github.com/qunova/protein-analyzer/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize services
	storageService := services.NewStorageService("./uploads")
	parserService := services.NewParserService()
	analysisService := services.NewAnalysisService(parserService)

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(storageService)
	analysisHandler := handlers.NewAnalysisHandler(analysisService, storageService)

	// Setup Gin router
	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Middleware
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestLogger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "protein-analyzer-api",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Upload endpoints
		api.POST("/upload", uploadHandler.UploadFile)
		
		// Analysis endpoints
		api.POST("/analyze/:fileId", analysisHandler.AnalyzeStructure)
		api.GET("/analysis/:analysisId", analysisHandler.GetAnalysisResult)
		api.GET("/analysis/:analysisId/chains", analysisHandler.GetChains)
		api.GET("/analysis/:analysisId/domains", analysisHandler.GetDomains)
		api.GET("/analysis/:analysisId/ligands", analysisHandler.GetLigands)
		
		// Export endpoints
		api.GET("/export/:analysisId", analysisHandler.ExportData)
		api.GET("/export/:analysisId/:dataType", analysisHandler.ExportSpecificData)
	}

	// Static file serving for uploads
	router.Static("/files", "./uploads")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
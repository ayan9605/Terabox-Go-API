package main

import (
	"log"
	"os"
	"terabox-api/handlers"
	"terabox-api/utils"
	"time"

	_ "terabox-api/docs" // Import generated docs

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           TeraBox API
// @version         1.0
// @description     High-performance TeraBox file downloader API with caching and proxy support
// @termsOfService  http://swagger.io/terms/

// @contact.name   Ayan Sayyad
// @contact.email  contact@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @schemes http https
func main() {
	// Initialize cache
	utils.InitCache(10 * time.Minute)

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Range"},
		ExposeHeaders:    []string{"Content-Length", "Content-Range"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
	})

	// Welcome page for GET /
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "TeraBox API",
			"version": "1.0.0",
			"status":  "online",
			"developer": "Ayan Sayyad",
			"endpoints": gin.H{
				"POST /":       "Get file info (JSON body with 'link' field)",
				"GET /api":     "Get file info (Query param ?url=terabox_link)",
				"GET /proxy":   "Download file proxy (?url=download_url&file_name=filename)",
				"GET /health":  "Health check",
				"GET /docs/index.html": "API Documentation (Swagger UI)",
			},
			"example": gin.H{
				"method": "GET",
				"url":    "/api?url=https://terabox.com/s/1abc123",
			},
		})
	})

	// API Routes
	router.POST("/", handlers.PostFileInfo)
	router.GET("/api", handlers.GetFileInfo)
	router.GET("/proxy", handlers.ProxyDownload)

	// Swagger documentation - Fixed route
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Default 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
			"available_endpoints": gin.H{
				"POST /":       "Send JSON body with 'link' field",
				"GET /api":     "Use ?url=your_terabox_link",
				"GET /proxy":   "Use ?url=download_url&file_name=filename",
				"GET /docs/index.html": "API Documentation",
			},
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📖 Swagger UI: http://localhost:%s/docs/index.html", port)
	log.Printf("💚 Health Check: http://localhost:%s/health", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

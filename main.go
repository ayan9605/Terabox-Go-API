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
	utils.InitCache(10 * time.Minute) // 10 minutes TTL

	// Set Gin mode (release for production)
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

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

	// API Routes
	router.POST("/", handlers.PostFileInfo)
	router.GET("/api", handlers.GetFileInfo)
	router.GET("/proxy", handlers.ProxyDownload)

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Default 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Method or path not allowed",
			"endpoints": gin.H{
				"POST /":       "Send JSON body with 'link' field",
				"GET /api":     "Use ?url=your_terabox_link",
				"GET /proxy":   "Use ?url=download_url&file_name=filename",
				"GET /docs":    "API Documentation (Swagger UI)",
			},
		})
	})

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

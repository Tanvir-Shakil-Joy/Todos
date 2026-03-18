// @title           Todos API
// @version         1.0
// @description     A REST API for managing todos with JWT authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@todos.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and the JWT token

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todos/config"
	"todos/database"
	"todos/middleware"
	"todos/routes"

	_ "todos/docs" // Swagger docs

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Initialize JWT middleware
	middleware.InitJWT(cfg)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("⚠️ Failed to close database connection: %v", err)
		}
	}()

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Create HTTP server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.App.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server running on port %s", cfg.App.Port)
		log.Printf("📄 Swagger UI: http://localhost:%s/swagger/index.html", cfg.App.Port)
		log.Printf("🏥 Health check: http://localhost:%s/health", cfg.App.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited")
}

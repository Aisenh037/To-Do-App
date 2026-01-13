package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/go-todo-api/internal/config"
	"github.com/user/go-todo-api/internal/database"
	"github.com/user/go-todo-api/internal/models"
	"github.com/user/go-todo-api/internal/routes"
	"github.com/user/go-todo-api/internal/worker"
)

// @title           Go Todo REST API
// @version         1.0
// @description     A comprehensive REST API for managing todos with authentication.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer " followed by your JWT token

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	database.Connect()

	// Initialize background worker
	worker.InitWorker()

	// Auto-migrate database schema
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Todo{},
		&models.RefreshToken{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	// Setup router
	router := routes.SetupRouter()

	// Create server with graceful shutdown
	srv := &http.Server{
		Addr:         ":" + config.AppConfig.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", config.AppConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give 5 seconds to finish processing requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

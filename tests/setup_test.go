package tests

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/user/go-todo-api/internal/config"
	"github.com/user/go-todo-api/internal/database"
	"github.com/user/go-todo-api/internal/models"
	"github.com/user/go-todo-api/internal/worker"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// Setup test database
	setupTestDB()

	// Initialize worker
	worker.InitWorker()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestDB()

	os.Exit(code)
}

func setupTestDB() {
	// Use in-memory SQLite for tests
	config.AppConfig = &config.Config{
		DBPath:         ":memory:",
		JWTSecret:      "test-secret",
		JWTExpiryHours: 1,
	}

	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Auto migrate
	database.DB.AutoMigrate(&models.User{}, &models.Todo{}, &models.RefreshToken{})
}

func cleanupTestDB() {
	// Close DB connection if possible, but in-memory cleans up on exit
}

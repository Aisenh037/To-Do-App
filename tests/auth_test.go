package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/user/go-todo-api/internal/handlers"
	"github.com/user/go-todo-api/internal/models"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
	}{
		{
			name:           "Success",
			payload:        map[string]string{"email": "test@example.com", "password": "password123", "name": "Test User"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Duplicate Email",
			payload:        map[string]string{"email": "test@example.com", "password": "password123", "name": "Test User"},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Invalid Input",
			payload:        map[string]string{"email": "", "password": "123", "name": ""},
			expectedStatus: http.StatusBadRequest,
		},
	}

	// Register route
	router := gin.New()
	authHandler := handlers.NewAuthHandler()
	router.POST("/register", authHandler.Register)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup router and handler
	router := gin.New()
	authHandler := handlers.NewAuthHandler()
	router.POST("/login", authHandler.Login)

	// Create a user first (manually or via endpoint) for login test
	// Note: In a real scenario, we might want to clean DB between tests or use fresh users
	// helper function to create user could be added here

	tests := []struct {
		name           string
		payload        models.LoginRequest
		expectedStatus int
	}{
		{
			name:           "Non-existent User",
			payload:        models.LoginRequest{Email: "missing@example.com", Password: "password"},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

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
)

func TestCreateTodoValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "Empty title",
			payload:        map[string]interface{}{"title": "", "description": "Some description"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing title",
			payload:        map[string]interface{}{"description": "Some description"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			handler := handlers.NewTodoHandler()

			// Set up mock user context
			router.POST("/todos", func(c *gin.Context) {
				c.Set("userID", uint(1))
				handler.Create(c)
			})

			body, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestGetTodoInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := handlers.NewTodoHandler()

	router.GET("/todos/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetByID(c)
	})

	req, _ := http.NewRequest("GET", "/todos/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

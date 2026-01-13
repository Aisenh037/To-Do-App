// Package docs provides API documentation
//
// # Go Todo REST API
//
// A comprehensive REST API for managing todos with authentication
//
//	Schemes: http, https
//	Host: localhost:8080
//	BasePath: /api
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	SecurityDefinitions:
//	  Bearer:
//	    type: apiKey
//	    name: Authorization
//	    in: header
//	    description: "Enter 'Bearer {token}'"
//
// swagger:meta
package docs

// swagger:model User
type UserDoc struct {
	// User ID
	// required: true
	// example: 1
	ID uint `json:"id"`

	// User email
	// required: true
	// example: user@example.com
	Email string `json:"email"`

	// User name
	// required: true
	// example: John Doe
	Name string `json:"name"`

	// Created timestamp
	// example: 2024-01-01T00:00:00Z
	CreatedAt string `json:"created_at"`
}

// swagger:model Todo
type TodoDoc struct {
	// Todo ID
	// required: true
	// example: 1
	ID uint `json:"id"`

	// Todo title
	// required: true
	// example: Learn Go
	Title string `json:"title"`

	// Todo description
	// example: Complete the todo API project
	Description string `json:"description"`

	// Todo status
	// enum: pending,in_progress,completed
	// example: pending
	Status string `json:"status"`

	// Due date
	// example: 2024-12-31T23:59:59Z
	DueDate *string `json:"due_date"`

	// Created timestamp
	CreatedAt string `json:"created_at"`

	// Updated timestamp
	UpdatedAt string `json:"updated_at"`
}

// swagger:model TokenPair
type TokenPairDoc struct {
	// JWT access token
	// required: true
	AccessToken string `json:"access_token"`

	// Refresh token for getting new access tokens
	// required: true
	RefreshToken string `json:"refresh_token"`

	// Seconds until access token expires
	// example: 86400
	ExpiresIn int `json:"expires_in"`
}

// swagger:model RegisterRequest
type RegisterRequestDoc struct {
	// User email
	// required: true
	// example: user@example.com
	Email string `json:"email"`

	// User password (min 6 characters)
	// required: true
	// example: password123
	Password string `json:"password"`

	// User name
	// required: true
	// example: John Doe
	Name string `json:"name"`
}

// swagger:model LoginRequest
type LoginRequestDoc struct {
	// User email
	// required: true
	// example: user@example.com
	Email string `json:"email"`

	// User password
	// required: true
	// example: password123
	Password string `json:"password"`
}

// swagger:model CreateTodoRequest
type CreateTodoRequestDoc struct {
	// Todo title
	// required: true
	// example: Learn Go
	Title string `json:"title"`

	// Todo description
	// example: Complete the todo API project
	Description string `json:"description"`

	// Todo status
	// enum: pending,in_progress,completed
	Status string `json:"status"`

	// Due date in RFC3339 format
	// example: 2024-12-31T23:59:59Z
	DueDate *string `json:"due_date"`
}

// swagger:model APIResponse
type APIResponseDoc struct {
	// Success flag
	// required: true
	Success bool `json:"success"`

	// Response message
	Message string `json:"message,omitempty"`

	// Response data
	Data interface{} `json:"data,omitempty"`

	// Error message
	Error string `json:"error,omitempty"`
}

// swagger:model PaginatedTodos
type PaginatedTodosDoc struct {
	// List of todos
	Todos []TodoDoc `json:"todos"`

	// Total number of todos
	// example: 50
	Total int64 `json:"total"`

	// Current page
	// example: 1
	Page int `json:"page"`

	// Items per page
	// example: 10
	PageSize int `json:"page_size"`

	// Total number of pages
	// example: 5
	TotalPages int `json:"total_pages"`
}

// ------- API Endpoints Documentation -------

// swagger:route POST /api/auth/register Auth registerUser
// Register a new user
//
// responses:
//   201: userTokenResponse
//   400: errorResponse
//   409: errorResponse

// swagger:route POST /api/auth/login Auth loginUser
// Login with email and password
//
// responses:
//   200: userTokenResponse
//   400: errorResponse
//   401: errorResponse

// swagger:route POST /api/auth/refresh Auth refreshToken
// Refresh access token using refresh token
//
// responses:
//   200: tokenPairResponse
//   401: errorResponse

// swagger:route GET /api/todos Todos getTodos
// Get all todos with pagination, filtering, and sorting
//
// Query parameters:
// - page: Page number (default: 1)
// - page_size: Items per page (default: 10, max: 100)
// - status: Filter by status (pending, in_progress, completed)
// - search: Search in title and description
// - sort_by: Sort field (created_at, updated_at, title, status, due_date)
// - sort_dir: Sort direction (ASC, DESC)
//
// security:
// - Bearer: []
// responses:
//   200: paginatedTodosResponse
//   401: errorResponse

// swagger:route POST /api/todos Todos createTodo
// Create a new todo
//
// security:
// - Bearer: []
// responses:
//   201: todoResponse
//   400: errorResponse
//   401: errorResponse

// swagger:route GET /api/todos/{id} Todos getTodo
// Get a single todo by ID
//
// security:
// - Bearer: []
// responses:
//   200: todoResponse
//   401: errorResponse
//   404: errorResponse

// swagger:route PUT /api/todos/{id} Todos updateTodo
// Update an existing todo
//
// security:
// - Bearer: []
// responses:
//   200: todoResponse
//   400: errorResponse
//   401: errorResponse
//   404: errorResponse

// swagger:route DELETE /api/todos/{id} Todos deleteTodo
// Delete a todo
//
// security:
// - Bearer: []
// responses:
//   200: successResponse
//   401: errorResponse
//   404: errorResponse

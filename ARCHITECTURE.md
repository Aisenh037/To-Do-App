# Architecture Documentation

## System Overview

The **Go Todo API** is a full-stack web application built with a modern, production-ready architecture. It demonstrates fundamental software engineering principles including separation of concerns, layered architecture, and security best practices.

![System Architecture](./assets/architecture.png)

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: [Gin](https://gin-gonic.com/) - High-performance HTTP web framework
- **Database**: SQLite (development) / PostgreSQL (production)
- **ORM**: [GORM](https://gorm.io/) - Feature-rich ORM library
- **Authentication**: JWT (JSON Web Tokens) with `golang-jwt/jwt`
- **API Documentation**: Swagger/OpenAPI with `swaggo`
- **Monitoring**: Prometheus metrics
- **Password Security**: bcrypt hashing

### Frontend
- **Framework**: Angular 21
- **Language**: TypeScript 5.9
- **HTTP Client**: Angular HttpClient with RxJS
- **Routing**: Angular Router with route guards
- **Styling**: Modern CSS with dark mode support

## Architecture Layers

### 1. Presentation Layer (Frontend - Angular)

**Location**: `/web`

**Components**:
- **Login Component**: Handles user authentication
- **Dashboard Component**: Main todo management interface
- **Todo List Component**: Displays and manages todos
- **JWT Service**: Manages authentication tokens

**Responsibilities**:
- User interface rendering
- User input validation
- HTTP requests to backend API
- JWT token management
- Client-side routing

### 2. Application Layer (Go Backend - API)

**Location**: `/internal`

#### 2.1 Handlers (`/internal/handlers`)
- **AuthHandler**: User registration, login, token refresh
- **TodoHandler**: CRUD operations for todos
- **ProfileHandler**: User profile management

**Responsibilities**:
- HTTP request/response handling
- Input validation
- Business logic orchestration
- Response formatting

#### 2.2 Middleware (`/internal/middleware`)
- **AuthMiddleware**: JWT token validation
- **CORSMiddleware**: Cross-origin resource sharing
- **LoggerMiddleware**: Request/response logging

**Responsibilities**:
- Request authentication
- Request logging
- Error handling
- CORS headers

#### 2.3 Repository Layer (`/internal/repository`)
- **UserRepository**: User database operations
- **TodoRepository**: Todo database operations

**Responsibilities**:
- Database query execution
- Data persistence
- Transaction management

### 3. Data Layer

**Database**: SQLite (local) / PostgreSQL (production)

**Models** (`/internal/models`):
```go
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Email     string `gorm:"uniqueIndex"`
    Password  string // bcrypt hashed
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Todo struct {
    ID          uint   `gorm:"primaryKey"`
    UserID      uint   `gorm:"index"`
    Title       string
    Description string
    Status      string // pending, in_progress, completed
    DueDate     *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## Design Patterns

### 1. Repository Pattern
Abstracts database operations behind interfaces, making the code testable and database-agnostic.

```go
type TodoRepository interface {
    Create(todo *models.Todo) error
    GetByID(id uint) (*models.Todo, error)
    GetByUserID(userID uint, filters TodoFilters) ([]models.Todo, error)
    Update(todo *models.Todo) error
    Delete(id uint) error
}
```

### 2. Middleware Pattern
Chains cross-cutting concerns like authentication, logging, and CORS.

```go
router.Use(middleware.Logger())
router.Use(middleware.CORS())

protected := router.Group("/api/todos")
protected.Use(middleware.AuthMiddleware())
```

### 3. DTO Pattern
Separates API contracts from internal models using Data Transfer Objects.

```go
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    User         UserDTO `json:"user"`
}
```

## Authentication Flow

```
1. User Registration:
   Frontend → POST /api/auth/register → Handler → Repository → Database
   ← JWT Tokens ← Hash Password ← Create User ←

2. User Login:
   Frontend → POST /api/auth/login → Handler → Verify Password
   ← JWT Tokens ← Generate Tokens ←

3. Protected Request:
   Frontend → GET /api/todos (with JWT) → Middleware (Validate JWT)
   → Handler → Repository → Database
   ← JSON Response ← Format Data ←

4. Token Refresh:
   Frontend → POST /api/auth/refresh (with Refresh Token)
   → Handler → Validate Refresh Token → Generate New Access Token
   ← New Access Token ←
```

## Data Flow

### Create Todo Example

```
┌─────────┐      HTTP POST       ┌──────────┐
│ Angular │ ──────────────────→ │   Gin    │
│Frontend │   /api/todos + JWT   │ Router   │
└─────────┘                      └──────────┘
                                       ↓
                                 ┌──────────┐
                                 │   Auth   │
                                 │Middleware│ (Verify JWT)
                                 └──────────┘
                                       ↓
                                 ┌──────────┐
                                 │   Todo   │
                                 │ Handler  │ (Validate Input)
                                 └──────────┘
                                       ↓
                                 ┌──────────┐
                                 │   Todo   │
                                 │Repository│ (Database Query)
                                 └──────────┘
                                       ↓
                                 ┌──────────┐
                                 │  GORM    │
                                 │  + DB    │ (Persist Data)
                                 └──────────┘
```

## Security Considerations

### 1. Password Security
- Passwords are hashed using bcrypt with cost factor 10
- Never stored in plain text
- Validated during login

### 2. JWT Security
- Short-lived access tokens (15 minutes)
- Long-lived refresh tokens (7 days, stored securely)
- Tokens include user ID and expiration claims
- Validated on every protected request

### 3. SQL Injection Prevention
- GORM ORM parameterizes all queries
- No raw SQL string concatenation

### 4. CORS
- Configured to allow specific origins
- Preflight requests handled properly

### 5. Input Validation
- Request binding with validation tags
- Email format validation
- Required field checks

## API Design Principles

### RESTful Conventions
- **GET** `/api/todos` - List todos
- **POST** `/api/todos` - Create todo
- **PUT** `/api/todos/:id` - Update todo
- **DELETE** `/api/todos/:id` - Delete todo

### Response Format
```json
{
  "status": "success",
  "data": {
    "todos": [...],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 45
    }
  }
}
```

### Error Handling
```json
{
  "status": "error",
  "message": "Todo not found",
  "code": "NOT_FOUND"
}
```

## Monitoring & Observability

### Prometheus Metrics
Available at `/metrics`:
- HTTP request duration
- Request count by endpoint
- Error rates
- Database query performance

### Health Check
Available at `/api/health`:
```json
{
  "status": "healthy",
  "database": "connected",
  "uptime": "2h 30m"
}
```

## Scalability Considerations

### Current Implementation
- SQLite for local development (single-file database)
- PostgreSQL ready for production (concurrent connections)
- Stateless API design (horizontal scaling ready)

### Future Enhancements
- Redis for session management and caching
- Database connection pooling
- Background job processing with worker pools
- Rate limiting middleware
- API versioning

## Testing Strategy

### Unit Tests
- Handler logic
- Middleware functionality
- Utility functions (JWT, password hashing)

### Integration Tests
- Full API endpoint workflows
- Database operations
- Authentication flows

**Test Location**: `/tests`

Run tests: `go test ./...`

## Deployment Architecture

```
                    ┌─────────────┐
                    │   Nginx /   │
                    │  Cloudflare │ (Reverse Proxy + SSL)
                    └─────────────┘
                          ↓
            ┌─────────────────────────┐
            │   Docker Container      │
            │  ┌──────────────────┐   │
            │  │  Angular Static  │   │ (Port 80)
            │  └──────────────────┘   │
            │  ┌──────────────────┐   │
            │  │   Go API Server  │   │ (Port 8080)
            │  └──────────────────┘   │
            └─────────────────────────┘
                          ↓
                    ┌─────────────┐
                    │  PostgreSQL │ (Managed DB)
                    │   Database  │
                    └─────────────┘
```

## Project Structure Rationale

```
go-todo-api/
├── cmd/api/               # Entry point - main.go
├── internal/              # Private application code
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   ├── handlers/          # HTTP handlers (controllers)
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Domain models & DTOs
│   ├── repository/        # Data access layer
│   └── routes/            # Route definitions
├── pkg/utils/             # Public utilities (reusable)
├── tests/                 # Integration tests
├── docs/                  # Swagger documentation
└── web/                   # Angular frontend
```

**Why this structure?**
- **`/cmd`**: Multiple entry points if needed (e.g., CLI, worker)
- **`/internal`**: Prevents import by external projects
- **`/pkg`**: Reusable packages (JWT, password utils)
- **Separation of layers**: Easy to test, maintain, and scale

## Conclusion

This architecture prioritizes:
✅ **Maintainability** - Clear separation of concerns  
✅ **Testability** - Repository pattern, dependency injection  
✅ **Security** - JWT auth, bcrypt, input validation  
✅ **Scalability** - Stateless API, ready for horizontal scaling  
✅ **Developer Experience** - Swagger docs, clear structure, comprehensive error handling

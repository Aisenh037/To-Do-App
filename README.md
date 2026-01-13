# Go Todo REST API

A comprehensive REST API for managing todos, built with Go to learn fundamental software development concepts.

## 🚀 Features

- **User Authentication** - JWT-based registration and login
- **Todo CRUD** - Create, read, update, delete todos
- **Protected Routes** - Middleware-based authentication
- **PostgreSQL** - Persistent data storage with GORM
- **Docker Ready** - Containerized deployment

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | [Gin](https://gin-gonic.com/) |
| Database | PostgreSQL |
| ORM | [GORM](https://gorm.io/) |
| Auth | JWT (golang-jwt/jwt) |
| Container | Docker & Docker Compose |

## 📁 Project Structure

```
go-todo-api/
├── cmd/api/            # Application entry point
├── internal/
│   ├── config/         # Environment configuration
│   ├── database/       # Database connection
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # Auth & logging middleware
│   ├── models/         # Data models & DTOs
│   ├── repository/     # Database operations
│   └── routes/         # Route definitions
├── pkg/utils/          # Shared utilities (JWT, password, responses)
├── tests/              # Unit tests
├── docker-compose.yml  # Docker services
└── Dockerfile          # Container build
```

## 🏃 Quick Start

### Prerequisites
- Go 1.21+

### Local Development (Zero Setup)

I've switched the database to **SQLite**, so you don't need to install PostgreSQL or start Docker anymore! Just follow these steps:

```bash
# 1. Clone/Navigate to the project
cd go-todo-api

# 2. Install dependencies
go mod download

# 3. Build and Run
go run ./cmd/api
```

The database file `todo.db` will be created automatically in your project folder.

### Frontend (Angular)

The frontend is located in the `web` directory.

```bash
# 1. Navigate to web directory
cd web

# 2. Install dependencies (if not done)
npm install

# 3. Start the dev server
npm start
```

The application will be available at `http://localhost:4200`.


## 📡 API Endpoints

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login and get tokens |
| POST | `/api/auth/refresh` | Refresh access token |

### Todos (Protected - requires Bearer token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/todos` | Get todos with pagination |
| GET | `/api/todos/:id` | Get single todo |
| POST | `/api/todos` | Create todo |
| PUT | `/api/todos/:id` | Update todo |
| DELETE | `/api/todos/:id` | Delete todo |

**Todo Query Parameters:**
- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)
- `status` - Filter: pending, in_progress, completed
- `search` - Search in title/description
- `sort_by` - created_at, title, status, due_date
- `sort_dir` - ASC or DESC

### User (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/profile` | Get user profile |
| POST | `/api/auth/logout` | Logout (revoke tokens) |

## 📝 Usage Examples

### Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Create a Todo (with token)
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Learn Go",
    "description": "Complete the todo API project",
    "status": "in_progress"
  }'
```

### Get All Todos
```bash
curl http://localhost:8080/api/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🧪 Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run integration tests
go test -v ./tests/...

# Run specific package tests
go test ./pkg/utils/...
```

## 🎓 Concepts Covered

This project teaches you:

**Phase 1 - Basics:**
- ✅ Go project structure and modules
- ✅ REST API design with Gin framework
- ✅ Database integration with GORM ORM
- ✅ JWT authentication and middleware
- ✅ Password hashing with bcrypt
- ✅ Request validation
- ✅ Error handling patterns
- ✅ Unit testing in Go
- ✅ Docker containerization

**Phase 3 - Production Extensions:**
- ✅ Interactive Swagger UI
- ✅ Background task worker (Goroutines + Channels)
- ✅ Scheduled reminders (time.Ticker)
- ✅ Prometheus monitoring metrics
- ✅ Centralized error mapping

## 📡 API Portal

- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Metrics**: [http://localhost:8080/metrics](http://localhost:8080/metrics)
- **Health**: [http://localhost:8080/api/health](http://localhost:8080/api/health)

## 📄 License

MIT License - feel free to use for learning!



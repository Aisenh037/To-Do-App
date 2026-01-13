# Deployment Guide

This guide provides instructions for deploying the Go Todo API application in various environments.

## Table of Contents
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Production Deployment](#production-deployment)
- [Environment Variables](#environment-variables)
- [CI/CD Setup](#cicd-setup)

## Local Development

### Prerequisites
- Go 1.21 or higher
- Node.js 18+ and npm
- Git

### Quick Start

#### 1. Clone the Repository
```bash
git clone https://github.com/Aisenh037/To-Do-App.git
cd To-Do-App
```

#### 2. Start Backend
```bash
# Install Go dependencies
go mod download

# Run the server (uses SQLite by default)
go run ./cmd/api
```

Backend will be available at `http://localhost:8080`

#### 3. Start Frontend
```bash
# Navigate to web directory
cd web

# Install dependencies
npm install

# Start development server
npm start
```

Frontend will be available at `http://localhost:4200`

#### 4. Access Application
- **Frontend**: http://localhost:4200
- **Swagger API Docs**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/api/health
- **Metrics**: http://localhost:8080/metrics

## Docker Deployment

### Using Docker Compose (Recommended)

#### 1. Prerequisites
- Docker Desktop
- Docker Compose

#### 2. Create Environment File
```bash
cp .env.example .env
# Edit .env with your configuration
```

#### 3. Start Services
```bash
docker-compose up -d
```

This will start:
- **Backend API** on port 8080
- **PostgreSQL Database** on port 5432
- **Angular Frontend** (if configured)

#### 4. Stop Services
```bash
docker-compose down
```

#### 5. View Logs
```bash
docker-compose logs -f
```

### Building Docker Image Manually

```bash
# Build backend image
docker build -t go-todo-api:latest .

# Run backend container
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=yourpassword \
  -e DB_NAME=tododb \
  -e JWT_SECRET=your-secret-key \
  --name todo-api \
  go-todo-api:latest
```

## Production Deployment

### Option 1: Railway (Recommended for Beginners)

1. **Create Railway Account**: https://railway.app
2. **Install Railway CLI**:
   ```bash
   npm install -g @railway/cli
   ```
3. **Login**:
   ```bash
   railway login
   ```
4. **Initialize Project**:
   ```bash
   railway init
   ```
5. **Add PostgreSQL**:
   ```bash
   railway add postgresql
   ```
6. **Deploy**:
   ```bash
   railway up
   ```
7. **Set Environment Variables** in Railway dashboard

### Option 2: Render

1. **Create Render Account**: https://render.com
2. **Create New Web Service**
3. **Connect GitHub Repository**
4. **Configure Build Settings**:
   - Build Command: `go build -o api ./cmd/api`
   - Start Command: `./api`
5. **Add PostgreSQL Database** from Render dashboard
6. **Set Environment Variables**
7. **Deploy**

### Option 3: DigitalOcean App Platform

1. **Create DigitalOcean Account**
2. **Create New App** → Import from GitHub
3. **Configure App**:
   - Detect Go application automatically
   - Add PostgreSQL database
4. **Set Environment Variables**
5. **Deploy**

### Option 4: AWS (Advanced)

#### Backend (Elastic Beanstalk)
```bash
# Install EB CLI
pip install awsebcli

# Initialize
eb init -p go go-todo-api

# Create environment
eb create production

# Deploy
eb deploy
```

#### Database (RDS PostgreSQL)
- Create RDS PostgreSQL instance
- Configure security groups
- Update environment variables in Elastic Beanstalk

#### Frontend (S3 + CloudFront)
```bash
cd web
npm run build

# Upload to S3
aws s3 sync dist/web s3://your-bucket-name

# Configure CloudFront distribution
```

## Environment Variables

### Required Variables

Create a `.env` file in the project root:

```bash
# Database Configuration
DB_HOST=localhost           # Database host
DB_PORT=5432               # Database port
DB_USER=postgres           # Database username
DB_PASSWORD=yourpassword   # Database password
DB_NAME=tododb             # Database name

# Use SQLite (for local development)
USE_SQLITE=true            # Set to false for PostgreSQL
SQLITE_PATH=./todo.db      # SQLite database file path

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this  # Change in production!
JWT_ACCESS_EXPIRY=15m      # Access token expiration (15 minutes)
JWT_REFRESH_EXPIRY=168h    # Refresh token expiration (7 days)

# Server Configuration
SERVER_PORT=8080           # API server port
GIN_MODE=release           # Use "release" for production

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:4200,https://yourdomain.com

# Feature Flags
ENABLE_SWAGGER=true        # Enable Swagger UI
ENABLE_METRICS=true        # Enable Prometheus metrics
```

### Production Security Checklist

- [ ] Change `JWT_SECRET` to a strong random string (min 32 characters)
- [ ] Use PostgreSQL instead of SQLite
- [ ] Set `GIN_MODE=release`
- [ ] Use HTTPS for all endpoints
- [ ] Restrict CORS to specific origins
- [ ] Use environment-specific secrets (not in code)
- [ ] Enable database connection pooling
- [ ] Set up proper logging and monitoring

## Database Migration

### PostgreSQL Setup (Production)

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE tododb;

# Exit psql
\q
```

The application will automatically create tables on first run using GORM auto-migration.

### Manual Migration (if needed)

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Todos table
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_status ON todos(status);
```

## CI/CD Setup

### GitHub Actions

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Production

on:
  push:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test -v ./...
      
      - name: Build
        run: go build -v ./cmd/api

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to Railway
        run: |
          npm install -g @railway/cli
          railway up
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
```

## Performance Optimization

### Backend
```go
// Enable GZIP compression
router.Use(gzip.Gzip(gzip.DefaultCompression))

// Database connection pooling
db.DB().SetMaxIdleConns(10)
db.DB().SetMaxOpenConns(100)
db.DB().SetConnMaxLifetime(time.Hour)

// Enable caching headers
router.Use(func(c *gin.Context) {
    c.Header("Cache-Control", "public, max-age=300")
})
```

### Frontend
```bash
# Build with production optimizations
ng build --configuration production

# Enable AOT compilation (default in production)
# Enable tree shaking
# Minify CSS and JS
```

## Monitoring

### Health Checks

```bash
# Check API health
curl http://localhost:8080/api/health

# Response:
{
  "status": "healthy",
  "database": "connected",
  "uptime": "2h30m15s"
}
```

### Prometheus Metrics

```bash
# Access metrics
curl http://localhost:8080/metrics
```

Set up Grafana dashboard for visualization.

## Troubleshooting

### Backend won't start
```bash
# Check if port 8080 is in use
netstat -ano | findstr :8080    # Windows
lsof -i :8080                   # Linux/Mac

# Kill process if needed
taskkill /PID <PID> /F          # Windows
kill -9 <PID>                   # Linux/Mac
```

### Database connection fails
- Verify PostgreSQL is running
- Check credentials in `.env`
- Ensure database exists
- Check firewall settings

### Frontend can't connect to backend
- Verify backend is running on port 8080
- Check CORS configuration
- Verify API base URL in Angular environment files

### Docker build fails
```bash
# Clear Docker cache
docker builder prune

# Rebuild without cache
docker-compose build --no-cache
```

## Backup & Recovery

### Database Backup
```bash
# PostgreSQL backup
pg_dump tododb > backup_$(date +%Y%m%d).sql

# Restore
psql tododb < backup_20260113.sql
```

## SSL/TLS Configuration

### Using Nginx as Reverse Proxy

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        proxy_pass http://localhost:4200;
    }
}
```

## Support

For deployment issues:
1. Check the logs: `docker-compose logs -f`
2. Verify environment variables
3. Test database connectivity
4. Review error messages in Swagger UI

## Production Checklist

- [ ] Database backups configured
- [ ] SSL/TLS enabled
- [ ] Environment variables secured
- [ ] Monitoring and alerts set up
- [ ] Rate limiting configured
- [ ] Database connection pooling enabled
- [ ] CORS properly configured
- [ ] Error logging to external service
- [ ] Health checks configured
- [ ] Documentation updated

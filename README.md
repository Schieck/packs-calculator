# Packs Calculator

A complete full-stack application for calculating optimal pack allocations.

## 🚀 Getting Started Locally

### Prerequisites

Before running the application, ensure you have the following installed:

- **Docker Desktop**: Latest version with Docker Compose
  - Download from [docker.com](https://www.docker.com/products/docker-desktop/)
  - Verify: `docker --version` and `docker compose version`
- **Go**: Version 1.24.3 or later (for local development)
  - Download from [golang.org](https://golang.org/dl/)
  - Verify: `go version`
- **Make**: For running automation commands
  - macOS: `brew install make` (if not already installed)
  - Linux: Usually pre-installed
  - Windows: Use WSL or install via chocolatey
- **Git**: For cloning the repository

### Complete Setup (Recommended)

**Step 1: Clone the Repository**
```bash
git clone https://github.com/Schieck/packs-calculator
cd packs-calculator
```

**Step 2: Complete Development Setup**
```bash
# This single command will:
# - Install development tools (Swagger CLI)
# - Download Go dependencies
# - Generate API documentation
# - Start all Docker services (database + backend + frontend)
make setup-dev
```

**Step 3: Verify Everything is Running**

Once `make setup-dev` completes, you should see:
- ✅ Backend API: http://localhost:8080
- ✅ API Documentation: http://localhost:8080/swagger/index.html
- ✅ Frontend: http://localhost:5173
- ✅ Database: Running on port 5432

**Step 4: Test the API**

Get an authentication token:
```bash
curl -X POST http://localhost:8080/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{"secret": "-your-auth-secret-change-in-production"}'
```

### Alternative Setup Methods

#### Option 1: Docker Only
```bash
# Start all services without local development tools
docker compose up -d

# View logs
docker compose logs -f
```

#### Option 2: Local Development with Containerized Database
```bash
# Start only the database
docker compose up -d db

# Run backend locally (in another terminal)
cd backend
go mod download
go run cmd/server/main.go
```

### Available Commands

```bash
# Essential commands
make help           # Show all available commands
make setup-dev      # Complete development setup (recommended for first time)

# Development modes
make dev            # Hot-reload development (stops Docker backend, keeps database)
make docker-dev     # Switch back to containerized development
make docker-up      # Start all services from scratch
make docker-down    # Stop all services

# Testing & building
make test           # Run all tests
make build          # Build the application
make clean          # Clean build artifacts
```

### Environment Configuration

The application uses the following default configuration (defined in `docker-compose.yml`):

```bash
# Database
POSTGRES_USER=packer
POSTGRES_PASSWORD=secret
POSTGRES_DB=packs

# Authentication (change in production!)
JWT_SECRET=your-jwt-secret-change-in-production
AUTH_SECRET=your-auth-secret-change-in-production

# Server
PORT=8080
```

For production deployment, create a `.env` file with your own values.

### Troubleshooting

**Problem: `docker compose: command not found`**
- Solution: Update Docker Desktop to latest version, or use `docker-compose` instead

**Problem: `Port already in use`**
- Solution: Stop existing services with `make docker-down` or change ports in `docker-compose.yml`

**Problem: Database connection failed**
- Solution: Ensure database is running with `docker compose up -d db` before running `make dev`

**Problem: `make: command not found`**
- Solution: Install make or run commands directly:
  ```bash
  # Instead of make setup-dev
  go install github.com/swaggo/swag/cmd/swag@latest
  cd backend && go mod download && go mod tidy
  docker compose up -d
  ```

**Problem: Hot-reload not working**
- Solution: Ensure `air` is installed with `go install github.com/air-verse/air@latest`
- Check that you're editing `.go` files (air ignores test files by default)
- Verify the `backend/.air.toml` configuration file exists

### Next Steps

1. **Explore the API**: Visit http://localhost:8080/swagger/index.html
2. **Test calculations**: Use the calculator endpoints
3. **Check the frontend**: Open http://localhost:5173
4. **Review the code**: Start with `backend/cmd/server/main.go`
5. **Run tests**: Execute `make test` to ensure everything works

## 📋 Features

- **🔐 Simple JWT Authentication**: Token-based auth using environment secrets
- **🧮 Pack Calculator**: Advanced algorithm for optimal pack allocation
- **⚙️ Configuration Management**: Save and reuse pack size configurations
- **📊 Real-time Calculations**: Instant pack optimization
- **🐳 Docker Support**: Full containerization with Docker Compose
- **📖 API Documentation**: Comprehensive Swagger/OpenAPI docs
- **🧪 Comprehensive Testing**: Unit tests with coverage reports
- **📝 Structured Logging**: JSON-based logging with request tracing

## 🏗️ Architecture

### Modular Backend Structure
```
backend/
├── cmd/server/          # Application entry point
├── internal/            # Private modules
│   ├── auth/           # JWT authentication
│   └── calculator/     # Pack calculation logic
└── pkg/                # Shared utilities
    ├── middleware/     # HTTP middleware
    ├── db/            # Database management
    └── response/      # Response utilities
```

### Technology Stack

**Backend:**
- **Language**: Go 1.23+
- **Framework**: Gin
- **Database**: PostgreSQL 16
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI
- **Testing**: Testify

**Frontend:**
- **Framework**: [Framework details in front-end directory]

**Infrastructure:**
- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL with auto-migration
- **Networking**: Docker networks with health checks

## 🔧 Environment Configuration

### Required Environment Variables

```bash
# Authentication
JWT_SECRET=your-jwt-secret-change-in-production
AUTH_SECRET=your-auth-secret-change-in-production

# Database
DB_DSN=postgres://packer:secret@localhost:5432/packs?sslmode=disable

# Server
PORT=8080
GIN_MODE=release
```

### Creating Environment File

For production deployment, create a `.env` file with your custom values:

```bash
# Copy the example and customize
cp docker-compose.yml .env.production
# Edit .env.production with your values
```

## 🛠️ Development

### Development Workflow

Once you've completed the initial setup with `make setup-dev`, you can choose between two development modes:

**🔥 Hot-Reload Development (Recommended for Active Development):**
```bash
# Stops backend container, keeps database running, starts local dev server
make dev
```

**🐳 Full Containerized Development:**
```bash
# Switch back to full Docker setup (both backend + database)
make docker-dev

# Or start everything from scratch
make docker-up
```

**Switching Between Modes:**
```bash
# Currently in Docker mode? Switch to hot-reload:
make dev

# Currently in hot-reload mode? Switch back to Docker:
make docker-dev
```

The `make dev` command uses **Air** for automatic hot-reloading:
- 🔥 **Instant Rebuilds**: Changes to `.go` files trigger automatic rebuilds
- ⚡ **Fast Restart**: Server restarts automatically after successful builds
- 📝 **Build Logs**: Compilation errors are logged to `build-errors.log`
- 🎯 **Smart Watching**: Only watches relevant files (excludes tests, docs, tmp files)
- 🛑 **No Port Conflicts**: Automatically stops Docker backend container

**Testing and Building:**
```bash
make test           # Run all tests
make build          # Build application
make clean          # Clean build artifacts
```

## 📡 API Endpoints

### Authentication
```http
POST /api/v1/auth/token          # Get JWT token
```

### Pack Calculator
```http
POST /api/v1/calculator/calculate           # Calculate optimal packs
GET  /api/v1/calculator/configurations      # List configurations
POST /api/v1/calculator/configurations      # Create configuration
PUT  /api/v1/calculator/configurations/{id} # Update configuration
DELETE /api/v1/calculator/configurations/{id} # Delete configuration
```

### System
```http
GET /health          # Health check
GET /swagger/*       # API documentation
```

## 🔐 Authentication

### Getting a Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{"secret": "your-auth-secret"}'
```

### Using the Token

```bash
curl -X GET http://localhost:8080/api/v1/calculator/configurations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🧮 Pack Calculation Algorithm

The calculator uses an advanced dynamic programming algorithm that:

1. **Minimizes Pack Count**: Finds the optimal combination using DP
2. **Prevents Undersupply**: Always meets or exceeds requirements  
3. **Handles Edge Cases**: Works with any pack size combination
4. **Optimizes Efficiency**: Falls back to greedy when needed

### Example Calculation

**Input**: 251 items, pack sizes [250, 500, 1000, 2000, 5000]

**Output**: 
- ✅ **Optimal**: 1 pack of 500 (500 items, 1 pack)
- ❌ **Suboptimal**: 2 packs of 250 (500 items, 2 packs)

### API Usage

```bash
curl -X POST http://localhost:8080/api/v1/calculator/calculate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "items": 251,
    "pack_sizes": [250, 500, 1000, 2000, 5000]
  }'
```

## 🧪 Testing

```bash
# Run all tests
make test

# Specific module
cd backend && go test -v ./internal/calculator
```

### Test Coverage

The project maintains high test coverage across:
- ✅ Pack calculation algorithms
- ✅ JWT middleware
- ✅ API endpoints
- ✅ Database operations
- ✅ Request validation

## 🚀 Deployment

### Using Docker (Recommended)

```bash
# Start all services
make docker-up

# Check health
curl http://localhost:8080/health
```

### Manual Deployment

```bash
# Build and run
make build
cd backend && ./app
```

## 🔧 Configuration Management

### Pack Configurations

Save and reuse common pack size combinations:

```bash
# Create configuration
curl -X POST http://localhost:8080/api/v1/calculator/configurations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Standard Packs",
    "description": "Common pack sizes",
    "pack_sizes": [250, 500, 1000, 2000, 5000]
  }'

# Use configuration in calculation
curl -X POST http://localhost:8080/api/v1/calculator/calculate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "items": 1250,
    "configuration_id": 1
  }'
```

## 📊 Monitoring & Logging

### Structured Logging

The application uses structured JSON logging:

```json
{
  "time": "2023-01-01T12:00:00Z",
  "level": "INFO",
  "msg": "HTTP Request",
  "method": "POST",
  "path": "/api/v1/calculator/calculate",
  "status": 200,
  "latency": "45ms",
  "ip": "127.0.0.1"
}
```

### Health Checks

- **Database**: Connection and query verification
- **Application**: Service availability
- **Docker**: Container health monitoring

## 🔒 Security Features

- **JWT Authentication**: Secure token-based auth
- **Input Validation**: Comprehensive request validation
- **CORS Protection**: Configurable cross-origin policies
- **SQL Injection Prevention**: Parameterized queries
- **Container Security**: Non-root user containers
- **Secret Management**: Environment-based secrets

## 🤝 Contributing

1. **Fork the repository**
2. **Setup development environment**: `make setup-dev`
3. **Make changes and add tests**
4. **Run tests**: `make test`
5. **Submit Pull Request**

## 📚 Documentation

- **API Documentation**: http://localhost:8080/swagger/index.html
- **Backend README**: [backend/README.md](backend/README.md)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Made with ❤️ for optimal pack calculations**

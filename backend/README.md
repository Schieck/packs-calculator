# Packs Calculator API

A RESTful API built with Go and Gin framework for calculating optimal pack allocations. The API includes simple JWT authentication, pack configuration management, and an advanced pack calculation engine.

## Features

- **Simple JWT Authentication**: Token-based authentication using a secret key
- **Pack Calculation**: Optimal pack allocation algorithm that minimizes waste
- **Configuration Management**: Save and reuse pack configurations (global)
- **OpenAPI Documentation**: Comprehensive API documentation with Swagger UI
- **Structured Logging**: JSON-based logging with slog
- **Request Validation**: Automatic request validation using go-playground/validator
- **Health Checks**: Built-in health check endpoints
- **Database**: PostgreSQL with automatic schema creation
- **Docker Support**: Multi-stage Docker build with security best practices

## Tech Stack

- **Framework**: Gin (Go web framework)
- **Database**: PostgreSQL with lib/pq driver
- **Authentication**: Simple JWT with golang-jwt/jwt
- **Validation**: go-playground/validator
- **Documentation**: Swagger with swaggo
- **Testing**: Testify framework
- **Logging**: Go's built-in slog package
- **Containerization**: Docker with Alpine Linux

## API Endpoints

### Authentication
- `POST /api/v1/auth/token` - Get JWT token using secret

### Pack Calculator
- `POST /api/v1/calculator/calculate` - Calculate optimal pack allocation
- `GET /api/v1/calculator/configurations` - Get all configurations
- `POST /api/v1/calculator/configurations` - Create new configuration
- `PUT /api/v1/calculator/configurations/{id}` - Update configuration
- `DELETE /api/v1/calculator/configurations/{id}` - Delete configuration

### Health Check
- `GET /health` - Health check endpoint

### Documentation
- `GET /swagger/index.html` - Swagger UI documentation

## Quick Start

### Prerequisites

- Go 1.23.4 or later
- PostgreSQL 16+
- Docker (optional)

### Using Make (Recommended)

```bash
# Complete development setup
make setup-dev

# Start with Docker
make docker-up

# Or run locally
make dev
```

### Manual Setup

1. **Set up environment variables**:

```bash
export DB_DSN="postgres://packer:secret@localhost:5432/packs?sslmode=disable"
export JWT_SECRET="your-jwt-secret-change-in-production"
export AUTH_SECRET="your-auth-secret-change-in-production"
export PORT="8080"
```

2. **Install dependencies**:

```bash
go mod download
```

3. **Generate Swagger documentation**:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o ./docs
```

4. **Run the application**:

```bash
go run cmd/server/main.go
```

## Authentication

The API uses simple JWT authentication. To get a token:

1. **Get JWT Token**:

```bash
curl -X POST http://localhost:8080/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "secret": "your-auth-secret-change-in-production"
  }'
```

2. **Use the token in subsequent requests**:

```bash
curl -X GET http://localhost:8080/api/v1/calculator/configurations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Pack Calculation Algorithm

The API implements an intelligent pack calculation algorithm that:

1. **Minimizes Pack Count**: Uses dynamic programming to find the optimal combination
2. **Prevents Undersupply**: Always meets or exceeds the required items
3. **Handles Edge Cases**: Works with any combination of pack sizes
4. **Optimizes for Efficiency**: Greedy approach with fallback to advanced DP when needed

### Example

For 251 items with pack sizes [250, 500, 1000, 2000, 5000]:
- **Optimal**: 1 pack of 500 (total: 500 items, 1 pack)
- **Suboptimal**: 2 packs of 250 (total: 500 items, 2 packs)

## Testing

### Run Unit Tests

```bash
make test
```

### Run Specific Test Suite

```bash
go test -v ./internal/calculator
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/                    # Private application modules
│   ├── auth/                    # Authentication module
│   │   ├── model.go            # JWT claims, request/response DTOs
│   │   ├── service.go          # Auth business logic
│   │   └── handler.go          # HTTP handlers for auth endpoints
│   └── calculator/             # Pack calculation module
│       ├── model.go            # Configuration, pack allocation DTOs
│       ├── service.go          # Calculation algorithms
│       ├── repo.go             # Configuration database operations
│       ├── handler.go          # HTTP handlers for calculator endpoints
│       └── calculator_test.go  # Comprehensive test suite
├── pkg/                        # Shared/reusable packages
│   ├── middleware/             # HTTP middleware
│   │   ├── middleware.go       # JWT, CORS, logging, recovery
│   │   └── middleware_test.go  # Middleware test suite
│   ├── db/                     # Database connection management
│   │   └── db.go               # Connection setup, table creation
│   └── response/               # HTTP response utilities
│       └── response.go         # Error/success responses, validation
├── docs/                       # Generated Swagger docs
├── Dockerfile                  # Docker build configuration
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
└── README.md                   # This file
```

## API Usage Examples

### Get Authentication Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "secret": "your-auth-secret-change-in-production"
  }'
```

### Calculate Packs

```bash
curl -X POST http://localhost:8080/api/v1/calculator/calculate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "items": 251,
    "pack_sizes": [250, 500, 1000, 2000, 5000]
  }'
```

### Create Configuration

```bash
curl -X POST http://localhost:8080/api/v1/calculator/configurations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Standard Packs",
    "description": "Standard pack sizes for general use",
    "pack_sizes": [250, 500, 1000, 2000, 5000]
  }'
```

## Database Schema

### Configurations Table
```sql
CREATE TABLE configurations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    pack_sizes INTEGER[] NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_DSN` | `postgres://packer:secret@localhost:5432/packs?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | `development-jwt-secret-change-in-production` | JWT signing secret |
| `AUTH_SECRET` | `development-auth-secret-change-in-production` | Authentication secret for token generation |
| `GIN_MODE` | `debug` | Gin mode (debug/release) |

## Security Features

- **JWT Authentication**: Secure token-based authentication
- **Input Validation**: Comprehensive request validation
- **CORS Support**: Configurable cross-origin resource sharing
- **SQL Injection Prevention**: Parameterized queries
- **Non-root Container**: Docker container runs as non-root user

## Development

### Essential Commands

```bash
make setup-dev      # Complete development setup
make dev            # Run in development mode
make test           # Run tests
make build          # Build application
make clean          # Clean artifacts
```

### Adding New Endpoints

1. Define models in `internal/MODULE/model.go`
2. Add repository methods in `internal/MODULE/repo.go`
3. Implement business logic in `internal/MODULE/service.go`
4. Create handlers in `internal/MODULE/handler.go`
5. Add routes in `cmd/server/main.go`
6. Update Swagger annotations
7. Write tests

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
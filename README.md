# Battleship Game Engine

A REST API backend service for the classic Battleship game, built with Go, Gin, and GORM.

## Overview

This is a **game engine only** - a backend API service that handles:
- Game state management
- Ship placement and tracking
- Shot processing and hit detection
- Turn alternation for multiplayer games
- Game statistics and persistence

**No frontend UI included** - this service is designed to be consumed by web, mobile, or desktop clients.

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT tokens
- **Monitoring**: Prometheus metrics, OpenTelemetry tracing, Zap logging
- **Containerization**: Docker & Docker Compose

## Features

- Single-player and two-player modes
- Configurable board sizes (5x5 to 100x100)
- Ship placement with rejection sampling
- Shot validation and duplicate detection
- Turn enforcement for multiplayer games
- Game statistics and metrics
- Rate limiting and security headers
- Health check endpoints

## API Documentation

Once the server is running, access:
- API documentation: `http://localhost:8080/api/v1`
- Health check (liveness): `http://localhost:8080/health/live`
- Health check (readiness): `http://localhost:8080/health/ready`

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Running with Docker

```bash
# Start the service and database
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

### Running Locally

```bash
# Copy environment file
cp config/.env.example config/.env

# Edit config/.env with your settings

# Install dependencies
go mod download

# Run migrations (if implemented)
go run ./cmd/migrate up

# Start the server
go run ./cmd/main.go
```

## Project Structure

```
src/
├── models/              # Core business entities (database-independent)
├── services/            # Business logic and use cases
├── adapters/            # Framework glue (Gin, GORM)
│   ├── gin/            # HTTP API layer
│   └── db/             # Database layer
└── validation/         # Input validation

lib/
├── logger/             # Logging infrastructure (Zap)
├── tracing/            # OpenTelemetry tracing
└── metrics/            # Prometheus metrics

tests/
├── unit/               # Unit tests
├── integration/        # Integration tests
└── contract/           # API contract tests
```

## API Endpoints

### Public
- `GET /health/live` - Liveness probe
- `GET /health/ready` - Readiness probe

### Protected (requires JWT)
- `POST /api/v1/games` - Start a new game
- `GET /api/v1/games/:game_id` - Get game state
- `POST /api/v1/games/:game_id/shoot` - Fire a shot
- `GET /api/v1/games/:game_id/stats` - Get game statistics

## Configuration

See `config/.env.example` for all configuration options:

```env
SERVER_PORT=8080
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=battleship
DATABASE_PASSWORD=battleship
DATABASE_NAME=battleship
JWT_SECRET=your-secret-key-change-in-production
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

```bash
golangci-lint run ./...
```

### Building

```bash
go build -o battleship-game-engine ./cmd/main.go
```

## Performance Targets

- API responses: <200ms for standard operations
- p95 latency: <100ms for typical operations
- Support: 1000+ sequential games per minute

## License

MIT

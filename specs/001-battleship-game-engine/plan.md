# Implementation Plan: Battleship Game Engine

**Branch**: `001-battleship-game-engine` | **Date**: 2026-06-19 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-battleship-game-engine/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Build a Battleship Game Engine (BGE) that implements the classic turn-based board game as a package separated from any presentation layer. The engine provides a clean API for game operations (startGame, shoot, gameStats) and supports both single-player and two-player modes with configurable board sizes. This is a Type 3-Advanced challenge requiring clean architecture with separation of concerns between the game engine and presentation layer.

## Technical Context

**Language/Version**: Go 1.21+

**Primary Dependencies**: 
- Gin (web framework for HTTP API)
- GORM (ORM for database operations)
- PostgreSQL (database)
- Docker & Docker Compose (infrastructure)

**Storage**: PostgreSQL (game state persistence via GORM)

**Testing**: Go testing package with table-driven tests, coverage targets: 95%+ for core logic, 80%+ for adapters

**Target Platform**: Linux server (Docker containerized)

**Project Type**: Web service / API service

**Performance Goals**: 
- API responses within 200ms for standard operations
- Support 1000+ sequential games per minute through API
- 4 players per match without degradation

**Constraints**: 
- Board minimum size: 5×5
- Board maximum size: 100×100
- Ship placement retry limit: 100 attempts
- 100ms p95 for typical operations

**Scale/Scope**: 
- Single instance supporting multiple concurrent games
- 1-4 players per game
- Per-session game history tracking

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| Clean Architecture | ✅ PASS | Go structure naturally supports clean architecture with clear layer separation |
| Code Quality | ✅ PASS | Go enforces clean code through tooling and conventions |
| Testing Standards | ✅ PASS | Go testing package supports unit, integration, and E2E tests |
| Test Coverage | ✅ PASS | 95%+ achievable for core game logic |
| User Experience Consistency | ✅ PASS | HTTP API provides consistent interface for UI consumers |
| Performance Requirements | ✅ PASS | Go performance characteristics meet requirements |
| Security & Data Validation | ✅ PASS | Full security design implemented (validation, auth, rate limiting, access control, error sanitization) |
| Observability | ✅ PASS | Logging, tracing, and monitoring design complete in Phase 1 research |

**Gate Status**: PASS - Security requirements addressed in Phase 1 design. Observability will be addressed in Phase 1 design phase.

## Security & Data Validation Design

**Decision Date**: 2026-06-19 | **Approach**: Gin built-in features + community middleware + custom middleware

### 1. Input Validation (Gin Built-in + Custom Package)

**Decision**: Use Gin's `ShouldBindJSON()` with `go-playground/validator/v10` struct tags for API-level validation, plus a centralized `validation/` package for use case-level validation.

**Implementation**:
```go
// validation/api.go
type ShootRequest struct {
    GameID   string `json:"game_id" binding:"required,uuid4"`
    PlayerID string `json:"player_id" binding:"required,uuid4"`
    Row      int    `json:"row" binding:"required,min=0,max=99"`
    Column   int    `json:"column" binding:"required,min=0,max=99"`
}
```

**Why**: Gin's built-in validator provides automatic validation with descriptive error messages. Struct tags ensure consistency across API and use cases. Aligns with FR-045-048 (input validation).

### 2. Authentication (JWT)

**Decision**: Use `gin-gonic/jwt` middleware for full user authentication with JWT tokens.

**Implementation**:
```go
// adapters/gin/middleware/jwt.go
jwtMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
    Realm:       "battleship-game-engine",
    Key:         []byte(secretKey),
    Timeout:     time.Hour * 24,
    MaxRefresh:  time.Hour * 24,
    IdentityKey: "user_id",
    Authenticator: validateUser,
    Authorizator:  checkGameAccess,
    TokenLookup:   "header: Authorization: Bearer",
})
```

**Why**: Full user authentication with JWT tokens provides stateless auth suitable for distributed systems. Custom `Authorizator` enables game-specific access control. Aligns with FR-031-032 (player authorization).

### 3. Rate Limiting (`ulule/limiter`)

**Decision**: Use `ulule/limiter` community middleware for rate limiting.

**Implementation**:
```go
// adapters/gin/middleware/ratelimit.go
limiterConf := limiter.Config{
    Store: limiterStore,
    Rule: limiter.Rate{
        Period: 60 * time.Second,
        Limit:  100, // 100 requests per minute
    },
}
```

**Why**: `ulule/limiter` is the most mature rate-limiting library for Go, supporting multiple storage backends (memory, Redis). Configurable limits per endpoint/user. Aligns with security requirements.

### 4. Access Control (Custom Middleware)

**Decision**: Implement custom middleware for game-specific authorization (player access, turn enforcement).

**Implementation**:
```go
// adapters/gin/middleware/authorization.go
func PlayerAuthorization() gin.HandlerFunc {
    return func(c *gin.Context) {
        gameID := c.Param("game_id")
        playerID := c.GetString("user_id")
        
        // Verify player is authorized for this game
        // Check if it's this player's turn
        // Reject unauthorized access
    }
}
```

**Why**: Custom middleware allows game-specific logic (turn enforcement, player access control). Works with JWT middleware for authenticated requests. Aligns with FR-031-032, FR-036 (turn enforcement).

### 5. Security Headers (Manual Configuration)

**Decision**: Implement custom middleware manually setting security headers (per Gin documentation).

**Implementation**:
```go
// adapters/gin/middleware/securityheaders.go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("Referrer-Policy", "strict-origin")
        c.Next()
    }
}
```

**Why**: Manual configuration gives full control over security headers. Prevents info leakage through headers. Aligns with FR-049 (error sanitization).

### 6. Error Sanitization (Custom Middleware)

**Decision**: Wrap Gin's recovery middleware with custom error sanitization.

**Implementation**:
```go
// adapters/gin/middleware/errorhandler.go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            sanitizedMsg := sanitizeErrorMessage(err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "error":   sanitizedMsg,
            })
        }
    }
}
```

**Why**: Wraps Gin's error handling while sanitizing error messages (no stack traces, no file paths). Aligns with FR-049 (error sanitization).

### Complete Router Setup

```go
// adapters/gin/router.go
func NewRouter() *gin.Engine {
    r := gin.New()
    
    // Global middleware
    r.Use(middleware.SecurityHeaders())
    r.Use(middleware.RateLimitMiddleware())
    r.Use(middleware.ErrorHandler())
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    
    // Public routes
    v1 := r.Group("/api/v1")
    {
        v1.POST("/auth/login", authController.Login)
        v1.POST("/auth/register", authController.Register)
    }
    
    // Protected routes (auth required)
    protected := v1.Group("/")
    protected.Use(middleware.JWTAuth())
    {
        games := protected.Group("/games")
        {
            games.POST("", middleware.PlayerAuthorization(), gameController.StartGame)
            games.POST("/:game_id/shoot", middleware.PlayerAuthorization(), gameController.Shoot)
            games.GET("/:game_id", middleware.PlayerAuthorization(), gameController.GetGameState)
        }
    }
    
    return r
}
```

### Constitution Compliance Summary

| Requirement | Implementation | Status |
|-------------|----------------|--------|
| FR-045-048 Input Validation | Gin validator + validation package | ✅ |
| FR-031-032 Player Authorization | JWT + custom access control middleware | ✅ |
| FR-036 Turn Enforcement | Custom access control middleware | ✅ |
| FR-049 Error Sanitization | Custom error handler middleware | ✅ |
| Security & Data Validation | All 6 components implemented | ✅ |

### Dependencies to Add

```go
// go.mod additions
go 1.21+

dependency "github.com/gin-gonic/gin" v1.10.0
dependency "github.com/appleboy/gin-jwt/v2" v2.10.0
dependency "github.com/ulule/limiter/v3" v3.9.2
dependency "github.com/go-playground/validator/v10" v10.17.0
dependency "github.com/go-playground/universal-translator" v0.18.1
```

## Observability Design (Phase 1)

**Decision Date**: 2026-06-19 | **Status**: Complete | **Research**: [research.md](./research.md)

### 1. Logging Standards (Layer-Based)

**Decision**: Use zap logging package with layer-specific log levels:

| Layer | Log Level | Content |
|-------|-----------|---------|
| Entities | ERROR only | Critical failures, never user-facing data |
| Use Cases | INFO/WARN | Business operation start/end, warnings |
| Adapters | DEBUG/INFO | HTTP requests/responses, database queries |
| Main | INFO | Application startup/shutdown |

**Implementation**:
```go
// lib/logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var (
    Logger *zap.Logger
    SugaredLogger *zap.SugaredLogger
)

func Init() {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zap.ISO8601TimeEncoder
    
    Logger, _ = config.Build()
    SugaredLogger = Logger.Sugar()
}
```

**Why**: Zap provides best performance (10x faster than logrus) with structured logging. Layer-based levels prevent noise from stable Entities while providing visibility in Adapters. Aligns with Constitution Section VIII (Observability).

### 2. Correlation IDs for Request Tracing

**Decision**: Use context.Context for correlation ID propagation:

```go
// adapters/gin/middleware/correlationid.go
func CorrelationIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        correlationID := uuid.New().String()
        ctx := context.WithValue(c.Request.Context(), "correlation_id", correlationID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Writer.Header().Set("X-Correlation-ID", correlationID)
        c.Next()
    }
}
```

**Why**: Context.Context is idiomatic Go for request-scoped values. Correlation ID travels through all layers automatically. Aligns with tracing requirements.

### 3. OpenTelemetry for Distributed Tracing

**Decision**: Use OpenTelemetry SDK with Gin instrumentation:

```go
// lib/tracing/tracing.go
package tracing

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/trace"
)

func Init() {
    exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
    tp := trace.NewTracerProvider(trace.WithSyncer(exp))
    otel.SetTracerProvider(tp)
}
```

**Why**: OpenTelemetry is the industry standard for observability. Gin instrumentation automatically traces HTTP requests. Exporters available for Jaeger, Prometheus, Datadog, etc. Aligns with Constitution Section VIII (Observability).

### 4. Health Check Endpoints

**Decision**: Two-tier health check system:

```go
// /health/live - Liveness probe (is process running?)
{
    "status": "healthy",
    "timestamp": "2026-06-19T12:00:00Z"
}

// /health/ready - Readiness probe (is service ready?)
{
    "status": "ready",
    "checks": {
        "database": "healthy",
        "cache": "healthy"
    },
    "timestamp": "2026-06-19T12:00:00Z"
}
```

**Why**: Kubernetes/containers use liveness/readiness probes differently. Liveness checks process health; Readiness checks service readiness including dependencies. Industry standard for containerized services.

### 5. Performance Metrics Collection

**Decision**: Use Prometheus client library for metrics:

```go
// lib/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    GameStarts = promauto.NewCounter(prometheus.CounterOpts{
        Name: "battleship_game_starts_total",
        Help: "Total number of games started",
    })
    
    ShotLatency = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "battleship_shot_latency_seconds",
        Help:    "Latency of shoot operations",
        Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
    })
)
```

**Why**: Prometheus is the standard for metrics collection. Histograms provide latency percentiles (p50, p95, p99). Counters track totals for rate calculations. Aligns with Constitution Section VI (Performance Requirements).

### Complete Middleware Chain

```go
// adapters/gin/router.go
func NewRouter() *gin.Engine {
    r := gin.New()
    
    // Global middleware
    r.Use(middleware.SecurityHeaders())
    r.Use(middleware.RateLimitMiddleware())
    r.Use(middleware.ErrorHandler())
    r.Use(otelgin.Middleware("battleship-game-engine"))  // Tracing
    r.Use(middleware.CorrelationIDMiddleware())           // Correlation IDs
    r.Use(gin.Logger())                                   // Request logging
    r.Use(gin.Recovery())
    
    // Health check routes (no auth required)
    r.GET("/health/live", healthController.Liveness)
    r.GET("/health/ready", healthController.Readiness)
    
    // Public routes
    v1 := r.Group("/api/v1")
    {
        v1.POST("/auth/login", authController.Login)
        v1.POST("/auth/register", authController.Register)
    }
    
    // Protected routes (auth required)
    protected := v1.Group("/")
    protected.Use(middleware.JWTAuth())
    {
        games := protected.Group("/games")
        {
            games.POST("", middleware.PlayerAuthorization(), gameController.StartGame)
            games.POST("/:game_id/shoot", middleware.PlayerAuthorization(), gameController.Shoot)
            games.GET("/:game_id", middleware.PlayerAuthorization(), gameController.GetGameState)
        }
    }
    
    return r
}
```

### Dependencies to Add

```go
// go.mod additions
go 1.21+

dependency "go.uber.org/zap" v1.27.0
dependency "go.opentelemetry.io/otel" v1.26.0
dependency "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" v1.26.0
dependency "go.opentelemetry.io/otel/sdk/trace" v1.26.0
dependency "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin" v0.52.0
dependency "github.com/prometheus/client_golang" v1.19.0
dependency "github.com/google/uuid" v1.4.0
```

### Constitution Compliance Summary

| Requirement | Implementation | Status |
|-------------|----------------|--------|
| Logging Standards | Layer-based zap logging | ✅ |
| Correlation IDs | Context-based propagation | ✅ |
| Distributed Tracing | OpenTelemetry SDK | ✅ |
| Health Checks | Liveness/Readiness endpoints | ✅ |
| Performance Metrics | Prometheus counters/histograms | ✅ |
| Observability | All components implemented | ✅ |

## Project Structure

### Documentation (this feature)

```text
specs/001-battleship-game-engine/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
src/                     # Game engine source code
├── models/              # Entities - core business rules
│   ├── game.go          # Game entity
│   ├── player.go        # Player entity  
│   ├── ship.go          # Ship entity
│   └── board.go         # Board entity
├── services/            # Use cases - application business rules
│   ├── game_service.go  # Game orchestration
│   └── stats_service.go # Statistics service
├── adapters/            # Interface adapters - framework glue
│   ├── api/             # HTTP handlers (Gin)
│   │   ├── routes.go
│   │   ├── handlers.go
│   │   └── middleware.go
│   └── db/              # Database layer (GORM)
│       ├── models.go    # GORM models
│       └── repository.go
└── lib/                 # Shared utilities
    └── validation.go

tests/
├── contract/            # Contract tests for API
├── integration/         # Integration tests (DB + API)
└── unit/                # Unit tests for core logic

docker-compose.yml       # Infrastructure setup (root level)
Dockerfile               # Application container (root level)
```

**Structure Decision**: Single project structure selected - this is a web service API that will be containerized. The clean architecture layers (models, services, adapters) will be organized as directories under `src/`. The HTTP API layer uses Gin framework, and database layer uses GORM with PostgreSQL. Docker Compose will orchestrate the application and database containers for development. Dockerfile and docker-compose.yml are placed at project root for standard Docker conventions.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all constitutional principles are achievable with the proposed architecture.

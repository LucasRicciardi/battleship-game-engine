# Research: Observability Design for Battleship Game Engine

**Date**: 2026-06-19  
**Feature**: Battleship Game Engine - Observability Implementation  
**Status**: Complete

---

## Research Tasks

### 1. Logging Standards for Clean Architecture (Go)

**Unknown**: What logging patterns should each Clean Architecture layer use?

**Decision**: Layer-based logging levels following Go conventions:

| Layer | Log Level | Content |
|-------|-----------|---------|
| Entities | ERROR only | Critical failures, never user-facing data |
| Use Cases | INFO/WARN | Business operation start/end, warnings |
| Adapters | DEBUG/INFO | HTTP requests/responses, database queries |
| Main | INFO | Application startup/shutdown |

**Rationale**: 
- Entities are foundational and stable; ERROR-only logging prevents noise
- Use Cases orchestrate business flows; INFO/WARN provides operational visibility
- Adapters interface with external systems; DEBUG/INFO captures infrastructure details
- Follows Go community standards (logrus/zap with structured logging)

**Alternatives considered**:
- Uniform INFO across all layers → Too noisy for Entities
- DEBUG across all layers → Excessive volume, performance impact

---

### 2. Correlation IDs for Request Tracing

**Unknown**: How to implement correlation IDs across layers?

**Decision**: Use context.Context for correlation ID propagation:

```go
// Generate correlation ID at adapter layer (HTTP handler)
correlationID := uuid.New().String()
ctx := context.WithValue(r.Context(), "correlation_id", correlationID)

// Log with correlation ID in all layers
log.WithContext(ctx).Info("Processing request")
```

**Rationale**:
- Context.Context is the idiomatic Go way to propagate request-scoped values
- Correlation ID travels through all layers automatically
- No need to pass correlation ID explicitly through function signatures

**Alternatives considered**:
- Global logger with correlation ID → Thread-safety issues in concurrent requests
- Explicit parameter in all functions → Clutters function signatures

---

### 3. OpenTelemetry for Distributed Tracing

**Unknown**: Should we use OpenTelemetry for tracing?

**Decision**: Use OpenTelemetry SDK with Gin instrumentation:

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Initialize tracer
exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
tp := trace.NewTracerProvider(trace.WithSyncer(exp))
otel.SetTracerProvider(tp)

// Add to Gin router
r.Use(otelgin.Middleware("battleship-game-engine"))
```

**Rationale**:
- OpenTelemetry is the industry standard for observability
- Gin instrumentation automatically traces HTTP requests
- Exporters available for Jaeger, Prometheus, Datadog, etc.
- Future-proof for distributed systems

**Alternatives considered**:
- Prometheus metrics only → No distributed tracing
- Custom tracing implementation → Reinventing the wheel

---

### 4. Structured Logging Package Selection

**Unknown**: Which logging package to use (logrus, zap, zerolog)?

**Decision**: Use `zap` for production, `logrus` for development:

```go
// Package: lib/logger
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

**Rationale**:
- Zap provides best performance (10x faster than logrus)
- Structured logging built-in
- SugaredLogger provides familiar API for gradual adoption
- Production-ready with configuration options

**Alternatives considered**:
- logrus → Slower, less structured
- zerolog → Slightly faster but less ecosystem support
- standard library `log` → Not structured, limited features

---

### 5. Health Check Endpoints

**Unknown**: What should health check endpoints return?

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

**Rationale**:
- Kubernetes/containers use liveness/readiness probes differently
- Liveness: process is running (simple check)
- Readiness: service can accept traffic (includes dependencies)
- Industry standard for containerized services

**Alternatives considered**:
- Single /health endpoint → Cannot distinguish process vs service health
- Complex health checks in liveness → Could restart healthy but slow services

---

### 6. Performance Metrics Collection

**Unknown**: What metrics should be collected and how?

**Decision**: Use Prometheus client library for metrics:

```go
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

**Rationale**:
- Prometheus is the standard for metrics collection
- Histograms provide latency percentiles (p50, p95, p99)
- Counters track totals for rate calculations
- Gin integration available for automatic HTTP metrics

**Alternatives considered**:
- Custom metrics implementation → More work, less standard
- OpenTelemetry metrics only → Can use, but Prometheus client is simpler

---

### 7. Error Tracking and Alerting

**Unknown**: How to track errors and set up alerting?

**Decision**: Use OpenTelemetry for error tracking with log aggregation:

```go
// Log errors with full context
Logger.Error("Shoot operation failed",
    zap.Error(err),
    zap.String("correlation_id", correlationID),
    zap.String("game_id", gameID),
    zap.Int("row", row),
    zap.Int("column", column),
)

// Export to log aggregation service (Datadog, ELK, etc.)
// Configure log rotation and retention
```

**Rationale**:
- OpenTelemetry unifies logs, metrics, and traces
- Log aggregation services provide alerting and dashboards
- Structured logging enables querying and filtering

**Alternatives considered**:
- Sentry for error tracking → More for frontend/exception tracking
- Custom error tracking → Reinventing the wheel

---

## Design Summary

### Observability Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Adapters (HTTP)                      │
│  ┌───────────────────────────────────────────────────┐  │
│  │  Request → Middleware Chain                       │  │
│  │  ├─ Security Headers                               │  │
│  │  ├─ Rate Limiting                                  │  │
│  │  ├─ JWT Auth                                       │  │
│  │  ├─ OpenTelemetry Tracing                          │  │
│  │  ├─ Correlation ID Injection                       │  │
│  │  ├─ Request Logging (DEBUG)                        │  │
│  │  └─ Error Handler                                  │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                    Use Cases                            │
│  ┌───────────────────────────────────────────────────┐  │
│  │  Business Logic                                   │  │
│  │  ├─ Info logging (operation start/end)            │  │
│  │  ├─ Warn logging (edge cases)                     │  │
│  │  └─ Metrics collection (counters, histograms)     │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                    Entities                             │
│  ┌───────────────────────────────────────────────────┐  │
│  │  Core Business Rules                              │  │
│  │  └─ Error logging (critical failures only)        │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### Middleware Chain

```go
r.Use(middleware.SecurityHeaders())
r.Use(middleware.RateLimitMiddleware())
r.Use(middleware.ErrorHandler())
r.Use(otelgin.Middleware("battleship-game-engine"))  // Tracing
r.Use(middleware.CorrelationIDMiddleware())           // Correlation IDs
r.Use(gin.Logger())                                   // Request logging
r.Use(gin.Recovery())
```

### Logging Configuration

```yaml
# config/logging.yaml
logging:
  level: "info"  # Development: "debug", Production: "info"
  format: "json"  # Development: "console", Production: "json"
  output: "stdout"
  rotation:
    max_size: 100  # MB
    max_age: 30    # days
    max_backups: 10
```

### Tracing Configuration

```yaml
# config/tracing.yaml
tracing:
  enabled: true
  sampler: "traceidratio"
  sampling_rate: 0.1  # 10% of requests
  exporter: "jaeger"  # Or "otlp", "zipkin"
  service_name: "battleship-game-engine"
```

---

## Dependencies to Add

```go
// go.mod additions
go 1.21+

// Logging
dependency "go.uber.org/zap" v1.27.0

// OpenTelemetry
dependency "go.opentelemetry.io/otel" v1.26.0
dependency "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" v1.26.0
dependency "go.opentelemetry.io/otel/sdk/trace" v1.26.0
dependency "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin" v0.52.0

// Prometheus
dependency "github.com/prometheus/client_golang" v1.19.0

// Utilities
dependency "github.com/google/uuid" v1.4.0  // For correlation IDs
```

---

## Implementation Checklist

### Phase 1 (Design) - Complete

- [x] Define logging standards per Clean Architecture layer
- [x] Design correlation ID propagation strategy
- [x] Select OpenTelemetry for distributed tracing
- [x] Choose zap logging package
- [x] Design health check endpoints
- [x] Define metrics to collect
- [x] Design error tracking strategy

### Phase 2 (Implementation) - Pending

- [ ] Create `lib/logger` package with zap configuration
- [ ] Create `lib/tracing` package with OpenTelemetry setup
- [ ] Create `lib/metrics` package with Prometheus counters/histograms
- [ ] Create `adapters/gin/middleware/correlationid.go`
- [ ] Create `adapters/gin/middleware/healthcheck.go`
- [ ] Update `adapters/gin/router.go` with observability middleware
- [ ] Add logging to Use Cases (INFO/WARN)
- [ ] Add logging to Entities (ERROR only)
- [ ] Add metrics collection to Use Cases
- [ ] Create health check endpoints (`/health/live`, `/health/ready`)
- [ ] Configure log rotation
- [ ] Create observability configuration file

---

## Conclusion

The observability design follows industry best practices with:
1. **Structured logging** using zap for performance and structured output
2. **Distributed tracing** using OpenTelemetry for request tracking
3. **Metrics collection** using Prometheus for monitoring
4. **Health checks** for container orchestration
5. **Correlation IDs** for request tracing across layers

This design satisfies all requirements from Constitution Section VIII (Observability) and provides comprehensive observability for production debugging and monitoring.

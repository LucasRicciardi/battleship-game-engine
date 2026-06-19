# Quickstart: Battleship Game Engine Validation

**Date**: 2026-06-19  
**Feature**: Battleship Game Engine - Observability Implementation  
**Status**: Complete

---

## Prerequisites

- Go 1.21+ installed
- PostgreSQL running (or Docker available)
- JWT secret key for testing

---

## Running the Engine

### 1. Start PostgreSQL (with Docker)

```bash
docker run -d \
  --name battleship-db \
  -e POSTGRES_USER=battleship \
  -e POSTGRES_PASSWORD=battleship \
  -e POSTGRES_DB=battleship \
  -p 5432:5432 \
  postgres:15
```

### 2. Run Database Migrations

```bash
cd src/adapters/db
go run migrate.go up
```

### 3. Start the API Server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.

---

## Validation Scenarios

### Scenario 1: Single Player Game (Complete Flow)

**Objective**: Verify a complete single-player game from start to victory.

#### Steps

1. **Start a new game**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/games \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "board_rows": 8,
       "board_columns": 8,
       "num_players": 1
     }'
   ```

   **Expected**: 201 Created with game_id and initial board state.

2. **Take a shot**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "player_id": 1,
       "row": 3,
       "column": 4
     }'
   ```

   **Expected**: 200 OK with hit/miss result and updated board.

3. **Check game state**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/games/{game_id} \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

   **Expected**: 200 OK with current game state.

4. **Check game statistics**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/games/{game_id}/stats \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

   **Expected**: 200 OK with game statistics.

5. **Repeat shots until victory**:
   Continue shooting until `ships_remaining` is 0.

   **Expected**: Final shot returns `ships_remaining: 0` and `winner: 1`.

#### Validation Commands

```bash
# Check health
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready

# Check metrics (if Prometheus enabled)
curl http://localhost:8080/metrics
```

---

### Scenario 2: Two Player Game (Turn Alternation)

**Objective**: Verify two-player game with correct turn alternation.

#### Steps

1. **Start a two-player game**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/games \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "board_rows": 8,
       "board_columns": 8,
       "num_players": 2
     }'
   ```

   **Expected**: 201 Created with `num_players: 2` and `current_player: 1`.

2. **Player 1 takes a shot**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "player_id": 1,
       "row": 2,
       "column": 3
     }'
   ```

   **Expected**: 200 OK with `current_player: 2` in response.

3. **Player 2 takes a shot**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "player_id": 2,
       "row": 5,
       "column": 6
     }'
   ```

   **Expected**: 200 OK with `current_player: 1` in response.

4. **Attempt out-of-turn shot (should fail)**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "player_id": 1,
       "row": 1,
       "column": 1
     }'
   ```

   **Expected**: 403 Forbidden with error "Not your turn".

#### Validation Commands

```bash
# Check game state after each turn
curl -X GET http://localhost:8080/api/v1/games/{game_id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### Scenario 3: Input Validation

**Objective**: Verify all input validation rules are enforced.

#### Steps

1. **Invalid board size (too small)**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/games \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "board_rows": 3,
       "board_columns": 3,
       "num_players": 1
     }'
   ```

   **Expected**: 400 Bad Request with error "Invalid board size".

2. **Invalid coordinates**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "player_id": 1,
       "row": 10,
       "column": 10
     }'
   ```

   **Expected**: 400 Bad Request with error "Invalid coordinates".

3. **Duplicate shot**:
   ```bash
   # First shot
   curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"player_id": 1, "row": 2, "column": 2}'

   # Duplicate shot
   curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"player_id": 1, "row": 2, "column": 2}'
   ```

   **Expected**: Second request returns 409 Conflict with error "Cell already targeted".

4. **Out of turn shot**:
   ```bash
   # In a two-player game, player 2 shoots out of turn
   curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"player_id": 2, "row": 3, "column": 3}'
   ```

   **Expected**: 403 Forbidden with error "Not your turn".

---

### Scenario 4: Observability Validation

**Objective**: Verify logging, tracing, and metrics are working.

#### Steps

1. **Check request logs**:
   ```bash
   # Look for correlation IDs in server logs
   grep "X-Correlation-ID" /var/log/battleship/server.log
   ```

   **Expected**: Each request has a unique correlation ID.

2. **Check metrics**:
   ```bash
   curl http://localhost:8080/metrics | grep battleship
   ```

   **Expected**: Metrics like `battleship_game_starts_total` and `battleship_shot_latency_seconds` are present.

3. **Check tracing**:
   ```bash
   # If using Jaeger or Zipkin UI
   open http://localhost:16686
   ```

   **Expected**: Traces for API requests are visible.

4. **Check health endpoints**:
   ```bash
   curl http://localhost:8080/health/live
   curl http://localhost:8080/health/ready
   ```

   **Expected**: Both return `{"status": "healthy"}` or `{"status": "ready"}`.

---

## Performance Validation

### Scenario 5: Performance Testing

**Objective**: Verify API performance meets requirements.

#### Steps

1. **Test single operation latency**:
   ```bash
   # Time a single shoot operation
   time curl -X POST http://localhost:8080/api/v1/games/{game_id}/shoot \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"player_id": 1, "row": 3, "column": 4}'
   ```

   **Expected**: Response time < 100ms.

2. **Test sequential games**:
   ```bash
   # Run 1000 sequential games
   for i in {1..1000}; do
     curl -X POST http://localhost:8080/api/v1/games \
       -H "Authorization: Bearer YOUR_JWT_TOKEN" \
       -H "Content-Type: application/json" \
       -d '{"board_rows": 8, "board_columns": 8, "num_players": 1}' > /dev/null
   done
   ```

   **Expected**: Complete in ~60 seconds (1000 games/minute).

---

## Debugging

### Common Issues

1. **Database connection failed**:
   ```bash
   # Check PostgreSQL is running
   docker ps | grep battleship-db
   
   # Check connection string
   export DATABASE_URL=postgres://battleship:battleship@localhost:5432/battleship?sslmode=disable
   ```

2. **JWT authentication failed**:
   ```bash
   # Generate a test token (replace with real auth in production)
   export JWT_SECRET=test_secret_key_12345
   ```

3. **Port already in use**:
   ```bash
   # Change port in config or kill existing process
   lsof -i :8080 | grep LISTEN
   ```

---

## Summary

This quickstart validates:
- ✅ Game creation with configurable board size and player count
- ✅ Shot processing with hit/miss detection
- ✅ Turn alternation in multiplayer games
- ✅ Input validation and error handling
- ✅ Observability (logging, tracing, metrics)
- ✅ Performance requirements

Run all scenarios to ensure the Battleship Game Engine is functioning correctly.

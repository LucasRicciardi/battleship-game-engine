# API Contract: Battleship Game Engine

## Overview

This document defines the HTTP API contract for the Battleship Game Engine. All endpoints return JSON responses with consistent structure.

**Base URL**: `/api/v1`

**Content-Type**: `application/json`

**Authentication**: None (game state is self-contained)

**Version**: v1 (current)

**Rate Limiting**: 100 requests per minute per IP address

---

## Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "detail": "Optional additional details"
  }
}
```

### Error Codes
| Code | Description |
|------|-------------|
| INVALID_REQUEST | Request body is invalid or malformed |
| GAME_NOT_FOUND | Game with specified ID does not exist |
| PLAYER_NOT_FOUND | Player with specified ID does not exist |
| INVALID_COORDINATES | Shot coordinates are out of bounds |
| DUPLICATE_SHOT | Shot at this position already exists |
| INVALID_TURN | It is not this player's turn |
| GAME_NOT_ACTIVE | Game is not in active state |
| SHIP_NOT_FOUND | Ship with specified ID does not exist |
| RATE_LIMIT_EXCEEDED | Too many requests from this IP |
| SERVER_ERROR | Internal server error |
| DATABASE_ERROR | Database operation failed |
| GAME_EXPIRED | Game has expired and been deleted |

---

## Endpoints

### 1. Create Game

Creates a new Battleship game.

**Endpoint**: `POST /api/v1/games`

**Request Body**:
```json
{
  "boardRows": 8,
  "boardCols": 8,
  "players": 1,
  "player1Name": "Player 1",
  "player2Name": "Player 2"
}
```

**Request Validation**:
- `boardRows`: integer, 5 ≤ value ≤ 100, required, default: 8
- `boardCols`: integer, 5 ≤ value ≤ 100, required, default: 8
- `players`: integer, 1 or 2, required, default: 1
- `player1Name`: string, 1-100 characters, required if players=1
- `player2Name`: string, 1-100 characters, required if players=2

**Success Response (201 Created)**:
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "boardRows": 8,
    "boardCols": 8,
    "players": [
      {
        "id": "110e8400-e29b-41d4-a716-446655440000",
        "playerNumber": 1,
        "name": "Player 1"
      }
    ],
    "ships": [
      {
        "id": "210e8400-e29b-41d4-a716-446655440000",
        "type": "destroyer",
        "length": 2,
        "orientation": "horizontal",
        "startX": 0,
        "startY": 0,
        "hits": 0,
        "sunk": false
      },
      {
        "id": "220e8400-e29b-41d4-a716-446655440001",
        "type": "cruiser",
        "length": 3,
        "orientation": "vertical",
        "startX": 3,
        "startY": 0,
        "hits": 0,
        "sunk": false
      },
      {
        "id": "230e8400-e29b-41d4-a716-446655440002",
        "type": "battleship",
        "length": 4,
        "orientation": "horizontal",
        "startX": 5,
        "startY": 0,
        "hits": 0,
        "sunk": false
      }
    ],
    "board": [
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "]
    ],
    "currentTurnPlayerId": "110e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request body or validation error
- `500 Internal Server Error`: Ship placement failed after 100 attempts

---

### 2. Get Game State

Retrieves the current state of a game.

**Endpoint**: `GET /api/v1/games/:gameId`

**Path Parameters**:
- `gameId`: UUID string, required

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "boardRows": 8,
    "boardCols": 8,
    "players": [
      {
        "id": "110e8400-e29b-41d4-a716-446655440000",
        "playerNumber": 1,
        "name": "Player 1"
      }
    ],
    "ships": [ ... ],
    "board": [ ... ],
    "currentTurnPlayerId": "110e8400-e29b-41d4-a716-446655440000",
    "status": "active"
  }
}
```

**Error Responses**:
- `404 Not Found`: Game with specified ID does not exist

---

### 3. Fire Shot

Fires a shot at the opponent's board.

**Endpoint**: `POST /api/v1/games/:gameId/shoot`

**Path Parameters**:
- `gameId`: UUID string, required

**Request Body**:
```json
{
  "playerId": "110e8400-e29b-41d4-a716-446655440000",
  "x": 3,
  "y": 5
}
```

**Request Validation**:
- `playerId`: UUID string, required
- `x`: integer, 0 ≤ value ≤ BoardCols-1, required
- `y`: integer, 0 ≤ value ≤ BoardRows-1, required

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "hit": true,
    "shipSunk": "battleship",
    "shipsRemaining": 2,
    "board": [
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", "X", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "]
    ],
    "currentTurnPlayerId": "110e8400-e29b-41d4-a716-446655440001"
  }
}
```

**Success Response (Miss)**:
```json
{
  "success": true,
  "data": {
    "hit": false,
    "shipsRemaining": 3,
    "board": [ ... ],
    "currentTurnPlayerId": "110e8400-e29b-41d4-a716-446655440001"
  }
}
```

**Success Response (Victory)**:
```json
{
  "success": true,
  "data": {
    "hit": true,
    "shipSunk": "cruiser",
    "shipsRemaining": 0,
    "board": [ ... ],
    "currentTurnPlayerId": "110e8400-e29b-41d4-a716-446655440001",
    "winner": "110e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request body or validation error
- `404 Not Found`: Game or player does not exist
- `409 Conflict`: Duplicate shot at this position
- `403 Forbidden`: It is not this player's turn
- `400 Bad Request`: Game is not active

---

### 4. Get Game Statistics

Retrieves statistics for a game.

**Endpoint**: `GET /api/v1/games/:gameId/stats`

**Path Parameters**:
- `gameId`: UUID string, required

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "gameId": "550e8400-e29b-41d4-a716-446655440000",
    "totalTurns": 15,
    "player1Hits": 8,
    "player1Misses": 7,
    "player2Hits": 6,
    "player2Misses": 9,
    "shipsRemaining": 2,
    "currentTurnPlayerId": "110e8400-e29b-41d4-a716-446655440001",
    "gameStatus": "active"
  }
}
```

**Error Responses**:
- `404 Not Found`: Game with specified ID does not exist

---

### 5. Health Check

Returns the health status of the service.

**Endpoint**: `GET /api/v1/health`

**Success Response (200 OK)**:
```json
{
  "status": "healthy",
  "timestamp": "2026-06-19T12:00:00Z"
}
```

---

## Contract Validation

All requests MUST be validated against the following rules:

1. **Request Body**: Must be valid JSON with required fields
2. **UUID Format**: All IDs must be valid UUID v4 format
3. **Coordinate Range**: X and Y must be within board dimensions
4. **Player Turn**: Only current player can fire shots
5. **Duplicate Prevention**: Same (gameId, playerId, x, y) combination cannot be reused
6. **Game State**: Only active games can accept shots

---

## Error Handling Examples

### Invalid Request
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Invalid request body",
    "detail": "boardRows must be between 5 and 100"
  }
}
```

### Game Not Found
```http
HTTP/1.1 404 Not Found
Content-Type: application/json

{
  "success": false,
  "error": {
    "code": "GAME_NOT_FOUND",
    "message": "Game with ID 123 not found"
  }
}
```

### Duplicate Shot
```http
HTTP/1.1 409 Conflict
Content-Type: application/json

{
  "success": false,
  "error": {
    "code": "DUPLICATE_SHOT",
    "message": "Shot at position (3, 5) already exists",
    "detail": "Previous shot was a hit"
  }
}
```

### Invalid Turn
```http
HTTP/1.1 403 Forbidden
Content-Type: application/json

{
  "success": false,
  "error": {
    "code": "INVALID_TURN",
    "message": "It is not Player 1's turn",
    "detail": "Current turn belongs to Player 2"
  }
}
```

## Rate Limiting

All authenticated endpoints are subject to rate limiting:

- **Limit**: 100 requests per minute per IP address
- **Window**: Sliding window (last 60 seconds)
- **Response**: 429 Too Many Requests when exceeded

### Rate Limit Response
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Please try again later.",
    "detail": "Retry after 45 seconds"
  }
}
```

---

## Versioning Strategy

The API uses URL-based versioning:

- **Current Version**: v1 (`/api/v1`)
- **Version Format**: `/api/v{number}`
- **Deprecation**: Versions are deprecated with 90-day notice
- **Backward Compatibility**: Minor versions maintain backward compatibility

### Version Header
Clients can also specify version via header:
```
API-Version: v1
```

---

## Game Expiration and Deletion

### Game Expiration
Games expire after 30 days of inactivity:

- **Inactivity**: No shots fired or game state retrieved
- **Expiration**: Automatic deletion after 30 days
- **Notification**: No notification (games are ephemeral)

### Game Deletion
Games can be explicitly deleted by the creator:

**Endpoint**: `DELETE /api/v1/games/:gameId`

**Authentication**: Required (game owner only)

**Success Response (204 No Content)**

**Error Responses**:
- `404 Not Found`: Game does not exist
- `403 Forbidden`: Not the game owner
- `409 Conflict`: Game is currently active (cannot delete active games)

---

## Concurrent Access Handling

### Race Condition Prevention
To handle concurrent access scenarios:

1. **Optimistic Locking**: Each game has a version number that increments on each update
2. **Conflict Detection**: Concurrent updates are rejected with 409 Conflict
3. **Retry Logic**: Clients should retry failed operations with fresh state

### Concurrent Shot Prevention
If multiple players shoot simultaneously:

1. First shot is processed normally
2. Second shot returns `409 Conflict` with error code `DUPLICATE_SHOT`
3. Client must retrieve fresh game state before retrying

---

## Database Error Handling

### Connection Failures
When database operations fail:

**Error Response (503 Service Unavailable)**:
```json
{
  "success": false,
  "error": {
    "code": "DATABASE_ERROR",
    "message": "Database connection failed",
    "detail": "Please try again later"
  }
}
```

### Retry Strategy
- **Initial Retry**: Wait 100ms, retry once
- **Exponential Backoff**: Double wait time up to 1 second
- **Max Retries**: 3 attempts total
- **Fallback**: Return 503 if all retries fail

---

## Performance Requirements

### Response Time Targets
- **p50**: <50ms for typical operations
- **p95**: <100ms for typical operations
- **p99**: <200ms for typical operations
- **Health Check**: <10ms

### Payload Size Limits
- **Maximum Response Size**: 1MB
- **Board Size Impact**: 100x100 board = 10,000 cells ≈ 10KB
- **Large Board Handling**: Boards >50x50 may require pagination in future versions

---

## Security Requirements

### Input Sanitization
All user inputs are sanitized to prevent injection attacks:

1. **Coordinate Validation**: Strict integer range checks
2. **UUID Validation**: RFC 4122 compliant UUID format
3. **String Sanitization**: Trim whitespace, reject control characters
4. **SQL Injection Prevention**: Parameterized queries only
5. **XSS Prevention**: No HTML rendering in API responses

### Authentication Requirements
- **JWT Tokens**: Required for protected endpoints (future)
- **Token Expiration**: 24 hours
- **Token Refresh**: Available via refresh endpoint (future)

---

## Observability Requirements

### Logging
All API requests are logged with:
- Request ID (correlation ID)
- Timestamp
- Client IP
- Endpoint path
- HTTP method
- Response status code
- Response time

### Tracing
Distributed tracing is enabled for all endpoints:
- Trace ID generated per request
- Span context propagated to database and cache
- Traces exported to observability platform

### Metrics
Key metrics collected:
- Request rate per endpoint
- Response time percentiles
- Error rate per error code
- Active game count
- Average game duration

---

## Error Handling Examples

### Rate Limit Exceeded
```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Please try again later.",
    "detail": "Retry after 45 seconds"
  }
}
```

### Database Error
```http
HTTP/1.1 503 Service Unavailable
Content-Type: application/json

{
  "success": false,
  "error": {
    "code": "DATABASE_ERROR",
    "message": "Database connection failed",
    "detail": "Please try again later"
  }
}
```

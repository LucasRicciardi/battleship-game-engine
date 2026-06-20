# Battleship Game Engine API Documentation

**Version**: v1  
**Base URL**: `/api/v1`  
**Content-Type**: `application/json`

---

## Overview

This document describes the REST API for the Battleship Game Engine. The API provides endpoints for managing Battleship games, including starting games, firing shots, and retrieving game state.

**Authentication**: JWT tokens required for protected endpoints  
**Rate Limiting**: 100 requests per minute per IP address  
**Versioning**: URL-based versioning (`/api/v1`)

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
  "error": "Error message",
  "error_code": "ERROR_CODE"
}
```

### Error Codes
| Code | Description | HTTP Status |
|------|-------------|-------------|
| INVALID_REQUEST | Request body is invalid or malformed | 400 |
| GAME_NOT_FOUND | Game with specified ID does not exist | 404 |
| PLAYER_NOT_FOUND | Player with specified ID does not exist | 404 |
| INVALID_COORDINATES | Shot coordinates are out of bounds | 400 |
| DUPLICATE_SHOT | Shot at this position already exists | 409 |
| INVALID_TURN | It is not this player's turn | 403 |
| GAME_NOT_ACTIVE | Game is not in active state | 400 |
| RATE_LIMIT_EXCEEDED | Too many requests from this IP | 429 |
| DATABASE_ERROR | Database operation failed | 503 |

---

## Endpoints

### 1. Health Check (Liveness)

Returns the health status of the service process.

**Endpoint**: `GET /health/live`

**Authentication**: Not required

**Rate Limiting**: Not applicable

**Success Response (200 OK)**:
```json
{
  "status": "healthy",
  "timestamp": "2026-06-19T12:00:00Z"
}
```

---

### 2. Health Check (Readiness)

Returns the health status of the service including dependencies.

**Endpoint**: `GET /health/ready`

**Authentication**: Not required

**Rate Limiting**: Not applicable

**Success Response (200 OK)**:
```json
{
  "status": "ready",
  "checks": {
    "database": "healthy",
    "cache": "healthy"
  },
  "timestamp": "2026-06-19T12:00:00Z"
}
```

**Failure Response (503 Service Unavailable)**:
```json
{
  "status": "not_ready",
  "checks": {
    "database": "unhealthy",
    "cache": "healthy"
  },
  "timestamp": "2026-06-19T12:00:00Z"
}
```

---

### 3. Start Game

Creates a new Battleship game with the specified configuration.

**Endpoint**: `POST /api/v1/games`

**Authentication**: Required (JWT token in `Authorization` header)

**Request Body**:
```json
{
  "board_rows": 8,
  "board_columns": 8,
  "num_players": 1
}
```

**Request Validation**:
- `board_rows`: integer, 5 ≤ value ≤ 100, required, default: 8
- `board_columns`: integer, 5 ≤ value ≤ 100, required, default: 8
- `num_players`: integer, 1 or 2, required, default: 1

**Success Response (201 Created)**:
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "board_rows": 8,
    "board_columns": 8,
    "num_players": 1,
    "turn": 1,
    "current_player": 1,
    "status": "active",
    "ships": [
      {
        "id": "111e8400-e29b-41d4-a716-446655440000",
        "type": "destroyer",
        "length": 2,
        "positions": [
          {"row": 0, "column": 0},
          {"row": 0, "column": 1}
        ],
        "hits": 0,
        "sunk": false
      },
      {
        "id": "222e8400-e29b-41d4-a716-446655440000",
        "type": "cruiser",
        "length": 3,
        "positions": [
          {"row": 2, "column": 2},
          {"row": 2, "column": 3},
          {"row": 2, "column": 4}
        ],
        "hits": 0,
        "sunk": false
      },
      {
        "id": "333e8400-e29b-41d4-a716-446655440000",
        "type": "battleship",
        "length": 4,
        "positions": [
          {"row": 5, "column": 0},
          {"row": 5, "column": 1},
          {"row": 5, "column": 2},
          {"row": 5, "column": 3}
        ],
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
    "players": [
      {
        "id": "110e8400-e29b-41d4-a716-446655440000",
        "player_number": 1,
        "name": "Player 1"
      }
    ]
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request body or validation error
- `401 Unauthorized`: Missing or invalid JWT token
- `409 Conflict`: Failed to place ships after 100 attempts
- `429 Too Many Requests`: Rate limit exceeded

---

### 4. Get Game State

Retrieves the current state of a game.

**Endpoint**: `GET /api/v1/games/:game_id`

**Authentication**: Required (JWT token in `Authorization` header)

**Path Parameters**:
- `game_id`: string (UUID) - Required

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "board_rows": 8,
    "board_columns": 8,
    "num_players": 2,
    "turn": 5,
    "current_player": 2,
    "status": "active",
    "winner": 0,
    "ships": [
      {
        "id": "111e8400-e29b-41d4-a716-446655440000",
        "type": "destroyer",
        "length": 2,
        "positions": [
          {"row": 0, "column": 0},
          {"row": 0, "column": 1}
        ],
        "hits": 0,
        "sunk": false
      },
      {
        "id": "222e8400-e29b-41d4-a716-446655440000",
        "type": "cruiser",
        "length": 3,
        "positions": [
          {"row": 2, "column": 2},
          {"row": 2, "column": 3},
          {"row": 2, "column": 4}
        ],
        "hits": 0,
        "sunk": false
      },
      {
        "id": "333e8400-e29b-41d4-a716-446655440000",
        "type": "battleship",
        "length": 4,
        "positions": [
          {"row": 5, "column": 0},
          {"row": 5, "column": 1},
          {"row": 5, "column": 2},
          {"row": 5, "column": 3}
        ],
        "hits": 4,
        "sunk": true
      }
    ],
    "board": [
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", "O", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      ["X", "X", "X", "X", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "]
    ],
    "players": [
      {
        "id": "110e8400-e29b-41d4-a716-446655440000",
        "player_number": 1,
        "name": "Player 1",
        "shots_fired": 10,
        "shots_hit": 3,
        "ships_sunk": 1
      },
      {
        "id": "220e8400-e29b-41d4-a716-446655440001",
        "player_number": 2,
        "name": "Player 2",
        "shots_fired": 9,
        "shots_hit": 2,
        "ships_sunk": 0
      }
    ]
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid game_id format
- `401 Unauthorized`: Missing or invalid JWT token
- `404 Not Found`: Game with specified ID does not exist
- `429 Too Many Requests`: Rate limit exceeded

---

### 5. Fire Shot

Fires a shot at the specified board coordinates.

**Endpoint**: `POST /api/v1/games/:game_id/shoot`

**Authentication**: Required (JWT token in `Authorization` header)

**Path Parameters**:
- `game_id`: string (UUID) - Required

**Request Body**:
```json
{
  "player_id": "110e8400-e29b-41d4-a716-446655440000",
  "row": 3,
  "column": 4
}
```

**Request Validation**:
- `player_id`: string (UUID) - Required
- `row`: integer, 0 ≤ value ≤ BoardRows-1, required
- `column`: integer, 0 ≤ value ≤ BoardColumns-1, required

**Success Response (200 OK) - Hit**:
```json
{
  "success": true,
  "data": {
    "hit": true,
    "ship_sunk": "battleship",
    "ships_remaining": 2,
    "turn": 5,
    "current_player": 2,
    "board": [
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", "O", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      ["X", "X", "X", "X", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "]
    ],
    "ships": [
      {
        "id": "111e8400-e29b-41d4-a716-446655440000",
        "type": "destroyer",
        "length": 2,
        "positions": [
          {"row": 0, "column": 0},
          {"row": 0, "column": 1}
        ],
        "hits": 0,
        "sunk": false
      },
      {
        "id": "222e8400-e29b-41d4-a716-446655440000",
        "type": "cruiser",
        "length": 3,
        "positions": [
          {"row": 2, "column": 2},
          {"row": 2, "column": 3},
          {"row": 2, "column": 4}
        ],
        "hits": 0,
        "sunk": false
      },
      {
        "id": "333e8400-e29b-41d4-a716-446655440000",
        "type": "battleship",
        "length": 4,
        "positions": [
          {"row": 5, "column": 0},
          {"row": 5, "column": 1},
          {"row": 5, "column": 2},
          {"row": 5, "column": 3}
        ],
        "hits": 4,
        "sunk": true
      }
    ]
  }
}
```

**Success Response (200 OK) - Miss**:
```json
{
  "success": true,
  "data": {
    "hit": false,
    "ships_remaining": 3,
    "turn": 5,
    "current_player": 2,
    "board": [
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", "O", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "],
      [" ", " ", " ", " ", " ", " ", " ", " "]
    ],
    "ships": [
      {
        "id": "111e8400-e29b-41d4-a716-446655440000",
        "type": "destroyer",
        "length": 2,
        "positions": [
          {"row": 0, "column": 0},
          {"row": 0, "column": 1}
        ],
        "hits": 0,
        "sunk": false
      },
      {
        "id": "222e8400-e29b-41d4-a716-446655440000",
        "type": "cruiser",
        "length": 3,
        "positions": [
          {"row": 2, "column": 2},
          {"row": 2, "column": 3},
          {"row": 2, "column": 4}
        ],
        "hits": 0,
        "sunk": false
      },
      {
        "id": "333e8400-e29b-41d4-a716-446655440000",
        "type": "battleship",
        "length": 4,
        "positions": [
          {"row": 5, "column": 0},
          {"row": 5, "column": 1},
          {"row": 5, "column": 2},
          {"row": 5, "column": 3}
        ],
        "hits": 0,
        "sunk": false
      }
    ]
  }
}
```

**Success Response (200 OK) - Victory**:
```json
{
  "success": true,
  "data": {
    "hit": true,
    "ship_sunk": "cruiser",
    "ships_remaining": 0,
    "turn": 5,
    "current_player": 2,
    "board": [ ... ],
    "winner": "110e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request body or validation error
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: It is not this player's turn
- `404 Not Found`: Game or player does not exist
- `409 Conflict`: Duplicate shot at this position
- `429 Too Many Requests`: Rate limit exceeded

---

### 6. Get Game Statistics

Retrieves statistics for a game.

**Endpoint**: `GET /api/v1/games/:game_id/stats`

**Authentication**: Required (JWT token in `Authorization` header)

**Path Parameters**:
- `game_id`: string (UUID) - Required

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "game_id": "550e8400-e29b-41d4-a716-446655440000",
    "total_turns": 15,
    "total_shots": 19,
    "total_hits": 5,
    "total_misses": 14,
    "ships_remaining": 2,
    "current_player": 2,
    "status": "active",
    "winner": 0,
    "player_stats": [
      {
        "player_id": "110e8400-e29b-41d4-a716-446655440000",
        "shots_fired": 10,
        "shots_hit": 3,
        "ships_sunk": 1
      },
      {
        "player_id": "220e8400-e29b-41d4-a716-446655440001",
        "shots_fired": 9,
        "shots_hit": 2,
        "ships_sunk": 0
      }
    ]
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid game_id format
- `401 Unauthorized`: Missing or invalid JWT token
- `404 Not Found`: Game with specified ID does not exist
- `429 Too Many Requests`: Rate limit exceeded

---

## Board Representation

The board is represented as a 2D array of strings:

| Value | Meaning |
|-------|---------|
| `" "` | Untargeted cell |
| `"O"` | Missed shot |
| `"X"` | Hit shot |
| `"S"` | Ship position (visible only to owner) |

Example 8x8 board:
```json
[
  [" ", " ", " ", " ", " ", " ", " ", " "],
  [" ", " ", " ", " ", " ", " ", " ", " "],
  [" ", " ", " ", " ", " ", " ", " ", " "],
  [" ", " ", " ", " ", "O", " ", " ", " "],
  [" ", " ", " ", " ", " ", " ", " ", " "],
  ["X", "X", "X", "X", " ", " ", " ", " "],
  [" ", " ", " ", " ", " ", " ", " ", " "],
  [" ", " ", " ", " ", " ", " ", " ", " "]
]
```

---

## Ship Types

| Type | Length | Count |
|------|--------|-------|
| destroyer | 2 | 1 |
| cruiser | 3 | 1 |
| battleship | 4 | 1 |
| aircraft | 5 | 1 |

---

## Performance Targets

- **p50**: <50ms for typical operations
- **p95**: <100ms for typical operations
- **p99**: <200ms for typical operations
- **Health Check**: <10ms

---

## Rate Limiting

All authenticated endpoints are subject to rate limiting:

- **Limit**: 100 requests per minute per IP address
- **Window**: Sliding window (last 60 seconds)
- **Response**: 429 Too Many Requests when exceeded

---

## Error Handling Examples

### Invalid Request
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "success": false,
  "error": "Invalid request body",
  "error_code": "INVALID_REQUEST"
}
```

### Game Not Found
```http
HTTP/1.1 404 Not Found
Content-Type: application/json

{
  "success": false,
  "error": "Game not found",
  "error_code": "GAME_NOT_FOUND"
}
```

### Duplicate Shot
```http
HTTP/1.1 409 Conflict
Content-Type: application/json

{
  "success": false,
  "error": "Cell already targeted",
  "error_code": "DUPLICATE_SHOT"
}
```

### Invalid Turn
```http
HTTP/1.1 403 Forbidden
Content-Type: application/json

{
  "success": false,
  "error": "Not this player's turn",
  "error_code": "INVALID_TURN"
}
```

### Rate Limit Exceeded
```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
  "success": false,
  "error": "Rate limit exceeded",
  "error_code": "RATE_LIMIT_EXCEEDED"
}
```

---

## Authentication

All protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

---

## Versioning

The API uses URL-based versioning:

- **Current Version**: v1 (`/api/v1`)
- **Version Format**: `/api/v{number}`
- **Deprecation**: Versions are deprecated with 90-day notice

---

## Support

For issues and questions, please open an issue on the project repository.

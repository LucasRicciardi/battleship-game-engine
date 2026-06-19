# API Contracts: Battleship Game Engine

**Date**: 2026-06-19  
**Feature**: Battleship Game Engine - Observability Implementation  
**Status**: Complete

---

## Overview

This document defines the HTTP API contracts for the Battleship Game Engine. The API follows REST conventions and uses JSON for request/response bodies.

**Base URL**: `/api/v1`

**Content-Type**: `application/json`

**Authentication**: JWT tokens required for protected endpoints

---

## API Response Format

All API responses follow this consistent format:

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

### Error Response
```json
{
  "success": false,
  "data": null,
  "error": "Descriptive error message"
}
```

### Error Response Codes
| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input or validation error |
| 401 | Unauthorized - Missing or invalid JWT token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource does not exist |
| 409 | Conflict - Duplicate operation (e.g., duplicate shot) |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Unexpected error |

---

## Endpoints

### 1. Start Game

Creates a new Battleship game with the specified configuration.

**Endpoint**: `POST /games`

**Authentication**: Required

**Request Body**:
```json
{
  "board_rows": 8,
  "board_columns": 8,
  "num_players": 1
}
```

**Request Validation**:
- `board_rows`: integer, min: 5, max: 100, required
- `board_columns`: integer, min: 5, max: 100, required
- `num_players`: integer, min: 1, max: 2, required

**Success Response (201 Created)**:
```json
{
  "success": true,
  "data": {
    "game_id": "550e8400-e29b-41d4-a716-446655440000",
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
    ]
  },
  "error": null
}
```

**Error Response (400 Bad Request)**:
```json
{
  "success": false,
  "data": null,
  "error": "Invalid board size: minimum is 5x5"
}
```

**Error Response (409 Conflict)**:
```json
{
  "success": false,
  "data": null,
  "error": "Failed to place ships after 100 attempts"
}
```

---

### 2. Shoot

Fires a shot at the specified board coordinates.

**Endpoint**: `POST /games/:game_id/shoot`

**Authentication**: Required

**Path Parameters**:
- `game_id`: string (UUID) - Required

**Request Body**:
```json
{
  "player_id": 1,
  "row": 3,
  "column": 4
}
```

**Request Validation**:
- `player_id`: integer, min: 1, max: NumPlayers, required
- `row`: integer, min: 0, max: BoardRows-1, required
- `column`: integer, min: 0, max: BoardColumns-1, required

**Success Response (200 OK)**:
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
  },
  "error": null
}
```

**Success Response (Miss)**:
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
  },
  "error": null
}
```

**Error Response (400 Bad Request)**:
```json
{
  "success": false,
  "data": null,
  "error": "Invalid coordinates: row must be between 0 and 7"
}
```

**Error Response (409 Conflict)**:
```json
{
  "success": false,
  "data": null,
  "error": "Cell already targeted: previous shot was a miss"
}
```

**Error Response (403 Forbidden)**:
```json
{
  "success": false,
  "data": null,
  "error": "Not your turn: it is player 2's turn"
}
```

---

### 3. Get Game State

Retrieves the current state of a game.

**Endpoint**: `GET /games/:game_id`

**Authentication**: Required

**Path Parameters**:
- `game_id`: string (UUID) - Required

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "game_id": "550e8400-e29b-41d4-a716-446655440000",
    "board_rows": 8,
    "board_columns": 8,
    "num_players": 2,
    "turn": 5,
    "current_player": 2,
    "status": "active",
    "winner": 0,
    "players": [
      {
        "id": 1,
        "shots_fired": 10,
        "shots_hit": 3,
        "ships_sunk": 1
      },
      {
        "id": 2,
        "shots_fired": 9,
        "shots_hit": 2,
        "ships_sunk": 0
      }
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
    ]
  },
  "error": null
}
```

**Error Response (404 Not Found)**:
```json
{
  "success": false,
  "data": null,
  "error": "Game not found"
}
```

---

### 4. Get Game Statistics

Retrieves statistics for a game.

**Endpoint**: `GET /games/:game_id/stats`

**Authentication**: Required

**Path Parameters**:
- `game_id`: string (UUID) - Required

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "game_id": "550e8400-e29b-41d4-a716-446655440000",
    "num_players": 2,
    "turn": 5,
    "active_player": 2,
    "status": "active",
    "total_shots": 19,
    "total_hits": 5,
    "total_misses": 14,
    "ships_remaining": 5,
    "winner": 0,
    "player_stats": [
      {
        "player_id": 1,
        "shots_fired": 10,
        "shots_hit": 3,
        "ships_sunk": 1
      },
      {
        "player_id": 2,
        "shots_fired": 9,
        "shots_hit": 2,
        "ships_sunk": 0
      }
    ]
  },
  "error": null
}
```

**Error Response (404 Not Found)**:
```json
{
  "success": false,
  "data": null,
  "error": "Game not found"
}
```

---

## Health Check Endpoints

### Liveness Probe

Checks if the process is running.

**Endpoint**: `GET /health/live`

**Authentication**: Not required

**Success Response (200 OK)**:
```json
{
  "status": "healthy",
  "timestamp": "2026-06-19T12:00:00Z"
}
```

### Readiness Probe

Checks if the service is ready to accept traffic.

**Endpoint**: `GET /health/ready`

**Authentication**: Not required

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

## Error Codes

| Error Code | Description | HTTP Status |
|------------|-------------|-------------|
| INVALID_BOARD_SIZE | Board dimensions outside valid range | 400 |
| INVALID_COORDINATES | Shot coordinates outside board | 400 |
| INVALID_PLAYER_ID | Player ID outside valid range | 400 |
| GAME_NOT_FOUND | Game ID does not exist | 404 |
| GAME_COMPLETED | Operation not allowed on completed game | 409 |
| DUPLICATE_SHOT | Shot at already-targeted cell | 409 |
| OUT_OF_TURN | Shot by wrong player in multiplayer | 403 |
| SHIP_PLACEMENT_FAILED | Could not place ships within retry limit | 409 |
| INTERNAL_ERROR | Unexpected server error | 500 |

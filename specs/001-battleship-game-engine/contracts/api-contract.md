# API Contract: Battleship Game Engine

## Overview

This document defines the HTTP API contract for the Battleship Game Engine. All endpoints return JSON responses with consistent structure.

**Base URL**: `/api/v1`

**Content-Type**: `application/json`

**Authentication**: None (game state is self-contained)

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
| SERVER_ERROR | Internal server error |

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

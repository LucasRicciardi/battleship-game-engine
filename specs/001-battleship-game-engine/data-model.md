# Data Model: Battleship Game Engine

**Date**: 2026-06-19  
**Feature**: Battleship Game Engine - Observability Implementation  
**Status**: Complete

---

## Core Entities

### Game

Represents an active Battleship game session containing board state, ship positions, hit/miss tracking, turn information, and player data.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| ID | string (UUID) | Unique game identifier | Required, unique |
| BoardRows | int | Number of rows in the board | Min: 5, Max: 100 |
| BoardColumns | int | Number of columns in the board | Min: 5, Max: 100 |
| NumPlayers | int | Number of players (1 or 2) | Min: 1, Max: 2 |
| CreatedAt | time.Time | Game creation timestamp | Required |
| Turn | int | Current turn number (1-indexed) | Min: 1 |
| CurrentPlayer | int | Player whose turn it is (1-indexed) | Min: 1, Max: NumPlayers |
| Status | string | Game status: "active", "completed" | Required |
| Winner | int | Winning player number (0 if no winner) | Min: 0, Max: NumPlayers |

**State Transitions**:
```
created → active → completed (with winner or draw)
```

---

### Board

An N×M grid representing the game area where ships are placed and shots are fired.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| GameID | string (UUID) | Associated game ID | Required |
| Rows | int | Number of rows | Same as Game.BoardRows |
| Columns | int | Number of columns | Same as Game.BoardColumns |
| Cells | [][]string | 2D array of cell states | Size: Rows × Columns |

**Cell States**:
- `" "` (space) - Untargeted cell
- `"O"` - Miss (no ship at location)
- `"X"` - Hit (ship occupied)

**Validation Rules**:
- Cells array must be exactly Rows × Columns
- All cells must be one of the three valid states
- No cell can be modified directly (only through shoot operation)

---

### Ship

A vessel with properties `{ id: string, type: string, length: number, positions: {row,col}[], hits: number, sunk: boolean }` placed horizontally or vertically on the board.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| ID | string (UUID) | Unique ship identifier | Required, unique |
| GameID | string (UUID) | Associated game ID | Required |
| Type | string | Ship type: "destroyer", "cruiser", "battleship" | Required |
| Length | int | Number of cells the ship occupies | Min: 1, Max: min(Rows, Columns) |
| Positions | []Position | Array of cell coordinates | Length must equal Length field |
| Hits | int | Number of hits on this ship | Min: 0, Max: Length |
| Sunk | bool | Whether the ship is sunk | Derived from Hits == Length |

**Position Structure**:
| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Row | int | Row coordinate (0-indexed) | Min: 0, Max: Rows-1 |
| Column | int | Column coordinate (0-indexed) | Min: 0, Max: Columns-1 |

**Validation Rules**:
- Positions array length must equal Length field
- All positions must be within board boundaries
- Ship cannot overlap with other ships (enforced at placement time)
- Ship must be placed horizontally OR vertically (not diagonal)

**State Transitions**:
```
placed (hits=0, sunk=false) → damaged (hits>0, sunk=false) → sunk (hits=Length, sunk=true)
```

---

### Player

A game participant with independent board state and ship positions.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| ID | int | Player number (1-indexed) | Min: 1, Max: NumPlayers |
| GameID | string (UUID) | Associated game ID | Required |
| BoardID | string (UUID) | Associated board ID | Required |
| Ships | []Ship | Player's ships | Array of Ship entities |
| ShotsFired | int | Total shots fired by this player | Min: 0 |
| ShotsHit | int | Shots that hit a ship | Min: 0, Max: ShotsFired |
| ShipsSunk | int | Number of opponent ships sunk | Min: 0, Max: NumOpponentShips |

**Validation Rules**:
- Each player must have exactly one board per game
- Ship count must match game configuration (3 ships by default)
- ShotsHit cannot exceed ShotsFired
- ShipsSunk cannot exceed opponent's ship count

---

### GameStats

A data structure containing metrics about current game state.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| GameID | string (UUID) | Associated game ID | Required |
| NumPlayers | int | Number of players | Required |
| Turn | int | Current turn number | Required |
| ActivePlayer | int | Player whose turn it is | Required |
| Status | string | Game status | Required |
| TotalShots | int | Total shots fired across all players | Min: 0 |
| TotalHits | int | Total shots that hit ships | Min: 0, Max: TotalShots |
| TotalMisses | int | Total shots that missed | Min: 0, Max: TotalShots |
| ShipsRemaining | int | Number of unsunk ships across all players | Min: 0, Max: TotalShips |
| Winner | int | Winning player number (0 if no winner) | Min: 0, Max: NumPlayers |

**Derived Values**:
- TotalMisses = TotalShots - TotalHits
- ShipsRemaining = Sum of (1 if ship.sunk == false else 0) for all ships

---

## Database Schema (GORM Models)

### games table

```go
type Game struct {
    ID             string    `gorm:"primaryKey;type:uuid"`
    BoardRows      int       `gorm:"not null;default:8"`
    BoardColumns   int       `gorm:"not null;default:8"`
    NumPlayers     int       `gorm:"not null;default:1"`
    CreatedAt      time.Time `gorm:"not null"`
    Turn           int       `gorm:"not null;default:1"`
    CurrentPlayer  int       `gorm:"not null;default:1"`
    Status         string    `gorm:"not null;default:'active'"`
    Winner         int       `gorm:"default:0"`
}
```

### boards table

```go
type Board struct {
    ID        string `gorm:"primaryKey;type:uuid"`
    GameID    string `gorm:"not null;type:uuid;index"`
    Rows      int    `gorm:"not null"`
    Columns   int    `gorm:"not null"`
    CellsJSON string `gorm:"type:text"` // Serialized 2D array
}
```

### ships table

```go
type Ship struct {
    ID        string `gorm:"primaryKey;type:uuid"`
    GameID    string `gorm:"not null;type:uuid;index"`
    Type      string `gorm:"not null"` // "destroyer", "cruiser", "battleship"
    Length    int    `gorm:"not null"`
    Positions string `gorm:"type:text"` // Serialized []Position
    Hits      int    `gorm:"not null;default:0"`
    Sunk      bool   `gorm:"not null;default:false"`
}
```

### players table

```go
type Player struct {
    ID           int    `gorm:"primaryKey"` // 1-indexed player number
    GameID       string `gorm:"not null;type:uuid;index"`
    BoardID      string `gorm:"not null;type:uuid"`
    ShotsFired   int    `gorm:"not null;default:0"`
    ShotsHit     int    `gorm:"not null;default:0"`
    ShipsSunk    int    `gorm:"not null;default:0"`
}
```

---

## Relationships

```
Game (1) ──< (N) Board
Game (1) ──< (N) Ship
Game (1) ──< (N) Player

Player (1) ──> (1) Board
Ship (N) ──> (1) Game
```

---

## Validation Summary

> **Note**: This data model references validation rules defined in `spec.md` (FR-018 through FR-021, FR-045 through FR-048). The spec serves as the source of truth for all validation requirements.

### Game Creation Validation
- BoardRows: 5-100 (per FR-018, FR-019)
- BoardColumns: 5-100 (per FR-018, FR-019)
- NumPlayers: 1-2 (per FR-029)
- Ship placement must succeed within 100 attempts (per FR-022)

### Ship Placement Validation
- Ship length: 1 to min(Rows, Columns) (per FR-020, FR-021)
- Ship must fit within board boundaries (per FR-021)
- No ship overlap allowed (per SPR-003)
- Horizontal OR vertical orientation only (per FR-003)

### Shot Validation
- Row: 0 to Rows-1 (per FR-046)
- Column: 0 to Columns-1 (per FR-047)
- Cell must not be previously targeted (per FR-016)
- Shot must match current player's turn (multiplayer) (per FR-036)

### Game Stats Validation
- TotalShots = Sum of ShotsFired for all players
- TotalHits + TotalMisses = TotalShots
- ShipsRemaining = TotalShips - Sum of sunk ships

package models

import (
	"fmt"
	"time"
)

// Game represents a Battleship game
type Game struct {
	ID             string    `json:"id"`
	BoardRows      int       `json:"board_rows"`
	BoardColumns   int       `json:"board_columns"`
	NumPlayers     int       `json:"num_players"`
	Turn           int       `json:"turn"`
	CurrentPlayer  int       `json:"current_player"`
	Status         string    `json:"status"`
	Winner         int       `json:"winner,omitempty"`
	Ships          []Ship    `json:"ships"`
	Boards         []Board   `json:"boards,omitempty"`
	Players        []Player  `json:"players,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Player represents a player in the game
type Player struct {
	ID        string `json:"id"`
	Number    int    `json:"player_number"`
	Name      string `json:"name"`
	ShotsFired int   `json:"shots_fired,omitempty"`
	ShotsHit  int    `json:"shots_hit,omitempty"`
	ShipsSunk int    `json:"ships_sunk,omitempty"`
}

// GameStatus represents the status of a game
type GameStatus string

const (
	// Game statuses
	Active   GameStatus = "active"
	Complete GameStatus = "complete"
	Abandoned GameStatus = "abandoned"
)

// NewGame creates a new game with the specified configuration
func NewGame(id string, boardRows, boardColumns, numPlayers int) *Game {
	boards := make([]Board, numPlayers)
	for i := 0; i < numPlayers; i++ {
		boards[i] = *NewBoard(boardRows, boardColumns)
		boards[i].OwnerID = fmt.Sprintf("player-%d", i+1)
	}

	players := make([]Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = Player{
			ID:        fmt.Sprintf("player-%d", i+1),
			Number:    i + 1,
			Name:      fmt.Sprintf("Player %d", i+1),
			ShotsFired: 0,
			ShotsHit:  0,
			ShipsSunk: 0,
		}
	}

	return &Game{
		ID:             id,
		BoardRows:      boardRows,
		BoardColumns:   boardColumns,
		NumPlayers:     numPlayers,
		Turn:           1,
		CurrentPlayer:  1,
		Status:         string(Active),
		Winner:         0,
		Ships:          []Ship{},
		Boards:         boards,
		Players:        players,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// IsPlayerTurn checks if it's the specified player's turn
func (g *Game) IsPlayerTurn(playerNumber int) bool {
	return g.CurrentPlayer == playerNumber
}

// NextTurn advances the game to the next player's turn
func (g *Game) NextTurn() {
	g.Turn++
	g.CurrentPlayer = ((g.CurrentPlayer - 1 + 1) % g.NumPlayers) + 1
	g.UpdatedAt = time.Now()
}

// GetCurrentBoard returns the board for the current player
func (g *Game) GetCurrentBoard() *Board {
	if g.CurrentPlayer > 0 && g.CurrentPlayer <= len(g.Boards) {
		return &g.Boards[g.CurrentPlayer-1]
	}
	return nil
}

// GetPlayerBoard returns the board for the specified player
func (g *Game) GetPlayerBoard(playerNumber int) *Board {
	if playerNumber > 0 && playerNumber <= len(g.Boards) {
		return &g.Boards[playerNumber-1]
	}
	return nil
}

// GetActiveShipsCount returns the number of active ships for a player
func (g *Game) GetActiveShipsCount(playerNumber int) int {
	if board := g.GetPlayerBoard(playerNumber); board != nil {
		return board.GetActiveShipsCount()
	}
	return 0
}

// IsGameOver checks if the game is over (any player has no ships remaining)
func (g *Game) IsGameOver() bool {
	for i := 0; i < g.NumPlayers; i++ {
		if g.GetActiveShipsCount(i+1) == 0 {
			return true
		}
	}
	return false
}

// GetWinner returns the winning player number, or 0 if game is not over
func (g *Game) GetWinner() int {
	if !g.IsGameOver() {
		return 0
	}

	// Find the player with ships remaining
	for i := 0; i < g.NumPlayers; i++ {
		if g.GetActiveShipsCount(i+1) > 0 {
			return i + 1
		}
	}
	return 0
}

// MarkShot marks a shot on the opponent's board
func (g *Game) MarkShot(playerNumber, row, column int, hit bool) {
	opponentNumber := ((playerNumber - 1 + 1) % g.NumPlayers) + 1
	opponentBoard := g.GetPlayerBoard(opponentNumber)
	
	if opponentBoard != nil {
		if hit {
			opponentBoard.MarkHit(row, column)
		} else {
			opponentBoard.MarkMiss(row, column)
		}
	}
	
	// Update player stats
	if playerNumber > 0 && playerNumber <= len(g.Players) {
		g.Players[playerNumber-1].ShotsFired++
		if hit {
			g.Players[playerNumber-1].ShotsHit++
		}
	}
}

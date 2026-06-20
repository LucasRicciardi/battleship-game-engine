package validation

import (
	"fmt"
	"regexp"
)

// Minimum and maximum board size
const (
	MinBoardSize = 5
	MaxBoardSize = 100
	MinPlayers   = 1
	MaxPlayers   = 2
	MaxShips     = 100
	MaxShots     = 10000
)

// StartGameRequest represents the request body for starting a new game
type StartGameRequest struct {
	BoardRows    int `json:"board_rows" binding:"required,min=5,max=100"`
	BoardColumns int `json:"board_columns" binding:"required,min=5,max=100"`
	NumPlayers   int `json:"num_players" binding:"required,min=1,max=2"`
}

// Validate validates the StartGameRequest
func (r *StartGameRequest) Validate() error {
	if r.BoardRows < MinBoardSize || r.BoardRows > MaxBoardSize {
		return fmt.Errorf("board_rows must be between %d and %d", MinBoardSize, MaxBoardSize)
	}
	if r.BoardColumns < MinBoardSize || r.BoardColumns > MaxBoardSize {
		return fmt.Errorf("board_columns must be between %d and %d", MinBoardSize, MaxBoardSize)
	}
	if r.NumPlayers < MinPlayers || r.NumPlayers > MaxPlayers {
		return fmt.Errorf("num_players must be between %d and %d", MinPlayers, MaxPlayers)
	}
	return nil
}

// ShootRequest represents the request body for firing a shot
type ShootRequest struct {
	GameID   string `json:"game_id" binding:"required,uuid4"`
	PlayerID string `json:"player_id" binding:"required,uuid4"`
	Row      int    `json:"row" binding:"required,min=0"`
	Column   int    `json:"column" binding:"required,min=0"`
}

// Validate validates the ShootRequest
func (r *ShootRequest) Validate(boardRows, boardColumns int) error {
	if r.Row < 0 || r.Row >= boardRows {
		return fmt.Errorf("row must be between 0 and %d", boardRows-1)
	}
	if r.Column < 0 || r.Column >= boardColumns {
		return fmt.Errorf("column must be between 0 and %d", boardColumns-1)
	}
	return nil
}

// GetGameStateRequest represents the request for getting game state
type GetGameStateRequest struct {
	GameID string `json:"game_id" binding:"required,uuid4"`
}

// Validate validates the GetGameStateRequest
func (r *GetGameStateRequest) Validate() error {
	if !isValidUUID(r.GameID) {
		return fmt.Errorf("game_id must be a valid UUID")
	}
	return nil
}

// GetGameStatsRequest represents the request for getting game statistics
type GetGameStatsRequest struct {
	GameID string `json:"game_id" binding:"required,uuid4"`
}

// Validate validates the GetGameStatsRequest
func (r *GetGameStatsRequest) Validate() error {
	if !isValidUUID(r.GameID) {
		return fmt.Errorf("game_id must be a valid UUID")
	}
	return nil
}

// isValidUUID checks if a string is a valid UUID v4
func isValidUUID(s string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(s)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"error_code,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

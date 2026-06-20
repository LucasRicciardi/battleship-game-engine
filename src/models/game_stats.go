package models

// GameStats represents statistics for a game
type GameStats struct {
	GameID         string `json:"game_id"`
	TotalTurns     int    `json:"total_turns"`
	TotalShots     int    `json:"total_shots"`
	TotalHits      int    `json:"total_hits"`
	TotalMisses    int    `json:"total_misses"`
	ShipsRemaining int    `json:"ships_remaining"`
	CurrentPlayer  int    `json:"current_player"`
	Status         string `json:"status"`
	Winner         int    `json:"winner,omitempty"`
}

// PlayerStats represents statistics for a specific player
type PlayerStats struct {
	PlayerID   string `json:"player_id"`
	ShotsFired int    `json:"shots_fired"`
	ShotsHit   int    `json:"shots_hit"`
	ShipsSunk  int    `json:"ships_sunk"`
}

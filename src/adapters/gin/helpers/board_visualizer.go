package helpers

import (
	"fmt"
	"strings"
)

// BoardVisualizer helps visualize the board state
type BoardVisualizer struct {
	Board [][]string
	Rows  int
	Cols  int
}

// NewBoardVisualizer creates a new board visualizer
func NewBoardVisualizer(board [][]string, rows, cols int) *BoardVisualizer {
	return &BoardVisualizer{
		Board: board,
		Rows:  rows,
		Cols:  cols,
	}
}

// Visualize returns a string representation of the board
func (v *BoardVisualizer) Visualize() string {
	var sb strings.Builder
	
	// Column headers
	sb.WriteString("   ")
	for i := 0; i < v.Cols; i++ {
		sb.WriteString(fmt.Sprintf("%2d ", i))
	}
	sb.WriteString("\n")
	
	// Board rows
	for i := 0; i < v.Rows; i++ {
		sb.WriteString(fmt.Sprintf("%2d ", i))
		for j := 0; j < v.Cols; j++ {
			sb.WriteString(fmt.Sprintf(" %s ", v.getCellDisplay(i, j)))
		}
		sb.WriteString("\n")
	}
	
	return sb.String()
}

// getCellDisplay returns the display character for a cell
func (v *BoardVisualizer) getCellDisplay(row, col int) string {
	if row >= v.Rows || col >= v.Cols {
		return "?"
	}
	
	cell := v.Board[row][col]
	switch cell {
	case " ":
		return " "
	case "O":
		return "O"
	case "X":
		return "X"
	case "S":
		return "S"
	default:
		return cell
	}
}

// VisualizeWithShips returns a string representation showing both player's boards
func (v *BoardVisualizer) VisualizeWithShips(myBoard, opponentBoard [][]string) string {
	var sb strings.Builder
	
	// Header
	sb.WriteString("My Board          Opponent Board\n")
	sb.WriteString("   ")
	
	// My board columns
	for i := 0; i < v.Cols; i++ {
		sb.WriteString(fmt.Sprintf("%2d ", i))
	}
	sb.WriteString("   ")
	
	// Opponent board columns
	for i := 0; i < v.Cols; i++ {
		sb.WriteString(fmt.Sprintf("%2d ", i))
	}
	sb.WriteString("\n")
	
	// Board rows
	for i := 0; i < v.Rows; i++ {
		// My board row
		sb.WriteString(fmt.Sprintf("%2d ", i))
		for j := 0; j < v.Cols; j++ {
			sb.WriteString(fmt.Sprintf(" %s ", v.getCellDisplay(myBoard, i, j)))
		}
		sb.WriteString("   ")
		
		// Opponent board row
		for j := 0; j < v.Cols; j++ {
			sb.WriteString(fmt.Sprintf(" %s ", v.getCellDisplay(opponentBoard, i, j)))
		}
		sb.WriteString("\n")
	}
	
	return sb.String()
}

// getCellDisplay returns the display character for a cell
func (v *BoardVisualizer) getCellDisplay(board [][]string, row, col int) string {
	if row >= len(board) || col >= len(board[0]) {
		return "?"
	}
	
	cell := board[row][col]
	switch cell {
	case " ":
		return " "
	case "O":
		return "O"
	case "X":
		return "X"
	case "S":
		return "S"
	default:
		return cell
	}
}

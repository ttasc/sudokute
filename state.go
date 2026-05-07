package main

import (
	"math/rand"
	"time"
)

// Constants defining the logical boundaries of the game.
const (
	GridSize    = 9
	SubgridSize = 3
)

// GameMode represents the current phase of the application.
type GameMode int

const (
	ModePlaying GameMode = iota
	ModeSolved
)

// Cell represents a single square on the Sudoku board.
type Cell struct {
	Value   uint8
	IsFixed bool
	IsError bool
}

// GameState holds all mutable data required to run the game.
type GameState struct {
	Board      [GridSize][GridSize]Cell
	CursorX    int
	CursorY    int
	ErrorCount int
	Mode       GameMode
	StartTime  time.Time
	EndTime    time.Time
}

// NewGameState initializes a new game session and allocates the state.
func NewGameState() *GameState {
	gs := &GameState{}
	gs.ResetGame()
	return gs
}

// ResetGame clears the current state, resets timers/counters, and generates a new puzzle.
func (gs *GameState) ResetGame() {
	gs.CursorX = GridSize / 2
	gs.CursorY = GridSize / 2
	gs.Mode = ModePlaying
	gs.ErrorCount = 0
	gs.StartTime = time.Now()
	gs.EndTime = time.Time{}
	gs.Board =[GridSize][GridSize]Cell{}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a completely random Sudoku board.
	var tempBoard [GridSize][GridSize]uint8
	fillGrid(&tempBoard, rng)

	// Punch holes in the board to create the test questions (Randomly delete 45 cells - Difficulty Level: Normal)
	holes := 45
	for range holes {
		for {
			x := rng.Intn(GridSize)
			y := rng.Intn(GridSize)
			if tempBoard[y][x] != 0 {
				tempBoard[y][x] = 0
				break
			}
		}
	}

	for y := range GridSize {
		for x := range GridSize {
			val := tempBoard[y][x]
			if val > 0 {
				gs.Board[y][x] = Cell{
					Value:   val,
					IsFixed: true, // The remaining numbers are marked as original numbers (cannot be edited).
					IsError: false,
				}
			}
		}
	}
}

// AUTOMATED SUDOKU GENERATION ALGORITHM (BACKTRACKING)
// fillGrid uses backtracking to completely fill a 9x9 board with valid numbers.
func fillGrid(b *[GridSize][GridSize]uint8, rng *rand.Rand) bool {
	for y := range GridSize {
		for x := range GridSize {
			// Find the first empty cell
			if b[y][x] == 0 {
				// Try random numbers from 1 to 9
				nums := rng.Perm(9)
				for _, n := range nums {
					val := uint8(n + 1)
					if isSafeToPlace(b, x, y, val) {
						b[y][x] = val

						// Recursively fill the next cell, if successful return true
						if fillGrid(b, rng) {
							return true
						}

						// If recursion fails (dead end), reset this cell and try another number
						b[y][x] = 0
					}
				}
				// No valid number found, report Backtrack error
				return false
			}
		}
	}
	return true // Table filled
}

// isSafeToPlace checks if a number can be placed at b[y][x] following Sudoku rules.
func isSafeToPlace(b *[GridSize][GridSize]uint8, x, y int, val uint8) bool {
	// Check Rows and Columns
	for i := range GridSize {
		if b[y][i] == val || b[i][x] == val {
			return false
		}
	}

	// Check Block 3x3
	startX := (x / SubgridSize) * SubgridSize
	startY := (y / SubgridSize) * SubgridSize
	for i := range SubgridSize {
		for j := range SubgridSize {
			if b[startY+i][startX+j] == val {
				return false
			}
		}
	}

	return true
}

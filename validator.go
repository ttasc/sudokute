package main

import "time"

// ValidateRules executes a full sweep of the board to find and mark duplicate numbers.
// It applies a deterministic approach: clear all flags, then mark from scratch.
func ValidateRules(state *GameState) {
	clearAllErrors(state)

	validateRows(state)
	validateColumns(state)
	validateSubgrids(state)

	countErrors(state)
}

// clearAllErrors resets the IsError flag for every cell on the board.
func clearAllErrors(state *GameState) {
	for y := range GridSize {
		for x := range GridSize {
			state.Board[y][x].IsError = false
		}
	}
}

// validateRows checks for duplicate numbers across horizontal lines.
func validateRows(state *GameState) {
	for y := range GridSize {
		// seen holds arrays of X coordinates, indexed by the cell's numeric value (1-9).
		var seen [10][]int

		for x := range GridSize {
			val := state.Board[y][x].Value
			if val > 0 {
				seen[val] = append(seen[val], x)
			}
		}

		// If a number appeared more than once, flag all its locations as errors.
		for _, coords := range seen {
			if len(coords) > 1 {
				for _, x := range coords {
					state.Board[y][x].IsError = true
				}
			}
		}
	}
}

// validateColumns checks for duplicate numbers across vertical lines.
func validateColumns(state *GameState) {
	for x := range GridSize {
		var seen[10][]int

		for y := range GridSize {
			val := state.Board[y][x].Value
			if val > 0 {
				seen[val] = append(seen[val], y)
			}
		}

		for _, coords := range seen {
			if len(coords) > 1 {
				for _, y := range coords {
					state.Board[y][x].IsError = true
				}
			}
		}
	}
}

// validateSubgrids checks for duplicates within the nine 3x3 blocks.
func validateSubgrids(state *GameState) {
	for blockY := range SubgridSize {
		for blockX := range SubgridSize {
			// seen holds arrays of [X, Y] coordinate pairs.
			var seen [10][][2]int

			// Iterate over the 9 cells within the current subgrid.
			for dy := range SubgridSize {
				for dx := range SubgridSize {
					y := (blockY * SubgridSize) + dy
					x := (blockX * SubgridSize) + dx
					val := state.Board[y][x].Value

					if val > 0 {
						seen[val] = append(seen[val], [2]int{x, y})
					}
				}
			}

			for _, coords := range seen {
				if len(coords) > 1 {
					for _, pos := range coords {
						state.Board[pos[1]][pos[0]].IsError = true
					}
				}
			}
		}
	}
}

// countErrors sums up all active IsError flags to update the game state.
func countErrors(state *GameState) {
	state.ErrorCount = 0
	for y := range GridSize {
		for x := range GridSize {
			if state.Board[y][x].IsError {
				state.ErrorCount++
			}
		}
	}
}

// CheckWinCondition evaluates if the board is completely and correctly filled.
// If true, it transitions the game mode to Solved.
func CheckWinCondition(state *GameState) {
	if state.Mode == ModeSolved {
		return
	}

	// Condition 1: No empty cells remaining
	for y := range GridSize {
		for x := range GridSize {
			if state.Board[y][x].Value == 0 {
				return
			}
		}
	}

	// Condition 2: No active rule violations
	if state.ErrorCount == 0 {
		state.Mode = ModeSolved
		state.EndTime = time.Now()
	}
}

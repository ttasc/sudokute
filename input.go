package main

import "github.com/ttasc/ttbox"

// handleInput processes a single terminal event and updates the GameState.
// It routes keyboard and mouse inputs to the appropriate game logic.
func handleInput(gs *GameState, evt ttbox.Event) {
	if gs.Mode == ModeSolved {
		if evt.Type == ttbox.EventKey && (evt.Ch == 'r' || evt.Ch == 'R') {
			gs.ResetGame()
			ValidateRules(gs) // Run a check on the new table to update the error count to 0.
		}
		return
	}

	switch evt.Type {
	case ttbox.EventKey:
		handleKeyEvent(gs, evt)
	case ttbox.EventMouse:
		handleMouseEvent(gs, evt)
	}
}

// handleKeyEvent modifies the cursor position or updates cell values based on key presses.
func handleKeyEvent(gs *GameState, evt ttbox.Event) {
	switch evt.Key {
	case ttbox.KeyArrowUp:
		gs.CursorY--
	case ttbox.KeyArrowDown:
		gs.CursorY++
	case ttbox.KeyArrowLeft:
		gs.CursorX--
	case ttbox.KeyArrowRight:
		gs.CursorX++
	case ttbox.KeyBackspace, ttbox.KeyDelete:
		updateCellValue(gs, 0)
	default:
		switch evt.Ch {
		case 'k', 'K':
			gs.CursorY--
		case 'j', 'J':
			gs.CursorY++
		case 'h', 'H':
			gs.CursorX--
		case 'l', 'L':
			gs.CursorX++
		case 'r', 'R':
			gs.ResetGame()
			ValidateRules(gs)
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			updateCellValue(gs, uint8(evt.Ch-'0'))
		case '0', ' ':
			updateCellValue(gs, 0)
		}
	}

	clampCursor(gs)
}

// handleMouseEvent maps physical mouse clicks to logical board coordinates.
func handleMouseEvent(gs *GameState, evt ttbox.Event) {
	if evt.Button != ttbox.MouseLeft || !evt.Press {
		return
	}

	offsetX, offsetY := GetBoardOffset()
	relX := evt.X - offsetX
	relY := evt.Y - offsetY

	// Only process if the mouse click is within the Sudoku grid area.
	if relX >= 0 && relX < BoardTotalWidth && relY >= 0 && relY < BoardTotalHeight {
		// Integer division will automatically convert pixels to array coordinates (0-8).
		targetX := relX / CellWidth
		targetY := relY / CellHeight

		gs.CursorX = targetX
		gs.CursorY = targetY
		clampCursor(gs)
	}
}

// updateCellValue mutates the selected cell if it is not fixed, then triggers validations.
func updateCellValue(gs *GameState, newValue uint8) {
	cell := &gs.Board[gs.CursorY][gs.CursorX]

	// Only allow corrections if the number is not the original number from the problem statement.
	if !cell.IsFixed {
		cell.Value = newValue
		ValidateRules(gs)
		CheckWinCondition(gs)
	}
}

// clampCursor ensures the cursor never leaves the 0-8 logical boundaries.
func clampCursor(gs *GameState) {
	if gs.CursorX < 0 {
		gs.CursorX = 0
	}
	if gs.CursorX > GridSize-1 {
		gs.CursorX = GridSize - 1
	}
	if gs.CursorY < 0 {
		gs.CursorY = 0
	}
	if gs.CursorY > GridSize-1 {
		gs.CursorY = GridSize - 1
	}
}

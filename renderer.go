package main

import (
	"fmt"
	"time"

	"github.com/ttasc/ttbox"
)

// UI Configuration & Board Dimensions
const (
	BoardTotalWidth  = 37 // (9 cells * 4 width) + 1 right border
	BoardTotalHeight = 19 // (9 cells * 2 height) + 1 bottom border
	CellWidth        = 4
	CellHeight       = 2

	CharLeftBracket  = '['
	CharRightBracket = ']'
)

// Modern and minimalist xterm-256 color palette (Shared with Gomoku)
const (
	ColorMajorBorder = 245 // Light gray for 3x3 subgrid borders
	ColorMinorBorder = 239 // Dark gray for standard grid lines

	ColorFixedNum    = 255 // Pure white for initial puzzle numbers
	ColorUserNum     = 39  // Bright blue for user-placed numbers
	ColorErrorNum    = 196 // Bright red for rule violations

	ColorSelValid    = 39  // Bright blue for the active cursor
	ColorSelInvalid  = 196 // Red for cursor if hovering over an error
	ColorWin         = 114 // Pastel green for the solved state

	ColorMatchBg     = 236 // Dark gray background for highlighting matching numbers
	ColorText        = 250 // Light gray text
	ColorTextDim     = 240 // Dim gray text for inactive elements
	ColorBgModal     = 235 // Background color for the popup modal
)

// RenderGame draws the entire game interface including board, cells, and status line.
func RenderGame(state *GameState) {
	ttbox.Clear()

	drawStatusline(state)
	drawBoard(state)

	if state.Mode == ModeSolved {
		drawEndgamePopup(state)
	}

	drawControlsGuide()

	ttbox.Present()
}

// GetBoardOffset calculates the top-left coordinate to perfectly center the board.
func GetBoardOffset() (int, int) {
	termW, termH := ttbox.Size()

	offsetX := (termW - BoardTotalWidth) / 2
	offsetY := (termH - BoardTotalHeight) / 2

	if offsetX < 0 { offsetX = 0 }
	if offsetY < 0 { offsetY = 0 }

	return offsetX, offsetY
}

// drawBoard combines drawing the grid lines and the numbers inside them.
func drawBoard(state *GameState) {
	offsetX, offsetY := GetBoardOffset()
	drawGrid(state, offsetX, offsetY)
	drawCells(state, offsetX, offsetY)
}

// drawGrid renders the empty Sudoku grid framework with major/minor borders.
func drawGrid(state *GameState, offsetX, offsetY int) {
	for gy := 0; gy <= BoardTotalHeight-1; gy++ {
		for gx := 0; gx <= BoardTotalWidth-1; gx++ {

			isMajorY := (gy % (CellHeight * SubgridSize)) == 0
			isMajorX := (gx % (CellWidth * SubgridSize)) == 0
			isYLine  := (gy % CellHeight) == 0
			isXLine  := (gx % CellWidth) == 0

			color := ColorMinorBorder
			if isMajorY || isMajorX {
				color = ColorMajorBorder
			}

			if state.Mode == ModeSolved {
				color = ColorWin
			}

			char := determineGridChar(gx, gy, isXLine, isYLine)
			if char != ' ' {
				ttbox.SetCell(offsetX+gx, offsetY+gy, char, color, ttbox.ColorDefault)
			}
		}
	}
}

// determineGridChar selects the correct box-drawing character based on intersection.
func determineGridChar(gx, gy int, isXLine, isYLine bool) rune {
	if isYLine && isXLine {
		if gy == 0 {
			if gx == 0 { return '┌' }
			if gx == BoardTotalWidth-1 { return '┐' }
			return '┬'
		}
		if gy == BoardTotalHeight-1 {
			if gx == 0 { return '└' }
			if gx == BoardTotalWidth-1 { return '┘' }
			return '┴'
		}
		if gx == 0 { return '├' }
		if gx == BoardTotalWidth-1 { return '┤' }
		return '┼'
	}
	if isYLine { return '─' }
	if isXLine { return '│' }
	return ' '
}

// drawCells renders the numbers, background highlights, and the cursor into the grid.
func drawCells(state *GameState, offsetX, offsetY int) {
	hoveredVal := state.Board[state.CursorY][state.CursorX].Value

	for y := range GridSize {
		for x := range GridSize {
			cell := state.Board[y][x]

			termX := offsetX + (x * CellWidth) + (CellWidth / 2)
			termY := offsetY + (y * CellHeight) + (CellHeight / 2)

			// Background Color
			bg := ttbox.ColorDefault
			if hoveredVal > 0 && cell.Value == hoveredVal {
				bg = ColorMatchBg // Highlight the numbers that are similar to the number being pointed to.
			}

			// Foreground Color
			fg := ColorUserNum
			if cell.IsFixed {
				fg = ColorFixedNum
			}
			if cell.IsError {
				fg = ColorErrorNum
			}
			if state.Mode == ModeSolved {
				fg = ColorWin
			}

			char := ' '
			if cell.Value > 0 {
				char = rune('0' + cell.Value)
			}

			// Cursor effects [ ]
			leftChar, rightChar := ' ', ' '
			bracketFg := ColorSelValid

			if x == state.CursorX && y == state.CursorY && state.Mode != ModeSolved {
				leftChar = CharLeftBracket
				rightChar = CharRightBracket
				if cell.IsError {
					bracketFg = ColorSelInvalid
				}
			}

			ttbox.SetCell(termX-1, termY, leftChar, bracketFg, bg)

			// Draw the number (Use bold if it is the original number from the problem.)
			if cell.IsFixed && cell.Value > 0 {
				ttbox.SetAttr(true, false, false, false)
				ttbox.SetCell(termX, termY, char, fg, bg)
				ttbox.ResetAttr()
			} else {
				ttbox.SetCell(termX, termY, char, fg, bg)
			}

			ttbox.SetCell(termX+1, termY, rightChar, bracketFg, bg)
		}
	}
}

// drawStatusline renders the top information bar matching the Gomoku UI style.
func drawStatusline(state *GameState) {
	w, h := ttbox.Size()
	if w == 0 || h == 0 {
		return
	}

	_, offsetY := GetBoardOffset()
	y := max(offsetY-2, 0)

	if y != 0 {
		ttbox.DrawTextCenter(1, " S U D O K U ", ColorText, ttbox.ColorDefault)
	}

	elapsed := getElapsedTime(state)
	hours := int(elapsed.Hours())
	mins := int(elapsed.Minutes()) % 60
	secs := int(elapsed.Seconds()) % 60
	timerText := fmt.Sprintf("  %02d:%02d:%02d  ", hours, mins, secs)

	centerX := w / 2

	// Setup Text & Colors for the two sides
	modeStr := " PLAYING "
	modeColor := ColorTextDim
	if state.Mode == ModeSolved {
		modeStr = " SOLVED "
		modeColor = ColorWin
	}

	errStr := fmt.Sprintf(" ERRORS: %d ", state.ErrorCount)
	errColor := ColorTextDim
	if state.ErrorCount > 0 {
		errColor = ColorErrorNum
	}

	// 1. Draw Left Indicator (Mode)
	modeX := centerX - (len(timerText) / 2) - len(modeStr)
	for i, ch := range modeStr {
		ttbox.SetCell(modeX+i, y, ch, modeColor, ttbox.ColorDefault)
	}

	// 2. Draw Center (Timer)
	ttbox.DrawTextCenter(y, timerText, ColorText, ttbox.ColorDefault)

	// 3. Draw Right Indicator (Errors)
	errX := centerX + (len(timerText) / 2) + (len(timerText) % 2)
	for i, ch := range errStr {
		ttbox.SetCell(errX+i, y, ch, errColor, ttbox.ColorDefault)
	}
}

// drawControlsGuide renders the bottom instructions bar to help players with keybindings.
func drawControlsGuide() {
	w, h := ttbox.Size()
	if w == 0 || h == 0 {
		return
	}

	guideText := " Move(h, j, k, l; arrows)   Input(1-9)   Clear(space, 0, backspace)   Quit(Q, Ctrl+C, Esc) "
	ttbox.DrawTextCenter(h-1, guideText, ColorText, ttbox.ColorDefault)
}

// drawEndgamePopup displays a modal dialogue when the puzzle is completely solved.
func drawEndgamePopup(state *GameState) {
	w, h := ttbox.Size()

	msgFg := ColorWin
	msg := " SUDOKU SOLVED! "
	subMsg := "[R] Play Again   [ESC] Exit "

	boxW := len(subMsg) + 6
	boxH := 4

	x := (w - boxW) / 2
	_, offsetY := GetBoardOffset()

	// Dynamic Y positioning: Shift the popup to the opposite vertical half of the screen relative
	// to the cursor. This prevents the modal from obscuring the cell that ended the game.
	var y int
	if state.CursorY < GridSize/2 {
		y = offsetY + BoardTotalHeight
		if y+boxH > h {
			y = h - boxH
		}
	} else {
		y = max(offsetY-boxH-1, 0)
	}

	ttbox.SetColor(ColorText, ColorBgModal)
	ttbox.Fill(x, y, boxW, boxH, ' ')

	// Box
	ttbox.SetColor(ColorMinorBorder, ColorBgModal)
	ttbox.Box(x, y, boxW, boxH)

	ttbox.SetAttr(true, false, false, false)
	ttbox.DrawTextCenter(y+1, msg, msgFg, ColorBgModal)
	ttbox.ResetAttr()

	ttbox.DrawTextCenter(y+2, subMsg, ColorTextDim, ColorBgModal)

	// Reset Color
	ttbox.SetColor(ttbox.ColorDefault, ttbox.ColorDefault)
}

// getElapsedTime safely calculates the duration of the current game.
func getElapsedTime(state *GameState) time.Duration {
	if state.Mode == ModeSolved {
		return state.EndTime.Sub(state.StartTime)
	}
	return time.Since(state.StartTime)
}

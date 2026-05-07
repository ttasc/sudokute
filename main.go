package main

import (
	"fmt"
	"time"

	"github.com/ttasc/ttbox"
)

// main initializes the terminal environment and runs the continuous game loop.
func main() {
	if err := ttbox.Init(); err != nil {
		fmt.Printf("Error initializing TUI: %v\n", err)
		return
	}
	defer ttbox.Close()

	ttbox.HideCursorFunc()
	ttbox.EnableMouseFunc()
	defer ttbox.DisableMouseFunc()

	state := NewGameState()

	// Run the rule check the first time to ensure the counters are initialized correctly.
	ValidateRules(state)

	isRunning := true
	for isRunning {
		// Timeout 100ms (~10 FPS)
		evt, err := ttbox.PollEventTimeout(100 * time.Millisecond)

		if err == nil {
			if evt.Type == ttbox.EventKey && (evt.Key == ttbox.KeyEscape || evt.Key == ttbox.KeyCtrlC || evt.Ch == 'q' || evt.Ch == 'Q') {
				isRunning = false
			} else {
				handleInput(state, evt)
			}
		}

		// Redraw the entire terminal frame based on the mutated state.
		RenderGame(state)
	}
}

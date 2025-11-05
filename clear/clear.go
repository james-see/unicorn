package clear

import (
	"time"

	tm "github.com/buger/goterm"
	"github.com/pterm/pterm"
)

// ClearIt clears the screen with a loading spinner
func ClearIt() {
	// Show spinner - make it visible
	spinner, _ := pterm.DefaultSpinner.
		WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").
		WithRemoveWhenDone(false).
		Start("Refreshing screen...")
	
	// Let spinner animate visibly - longer duration so user can see it
	time.Sleep(500 * time.Millisecond)
	
	// Stop spinner
	spinner.Stop()
	
	// Clear the screen
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Println("Current Time:", time.Now().Format(time.RFC1123))
	tm.Flush()
	
	// Small delay after clearing
	time.Sleep(100 * time.Millisecond)
}

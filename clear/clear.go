package clear

import (
	"time"

	tm "github.com/buger/goterm"
	"github.com/pterm/pterm"
)

// ClearIt clears the screen with a loading spinner
func ClearIt() {
	// Show spinner before clearing
	spinner, _ := pterm.DefaultSpinner.
		WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").
		WithRemoveWhenDone(true).
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
	
	// Show spinner again after clearing (on the "right side")
	spinner2, _ := pterm.DefaultSpinner.
		WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").
		WithRemoveWhenDone(true).
		Start("Loading...")
	
	time.Sleep(500 * time.Millisecond)
	spinner2.Stop()
}

package main

import (
	"fmt"
	"os"

	"github.com/jamesacampbell/unicorn/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
		os.Exit(1)
	}
}

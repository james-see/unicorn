package main

import (
	"fmt"
	"os"

	"github.com/jamesacampbell/unicorn/tui"
	"github.com/jamesacampbell/unicorn/version"
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("%s\n%s\n", version.String(), version.ReleaseInfoURL())
			os.Exit(0)
		}
	}
	if err := tui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
		os.Exit(1)
	}
}

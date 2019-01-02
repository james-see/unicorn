package main

import (
	"cmd/version"

	"cmd/menu"
)

func main() {
	menu.Execute()
	version.Execute()
}

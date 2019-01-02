package menu

import "fmt"

// DisplayMenu ...
func DisplayMenu(username string) {
	fmt.Printf("%s, you have $250,000 to build your initial portfolio.\nEach turn is 1 month of time and you will face questions at each turn that will either increase or decrease your wealth. \nYour game ends after 10 years (roughly 120 turns).\n", username)
}

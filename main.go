package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	logo "github.com/jamesacampbell/unicorn/logo"
	menu "github.com/jamesacampbell/unicorn/menu"
)

type userData struct {
	username string
}

func initMenu() (username string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your Name: ")
	text, _ := reader.ReadString('\n')
	fmt.Printf("Welcome %s\n", text)
	return text
}

func main() {
	// init with the Cyan unicorn logo b*tch
	c := color.New(color.FgCyan)
	logo.InitLogo(c)
	// get the username
	s := userData{initMenu()}
	// display the intro menu
	menu.DisplayMenu(s.username)

}

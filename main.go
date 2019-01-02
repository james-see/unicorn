package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	clear "github.com/jamesacampbell/unicorn/clear"
	logo "github.com/jamesacampbell/unicorn/logo"
	menu "github.com/jamesacampbell/unicorn/menu"
)

var pot int64

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
	c := color.New(color.FgCyan).Add(color.Bold)
	logo.InitLogo(c)
	// get the username
	s := userData{initMenu()}
	// display the intro menu
	clear.ClearIt()
	pot = 250000
	menu.DisplayMenu(s.username, pot)
	clear.ClearIt()
	for i := 1; i <= 2; i++ {
		fmt.Printf("Company %d:", i)
		menu.DisplayStartups(s.username, pot, i)
		fmt.Print("Press 'Enter' to see next company...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		clear.ClearIt()
	}
}

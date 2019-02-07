package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	clear "github.com/jamesacampbell/unicorn/clear"
	logo "github.com/jamesacampbell/unicorn/logo"
	menu "github.com/jamesacampbell/unicorn/menu"
	yaml "gopkg.in/yaml.v2"
)

type userData struct {
	username string
}

type gameData struct {
	Pot        int64  `yaml:"starting-cash"`
	BadThings  int64  `yaml:"number-of-bad-things-per-year"`
	Foreground string `yaml:"foreground-color"`
}

func initMenu() (username string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your Name: ")
	text, _ := reader.ReadString('\n')
	fmt.Printf("Welcome %s\n", text)
	return text
}

func initGame() gameData {
	var gd gameData
	yamlFile, err := os.Open("config/data.yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer yamlFile.Close()
	byteValue, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal(byteValue, &gd)
	if err != nil {
		fmt.Println(err)
	}
	return gd
}

func main() {
	// init with the unicorn logo b*tch
	gamesettings := initGame()
	c := color.New(color.FgCyan).Add(color.Bold)
	logo.InitLogo(c)
	// get the username
	s := userData{initMenu()}
	// display the intro menu
	clear.ClearIt()
	pot := gamesettings.Pot
	menu.DisplayMenu(s.username, pot)
	clear.ClearIt()
	for i := 1; i <= 2; i++ {
		fmt.Printf("Company %d:", i)
		menu.DisplayStartups(s.username, pot, i)
		if i == 2 {
			fmt.Print("Press 'Enter' to finish round")
		} else {
			fmt.Print("Press 'Enter' to see next company...")
		}
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		clear.ClearIt()
	}
}

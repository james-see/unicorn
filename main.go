package main

import (
	"bufio"
	"fmt"
	"os"

	Logo "unicorn/functions/Logo"
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
	Logo.initLogo()
	s := userData{initMenu()}
	fmt.Println(s.username)

}

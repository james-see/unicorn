package main

import (
	"bufio"
	"fmt"
	"os"

	logo "github.com/jamesacampbell/unicorn/logo"
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
	logo.InitLogo()
	s := userData{initMenu()}
	fmt.Println(s.username)

}

package menu

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	clear "github.com/jamesacampbell/unicorn/clear"
)

// StartupData ...
type StartupData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DisplayMenu ...
func DisplayMenu(username string, pot int64) {
	fmt.Printf("%v", username)
	fmt.Printf("you have $%d to build your initial portfolio.\nEach turn is 1 month of time and you will face questions at each turn that will either increase or decrease your wealth. \nYour game ends after 10 years (roughly 120 turns).\n", pot)
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	clear.ClearIt()
	fmt.Printf("That is about it, %v", username)
	fmt.Println("\nYou will now get to read and decide which startups you will back and how much.")
}

func loadJSONData(companyid int) StartupData {
	var t StartupData
	jsonFile, err := os.Open(fmt.Sprintf("startups/%d.json", companyid))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &t)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

// DisplayStartups ...
func DisplayStartups(username string, pot int64, companyid int) {
	t := loadJSONData(companyid)
	fmt.Println(t.Name)
	return
}

package menu

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/jamesacampbell/unicorn/assets"
	clear "github.com/jamesacampbell/unicorn/clear"
)

// StartupData ...
type StartupData struct {
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Category               string `json:"category"`
	Valuation              int    `json:"valuation"`
	GrossBurnRate          int    `json:"grossburnrate"`
	MonthlyActivationRate  int    `json:"Monthly Activation Rate"`
	MonthlyWebsiteVisitors int    `json:"Monthly Active Visitors"`
	MonthlySales           int    `json:"Monthly Sales"`
	Cost                   int    `json:"Cost"`
	SalePrice              int    `json:"Sale Price"`
	PercentMargin          int    `json:"Percent Margin Per Unit"`
}

// DisplayMenu ...
func DisplayMenu(username string, pot int64) {
	fmt.Printf("%v", username)
	fmt.Printf("you have $%d to build your initial portfolio.\nEach turn is 1 month of time and you will face questions at each turn that will either increase or decrease your wealth. \nYour game ends after 5 years (60 turns).\n", pot)
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	clear.ClearIt()
	fmt.Printf("That is about it, %v", username)
	fmt.Println("\nYou will now get to read and decide which startups you will back and how much.")
}

func loadJSONData(companyid int) StartupData {
	var t StartupData
	byteValue, err := assets.ReadStartupFile(companyid)
	if err != nil {
		fmt.Println(err)
		return t
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
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s %v\n",
			typeOfT.Field(i).Name, f.Interface())
	}
	return
}

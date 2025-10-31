package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	clear "github.com/jamesacampbell/unicorn/clear"
	game "github.com/jamesacampbell/unicorn/game"
	logo "github.com/jamesacampbell/unicorn/logo"
	yaml "gopkg.in/yaml.v2"
)

type gameData struct {
	Pot        int64  `yaml:"starting-cash"`
	BadThings  int64  `yaml:"number-of-bad-things-per-year"`
	Foreground string `yaml:"foreground-color"`
}

func initMenu() (username string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your Name: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	fmt.Printf("\nWelcome %s!\n", text)
	return text
}

func loadConfig() gameData {
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

func displayWelcome(username string, startingCash int64) {
	fmt.Printf("\n%s, you have $%s to invest over the next 10 years.\n", username, formatMoney(startingCash))
	fmt.Println("Each turn = 1 month. You have 120 turns to build your fortune.")
	fmt.Println("Choose your investments wisely - events will affect valuations!")
	fmt.Print("\nPress 'Enter' to see available startups...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func displayStartup(s game.Startup, index int) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	cyan.Printf("\n[%d] %s\n", index+1, s.Name)
	fmt.Printf("    %s\n", s.Description)
	yellow.Printf("    Category: %s\n", s.Category)
	fmt.Printf("    Valuation: $%s\n", formatMoney(s.Valuation))
	fmt.Printf("    Monthly Sales: %d units\n", s.MonthlySales)
	fmt.Printf("    Margin: %d%%\n", s.PercentMargin)
	fmt.Printf("    Website Visitors: %s/month\n", formatNumber(s.MonthlyWebsiteVisitors))
	
	// Risk indicator
	riskColor := color.New(color.FgGreen)
	riskLabel := "Low"
	if s.RiskScore > 0.6 {
		riskColor = color.New(color.FgRed)
		riskLabel = "High"
	} else if s.RiskScore > 0.4 {
		riskColor = color.New(color.FgYellow)
		riskLabel = "Medium"
	}
	riskColor.Printf("    Risk: %s", riskLabel)
	
	// Growth indicator
	growthColor := color.New(color.FgGreen)
	growthLabel := "High"
	if s.GrowthPotential < 0.4 {
		growthColor = color.New(color.FgRed)
		growthLabel = "Low"
	} else if s.GrowthPotential < 0.6 {
		growthColor = color.New(color.FgYellow)
		growthLabel = "Medium"
	}
	growthColor.Printf(" | Growth Potential: %s\n", growthLabel)
}

func investmentPhase(gs *game.GameState) {
	clear.ClearIt()
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("\n?? INVESTMENT PHASE - Turn %d/%d\n", gs.Portfolio.Turn, gs.Portfolio.MaxTurns)
	fmt.Printf("Cash Available: $%s\n", formatMoney(gs.Portfolio.Cash))
	fmt.Printf("Portfolio Value: $%s\n", formatMoney(gs.GetPortfolioValue()))
	fmt.Printf("Net Worth: $%s\n", formatMoney(gs.Portfolio.NetWorth))
	
	// Show available startups
	fmt.Println("\n????????????????????????????????????????????")
	fmt.Println("AVAILABLE STARTUPS:")
	for i, startup := range gs.AvailableStartups {
		displayStartup(startup, i)
	}
	fmt.Println("????????????????????????????????????????????")
	
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Printf("\nEnter company number (1-%d) to invest, or 'done' to continue: ", len(gs.AvailableStartups))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "done" || input == "" {
			break
		}
		
		companyNum, err := strconv.Atoi(input)
		if err != nil || companyNum < 1 || companyNum > len(gs.AvailableStartups) {
			color.Red("Invalid company number!")
			continue
		}
		
		fmt.Printf("Enter investment amount ($): ")
		amountStr, _ := reader.ReadString('\n')
		amountStr = strings.TrimSpace(amountStr)
		amount, err := strconv.ParseInt(amountStr, 10, 64)
		
		if err != nil {
			color.Red("Invalid amount!")
			continue
		}
		
		err = gs.MakeInvestment(companyNum-1, amount)
		if err != nil {
			color.Red("Error: %v", err)
		} else {
			color.Green("? Investment successful!")
			fmt.Printf("Cash remaining: $%s\n", formatMoney(gs.Portfolio.Cash))
		}
	}
}

func playTurn(gs *game.GameState) {
	clear.ClearIt()
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("\n?? MONTH %d of %d\n", gs.Portfolio.Turn, gs.Portfolio.MaxTurns)
	
	messages := gs.ProcessTurn()
	
	if len(messages) > 0 {
		fmt.Println("\n????????????????????????????????????????????")
		fmt.Println("COMPANY NEWS:")
		for _, msg := range messages {
			fmt.Println(msg)
		}
		fmt.Println("????????????????????????????????????????????")
	}
	
	// Show portfolio status
	fmt.Println("\n?? YOUR PORTFOLIO:")
	if len(gs.Portfolio.Investments) == 0 {
		fmt.Println("   No investments yet")
	} else {
		for _, inv := range gs.Portfolio.Investments {
			value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
			profit := value - inv.AmountInvested
			profitColor := color.New(color.FgGreen)
			profitSign := "+"
			if profit < 0 {
				profitColor = color.New(color.FgRed)
				profitSign = ""
			}
			
			fmt.Printf("   %s: $%s invested, %.2f%% equity\n", 
				inv.CompanyName, formatMoney(inv.AmountInvested), inv.EquityPercent)
			fmt.Printf("      Current Value: $%s ", formatMoney(value))
			profitColor.Printf("(%s$%s)\n", profitSign, formatMoney(abs(profit)))
		}
	}
	
	fmt.Printf("\n?? Net Worth: $%s\n", formatMoney(gs.Portfolio.NetWorth))
	
	fmt.Print("\nPress 'Enter' to continue to next month...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func displayFinalScore(gs *game.GameState) {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	
	cyan.Println("\n" + strings.Repeat("?", 50))
	cyan.Println("           ?? GAME OVER - FINAL RESULTS ??")
	cyan.Println(strings.Repeat("?", 50))
	
	netWorth, roi, successfulExits := gs.GetFinalScore()
	
	fmt.Printf("\n?? Player: %s\n", gs.PlayerName)
	fmt.Printf("?? Turns Played: %d\n\n", gs.Portfolio.Turn-1)
	
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("?? Final Net Worth: $%s\n", formatMoney(netWorth))
	
	roiColor := color.New(color.FgGreen)
	if roi < 0 {
		roiColor = color.New(color.FgRed)
	}
	roiColor.Printf("?? Return on Investment: %.2f%%\n", roi)
	fmt.Printf("?? Successful Exits (5x+): %d\n", successfulExits)
	
	fmt.Println("\n?? FINAL PORTFOLIO:")
	for _, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		fmt.Printf("   %s: $%s ? $%s\n", 
			inv.CompanyName, formatMoney(inv.AmountInvested), formatMoney(value))
	}
	
	// Performance rating
	fmt.Println("\n" + strings.Repeat("?", 50))
	var rating string
	if roi >= 1000 {
		rating = "?? UNICORN HUNTER - Legendary!"
	} else if roi >= 500 {
		rating = "?? Elite VC - Outstanding!"
	} else if roi >= 200 {
		rating = "? Great Investor - Excellent!"
	} else if roi >= 50 {
		rating = "?? Solid Performance - Good!"
	} else if roi >= 0 {
		rating = "?? Break Even - Not Bad"
	} else {
		rating = "?? Lost Money - Better Luck Next Time"
	}
	
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("Rating: %s\n", rating)
	fmt.Println(strings.Repeat("?", 50) + "\n")
}

func main() {
	// Init with the unicorn logo
	config := loadConfig()
	c := color.New(color.FgCyan).Add(color.Bold)
	logo.InitLogo(c)
	
	// Get username
	username := initMenu()
	clear.ClearIt()
	
	// Display welcome and rules
	displayWelcome(username, config.Pot)
	
	// Initialize game
	gs := game.NewGame(username, config.Pot)
	
	// Investment phase at start
	investmentPhase(gs)
	
	// Main game loop
	for !gs.IsGameOver() {
		playTurn(gs)
	}
	
	// Show final score
	displayFinalScore(gs)
}

// Helper functions
func formatMoney(amount int64) string {
	abs := amount
	if abs < 0 {
		abs = -abs
	}
	
	s := strconv.FormatInt(abs, 10)
	
	// Add commas
	n := len(s)
	if n <= 3 {
		if amount < 0 {
			return "-" + s
		}
		return s
	}
	
	result := ""
	for i, digit := range s {
		if i > 0 && (n-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	
	if amount < 0 {
		return "-" + result
	}
	return result
}

func formatNumber(n int) string {
	return formatMoney(int64(n))
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

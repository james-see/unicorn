package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/animations"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/leaderboard"
)


func ShowTopScores(sortBy string, difficulty string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	var scores []database.GameScore
	var err error
	var title string

	// Show spinner while loading
	spinner, _ := animations.StartSpinner("Loading leaderboard...")
	defer spinner.Stop()

	if sortBy == "roi" {
		scores, err = database.GetTopScoresByROI(10, difficulty)
		title = "TOP 10 BY ROI"
	} else {
		scores, err = database.GetTopScoresByNetWorth(10, difficulty)
		title = "TOP 10 BY NET WORTH"
	}

	if difficulty != "all" && difficulty != "" {
		title += fmt.Sprintf(" (%s)", strings.ToUpper(difficulty))
	}

	if err != nil {
		color.Red("Error loading leaderboard: %v", err)
		return
	}

	cyan.Println("\n" + strings.Repeat("=", 90))
	cyan.Printf("%-40s\n", title)
	cyan.Println(strings.Repeat("=", 90))

	if len(scores) == 0 {
		yellow.Println("\nNo games played yet! Be the first!")
		return
	}

	fmt.Printf("\n%-5s %-20s %-15s %-15s %-10s %-12s\n",
		"RANK", "PLAYER", "NET WORTH", "ROI", "EXITS", "DIFFICULTY")
	fmt.Println(strings.Repeat("-", 90))

	for i, score := range scores {
		rankColor := color.New(color.FgWhite)
		if i == 0 {
			rankColor = color.New(color.FgYellow, color.Bold)
		} else if i == 1 {
			rankColor = color.New(color.FgCyan)
		} else if i == 2 {
			rankColor = color.New(color.FgGreen)
		}

		roiColor := color.New(color.FgGreen)
		if score.ROI < 0 {
			roiColor = color.New(color.FgRed)
		}

		rankColor.Printf("%-5d ", i+1)
		fmt.Printf("%-20s ", truncateString(score.PlayerName, 20))
		fmt.Printf("$%-14s ", FormatMoney(score.FinalNetWorth))
		roiColor.Printf("%-15s ", fmt.Sprintf("%.1f%%", score.ROI))
		fmt.Printf("%-10d ", score.SuccessfulExits)
		fmt.Printf("%-12s\n", score.Difficulty)
	}
}

func ShowRecentGames() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	// Show spinner while loading
	spinner, _ := animations.StartSpinner("Loading recent games...")
	defer spinner.Stop()

	scores, err := database.GetRecentGames(10)
	if err != nil {
		color.Red("Error loading recent games: %v", err)
		return
	}

	cyan.Println("\n" + strings.Repeat("=", 90))
	cyan.Println("                           RECENT GAMES")
	cyan.Println(strings.Repeat("=", 90))

	if len(scores) == 0 {
		yellow.Println("\nNo games played yet!")
		return
	}

	fmt.Printf("\n%-20s %-15s %-15s %-12s %-20s\n",
		"PLAYER", "NET WORTH", "ROI", "DIFFICULTY", "DATE")
	fmt.Println(strings.Repeat("-", 90))

	for _, score := range scores {
		roiColor := color.New(color.FgGreen)
		if score.ROI < 0 {
			roiColor = color.New(color.FgRed)
		}

		fmt.Printf("%-20s ", truncateString(score.PlayerName, 20))
		fmt.Printf("$%-14s ", FormatMoney(score.FinalNetWorth))
		roiColor.Printf("%-15s ", fmt.Sprintf("%.1f%%", score.ROI))
		fmt.Printf("%-12s ", score.Difficulty)
		fmt.Printf("%-20s\n", score.PlayedAt.Format("2006-01-02 15:04"))
	}
}

func AskToSubmitToGlobalLeaderboard(score database.GameScore) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("           ?? GLOBAL LEADERBOARD")
	fmt.Println(strings.Repeat("=", 60))

	yellow.Println("\nWould you like to submit your score to the global leaderboard?")
	fmt.Println("Your score will be visible to all players worldwide!")
	fmt.Print("\nSubmit to global leaderboard? (y/n, default n): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	if choice != "y" && choice != "yes" {
		fmt.Println("Okay, score saved locally only.")
		return
	}

	// Check API availability
	fmt.Print("\nChecking global leaderboard service...")
	if !leaderboard.IsAPIAvailable("") {
		color.Yellow("\n??  Global leaderboard service is not available right now.")
		color.Yellow("Your score has been saved locally.")
		return
	}
	color.Green(" ?")

	// Submit score
	fmt.Print("Submitting your score...")
	submission := leaderboard.ScoreSubmission{
		PlayerName:      score.PlayerName,
		FinalNetWorth:   score.FinalNetWorth,
		ROI:             score.ROI,
		SuccessfulExits: score.SuccessfulExits,
		TurnsPlayed:     score.TurnsPlayed,
		Difficulty:      score.Difficulty,
	}

	err := leaderboard.SubmitScore(submission, "")
	if err != nil {
		color.Red("\n? Failed to submit score: %v", err)
		color.Yellow("Your score has been saved locally.")
		return
	}

	color.Green(" ?")
	cyan.Println("\n?? Success! Your score has been submitted to the global leaderboard!")
	yellow.Println("\nView the global leaderboard at:")
	yellow.Println("https://james-see.github.io/unicorn")
}
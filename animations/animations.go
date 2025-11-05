package animations

import (
	"time"

	"github.com/pterm/pterm"
)

// ShowSplashScreen displays an animated splash screen
func ShowSplashScreen() {
	// Clear screen
	pterm.DefaultBox.WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).
		WithTitle("🦄 UNICORN 🦄").
		WithTitleTopCenter().
		Println("The Ultimate VC Simulation Game\n\n" +
			"Where Dreams Become Unicorns... Or Don't\n\n" +
			"Copyright 2025")

	time.Sleep(500 * time.Millisecond)

	// Create a gradient effect for loading
	spinner, _ := pterm.DefaultSpinner.Start("Initializing game...")
	time.Sleep(1 * time.Second)
	spinner.Success("Ready to invest! 💰")
}

// ShowLoadingSpinner shows a loading animation with custom text
func ShowLoadingSpinner(text string, duration time.Duration) {
	spinner, _ := pterm.DefaultSpinner.Start(text)
	time.Sleep(duration)
	spinner.Stop()
}

// ShowSuccessMessage displays an animated success message
func ShowSuccessMessage(message string) {
	pterm.Success.Println(message)
}

// ShowErrorMessage displays an animated error message
func ShowErrorMessage(message string) {
	pterm.Error.Println(message)
}

// ShowWarningMessage displays an animated warning message
func ShowWarningMessage(message string) {
	pterm.Warning.Println(message)
}

// ShowInfoMessage displays an animated info message
func ShowInfoMessage(message string) {
	pterm.Info.Println(message)
}

// AnimatedTitle shows an animated title with big letters
func AnimatedTitle(text string) {
	s, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle(text, pterm.NewStyle(pterm.FgCyan))).
		Srender()
	pterm.DefaultCenter.Println(s)
}

// ShowProgressBar displays a progress bar
func ShowProgressBar(title string, total int) {
	progressbar, _ := pterm.DefaultProgressbar.WithTotal(total).WithTitle(title).Start()

	for i := 0; i < total; i++ {
		progressbar.Increment()
		time.Sleep(50 * time.Millisecond)
	}
}

// TypewriterEffect prints text with a typewriter effect
func TypewriterEffect(text string, delay time.Duration) {
	for _, char := range text {
		pterm.Print(string(char))
		time.Sleep(delay)
	}
	pterm.Println()
}

// ShowBigUnicorn displays a big UNICORN title
func ShowBigUnicorn() {
	s, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("UNICORN", pterm.NewStyle(pterm.FgLightMagenta))).
		Srender()
	pterm.DefaultCenter.Println(s)
	time.Sleep(500 * time.Millisecond)
}

// ShowGameStartAnimation displays an animated game start sequence
func ShowGameStartAnimation() {
	// Clear the screen first
	pterm.Print("\033[H\033[2J")

	// Show big UNICORN title
	ShowBigUnicorn()

	// Create an animated box
	pterm.DefaultBox.WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).
		WithTitle("🚀 THE GAME 🚀").
		WithTitleTopCenter().
		Println("Build the next billion-dollar startup\n" +
			"Invest wisely, avoid disasters\n" +
			"Become a legendary VC")

	time.Sleep(300 * time.Millisecond)

	// Show loading animation
	spinner, _ := pterm.DefaultSpinner.
		WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").
		Start("Loading your portfolio...")
	time.Sleep(1200 * time.Millisecond)
	spinner.Success("Ready to make billions! 💰")

	time.Sleep(300 * time.Millisecond)
}

// ShowRoundTransition shows an animated transition between rounds
func ShowRoundTransition(roundNum int) {
	pterm.Print("\n")
	pterm.DefaultHeader.WithFullWidth().
		WithBackgroundStyle(pterm.NewStyle(pterm.BgLightMagenta)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
		Printfln("🎯 ROUND %d 🎯", roundNum)
	time.Sleep(500 * time.Millisecond)
}

// ShowInvestmentAnimation shows an animation when making an investment
func ShowInvestmentAnimation(startupName string, amount int64) {
	spinner, _ := pterm.DefaultSpinner.
		WithSequence("💰", "💸", "💵", "💴", "💶", "💷").
		Start("Processing investment...")
	time.Sleep(800 * time.Millisecond)
	spinner.Success("Invested in " + startupName + "!")
}

// ShowExitAnimation shows an animation for a successful exit
func ShowExitAnimation(startupName string, returns int64) {
	pterm.Print("\n")

	// Create fireworks effect with colors
	colors := []pterm.Color{pterm.FgYellow, pterm.FgGreen, pterm.FgCyan, pterm.FgMagenta}
	fireworks := []string{"✨", "🎆", "🎇", "💫", "⭐", "🌟"}

	for i := 0; i < 3; i++ {
		for _, fw := range fireworks {
			pterm.Println(pterm.NewStyle(colors[i%len(colors)]).Sprint(fw + " " + fw + " " + fw))
		}
	}

	pterm.DefaultBox.WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
		WithTitle("🚀 SUCCESSFUL EXIT! 🚀").
		WithTitleTopCenter().
		Printfln("%s has been acquired!\n\nReturns: $%d", startupName, returns)

	time.Sleep(1000 * time.Millisecond)
}

// ShowGameOverAnimation shows an animated game over screen
func ShowGameOverAnimation(won bool, finalNetWorth int64) {
	pterm.Print("\033[H\033[2J") // Clear screen

	if won {
		// Victory animation
		s, _ := pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromStringWithStyle("VICTORY", pterm.NewStyle(pterm.FgGreen))).
			Srender()
		pterm.DefaultCenter.Println(s)

		pterm.DefaultBox.WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
			WithTitle("🏆 LEGENDARY VC 🏆").
			WithTitleTopCenter().
			Printfln("You've built an empire!\n\nFinal Net Worth: $%d\n\nYou are a unicorn whisperer! 🦄", finalNetWorth)
	} else {
		// Game over animation
		s, _ := pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromStringWithStyle("GAME OVER", pterm.NewStyle(pterm.FgRed))).
			Srender()
		pterm.DefaultCenter.Println(s)

		pterm.DefaultBox.WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			WithTitle("💸 BANKRUPT 💸").
			WithTitleTopCenter().
			Printfln("Your fund has failed\n\nFinal Net Worth: $%d\n\nBetter luck next time!", finalNetWorth)
	}

	time.Sleep(2000 * time.Millisecond)
}

// ShowAchievementUnlock shows an animated achievement unlock
func ShowAchievementUnlock(achievementName, description string) {
	pterm.Print("\n")

	// Flash effect
	for i := 0; i < 3; i++ {
		pterm.Println(pterm.NewStyle(pterm.FgYellow, pterm.Bold).Sprint("⭐ ✨ 🏆 ✨ ⭐"))
		time.Sleep(150 * time.Millisecond)
	}

	pterm.DefaultBox.WithBoxStyle(pterm.NewStyle(pterm.FgYellow)).
		WithTitle("🏆 ACHIEVEMENT UNLOCKED 🏆").
		WithTitleTopCenter().
		Printfln("%s\n\n%s", achievementName, description)

	time.Sleep(1500 * time.Millisecond)
}


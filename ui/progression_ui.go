package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/clear"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/progression"
)

// DisplayLevelUp shows a celebration screen when player levels up
func DisplayLevelUp(playerName string, oldLevel, newLevel int, unlocks []string, pointsEarned int) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	
	// Level up banner
	cyan.Println("\n" + strings.Repeat("=", 70))
	yellow.Println("                    🎉 LEVEL UP! 🎉")
	cyan.Println(strings.Repeat("=", 70))
	
	// Show level progression
	fmt.Println()
	green.Printf("   Level %d ", oldLevel)
	fmt.Print("→ ")
	yellow.Printf("Level %d\n", newLevel)
	
	// Show points earned
	if pointsEarned > 0 {
		fmt.Println()
		magenta.Printf("   💰 Points Earned: +%d points\n", pointsEarned)
	}
	
	// Show new title
	levelInfo := progression.GetLevelInfo(newLevel)
	fmt.Println()
	cyan.Printf("   New Rank: ")
	yellow.Printf("%s\n", levelInfo.Title)
	
	// Show unlocks
	if len(unlocks) > 0 {
		fmt.Println()
		green.Println("   🔓 NEW UNLOCKS:")
		for _, unlock := range unlocks {
			fmt.Printf("      • %s\n", unlock)
		}
	}
	
	// Show next milestone
	nextMilestone := getNextMilestone(newLevel)
	if nextMilestone != "" {
		fmt.Println()
		cyan.Printf("   Next Milestone: ")
		fmt.Printf("%s\n", nextMilestone)
	}
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// DisplayProgressionStats shows current XP and level progress
func DisplayProgressionStats(playerName string) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	
	// Get player profile
	profile, err := database.GetPlayerProfile(playerName)
	if err != nil {
		color.Red("Error loading profile: %v", err)
		return
	}
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                    PLAYER PROGRESSION")
	cyan.Println(strings.Repeat("=", 70))
	
	// Current level and rank
	fmt.Println()
	levelInfo := progression.GetLevelInfo(profile.Level)
	yellow.Printf("   Level: ")
	green.Printf("%d - %s\n", profile.Level, levelInfo.Title)
	
	// XP Progress
	fmt.Println()
	yellow.Printf("   Experience: ")
	fmt.Printf("%d / %d XP\n", profile.ExperiencePoints, profile.NextLevelXP)
	
	// Progress bar
	progressBar := progression.FormatXPBar(profile.ExperiencePoints, profile.NextLevelXP, 40)
	fmt.Printf("   %s %.1f%%\n", progressBar, profile.ProgressPercent)
	
	// Total XP earned
	fmt.Println()
	yellow.Printf("   Total XP Earned: ")
	fmt.Printf("%d XP\n", profile.TotalPointsEarned)
	
	// Next level info
	if profile.Level < 50 {
		fmt.Println()
		nextLevelInfo := progression.GetLevelInfo(profile.Level + 1)
		cyan.Println("   NEXT LEVEL:")
		fmt.Printf("      Level %d: %s\n", nextLevelInfo.Level, nextLevelInfo.Title)
		
		xpNeeded := profile.NextLevelXP - profile.ExperiencePoints
		fmt.Printf("      XP needed: %d\n", xpNeeded)
		
		if len(nextLevelInfo.Unlocks) > 0 {
			fmt.Println("      Unlocks:")
			for _, unlock := range nextLevelInfo.Unlocks {
				fmt.Printf("         • %s\n", unlock)
			}
		}
	} else {
		fmt.Println()
		green.Println("   🌟 MAX LEVEL REACHED! 🌟")
	}
	
	// Show upcoming milestones
	fmt.Println()
	cyan.Println("   UPCOMING MILESTONES:")
	milestones := getUpcomingMilestones(profile.Level)
	if len(milestones) > 0 {
		for level, description := range milestones {
			fmt.Printf("      Level %d: %s\n", level, description)
		}
	} else {
		fmt.Println("      All milestones achieved!")
	}
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// DisplayXPGained shows XP breakdown after a game
func DisplayXPGained(xpBreakdown map[string]int, totalXP int) {
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen)
	
	fmt.Println()
	yellow.Println("   📊 EXPERIENCE EARNED:")
	
	for source, amount := range xpBreakdown {
		if amount > 0 {
			green.Printf("      +%d XP ", amount)
			fmt.Printf("- %s\n", source)
		}
	}
	
	fmt.Println()
	yellow.Printf("   Total XP Gained: ")
	green.Printf("+%d XP\n", totalXP)
}

// DisplayLevelUnlocks shows what unlocks at different levels (preview)
func DisplayLevelUnlocks() {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                    LEVEL UNLOCK PREVIEW")
	cyan.Println(strings.Repeat("=", 70))
	
	fmt.Println()
	
	// Show major milestones
	milestones := []int{1, 5, 10, 15, 20, 25, 30, 50}
	
	for _, level := range milestones {
		unlocks := progression.GetUnlocksForLevel(level)
		if len(unlocks) > 0 {
			levelInfo := progression.GetLevelInfo(level)
			
			yellow.Printf("\n   Level %d - %s:\n", level, levelInfo.Title)
			for _, unlock := range unlocks {
				fmt.Printf("      %s %s\n", unlock.Icon, unlock.Name)
				fmt.Printf("         %s\n", unlock.Description)
			}
		}
	}
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Helper functions

func getNextMilestone(currentLevel int) string {
	milestones := map[int]string{
		5:  "Unlock Hard Difficulty",
		10: "Unlock Expert Difficulty & Advanced Analytics",
		15: "Unlock Nightmare Mode & Achievement Chains",
		20: "Unlock Game Modifiers & Secondary Markets",
		25: "Unlock Prestige System & Legendary Achievements",
		30: "Unlock Master Difficulty",
		50: "Achieve Titan Status",
	}
	
	for level := currentLevel + 1; level <= 50; level++ {
		if milestone, exists := milestones[level]; exists {
			return fmt.Sprintf("Level %d - %s", level, milestone)
		}
	}
	
	return ""
}

func getUpcomingMilestones(currentLevel int) map[int]string {
	allMilestones := map[int]string{
		5:  "🔥 Hard Difficulty",
		10: "💎 Expert Difficulty + Analytics",
		15: "👹 Nightmare Mode + Achievement Chains",
		20: "⚙️ Game Modifiers + Secondary Markets",
		25: "✨ Prestige System + Legendary Achievements",
		30: "👑 Master Difficulty",
		50: "🌟 Titan Status",
	}
	
	upcoming := make(map[int]string)
	count := 0
	
	for level := currentLevel + 1; level <= 50 && count < 3; level++ {
		if milestone, exists := allMilestones[level]; exists {
			upcoming[level] = milestone
			count++
		}
	}
	
	return upcoming
}

// DisplayLevelProgress shows a compact level indicator for the main menu
func DisplayLevelProgress(playerName string) {
	profile, err := database.GetPlayerProfile(playerName)
	if err != nil {
		return
	}
	
	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)
	
	levelInfo := progression.GetLevelInfo(profile.Level)
	cyan.Printf("   Level %d - %s", profile.Level, levelInfo.Title)
	
	if profile.Level < 50 {
		progressBar := progression.FormatXPBar(profile.ExperiencePoints, profile.NextLevelXP, 20)
		fmt.Printf("\n   %s ", progressBar)
		yellow.Printf("%d/%d XP", profile.ExperiencePoints, profile.NextLevelXP)
	} else {
		yellow.Print(" [MAX]")
	}
	
	fmt.Println()
}


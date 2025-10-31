package ascii

import "strings"

// ASCII art icons for the game (text-based for terminal compatibility)
const (
	Money     = "$"
	Calendar  = "[📅]"
	Portfolio = "📊"
	News      = "📰"
	Trophy    = "🏆"
	Star      = "★"
	Fire      = "🔥"
	Target    = "◎"
	Rocket    = "🚀"
	Chart     = "📈"
	Medal     = "🥇"
	Silver    = "🥈"
	Bronze    = "🥉"
	Check     = "✓"
	Warning   = "!"
	GameOver  = "🏁"
	Invest    = "$"
	Exit      = "→"
	Level     = "↑"
	Coin      = "●"
	Bank      = "🏦"
	Building  = "🏢"
	Lightbulb = "💡"
	Shield    = "🛡"
	Crown     = "👑"
	Diamond   = "◆"
	Zap       = "⚡"
	Star2     = "✦"
)

// News header ASCII art
const NewsHeader = `
╔══════════════════════════════════════════════════════════╗
║                    📰 COMPANY NEWS 📰                     ║
╚══════════════════════════════════════════════════════════╝
`

// Investment phase header
const InvestmentHeader = `
╔══════════════════════════════════════════════════════════╗
║                  💵 INVESTMENT PHASE 💵                   ║
╚══════════════════════════════════════════════════════════╝
`

// Portfolio header
const PortfolioHeader = `
╔══════════════════════════════════════════════════════════╗
║                    📊 YOUR PORTFOLIO 📊                   ║
╚══════════════════════════════════════════════════════════╝
`

// Game over header
const GameOverHeader = `
╔══════════════════════════════════════════════════════════╗
║              🏁 GAME OVER - FINAL RESULTS 🏁              ║
╚══════════════════════════════════════════════════════════╝
`

// Leaderboard header
const LeaderboardHeader = `
╔══════════════════════════════════════════════════════════╗
║                  🏆 LEADERBOARDS 🏆                       ║
╚══════════════════════════════════════════════════════════╝
`

// Achievements header
const AchievementsHeader = `
╔══════════════════════════════════════════════════════════╗
║                ⭐ ACHIEVEMENTS ⭐                         ║
╚══════════════════════════════════════════════════════════╝
`

// ASCII art for section headers
func GetSectionHeader(title string) string {
	return "\n" + strings.Repeat("═", 60) + "\n" +
		centerText(title, 60) + "\n" +
		strings.Repeat("═", 60)
}

func centerText(text string, width int) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}
	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text
}

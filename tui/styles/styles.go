package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Color palette - matching existing game theme
var (
	// Primary colors
	Cyan       = lipgloss.Color("#00FFFF")
	Green      = lipgloss.Color("#00FF00")
	Yellow     = lipgloss.Color("#FFFF00")
	Magenta    = lipgloss.Color("#FF00FF")
	Red        = lipgloss.Color("#FF0000")
	White      = lipgloss.Color("#FFFFFF")
	Gray       = lipgloss.Color("#808080")
	DarkGray   = lipgloss.Color("#404040")
	Black      = lipgloss.Color("#000000")

	// Accent colors
	Gold       = lipgloss.Color("#FFD700")
	Orange     = lipgloss.Color("#FFA500")
	LightCyan  = lipgloss.Color("#E0FFFF")
	LightGreen = lipgloss.Color("#90EE90")
	DimWhite   = lipgloss.Color("#A0A0A0")
)

// Base styles
var (
	// App container
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Italic(true)

	// Header bar
	HeaderStyle = lipgloss.NewStyle().
			Foreground(Black).
			Background(Cyan).
			Bold(true).
			Padding(0, 2).
			Width(70)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Cyan).
			Padding(1, 2)

	FocusedBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Magenta).
			Padding(1, 2)

	// Panel styles for split layouts
	LeftPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Cyan).
			Padding(0, 1).
			Width(35)

	RightPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Cyan).
			Padding(0, 1).
			Width(35)
)

// Menu styles
var (
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(White).
			Padding(0, 2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(Black).
				Background(Cyan).
				Bold(true).
				Padding(0, 2)

	MenuDescriptionStyle = lipgloss.NewStyle().
				Foreground(Gray).
				Italic(true).
				PaddingLeft(4)
)

// Status and info styles
var (
	// Status bar at bottom
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(DimWhite).
			Background(DarkGray).
			Padding(0, 1)

	// Help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(Gray)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Gray)
)

// Game-specific styles
var (
	// Money/currency
	MoneyPositiveStyle = lipgloss.NewStyle().
				Foreground(Green).
				Bold(true)

	MoneyNegativeStyle = lipgloss.NewStyle().
				Foreground(Red).
				Bold(true)

	MoneyNeutralStyle = lipgloss.NewStyle().
				Foreground(Yellow)

	// Risk indicators
	RiskLowStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	RiskMediumStyle = lipgloss.NewStyle().
			Foreground(Yellow)

	RiskHighStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)

	// Growth indicators
	GrowthHighStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	GrowthMediumStyle = lipgloss.NewStyle().
				Foreground(Yellow)

	GrowthLowStyle = lipgloss.NewStyle().
			Foreground(Red)
)

// News and events
var (
	NewsHeaderStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(Cyan)

	NewsGoodStyle = lipgloss.NewStyle().
			Foreground(Green)

	NewsBadStyle = lipgloss.NewStyle().
			Foreground(Red)

	NewsNeutralStyle = lipgloss.NewStyle().
				Foreground(Yellow)
)

// Dialog styles
var (
	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(Magenta).
			Padding(1, 2).
			Width(50)

	DialogTitleStyle = lipgloss.NewStyle().
				Foreground(Magenta).
				Bold(true).
				Align(lipgloss.Center)

	DialogButtonStyle = lipgloss.NewStyle().
				Foreground(White).
				Background(DarkGray).
				Padding(0, 2).
				Margin(0, 1)

	DialogButtonActiveStyle = lipgloss.NewStyle().
				Foreground(Black).
				Background(Cyan).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)
)

// Table styles
var (
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(Cyan).
				Bold(true).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(Cyan)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(White)

	TableSelectedRowStyle = lipgloss.NewStyle().
				Foreground(Black).
				Background(Cyan).
				Bold(true)

	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)
)

// Input styles
var (
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Cyan).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(Magenta).
				Padding(0, 1)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(Cyan)

	InputPlaceholderStyle = lipgloss.NewStyle().
				Foreground(Gray)
)

// Achievement and progression
var (
	AchievementStyle = lipgloss.NewStyle().
				Foreground(Gold).
				Bold(true)

	AchievementLockedStyle = lipgloss.NewStyle().
				Foreground(Gray)

	LevelStyle = lipgloss.NewStyle().
			Foreground(Magenta).
			Bold(true)

	XPStyle = lipgloss.NewStyle().
		Foreground(LightGreen)
)

// Leaderboard styles
var (
	LeaderboardFirstStyle = lipgloss.NewStyle().
				Foreground(Gold).
				Bold(true)

	LeaderboardSecondStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#C0C0C0")).
				Bold(true)

	LeaderboardThirdStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CD7F32")).
				Bold(true)

	LeaderboardPlayerStyle = lipgloss.NewStyle().
				Foreground(Cyan).
				Bold(true)
)

// Spinner style
var (
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(Cyan)
)

// Logo/splash styles
var (
	LogoStyle = lipgloss.NewStyle().
			Foreground(Magenta).
			Bold(true).
			Align(lipgloss.Center)

	SplashBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(Cyan).
			Padding(2, 4).
			Align(lipgloss.Center)
)

// Helper functions

// FormatMoney returns styled money string based on value
func FormatMoney(amount int64) string {
	if amount > 0 {
		return MoneyPositiveStyle.Render(formatCurrency(amount))
	} else if amount < 0 {
		return MoneyNegativeStyle.Render(formatCurrency(amount))
	}
	return MoneyNeutralStyle.Render(formatCurrency(amount))
}

// FormatProfit returns styled profit/loss string
func FormatProfit(amount int64) string {
	if amount > 0 {
		return MoneyPositiveStyle.Render("+" + formatCurrency(amount))
	} else if amount < 0 {
		return MoneyNegativeStyle.Render(formatCurrency(amount))
	}
	return MoneyNeutralStyle.Render(formatCurrency(amount))
}

// formatCurrency formats an int64 as currency with commas
func formatCurrency(amount int64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}

	str := ""
	for amount > 0 {
		if str != "" {
			str = "," + str
		}
		if amount >= 1000 {
			str = fmt.Sprintf("%03d", amount%1000) + str
		} else {
			str = fmt.Sprintf("%d", amount%1000) + str
		}
		amount /= 1000
	}

	if str == "" {
		str = "0"
	}

	if negative {
		return "-$" + str
	}
	return "$" + str
}

// RiskStyle returns appropriate style for risk level
func RiskStyle(riskScore float64) lipgloss.Style {
	if riskScore > 0.6 {
		return RiskHighStyle
	} else if riskScore > 0.4 {
		return RiskMediumStyle
	}
	return RiskLowStyle
}

// GrowthStyle returns appropriate style for growth potential
func GrowthStyle(growthPotential float64) lipgloss.Style {
	if growthPotential > 0.6 {
		return GrowthHighStyle
	} else if growthPotential > 0.4 {
		return GrowthMediumStyle
	}
	return GrowthLowStyle
}

// LeaderboardRankStyle returns style based on rank
func LeaderboardRankStyle(rank int, isPlayer bool) lipgloss.Style {
	if isPlayer {
		return LeaderboardPlayerStyle
	}
	switch rank {
	case 1:
		return LeaderboardFirstStyle
	case 2:
		return LeaderboardSecondStyle
	case 3:
		return LeaderboardThirdStyle
	default:
		return TableRowStyle
	}
}

// Width helpers for responsive layouts
func SetWidth(style lipgloss.Style, width int) lipgloss.Style {
	return style.Width(width)
}

func SetHeight(style lipgloss.Style, height int) lipgloss.Style {
	return style.Height(height)
}

// Center text helper
func Center(width int, s string) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(s)
}

// Horizontal join with gap
func JoinHorizontal(gap int, strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

// Vertical join
func JoinVertical(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, strs...)
}

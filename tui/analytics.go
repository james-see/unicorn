package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/analytics"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// AnalyticsScreen shows analytics dashboard
type AnalyticsScreen struct {
	width       int
	height      int
	playerName  string
	needsName   bool
	nameInput   textinput.Model
	stats       *database.PlayerStats
	trendReport *analytics.TrendReport
	monthlyData []*analytics.MonthlyReport
	selectedTab int
	tabs        []string
}

// NewAnalyticsScreen creates a new analytics screen
func NewAnalyticsScreen(width, height int) *AnalyticsScreen {
	ti := textinput.New()
	ti.Placeholder = "Enter player name"
	ti.CharLimit = 30
	ti.Width = 30
	ti.Focus()

	return &AnalyticsScreen{
		width:     width,
		height:    height,
		needsName: true,
		nameInput: ti,
		tabs:      []string{"Overview", "Heatmap", "Difficulty"},
	}
}

func (s *AnalyticsScreen) loadData() {
	// Load player stats
	stats, err := database.GetPlayerStats(s.playerName)
	if err == nil && stats.TotalGames > 0 {
		s.stats = stats
	}

	// Load trend report
	trendReport, err := analytics.GenerateTrendAnalysis(s.playerName, 30)
	if err == nil {
		s.trendReport = trendReport
	}

	// Load monthly data for heatmap (last 6 months)
	now := time.Now()
	for i := 5; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		monthReport, err := analytics.GetMonthlyStats(s.playerName, month.Year(), int(month.Month()))
		if err == nil && monthReport != nil {
			s.monthlyData = append(s.monthlyData, monthReport)
		}
	}
}

// Init initializes the analytics screen
func (s *AnalyticsScreen) Init() tea.Cmd {
	if s.needsName {
		return textinput.Blink
	}
	return nil
}

// Update handles analytics input
func (s *AnalyticsScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	if s.needsName {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				name := strings.TrimSpace(s.nameInput.Value())
				if name != "" {
					s.playerName = name
					s.needsName = false
					s.loadData()
				}
				return s, nil
			case "esc":
				return s, PopScreen()
			}
		}
		var cmd tea.Cmd
		s.nameInput, cmd = s.nameInput.Update(msg)
		return s, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Global.Back):
			return s, PopScreen()
		case key.Matches(msg, keys.Global.Left):
			if s.selectedTab > 0 {
				s.selectedTab--
			}
		case key.Matches(msg, keys.Global.Right):
			if s.selectedTab < len(s.tabs)-1 {
				s.selectedTab++
			}
		}
	}
	return s, nil
}

// View renders analytics
func (s *AnalyticsScreen) View() string {
	if s.needsName {
		return s.renderNameInput()
	}

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📊 ANALYTICS DASHBOARD 📊")))
	b.WriteString("\n\n")

	// Tabs
	var tabs []string
	for i, tab := range s.tabs {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == s.selectedTab {
			style = style.Foreground(styles.Black).Background(styles.Cyan).Bold(true)
		} else {
			style = style.Foreground(styles.Gray)
		}
		tabs = append(tabs, style.Render(tab))
	}
	tabRow := strings.Join(tabs, " ")
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(tabRow))
	b.WriteString("\n\n")

	// Content based on selected tab
	switch s.tabs[s.selectedTab] {
	case "Overview":
		b.WriteString(s.renderOverview())
	case "Heatmap":
		b.WriteString(s.renderHeatmap())
	case "Difficulty":
		b.WriteString(s.renderDifficultyBreakdown())
	}

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("←/→ switch tabs • esc back"))

	return b.String()
}

func (s *AnalyticsScreen) renderNameInput() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📊 ANALYTICS DASHBOARD 📊")))
	b.WriteString("\n\n")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(40)

	var content strings.Builder
	content.WriteString(lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true).Render("Enter Player Name"))
	content.WriteString("\n\n")
	content.WriteString(s.nameInput.View())

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(boxStyle.Render(content.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter submit • esc back"))

	return b.String()
}

func (s *AnalyticsScreen) renderOverview() string {
	var b strings.Builder

	contentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(65)

	var content strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(styles.White)

	content.WriteString(titleStyle.Render(fmt.Sprintf("Player: %s", s.playerName)))
	content.WriteString("\n\n")

	if s.stats == nil || s.stats.TotalGames == 0 {
		content.WriteString(labelStyle.Render("No game data available. Play some games first!"))
	} else {
		// Performance Overview
		content.WriteString(titleStyle.Render("PERFORMANCE OVERVIEW"))
		content.WriteString("\n\n")

		content.WriteString(labelStyle.Render("Total Games:     "))
		content.WriteString(valueStyle.Render(fmt.Sprintf("%d", s.stats.TotalGames)))
		content.WriteString("\n")

		content.WriteString(labelStyle.Render("Win Rate:        "))
		content.WriteString(valueStyle.Render(fmt.Sprintf("%.1f%%", s.stats.WinRate)))
		content.WriteString("\n")

		content.WriteString(labelStyle.Render("Best Net Worth:  "))
		content.WriteString(valueStyle.Render(fmt.Sprintf("$%s", formatCompactMoney(s.stats.BestNetWorth))))
		content.WriteString("\n")

		content.WriteString(labelStyle.Render("Best ROI:        "))
		content.WriteString(valueStyle.Render(fmt.Sprintf("%.1f%%", s.stats.BestROI)))
		content.WriteString("\n")

		content.WriteString(labelStyle.Render("Total Exits:     "))
		content.WriteString(valueStyle.Render(fmt.Sprintf("%d", s.stats.TotalExits)))
		content.WriteString("\n\n")

		// Trend info
		if s.trendReport != nil {
			content.WriteString(titleStyle.Render("TREND ANALYSIS"))
			content.WriteString("\n\n")

			trendStyle := lipgloss.NewStyle().Foreground(styles.Cyan)
			if strings.Contains(s.trendReport.TrendVector, "Improving") {
				trendStyle = trendStyle.Foreground(styles.Green)
			} else if strings.Contains(s.trendReport.TrendVector, "Declining") {
				trendStyle = trendStyle.Foreground(styles.Red)
			}

			content.WriteString(labelStyle.Render("Current Trend: "))
			content.WriteString(trendStyle.Render(s.trendReport.TrendVector))
			content.WriteString("\n")

			if s.trendReport.Last7Days.GamesPlayed > 0 {
				content.WriteString(labelStyle.Render("Last 7 Days:   "))
				content.WriteString(valueStyle.Render(fmt.Sprintf("%d games, %.0f%% win rate",
					s.trendReport.Last7Days.GamesPlayed, s.trendReport.Last7Days.WinRate)))
				content.WriteString("\n")
			}
		}
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(contentBox.Render(content.String())))
	b.WriteString("\n\n")

	return b.String()
}

func (s *AnalyticsScreen) renderHeatmap() string {
	var b strings.Builder

	contentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2).
		Width(70)

	var content strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)

	content.WriteString(titleStyle.Render("📈 PERFORMANCE HEATMAP (Last 6 Months)"))
	content.WriteString("\n\n")

	if len(s.monthlyData) == 0 {
		content.WriteString("No monthly data available yet.\n")
	} else {
		// Month headers
		headerStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
		for _, month := range s.monthlyData {
			content.WriteString(headerStyle.Render(fmt.Sprintf(" %-8s", month.Month[:3])))
		}
		content.WriteString("\n")

		// Games played row
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Yellow).Render("Games: "))
		for _, month := range s.monthlyData {
			cellStyle := s.getHeatmapStyle(month.GamesPlayed, 10)
			content.WriteString(cellStyle.Render(fmt.Sprintf(" %-8d", month.GamesPlayed)))
		}
		content.WriteString("\n")

		// Wins row
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Yellow).Render("Wins:  "))
		for _, month := range s.monthlyData {
			cellStyle := s.getHeatmapStyle(month.Wins, 5)
			content.WriteString(cellStyle.Render(fmt.Sprintf(" %-8d", month.Wins)))
		}
		content.WriteString("\n")

		// Avg ROI row
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Yellow).Render("Avg ROI:"))
		for _, month := range s.monthlyData {
			roiStyle := lipgloss.NewStyle()
			if month.AvgROI > 50 {
				roiStyle = roiStyle.Foreground(styles.Green).Bold(true)
			} else if month.AvgROI > 0 {
				roiStyle = roiStyle.Foreground(styles.Green)
			} else if month.AvgROI > -50 {
				roiStyle = roiStyle.Foreground(styles.Yellow)
			} else {
				roiStyle = roiStyle.Foreground(styles.Red)
			}
			content.WriteString(roiStyle.Render(fmt.Sprintf(" %-8.0f%%", month.AvgROI)))
		}
		content.WriteString("\n\n")

		// Visual heatmap grid
		content.WriteString(titleStyle.Render("Activity Grid:"))
		content.WriteString("\n")
		for _, month := range s.monthlyData {
			intensity := s.getActivityIntensity(month.GamesPlayed)
			content.WriteString(intensity)
		}
		content.WriteString("\n\n")

		// Legend
		legendStyle := lipgloss.NewStyle().Foreground(styles.Gray)
		content.WriteString(legendStyle.Render("Legend: "))
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Gray).Render("░ "))
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Green).Render("▒ "))
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Cyan).Render("▓ "))
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true).Render("█ "))
		content.WriteString(legendStyle.Render("(Low → High activity)"))
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(contentBox.Render(content.String())))
	b.WriteString("\n\n")

	return b.String()
}

func (s *AnalyticsScreen) getHeatmapStyle(value, threshold int) lipgloss.Style {
	if value == 0 {
		return lipgloss.NewStyle().Foreground(styles.Gray)
	} else if value < threshold/2 {
		return lipgloss.NewStyle().Foreground(styles.Yellow)
	} else if value < threshold {
		return lipgloss.NewStyle().Foreground(styles.Green)
	}
	return lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
}

func (s *AnalyticsScreen) getActivityIntensity(games int) string {
	if games == 0 {
		return lipgloss.NewStyle().Foreground(styles.Gray).Render(" ░░░░░░░░ ")
	} else if games < 3 {
		return lipgloss.NewStyle().Foreground(styles.Green).Render(" ▒▒▒▒░░░░ ")
	} else if games < 6 {
		return lipgloss.NewStyle().Foreground(styles.Cyan).Render(" ▓▓▓▓▒▒░░ ")
	} else if games < 10 {
		return lipgloss.NewStyle().Foreground(styles.Magenta).Render(" █▓▓▓▓▒▒░ ")
	}
	return lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true).Render(" ████▓▓▓▒ ")
}

func (s *AnalyticsScreen) renderDifficultyBreakdown() string {
	var b strings.Builder

	contentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2).
		Width(65)

	var content strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)

	content.WriteString(titleStyle.Render("📊 DIFFICULTY BREAKDOWN"))
	content.WriteString("\n\n")

	difficulties := []string{"Easy", "Medium", "Hard", "Expert"}
	headerStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-10s %8s %8s %12s", "DIFF", "GAMES", "WINS", "AVG NET")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 45))
	content.WriteString("\n")

	for _, diff := range difficulties {
		scores, err := database.GetTopScoresByNetWorth(1000, diff)
		if err != nil {
			continue
		}

		// Filter for this player
		playerGames := 0
		playerWins := 0
		totalNet := int64(0)

		for _, score := range scores {
			if score.PlayerName == s.playerName {
				playerGames++
				if score.ROI > 0 {
					playerWins++
				}
				totalNet += score.FinalNetWorth
			}
		}

		if playerGames == 0 {
			grayStyle := lipgloss.NewStyle().Foreground(styles.Gray)
			content.WriteString(grayStyle.Render(fmt.Sprintf("%-10s %8s %8s %12s", diff, "-", "-", "-")))
		} else {
			avgNet := totalNet / int64(playerGames)
			diffStyle := lipgloss.NewStyle().Foreground(styles.White)
			valueStyle := lipgloss.NewStyle().Foreground(styles.Green)

			content.WriteString(diffStyle.Render(fmt.Sprintf("%-10s ", diff)))
			content.WriteString(valueStyle.Render(fmt.Sprintf("%8d %8d %12s",
				playerGames, playerWins, "$"+formatCompactMoney(avgNet))))
		}
		content.WriteString("\n")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(contentBox.Render(content.String())))
	b.WriteString("\n\n")

	return b.String()
}

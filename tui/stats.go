package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// StatsScreen shows player statistics
type StatsScreen struct {
	width       int
	height      int
	nameInput   textinput.Model
	stats       *database.PlayerStats
	playerName  string
	inputMode   bool
}

// NewStatsScreen creates a new stats screen
func NewStatsScreen(width, height int) *StatsScreen {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter player name"
	nameInput.Focus()
	nameInput.CharLimit = 30
	nameInput.Width = 25
	
	return &StatsScreen{
		width:     width,
		height:    height,
		nameInput: nameInput,
		inputMode: true,
	}
}

// Init initializes the stats screen
func (s *StatsScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles stats input
func (s *StatsScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Global.Back) {
			if s.inputMode {
				return s, PopScreen()
			}
			s.inputMode = true
			s.stats = nil
			return s, textinput.Blink
		}
		
		if s.inputMode && msg.Type == tea.KeyEnter {
			name := strings.TrimSpace(s.nameInput.Value())
			if name != "" {
				s.playerName = name
				stats, err := database.GetPlayerStats(name)
				if err == nil {
					s.stats = stats
				}
				s.inputMode = false
			}
			return s, nil
		}
	}
	
	if s.inputMode {
		var cmd tea.Cmd
		s.nameInput, cmd = s.nameInput.Update(msg)
		return s, cmd
	}
	
	return s, nil
}

// View renders stats
func (s *StatsScreen) View() string {
	var b strings.Builder
	
	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)
	
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📊 PLAYER STATISTICS 📊")))
	b.WriteString("\n\n")
	
	if s.inputMode {
		// Name input
		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Cyan).
			Padding(1, 2).
			Width(40)
		
		var input strings.Builder
		input.WriteString("Enter player name:\n\n")
		input.WriteString(s.nameInput.View())
		
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render(input.String())))
	} else {
		// Stats display
		statsBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Cyan).
			Padding(1, 2).
			Width(50)
		
		var stats strings.Builder
		titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
		stats.WriteString(titleStyle.Render(fmt.Sprintf("Stats for %s", s.playerName)))
		stats.WriteString("\n\n")
		
		if s.stats == nil || s.stats.TotalGames == 0 {
			stats.WriteString("No games played yet")
		} else {
			labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
			
			stats.WriteString(labelStyle.Render("Total Games: "))
			stats.WriteString(fmt.Sprintf("%d\n", s.stats.TotalGames))
			
			stats.WriteString(labelStyle.Render("Best Net Worth: "))
			stats.WriteString(fmt.Sprintf("$%s\n", formatCompactMoney(s.stats.BestNetWorth)))
			
			stats.WriteString(labelStyle.Render("Best ROI: "))
			roiStyle := lipgloss.NewStyle().Foreground(styles.Green)
			stats.WriteString(roiStyle.Render(fmt.Sprintf("%.1f%%\n", s.stats.BestROI*100)))
			
			stats.WriteString(labelStyle.Render("Total Exits: "))
			stats.WriteString(fmt.Sprintf("%d\n", s.stats.TotalExits))
			
			stats.WriteString(labelStyle.Render("Average Net Worth: "))
			stats.WriteString(fmt.Sprintf("$%s\n", formatCompactMoney(int64(s.stats.AverageNetWorth))))
			
			stats.WriteString(labelStyle.Render("Win Rate: "))
			stats.WriteString(fmt.Sprintf("%.1f%%\n", s.stats.WinRate))
		}
		
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(statsBox.Render(stats.String())))
	}
	
	b.WriteString("\n\n")
	
	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))
	
	return b.String()
}

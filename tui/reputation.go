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

// ReputationScreen displays VC reputation
type ReputationScreen struct {
	width      int
	height     int
	playerName string
	reputation *database.VCReputation
	nameInput  textinput.Model
	inputMode  bool
}

// NewReputationScreen creates a new reputation screen
func NewReputationScreen(width, height int, playerName string) *ReputationScreen {
	ti := textinput.New()
	ti.Placeholder = "Enter player name"
	ti.CharLimit = 30
	ti.Width = 30

	s := &ReputationScreen{
		width:      width,
		height:     height,
		playerName: playerName,
		nameInput:  ti,
		inputMode:  playerName == "",
	}

	if playerName == "" {
		ti.Focus()
	} else {
		s.loadReputation()
	}

	return s
}

func (s *ReputationScreen) loadReputation() {
	rep, err := database.GetVCReputation(s.playerName)
	if err == nil {
		s.reputation = rep
	}
}

// Init initializes the reputation screen
func (s *ReputationScreen) Init() tea.Cmd {
	if s.inputMode {
		return textinput.Blink
	}
	return nil
}

// Update handles reputation screen input
func (s *ReputationScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.inputMode {
			switch msg.String() {
			case "enter":
				name := strings.TrimSpace(s.nameInput.Value())
				if name != "" {
					s.playerName = name
					s.inputMode = false
					s.loadReputation()
				}
				return s, nil
			case "esc":
				return s, PopScreen()
			}
			var cmd tea.Cmd
			s.nameInput, cmd = s.nameInput.Update(msg)
			return s, cmd
		}

		if key.Matches(msg, keys.Global.Back) {
			return s, PopScreen()
		}
	}

	return s, nil
}

// View renders the reputation screen
func (s *ReputationScreen) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Gold).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⭐ VC REPUTATION ⭐")))
	b.WriteString("\n\n")

	if s.inputMode {
		return s.renderNameInput()
	}

	if s.reputation == nil {
		infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
		b.WriteString(infoStyle.Render("No reputation data found for this player."))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center).Render("esc back"))
		return b.String()
	}

	rep := s.reputation

	// Calculate aggregate score
	aggregate := (rep.PerformanceScore * 0.4) + (rep.FounderScore * 0.3) + (rep.MarketScore * 0.3)
	level := getReputationLevel(aggregate)
	tier := getDealQualityTier(aggregate)

	// Main content box
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Gold).
		Padding(1, 2).
		Width(60)

	var content strings.Builder

	// Overall reputation
	titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)

	content.WriteString(titleStyle.Render(fmt.Sprintf("Player: %s", s.playerName)))
	content.WriteString("\n\n")

	content.WriteString(valueStyle.Render(fmt.Sprintf("Overall Reputation: %.1f/100 - %s", aggregate, level)))
	content.WriteString("\n\n")

	// Component scores with bars
	content.WriteString(titleStyle.Render("REPUTATION COMPONENTS"))
	content.WriteString("\n\n")

	content.WriteString(s.renderScoreBar("Performance", rep.PerformanceScore))
	content.WriteString(s.renderScoreBar("Founder Relations", rep.FounderScore))
	content.WriteString(s.renderScoreBar("Market Standing", rep.MarketScore))

	content.WriteString("\n")
	content.WriteString(titleStyle.Render("CAREER STATS"))
	content.WriteString("\n\n")

	statStyle := lipgloss.NewStyle().Foreground(styles.White)
	content.WriteString(statStyle.Render(fmt.Sprintf("  Games Played:     %d", rep.TotalGamesPlayed)))
	content.WriteString("\n")
	content.WriteString(statStyle.Render(fmt.Sprintf("  Successful Exits: %d", rep.SuccessfulExits)))
	content.WriteString("\n")
	content.WriteString(statStyle.Render(fmt.Sprintf("  Avg ROI (Last 5): %.1f%%", rep.AvgROILast5)))
	content.WriteString("\n\n")

	// Deal flow quality
	content.WriteString(titleStyle.Render("DEAL FLOW QUALITY"))
	content.WriteString("\n\n")

	tierStyle := lipgloss.NewStyle().Bold(true)
	if aggregate >= 70 {
		tierStyle = tierStyle.Foreground(styles.Green)
	} else if aggregate >= 40 {
		tierStyle = tierStyle.Foreground(styles.Yellow)
	} else {
		tierStyle = tierStyle.Foreground(styles.Red)
	}

	content.WriteString(tierStyle.Render(fmt.Sprintf("  %s", tier)))
	content.WriteString("\n")
	content.WriteString(statStyle.Render(fmt.Sprintf("  %s", getDealQualityDescription(aggregate))))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(contentStyle.Render(content.String())))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))

	return b.String()
}

func (s *ReputationScreen) renderNameInput() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Gold).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⭐ VC REPUTATION ⭐")))
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

func (s *ReputationScreen) renderScoreBar(label string, score float64) string {
	bars := int(score / 10)
	barStr := ""
	for i := 0; i < 10; i++ {
		if i < bars {
			barStr += "█"
		} else {
			barStr += "░"
		}
	}

	labelStyle := lipgloss.NewStyle().Foreground(styles.White).Width(18)
	scoreStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	barStyle := lipgloss.NewStyle().Foreground(styles.Green)

	return fmt.Sprintf("  %s %s %s\n",
		labelStyle.Render(label+":"),
		scoreStyle.Render(fmt.Sprintf("%5.1f/100", score)),
		barStyle.Render("["+barStr+"]"))
}

// Helper functions
func getReputationLevel(score float64) string {
	if score >= 90 {
		return "Legendary VC"
	} else if score >= 80 {
		return "Top-Tier VC"
	} else if score >= 70 {
		return "Established VC"
	} else if score >= 60 {
		return "Rising VC"
	} else if score >= 50 {
		return "Competent VC"
	} else if score >= 40 {
		return "Developing VC"
	}
	return "Emerging VC"
}

func getDealQualityTier(score float64) string {
	if score >= 70 {
		return "Tier 1 (Hot Deals) 🔥"
	} else if score >= 40 {
		return "Tier 2 (Standard Deals)"
	}
	return "Tier 3 (Struggling Deals)"
}

func getDealQualityDescription(score float64) string {
	if score >= 70 {
		return "Access to high-quality startups with lower risk"
	} else if score >= 40 {
		return "Access to standard startup opportunities"
	}
	return "Limited to higher-risk startups"
}

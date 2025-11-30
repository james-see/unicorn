package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/achievements"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// AchievementsScreen shows player achievements
type AchievementsScreen struct {
	width            int
	height           int
	playerName       string
	unlockedAchs     []string
	totalPoints      int
	scrollOffset     int
	selectedCategory int
	categories       []string

	// Name entry
	needsName bool
	nameInput textinput.Model
}

// NewAchievementsScreen creates a new achievements screen
func NewAchievementsScreen(width, height int, playerName string) *AchievementsScreen {
	// If no player name, show name input first
	if playerName == "" {
		nameInput := textinput.New()
		nameInput.Placeholder = "Enter your name"
		nameInput.Focus()
		nameInput.CharLimit = 30
		nameInput.Width = 30

		return &AchievementsScreen{
			width:     width,
			height:    height,
			needsName: true,
			nameInput: nameInput,
		}
	}

	// Load unlocked achievements
	unlocked, _ := database.GetPlayerAchievements(playerName)

	// Calculate total points
	totalPoints := 0
	for _, id := range unlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			totalPoints += ach.Points
		}
	}

	// Get categories
	categories := []string{"All", "Investing", "Performance", "Events", "Meta"}

	return &AchievementsScreen{
		width:        width,
		height:       height,
		playerName:   playerName,
		unlockedAchs: unlocked,
		totalPoints:  totalPoints,
		categories:   categories,
	}
}

// Init initializes the achievements screen
func (s *AchievementsScreen) Init() tea.Cmd {
	if s.needsName {
		return textinput.Blink
	}
	return nil
}

// Update handles achievements input
func (s *AchievementsScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	// Handle name input state
	if s.needsName {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Global.Back):
				return s, PopScreen()
			case msg.Type == tea.KeyEnter:
				name := strings.TrimSpace(s.nameInput.Value())
				if name != "" {
					return NewAchievementsScreen(s.width, s.height, name), textinput.Blink
				}
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
		case key.Matches(msg, keys.Global.Down):
			s.scrollOffset++
		case key.Matches(msg, keys.Global.Up):
			if s.scrollOffset > 0 {
				s.scrollOffset--
			}
		case key.Matches(msg, keys.Global.Left):
			if s.selectedCategory > 0 {
				s.selectedCategory--
			}
		case key.Matches(msg, keys.Global.Right):
			if s.selectedCategory < len(s.categories)-1 {
				s.selectedCategory++
			}
		}
	}

	return s, nil
}

// View renders achievements
func (s *AchievementsScreen) View() string {
	// Name input screen
	if s.needsName {
		return s.renderNameInput()
	}

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Gold).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🏆 ACHIEVEMENTS 🏆")))
	b.WriteString("\n\n")

	// Player stats
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Gold).
		Padding(0, 2).
		Width(50)

	unlockedCount := len(s.unlockedAchs)
	totalCount := len(achievements.AllAchievements)

	stats := fmt.Sprintf("Player: %s\nUnlocked: %d/%d\nTotal Points: %d",
		s.playerName, unlockedCount, totalCount, s.totalPoints)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(statsBox.Render(stats)))
	b.WriteString("\n\n")

	// Category tabs
	var tabs []string
	for i, cat := range s.categories {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == s.selectedCategory {
			style = style.Foreground(styles.Black).Background(styles.Cyan).Bold(true)
		} else {
			style = style.Foreground(styles.Gray)
		}
		tabs = append(tabs, style.Render(cat))
	}
	tabRow := strings.Join(tabs, " ")
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(tabRow))
	b.WriteString("\n\n")

	// Achievement list
	achBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(65).
		Height(15)

	var achList strings.Builder
	displayCount := 0

	for _, ach := range achievements.AllAchievements {
		// Filter by category
		if s.categories[s.selectedCategory] != "All" {
			if ach.Category != s.categories[s.selectedCategory] {
				continue
			}
		}

		// Skip if before scroll offset
		if displayCount < s.scrollOffset {
			displayCount++
			continue
		}

		// Check if unlocked
		isUnlocked := false
		for _, id := range s.unlockedAchs {
			if id == ach.ID {
				isUnlocked = true
				break
			}
		}

		// Render achievement
		if isUnlocked {
			nameStyle := lipgloss.NewStyle().Foreground(styles.Gold).Bold(true)
			achList.WriteString(nameStyle.Render(fmt.Sprintf("⭐ %s", ach.Name)))
		} else {
			nameStyle := lipgloss.NewStyle().Foreground(styles.Gray)
			achList.WriteString(nameStyle.Render(fmt.Sprintf("🔒 %s", ach.Name)))
		}
		achList.WriteString("\n")

		descStyle := lipgloss.NewStyle().Foreground(styles.DimWhite).PaddingLeft(3)
		achList.WriteString(descStyle.Render(ach.Description))
		achList.WriteString("\n")

		pointStyle := lipgloss.NewStyle().Foreground(styles.Yellow).PaddingLeft(3)
		achList.WriteString(pointStyle.Render(fmt.Sprintf("%d points", ach.Points)))
		achList.WriteString("\n\n")

		displayCount++
		if displayCount-s.scrollOffset >= 4 { // Show 4 achievements at a time
			break
		}
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(achBox.Render(achList.String())))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("←/→ category • ↑/↓ scroll • esc back"))

	return b.String()
}

func (s *AchievementsScreen) renderNameInput() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Gold).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🏆 ACHIEVEMENTS 🏆")))
	b.WriteString("\n\n")

	// Explanation
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2).
		Width(50)

	info := "Enter your player name to view your achievements."
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(infoBox.Render(info)))
	b.WriteString("\n\n")

	// Name input
	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(titleStyle.Render("PLAYER NAME")))
	b.WriteString("\n")

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(35)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render(s.nameInput.View())))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter confirm • esc back"))

	return b.String()
}

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
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
	"github.com/jamesacampbell/unicorn/upgrades"
)

// UpgradesScreen shows the upgrade store
type UpgradesScreen struct {
	width           int
	height          int
	playerName      string
	availablePoints int
	ownedUpgrades   []string
	menu            *components.Menu
	selectedUpgrade *upgrades.Upgrade
	confirmPurchase bool

	// Name entry (when player name not provided)
	needsName bool
	nameInput textinput.Model
}

// NewUpgradesScreen creates a new upgrades screen
func NewUpgradesScreen(width, height int, playerName string) *UpgradesScreen {
	// If no player name, show name input first
	if playerName == "" {
		nameInput := textinput.New()
		nameInput.Placeholder = "Enter your name"
		nameInput.Focus()
		nameInput.CharLimit = 30
		nameInput.Width = 30

		return &UpgradesScreen{
			width:     width,
			height:    height,
			needsName: true,
			nameInput: nameInput,
		}
	}

	// Load player data
	allUnlocked, _ := database.GetPlayerAchievements(playerName)
	ownedUpgrades, _ := database.GetPlayerUpgrades(playerName)

	// Calculate available points
	totalPoints := 0
	for _, id := range allUnlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			totalPoints += ach.Points
		}
	}

	// Add level-up points
	profile, _ := database.GetPlayerProfile(playerName)
	if profile != nil {
		totalPoints += profile.LevelUpPoints
	}

	// Subtract spent points
	for _, id := range ownedUpgrades {
		if up, exists := upgrades.AllUpgrades[id]; exists {
			totalPoints -= up.Cost
		}
	}

	// Build menu items
	var menuItems []components.MenuItem
	for _, upgrade := range upgrades.AllUpgrades {
		owned := false
		for _, id := range ownedUpgrades {
			if id == upgrade.ID {
				owned = true
				break
			}
		}

		var title, desc string
		if owned {
			title = fmt.Sprintf("✓ %s %s", upgrade.Icon, upgrade.Name)
			desc = "OWNED"
		} else {
			title = fmt.Sprintf("%s %s (%d pts)", upgrade.Icon, upgrade.Name, upgrade.Cost)
			desc = upgrade.Description
		}

		menuItems = append(menuItems, components.MenuItem{
			ID:          upgrade.ID,
			Title:       title,
			Description: desc,
			Disabled:    owned || totalPoints < upgrade.Cost,
		})
	}

	menu := components.NewMenu("", menuItems)
	menu.SetSize(70, 20)
	menu.SetHideHelp(true)

	return &UpgradesScreen{
		width:           width,
		height:          height,
		playerName:      playerName,
		availablePoints: totalPoints,
		ownedUpgrades:   ownedUpgrades,
		menu:            menu,
	}
}

// Init initializes the upgrades screen
func (s *UpgradesScreen) Init() tea.Cmd {
	if s.needsName {
		return textinput.Blink
	}
	return nil
}

// Update handles upgrades input
func (s *UpgradesScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
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
					// Reload screen with player name
					return NewUpgradesScreen(s.width, s.height, name), textinput.Blink
				}
			}
		}

		var cmd tea.Cmd
		s.nameInput, cmd = s.nameInput.Update(msg)
		return s, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.confirmPurchase {
			switch {
			case msg.String() == "y", msg.String() == "Y":
				return s.doPurchase()
			case msg.String() == "n", msg.String() == "N", key.Matches(msg, keys.Global.Back):
				s.confirmPurchase = false
				s.selectedUpgrade = nil
			}
			return s, nil
		}

		if key.Matches(msg, keys.Global.Back) {
			return s, PopScreen()
		}

	case components.MenuSelectedMsg:
		// Find the upgrade
		if upgrade, exists := upgrades.AllUpgrades[msg.ID]; exists {
			// Check if can purchase
			owned := false
			for _, id := range s.ownedUpgrades {
				if id == msg.ID {
					owned = true
					break
				}
			}

			if !owned && s.availablePoints >= upgrade.Cost {
				s.selectedUpgrade = &upgrade
				s.confirmPurchase = true
			}
		}
	}

	var cmd tea.Cmd
	s.menu, cmd = s.menu.Update(msg)
	return s, cmd
}

func (s *UpgradesScreen) doPurchase() (ScreenModel, tea.Cmd) {
	if s.selectedUpgrade == nil {
		return s, nil
	}

	err := database.PurchaseUpgrade(s.playerName, s.selectedUpgrade.ID)
	if err != nil {
		s.confirmPurchase = false
		s.selectedUpgrade = nil
		return s, nil
	}

	// Update state
	s.ownedUpgrades = append(s.ownedUpgrades, s.selectedUpgrade.ID)
	s.availablePoints -= s.selectedUpgrade.Cost
	s.confirmPurchase = false
	s.selectedUpgrade = nil

	// Refresh screen
	return NewUpgradesScreen(s.width, s.height, s.playerName), nil
}

// View renders the upgrades screen
func (s *UpgradesScreen) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🎁 UPGRADE STORE 🎁")))
	b.WriteString("\n\n")

	// Name input screen
	if s.needsName {
		return s.renderNameInput()
	}

	// Points display
	pointsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(0, 2).
		Width(40)

	pointsStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	points := fmt.Sprintf("Player: %s\nAvailable Points: %s",
		s.playerName, pointsStyle.Render(fmt.Sprintf("%d", s.availablePoints)))
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(pointsBox.Render(points)))
	b.WriteString("\n\n")

	// Confirmation dialog
	if s.confirmPurchase && s.selectedUpgrade != nil {
		dialogBox := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(styles.Yellow).
			Padding(1, 2).
			Width(50)

		dialog := fmt.Sprintf("Purchase %s %s for %d points?\n\n[Y]es / [N]o",
			s.selectedUpgrade.Icon, s.selectedUpgrade.Name, s.selectedUpgrade.Cost)
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(dialogBox.Render(dialog)))
		b.WriteString("\n\n")
	} else {
		// Menu
		menuBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Magenta).
			Padding(0, 1)

		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(menuBox.Render(s.menu.View())))
		b.WriteString("\n\n")
	}

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("↑/↓ navigate • enter purchase • esc back"))

	return b.String()
}

func (s *UpgradesScreen) renderNameInput() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🎁 UPGRADE STORE 🎁")))
	b.WriteString("\n\n")

	// Explanation
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2).
		Width(50)

	info := "Enter your player name to view your upgrades and achievement points."
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

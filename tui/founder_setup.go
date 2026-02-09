package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/founder"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
	"github.com/jamesacampbell/unicorn/upgrades"
)

// FounderSetupStep represents steps in founder setup
type FounderSetupStep int

const (
	FounderStepName FounderSetupStep = iota
	FounderStepCompany
	FounderStepCategory
	FounderStepDifficulty
	FounderStepReady
)

// FounderSetupScreen handles founder game setup
type FounderSetupScreen struct {
	width    int
	height   int
	gameData *GameData
	step     FounderSetupStep

	// Inputs
	nameInput      textinput.Model
	companyInput   textinput.Model
	categoryMenu   *components.Menu
	difficultyMenu *components.Menu

	// Data
	playerName  string
	companyName string
	category    string
	difficulty  string

	// Player state
	playerLevel int
	welcomeBack bool
	playerStats *database.PlayerStats
}

// NewFounderSetupScreen creates a new founder setup screen
func NewFounderSetupScreen(width, height int, gameData *GameData) *FounderSetupScreen {
	// Name input
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter your name"
	nameInput.Focus()
	nameInput.CharLimit = 30
	nameInput.Width = 30

	// Company input
	companyInput := textinput.New()
	companyInput.Placeholder = "e.g., Acme Inc."
	companyInput.CharLimit = 40
	companyInput.Width = 30

	// Category menu
	categoryItems := []components.MenuItem{
		{ID: "SaaS", Title: "SaaS", Description: "Software as a Service", Icon: "💻"},
		{ID: "Fintech", Title: "Fintech", Description: "Financial Technology", Icon: "💳"},
		{ID: "E-commerce", Title: "E-commerce", Description: "Online Retail", Icon: "🛒"},
		{ID: "AI/ML", Title: "AI/ML", Description: "Artificial Intelligence", Icon: "🤖"},
		{ID: "Healthcare", Title: "Healthcare", Description: "Health Technology", Icon: "🏥"},
		{ID: "Consumer", Title: "Consumer", Description: "Consumer Apps", Icon: "📱"},
		{ID: "Enterprise", Title: "Enterprise", Description: "B2B Solutions", Icon: "🏢"},
		{ID: "Crypto", Title: "Crypto/Web3", Description: "Blockchain Technology", Icon: "🔗"},
	}
	categoryMenu := components.NewMenu("SELECT CATEGORY", categoryItems)
	categoryMenu.SetSize(50, 15)
	categoryMenu.SetHideHelp(true)

	// Get player level for difficulty unlocks
	playerLevel := 1
	profile, err := database.GetPlayerProfile(gameData.PlayerName)
	if err == nil && profile != nil {
		playerLevel = profile.Level
	}

	// Difficulty menu
	difficultyItems := []components.MenuItem{
		{
			ID:          "easy",
			Title:       "Easy",
			Description: "More cash, slower burn, friendlier market",
			Icon:        "🟢",
		},
		{
			ID:          "medium",
			Title:       "Medium",
			Description: "Balanced challenge, realistic conditions",
			Icon:        "🟡",
		},
		{
			ID:          "hard",
			Title:       "Hard",
			Description: "Less cash, aggressive competition, higher churn",
			Icon:        "🔴",
			Disabled:    playerLevel < 5,
		},
		{
			ID:          "expert",
			Title:       "Expert",
			Description: "Minimal runway, brutal market, investor pressure",
			Icon:        "💀",
			Disabled:    playerLevel < 10,
		},
	}
	difficultyMenu := components.NewMenu("SELECT DIFFICULTY", difficultyItems)
	difficultyMenu.SetSize(60, 15)
	difficultyMenu.SetHideHelp(true)

	return &FounderSetupScreen{
		width:          width,
		height:         height,
		gameData:       gameData,
		step:           FounderStepName,
		nameInput:      nameInput,
		companyInput:   companyInput,
		categoryMenu:   categoryMenu,
		difficultyMenu: difficultyMenu,
		playerLevel:    playerLevel,
		difficulty:     "easy",
	}
}

// Init initializes the founder setup screen
func (s *FounderSetupScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles founder setup input
func (s *FounderSetupScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Global.Back) {
			if s.step > FounderStepName {
				s.step--
				return s, nil
			}
			return s, SwitchTo(ScreenVCSetup)
		}

		switch s.step {
		case FounderStepName:
			if msg.Type == tea.KeyEnter {
				name := strings.TrimSpace(s.nameInput.Value())
				if name != "" {
					s.playerName = name
					s.gameData.PlayerName = name

					// Check for returning player
					stats, err := database.GetPlayerStats(name)
					if err == nil && stats != nil && stats.TotalGames > 0 {
						s.welcomeBack = true
						s.playerStats = stats
					}

					// Refresh player level
					profile, err := database.GetPlayerProfile(name)
					if err == nil && profile != nil {
						s.playerLevel = profile.Level
						// Refresh difficulty menu locks
						s.refreshDifficultyMenu()
					}

					// Load player upgrades from DB
					playerUpgrades, err := database.GetPlayerUpgrades(name)
					if err == nil {
						s.gameData.PlayerUpgrades = playerUpgrades
					}

					s.step = FounderStepCompany
					s.companyInput.Focus()
					return s, textinput.Blink
				}
			}

		case FounderStepCompany:
			if msg.Type == tea.KeyEnter {
				company := strings.TrimSpace(s.companyInput.Value())
				if company != "" {
					s.companyName = company
					s.step = FounderStepCategory
					return s, nil
				}
			}
		}

	case components.MenuSelectedMsg:
		switch s.step {
		case FounderStepCategory:
			s.category = msg.ID
			s.step = FounderStepDifficulty
			return s, nil

		case FounderStepDifficulty:
			switch msg.ID {
			case "easy":
				s.difficulty = "easy"
			case "medium":
				s.difficulty = "medium"
			case "hard":
				if s.playerLevel >= 5 {
					s.difficulty = "hard"
				} else {
					return s, nil
				}
			case "expert":
				if s.playerLevel >= 10 {
					s.difficulty = "expert"
				} else {
					return s, nil
				}
			}
			return s, s.startGame()
		}
	}

	// Update current component
	var cmd tea.Cmd
	switch s.step {
	case FounderStepName:
		s.nameInput, cmd = s.nameInput.Update(msg)
	case FounderStepCompany:
		s.companyInput, cmd = s.companyInput.Update(msg)
	case FounderStepCategory:
		s.categoryMenu, cmd = s.categoryMenu.Update(msg)
	case FounderStepDifficulty:
		s.difficultyMenu, cmd = s.difficultyMenu.Update(msg)
	}

	return s, cmd
}

func (s *FounderSetupScreen) refreshDifficultyMenu() {
	difficultyItems := []components.MenuItem{
		{
			ID:          "easy",
			Title:       "Easy",
			Description: "More cash, slower burn, friendlier market",
			Icon:        "🟢",
		},
		{
			ID:          "medium",
			Title:       "Medium",
			Description: "Balanced challenge, realistic conditions",
			Icon:        "🟡",
		},
		{
			ID:          "hard",
			Title:       "Hard",
			Description: "Less cash, aggressive competition, higher churn",
			Icon:        "🔴",
			Disabled:    s.playerLevel < 5,
		},
		{
			ID:          "expert",
			Title:       "Expert",
			Description: "Minimal runway, brutal market, investor pressure",
			Icon:        "💀",
			Disabled:    s.playerLevel < 10,
		},
	}
	s.difficultyMenu = components.NewMenu("SELECT DIFFICULTY", difficultyItems)
	s.difficultyMenu.SetSize(60, 15)
	s.difficultyMenu.SetHideHelp(true)
}

func (s *FounderSetupScreen) startGame() tea.Cmd {
	return func() tea.Msg {
		// Load startup templates
		templates, err := founder.LoadFounderStartups("founder/startups.json")
		if err != nil || len(templates) == 0 {
			// Create a default template if loading fails
			template := founder.StartupTemplate{
				ID:               "custom",
				Name:             s.companyName,
				Type:             s.category,
				Description:      "A new startup",
				InitialCash:      500000,
				MonthlyBurn:      50000,
				InitialCustomers: 10,
				InitialMRR:       5000,
				AvgDealSize:      500,
				BaseChurnRate:    0.05,
				BaseCAC:          1000,
				TargetMarketSize: 10000,
			}

			// Apply difficulty modifiers
			s.applyDifficultyToTemplate(&template)

			s.gameData.FounderState = founder.NewFounderGame(s.playerName, template, s.gameData.PlayerUpgrades)
		} else {
			// Find template matching category, or use first one
			var selectedTemplate founder.StartupTemplate
			for _, t := range templates {
				if t.Type == s.category {
					selectedTemplate = t
					selectedTemplate.Name = s.companyName
					break
				}
			}
			if selectedTemplate.ID == "" {
				selectedTemplate = templates[0]
				selectedTemplate.Name = s.companyName
			}

			// Apply difficulty modifiers
			s.applyDifficultyToTemplate(&selectedTemplate)

			s.gameData.FounderState = founder.NewFounderGame(s.playerName, selectedTemplate, s.gameData.PlayerUpgrades)
		}

		return SwitchScreenMsg{Screen: ScreenFounderGame}
	}
}

// applyDifficultyToTemplate modifies the template based on difficulty
func (s *FounderSetupScreen) applyDifficultyToTemplate(t *founder.StartupTemplate) {
	switch s.difficulty {
	case "easy":
		t.InitialCash = int64(float64(t.InitialCash) * 1.5)
		t.BaseChurnRate = t.BaseChurnRate * 0.7
		t.BaseCAC = int64(float64(t.BaseCAC) * 0.8)
	case "medium":
		// Default - no changes
	case "hard":
		t.InitialCash = int64(float64(t.InitialCash) * 0.7)
		t.BaseChurnRate = t.BaseChurnRate * 1.3
		t.BaseCAC = int64(float64(t.BaseCAC) * 1.3)
	case "expert":
		t.InitialCash = int64(float64(t.InitialCash) * 0.5)
		t.BaseChurnRate = t.BaseChurnRate * 1.6
		t.BaseCAC = int64(float64(t.BaseCAC) * 1.6)
		t.MonthlyBurn = int64(float64(t.MonthlyBurn) * 1.3)
	}
}

// View renders the founder setup screen
func (s *FounderSetupScreen) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(60).
		Align(lipgloss.Center).
		Padding(0, 2)

	header := headerStyle.Render("🚀 FOUNDER MODE SETUP")
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(header))
	b.WriteString("\n\n")

	// Welcome back message
	if s.welcomeBack && s.playerStats != nil && s.step == FounderStepCompany {
		welcomeBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Green).
			Padding(0, 2).
			Width(55)

		welcomeStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		var welcome strings.Builder
		welcome.WriteString(welcomeStyle.Render(fmt.Sprintf("Welcome back, %s!", s.playerName)))
		welcome.WriteString("\n")
		welcome.WriteString(fmt.Sprintf("Games Played: %d | Best Net Worth: $%s | Win Rate: %.0f%%",
			s.playerStats.TotalGames,
			formatCompactMoney(s.playerStats.BestNetWorth),
			s.playerStats.WinRate))

		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(welcomeBox.Render(welcome.String())))
		b.WriteString("\n\n")
	}

	// Active upgrades display
	if s.step == FounderStepDifficulty && len(s.gameData.PlayerUpgrades) > 0 {
		founderUpgrades := upgrades.FilterUpgradeIDsForGameMode(s.gameData.PlayerUpgrades, "founder")
		if len(founderUpgrades) > 0 {
			upgradeBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(styles.Green).
				Padding(0, 2).
				Width(55)

			var upgradeContent strings.Builder
			upgradeStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
			upgradeContent.WriteString(upgradeStyle.Render("✨ Active Upgrades:"))
			upgradeContent.WriteString("\n")
			for _, upgradeID := range founderUpgrades {
				if upgrade, exists := upgrades.AllUpgrades[upgradeID]; exists {
					upgradeContent.WriteString(fmt.Sprintf("  %s %s", upgrade.Icon, upgrade.Name))
					upgradeContent.WriteString("\n")
				}
			}

			b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(upgradeBox.Render(upgradeContent.String())))
			b.WriteString("\n")
		}
	}

	// Progress
	b.WriteString(s.renderProgress())
	b.WriteString("\n\n")

	// Content box
	contentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2).
		Width(55)

	var content string
	switch s.step {
	case FounderStepName:
		content = s.renderNameStep()
	case FounderStepCompany:
		content = s.renderCompanyStep()
	case FounderStepCategory:
		content = s.categoryMenu.View()
	case FounderStepDifficulty:
		content = s.renderDifficultyStep()
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(contentBox.Render(content)))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter continue"))

	return b.String()
}

func (s *FounderSetupScreen) renderProgress() string {
	steps := []string{"Name", "Company", "Category", "Difficulty"}
	var parts []string

	for i, step := range steps {
		var style lipgloss.Style
		if i < int(s.step) {
			style = lipgloss.NewStyle().Foreground(styles.Green)
			parts = append(parts, style.Render("✓ "+step))
		} else if i == int(s.step) {
			style = lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
			parts = append(parts, style.Render("● "+step))
		} else {
			style = lipgloss.NewStyle().Foreground(styles.Gray)
			parts = append(parts, style.Render("○ "+step))
		}
	}

	return lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(strings.Join(parts, "  →  "))
}

func (s *FounderSetupScreen) renderNameStep() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(titleStyle.Render("ENTER YOUR NAME"))
	b.WriteString("\n\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Magenta).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(s.nameInput.View()))

	return b.String()
}

func (s *FounderSetupScreen) renderCompanyStep() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(titleStyle.Render("NAME YOUR STARTUP"))
	b.WriteString("\n\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Magenta).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(s.companyInput.View()))
	b.WriteString("\n\n")

	hintStyle := lipgloss.NewStyle().Foreground(styles.Gray).Italic(true)
	b.WriteString(hintStyle.Render(fmt.Sprintf("Founded by %s", s.playerName)))

	return b.String()
}

func (s *FounderSetupScreen) renderDifficultyStep() string {
	var b strings.Builder

	// Show level info
	if s.playerLevel > 1 {
		levelStyle := lipgloss.NewStyle().Foreground(styles.Cyan)
		b.WriteString(levelStyle.Render(fmt.Sprintf("Your Level: %d", s.playerLevel)))
		if s.playerLevel < 5 {
			lockStyle := lipgloss.NewStyle().Foreground(styles.Gray)
			b.WriteString(lockStyle.Render(" (Hard unlocks at Lvl 5, Expert at Lvl 10)"))
		} else if s.playerLevel < 10 {
			lockStyle := lipgloss.NewStyle().Foreground(styles.Gray)
			b.WriteString(lockStyle.Render(" (Expert unlocks at Lvl 10)"))
		}
		b.WriteString("\n\n")
	}

	b.WriteString(s.difficultyMenu.View())

	return b.String()
}

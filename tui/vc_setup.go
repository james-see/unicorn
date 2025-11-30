package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/game"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// SetupStep represents the current step in the setup wizard
type SetupStep int

const (
	StepGameMode SetupStep = iota
	StepPlayerName
	StepDifficulty
	StepFirmName
	StepPlayMode
	StepReady
)

// VCSetupScreen handles the game setup flow
type VCSetupScreen struct {
	width    int
	height   int
	gameData *GameData
	step     SetupStep

	// Components for each step
	gameModeMenu   *components.Menu
	nameInput      textinput.Model
	difficultyMenu *components.Menu
	firmInput      textinput.Model
	playModeMenu   *components.Menu

	// Computed values
	playerLevel int
	welcomeBack bool
	playerStats *database.PlayerStats
}

// NewVCSetupScreen creates a new setup screen
func NewVCSetupScreen(width, height int, gameData *GameData) *VCSetupScreen {
	// Game mode menu
	gameModeItems := []components.MenuItem{
		{
			ID:          "vc",
			Title:       "VC Investor Mode (Classic)",
			Description: "Build a portfolio of startups and compete against AI investors",
			Icon:        "💼",
		},
		{
			ID:          "founder",
			Title:       "Startup Founder Mode",
			Description: "Build your own startup from the ground up",
			Icon:        "🚀",
		},
	}
	gameModeMenu := components.NewMenu("SELECT GAME MODE", gameModeItems)
	gameModeMenu.SetSize(60, 10)
	gameModeMenu.SetHideHelp(true)

	// Player name input
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter your name"
	nameInput.Focus()
	nameInput.CharLimit = 30
	nameInput.Width = 30

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
			Description: game.EasyDifficulty.Description,
			Icon:        "🟢",
		},
		{
			ID:          "medium",
			Title:       "Medium",
			Description: game.MediumDifficulty.Description,
			Icon:        "🟡",
		},
		{
			ID:          "hard",
			Title:       "Hard",
			Description: game.HardDifficulty.Description,
			Icon:        "🔴",
			Disabled:    playerLevel < 5,
		},
		{
			ID:          "expert",
			Title:       "Expert",
			Description: game.ExpertDifficulty.Description,
			Icon:        "💀",
			Disabled:    playerLevel < 10,
		},
	}
	difficultyMenu := components.NewMenu("SELECT DIFFICULTY", difficultyItems)
	difficultyMenu.SetSize(60, 15)
	difficultyMenu.SetHideHelp(true)

	// Firm name input
	firmInput := textinput.New()
	firmInput.Placeholder = "e.g., Sequoia Capital"
	firmInput.CharLimit = 40
	firmInput.Width = 30

	// Play mode menu
	playModeItems := []components.MenuItem{
		{
			ID:          "manual",
			Title:       "Manual Mode (Recommended)",
			Description: "Press Enter each turn, full access to all features",
			Icon:        "🎮",
		},
		{
			ID:          "auto",
			Title:       "Automated Mode",
			Description: "1 second per turn, simplified gameplay",
			Icon:        "⏩",
		},
	}
	playModeMenu := components.NewMenu("SELECT PLAY MODE", playModeItems)
	playModeMenu.SetSize(60, 10)
	playModeMenu.SetHideHelp(true)

	return &VCSetupScreen{
		width:          width,
		height:         height,
		gameData:       gameData,
		step:           StepGameMode,
		gameModeMenu:   gameModeMenu,
		nameInput:      nameInput,
		difficultyMenu: difficultyMenu,
		firmInput:      firmInput,
		playModeMenu:   playModeMenu,
		playerLevel:    playerLevel,
	}
}

// Init initializes the setup screen
func (s *VCSetupScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles setup screen input
func (s *VCSetupScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle escape to go back
		if key.Matches(msg, keys.Global.Back) {
			if s.step > StepGameMode {
				s.step--
				return s, nil
			}
			return s, SwitchTo(ScreenMainMenu)
		}

	case components.MenuSelectedMsg:
		return s.handleMenuSelection(msg.ID)
	}

	// Update current step's component
	var cmd tea.Cmd
	switch s.step {
	case StepGameMode:
		s.gameModeMenu, cmd = s.gameModeMenu.Update(msg)
	case StepPlayerName:
		s.nameInput, cmd = s.nameInput.Update(msg)
		// Check for enter to proceed
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
			return s.handleNameSubmit()
		}
	case StepDifficulty:
		s.difficultyMenu, cmd = s.difficultyMenu.Update(msg)
	case StepFirmName:
		s.firmInput, cmd = s.firmInput.Update(msg)
		// Check for enter to proceed
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
			return s.handleFirmSubmit()
		}
	case StepPlayMode:
		s.playModeMenu, cmd = s.playModeMenu.Update(msg)
	}

	return s, cmd
}

func (s *VCSetupScreen) handleMenuSelection(id string) (ScreenModel, tea.Cmd) {
	switch s.step {
	case StepGameMode:
		if id == "founder" {
			return s, SwitchTo(ScreenFounderSetup)
		}
		// VC mode selected, move to name
		s.step = StepPlayerName
		s.nameInput.Focus()
		return s, textinput.Blink

	case StepDifficulty:
		// Set difficulty
		switch id {
		case "easy":
			s.gameData.Difficulty = game.EasyDifficulty
		case "medium":
			s.gameData.Difficulty = game.MediumDifficulty
		case "hard":
			if s.playerLevel >= 5 {
				s.gameData.Difficulty = game.HardDifficulty
			} else {
				return s, nil // Locked
			}
		case "expert":
			if s.playerLevel >= 10 {
				s.gameData.Difficulty = game.ExpertDifficulty
			} else {
				return s, nil // Locked
			}
		}
		// Move to firm name
		s.step = StepFirmName
		s.firmInput.SetValue(game.GenerateDefaultFirmName(s.gameData.PlayerName))
		s.firmInput.Focus()
		return s, textinput.Blink

	case StepPlayMode:
		s.gameData.AutoMode = (id == "auto")
		// Setup complete, start game
		return s, s.startGame()
	}

	return s, nil
}

func (s *VCSetupScreen) handleNameSubmit() (ScreenModel, tea.Cmd) {
	name := strings.TrimSpace(s.nameInput.Value())
	if name == "" {
		return s, nil
	}

	s.gameData.PlayerName = name

	// Check for returning player
	stats, err := database.GetPlayerStats(name)
	if err == nil && stats.TotalGames > 0 {
		s.welcomeBack = true
		s.playerStats = stats

		// Update player level
		profile, err := database.GetPlayerProfile(name)
		if err == nil && profile != nil {
			s.playerLevel = profile.Level
		}

		// Update difficulty menu for unlocks
		s.updateDifficultyUnlocks()
	}

	// Load player upgrades
	upgrades, _ := database.GetPlayerUpgrades(name)
	s.gameData.PlayerUpgrades = upgrades

	s.step = StepDifficulty
	return s, nil
}

func (s *VCSetupScreen) handleFirmSubmit() (ScreenModel, tea.Cmd) {
	firm := strings.TrimSpace(s.firmInput.Value())
	if firm == "" {
		firm = game.GenerateDefaultFirmName(s.gameData.PlayerName)
	}
	s.gameData.FirmName = firm

	s.step = StepPlayMode
	return s, nil
}

func (s *VCSetupScreen) updateDifficultyUnlocks() {
	// Rebuild difficulty menu with current unlock status
	difficultyItems := []components.MenuItem{
		{
			ID:          "easy",
			Title:       "Easy",
			Description: fmt.Sprintf("%s | $%s", game.EasyDifficulty.Description, formatSetupMoney(game.EasyDifficulty.StartingCash)),
			Icon:        "🟢",
		},
		{
			ID:          "medium",
			Title:       "Medium",
			Description: fmt.Sprintf("%s | $%s", game.MediumDifficulty.Description, formatSetupMoney(game.MediumDifficulty.StartingCash)),
			Icon:        "🟡",
		},
		{
			ID:          "hard",
			Title:       fmt.Sprintf("Hard %s", unlockText(s.playerLevel >= 5, 5)),
			Description: fmt.Sprintf("%s | $%s", game.HardDifficulty.Description, formatSetupMoney(game.HardDifficulty.StartingCash)),
			Icon:        "🔴",
			Disabled:    s.playerLevel < 5,
		},
		{
			ID:          "expert",
			Title:       fmt.Sprintf("Expert %s", unlockText(s.playerLevel >= 10, 10)),
			Description: fmt.Sprintf("%s | $%s", game.ExpertDifficulty.Description, formatSetupMoney(game.ExpertDifficulty.StartingCash)),
			Icon:        "💀",
			Disabled:    s.playerLevel < 10,
		},
	}
	s.difficultyMenu = components.NewMenu("SELECT DIFFICULTY", difficultyItems)
	s.difficultyMenu.SetSize(60, 15)
	s.difficultyMenu.SetHideHelp(true)
}

func unlockText(unlocked bool, level int) string {
	if unlocked {
		return ""
	}
	return fmt.Sprintf("🔒 Lvl %d", level)
}

func (s *VCSetupScreen) startGame() tea.Cmd {
	return func() tea.Msg {
		// Create the game state
		s.gameData.GameState = game.NewGame(
			s.gameData.PlayerName,
			s.gameData.FirmName,
			s.gameData.Difficulty,
			s.gameData.PlayerUpgrades,
		)

		// Load reputation
		dbRep, err := database.GetVCReputation(s.gameData.PlayerName)
		if err == nil && dbRep != nil {
			s.gameData.GameState.PlayerReputation = &game.VCReputation{
				PlayerName:       dbRep.PlayerName,
				PerformanceScore: dbRep.PerformanceScore,
				FounderScore:     dbRep.FounderScore,
				MarketScore:      dbRep.MarketScore,
				TotalGamesPlayed: dbRep.TotalGamesPlayed,
				SuccessfulExits:  dbRep.SuccessfulExits,
				AvgROILast5:      dbRep.AvgROILast5,
			}
		}

		return SwitchScreenMsg{Screen: ScreenVCInvest}
	}
}

// View renders the setup screen
func (s *VCSetupScreen) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center).
		Padding(0, 2)

	header := headerStyle.Render("🦄 GAME SETUP")
	headerContainer := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(headerContainer.Render(header))
	b.WriteString("\n\n")

	// Progress indicator
	b.WriteString(s.renderProgress())
	b.WriteString("\n\n")

	// Current step content
	contentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(65)

	var content string
	switch s.step {
	case StepGameMode:
		content = s.gameModeMenu.View()
	case StepPlayerName:
		content = s.renderNameStep()
	case StepDifficulty:
		content = s.difficultyMenu.View()
	case StepFirmName:
		content = s.renderFirmStep()
	case StepPlayMode:
		content = s.playModeMenu.View()
	}

	contentContainer := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(contentContainer.Render(contentBox.Render(content)))

	// Help text
	b.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(styles.Gray).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • ↑/↓ navigate • enter select"))

	return b.String()
}

func (s *VCSetupScreen) renderProgress() string {
	steps := []string{"Mode", "Name", "Difficulty", "Firm", "Play Style"}
	var parts []string

	for i, step := range steps {
		var style lipgloss.Style
		if i < int(s.step) {
			// Completed
			style = lipgloss.NewStyle().Foreground(styles.Green)
			parts = append(parts, style.Render("✓ "+step))
		} else if i == int(s.step) {
			// Current
			style = lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
			parts = append(parts, style.Render("● "+step))
		} else {
			// Upcoming
			style = lipgloss.NewStyle().Foreground(styles.Gray)
			parts = append(parts, style.Render("○ "+step))
		}
	}

	progressStyle := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center)
	return progressStyle.Render(strings.Join(parts, "  →  "))
}

func (s *VCSetupScreen) renderNameStep() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Yellow).
		Bold(true)
	b.WriteString(titleStyle.Render("ENTER YOUR NAME"))
	b.WriteString("\n\n")

	// Input
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(s.nameInput.View()))
	b.WriteString("\n\n")

	// Welcome back message if returning player
	if s.welcomeBack && s.playerStats != nil {
		welcomeStyle := lipgloss.NewStyle().
			Foreground(styles.Green)
		b.WriteString(welcomeStyle.Render(fmt.Sprintf("🎉 Welcome back! Games: %d | Best: $%s",
			s.playerStats.TotalGames, formatSetupMoney(s.playerStats.BestNetWorth))))
	}

	return b.String()
}

func (s *VCSetupScreen) renderFirmStep() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Yellow).
		Bold(true)
	b.WriteString(titleStyle.Render("NAME YOUR VC FIRM"))
	b.WriteString("\n\n")

	// Input
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(s.firmInput.View()))
	b.WriteString("\n\n")

	hintStyle := lipgloss.NewStyle().
		Foreground(styles.Gray).
		Italic(true)
	b.WriteString(hintStyle.Render("Press Enter to use default or type your own"))

	return b.String()
}

func formatSetupMoney(amount int64) string {
	if amount >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(amount)/1000000)
	}
	if amount >= 1000 {
		return fmt.Sprintf("%.0fK", float64(amount)/1000)
	}
	return fmt.Sprintf("%d", amount)
}

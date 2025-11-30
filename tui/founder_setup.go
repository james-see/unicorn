package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/founder"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// FounderSetupStep represents steps in founder setup
type FounderSetupStep int

const (
	FounderStepName FounderSetupStep = iota
	FounderStepCompany
	FounderStepCategory
	FounderStepReady
)

// FounderSetupScreen handles founder game setup
type FounderSetupScreen struct {
	width    int
	height   int
	gameData *GameData
	step     FounderSetupStep

	// Inputs
	nameInput    textinput.Model
	companyInput textinput.Model
	categoryMenu *components.Menu

	// Data
	playerName  string
	companyName string
	category    string
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

	return &FounderSetupScreen{
		width:        width,
		height:       height,
		gameData:     gameData,
		step:         FounderStepName,
		nameInput:    nameInput,
		companyInput: companyInput,
		categoryMenu: categoryMenu,
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
		if s.step == FounderStepCategory {
			s.category = msg.ID
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
	}

	return s, cmd
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
			s.gameData.FounderState = founder.NewFounderGame(s.playerName, template, s.gameData.PlayerUpgrades)
		} else {
			// Find template matching category, or use first one
			var selectedTemplate founder.StartupTemplate
			for _, t := range templates {
				if t.Type == s.category {
					selectedTemplate = t
					selectedTemplate.Name = s.companyName // Use custom company name
					break
				}
			}
			if selectedTemplate.ID == "" {
				selectedTemplate = templates[0]
				selectedTemplate.Name = s.companyName
			}
			s.gameData.FounderState = founder.NewFounderGame(s.playerName, selectedTemplate, s.gameData.PlayerUpgrades)
		}

		return SwitchScreenMsg{Screen: ScreenFounderGame}
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
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(contentBox.Render(content)))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter continue"))

	return b.String()
}

func (s *FounderSetupScreen) renderProgress() string {
	steps := []string{"Name", "Company", "Category"}
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

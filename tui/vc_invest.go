package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/game"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// InvestPhase tracks what part of investment we're in
type InvestPhase int

const (
	PhaseStartupList InvestPhase = iota
	PhaseAmountInput
	PhaseTermsSelect
	PhaseDueDiligence
	PhaseDDResults
	PhaseSyndicateList
	PhaseSyndicateAmount
)

// VCInvestScreen handles the investment phase
type VCInvestScreen struct {
	width    int
	height   int
	gameData *GameData
	phase    InvestPhase

	// Components
	startupTable   *components.GameTable
	amountInput    textinput.Model
	termsMenu      *components.Menu
	ddMenu         *components.Menu
	syndicateTable *components.GameTable

	// State
	selectedIdx     int
	selectedStartup *game.Startup
	investAmount    int64
	errorMsg        string
	selectedTerms   game.InvestmentTerms

	// Maps table row index to AvailableStartups index
	rowToStartupIdx []int

	// Due diligence state
	ddFindings    []game.DDFinding
	ddLevel       string
	ddShouldBlock bool
	ddBlockReason string

	// Syndicate state
	selectedSyndicate int
}

// NewVCInvestScreen creates a new investment screen
func NewVCInvestScreen(width, height int, gameData *GameData) *VCInvestScreen {
	// Amount input
	amountInput := textinput.New()
	amountInput.Placeholder = "Enter amount (e.g., 100000)"
	amountInput.CharLimit = 15
	amountInput.Width = 20

	s := &VCInvestScreen{
		width:       width,
		height:      height,
		gameData:    gameData,
		phase:       PhaseStartupList,
		amountInput: amountInput,
	}

	s.refreshStartupTable()
	return s
}

func (s *VCInvestScreen) refreshStartupTable() {
	gs := s.gameData.GameState

	// Build a set of already-invested company names
	investedNames := make(map[string]bool)
	for _, inv := range gs.Portfolio.Investments {
		investedNames[inv.CompanyName] = true
	}

	// Build startup rows, skipping already-invested companies
	var rows []table.Row
	s.rowToStartupIdx = nil // Reset the mapping
	rowNum := 1

	for i, startup := range gs.AvailableStartups {
		// Skip if already invested (e.g., via syndicate)
		if investedNames[startup.Name] {
			continue
		}

		riskLabel := "Low"
		if startup.RiskScore > 0.85 {
			riskLabel = "V.High"
		} else if startup.RiskScore > 0.6 {
			riskLabel = "High"
		} else if startup.RiskScore > 0.4 {
			riskLabel = "Med"
		}

		growthLabel := "High"
		if startup.GrowthPotential > 0.85 {
			growthLabel = "V.High"
		} else if startup.GrowthPotential < 0.4 {
			growthLabel = "Low"
		} else if startup.GrowthPotential < 0.6 {
			growthLabel = "Med"
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", rowNum),
			truncate(startup.Name, 14),
			truncate(startup.Category, 11),
			formatCompactMoney(startup.Valuation),
			riskLabel,
			growthLabel,
		})
		s.rowToStartupIdx = append(s.rowToStartupIdx, i)
		rowNum++
	}

	columns := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Name", Width: 14},
		{Title: "Category", Width: 11},
		{Title: "Valuation", Width: 9},
		{Title: "Risk", Width: 4},
		{Title: "Growth", Width: 6},
	}

	s.startupTable = components.NewGameTable("", columns, rows)
	s.startupTable.SetSize(70, 12)
}

// Init initializes the investment screen
func (s *VCInvestScreen) Init() tea.Cmd {
	gs := s.gameData.GameState

	// Get player level for syndicate unlock
	playerLevel := 1
	if s.gameData.PlayerName != "" {
		profile, err := database.GetPlayerProfile(s.gameData.PlayerName)
		if err == nil && profile != nil {
			playerLevel = profile.Level
		}
	}

	// Generate syndicate opportunities if unlocked (level 2+)
	if playerLevel >= 2 {
		gs.GenerateSyndicateOpportunities(playerLevel)
	}

	return nil
}

// Update handles investment screen input
func (s *VCInvestScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s.phase {
		case PhaseStartupList:
			switch {
			case key.Matches(msg, keys.Global.Back):
				return s, SwitchTo(ScreenMainMenu)
			case msg.String() == "d":
				// Done investing, start game
				gs.AIPlayerMakeInvestments()
				return s, SwitchTo(ScreenVCTurn)
			case msg.String() == "s":
				// Show syndicate opportunities
				if len(gs.SyndicateOpportunities) > 0 {
					s.buildSyndicateTable()
					s.phase = PhaseSyndicateList
					return s, nil
				}
			case key.Matches(msg, keys.Global.Enter):
				// Select startup for investment
				tableIdx := s.startupTable.Cursor()
				if tableIdx >= 0 && tableIdx < len(s.rowToStartupIdx) {
					// Map table row to actual startup index
					startupIdx := s.rowToStartupIdx[tableIdx]
					s.selectedIdx = startupIdx
					s.selectedStartup = &gs.AvailableStartups[startupIdx]
					s.phase = PhaseAmountInput
					s.amountInput.Focus()
					s.amountInput.SetValue("")
					s.errorMsg = ""
					return s, textinput.Blink
				}
			}

		case PhaseAmountInput:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.phase = PhaseStartupList
				s.errorMsg = ""
				return s, nil
			case key.Matches(msg, keys.Global.Enter):
				return s.handleAmountSubmit()
			}

		case PhaseTermsSelect:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.phase = PhaseAmountInput
				return s, textinput.Blink
			}

		case PhaseDueDiligence:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.phase = PhaseTermsSelect
				return s, nil
			}

		case PhaseDDResults:
			switch {
			case key.Matches(msg, keys.Global.Back), msg.String() == "n":
				// Cancel investment
				s.phase = PhaseStartupList
				s.ddFindings = nil
				return s, nil
			case msg.String() == "y":
				// Proceed with investment
				return s.finalizeInvestment()
			}

		case PhaseSyndicateList:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.phase = PhaseStartupList
				return s, nil
			}

		case PhaseSyndicateAmount:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.phase = PhaseSyndicateList
				s.errorMsg = ""
				return s, nil
			case key.Matches(msg, keys.Global.Enter):
				return s.handleSyndicateAmountSubmit()
			}
		}

	case components.MenuSelectedMsg:
		switch s.phase {
		case PhaseTermsSelect:
			return s.handleTermsSelection(msg.ID)
		case PhaseDueDiligence:
			return s.handleDDSelection(msg.ID)
		}

	case components.TableRowSelectedMsg:
		switch s.phase {
		case PhaseStartupList:
			tableIdx := msg.Index
			if tableIdx >= 0 && tableIdx < len(s.rowToStartupIdx) {
				// Map table row to actual startup index
				startupIdx := s.rowToStartupIdx[tableIdx]
				s.selectedIdx = startupIdx
				s.selectedStartup = &gs.AvailableStartups[startupIdx]
				s.phase = PhaseAmountInput
				s.amountInput.Focus()
				s.amountInput.SetValue("")
				s.errorMsg = ""
				return s, textinput.Blink
			}
		case PhaseSyndicateList:
			idx := msg.Index
			if idx >= 0 && idx < len(gs.SyndicateOpportunities) {
				s.selectedSyndicate = idx
				s.phase = PhaseSyndicateAmount
				s.amountInput.Focus()
				s.amountInput.SetValue("")
				s.errorMsg = ""
				return s, textinput.Blink
			}
		}
	}

	// Update current component
	var cmd tea.Cmd
	switch s.phase {
	case PhaseStartupList:
		s.startupTable, cmd = s.startupTable.Update(msg)
	case PhaseAmountInput, PhaseSyndicateAmount:
		s.amountInput, cmd = s.amountInput.Update(msg)
	case PhaseTermsSelect:
		if s.termsMenu != nil {
			s.termsMenu, cmd = s.termsMenu.Update(msg)
		}
	case PhaseDueDiligence:
		if s.ddMenu != nil {
			s.ddMenu, cmd = s.ddMenu.Update(msg)
		}
	case PhaseSyndicateList:
		if s.syndicateTable != nil {
			s.syndicateTable, cmd = s.syndicateTable.Update(msg)
		}
	}

	return s, cmd
}

func (s *VCInvestScreen) handleAmountSubmit() (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState
	amountStr := strings.TrimSpace(s.amountInput.Value())

	if amountStr == "" || amountStr == "0" {
		s.phase = PhaseStartupList
		return s, nil
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		s.errorMsg = "Invalid amount"
		return s, nil
	}

	// Validate amount
	minInvest := int64(10000)
	maxInvest := int64(float64(s.selectedStartup.Valuation) * 0.20)
	if maxInvest > gs.Portfolio.Cash {
		maxInvest = gs.Portfolio.Cash
	}

	if amount < minInvest {
		s.errorMsg = fmt.Sprintf("Minimum investment is $%d", minInvest)
		return s, nil
	}

	if amount > maxInvest {
		s.errorMsg = fmt.Sprintf("Maximum investment is $%d", maxInvest)
		return s, nil
	}

	s.investAmount = amount

	// For larger investments, show terms selection
	if amount >= 50000 {
		s.buildTermsMenu()
		s.phase = PhaseTermsSelect
		return s, nil
	}

	// Small investment - use default terms
	return s.makeInvestment(game.InvestmentTerms{
		Type:             "Common Stock",
		HasProRataRights: false,
		HasInfoRights:    false,
		HasBoardSeat:     false,
	})
}

func (s *VCInvestScreen) buildTermsMenu() {
	gs := s.gameData.GameState
	options := gs.GenerateTermOptions(s.selectedStartup, s.investAmount)

	items := make([]components.MenuItem, len(options))
	for i, opt := range options {
		desc := ""
		if opt.HasProRataRights {
			desc += "Pro-Rata "
		}
		if opt.HasBoardSeat {
			desc += "Board "
		}
		if opt.LiquidationPref > 0 {
			desc += fmt.Sprintf("%.0fx Liq ", opt.LiquidationPref)
		}
		if desc == "" {
			desc = "Basic terms"
		}

		items[i] = components.MenuItem{
			ID:          fmt.Sprintf("%d", i),
			Title:       opt.Type,
			Description: desc,
			Icon:        "📜",
		}
	}

	s.termsMenu = components.NewMenu("SELECT INVESTMENT TERMS", items)
	s.termsMenu.SetSize(50, 15)
	s.termsMenu.SetHideHelp(true)
}

func (s *VCInvestScreen) handleTermsSelection(id string) (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState
	idx, _ := strconv.Atoi(id)
	options := gs.GenerateTermOptions(s.selectedStartup, s.investAmount)

	if idx >= 0 && idx < len(options) {
		s.selectedTerms = options[idx]

		// Skip DD in auto mode
		if s.gameData.AutoMode {
			s.ddLevel = "none"
			return s.finalizeInvestment()
		}

		// Show due diligence options
		s.buildDDMenu()
		s.phase = PhaseDueDiligence
		return s, nil
	}

	return s, nil
}

func (s *VCInvestScreen) buildDDMenu() {
	levels := game.GetDDLevels()
	items := make([]components.MenuItem, len(levels))

	for i, level := range levels {
		var title, desc string
		if level.ID == "none" {
			title = fmt.Sprintf("⏭️  %s (Free)", level.Name)
			desc = level.Description
		} else {
			title = fmt.Sprintf("🔍 %s ($%s, %d days)", level.Name, formatCompactMoney(level.Cost), level.Duration)
			desc = fmt.Sprintf("Reveals: %s", strings.Join(level.Reveals, ", "))
		}

		items[i] = components.MenuItem{
			ID:          level.ID,
			Title:       title,
			Description: desc,
		}
	}

	s.ddMenu = components.NewMenu("DUE DILIGENCE OPTIONS", items)
	s.ddMenu.SetSize(60, 15)
	s.ddMenu.SetHideHelp(true)
}

func (s *VCInvestScreen) handleDDSelection(id string) (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState

	if id == "none" {
		s.ddLevel = "none"
		return s.finalizeInvestment()
	}

	// Find the DD level
	levels := game.GetDDLevels()
	var selectedLevel game.DDLevel
	for _, level := range levels {
		if level.ID == id {
			selectedLevel = level
			break
		}
	}

	// Check if can afford
	if selectedLevel.Cost > gs.Portfolio.Cash {
		s.errorMsg = "Insufficient funds for due diligence"
		return s, nil
	}

	// Perform DD
	gs.Portfolio.Cash -= selectedLevel.Cost
	s.ddFindings = game.PerformDueDiligence(s.selectedStartup, id)
	s.ddLevel = id

	// Check if should block
	s.ddShouldBlock, s.ddBlockReason = game.ShouldBlockInvestment(s.ddFindings)

	// Apply findings
	game.ApplyDDFindings(s.selectedStartup, s.ddFindings)

	s.phase = PhaseDDResults
	return s, nil
}

func (s *VCInvestScreen) finalizeInvestment() (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState

	err := gs.MakeInvestmentWithTerms(s.selectedIdx, s.investAmount, s.selectedTerms)
	if err != nil {
		s.errorMsg = err.Error()
		s.phase = PhaseAmountInput
		return s, textinput.Blink
	}

	// Initialize founder relationship with DD bonus
	if len(gs.Portfolio.Investments) > 0 {
		lastInv := &gs.Portfolio.Investments[len(gs.Portfolio.Investments)-1]
		lastInv.FounderName = game.GenerateFounderName()
		lastInv.HasDueDiligence = s.ddLevel != "none"
		lastInv.RelationshipScore = game.CalculateInitialRelationship(s.selectedTerms, lastInv.HasDueDiligence, s.investAmount)
		lastInv.LastInteraction = gs.Portfolio.Turn
	}

	// Reset state
	s.phase = PhaseStartupList
	s.selectedStartup = nil
	s.investAmount = 0
	s.errorMsg = ""
	s.ddFindings = nil
	s.ddLevel = ""
	s.refreshStartupTable()

	// Check if out of money
	if gs.Portfolio.Cash < 10000 {
		gs.AIPlayerMakeInvestments()
		return s, SwitchTo(ScreenVCTurn)
	}

	return s, nil
}

func (s *VCInvestScreen) buildSyndicateTable() {
	gs := s.gameData.GameState

	rows := make([]table.Row, len(gs.SyndicateOpportunities))
	for i, opp := range gs.SyndicateOpportunities {
		rows[i] = table.Row{
			fmt.Sprintf("%d", i+1),
			truncate(opp.CompanyName, 12),
			truncate(opp.LeadInvestor, 10),
			formatCompactMoney(opp.TotalRoundSize),
			fmt.Sprintf("$%s-$%s", formatCompactMoney(opp.YourMinShare), formatCompactMoney(opp.YourMaxShare)),
		}
	}

	columns := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Company", Width: 12},
		{Title: "Lead", Width: 10},
		{Title: "Round", Width: 8},
		{Title: "Your Range", Width: 15},
	}

	s.syndicateTable = components.NewGameTable("", columns, rows)
	s.syndicateTable.SetSize(55, 10)
}

func (s *VCInvestScreen) handleSyndicateAmountSubmit() (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState
	amountStr := strings.TrimSpace(s.amountInput.Value())

	if amountStr == "" || amountStr == "0" {
		s.phase = PhaseSyndicateList
		return s, nil
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		s.errorMsg = "Invalid amount"
		return s, nil
	}

	opp := gs.SyndicateOpportunities[s.selectedSyndicate]

	if amount < opp.YourMinShare {
		s.errorMsg = fmt.Sprintf("Minimum investment is $%d", opp.YourMinShare)
		return s, nil
	}
	if amount > opp.YourMaxShare {
		s.errorMsg = fmt.Sprintf("Maximum investment is $%d", opp.YourMaxShare)
		return s, nil
	}
	if amount > gs.Portfolio.Cash {
		s.errorMsg = fmt.Sprintf("Insufficient funds (have $%d)", gs.Portfolio.Cash)
		return s, nil
	}

	err = gs.MakeSyndicateInvestment(s.selectedSyndicate, amount)
	if err != nil {
		s.errorMsg = err.Error()
		return s, nil
	}

	// Reset and go back
	s.phase = PhaseStartupList
	s.errorMsg = ""
	s.refreshStartupTable()
	s.buildSyndicateTable()

	return s, nil
}

func (s *VCInvestScreen) makeInvestment(terms game.InvestmentTerms) (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState

	err := gs.MakeInvestmentWithTerms(s.selectedIdx, s.investAmount, terms)
	if err != nil {
		s.errorMsg = err.Error()
		s.phase = PhaseAmountInput
		return s, textinput.Blink
	}

	// Initialize founder relationship
	if len(gs.Portfolio.Investments) > 0 {
		lastInv := &gs.Portfolio.Investments[len(gs.Portfolio.Investments)-1]
		lastInv.FounderName = game.GenerateFounderName()
		lastInv.RelationshipScore = game.CalculateInitialRelationship(terms, false, s.investAmount)
		lastInv.LastInteraction = gs.Portfolio.Turn
	}

	// Reset and go back to startup list
	s.phase = PhaseStartupList
	s.selectedStartup = nil
	s.investAmount = 0
	s.errorMsg = ""
	s.refreshStartupTable()

	// Check if out of money
	if gs.Portfolio.Cash < 10000 {
		gs.AIPlayerMakeInvestments()
		return s, SwitchTo(ScreenVCTurn)
	}

	return s, nil
}

// View renders the investment screen
func (s *VCInvestScreen) View() string {
	gs := s.gameData.GameState
	var b strings.Builder

	// Header
	header := s.renderHeader()
	b.WriteString(header)
	b.WriteString("\n")

	// Main content based on phase
	switch s.phase {
	case PhaseStartupList:
		b.WriteString(s.renderStartupList())
	case PhaseAmountInput:
		b.WriteString(s.renderAmountInput())
	case PhaseTermsSelect:
		b.WriteString(s.renderTermsSelect())
	case PhaseDueDiligence:
		b.WriteString(s.renderDueDiligence())
	case PhaseDDResults:
		b.WriteString(s.renderDDResults())
	case PhaseSyndicateList:
		b.WriteString(s.renderSyndicateList())
	case PhaseSyndicateAmount:
		b.WriteString(s.renderSyndicateAmount())
	}

	// Status bar - but don't show turn screen help
	b.WriteString("\n")
	statusBar := components.GameStatusBar(
		s.width,
		gs.Portfolio.Turn,
		gs.Portfolio.MaxTurns,
		gs.Portfolio.Cash,
		gs.Portfolio.NetWorth,
	)
	// Override help text for investment phase
	statusBar.SetShowHelp(false)
	b.WriteString(statusBar.View())

	return b.String()
}

func (s *VCInvestScreen) renderHeader() string {
	gs := s.gameData.GameState

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(70).
		Align(lipgloss.Center).
		Padding(0, 2)

	title := fmt.Sprintf("🦄 %s - INVESTMENT PHASE", s.gameData.FirmName)

	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render(title)))
	b.WriteString("\n")

	// Fund info
	infoStyle := lipgloss.NewStyle().
		Foreground(styles.Yellow).
		Width(s.width).
		Align(lipgloss.Center)
	info := fmt.Sprintf("Fund: $%s | Cash: $%s | Investments: %d",
		formatCompactMoney(gs.Portfolio.InitialFundSize),
		formatCompactMoney(gs.Portfolio.Cash),
		len(gs.Portfolio.Investments))
	b.WriteString(infoStyle.Render(info))

	return b.String()
}

func (s *VCInvestScreen) renderStartupList() string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Cyan).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(titleStyle.Render("AVAILABLE STARTUPS"))
	b.WriteString("\n\n")

	// Table
	tableContainer := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(tableContainer.Render(s.startupTable.View()))
	b.WriteString("\n\n")

	// Current investments
	gs := s.gameData.GameState
	if len(gs.Portfolio.Investments) > 0 {
		investStyle := lipgloss.NewStyle().
			Foreground(styles.Green).
			Width(s.width).
			Align(lipgloss.Center)
		b.WriteString(investStyle.Render(fmt.Sprintf("✓ %d investments made", len(gs.Portfolio.Investments))))
		b.WriteString("\n")
	}

	// Syndicate hint
	if len(gs.SyndicateOpportunities) > 0 {
		syndicateStyle := lipgloss.NewStyle().
			Foreground(styles.Magenta).
			Width(s.width).
			Align(lipgloss.Center)
		b.WriteString(syndicateStyle.Render(fmt.Sprintf("🤝 %d syndicate opportunities available (press 's')", len(gs.SyndicateOpportunities))))
		b.WriteString("\n")
	}

	// Help
	helpStyle := lipgloss.NewStyle().
		Foreground(styles.Gray).
		Width(s.width).
		Align(lipgloss.Center)

	helpText := "↑/↓ select • enter invest • d done • esc back"
	if len(gs.SyndicateOpportunities) > 0 {
		helpText = "↑/↓ select • enter invest • s syndicate • d done • esc back"
	}
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

func (s *VCInvestScreen) renderAmountInput() string {
	gs := s.gameData.GameState
	startup := s.selectedStartup
	var b strings.Builder

	// Startup details box
	detailBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(60)

	var details strings.Builder
	nameStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	details.WriteString(nameStyle.Render(startup.Name))
	details.WriteString("\n")
	details.WriteString(startup.Description)
	details.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	details.WriteString(labelStyle.Render("Category: "))
	details.WriteString(startup.Category)
	details.WriteString("\n")
	details.WriteString(labelStyle.Render("Valuation: "))
	details.WriteString(fmt.Sprintf("$%s", formatCompactMoney(startup.Valuation)))
	details.WriteString("\n")

	maxInvest := int64(float64(startup.Valuation) * 0.20)
	if maxInvest > gs.Portfolio.Cash {
		maxInvest = gs.Portfolio.Cash
	}
	details.WriteString(labelStyle.Render("Max Investment: "))
	details.WriteString(fmt.Sprintf("$%d (20%% of valuation)", maxInvest))

	detailContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	b.WriteString(detailContainer.Render(detailBox.Render(details.String())))
	b.WriteString("\n\n")

	// Amount input
	inputLabel := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputLabel.Render("INVESTMENT AMOUNT")))
	b.WriteString("\n")

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(25)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render("$ " + s.amountInput.View())))
	b.WriteString("\n")

	// Error message
	if s.errorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(styles.Red).Width(s.width).Align(lipgloss.Center)
		b.WriteString(errStyle.Render("⚠ " + s.errorMsg))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter confirm • esc cancel • 0 skip"))

	return b.String()
}

func (s *VCInvestScreen) renderTermsSelect() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Cyan).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(titleStyle.Render(fmt.Sprintf("INVESTING $%s IN %s", formatCompactMoney(s.investAmount), s.selectedStartup.Name)))
	b.WriteString("\n\n")

	if s.termsMenu != nil {
		menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
		menuBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Cyan).
			Padding(1, 2)
		b.WriteString(menuContainer.Render(menuBox.Render(s.termsMenu.View())))
	}

	return b.String()
}

func (s *VCInvestScreen) renderDueDiligence() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Cyan).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)

	b.WriteString(titleStyle.Render(fmt.Sprintf("🔍 DUE DILIGENCE - %s", s.selectedStartup.Name)))
	b.WriteString("\n\n")

	// Info box
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2).
		Width(60)

	info := fmt.Sprintf("Investment: $%s in %s\n\nDue diligence can reveal red flags, hidden gems, and improve founder relationships.",
		formatCompactMoney(s.investAmount), s.selectedStartup.Name)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(infoBox.Render(info)))
	b.WriteString("\n\n")

	// DD Menu
	if s.ddMenu != nil {
		menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
		menuBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Cyan).
			Padding(1, 2)
		b.WriteString(menuContainer.Render(menuBox.Render(s.ddMenu.View())))
	}

	// Error
	if s.errorMsg != "" {
		b.WriteString("\n")
		errStyle := lipgloss.NewStyle().Foreground(styles.Red).Width(s.width).Align(lipgloss.Center)
		b.WriteString(errStyle.Render(s.errorMsg))
	}

	b.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter select • esc back"))

	return b.String()
}

func (s *VCInvestScreen) renderDDResults() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Cyan).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)

	b.WriteString(titleStyle.Render("🔍 DUE DILIGENCE FINDINGS"))
	b.WriteString("\n\n")

	// Findings box
	findingsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(65)

	var findings strings.Builder
	for _, finding := range s.ddFindings {
		switch finding.Type {
		case "red_flag":
			redStyle := lipgloss.NewStyle().Foreground(styles.Red)
			findings.WriteString(redStyle.Render(fmt.Sprintf("🚩 RED FLAG - %s: %s", finding.Category, finding.Description)))
		case "green_flag":
			greenStyle := lipgloss.NewStyle().Foreground(styles.Green)
			findings.WriteString(greenStyle.Render(fmt.Sprintf("✓ POSITIVE - %s: %s", finding.Category, finding.Description)))
		default:
			findings.WriteString(fmt.Sprintf("ℹ️  INFO - %s: %s", finding.Category, finding.Description))
		}
		findings.WriteString("\n")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(findingsBox.Render(findings.String())))
	b.WriteString("\n")

	// Warning if should block
	if s.ddShouldBlock {
		warnBox := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(styles.Red).
			Foreground(styles.Red).
			Padding(0, 2).
			Width(50)
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(warnBox.Render("⚠️  WARNING: " + s.ddBlockReason)))
	}

	b.WriteString("\n\n")

	// Confirm
	confirmStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Width(s.width).Align(lipgloss.Center)
	if s.ddShouldBlock {
		b.WriteString(confirmStyle.Render("Still proceed with investment? [y]es / [n]o"))
	} else {
		b.WriteString(confirmStyle.Render("Proceed with investment? [y]es / [n]o"))
	}

	return b.String()
}

func (s *VCInvestScreen) renderSyndicateList() string {
	gs := s.gameData.GameState
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Magenta).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)

	b.WriteString(titleStyle.Render("🤝 SYNDICATE OPPORTUNITIES"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center).Render("Co-invest with other VCs for better terms"))
	b.WriteString("\n\n")

	if len(gs.SyndicateOpportunities) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center).Render("No syndicate opportunities available"))
	} else {
		// Syndicate table
		tableContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
		b.WriteString(tableContainer.Render(s.syndicateTable.View()))
	}

	b.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("↑/↓ select • enter invest • esc back"))

	return b.String()
}

func (s *VCInvestScreen) renderSyndicateAmount() string {
	gs := s.gameData.GameState
	opp := gs.SyndicateOpportunities[s.selectedSyndicate]
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Magenta).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)

	b.WriteString(titleStyle.Render(fmt.Sprintf("🤝 SYNDICATE: %s", opp.CompanyName)))
	b.WriteString("\n\n")

	// Details box
	detailBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2).
		Width(55)

	var details strings.Builder
	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	details.WriteString(labelStyle.Render("Lead Investor: "))
	details.WriteString(fmt.Sprintf("%s (%s)\n", opp.LeadInvestor, opp.LeadInvestorFirm))
	details.WriteString(labelStyle.Render("Total Round: "))
	details.WriteString(fmt.Sprintf("$%s\n", formatCompactMoney(opp.TotalRoundSize)))
	details.WriteString(labelStyle.Render("Valuation: "))
	details.WriteString(fmt.Sprintf("$%s\n", formatCompactMoney(opp.Valuation)))
	details.WriteString(labelStyle.Render("Your Range: "))
	greenStyle := lipgloss.NewStyle().Foreground(styles.Green)
	details.WriteString(greenStyle.Render(fmt.Sprintf("$%s - $%s\n", formatCompactMoney(opp.YourMinShare), formatCompactMoney(opp.YourMaxShare))))
	details.WriteString("\n")
	details.WriteString(labelStyle.Render("Benefits:\n"))
	for _, benefit := range opp.Benefits {
		details.WriteString(fmt.Sprintf("  • %s\n", benefit))
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(detailBox.Render(details.String())))
	b.WriteString("\n\n")

	// Amount input
	inputLabel := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputLabel.Render("INVESTMENT AMOUNT")))
	b.WriteString("\n")

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(30)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render("$ " + s.amountInput.View())))
	b.WriteString("\n")

	// Error
	if s.errorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(styles.Red).Width(s.width).Align(lipgloss.Center)
		b.WriteString(errStyle.Render(s.errorMsg))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter confirm • esc cancel"))

	return b.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatCompactMoney(amount int64) string {
	if amount >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(amount)/1000000000)
	}
	if amount >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(amount)/1000000)
	}
	if amount >= 1000 {
		return fmt.Sprintf("%.0fK", float64(amount)/1000)
	}
	return fmt.Sprintf("%d", amount)
}

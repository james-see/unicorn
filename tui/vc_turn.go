package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/game"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// TurnView represents what we're currently showing
type TurnView int

const (
	ViewTurnSummary TurnView = iota
	ViewDashboard
	ViewValueAdd
	ViewSecondaryMarket
	ViewFollowOn       // Follow-on investment opportunity
	ViewBoardVote      // Board vote required
	ViewFollowOnAmount // Entering follow-on amount
)

// VCTurnScreen handles the main game turn loop
type VCTurnScreen struct {
	width    int
	height   int
	gameData *GameData
	view     TurnView

	// Turn state
	turnMessages     []string
	portfolioTable   *components.GameTable
	leaderboardTable *components.GameTable

	// Follow-on state
	followOnOpps    []game.FollowOnOpportunity
	currentFollowOn int
	followOnAmount  textinput.Model
	followOnMsg     string

	// Board vote state
	pendingVotes []game.BoardVote
	currentVote  int
	boardVoteMsg string

	// Value-add state
	valueAddCompany string // Selected company for value-add
	valueAddPhase   int    // 0 = select company, 1 = select action
	valueAddMsg     string // Feedback message

	// Auto mode
	autoTicker *time.Ticker
}

// NewVCTurnScreen creates a new turn screen
func NewVCTurnScreen(width, height int, gameData *GameData) *VCTurnScreen {
	// Follow-on amount input
	followOnInput := textinput.New()
	followOnInput.Placeholder = "Enter amount"
	followOnInput.CharLimit = 15
	followOnInput.Width = 20

	s := &VCTurnScreen{
		width:          width,
		height:         height,
		gameData:       gameData,
		view:           ViewTurnSummary,
		followOnAmount: followOnInput,
	}

	s.refreshPortfolioTable()
	s.refreshLeaderboard()

	return s
}

func (s *VCTurnScreen) refreshPortfolioTable() {
	gs := s.gameData.GameState

	rows := make([]table.Row, len(gs.Portfolio.Investments))
	for i, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		profit := value - inv.AmountInvested

		profitStr := fmt.Sprintf("+$%s", formatCompactMoney(profit))
		if profit < 0 {
			profitStr = fmt.Sprintf("-$%s", formatCompactMoney(-profit))
		}

		rows[i] = table.Row{
			truncate(inv.CompanyName, 18),
			fmt.Sprintf("$%s", formatCompactMoney(inv.AmountInvested)),
			fmt.Sprintf("$%s", formatCompactMoney(value)),
			fmt.Sprintf("%.1f%%", inv.EquityPercent),
			profitStr,
		}
	}

	columns := []table.Column{
		{Title: "Company", Width: 18},
		{Title: "Invested", Width: 10},
		{Title: "Value", Width: 10},
		{Title: "Equity", Width: 8},
		{Title: "P/L", Width: 10},
	}

	s.portfolioTable = components.NewGameTable("", columns, rows)
	s.portfolioTable.SetSize(60, 10)
}

func (s *VCTurnScreen) refreshLeaderboard() {
	gs := s.gameData.GameState
	leaderboard := gs.GetLeaderboard()

	rows := make([]table.Row, len(leaderboard))
	for i, entry := range leaderboard {
		marker := ""
		if entry.IsPlayer {
			marker = "→"
		}

		rows[i] = table.Row{
			fmt.Sprintf("%s%d", marker, i+1),
			truncate(entry.Name, 15),
			fmt.Sprintf("$%s", formatCompactMoney(entry.NetWorth)),
			fmt.Sprintf("%.0f%%", entry.ROI),
		}
	}

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Investor", Width: 15},
		{Title: "Net Worth", Width: 12},
		{Title: "ROI", Width: 8},
	}

	s.leaderboardTable = components.NewGameTable("", columns, rows)
	s.leaderboardTable.SetSize(45, 8)
}

// tickMsg is sent for auto mode
type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the turn screen
func (s *VCTurnScreen) Init() tea.Cmd {
	// Process first turn
	s.processTurn()

	if s.gameData.AutoMode {
		return doTick()
	}
	return nil
}

// Update handles turn screen input
func (s *VCTurnScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s.view {
		case ViewTurnSummary:
			switch {
			case key.Matches(msg, keys.Global.Enter):
				if gs.IsGameOver() {
					return s, SwitchTo(ScreenVCResults)
				}
				s.processTurn()
				s.refreshPortfolioTable()
				s.refreshLeaderboard()
				return s, nil

			case msg.String() == "d":
				s.view = ViewDashboard
				return s, nil

			case msg.String() == "v":
				s.view = ViewValueAdd
				return s, nil

			case msg.String() == "s":
				s.view = ViewSecondaryMarket
				return s, nil

			case key.Matches(msg, keys.Global.Back), msg.String() == "q":
				// Quit confirmation would go here
				return s, SwitchTo(ScreenMainMenu)
			}

		case ViewDashboard, ViewSecondaryMarket:
			if key.Matches(msg, keys.Global.Back) {
				s.view = ViewTurnSummary
				return s, nil
			}

		case ViewValueAdd:
			if key.Matches(msg, keys.Global.Back) {
				if s.valueAddPhase == 1 {
					// Go back to company selection
					s.valueAddPhase = 0
					s.valueAddCompany = ""
					return s, nil
				}
				s.view = ViewTurnSummary
				s.valueAddPhase = 0
				s.valueAddCompany = ""
				s.valueAddMsg = ""
				return s, nil
			}

			// Handle number key presses
			keyStr := msg.String()
			if len(keyStr) == 1 && keyStr[0] >= '1' && keyStr[0] <= '9' {
				num := int(keyStr[0] - '0')
				return s.handleValueAddSelection(num)
			}

		case ViewFollowOn:
			switch {
			case msg.String() == "y" || msg.String() == "i":
				// Invest - show amount input
				s.view = ViewFollowOnAmount
				s.followOnAmount.Focus()
				s.followOnAmount.SetValue("")
				return s, textinput.Blink
			case msg.String() == "n" || msg.String() == "s":
				// Skip this opportunity
				return s.skipFollowOn()
			}

		case ViewFollowOnAmount:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.view = ViewFollowOn
				s.followOnMsg = ""
				return s, nil
			case key.Matches(msg, keys.Global.Enter):
				return s.handleFollowOnInvest()
			}

		case ViewBoardVote:
			switch {
			case msg.String() == "a" || msg.String() == "1":
				return s.handleBoardVote(true) // Accept/Approve
			case msg.String() == "b" || msg.String() == "2":
				return s.handleBoardVote(false) // Reject/Decline
			}
		}

	case tickMsg:
		if s.gameData.AutoMode && !gs.IsGameOver() {
			s.processTurn()
			s.refreshPortfolioTable()
			s.refreshLeaderboard()

			if gs.IsGameOver() {
				return s, SwitchTo(ScreenVCResults)
			}
			return s, doTick()
		}
	}

	// Update sub-components
	var cmd tea.Cmd
	switch s.view {
	case ViewTurnSummary:
		s.portfolioTable, cmd = s.portfolioTable.Update(msg)
	case ViewDashboard:
		s.portfolioTable, cmd = s.portfolioTable.Update(msg)
	case ViewFollowOnAmount:
		s.followOnAmount, cmd = s.followOnAmount.Update(msg)
	}

	return s, cmd
}

func (s *VCTurnScreen) skipFollowOn() (ScreenModel, tea.Cmd) {
	s.currentFollowOn++
	if s.currentFollowOn >= len(s.followOnOpps) {
		// Done with all follow-ons, continue to process turn
		s.followOnOpps = nil
		s.view = ViewTurnSummary
		s.continueProcessTurn()
	}
	return s, nil
}

func (s *VCTurnScreen) handleFollowOnInvest() (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState
	opp := s.followOnOpps[s.currentFollowOn]

	amountStr := strings.TrimSpace(s.followOnAmount.Value())
	if amountStr == "" || amountStr == "0" {
		return s.skipFollowOn()
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount < 0 {
		s.followOnMsg = "Invalid amount"
		return s, nil
	}

	if amount < opp.MinInvestment {
		s.followOnMsg = fmt.Sprintf("Minimum investment is $%d", opp.MinInvestment)
		return s, nil
	}

	if amount > opp.MaxInvestment {
		s.followOnMsg = fmt.Sprintf("Maximum investment is $%d", opp.MaxInvestment)
		return s, nil
	}

	err = gs.MakeFollowOnInvestment(opp.CompanyName, amount)
	if err != nil {
		s.followOnMsg = err.Error()
		return s, nil
	}

	s.followOnMsg = fmt.Sprintf("✓ Invested $%s in %s!", formatCompactMoney(amount), opp.CompanyName)
	s.currentFollowOn++

	if s.currentFollowOn >= len(s.followOnOpps) {
		// Done with all follow-ons
		s.followOnOpps = nil
		s.view = ViewTurnSummary
		s.continueProcessTurn()
	} else {
		s.view = ViewFollowOn
		s.followOnAmount.SetValue("")
	}

	return s, nil
}

func (s *VCTurnScreen) handleBoardVote(approve bool) (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState

	if s.currentVote >= len(s.pendingVotes) {
		s.view = ViewTurnSummary
		return s, nil
	}

	vote := s.pendingVotes[s.currentVote]

	// Find the vote index in pending votes
	voteIndex := -1
	for i, v := range gs.GetPendingBoardVotes() {
		if v.CompanyName == vote.CompanyName && v.VoteType == vote.VoteType {
			voteIndex = i
			break
		}
	}

	if voteIndex == -1 {
		s.currentVote++
		if s.currentVote >= len(s.pendingVotes) {
			s.view = ViewTurnSummary
		}
		return s, nil
	}

	voteChoice := "B"
	if approve {
		voteChoice = "A"
	}

	result, passed, err := gs.ProcessBoardVote(voteIndex, voteChoice)
	if err != nil {
		s.boardVoteMsg = err.Error()
		return s, nil
	}

	// Execute outcome
	outcomeMessages := gs.ExecuteBoardVoteOutcome(vote, passed)
	s.turnMessages = append(s.turnMessages, result)
	s.turnMessages = append(s.turnMessages, outcomeMessages...)

	s.currentVote++
	if s.currentVote >= len(s.pendingVotes) {
		s.pendingVotes = nil
		s.view = ViewTurnSummary
		s.refreshPortfolioTable()
		s.refreshLeaderboard()
	}

	return s, nil
}

func (s *VCTurnScreen) continueProcessTurn() {
	gs := s.gameData.GameState

	// Process the turn
	messages := gs.ProcessTurn()

	// Process value-add actions
	valueAddMsgs := gs.ProcessActiveValueAddActions()
	messages = append(messages, valueAddMsgs...)

	// Process relationships
	for i := range gs.Portfolio.Investments {
		inv := &gs.Portfolio.Investments[i]
		if inv.FounderName != "" {
			event := game.GenerateRelationshipEvent(inv, gs.Portfolio.Turn)
			if event != nil {
				inv.RelationshipScore = game.ApplyRelationshipChange(
					inv.RelationshipScore,
					event.ScoreChange)
				messages = append(messages, event.Description)
			}
		}
	}

	// Generate secondary offers
	newOffers := gs.GenerateSecondaryOffers()
	gs.SecondaryMarketOffers = append(gs.SecondaryMarketOffers, newOffers...)

	// Process expirations
	expiredMsgs := gs.ProcessSecondaryOfferExpirations()
	messages = append(messages, expiredMsgs...)

	// Check for board votes
	pendingVotes := gs.GetPendingBoardVotes()
	if len(pendingVotes) > 0 {
		s.pendingVotes = pendingVotes
		s.currentVote = 0
		s.boardVoteMsg = ""
		s.view = ViewBoardVote
	}

	s.turnMessages = messages
	s.refreshPortfolioTable()
	s.refreshLeaderboard()
}

func (s *VCTurnScreen) handleValueAddSelection(num int) (ScreenModel, tea.Cmd) {
	gs := s.gameData.GameState
	opportunities := gs.GetValueAddOpportunities()

	if s.valueAddPhase == 0 {
		// Phase 0: Selecting a company
		if num < 1 || num > len(opportunities) {
			s.valueAddMsg = "Invalid company selection"
			return s, nil
		}

		// Check if we have attention points
		actionsThisTurn := 0
		for _, action := range gs.ActiveValueAddActions {
			if action.AppliedTurn == gs.Portfolio.Turn {
				actionsThisTurn++
			}
		}
		if actionsThisTurn >= 2 {
			s.valueAddMsg = "No attention points remaining this turn"
			return s, nil
		}

		s.valueAddCompany = opportunities[num-1]
		s.valueAddPhase = 1
		s.valueAddMsg = ""
		return s, nil

	} else {
		// Phase 1: Selecting an action type
		actionTypes := game.GetAvailableValueAddTypes()
		if num < 1 || num > len(actionTypes) {
			s.valueAddMsg = "Invalid action selection"
			return s, nil
		}

		selectedAction := actionTypes[num-1]

		// Check cost
		if selectedAction.Cost > gs.Portfolio.Cash {
			s.valueAddMsg = fmt.Sprintf("Insufficient funds (need $%s)", formatCompactMoney(selectedAction.Cost))
			return s, nil
		}

		// Find the investment
		var inv *game.Investment
		for i := range gs.Portfolio.Investments {
			if gs.Portfolio.Investments[i].CompanyName == s.valueAddCompany {
				inv = &gs.Portfolio.Investments[i]
				break
			}
		}

		if inv == nil {
			s.valueAddMsg = "Company not found"
			return s, nil
		}

		// Check requirements
		if selectedAction.RequiresBoardSeat && !inv.Terms.HasBoardSeat {
			s.valueAddMsg = "This action requires a board seat"
			return s, nil
		}
		if inv.EquityPercent < selectedAction.MinEquityPct {
			s.valueAddMsg = fmt.Sprintf("Need %.1f%% equity (have %.1f%%)", selectedAction.MinEquityPct, inv.EquityPercent)
			return s, nil
		}

		// Apply the action
		gs.Portfolio.Cash -= selectedAction.Cost

		action := game.ValueAddAction{
			ActionType:     selectedAction.ID,
			CompanyName:    s.valueAddCompany,
			Cost:           selectedAction.Cost,
			AppliedTurn:    gs.Portfolio.Turn,
			Duration:       selectedAction.Duration,
			ValuationBoost: selectedAction.MinValBoost,
			RiskReduction:  selectedAction.RiskReduction,
			Description:    selectedAction.Description,
		}
		gs.ActiveValueAddActions = append(gs.ActiveValueAddActions, action)

		// Improve relationship
		inv.RelationshipScore = game.ApplyRelationshipChange(inv.RelationshipScore, 5.0)

		s.valueAddMsg = fmt.Sprintf("✓ Applied %s to %s!", selectedAction.Name, s.valueAddCompany)
		s.valueAddPhase = 0
		s.valueAddCompany = ""
		return s, nil
	}
}

func (s *VCTurnScreen) processTurn() {
	gs := s.gameData.GameState

	// Check for follow-on opportunities BEFORE processing turn
	opportunities := gs.GetFollowOnOpportunities()
	if len(opportunities) > 0 {
		s.followOnOpps = opportunities
		s.currentFollowOn = 0
		s.followOnMsg = ""
		s.followOnAmount.SetValue("")
		s.view = ViewFollowOn
		return // Don't process turn yet - wait for follow-on decisions
	}

	// Process the turn
	messages := gs.ProcessTurn()

	// Process value-add actions
	valueAddMsgs := gs.ProcessActiveValueAddActions()
	messages = append(messages, valueAddMsgs...)

	// Process relationships
	for i := range gs.Portfolio.Investments {
		inv := &gs.Portfolio.Investments[i]
		if inv.FounderName != "" {
			event := game.GenerateRelationshipEvent(inv, gs.Portfolio.Turn)
			if event != nil {
				inv.RelationshipScore = game.ApplyRelationshipChange(
					inv.RelationshipScore,
					event.ScoreChange)
				messages = append(messages, event.Description)
			}
		}
	}

	// Generate secondary offers
	newOffers := gs.GenerateSecondaryOffers()
	gs.SecondaryMarketOffers = append(gs.SecondaryMarketOffers, newOffers...)

	// Process expirations
	expiredMsgs := gs.ProcessSecondaryOfferExpirations()
	messages = append(messages, expiredMsgs...)

	// Check for board votes AFTER processing turn
	pendingVotes := gs.GetPendingBoardVotes()
	if len(pendingVotes) > 0 {
		s.pendingVotes = pendingVotes
		s.currentVote = 0
		s.boardVoteMsg = ""
		s.view = ViewBoardVote
	}

	s.turnMessages = messages
}

// View renders the turn screen
func (s *VCTurnScreen) View() string {
	switch s.view {
	case ViewDashboard:
		return s.renderDashboard()
	case ViewValueAdd:
		return s.renderValueAdd()
	case ViewSecondaryMarket:
		return s.renderSecondaryMarket()
	case ViewFollowOn, ViewFollowOnAmount:
		return s.renderFollowOn()
	case ViewBoardVote:
		return s.renderBoardVote()
	default:
		return s.renderTurnSummary()
	}
}

func (s *VCTurnScreen) renderTurnSummary() string {
	gs := s.gameData.GameState
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	header := fmt.Sprintf("🦄 %s - MONTH %d/%d", s.gameData.FirmName, gs.Portfolio.Turn, gs.Portfolio.MaxTurns)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render(header)))
	b.WriteString("\n\n")

	// Main layout: Portfolio on left, Standings on right
	leftPanel := s.renderPortfolioPanel()
	rightPanel := s.renderStandingsPanel()

	// Join panels side by side
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, "  ", rightPanel)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(content))
	b.WriteString("\n")

	// News as full-width panel below
	b.WriteString(s.renderNewsPanel())
	b.WriteString("\n")

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Foreground(styles.White).
		Background(styles.DarkGray).
		Width(70).
		Padding(0, 2)

	// Calculate ROI
	invested := gs.GetTotalInvested()
	roi := 0.0
	if invested > 0 {
		roi = float64(gs.Portfolio.NetWorth-gs.Portfolio.InitialFundSize) / float64(gs.Portfolio.InitialFundSize) * 100
	}

	status := fmt.Sprintf("💰 Cash: $%s  |  📊 Net Worth: $%s  |  📈 ROI: %.1f%%",
		formatCompactMoney(gs.Portfolio.Cash),
		formatCompactMoney(gs.Portfolio.NetWorth),
		roi)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(statusStyle.Render(status)))
	b.WriteString("\n\n")

	// Help text
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	if gs.IsGameOver() {
		b.WriteString(helpStyle.Render("🏁 GAME OVER - Press Enter to see results"))
	} else {
		b.WriteString(helpStyle.Render("enter next • d dashboard • v value-add • s secondary • q quit"))
	}

	return b.String()
}

func (s *VCTurnScreen) renderPortfolioPanel() string {
	gs := s.gameData.GameState

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(35)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	b.WriteString(titleStyle.Render("📊 PORTFOLIO"))
	b.WriteString("\n")

	if len(gs.Portfolio.Investments) == 0 {
		b.WriteString("\n  No investments yet\n")
	} else {
		for _, inv := range gs.Portfolio.Investments {
			value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
			profit := value - inv.AmountInvested

			// Company name
			nameStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
			b.WriteString(nameStyle.Render(truncate(inv.CompanyName, 20)))

			// Profit/loss indicator
			if profit >= 0 {
				profitStyle := lipgloss.NewStyle().Foreground(styles.Green)
				b.WriteString(profitStyle.Render(fmt.Sprintf(" +$%s", formatCompactMoney(profit))))
			} else {
				lossStyle := lipgloss.NewStyle().Foreground(styles.Red)
				b.WriteString(lossStyle.Render(fmt.Sprintf(" -$%s", formatCompactMoney(-profit))))
			}
			b.WriteString("\n")

			// Details
			detailStyle := lipgloss.NewStyle().Foreground(styles.Gray)
			b.WriteString(detailStyle.Render(fmt.Sprintf("  $%s → $%s (%.1f%%)",
				formatCompactMoney(inv.AmountInvested),
				formatCompactMoney(value),
				inv.EquityPercent)))
			b.WriteString("\n")
		}
	}

	return panelStyle.Render(b.String())
}

func (s *VCTurnScreen) renderNewsPanel() string {
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(0, 1).
		Width(70)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(titleStyle.Render("📰 NEWS"))
	b.WriteString("\n")

	if len(s.turnMessages) == 0 {
		b.WriteString("  No news this turn\n")
	} else {
		// Show last 4 messages
		start := 0
		if len(s.turnMessages) > 4 {
			start = len(s.turnMessages) - 4
		}

		for _, msg := range s.turnMessages[start:] {
			// Color code based on content
			msgStyle := lipgloss.NewStyle().Foreground(styles.White)
			if strings.Contains(msg, "raised") || strings.Contains(msg, "📈") || strings.Contains(msg, "increased") || strings.Contains(msg, "growth") {
				msgStyle = lipgloss.NewStyle().Foreground(styles.Green)
			} else if strings.Contains(msg, "down") || strings.Contains(msg, "📉") || strings.Contains(msg, "decreased") || strings.Contains(msg, "scandal") {
				msgStyle = lipgloss.NewStyle().Foreground(styles.Red)
			}

			// Allow longer messages now that we have full width
			displayMsg := msg
			if len(displayMsg) > 65 {
				displayMsg = displayMsg[:62] + "..."
			}

			b.WriteString(msgStyle.Render("• " + displayMsg))
			b.WriteString("\n")
		}
	}

	return lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(panelStyle.Render(b.String()))
}

func (s *VCTurnScreen) renderStandingsPanel() string {
	gs := s.gameData.GameState
	leaderboard := gs.GetLeaderboard()

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(0, 1).
		Width(35)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
	b.WriteString(titleStyle.Render("🏆 STANDINGS"))
	b.WriteString("\n")

	for i, entry := range leaderboard {
		marker := "  "
		nameStyle := lipgloss.NewStyle().Foreground(styles.White)
		if entry.IsPlayer {
			marker = "→ "
			nameStyle = lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
		}

		b.WriteString(fmt.Sprintf("%s%d. ", marker, i+1))
		b.WriteString(nameStyle.Render(truncate(entry.Name, 12)))
		b.WriteString(fmt.Sprintf(" $%s (%.0f%%)\n", formatCompactMoney(entry.NetWorth), entry.ROI))
	}

	return panelStyle.Render(b.String())
}

func (s *VCTurnScreen) renderMiniLeaderboard() string {
	gs := s.gameData.GameState
	leaderboard := gs.GetLeaderboard()

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(0, 1).
		Width(50)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
	b.WriteString(titleStyle.Render("🏆 STANDINGS"))
	b.WriteString("\n")

	for i, entry := range leaderboard {
		marker := "  "
		nameStyle := lipgloss.NewStyle().Foreground(styles.White)
		if entry.IsPlayer {
			marker = "→ "
			nameStyle = lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
		}

		b.WriteString(fmt.Sprintf("%s%d. ", marker, i+1))
		b.WriteString(nameStyle.Render(truncate(entry.Name, 12)))
		b.WriteString(fmt.Sprintf(" $%s (%.0f%%)\n", formatCompactMoney(entry.NetWorth), entry.ROI))
	}

	return lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(panelStyle.Render(b.String()))
}

func (s *VCTurnScreen) renderDashboard() string {
	gs := s.gameData.GameState
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📊 PORTFOLIO DASHBOARD")))
	b.WriteString("\n\n")

	// Detailed portfolio table
	tableContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	b.WriteString(tableContainer.Render(s.portfolioTable.View()))
	b.WriteString("\n\n")

	// Summary stats
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2).
		Width(50)

	// Calculate ROI
	dashRoi := 0.0
	if gs.Portfolio.InitialFundSize > 0 {
		dashRoi = float64(gs.Portfolio.NetWorth-gs.Portfolio.InitialFundSize) / float64(gs.Portfolio.InitialFundSize) * 100
	}

	var stats strings.Builder
	stats.WriteString(fmt.Sprintf("Total Invested: $%s\n", formatCompactMoney(gs.GetTotalInvested())))
	stats.WriteString(fmt.Sprintf("Portfolio Value: $%s\n", formatCompactMoney(gs.GetPortfolioValue())))
	stats.WriteString(fmt.Sprintf("Cash: $%s\n", formatCompactMoney(gs.Portfolio.Cash)))
	stats.WriteString(fmt.Sprintf("Follow-on Reserve: $%s\n", formatCompactMoney(gs.Portfolio.FollowOnReserve)))
	stats.WriteString(fmt.Sprintf("Net Worth: $%s\n", formatCompactMoney(gs.Portfolio.NetWorth)))
	stats.WriteString(fmt.Sprintf("ROI: %.1f%%", dashRoi))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(statsBox.Render(stats.String())))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))

	return b.String()
}

func (s *VCTurnScreen) renderValueAdd() string {
	gs := s.gameData.GameState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⚡ VALUE-ADD ACTIONS")))
	b.WriteString("\n\n")

	// Get eligible companies
	opportunities := gs.GetValueAddOpportunities()

	// Count actions taken this turn
	actionsThisTurn := 0
	for _, action := range gs.ActiveValueAddActions {
		if action.AppliedTurn == gs.Portfolio.Turn {
			actionsThisTurn++
		}
	}
	remainingPoints := 2 - actionsThisTurn

	// Show feedback message if any
	if s.valueAddMsg != "" {
		msgStyle := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
		if strings.HasPrefix(s.valueAddMsg, "✓") {
			msgStyle = msgStyle.Foreground(styles.Green).Bold(true)
		} else {
			msgStyle = msgStyle.Foreground(styles.Red)
		}
		b.WriteString(msgStyle.Render(s.valueAddMsg))
		b.WriteString("\n\n")
	}

	if len(opportunities) == 0 {
		infoStyle := lipgloss.NewStyle().
			Foreground(styles.Yellow).
			Width(s.width).
			Align(lipgloss.Center)
		b.WriteString(infoStyle.Render("No companies eligible for value-add actions"))
		b.WriteString("\n")
		b.WriteString(infoStyle.Render("(Need ≥5% equity or board seat)"))
		b.WriteString("\n\n")
	} else if s.valueAddPhase == 1 {
		// Phase 1: Select action type for chosen company
		titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true).Width(s.width).Align(lipgloss.Center)
		b.WriteString(titleStyle.Render(fmt.Sprintf("Select action for: %s", s.valueAddCompany)))
		b.WriteString("\n\n")

		// Show action types
		actionBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Magenta).
			Padding(0, 2).
			Width(65)

		var actions strings.Builder
		actions.WriteString(fmt.Sprintf("Attention Points: %d remaining | Cash: $%s\n\n", remainingPoints, formatCompactMoney(gs.Portfolio.Cash)))

		actionTypes := game.GetAvailableValueAddTypes()
		for i, at := range actionTypes {
			numStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
			actions.WriteString(numStyle.Render(fmt.Sprintf("%d. ", i+1)))
			actions.WriteString(fmt.Sprintf("%s ($%s)\n", at.Name, formatCompactMoney(at.Cost)))
			actions.WriteString(fmt.Sprintf("   %s\n", at.Description))
			if at.RequiresBoardSeat {
				actions.WriteString("   ⚠️ Requires board seat\n")
			}
		}

		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(actionBox.Render(actions.String())))
		b.WriteString("\n")

	} else {
		// Phase 0: Select company
		infoStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
		b.WriteString(infoStyle.Render(fmt.Sprintf("Attention Points: %d/2 remaining this turn", remainingPoints)))
		b.WriteString("\n\n")

		// Eligible companies
		companyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Green).
			Padding(0, 2).
			Width(60)

		var companies strings.Builder
		titleStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		companies.WriteString(titleStyle.Render("SELECT A COMPANY"))
		companies.WriteString("\n\n")

		for i, companyName := range opportunities {
			for _, inv := range gs.Portfolio.Investments {
				if inv.CompanyName == companyName {
					emoji := game.GetRelationshipEmoji(inv.RelationshipScore)
					numStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
					nameStyle := lipgloss.NewStyle().Foreground(styles.Cyan)
					companies.WriteString(numStyle.Render(fmt.Sprintf("%d. ", i+1)))
					companies.WriteString(nameStyle.Render(companyName))
					companies.WriteString(fmt.Sprintf(" (%.1f%% equity, %s %.0f)\n", inv.EquityPercent, emoji, inv.RelationshipScore))
					break
				}
			}
		}

		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(companyBox.Render(companies.String())))
		b.WriteString("\n")

		if remainingPoints <= 0 {
			warnStyle := lipgloss.NewStyle().Foreground(styles.Red).Width(s.width).Align(lipgloss.Center)
			b.WriteString(warnStyle.Render("⚠️ No attention points remaining this turn"))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	if s.valueAddPhase == 1 {
		b.WriteString(helpStyle.Render("1-5 select action • esc back to companies"))
	} else {
		b.WriteString(helpStyle.Render("1-9 select company • esc back"))
	}

	return b.String()
}

func (s *VCTurnScreen) renderSecondaryMarket() string {
	gs := s.gameData.GameState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Orange).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("💱 SECONDARY MARKET")))
	b.WriteString("\n\n")

	if len(gs.SecondaryMarketOffers) == 0 {
		infoStyle := lipgloss.NewStyle().
			Foreground(styles.Yellow).
			Width(s.width).
			Align(lipgloss.Center)
		b.WriteString(infoStyle.Render("No secondary market offers available"))
	} else {
		// List offers
		for i, offer := range gs.SecondaryMarketOffers {
			offerStyle := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(styles.Orange).
				Padding(0, 1).
				Width(50)

			offerText := fmt.Sprintf("%d. %s - %.1f%% stake @ $%s",
				i+1, offer.CompanyName, offer.EquityOffered, formatCompactMoney(offer.OfferAmount))

			b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(offerStyle.Render(offerText)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))

	return b.String()
}

func (s *VCTurnScreen) renderFollowOn() string {
	gs := s.gameData.GameState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🚀 FOLLOW-ON INVESTMENT OPPORTUNITY")))
	b.WriteString("\n\n")

	if s.currentFollowOn >= len(s.followOnOpps) {
		return b.String()
	}

	opp := s.followOnOpps[s.currentFollowOn]

	// Info
	infoStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Width(s.width).Align(lipgloss.Center)
	b.WriteString(infoStyle.Render("One of your portfolio companies is raising a new funding round!"))
	b.WriteString("\n")
	b.WriteString(infoStyle.Render("Invest more to avoid dilution and increase ownership."))
	b.WriteString("\n\n")

	// Opportunity details box
	detailBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2).
		Width(60)

	var details strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
	details.WriteString(titleStyle.Render(fmt.Sprintf("🏢 %s", opp.CompanyName)))
	details.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	details.WriteString(labelStyle.Render("Round: "))
	details.WriteString(fmt.Sprintf("%s\n", opp.RoundName))
	details.WriteString(labelStyle.Render("Pre-money Valuation: "))
	details.WriteString(fmt.Sprintf("$%s\n", formatCompactMoney(opp.PreMoneyVal)))
	details.WriteString(labelStyle.Render("Post-money Valuation: "))
	details.WriteString(fmt.Sprintf("$%s\n", formatCompactMoney(opp.PostMoneyVal)))

	equityStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	details.WriteString(labelStyle.Render("Your Current Equity: "))
	details.WriteString(equityStyle.Render(fmt.Sprintf("%.2f%%\n", opp.CurrentEquity)))

	details.WriteString("\n")
	availableFunds := gs.Portfolio.Cash + gs.Portfolio.FollowOnReserve
	fundStyle := lipgloss.NewStyle().Foreground(styles.Green)
	details.WriteString(labelStyle.Render("Available Funds: "))
	details.WriteString(fundStyle.Render(fmt.Sprintf("$%s", formatCompactMoney(availableFunds))))
	details.WriteString(fmt.Sprintf(" (Cash: $%s + Reserve: $%s)\n",
		formatCompactMoney(gs.Portfolio.Cash),
		formatCompactMoney(gs.Portfolio.FollowOnReserve)))

	details.WriteString("\n")
	details.WriteString(labelStyle.Render("Investment Range: "))
	details.WriteString(fmt.Sprintf("$%d - $%d\n", opp.MinInvestment, opp.MaxInvestment))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(detailBox.Render(details.String())))
	b.WriteString("\n")

	// Message
	if s.followOnMsg != "" {
		msgStyle := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
		if strings.HasPrefix(s.followOnMsg, "✓") {
			msgStyle = msgStyle.Foreground(styles.Green).Bold(true)
		} else {
			msgStyle = msgStyle.Foreground(styles.Red)
		}
		b.WriteString(msgStyle.Render(s.followOnMsg))
		b.WriteString("\n")
	}

	// Amount input if in that phase
	if s.view == ViewFollowOnAmount {
		b.WriteString("\n")
		inputLabel := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Width(s.width).Align(lipgloss.Center)
		b.WriteString(inputLabel.Render("ENTER INVESTMENT AMOUNT"))
		b.WriteString("\n")

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(styles.Cyan).
			Padding(0, 1).
			Width(30)

		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render("$ " + s.followOnAmount.View())))
		b.WriteString("\n\n")

		helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
		b.WriteString(helpStyle.Render("enter confirm • esc back"))
	} else {
		b.WriteString("\n")
		helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
		b.WriteString(helpStyle.Render("[i]nvest • [s]kip"))
	}

	return b.String()
}

func (s *VCTurnScreen) renderBoardVote() string {
	gs := s.gameData.GameState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🏛️ BOARD VOTE REQUIRED")))
	b.WriteString("\n\n")

	if s.currentVote >= len(s.pendingVotes) {
		return b.String()
	}

	vote := s.pendingVotes[s.currentVote]

	// Vote details box
	voteBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2).
		Width(65)

	var details strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	details.WriteString(titleStyle.Render(fmt.Sprintf("Company: %s", vote.CompanyName)))
	details.WriteString("\n\n")

	voteTitle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	details.WriteString(voteTitle.Render(vote.Title))
	details.WriteString("\n\n")

	details.WriteString(vote.Description)
	details.WriteString("\n")

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(voteBox.Render(details.String())))
	b.WriteString("\n\n")

	// Options
	optionAStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Green).
		Padding(0, 2).
		Width(60)

	optionBStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Red).
		Padding(0, 2).
		Width(60)

	var optA strings.Builder
	optA.WriteString(lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("[A] " + vote.OptionA))
	optA.WriteString("\n")
	optA.WriteString(lipgloss.NewStyle().Foreground(styles.Gray).Render("→ " + vote.ConsequenceA))

	var optB strings.Builder
	optB.WriteString(lipgloss.NewStyle().Foreground(styles.Red).Bold(true).Render("[B] " + vote.OptionB))
	optB.WriteString("\n")
	optB.WriteString(lipgloss.NewStyle().Foreground(styles.Gray).Render("→ " + vote.ConsequenceB))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(optionAStyle.Render(optA.String())))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(optionBStyle.Render(optB.String())))
	b.WriteString("\n")

	// Voting power
	voteWeight := 1
	for _, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == vote.CompanyName && inv.Terms.HasBoardSeat {
			voteWeight = inv.Terms.BoardSeatMultiplier
			if voteWeight == 0 {
				voteWeight = 1
			}
			break
		}
	}

	b.WriteString("\n")
	powerStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	if voteWeight > 1 {
		b.WriteString(powerStyle.Render(fmt.Sprintf("⚡ Voting Power: %d votes (Double Board Seat!)", voteWeight)))
	} else {
		b.WriteString(powerStyle.Render("Voting Power: 1 vote"))
	}

	// Message
	if s.boardVoteMsg != "" {
		b.WriteString("\n")
		msgStyle := lipgloss.NewStyle().Foreground(styles.Red).Width(s.width).Align(lipgloss.Center)
		b.WriteString(msgStyle.Render(s.boardVoteMsg))
	}

	b.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("Press A to Accept/Approve • Press B to Reject/Decline"))

	return b.String()
}

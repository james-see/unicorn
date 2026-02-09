package tui

import (
	"fmt"
	"math/rand"
	"strconv"
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

// FounderView represents what we're showing in founder mode
type FounderView int

const (
	FounderViewMain FounderView = iota
	FounderViewActions
	FounderViewHiring
	FounderViewHiringMarket
	FounderViewFiring
	FounderViewMarketing
	FounderViewFunding
	FounderViewFundingTerms
	FounderViewPartnership
	FounderViewAffiliate
	FounderViewCompetitors
	FounderViewExpansion
	FounderViewPivot
	FounderViewBoard
	FounderViewBoardAction
	FounderViewBuyback
	FounderViewBuybackConfirm
	FounderViewTeamRoster
	FounderViewCustomers
	FounderViewFinancials
	FounderViewExit
	FounderViewConfirmExit
	// Phase 4: Advanced features
	FounderViewRoadmap
	FounderViewRoadmapStart
	FounderViewSegments
	FounderViewPricing
	FounderViewAcquisitions
	FounderViewPlatform
	FounderViewSecurity
	FounderViewPRCrisis
	FounderViewEconomy
	FounderViewSuccession
	FounderViewSalesPipeline
	// New features for parity
	FounderViewStrategicOpportunity
	FounderViewContentMarketing
	FounderViewCSPlaybooks
	FounderViewCompetitiveIntel
	FounderViewReferralProgram
	FounderViewTechDebt
	FounderViewEndAffiliate
	FounderViewAcquisitionOffer
	FounderViewAdvisorExpertise
	FounderViewAdvisorConfirm
	FounderViewRemoveAdvisor
	FounderViewBoardTable
	FounderViewConfirmQuit
	FounderViewEngineerRealloc
	FounderViewCapTable
	FounderViewExecOffer
)

// FounderGameScreen handles the founder game
type FounderGameScreen struct {
	width    int
	height   int
	gameData *GameData
	view     FounderView

	// Turn state
	turnMessages []string

	// Menus
	actionsMenu      *components.Menu
	hiringMenu       *components.Menu
	firingMenu       *components.Menu
	partnershipMenu  *components.Menu
	fundingMenu      *components.Menu
	fundingTerms     []founder.TermSheetOption
	exitMenu         *components.Menu
	competitorMenu   *components.Menu
	expansionMenu    *components.Menu
	pivotMenu        *components.Menu
	boardMenu        *components.Menu
	buybackMenu      *components.Menu
	competitorAction *components.Menu

	// Input state
	marketingInput textinput.Model
	affiliateInput textinput.Model
	inputMessage   string

	// Selected state
	selectedRole          founder.EmployeeRole
	selectedIsExec        bool
	selectedRoundName     string
	selectedTermIndex     int
	selectedExitType      string
	selectedCompetitorIdx int
	selectedBuybackRound  string
	selectedBuybackEquity float64

	// Market selection for hiring
	marketOptions []string

	// Input for buyback
	buybackInput textinput.Model

	// Phase 4: Advanced menus
	roadmapMenu      *components.Menu
	roadmapFeatures  []founder.ProductFeature
	segmentsMenu     *components.Menu
	pricingMenu      *components.Menu
	acquisitionsMenu *components.Menu
	platformMenu     *components.Menu
	securityMenu     *components.Menu
	prCrisisMenu     *components.Menu
	economyMenu      *components.Menu
	successionMenu   *components.Menu

	// Advanced inputs
	pricingInput    textinput.Model
	securityInput   textinput.Model
	successionInput textinput.Model

	// New feature menus for parity
	strategicOpportunityMenu *components.Menu
	contentMarketingMenu     *components.Menu
	csPlaybooksMenu          *components.Menu
	competitiveIntelMenu     *components.Menu
	referralProgramMenu      *components.Menu
	techDebtMenu             *components.Menu
	endAffiliateMenu         *components.Menu

	// New inputs
	contentBudgetInput    textinput.Model
	csPlaybookInput       textinput.Model
	referralRewardInput   textinput.Model
	intelReportInput      textinput.Model
	techDebtRefactorInput textinput.Model

	// Acquisition offer state
	pendingAcquisition *founder.AcquisitionOffer
	acquisitionMenu    *components.Menu

	// Board management state
	advisorExpertiseMenu *components.Menu
	advisorConfirmMenu   *components.Menu
	removeAdvisorMenu    *components.Menu
	selectedExpertise    string
	pendingAdvisorName   string
	pendingAdvisorCost   float64
	pendingAdvisorSetup  int64

	// Equity pool input
	equityPoolInput textinput.Model

	// Engineer reallocation state
	reallocMenu     *components.Menu
	reallocFeatures []founder.ProductFeature

	// Confirm quit state
	confirmQuitMenu *components.Menu

	// Board sub-action tracking (for chairman selection vs fire)
	pendingBoardSubAction string

	// Executive offer negotiation
	execOffers   []founder.ExecOffer
	execOfferMenu *components.Menu
}

// NewFounderGameScreen creates a new founder game screen
func NewFounderGameScreen(width, height int, gameData *GameData) *FounderGameScreen {
	// Marketing input
	marketingInput := textinput.New()
	marketingInput.Placeholder = "Enter amount (e.g. 50000)"
	marketingInput.CharLimit = 15
	marketingInput.Width = 20

	// Affiliate input
	affiliateInput := textinput.New()
	affiliateInput.Placeholder = "Enter commission % (5-30)"
	affiliateInput.CharLimit = 5
	affiliateInput.Width = 10

	// Buyback input
	buybackInput := textinput.New()
	buybackInput.Placeholder = "Enter equity % to buy back"
	buybackInput.CharLimit = 10
	buybackInput.Width = 15

	// Pricing input
	pricingInput := textinput.New()
	pricingInput.Placeholder = "Enter discount %"
	pricingInput.CharLimit = 5
	pricingInput.Width = 10

	// Security input
	securityInput := textinput.New()
	securityInput.Placeholder = "Monthly budget"
	securityInput.CharLimit = 10
	securityInput.Width = 15

	// Succession input
	successionInput := textinput.New()
	successionInput.Placeholder = "Backup person name"
	successionInput.CharLimit = 30
	successionInput.Width = 25

	// Content budget input
	contentBudgetInput := textinput.New()
	contentBudgetInput.Placeholder = "Monthly budget (10000-50000)"
	contentBudgetInput.CharLimit = 10
	contentBudgetInput.Width = 15

	// CS Playbook input
	csPlaybookInput := textinput.New()
	csPlaybookInput.Placeholder = "Budget per month"
	csPlaybookInput.CharLimit = 10
	csPlaybookInput.Width = 15

	// Referral reward input
	referralRewardInput := textinput.New()
	referralRewardInput.Placeholder = "Reward amount (100-1000)"
	referralRewardInput.CharLimit = 10
	referralRewardInput.Width = 15

	// Intel report input
	intelReportInput := textinput.New()
	intelReportInput.Placeholder = "Competitor name"
	intelReportInput.CharLimit = 30
	intelReportInput.Width = 25

	// Tech debt refactor input
	techDebtRefactorInput := textinput.New()
	techDebtRefactorInput.Placeholder = "Budget (50000-200000)"
	techDebtRefactorInput.CharLimit = 10
	techDebtRefactorInput.Width = 15

	// Equity pool input
	equityPoolInput := textinput.New()
	equityPoolInput.Placeholder = "Enter % (1-10)"
	equityPoolInput.CharLimit = 5
	equityPoolInput.Width = 10

	s := &FounderGameScreen{
		width:                 width,
		height:                height,
		gameData:              gameData,
		view:                  FounderViewMain,
		marketingInput:        marketingInput,
		affiliateInput:        affiliateInput,
		buybackInput:          buybackInput,
		pricingInput:          pricingInput,
		securityInput:         securityInput,
		successionInput:       successionInput,
		contentBudgetInput:    contentBudgetInput,
		csPlaybookInput:       csPlaybookInput,
		referralRewardInput:   referralRewardInput,
		intelReportInput:      intelReportInput,
		techDebtRefactorInput: techDebtRefactorInput,
		equityPoolInput:       equityPoolInput,
	}

	s.rebuildActionsMenu()
	s.rebuildHiringMenu()
	s.rebuildFiringMenu()

	// Generate initial founding story
	fg := gameData.FounderState
	totalTeam := len(fg.Team.Engineers) + len(fg.Team.Sales) + len(fg.Team.CustomerSuccess) + len(fg.Team.Marketing) + len(fg.Team.Executives)
	runwayMonths := 0
	if fg.MonthlyTeamCost+fg.FounderSalary > 0 {
		totalBurn := fg.MonthlyTeamCost + fg.FounderSalary
		runwayMonths = int(fg.Cash / totalBurn)
	}

	story := []string{
		fmt.Sprintf("🚀 FOUNDING STORY: %s", fg.CompanyName),
		"",
		fmt.Sprintf("   %s and a small team of %d co-founders pooled their resources", fg.FounderName, totalTeam),
		fmt.Sprintf("   to bootstrap %s, a %s startup.", fg.CompanyName, fg.Category),
		fmt.Sprintf("   \"%s\"", fg.Description),
		"",
		fmt.Sprintf("   💰 Initial Capital: $%s (self-funded)", formatCompactMoney(fg.Cash)),
		fmt.Sprintf("   👥 Founding Team: %d members", totalTeam),
		fmt.Sprintf("   📊 Starting Customers: %d", fg.Customers),
		fmt.Sprintf("   💵 Monthly Burn: $%s/mo", formatCompactMoney(fg.MonthlyTeamCost+fg.FounderSalary)),
		fmt.Sprintf("   ⏱️  Estimated Runway: %d months", runwayMonths),
		"",
		"   The journey begins. Build your product, grow your team, and scale to unicorn status!",
		"   Press ENTER for actions or N to advance to next month.",
	}
	s.turnMessages = story

	return s
}

func (s *FounderGameScreen) rebuildActionsMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{
		// Team & Operations
		{ID: "header_team", Title: "── TEAM & OPERATIONS ──", Disabled: true, Icon: ""},
		{ID: "hiring", Title: "Hire Team Member", Description: "Engineers, Sales, CS, Marketing, C-Suite", Icon: "💼"},
		{ID: "firing", Title: "Let Go Team Member", Description: "Reduce headcount to cut costs", Icon: "⚠️"},
		{ID: "marketing", Title: "Spend on Marketing", Description: "Acquire customers with ad spend", Icon: "📣"},

		// Funding & Equity
		{ID: "header_funding", Title: "── FUNDING & EQUITY ──", Disabled: true, Icon: ""},
		{ID: "funding", Title: "Raise Funding Round", Description: "Seed, Series A, B, C", Icon: "💰"},
	}

	// Buyback only if profitable and have raised
	if len(fg.FundingRounds) > 0 && fg.MRR > fg.MonthlyTeamCost {
		items = append(items, components.MenuItem{
			ID: "buyback", Title: "Buy Back Equity", Description: "Repurchase shares from investors", Icon: "📈",
		})
	}

	items = append(items, components.MenuItem{
		ID: "board", Title: "Manage Board & Equity", Description: "Advisors, cap table, equity pool", Icon: "👔",
	})

	// Strategic options
	items = append(items, components.MenuItem{
		ID: "header_strategic", Title: "── STRATEGIC ──", Disabled: true, Icon: "",
	})
	items = append(items, components.MenuItem{
		ID: "partnership", Title: "Start Partnership", Description: "Distribution, tech, co-marketing", Icon: "🤝",
	})

	if fg.AffiliateProgram == nil {
		items = append(items, components.MenuItem{
			ID: "affiliate", Title: "Launch Affiliate Program", Description: "Let partners sell for commission", Icon: "💸",
		})
	} else {
		items = append(items, components.MenuItem{
			ID: "affiliate_view", Title: "View Affiliate Program", Description: fmt.Sprintf("%d affiliates, $%s/mo", fg.AffiliateProgram.Affiliates, formatCompactMoney(fg.AffiliateMRR)), Icon: "💸",
		})
		items = append(items, components.MenuItem{
			ID: "end_affiliate", Title: "End Affiliate Program", Description: "Shut down affiliate program", Icon: "🚫",
		})
	}

	// Referral Program
	if fg.ReferralProgram == nil && fg.Customers >= 10 {
		items = append(items, components.MenuItem{
			ID: "referral", Title: "Launch Referral Program", Description: "Customers refer new customers", Icon: "🎁",
		})
	} else if fg.ReferralProgram != nil {
		items = append(items, components.MenuItem{
			ID: "referral_view", Title: "View Referral Program", Description: fmt.Sprintf("%d referrals total", fg.ReferralProgram.TotalReferrals), Icon: "🎁",
		})
	}

	activeCompetitors := 0
	for _, c := range fg.Competitors {
		if c.Active {
			activeCompetitors++
		}
	}
	if activeCompetitors > 0 {
		items = append(items, components.MenuItem{
			ID: "competitors", Title: "Handle Competitors", Description: fmt.Sprintf("%d active competitors", activeCompetitors), Icon: "⚔️",
		})
	}

	items = append(items, components.MenuItem{
		ID: "expansion", Title: "Expand to New Market", Description: "EU, APAC, LATAM expansion", Icon: "🌍",
	})
	items = append(items, components.MenuItem{
		ID: "pivot", Title: "Execute Pivot", Description: "Change market or strategy", Icon: "🔄",
	})

	// Advanced strategic options (Phase 4)
	items = append(items, components.MenuItem{
		ID: "roadmap", Title: "Product Roadmap", Description: "Manage feature development", Icon: "🔨",
	})

	if fg.Customers >= 50 {
		items = append(items, components.MenuItem{
			ID: "segments", Title: "Customer Segments", Description: "Target customer verticals", Icon: "🎯",
		})
	}

	if len(fg.FundingRounds) > 0 || fg.MRR >= 100000 {
		items = append(items, components.MenuItem{
			ID: "pricing", Title: "Pricing Strategy", Description: "Adjust pricing model", Icon: "💲",
		})
	}

	if fg.CanAcquire() {
		items = append(items, components.MenuItem{
			ID: "acquisitions", Title: "Acquisitions", Description: "Acquire other companies", Icon: "🏢",
		})
	}

	if fg.CanLaunchPlatform() {
		items = append(items, components.MenuItem{
			ID: "platform", Title: "Platform Strategy", Description: "Build developer ecosystem", Icon: "🌐",
		})
	}

	if fg.SecurityPosture != nil && fg.CanHaveSecurityIncidents() {
		items = append(items, components.MenuItem{
			ID: "security", Title: "Security & Compliance", Description: "Manage security posture", Icon: "🔒",
		})
	}

	if fg.PRProgram != nil && fg.PRProgram.HasPRFirm {
		items = append(items, components.MenuItem{
			ID: "pr_crisis", Title: "PR Crisis Management", Description: "Handle PR issues", Icon: "📰",
		})
	}

	if fg.Turn >= 12 {
		items = append(items, components.MenuItem{
			ID: "economy", Title: "Economic Strategy", Description: "Respond to market conditions", Icon: "📉",
		})
	}

	if fg.CanHaveKeyPersonRisk() {
		items = append(items, components.MenuItem{
			ID: "succession", Title: "Succession Planning", Description: "Manage key person risk", Icon: "👤",
		})
	}

	// Content Marketing (unlock: marketing hire OR $200k MRR)
	hasMarketing := len(fg.Team.Marketing) > 0
	if hasMarketing || fg.MRR >= 200000 {
		if fg.ContentProgram == nil {
			items = append(items, components.MenuItem{
				ID: "content_marketing", Title: "Launch Content Marketing", Description: "SEO and content program", Icon: "📝",
			})
		} else {
			items = append(items, components.MenuItem{
				ID: "content_marketing", Title: "Manage Content Marketing", Description: fmt.Sprintf("SEO Score: %d, Traffic: %d", fg.ContentProgram.SEOScore, fg.ContentProgram.OrganicTraffic), Icon: "📝",
			})
		}
	}

	// CS Playbooks (unlock: CS hire OR 100+ customers)
	hasCS := len(fg.Team.CustomerSuccess) > 0
	if hasCS || fg.Customers >= 100 {
		items = append(items, components.MenuItem{
			ID: "cs_playbooks", Title: "CS Playbooks", Description: fmt.Sprintf("%d active playbooks", len(fg.CSPlaybooks)), Icon: "📋",
		})
	}

	// Competitive Intelligence (unlock: Series A OR 5+ competitors)
	hasSeriesA := false
	for _, round := range fg.FundingRounds {
		if round.RoundName == "Series A" {
			hasSeriesA = true
			break
		}
	}
	if hasSeriesA || len(fg.Competitors) >= 5 {
		if fg.CompetitiveIntel == nil {
			items = append(items, components.MenuItem{
				ID: "competitive_intel", Title: "Launch Competitive Intel", Description: "Hire analyst, gather intel", Icon: "🕵️",
			})
		} else {
			items = append(items, components.MenuItem{
				ID: "competitive_intel", Title: "Competitive Intelligence", Description: fmt.Sprintf("%d reports", len(fg.CompetitiveIntel.IntelReports)), Icon: "🕵️",
			})
		}
	}

	// Technical Debt (unlock: 5+ engineers OR $1M+ MRR)
	if len(fg.Team.Engineers) >= 5 || fg.MRR >= 1000000 {
		debtLevel := 0
		if fg.TechnicalDebt != nil {
			debtLevel = fg.TechnicalDebt.CurrentLevel
		}
		items = append(items, components.MenuItem{
			ID: "tech_debt", Title: "Technical Debt", Description: fmt.Sprintf("Current debt: %d/100", debtLevel), Icon: "🔧",
		})
	}

	// View data
	items = append(items, components.MenuItem{
		ID: "header_view", Title: "── VIEW DATA ──", Disabled: true, Icon: "",
	})
	items = append(items, components.MenuItem{
		ID: "team_roster", Title: "View Team Roster", Description: "See all employees and their impact", Icon: "👥",
	})
	items = append(items, components.MenuItem{
		ID: "customers", Title: "View Customers", Description: "Customer deals and health", Icon: "🏢",
	})
	items = append(items, components.MenuItem{
		ID: "financials", Title: "View Financials", Description: "Cash flow breakdown", Icon: "📊",
	})
	items = append(items, components.MenuItem{
		ID: "cap_table", Title: "View Cap Table", Description: "Equity ownership breakdown", Icon: "📋",
	})

	if fg.SalesPipeline != nil && len(fg.SalesPipeline.ActiveDeals) > 0 {
		items = append(items, components.MenuItem{
			ID: "pipeline", Title: "View Sales Pipeline", Description: "Track deal progress", Icon: "📈",
		})
	}

	// Solicit Customer Feedback (when customers > 0)
	if fg.Customers > 0 {
		items = append(items, components.MenuItem{
			ID: "feedback", Title: "Solicit Customer Feedback", Description: "Improve product, reduce churn", Icon: "💬",
		})
	}

	// Strategic Opportunity (when pending)
	if fg.PendingOpportunity != nil {
		items = append(items, components.MenuItem{
			ID:          "opportunity",
			Title:       "Strategic Opportunity!",
			Description: fmt.Sprintf("%s (expires in %d months)", fg.PendingOpportunity.Title, fg.PendingOpportunity.ExpiresIn),
			Icon:        "💡",
		})
	}

	// Exit options
	exits := fg.GetAvailableExits()
	hasAvailableExit := false
	for _, exit := range exits {
		if exit.CanExit && exit.Type != "continue" {
			hasAvailableExit = true
			break
		}
	}
	if hasAvailableExit {
		items = append(items, components.MenuItem{
			ID: "exit", Title: "Consider Exit Options", Description: "IPO, Acquisition, Secondary sale", Icon: "🚪",
		})
	}

	// Continue / Skip
	items = append(items, components.MenuItem{
		ID: "header_action", Title: "── ACTION ──", Disabled: true, Icon: "",
	})
	items = append(items, components.MenuItem{
		ID: "continue", Title: "Continue to Next Month", Description: "Process this month's events", Icon: "⏭️",
	})

	s.actionsMenu = components.NewMenu("DECISIONS", items)
	s.actionsMenu.SetSize(55, 20)
	s.actionsMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) rebuildHiringMenu() {
	items := []components.MenuItem{
		{ID: "header_ic", Title: "── INDIVIDUAL CONTRIBUTORS ($100k/yr) ──", Disabled: true},
		{ID: "engineer", Title: "Engineer", Description: "Builds product, reduces churn", Icon: "👨‍💻"},
		{ID: "sales", Title: "Sales Rep", Description: "Increases customer acquisition", Icon: "📞"},
		{ID: "cs", Title: "Customer Success", Description: "Reduces churn rate", Icon: "🤝"},
		{ID: "marketing", Title: "Marketing", Description: "Supports acquisition campaigns", Icon: "📢"},
		{ID: "header_exec", Title: "── C-LEVEL EXECUTIVES ($300k/yr, 3x impact) ──", Disabled: true},
		{ID: "cto", Title: "CTO", Description: "Like hiring 3 engineers", Icon: "⚡"},
		{ID: "cgo", Title: "CGO (Growth)", Description: "Like hiring 3 sales reps", Icon: "📈"},
		{ID: "coo", Title: "COO (Operations)", Description: "Like hiring 3 CS reps", Icon: "⚙️"},
		{ID: "cfo", Title: "CFO (Finance)", Description: "Reduces burn by 10%", Icon: "💹"},
		{ID: "cancel", Title: "Cancel", Description: "Go back", Icon: "←"},
	}
	s.hiringMenu = components.NewMenu("HIRE", items)
	s.hiringMenu.SetSize(50, 15)
	s.hiringMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) rebuildFiringMenu() {
	fg := s.gameData.FounderState

	execCount := make(map[founder.EmployeeRole]int)
	for _, exec := range fg.Team.Executives {
		execCount[exec.Role]++
	}

	items := []components.MenuItem{
		{ID: "header_ic", Title: "── INDIVIDUAL CONTRIBUTORS ──", Disabled: true},
		{ID: "engineer", Title: fmt.Sprintf("Engineer (current: %d)", len(fg.Team.Engineers)), Icon: "👨‍💻",
			Disabled: len(fg.Team.Engineers) == 0},
		{ID: "sales", Title: fmt.Sprintf("Sales Rep (current: %d)", len(fg.Team.Sales)), Icon: "📞",
			Disabled: len(fg.Team.Sales) == 0},
		{ID: "cs", Title: fmt.Sprintf("Customer Success (current: %d)", len(fg.Team.CustomerSuccess)), Icon: "🤝",
			Disabled: len(fg.Team.CustomerSuccess) == 0},
		{ID: "marketing", Title: fmt.Sprintf("Marketing (current: %d)", len(fg.Team.Marketing)), Icon: "📢",
			Disabled: len(fg.Team.Marketing) == 0},
		{ID: "header_exec", Title: "── EXECUTIVES ──", Disabled: true},
		{ID: "cto", Title: fmt.Sprintf("CTO (current: %d)", execCount[founder.RoleCTO]), Icon: "⚡",
			Disabled: execCount[founder.RoleCTO] == 0},
		{ID: "cgo", Title: fmt.Sprintf("CGO (current: %d)", execCount[founder.RoleCGO]), Icon: "📈",
			Disabled: execCount[founder.RoleCGO] == 0},
		{ID: "coo", Title: fmt.Sprintf("COO (current: %d)", execCount[founder.RoleCOO]), Icon: "⚙️",
			Disabled: execCount[founder.RoleCOO] == 0},
		{ID: "cfo", Title: fmt.Sprintf("CFO (current: %d)", execCount[founder.RoleCFO]), Icon: "💹",
			Disabled: execCount[founder.RoleCFO] == 0},
		{ID: "cancel", Title: "Cancel", Description: "Go back", Icon: "←"},
	}
	s.firingMenu = components.NewMenu("LET GO", items)
	s.firingMenu.SetSize(50, 15)
	s.firingMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) rebuildFundingMenu() {
	fg := s.gameData.FounderState

	hasSeed := false
	hasSeriesA := false
	hasSeriesB := false

	for _, round := range fg.FundingRounds {
		switch round.RoundName {
		case "Seed":
			hasSeed = true
		case "Series A":
			hasSeriesA = true
		case "Series B":
			hasSeriesB = true
		}
	}

	items := []components.MenuItem{}

	if !hasSeed {
		items = append(items, components.MenuItem{
			ID: "seed", Title: "Seed Round ($2-5M)", Description: "Early stage funding", Icon: "🌱",
		})
	}
	if hasSeed && !hasSeriesA {
		items = append(items, components.MenuItem{
			ID: "series_a", Title: "Series A ($10-20M)", Description: "Growth stage funding", Icon: "🚀",
		})
	}
	if hasSeriesA && !hasSeriesB {
		items = append(items, components.MenuItem{
			ID: "series_b", Title: "Series B ($30-50M)", Description: "Scale stage funding", Icon: "📈",
		})
	}

	if len(items) == 0 {
		items = append(items, components.MenuItem{
			ID: "none", Title: "No rounds available", Disabled: true, Icon: "❌",
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Cancel", Icon: "←",
	})

	s.fundingMenu = components.NewMenu("RAISE FUNDING", items)
	s.fundingMenu.SetSize(50, 12)
	s.fundingMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) rebuildExitMenu() {
	fg := s.gameData.FounderState
	exits := fg.GetAvailableExits()

	items := []components.MenuItem{}

	for _, exit := range exits {
		if exit.Type == "continue" {
			continue
		}
		icon := "🔒"
		if exit.CanExit {
			switch exit.Type {
			case "ipo":
				icon = "🏛️"
			case "acquisition":
				icon = "🤝"
			case "secondary":
				icon = "💼"
			}
		}
		items = append(items, components.MenuItem{
			ID:          exit.Type,
			Title:       fmt.Sprintf("%s: $%s valuation", strings.ToUpper(exit.Type), formatCompactMoney(exit.Valuation)),
			Description: exit.Description,
			Icon:        icon,
			Disabled:    !exit.CanExit,
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Keep Building", Icon: "←",
	})

	s.exitMenu = components.NewMenu("EXIT OPTIONS", items)
	s.exitMenu.SetSize(60, 15)
	s.exitMenu.SetHideHelp(true)
}

// Init initializes the founder game screen
func (s *FounderGameScreen) Init() tea.Cmd {
	return nil
}

// Update handles founder game input
func (s *FounderGameScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s.view {
		case FounderViewMain:
			switch {
			case key.Matches(msg, keys.Global.Enter):
				s.rebuildActionsMenu()
				s.view = FounderViewActions
				return s, nil
			case msg.String() == "n":
				// Quick shortcut: advance to next month
				fg := s.gameData.FounderState
				preDecisionMRR := fg.MRR
				msgs := fg.ProcessMonthWithBaseline(preDecisionMRR)
				s.turnMessages = msgs
				// Check for acquisition offers
				offer := fg.CheckForAcquisition()
				if offer != nil {
					s.pendingAcquisition = offer
					s.rebuildAcquisitionMenu()
					s.view = FounderViewAcquisitionOffer
					return s, nil
				}
				s.view = FounderViewMain
				return s, nil
			case key.Matches(msg, keys.Global.Back), msg.String() == "q":
				// Show confirmation instead of immediately quitting
				s.rebuildConfirmQuitMenu()
				s.view = FounderViewConfirmQuit
				return s, nil
			}

		case FounderViewActions:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewMain
				return s, nil
			}

		case FounderViewHiring:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewExecOffer:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewHiring
				return s, nil
			}

		case FounderViewHiringMarket:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewHiring
				return s, nil
			}
			// Handle number keys for market selection
			keyStr := msg.String()
			if len(keyStr) == 1 && keyStr[0] >= '1' && keyStr[0] <= '9' {
				num := int(keyStr[0] - '0')
				return s.handleMarketSelection(num)
			}

		case FounderViewFiring:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewMarketing:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.view = FounderViewActions
				s.inputMessage = ""
				return s, nil
			case key.Matches(msg, keys.Global.Enter):
				return s.handleMarketingConfirm()
			}

		case FounderViewFunding:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewFundingTerms:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewFunding
				return s, nil
			}
			// Handle number keys for term selection
			keyStr := msg.String()
			if len(keyStr) == 1 && keyStr[0] >= '1' && keyStr[0] <= '4' {
				num := int(keyStr[0] - '0')
				return s.handleTermSelection(num)
			}

		case FounderViewPartnership:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewAffiliate:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.view = FounderViewActions
				s.inputMessage = ""
				return s, nil
			case key.Matches(msg, keys.Global.Enter):
				return s.handleAffiliateConfirm()
			}

		case FounderViewCompetitors:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}
			// Handle number selection for competitor
			keyStr := msg.String()
			if len(keyStr) == 1 && keyStr[0] >= '1' && keyStr[0] <= '9' {
				num := int(keyStr[0] - '0')
				return s.handleCompetitorSelect(num)
			}

		case FounderViewExpansion:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewPivot:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewBoard:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewBoardAction:
			if key.Matches(msg, keys.Global.Back) {
				s.rebuildBoardMenu()
				s.view = FounderViewBoard
				return s, nil
			}
			// Handle equity pool input enter
			if msg.Type == tea.KeyEnter {
				val := strings.TrimSpace(s.equityPoolInput.Value())
				if val != "" {
					pct, err := strconv.ParseFloat(val, 64)
					if err == nil && pct >= 1 && pct <= 10 {
						fg := s.gameData.FounderState
						fg.ExpandEquityPool(pct)
						avail := fg.EquityPool - fg.EquityAllocated
						if avail < 0 {
							avail = 0
						}
						s.turnMessages = []string{
							fmt.Sprintf("✓ Expanded equity pool by %.1f%%", pct),
							fmt.Sprintf("   Total pool: %.1f%% (%.1f%% available)", fg.EquityPool, avail),
							fmt.Sprintf("   Your equity: %.1f%%", 100.0-fg.EquityPool-fg.EquityGivenAway),
						}
						s.rebuildBoardMenu()
						s.view = FounderViewBoard
						return s, nil
					}
				}
			}

		case FounderViewBuyback:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewBuybackConfirm:
			switch {
			case key.Matches(msg, keys.Global.Back):
				s.view = FounderViewBuyback
				s.inputMessage = ""
				return s, nil
			case key.Matches(msg, keys.Global.Enter):
				return s.handleBuybackConfirm()
			}

		case FounderViewTeamRoster, FounderViewCustomers, FounderViewFinancials, FounderViewCapTable:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewExit:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewRoadmap:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewRoadmapStart:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewRoadmap
				return s, nil
			}
			// Handle number keys for feature selection
			keyStr := msg.String()
			if len(keyStr) == 1 && keyStr[0] >= '1' && keyStr[0] <= '9' {
				num := int(keyStr[0] - '0')
				return s.handleFeatureStart(num)
			}

		case FounderViewSegments:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewPricing:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewAcquisitions:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewPlatform:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewSecurity:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewPRCrisis:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewEconomy:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewSuccession:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewSalesPipeline:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		// New feature views
		case FounderViewStrategicOpportunity:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewContentMarketing:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewCSPlaybooks:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewCompetitiveIntel:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewReferralProgram:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewTechDebt:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewAcquisitionOffer:
			if key.Matches(msg, keys.Global.Back) {
				s.pendingAcquisition = nil
				s.view = FounderViewMain
				return s, nil
			}

		case FounderViewAdvisorExpertise:
			if key.Matches(msg, keys.Global.Back) {
				s.rebuildBoardMenu()
				s.view = FounderViewBoard
				return s, nil
			}

		case FounderViewAdvisorConfirm:
			if key.Matches(msg, keys.Global.Back) {
				s.rebuildAdvisorExpertiseMenu()
				s.view = FounderViewAdvisorExpertise
				return s, nil
			}

		case FounderViewRemoveAdvisor:
			if key.Matches(msg, keys.Global.Back) {
				s.rebuildBoardMenu()
				s.view = FounderViewBoard
				return s, nil
			}

		case FounderViewBoardTable:
			if key.Matches(msg, keys.Global.Back) || key.Matches(msg, keys.Global.Enter) {
				s.rebuildBoardMenu()
				s.view = FounderViewBoard
				return s, nil
			}

		case FounderViewConfirmQuit:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewMain
				return s, nil
			}

		case FounderViewEngineerRealloc:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewRoadmap
				return s, nil
			}

		case FounderViewEndAffiliate:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewActions
				return s, nil
			}

		case FounderViewConfirmExit:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewExit
				return s, nil
			}
			if msg.String() == "y" {
				return s.handleExitConfirm()
			}
			if msg.String() == "n" {
				s.view = FounderViewActions
				return s, nil
			}
		}

	case components.MenuSelectedMsg:
		switch s.view {
		case FounderViewActions:
			return s.handleAction(msg.ID)
		case FounderViewHiring:
			return s.handleHiringSelection(msg.ID)
		case FounderViewExecOffer:
			return s.handleExecOfferSelection(msg.ID)
		case FounderViewFiring:
			return s.handleFiringSelection(msg.ID)
		case FounderViewFunding:
			return s.handleFundingSelection(msg.ID)
		case FounderViewPartnership:
			return s.handlePartnershipSelection(msg.ID)
		case FounderViewExit:
			return s.handleExitSelection(msg.ID)
		case FounderViewCompetitors:
			return s.handleCompetitorStrategy(msg.ID)
		case FounderViewExpansion:
			return s.handleExpansionSelection(msg.ID)
		case FounderViewPivot:
			return s.handlePivotSelection(msg.ID)
		case FounderViewBoard:
			return s.handleBoardSelection(msg.ID)
		case FounderViewBuyback:
			return s.handleBuybackSelection(msg.ID)
		case FounderViewRoadmap:
			return s.handleRoadmapSelection(msg.ID)
		case FounderViewSegments:
			return s.handleSegmentsSelection(msg.ID)
		case FounderViewPricing:
			return s.handlePricingSelection(msg.ID)
		case FounderViewAcquisitions:
			return s.handleAcquisitionsSelection(msg.ID)
		case FounderViewPlatform:
			return s.handlePlatformSelection(msg.ID)
		case FounderViewSecurity:
			return s.handleSecuritySelection(msg.ID)
		case FounderViewPRCrisis:
			return s.handlePRCrisisSelection(msg.ID)
		case FounderViewEconomy:
			return s.handleEconomySelection(msg.ID)
		case FounderViewSuccession:
			return s.handleSuccessionSelection(msg.ID)
		// New feature handlers
		case FounderViewStrategicOpportunity:
			return s.handleStrategicOpportunitySelection(msg.ID)
		case FounderViewContentMarketing:
			return s.handleContentMarketingSelection(msg.ID)
		case FounderViewCSPlaybooks:
			return s.handleCSPlaybooksSelection(msg.ID)
		case FounderViewCompetitiveIntel:
			return s.handleCompetitiveIntelSelection(msg.ID)
		case FounderViewReferralProgram:
			return s.handleReferralProgramSelection(msg.ID)
		case FounderViewTechDebt:
			return s.handleTechDebtSelection(msg.ID)
		case FounderViewEndAffiliate:
			return s.handleEndAffiliateSelection(msg.ID)
		case FounderViewAcquisitionOffer:
			return s.handleAcquisitionSelection(msg.ID)
		case FounderViewAdvisorExpertise:
			return s.handleAdvisorExpertiseSelection(msg.ID)
		case FounderViewAdvisorConfirm:
			return s.handleAdvisorConfirmSelection(msg.ID)
		case FounderViewRemoveAdvisor:
			return s.handleRemoveAdvisorSelection(msg.ID)
		case FounderViewConfirmQuit:
			return s.handleConfirmQuitSelection(msg.ID)
		case FounderViewEngineerRealloc:
			return s.handleEngineerReallocSelection(msg.ID)
		}
	}

	// Check for game over
	if fg != nil && fg.IsGameOver() {
		return s, SwitchTo(ScreenFounderResults)
	}

	// Update current component
	var cmd tea.Cmd
	switch s.view {
	case FounderViewActions:
		s.actionsMenu, cmd = s.actionsMenu.Update(msg)
	case FounderViewHiring:
		s.hiringMenu, cmd = s.hiringMenu.Update(msg)
	case FounderViewExecOffer:
		if s.execOfferMenu != nil {
			s.execOfferMenu, cmd = s.execOfferMenu.Update(msg)
		}
	case FounderViewFiring:
		s.firingMenu, cmd = s.firingMenu.Update(msg)
	case FounderViewFunding:
		s.fundingMenu, cmd = s.fundingMenu.Update(msg)
	case FounderViewMarketing:
		s.marketingInput, cmd = s.marketingInput.Update(msg)
	case FounderViewAffiliate:
		s.affiliateInput, cmd = s.affiliateInput.Update(msg)
	case FounderViewPartnership:
		s.partnershipMenu, cmd = s.partnershipMenu.Update(msg)
	case FounderViewExit:
		s.exitMenu, cmd = s.exitMenu.Update(msg)
	case FounderViewCompetitors:
		if s.competitorAction != nil {
			s.competitorAction, cmd = s.competitorAction.Update(msg)
		}
	case FounderViewExpansion:
		if s.expansionMenu != nil {
			s.expansionMenu, cmd = s.expansionMenu.Update(msg)
		}
	case FounderViewPivot:
		if s.pivotMenu != nil {
			s.pivotMenu, cmd = s.pivotMenu.Update(msg)
		}
	case FounderViewBoard:
		if s.boardMenu != nil {
			s.boardMenu, cmd = s.boardMenu.Update(msg)
		}
	case FounderViewBoardAction:
		s.equityPoolInput, cmd = s.equityPoolInput.Update(msg)
	case FounderViewBuyback:
		if s.buybackMenu != nil {
			s.buybackMenu, cmd = s.buybackMenu.Update(msg)
		}
	case FounderViewBuybackConfirm:
		s.buybackInput, cmd = s.buybackInput.Update(msg)
	case FounderViewRoadmap:
		if s.roadmapMenu != nil {
			s.roadmapMenu, cmd = s.roadmapMenu.Update(msg)
		}
	case FounderViewSegments:
		if s.segmentsMenu != nil {
			s.segmentsMenu, cmd = s.segmentsMenu.Update(msg)
		}
	case FounderViewPricing:
		if s.pricingMenu != nil {
			s.pricingMenu, cmd = s.pricingMenu.Update(msg)
		}
	case FounderViewAcquisitions:
		if s.acquisitionsMenu != nil {
			s.acquisitionsMenu, cmd = s.acquisitionsMenu.Update(msg)
		}
	case FounderViewPlatform:
		if s.platformMenu != nil {
			s.platformMenu, cmd = s.platformMenu.Update(msg)
		}
	case FounderViewSecurity:
		if s.securityMenu != nil {
			s.securityMenu, cmd = s.securityMenu.Update(msg)
		}
	case FounderViewPRCrisis:
		if s.prCrisisMenu != nil {
			s.prCrisisMenu, cmd = s.prCrisisMenu.Update(msg)
		}
	case FounderViewEconomy:
		if s.economyMenu != nil {
			s.economyMenu, cmd = s.economyMenu.Update(msg)
		}
	case FounderViewSuccession:
		if s.successionMenu != nil {
			s.successionMenu, cmd = s.successionMenu.Update(msg)
		}
	// New feature menu updates
	case FounderViewStrategicOpportunity:
		if s.strategicOpportunityMenu != nil {
			s.strategicOpportunityMenu, cmd = s.strategicOpportunityMenu.Update(msg)
		}
	case FounderViewContentMarketing:
		if s.contentMarketingMenu != nil {
			s.contentMarketingMenu, cmd = s.contentMarketingMenu.Update(msg)
		}
	case FounderViewCSPlaybooks:
		if s.csPlaybooksMenu != nil {
			s.csPlaybooksMenu, cmd = s.csPlaybooksMenu.Update(msg)
		}
	case FounderViewCompetitiveIntel:
		if s.competitiveIntelMenu != nil {
			s.competitiveIntelMenu, cmd = s.competitiveIntelMenu.Update(msg)
		}
	case FounderViewReferralProgram:
		if s.referralProgramMenu != nil {
			s.referralProgramMenu, cmd = s.referralProgramMenu.Update(msg)
		}
	case FounderViewTechDebt:
		if s.techDebtMenu != nil {
			s.techDebtMenu, cmd = s.techDebtMenu.Update(msg)
		}
	case FounderViewAcquisitionOffer:
		if s.acquisitionMenu != nil {
			s.acquisitionMenu, cmd = s.acquisitionMenu.Update(msg)
		}
	case FounderViewAdvisorExpertise:
		if s.advisorExpertiseMenu != nil {
			s.advisorExpertiseMenu, cmd = s.advisorExpertiseMenu.Update(msg)
		}
	case FounderViewAdvisorConfirm:
		if s.advisorConfirmMenu != nil {
			s.advisorConfirmMenu, cmd = s.advisorConfirmMenu.Update(msg)
		}
	case FounderViewRemoveAdvisor:
		if s.removeAdvisorMenu != nil {
			s.removeAdvisorMenu, cmd = s.removeAdvisorMenu.Update(msg)
		}
	case FounderViewEndAffiliate:
		if s.endAffiliateMenu != nil {
			s.endAffiliateMenu, cmd = s.endAffiliateMenu.Update(msg)
		}
	case FounderViewConfirmQuit:
		if s.confirmQuitMenu != nil {
			s.confirmQuitMenu, cmd = s.confirmQuitMenu.Update(msg)
		}
	case FounderViewEngineerRealloc:
		if s.reallocMenu != nil {
			s.reallocMenu, cmd = s.reallocMenu.Update(msg)
		}
	}

	return s, cmd
}

func (s *FounderGameScreen) handleAction(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "continue":
		// Capture MRR before processing
		preDecisionMRR := fg.MRR
		msgs := fg.ProcessMonthWithBaseline(preDecisionMRR)
		s.turnMessages = msgs

		// Check for acquisition offers
		offer := fg.CheckForAcquisition()
		if offer != nil {
			s.pendingAcquisition = offer
			s.rebuildAcquisitionMenu()
			s.view = FounderViewAcquisitionOffer
			return s, nil
		}

		s.view = FounderViewMain
		return s, nil

	case "hiring":
		s.rebuildHiringMenu()
		s.view = FounderViewHiring
		return s, nil

	case "firing":
		s.rebuildFiringMenu()
		s.view = FounderViewFiring
		return s, nil

	case "marketing":
		s.marketingInput.SetValue("")
		s.marketingInput.Focus()
		s.inputMessage = ""
		s.view = FounderViewMarketing
		return s, textinput.Blink

	case "funding":
		s.rebuildFundingMenu()
		s.view = FounderViewFunding
		return s, nil

	case "buyback":
		s.rebuildBuybackMenu()
		s.view = FounderViewBuyback
		return s, nil

	case "board":
		s.rebuildBoardMenu()
		s.view = FounderViewBoard
		return s, nil

	case "partnership":
		s.rebuildPartnershipMenu()
		s.view = FounderViewPartnership
		return s, nil

	case "affiliate":
		s.affiliateInput.SetValue("")
		s.affiliateInput.Focus()
		s.inputMessage = ""
		s.view = FounderViewAffiliate
		return s, textinput.Blink

	case "affiliate_view":
		monthsActive := fg.Turn - fg.AffiliateProgram.LaunchedMonth
		avgPerAffiliate := int64(0)
		if fg.AffiliateProgram.Affiliates > 0 {
			avgPerAffiliate = fg.AffiliateMRR / int64(fg.AffiliateProgram.Affiliates)
		}
		s.turnMessages = []string{
			"┌─────────── AFFILIATE PROGRAM STATS ───────────┐",
			fmt.Sprintf("│ Launched: Month %d (%d months active)", fg.AffiliateProgram.LaunchedMonth, monthsActive),
			fmt.Sprintf("│ Commission Rate: %.1f%%", fg.AffiliateProgram.Commission*100),
			fmt.Sprintf("│ Active Affiliates: %d", fg.AffiliateProgram.Affiliates),
			fmt.Sprintf("│ Setup Cost: $%s", formatCompactMoney(fg.AffiliateProgram.SetupCost)),
			fmt.Sprintf("│ Monthly Platform Fee: $%s", formatCompactMoney(fg.AffiliateProgram.MonthlyPlatformFee)),
			"│",
			fmt.Sprintf("│ Customers Acquired: %d", fg.AffiliateProgram.CustomersAcquired),
			fmt.Sprintf("│ Affiliate MRR: $%s", formatCompactMoney(fg.AffiliateMRR)),
			fmt.Sprintf("│ Monthly Revenue: $%s", formatCompactMoney(fg.AffiliateProgram.MonthlyRevenue)),
			fmt.Sprintf("│ Avg Revenue/Affiliate: $%s", formatCompactMoney(avgPerAffiliate)),
			"└────────────────────────────────────────────────┘",
		}
		s.view = FounderViewMain
		return s, nil

	case "end_affiliate":
		s.rebuildEndAffiliateMenu()
		s.view = FounderViewEndAffiliate
		return s, nil

	case "referral":
		s.rebuildReferralProgramMenu()
		s.view = FounderViewReferralProgram
		return s, nil

	case "referral_view":
		if fg.ReferralProgram != nil {
			s.turnMessages = []string{
				fmt.Sprintf("Referral Program Active since month %d", fg.ReferralProgram.LaunchedMonth),
				fmt.Sprintf("Reward per Referral: $%s (%s)", formatCompactMoney(fg.ReferralProgram.RewardPerReferral), fg.ReferralProgram.RewardType),
				fmt.Sprintf("Total Referrals: %d", fg.ReferralProgram.TotalReferrals),
				fmt.Sprintf("Customers Acquired: %d", fg.ReferralProgram.CustomersAcquired),
				fmt.Sprintf("Monthly Budget: $%s", formatCompactMoney(fg.ReferralProgram.MonthlyBudget)),
			}
		}
		s.view = FounderViewMain
		return s, nil

	case "competitors":
		s.rebuildCompetitorMenu()
		s.view = FounderViewCompetitors
		return s, nil

	case "expansion":
		s.rebuildExpansionMenu()
		s.view = FounderViewExpansion
		return s, nil

	case "pivot":
		s.rebuildPivotMenu()
		s.view = FounderViewPivot
		return s, nil

	case "team_roster":
		s.view = FounderViewTeamRoster
		return s, nil

	case "customers":
		s.view = FounderViewCustomers
		return s, nil

	case "financials":
		s.view = FounderViewFinancials
		return s, nil

	case "cap_table":
		s.view = FounderViewCapTable
		return s, nil

	case "exit":
		s.rebuildExitMenu()
		s.view = FounderViewExit
		return s, nil

	case "roadmap":
		s.rebuildRoadmapMenu()
		s.view = FounderViewRoadmap
		return s, nil

	case "segments":
		s.rebuildSegmentsMenu()
		s.view = FounderViewSegments
		return s, nil

	case "pricing":
		s.rebuildPricingMenu()
		s.view = FounderViewPricing
		return s, nil

	case "acquisitions":
		s.rebuildAcquisitionsMenu()
		s.view = FounderViewAcquisitions
		return s, nil

	case "platform":
		s.rebuildPlatformMenu()
		s.view = FounderViewPlatform
		return s, nil

	case "security":
		s.rebuildSecurityMenu()
		s.view = FounderViewSecurity
		return s, nil

	case "pr_crisis":
		s.rebuildPRCrisisMenu()
		s.view = FounderViewPRCrisis
		return s, nil

	case "economy":
		s.rebuildEconomyMenu()
		s.view = FounderViewEconomy
		return s, nil

	case "succession":
		s.rebuildSuccessionMenu()
		s.view = FounderViewSuccession
		return s, nil

	case "pipeline":
		s.view = FounderViewSalesPipeline
		return s, nil

	case "feedback":
		err := fg.SolicitCustomerFeedback()
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				"✓ Solicited customer feedback",
				"  Product maturity improved",
				"  Customer churn reduced",
			}
		}
		s.view = FounderViewMain
		return s, nil

	case "opportunity":
		s.rebuildStrategicOpportunityMenu()
		s.view = FounderViewStrategicOpportunity
		return s, nil

	case "content_marketing":
		s.rebuildContentMarketingMenu()
		s.view = FounderViewContentMarketing
		return s, nil

	case "cs_playbooks":
		s.rebuildCSPlaybooksMenu()
		s.view = FounderViewCSPlaybooks
		return s, nil

	case "competitive_intel":
		s.rebuildCompetitiveIntelMenu()
		s.view = FounderViewCompetitiveIntel
		return s, nil

	case "tech_debt":
		s.rebuildTechDebtMenu()
		s.view = FounderViewTechDebt
		return s, nil
	}

	return s, nil
}

func (s *FounderGameScreen) handleHiringSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	var role founder.EmployeeRole
	var isExec bool

	switch id {
	case "engineer":
		role = founder.RoleEngineer
	case "sales":
		role = founder.RoleSales
	case "cs":
		role = founder.RoleCustomerSuccess
	case "marketing":
		role = founder.RoleMarketing
	case "cto":
		role = founder.RoleCTO
		isExec = true
	case "cgo":
		role = founder.RoleCGO
		isExec = true
	case "coo":
		role = founder.RoleCOO
		isExec = true
	case "cfo":
		role = founder.RoleCFO
		isExec = true
	default:
		return s, nil
	}

	s.selectedRole = role
	s.selectedIsExec = isExec

	// Executive hires go through offer negotiation
	if isExec {
		offers, err := fg.GenerateExecOffers(role)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ %v", err)}
			s.view = FounderViewMain
			return s, nil
		}
		s.execOffers = offers
		s.rebuildExecOfferMenu()
		s.view = FounderViewExecOffer
		return s, nil
	}

	// Check if market selection is needed
	if (role == founder.RoleSales || role == founder.RoleMarketing || role == founder.RoleCustomerSuccess) && len(fg.GlobalMarkets) > 0 {
		s.marketOptions = []string{"USA"}
		for _, m := range fg.GlobalMarkets {
			s.marketOptions = append(s.marketOptions, m.Region)
		}
		s.marketOptions = append(s.marketOptions, "All")
		s.view = FounderViewHiringMarket
		return s, nil
	}

	// Direct hire for regular employees
	err := fg.HireEmployee(role)
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		s.turnMessages = []string{
			fmt.Sprintf("✓ Hired a new %s!", role),
			fmt.Sprintf("   Cost: $100k/year"),
			fmt.Sprintf("   New runway: %d months", fg.CashRunwayMonths),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// Executive offer negotiation
func (s *FounderGameScreen) rebuildExecOfferMenu() {
	items := []components.MenuItem{}
	for i, offer := range s.execOffers {
		equityStr := fmt.Sprintf("%.1f%% equity", offer.Equity)
		salaryStr := fmt.Sprintf("$%dk/yr", offer.AnnualCost/1000)
		impactStr := fmt.Sprintf("%.1fx impact", offer.Impact)
		items = append(items, components.MenuItem{
			ID:          fmt.Sprintf("offer_%d", i),
			Title:       fmt.Sprintf("%s — %s + %s", offer.Label, equityStr, salaryStr),
			Description: fmt.Sprintf("%s (%s)", offer.Description, impactStr),
			Icon:        []string{"📋", "🚀", "💵"}[i%3],
		})
	}
	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Pass on this candidate", Icon: "←",
	})
	s.execOfferMenu = components.NewMenu("CHOOSE COMPENSATION PACKAGE", items)
	s.execOfferMenu.SetSize(65, 12)
	s.execOfferMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleExecOfferSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewHiring
		return s, nil
	}

	if strings.HasPrefix(id, "offer_") {
		idxStr := strings.TrimPrefix(id, "offer_")
		idx := 0
		fmt.Sscanf(idxStr, "%d", &idx)

		if idx >= 0 && idx < len(s.execOffers) {
			offer := s.execOffers[idx]
			err := fg.HireExecWithOffer(offer)
			if err != nil {
				s.turnMessages = []string{fmt.Sprintf("❌ %v", err)}
			} else {
				s.turnMessages = []string{
					fmt.Sprintf("✓ Hired %s as %s!", offer.Name, strings.ToUpper(string(offer.Role))),
					fmt.Sprintf("   Package: %s", offer.Label),
					fmt.Sprintf("   Equity: %.1f%% (4yr vest, 1yr cliff)", offer.Equity),
					fmt.Sprintf("   Salary: $%dk/year ($%dk/mo)", offer.AnnualCost/1000, offer.MonthlyCost/1000),
					fmt.Sprintf("   Impact: %.1fx", offer.Impact),
					fmt.Sprintf("   New runway: %d months", fg.CashRunwayMonths),
				}
			}
			s.view = FounderViewMain
			return s, nil
		}
	}

	return s, nil
}

func (s *FounderGameScreen) renderExecOffer() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(65).
		Align(lipgloss.Center)

	if len(s.execOffers) > 0 {
		offer := s.execOffers[0]
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(
			headerStyle.Render(fmt.Sprintf("🤝 %s WANTS TO JOIN AS %s", strings.ToUpper(offer.Name), strings.ToUpper(string(offer.Role))))))
	} else {
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🤝 EXECUTIVE OFFER NEGOTIATION")))
	}
	b.WriteString("\n\n")

	fg := s.gameData.FounderState
	availPool := fg.EquityPool - fg.EquityAllocated
	if availPool < 0 {
		availPool = 0
	}
	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	b.WriteString(infoStyle.Render(fmt.Sprintf("Equity Pool: %.1f%% available | Cash: $%s | Runway: %d months",
		availPool, formatCompactMoney(fg.Cash), fg.CashRunwayMonths)))
	b.WriteString("\n\n")

	// Offer comparison table
	if len(s.execOffers) > 0 {
		tableBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Magenta).
			Padding(1, 2).
			Width(65)

		var table strings.Builder
		colHead := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
		table.WriteString(colHead.Render(fmt.Sprintf("  %-18s %8s  %10s  %7s", "PACKAGE", "EQUITY", "SALARY", "IMPACT")))
		table.WriteString("\n")
		dimStyle := lipgloss.NewStyle().Foreground(styles.Gray)
		table.WriteString(dimStyle.Render("  " + strings.Repeat("─", 50)))
		table.WriteString("\n")

		colors := []lipgloss.Color{styles.White, styles.Green, styles.Yellow}
		for i, offer := range s.execOffers {
			style := lipgloss.NewStyle().Foreground(colors[i%len(colors)])
			table.WriteString(style.Render(fmt.Sprintf("  %-18s %7.1f%%  $%dk/yr  %5.1fx",
				offer.Label, offer.Equity, offer.AnnualCost/1000, offer.Impact)))
			table.WriteString("\n")
		}

		table.WriteString("\n")
		noteStyle := lipgloss.NewStyle().Foreground(styles.Gray).Italic(true)
		table.WriteString(noteStyle.Render("  More equity → lower salary, higher motivation"))
		table.WriteString("\n")
		table.WriteString(noteStyle.Render("  More cash → less dilution, but higher burn rate"))

		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(tableBox.Render(table.String())))
		b.WriteString("\n\n")
	}

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	if s.execOfferMenu != nil {
		b.WriteString(menuContainer.Render(s.execOfferMenu.View()))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) handleMarketSelection(num int) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if num < 1 || num > len(s.marketOptions) {
		return s, nil
	}

	market := s.marketOptions[num-1]

	var err error
	if market != "USA" {
		err = fg.HireEmployeeWithMarket(s.selectedRole, market)
	} else {
		err = fg.HireEmployee(s.selectedRole)
	}

	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		s.turnMessages = []string{
			fmt.Sprintf("✓ Hired a new %s!", s.selectedRole),
			fmt.Sprintf("   Assigned to: %s", market),
			fmt.Sprintf("   New runway: %d months", fg.CashRunwayMonths),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

func (s *FounderGameScreen) handleFiringSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	var role founder.EmployeeRole
	switch id {
	case "engineer":
		role = founder.RoleEngineer
	case "sales":
		role = founder.RoleSales
	case "cs":
		role = founder.RoleCustomerSuccess
	case "marketing":
		role = founder.RoleMarketing
	case "cto":
		role = founder.RoleCTO
	case "cgo":
		role = founder.RoleCGO
	case "coo":
		role = founder.RoleCOO
	case "cfo":
		role = founder.RoleCFO
	default:
		return s, nil
	}

	err := fg.FireEmployee(role)
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		runway := fmt.Sprintf("%d months", fg.CashRunwayMonths)
		if fg.CashRunwayMonths < 0 {
			runway = "∞ (profitable!)"
		}
		s.turnMessages = []string{
			fmt.Sprintf("✓ Let go one %s", role),
			fmt.Sprintf("   New runway: %s", runway),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

func (s *FounderGameScreen) handleMarketingConfirm() (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	amountStr := strings.TrimSpace(strings.ReplaceAll(s.marketingInput.Value(), ",", ""))
	if amountStr == "" || amountStr == "0" {
		s.view = FounderViewActions
		return s, nil
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount < 0 {
		s.inputMessage = "Invalid amount"
		return s, nil
	}

	if amount > fg.Cash {
		s.inputMessage = "Not enough cash!"
		return s, nil
	}

	newCustomers := fg.SpendOnMarketing(amount)
	s.turnMessages = []string{
		"✓ Marketing campaign launched!",
		fmt.Sprintf("   Acquired %d new customers!", newCustomers),
		fmt.Sprintf("   New MRR: $%s", formatCompactMoney(fg.MRR)),
	}

	s.view = FounderViewMain
	return s, nil
}

func (s *FounderGameScreen) handleFundingSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" || id == "none" {
		s.view = FounderViewActions
		return s, nil
	}

	var roundName string
	switch id {
	case "seed":
		roundName = "Seed"
	case "series_a":
		roundName = "Series A"
	case "series_b":
		roundName = "Series B"
	default:
		return s, nil
	}

	s.selectedRoundName = roundName
	s.fundingTerms = fg.GenerateTermSheetOptions(roundName)

	if len(s.fundingTerms) == 0 {
		s.turnMessages = []string{"❌ Unable to generate term sheets!"}
		s.view = FounderViewMain
		return s, nil
	}

	s.view = FounderViewFundingTerms
	return s, nil
}

func (s *FounderGameScreen) handleTermSelection(num int) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if num < 1 || num > len(s.fundingTerms) {
		return s, nil
	}

	selectedSheet := s.fundingTerms[num-1]
	success := fg.RaiseFundingWithTerms(s.selectedRoundName, selectedSheet)

	if !success {
		s.turnMessages = []string{"❌ Failed to raise funding!"}
	} else {
		runway := fmt.Sprintf("%d months", fg.CashRunwayMonths)
		if fg.CashRunwayMonths < 0 {
			runway = "∞ (profitable!)"
		}
		s.turnMessages = []string{
			fmt.Sprintf("✓ Successfully raised %s!", s.selectedRoundName),
			fmt.Sprintf("   Amount: $%s", formatCompactMoney(selectedSheet.Amount)),
			fmt.Sprintf("   Valuation: $%s", formatCompactMoney(selectedSheet.PostValuation)),
			fmt.Sprintf("   Equity Given: %.1f%%", selectedSheet.Equity),
			fmt.Sprintf("   Your equity: %.1f%%", 100.0-fg.EquityGivenAway-fg.EquityPool),
			fmt.Sprintf("   New runway: %s", runway),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

func (s *FounderGameScreen) rebuildPartnershipMenu() {
	items := []components.MenuItem{
		{ID: "distribution", Title: "Distribution Partnership ($50-150k)", Description: "10-30% MRR boost", Icon: "🚚"},
		{ID: "technology", Title: "Technology Partnership ($30-100k)", Description: "Product integration, reduces churn", Icon: "🔧"},
		{ID: "co-marketing", Title: "Co-Marketing Partnership ($25-75k)", Description: "15-40% MRR boost", Icon: "📣"},
		{ID: "data", Title: "Data Partnership ($40-100k)", Description: "Analytics/insights, reduces churn", Icon: "📊"},
		{ID: "cancel", Title: "Cancel", Icon: "←"},
	}
	s.partnershipMenu = components.NewMenu("PARTNERSHIPS", items)
	s.partnershipMenu.SetSize(55, 12)
	s.partnershipMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handlePartnershipSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	partnership, err := fg.StartPartnership(id)
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		s.turnMessages = []string{
			fmt.Sprintf("✓ Partnership with %s started!", partnership.Partner),
			fmt.Sprintf("   Type: %s", partnership.Type),
			fmt.Sprintf("   Cost: $%s", formatCompactMoney(partnership.Cost)),
			fmt.Sprintf("   Duration: %d months", partnership.Duration),
			fmt.Sprintf("   Expected MRR Boost: $%s/mo", formatCompactMoney(partnership.MRRBoost)),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

func (s *FounderGameScreen) handleAffiliateConfirm() (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	commStr := strings.TrimSpace(s.affiliateInput.Value())
	commission, err := strconv.ParseFloat(commStr, 64)
	if err != nil || commission < 5 || commission > 30 {
		s.inputMessage = "Invalid commission rate (5-30%)"
		return s, nil
	}

	err = fg.LaunchAffiliateProgram(commission)
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		s.turnMessages = []string{
			"✓ Affiliate program launched!",
			fmt.Sprintf("   Commission: %.1f%%", commission),
			fmt.Sprintf("   Starting Affiliates: %d", fg.AffiliateProgram.Affiliates),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// Competitor management
func (s *FounderGameScreen) rebuildCompetitorMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{}
	for i, comp := range fg.Competitors {
		if !comp.Active {
			continue
		}
		threatIcon := "🟡"
		if comp.Threat == "high" || comp.Threat == "critical" {
			threatIcon = "🔴"
		} else if comp.Threat == "low" {
			threatIcon = "🟢"
		}
		items = append(items, components.MenuItem{
			ID:          fmt.Sprintf("comp_%d", i),
			Title:       comp.Name,
			Description: fmt.Sprintf("%s | Market: %.1f%% | Strategy: %s", comp.Threat, comp.MarketShare*100, comp.Strategy),
			Icon:        threatIcon,
		})
	}

	if len(items) == 0 {
		items = append(items, components.MenuItem{
			ID:       "none",
			Title:    "No active competitors",
			Disabled: true,
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Back", Icon: "←",
	})

	s.competitorMenu = components.NewMenu("COMPETITORS", items)
	s.competitorMenu.SetSize(55, 15)
	s.competitorMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleCompetitorSelect(num int) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	activeCount := 0
	for i, comp := range fg.Competitors {
		if !comp.Active {
			continue
		}
		activeCount++
		if activeCount == num {
			s.selectedCompetitorIdx = i
			s.rebuildCompetitorActionMenu()
			return s, nil
		}
	}
	return s, nil
}

func (s *FounderGameScreen) rebuildCompetitorActionMenu() {
	items := []components.MenuItem{
		{ID: "ignore", Title: "Ignore", Description: "No cost, they may take market share", Icon: "😐"},
		{ID: "compete", Title: "Compete Aggressively", Description: "$50-150k, reduce their threat", Icon: "⚔️"},
		{ID: "partner", Title: "Partner With Them", Description: "$100-250k, merge customer bases", Icon: "🤝"},
		{ID: "cancel", Title: "Cancel", Icon: "←"},
	}
	s.competitorAction = components.NewMenu("STRATEGY", items)
	s.competitorAction.SetSize(50, 10)
	s.competitorAction.SetHideHelp(true)
}

func (s *FounderGameScreen) handleCompetitorStrategy(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" || id == "none" {
		// Go back to competitor list, not main menu
		s.competitorAction = nil
		s.selectedCompetitorIdx = -1
		s.view = FounderViewCompetitors
		return s, nil
	}

	// If it's a competitor selection
	if strings.HasPrefix(id, "comp_") {
		idxStr := strings.TrimPrefix(id, "comp_")
		idx, _ := strconv.Atoi(idxStr)
		// Validate competitor is still active
		if idx >= 0 && idx < len(fg.Competitors) && fg.Competitors[idx].Active {
			s.selectedCompetitorIdx = idx
			s.rebuildCompetitorActionMenu()
		} else {
			s.turnMessages = []string{"⚠️ That competitor is no longer active"}
			s.view = FounderViewCompetitors
		}
		return s, nil
	}

	// Validate competitor is still active before taking action
	if s.selectedCompetitorIdx < 0 || s.selectedCompetitorIdx >= len(fg.Competitors) || !fg.Competitors[s.selectedCompetitorIdx].Active {
		s.turnMessages = []string{"⚠️ That competitor is no longer active"}
		s.competitorAction = nil
		s.selectedCompetitorIdx = -1
		s.view = FounderViewCompetitors
		return s, nil
	}

	// Handle strategy selection
	message, err := fg.HandleCompetitor(s.selectedCompetitorIdx, id)
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		s.turnMessages = []string{fmt.Sprintf("✓ %s", message)}
	}

	// Return to competitor list after action, not main view
	s.competitorAction = nil
	s.selectedCompetitorIdx = -1
	s.view = FounderViewCompetitors
	return s, nil
}

// Global expansion
func (s *FounderGameScreen) rebuildExpansionMenu() {
	fg := s.gameData.FounderState

	activeMarkets := make(map[string]bool)
	for _, m := range fg.GlobalMarkets {
		activeMarkets[m.Region] = true
	}

	type marketOption struct {
		region, cost, monthlyCost, competition string
	}
	allMarkets := []marketOption{
		{"Europe", "$200k", "$30k/mo", "high"},
		{"Asia", "$250k", "$40k/mo", "very high"},
		{"LATAM", "$150k", "$20k/mo", "medium"},
		{"Middle East", "$180k", "$25k/mo", "low"},
		{"Africa", "$120k", "$15k/mo", "low"},
		{"Australia", "$100k", "$18k/mo", "medium"},
	}

	items := []components.MenuItem{}
	for _, m := range allMarkets {
		if activeMarkets[m.region] {
			continue
		}
		items = append(items, components.MenuItem{
			ID:          strings.ToLower(strings.ReplaceAll(m.region, " ", "_")),
			Title:       m.region,
			Description: fmt.Sprintf("%s setup, %s, %s competition", m.cost, m.monthlyCost, m.competition),
			Icon:        "🌍",
		})
	}

	if len(items) == 0 {
		items = append(items, components.MenuItem{
			ID:       "none",
			Title:    "All markets expanded!",
			Disabled: true,
			Icon:     "✓",
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Cancel", Icon: "←",
	})

	s.expansionMenu = components.NewMenu("EXPAND TO", items)
	s.expansionMenu.SetSize(55, 12)
	s.expansionMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleExpansionSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" || id == "none" {
		s.view = FounderViewActions
		return s, nil
	}

	regionMap := map[string]string{
		"europe":      "Europe",
		"asia":        "Asia",
		"latam":       "LATAM",
		"middle_east": "Middle East",
		"africa":      "Africa",
		"australia":   "Australia",
	}

	region := regionMap[id]
	if region == "" {
		s.view = FounderViewActions
		return s, nil
	}

	market, err := fg.ExpandToMarket(region)
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		s.turnMessages = []string{
			fmt.Sprintf("✓ Launched in %s!", market.Region),
			fmt.Sprintf("   Setup Cost: $%s", formatCompactMoney(market.SetupCost)),
			fmt.Sprintf("   Initial Customers: %d", market.CustomerCount),
			fmt.Sprintf("   Monthly Cost: $%s", formatCompactMoney(market.MonthlyCost)),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// Pivot
func (s *FounderGameScreen) rebuildPivotMenu() {
	items := []components.MenuItem{
		{ID: "header", Title: "⚠️ WARNING: Pivots are risky!", Disabled: true},
		{ID: "enterprise_b2b", Title: "Enterprise B2B", Icon: "🏢"},
		{ID: "smb_b2b", Title: "SMB B2B", Icon: "🏪"},
		{ID: "b2c", Title: "B2C", Icon: "👤"},
		{ID: "marketplace", Title: "Marketplace", Icon: "🛒"},
		{ID: "platform", Title: "Platform", Icon: "🔧"},
		{ID: "vertical_saas", Title: "Vertical SaaS", Icon: "📊"},
		{ID: "horizontal_saas", Title: "Horizontal SaaS", Icon: "📈"},
		{ID: "deep_tech", Title: "Deep Tech", Icon: "🔬"},
		{ID: "consumer_apps", Title: "Consumer Apps", Icon: "📱"},
		{ID: "cancel", Title: "Cancel", Icon: "←"},
	}

	s.pivotMenu = components.NewMenu("PIVOT TO", items)
	s.pivotMenu.SetSize(50, 15)
	s.pivotMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handlePivotSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" || id == "header" {
		s.view = FounderViewActions
		return s, nil
	}

	strategyMap := map[string]string{
		"enterprise_b2b":  "Enterprise B2B",
		"smb_b2b":         "SMB B2B",
		"b2c":             "B2C",
		"marketplace":     "Marketplace",
		"platform":        "Platform",
		"vertical_saas":   "Vertical SaaS",
		"horizontal_saas": "Horizontal SaaS",
		"deep_tech":       "Deep Tech",
		"consumer_apps":   "Consumer Apps",
	}

	toStrategy := strategyMap[id]
	if toStrategy == "" {
		s.view = FounderViewActions
		return s, nil
	}

	pivot, err := fg.ExecutePivot(toStrategy, "Strategic repositioning")
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		if pivot.Success {
			s.turnMessages = []string{
				"🎉 Pivot SUCCESSFUL!",
				fmt.Sprintf("   New Strategy: %s", pivot.ToStrategy),
				fmt.Sprintf("   Cost: $%s", formatCompactMoney(pivot.Cost)),
				fmt.Sprintf("   Customers Lost: %d", pivot.CustomersLost),
			}
		} else {
			s.turnMessages = []string{
				"😞 Pivot FAILED",
				"   The market didn't respond well",
				fmt.Sprintf("   Cost: $%s", formatCompactMoney(pivot.Cost)),
				fmt.Sprintf("   Customers Lost: %d", pivot.CustomersLost),
			}
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// Board management
func (s *FounderGameScreen) rebuildBoardMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{
		{ID: "view", Title: "View Board & Advisors", Description: "See current board members", Icon: "👁️"},
		{ID: "add_seat", Title: "Add Board Seat", Description: "~2% from equity pool", Icon: "➕"},
		{ID: "expand_pool", Title: "Expand Equity Pool", Description: "Dilute 1-10%", Icon: "📊"},
		{ID: "add_advisor", Title: "Add Advisor", Description: "0.25-1% equity for guidance", Icon: "🧠"},
	}

	// Check for advisors who can be promoted to chairman
	hasActiveAdvisor := false
	hasChairman := false
	for _, m := range fg.BoardMembers {
		if m.IsActive {
			hasActiveAdvisor = true
			if m.IsChairman {
				hasChairman = true
			}
		}
	}

	if hasActiveAdvisor && !hasChairman {
		items = append(items, components.MenuItem{
			ID: "set_chairman", Title: "Set Chairman", Description: "Promote advisor to chairman (2x impact, +0.25% equity)", Icon: "👑",
		})
	}

	if hasChairman {
		items = append(items, components.MenuItem{
			ID: "remove_chairman", Title: "Remove Chairman", Description: "Demote chairman (board pressure increase)", Icon: "⬇️",
		})
	}

	if hasActiveAdvisor {
		items = append(items, components.MenuItem{
			ID: "remove_advisor", Title: "Remove Advisor", Description: "With equity buyback option", Icon: "❌",
		})
	}

	// Fire investor board member (requires 51%+ ownership)
	founderEquity := 100.0 - fg.EquityGivenAway - fg.EquityPool
	hasInvestorBoardMember := false
	for _, m := range fg.BoardMembers {
		if m.IsActive && m.Type == "investor" {
			hasInvestorBoardMember = true
			break
		}
	}
	if hasInvestorBoardMember && founderEquity >= 51.0 {
		items = append(items, components.MenuItem{
			ID: "fire_board_member", Title: "Fire Board Member", Description: "Remove investor director (requires 51%+)", Icon: "🔥",
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Cancel", Icon: "←",
	})

	s.boardMenu = components.NewMenu("BOARD & EQUITY", items)
	s.boardMenu.SetSize(55, 15)
	s.boardMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleBoardSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "cancel":
		s.view = FounderViewActions
		return s, nil

	case "view":
		s.view = FounderViewBoardTable
		return s, nil

	case "add_seat":
		fg.AddBoardSeat("Strategic advisor")
		s.turnMessages = []string{
			"✓ Added board seat",
			fmt.Sprintf("   New board seats: %d", fg.BoardSeats),
			fmt.Sprintf("   Remaining equity pool: %.1f%%", fg.EquityPool),
		}
		s.rebuildBoardMenu()
		s.view = FounderViewBoard
		return s, nil

	case "expand_pool":
		s.equityPoolInput.SetValue("")
		s.equityPoolInput.Focus()
		avail := fg.EquityPool - fg.EquityAllocated
		if avail < 0 {
			avail = 0
		}
		s.inputMessage = fmt.Sprintf("Pool: %.1f%% (%.1f%% available) | Your equity: %.1f%%",
			fg.EquityPool, avail, 100.0-fg.EquityPool-fg.EquityGivenAway)
		s.view = FounderViewBoardAction
		return s, textinput.Blink

	case "add_advisor":
		// Show expertise selection
		s.rebuildAdvisorExpertiseMenu()
		s.view = FounderViewAdvisorExpertise
		return s, nil

	case "remove_advisor":
		s.rebuildRemoveAdvisorMenu()
		s.view = FounderViewRemoveAdvisor
		return s, nil

	case "remove_chairman":
		chairman := fg.GetChairman()
		if chairman != nil {
			fg.RemoveChairman()
			s.turnMessages = []string{
				fmt.Sprintf("✓ Removed %s as Chairman", chairman.Name),
				"   ⚠️  Board pressure increased",
				"   ⚠️  Possible negative PR impact",
			}
		}
		s.rebuildBoardMenu()
		s.view = FounderViewBoard
		return s, nil

	case "set_chairman":
		// Build menu of advisors to pick from
		items := []components.MenuItem{}
		for _, m := range fg.BoardMembers {
			if m.IsActive && !m.IsChairman && m.Type == "advisor" {
				items = append(items, components.MenuItem{
					ID:          "chairman_" + m.Name,
					Title:       m.Name,
					Description: fmt.Sprintf("Expertise: %s | Score: %.0f%%", m.Expertise, m.ContributionScore*100),
					Icon:        "👑",
				})
			}
		}
		if len(items) == 1 {
			// Only one advisor, promote directly
			name := fg.BoardMembers[0].Name
			for _, m := range fg.BoardMembers {
				if m.IsActive && !m.IsChairman && m.Type == "advisor" {
					name = m.Name
					break
				}
			}
			err := fg.SetChairman(name)
			if err != nil {
				s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
			} else {
				s.turnMessages = []string{
					fmt.Sprintf("✓ %s promoted to Chairman of the Board", name),
					"   Chairman has 2x impact and can mitigate crises",
					"   Additional 0.25% equity granted",
				}
			}
			s.rebuildBoardMenu()
			s.view = FounderViewBoard
			return s, nil
		}
		if len(items) == 0 {
			s.turnMessages = []string{"⚠️ No eligible advisors to set as chairman. Hire an advisor first."}
			s.rebuildBoardMenu()
			s.view = FounderViewBoard
			return s, nil
		}
		// Multiple advisors - show menu
		s.pendingBoardSubAction = "set_chairman"
		items = append(items, components.MenuItem{ID: "cancel", Title: "Cancel", Icon: "←"})
		s.boardMenu = components.NewMenu("SELECT CHAIRMAN", items)
		s.boardMenu.SetSize(55, 15)
		s.boardMenu.SetHideHelp(true)
		// Stay in Board view with new menu
		return s, nil

	case "fire_board_member":
		// Build menu of investor board members
		items := []components.MenuItem{}
		for _, m := range fg.BoardMembers {
			if m.IsActive && m.Type == "investor" {
				items = append(items, components.MenuItem{
					ID:          "fire_" + m.Name,
					Title:       m.Name,
					Description: fmt.Sprintf("Equity: %.2f%% | ⚠️ Serious consequences", m.EquityCost),
					Icon:        "🔥",
				})
			}
		}
		if len(items) == 1 {
			name := strings.TrimPrefix(items[0].ID, "fire_")
			err := fg.FireBoardMember(name)
			if err != nil {
				s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
			} else {
				s.turnMessages = []string{
					fmt.Sprintf("✓ Fired board member: %s", name),
					"   ⚠️  Board pressure increased significantly",
					"   ⚠️  Board sentiment worsened",
				}
			}
			s.rebuildBoardMenu()
			s.view = FounderViewBoard
			return s, nil
		}
		if len(items) == 0 {
			s.turnMessages = []string{"⚠️ No investor board members to fire"}
			s.rebuildBoardMenu()
			s.view = FounderViewBoard
			return s, nil
		}
		s.pendingBoardSubAction = "fire_board_member"
		items = append(items, components.MenuItem{ID: "cancel", Title: "Cancel", Icon: "←"})
		s.boardMenu = components.NewMenu("FIRE BOARD MEMBER", items)
		s.boardMenu.SetSize(55, 15)
		s.boardMenu.SetHideHelp(true)
		// Stay in Board view with new menu
		return s, nil
	}

	// Handle chairman/fire sub-menu selections
	if strings.HasPrefix(id, "chairman_") {
		name := strings.TrimPrefix(id, "chairman_")
		err := fg.SetChairman(name)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				fmt.Sprintf("✓ %s promoted to Chairman of the Board", name),
				"   Chairman has 2x impact and can mitigate crises",
				"   Additional 0.25% equity granted",
			}
		}
		s.pendingBoardSubAction = ""
		s.rebuildBoardMenu()
		s.view = FounderViewBoard
		return s, nil
	}

	if strings.HasPrefix(id, "fire_") {
		name := strings.TrimPrefix(id, "fire_")
		err := fg.FireBoardMember(name)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				fmt.Sprintf("✓ Fired board member: %s", name),
				"   ⚠️  Board pressure increased significantly",
				"   ⚠️  Board sentiment worsened",
			}
		}
		s.pendingBoardSubAction = ""
		s.rebuildBoardMenu()
		s.view = FounderViewBoard
		return s, nil
	}

	return s, nil
}

// Buyback
func (s *FounderGameScreen) rebuildBuybackMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{}

	for _, round := range fg.FundingRounds {
		if round.EquityGiven > 0 {
			items = append(items, components.MenuItem{
				ID:          round.RoundName,
				Title:       round.RoundName,
				Description: fmt.Sprintf("%.1f%% equity available", round.EquityGiven),
				Icon:        "💰",
			})
		}
	}

	if len(items) == 0 {
		items = append(items, components.MenuItem{
			ID:       "none",
			Title:    "No rounds to buy back from",
			Disabled: true,
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Cancel", Icon: "←",
	})

	s.buybackMenu = components.NewMenu("BUY BACK FROM", items)
	s.buybackMenu.SetSize(50, 10)
	s.buybackMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleBuybackSelection(id string) (ScreenModel, tea.Cmd) {
	if id == "cancel" || id == "none" {
		s.view = FounderViewActions
		return s, nil
	}

	s.selectedBuybackRound = id
	s.buybackInput.SetValue("")
	s.buybackInput.Focus()
	s.inputMessage = ""
	s.view = FounderViewBuybackConfirm
	return s, textinput.Blink
}

func (s *FounderGameScreen) handleBuybackConfirm() (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	equityStr := strings.TrimSpace(s.buybackInput.Value())
	equity, err := strconv.ParseFloat(equityStr, 64)
	if err != nil || equity <= 0 {
		s.inputMessage = "Invalid equity percentage"
		return s, nil
	}

	buyback, err := fg.BuybackEquity(s.selectedBuybackRound, equity)
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		s.turnMessages = []string{
			fmt.Sprintf("✓ Bought back %.1f%% equity!", buyback.EquityBought),
			fmt.Sprintf("   Paid: $%s", formatCompactMoney(buyback.PricePaid)),
			fmt.Sprintf("   Your new ownership: %.1f%%", 100.0-fg.EquityGivenAway-fg.EquityPool),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

func (s *FounderGameScreen) handleExitSelection(id string) (ScreenModel, tea.Cmd) {
	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	s.selectedExitType = id
	s.view = FounderViewConfirmExit
	return s, nil
}

func (s *FounderGameScreen) handleExitConfirm() (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	fg.ExecuteExit(s.selectedExitType)

	// At exit, unallocated equity pool cancels — only allocated equity counts
	founderEquity := 100.0 - fg.EquityGivenAway - fg.EquityAllocated
	payout := int64(float64(fg.ExitValuation) * founderEquity / 100.0)

	s.turnMessages = []string{
		fmt.Sprintf("🎉 Congratulations! %s exit completed!", strings.ToUpper(s.selectedExitType)),
		fmt.Sprintf("   Exit Valuation: $%s", formatCompactMoney(fg.ExitValuation)),
		fmt.Sprintf("   Your Equity: %.1f%%", founderEquity),
		fmt.Sprintf("   Your Payout: $%s", formatCompactMoney(payout)),
	}

	s.view = FounderViewMain
	return s, nil
}

// ============================================================================
// PHASE 4: ADVANCED FEATURES
// ============================================================================

// Product Roadmap
func (s *FounderGameScreen) rebuildRoadmapMenu() {
	fg := s.gameData.FounderState

	// Initialize roadmap if needed
	if fg.ProductRoadmap == nil {
		fg.InitializeProductRoadmap()
	}

	items := []components.MenuItem{
		{ID: "view", Title: "View Roadmap Status", Description: "See in-progress and completed features", Icon: "👁️"},
		{ID: "start", Title: "Start New Feature", Description: "Assign engineers to build features", Icon: "🚀"},
	}

	inProgress := fg.GetInProgressFeatures()
	if len(inProgress) > 0 {
		items = append(items, components.MenuItem{
			ID: "reallocate", Title: "Reallocate Engineers", Description: "Adjust team assignments", Icon: "👥",
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Back", Icon: "←",
	})

	s.roadmapMenu = components.NewMenu("PRODUCT ROADMAP", items)
	s.roadmapMenu.SetSize(55, 12)
	s.roadmapMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleRoadmapSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "cancel":
		s.view = FounderViewActions
		return s, nil

	case "view":
		var msgs []string
		inProgress := fg.GetInProgressFeatures()
		completed := fg.GetCompletedFeatures()

		if len(inProgress) > 0 {
			msgs = append(msgs, "🔨 IN PROGRESS:")
			for _, f := range inProgress {
				msgs = append(msgs, fmt.Sprintf("  • %s (%d%% done, %d engineers)", f.Name, f.DevelopmentProgress, f.AllocatedEngineers))
			}
		}

		if len(completed) > 0 {
			msgs = append(msgs, "✅ COMPLETED:")
			for _, f := range completed {
				msgs = append(msgs, fmt.Sprintf("  • %s (Category: %s)", f.Name, f.Category))
			}
		}

		if len(msgs) == 0 {
			msgs = append(msgs, "No features in roadmap yet. Start building!")
		}

		s.turnMessages = msgs
		s.view = FounderViewMain
		return s, nil

	case "start":
		s.roadmapFeatures = fg.GetAvailableFeaturesToStart()
		s.view = FounderViewRoadmapStart
		return s, nil

	case "reallocate":
		s.rebuildEngineerReallocMenu()
		if s.reallocMenu == nil {
			s.turnMessages = []string{"⚠️ No features currently in progress to reallocate engineers"}
			s.view = FounderViewMain
			return s, nil
		}
		s.view = FounderViewEngineerRealloc
		return s, nil
	}

	return s, nil
}

func (s *FounderGameScreen) handleFeatureStart(num int) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if num < 1 || num > len(s.roadmapFeatures) {
		return s, nil
	}

	feature := s.roadmapFeatures[num-1]
	engineers := 1 // Start with 1 engineer
	if len(fg.Team.Engineers) > 1 {
		engineers = 2
	}

	err := fg.StartFeature(feature.Name, engineers)
	if err != nil {
		s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
	} else {
		s.turnMessages = []string{
			fmt.Sprintf("✓ Started feature: %s", feature.Name),
			fmt.Sprintf("   Engineers assigned: %d", engineers),
			fmt.Sprintf("   Estimated time: %d months", feature.EngineerMonths/engineers),
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// Customer Segments
func (s *FounderGameScreen) rebuildSegmentsMenu() {
	fg := s.gameData.FounderState

	if fg.CustomerSegments == nil {
		fg.InitializeSegments()
	}

	items := []components.MenuItem{
		{ID: "view", Title: "View Current Segments", Description: "See customer distribution", Icon: "👁️"},
	}

	// Add segment options with current focus indicator
	segments := []string{"SMB", "Mid-Market", "Enterprise", "Strategic"}
	actionLabel := "Select"
	if fg.SelectedICP != "" {
		actionLabel = "Switch to"
	}
	for _, seg := range segments {
		icon := "🎯"
		desc := ""
		if seg == fg.SelectedICP {
			icon = "✅"
			desc = "Currently focused"
		} else {
			desc = fmt.Sprintf("%s %s segment", actionLabel, seg)
		}
		items = append(items, components.MenuItem{
			ID:          "focus_" + strings.ToLower(strings.ReplaceAll(seg, "-", "_")),
			Title:       fmt.Sprintf("Focus on %s", seg),
			Description: desc,
			Icon:        icon,
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Back", Icon: "←",
	})

	s.segmentsMenu = components.NewMenu("CUSTOMER SEGMENTS", items)
	s.segmentsMenu.SetSize(55, 15)
	s.segmentsMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleSegmentsSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	if id == "view" {
		// Ensure segments are initialized and volumes are up to date
		if fg.CustomerSegments == nil {
			fg.InitializeSegments()
		}
		fg.UpdateSegmentVolumes()

		var msgs []string
		if len(fg.CustomerSegments) > 0 {
			icp := fg.SelectedICP
			if icp == "" {
				icp = "Not Set"
			}
			msgs = append(msgs, fmt.Sprintf("Current ICP: %s", icp))
			msgs = append(msgs, "")
			for _, seg := range fg.CustomerSegments {
				msgs = append(msgs, fmt.Sprintf("%s: %d customers, $%s/mo avg, %.0f%% churn",
					seg.Name, seg.Volume, formatCompactMoney(seg.AvgDealSize), seg.ChurnRate*100))
			}
		} else {
			msgs = append(msgs, "No segment data available")
		}
		s.turnMessages = msgs
		s.view = FounderViewMain
		return s, nil
	}

	if strings.HasPrefix(id, "focus_") {
		segmentMap := map[string]string{
			"focus_smb":        "SMB",
			"focus_mid_market": "Mid-Market",
			"focus_enterprise": "Enterprise",
			"focus_strategic":  "Strategic",
		}
		segment := segmentMap[id]
		if segment != "" {
			var err error
			if fg.SelectedICP == "" {
				// First time selecting an ICP
				err = fg.SelectICP(segment)
			} else if fg.SelectedICP == segment {
				s.turnMessages = []string{fmt.Sprintf("Already focused on %s", segment)}
				s.view = FounderViewMain
				return s, nil
			} else {
				// Changing from one ICP to another
				err = fg.ChangeICP(segment)
			}
			if err != nil {
				s.turnMessages = []string{fmt.Sprintf("❌ %v", err)}
			} else {
				s.turnMessages = []string{
					fmt.Sprintf("✓ Now focusing on %s segment!", segment),
					"   Benefits: -20%% CAC, +15%% close rate, +10%% deal size",
				}
			}
		}
		s.view = FounderViewMain
		return s, nil
	}

	return s, nil
}

// Pricing Strategy
func (s *FounderGameScreen) rebuildPricingMenu() {
	fg := s.gameData.FounderState

	if fg.PricingStrategy == nil {
		fg.InitializePricingStrategy()
	}

	items := []components.MenuItem{
		{ID: "view", Title: "View Current Pricing", Description: "See pricing model details", Icon: "👁️"},
		{ID: "freemium", Title: "Freemium Model", Description: "Free tier + paid upgrades", Icon: "🆓"},
		{ID: "trial", Title: "Free Trial Model", Description: "Trial period then paid", Icon: "⏳"},
		{ID: "annual_upfront", Title: "Annual Upfront", Description: "Annual billing with discount", Icon: "📆"},
		{ID: "usage_based", Title: "Usage-Based Pricing", Description: "Pay-as-you-go model", Icon: "📊"},
		{ID: "tiered", Title: "Tiered Pricing", Description: "Multiple tiers for segments", Icon: "📶"},
		{ID: "cancel", Title: "Back", Icon: "←"},
	}

	s.pricingMenu = components.NewMenu("PRICING STRATEGY", items)
	s.pricingMenu.SetSize(55, 14)
	s.pricingMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handlePricingSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "cancel":
		s.view = FounderViewActions
		return s, nil

	case "view":
		var msgs []string
		if fg.PricingStrategy != nil {
			msgs = append(msgs, fmt.Sprintf("Current Model: %s", fg.PricingStrategy.Model))
			msgs = append(msgs, fmt.Sprintf("Annual Billing: %v", fg.PricingStrategy.IsAnnual))
			msgs = append(msgs, fmt.Sprintf("Discount: %.0f%%", fg.PricingStrategy.Discount*100))
		} else {
			msgs = append(msgs, "No pricing strategy configured")
		}
		s.turnMessages = msgs
		s.view = FounderViewMain
		return s, nil

	case "freemium":
		err := fg.ChangePricingModel("freemium", false, 0)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Switched to freemium model", "Free tier attracts users, paid upgrades drive revenue"}
		}
		s.view = FounderViewMain
		return s, nil

	case "trial":
		err := fg.ChangePricingModel("trial", false, 0)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Switched to free trial model", "Users try before buying — higher conversion expected"}
		}
		s.view = FounderViewMain
		return s, nil

	case "annual_upfront":
		err := fg.ChangePricingModel("annual_upfront", true, 0.15)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Switched to annual upfront billing", "15% discount for commitment, lower churn"}
		}
		s.view = FounderViewMain
		return s, nil

	case "usage_based":
		err := fg.ChangePricingModel("usage_based", false, 0)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Switched to usage-based pricing", "Revenue scales with customer usage"}
		}
		s.view = FounderViewMain
		return s, nil

	case "tiered":
		err := fg.ChangePricingModel("tiered", false, 0)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Switched to tiered pricing", "Multiple tiers serve different customer segments"}
		}
		s.view = FounderViewMain
		return s, nil
	}

	return s, nil
}

// Acquisitions
func (s *FounderGameScreen) rebuildAcquisitionsMenu() {
	fg := s.gameData.FounderState

	// Generate targets if needed
	if len(fg.AcquisitionTargets) == 0 {
		fg.GenerateAcquisitionTargets()
	}

	items := []components.MenuItem{
		{ID: "view", Title: "View Acquisition History", Description: "Past and in-progress acquisitions", Icon: "👁️"},
	}

	for i, target := range fg.AcquisitionTargets {
		items = append(items, components.MenuItem{
			ID:          fmt.Sprintf("acquire_%d", i),
			Title:       fmt.Sprintf("Acquire %s", target.Name),
			Description: fmt.Sprintf("$%s, +$%s MRR, %s", formatCompactMoney(target.AcquisitionCost), formatCompactMoney(target.MRR), target.Category),
			Icon:        "🏢",
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Back", Icon: "←",
	})

	s.acquisitionsMenu = components.NewMenu("ACQUISITIONS", items)
	s.acquisitionsMenu.SetSize(60, 15)
	s.acquisitionsMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleAcquisitionsSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	if id == "view" {
		completed, inProgress, totalMRR, totalCost := fg.GetAcquisitionSummary()
		s.turnMessages = []string{
			fmt.Sprintf("Completed: %d acquisitions", completed),
			fmt.Sprintf("In Progress: %d integrations", inProgress),
			fmt.Sprintf("Total MRR Gained: $%s", formatCompactMoney(totalMRR)),
			fmt.Sprintf("Total Spent: $%s", formatCompactMoney(totalCost)),
		}
		s.view = FounderViewMain
		return s, nil
	}

	if strings.HasPrefix(id, "acquire_") {
		idxStr := strings.TrimPrefix(id, "acquire_")
		idx, _ := strconv.Atoi(idxStr)

		acq, err := fg.AcquireCompany(idx)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				fmt.Sprintf("✓ Acquired %s!", acq.TargetName),
				fmt.Sprintf("   Paid: $%s", formatCompactMoney(acq.Cost)),
				fmt.Sprintf("   MRR Added: $%s", formatCompactMoney(acq.MRRGained)),
				fmt.Sprintf("   Integration: %d months", acq.IntegrationMonths),
			}
		}
		s.view = FounderViewMain
		return s, nil
	}

	return s, nil
}

// Platform Strategy
func (s *FounderGameScreen) rebuildPlatformMenu() {
	fg := s.gameData.FounderState

	// Initialize platform if needed
	if fg.PlatformMetrics == nil {
		founder.InitializePlatform(fg)
	}

	items := []components.MenuItem{}

	if !fg.PlatformMetrics.IsPlatform {
		// Platform not yet launched - show launch options
		if fg.CanLaunchPlatform() {
			items = append(items, components.MenuItem{
				ID: "launch_marketplace", Title: "Launch Marketplace", Description: "Third-party app ecosystem", Icon: "🛒",
			})
			items = append(items, components.MenuItem{
				ID: "launch_social", Title: "Launch Social Platform", Description: "Community-driven network effects", Icon: "👥",
			})
			items = append(items, components.MenuItem{
				ID: "launch_data", Title: "Launch Data Platform", Description: "Data exchange and insights", Icon: "📊",
			})
			items = append(items, components.MenuItem{
				ID: "launch_infrastructure", Title: "Launch Infrastructure Platform", Description: "Developer tools and APIs", Icon: "🔌",
			})
		} else {
			items = append(items, components.MenuItem{
				ID: "locked", Title: "Platform Locked", Description: "Requires $1M+ ARR or 500+ customers", Icon: "🔒", Disabled: true,
			})
		}
	} else {
		// Platform launched - show management options
		items = append(items, components.MenuItem{
			ID: "view", Title: "View Platform Metrics", Description: "Developer stats and API usage", Icon: "👁️",
		})
		items = append(items, components.MenuItem{
			ID: "invest", Title: "Invest in Developer Program", Description: "Increase platform adoption", Icon: "💰",
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Back", Icon: "←",
	})

	s.platformMenu = components.NewMenu("PLATFORM STRATEGY", items)
	s.platformMenu.SetSize(55, 12)
	s.platformMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handlePlatformSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "cancel":
		s.view = FounderViewActions
		return s, nil

	case "view":
		isPlatform, devs, apps, apiUsage, revenue, networkScore := fg.GetPlatformSummary()
		if isPlatform {
			s.turnMessages = []string{
				fmt.Sprintf("Developers: %d", devs),
				fmt.Sprintf("Apps: %d", apps),
				fmt.Sprintf("API Calls/mo: %d", apiUsage),
				fmt.Sprintf("Marketplace Revenue: $%s/mo", formatCompactMoney(revenue)),
				fmt.Sprintf("Network Score: %.2f", networkScore),
			}
		} else {
			s.turnMessages = []string{"Platform not launched yet"}
		}
		s.view = FounderViewMain
		return s, nil

	case "launch_marketplace":
		err := fg.LaunchPlatform("marketplace")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Marketplace launched!", "Third-party apps can now be sold"}
		}
		s.view = FounderViewMain
		return s, nil

	case "launch_social":
		err := fg.LaunchPlatform("social")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Social Platform launched!", "Community-driven network effects activated"}
		}
		s.view = FounderViewMain
		return s, nil

	case "launch_data":
		err := fg.LaunchPlatform("data")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Data Platform launched!", "Data exchange and insights marketplace active"}
		}
		s.view = FounderViewMain
		return s, nil

	case "launch_infrastructure":
		err := fg.LaunchPlatform("infrastructure")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Infrastructure Platform launched!", "Developers can now build on your APIs"}
		}
		s.view = FounderViewMain
		return s, nil

	case "invest":
		budget := int64(50000)
		err := fg.InvestInDeveloperProgram(budget)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{fmt.Sprintf("✓ Invested $%s/mo in developer program", formatCompactMoney(budget))}
		}
		s.view = FounderViewMain
		return s, nil
	}

	return s, nil
}

// Security & Compliance
func (s *FounderGameScreen) rebuildSecurityMenu() {
	items := []components.MenuItem{
		{ID: "view", Title: "View Security Posture", Description: "Current security status", Icon: "👁️"},
		{ID: "invest", Title: "Increase Security Budget", Description: "More resources for security", Icon: "💰"},
		{ID: "hire", Title: "Hire Security Team", Description: "Dedicated security engineers", Icon: "👥"},
		{ID: "soc2", Title: "Get SOC2 Certification", Description: "Enterprise compliance ($150k)", Icon: "📜"},
		{ID: "iso27001", Title: "Get ISO 27001", Description: "International security standard ($200k)", Icon: "🌐"},
		{ID: "cancel", Title: "Back", Icon: "←"},
	}

	s.securityMenu = components.NewMenu("SECURITY & COMPLIANCE", items)
	s.securityMenu.SetSize(55, 15)
	s.securityMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleSecuritySelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "cancel":
		s.view = FounderViewActions
		return s, nil

	case "view":
		if fg.SecurityPosture != nil {
			s.turnMessages = []string{
				fmt.Sprintf("Security Score: %d/100", fg.SecurityPosture.SecurityScore),
				fmt.Sprintf("Monthly Budget: $%s", formatCompactMoney(fg.SecurityPosture.SecurityBudget)),
				fmt.Sprintf("Security Team: %d", fg.SecurityPosture.SecurityTeamSize),
				fmt.Sprintf("Vulnerabilities: %d", fg.SecurityPosture.Vulnerabilities),
			}
			if len(fg.SecurityPosture.ComplianceCerts) > 0 {
				s.turnMessages = append(s.turnMessages, fmt.Sprintf("Certifications: %v", fg.SecurityPosture.ComplianceCerts))
			}
		} else {
			s.turnMessages = []string{"No security posture established"}
		}
		s.view = FounderViewMain
		return s, nil

	case "invest":
		budget := int64(25000)
		err := fg.InvestInSecurity(budget)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{fmt.Sprintf("✓ Security budget set to $%s/mo", formatCompactMoney(budget))}
		}
		s.view = FounderViewMain
		return s, nil

	case "hire":
		err := fg.HireSecurityTeam(1)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Hired 1 security engineer"}
		}
		s.view = FounderViewMain
		return s, nil

	case "soc2":
		err := fg.GetComplianceCertification("SOC2")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ SOC2 certification obtained!", "Enterprise customers now more likely to buy"}
		}
		s.view = FounderViewMain
		return s, nil

	case "iso27001":
		err := fg.GetComplianceCertification("ISO27001")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ ISO 27001 certification obtained!", "International sales now easier"}
		}
		s.view = FounderViewMain
		return s, nil
	}

	return s, nil
}

// PR Crisis Management
func (s *FounderGameScreen) rebuildPRCrisisMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{
		{ID: "view", Title: "View PR Status", Description: "Current PR situation", Icon: "👁️"},
	}

	if fg.ActivePRCrisis != nil {
		items = append(items, components.MenuItem{
			ID: "respond_apologize", Title: "Issue Public Apology", Description: "Accept responsibility", Icon: "🙏",
		})
		items = append(items, components.MenuItem{
			ID: "respond_deny", Title: "Deny Allegations", Description: "Dispute the claims", Icon: "❌",
		})
		items = append(items, components.MenuItem{
			ID: "respond_ignore", Title: "No Comment", Description: "Wait for it to blow over", Icon: "🤐",
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Back", Icon: "←",
	})

	s.prCrisisMenu = components.NewMenu("PR CRISIS MANAGEMENT", items)
	s.prCrisisMenu.SetSize(55, 12)
	s.prCrisisMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handlePRCrisisSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "cancel":
		s.view = FounderViewActions
		return s, nil

	case "view":
		if fg.ActivePRCrisis != nil {
			s.turnMessages = []string{
				fmt.Sprintf("Crisis: %s", fg.ActivePRCrisis.Type),
				fmt.Sprintf("Severity: %s", fg.ActivePRCrisis.Severity),
				fmt.Sprintf("Duration: %d months", fg.ActivePRCrisis.DurationMonths),
			}
		} else {
			s.turnMessages = []string{"No active PR crisis. Good news!"}
		}
		s.view = FounderViewMain
		return s, nil

	case "respond_apologize":
		err := fg.RespondToPRCrisis("apologize")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Public apology issued", "Reputation impact minimized"}
		}
		s.view = FounderViewMain
		return s, nil

	case "respond_deny":
		err := fg.RespondToPRCrisis("deny")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Allegations denied", "May backfire if evidence surfaces"}
		}
		s.view = FounderViewMain
		return s, nil

	case "respond_ignore":
		err := fg.RespondToPRCrisis("ignore")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ No comment issued", "Crisis may escalate or fade away"}
		}
		s.view = FounderViewMain
		return s, nil
	}

	return s, nil
}

// Economic Strategy
func (s *FounderGameScreen) rebuildEconomyMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{
		{ID: "view", Title: "View Economic Conditions", Description: "Current market situation", Icon: "👁️"},
	}

	if fg.EconomicEvent != nil && fg.EconomicEvent.Active {
		items = append(items, components.MenuItem{
			ID: "survive_cut", Title: "Cut Costs Aggressively", Description: "Reduce burn by 40%", Icon: "✂️",
		})
		items = append(items, components.MenuItem{
			ID: "survive_raise", Title: "Emergency Fundraise", Description: "Down round if needed", Icon: "💰",
		})
		items = append(items, components.MenuItem{
			ID: "survive_pivot", Title: "Pivot to Profitability", Description: "Focus on unit economics", Icon: "📈",
		})
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Back", Icon: "←",
	})

	s.economyMenu = components.NewMenu("ECONOMIC STRATEGY", items)
	s.economyMenu.SetSize(55, 12)
	s.economyMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleEconomySelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "cancel":
		s.view = FounderViewActions
		return s, nil

	case "view":
		if fg.EconomicEvent != nil && fg.EconomicEvent.Active {
			s.turnMessages = []string{
				fmt.Sprintf("Event: %s", fg.EconomicEvent.Type),
				fmt.Sprintf("Severity: %s", fg.EconomicEvent.Severity),
				fmt.Sprintf("Duration: %d months", fg.EconomicEvent.DurationMonths),
			}
		} else {
			s.turnMessages = []string{"Economy is stable. No major events."}
		}
		s.view = FounderViewMain
		return s, nil

	case "survive_cut":
		err := fg.ExecuteSurvivalStrategy("cut_costs")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Cost cutting measures implemented", "Burn rate reduced significantly"}
		}
		s.view = FounderViewMain
		return s, nil

	case "survive_raise":
		err := fg.ExecuteSurvivalStrategy("emergency_raise")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Emergency funding secured", "May have taken dilution"}
		}
		s.view = FounderViewMain
		return s, nil

	case "survive_pivot":
		err := fg.ExecuteSurvivalStrategy("profitability_pivot")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Pivoted to profitability focus", "Growth may slow but runway extended"}
		}
		s.view = FounderViewMain
		return s, nil
	}

	return s, nil
}

// Succession Planning
func (s *FounderGameScreen) rebuildSuccessionMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{
		{ID: "view", Title: "View Key Person Risks", Description: "Identify critical dependencies", Icon: "👁️"},
		{ID: "assess", Title: "Assess Risks", Description: "Evaluate key person dependencies", Icon: "🔍"},
	}

	if fg.KeyPersonRisks != nil && len(fg.KeyPersonRisks) > 0 {
		for i, risk := range fg.KeyPersonRisks {
			if risk.SuccessionReady {
				continue
			}
			items = append(items, components.MenuItem{
				ID:          fmt.Sprintf("plan_%d", i),
				Title:       fmt.Sprintf("Create Plan for %s", risk.PersonName),
				Description: fmt.Sprintf("Risk: %s", risk.RiskLevel),
				Icon:        "📋",
			})
		}
	}

	items = append(items, components.MenuItem{
		ID: "cancel", Title: "Back", Icon: "←",
	})

	s.successionMenu = components.NewMenu("SUCCESSION PLANNING", items)
	s.successionMenu.SetSize(55, 15)
	s.successionMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleSuccessionSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "cancel":
		s.view = FounderViewActions
		return s, nil

	case "view":
		if fg.KeyPersonRisks != nil && len(fg.KeyPersonRisks) > 0 {
			var msgs []string
			for _, risk := range fg.KeyPersonRisks {
				status := "⚠️ No plan"
				if risk.SuccessionReady {
					status = "✓ Plan in place"
				}
				msgs = append(msgs, fmt.Sprintf("%s (%s): %s - %s", risk.PersonName, risk.Role, risk.RiskLevel, status))
			}
			s.turnMessages = msgs
		} else {
			s.turnMessages = []string{"No key person risks identified yet"}
		}
		s.view = FounderViewMain
		return s, nil

	case "assess":
		fg.AssessKeyPersonRisks()
		s.turnMessages = []string{"✓ Key person risks assessed", "Check View for details"}
		s.view = FounderViewMain
		return s, nil
	}

	if strings.HasPrefix(id, "plan_") {
		idxStr := strings.TrimPrefix(id, "plan_")
		idx, _ := strconv.Atoi(idxStr)
		if idx >= 0 && idx < len(fg.KeyPersonRisks) {
			risk := fg.KeyPersonRisks[idx]
			err := fg.CreateSuccessionPlan(risk.PersonName, "Hired backup")
			if err != nil {
				s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
			} else {
				s.turnMessages = []string{fmt.Sprintf("✓ Succession plan created for %s", risk.PersonName)}
			}
		}
		s.view = FounderViewMain
		return s, nil
	}

	return s, nil
}

// ============================================================================
// NEW FEATURE MENUS AND HANDLERS
// ============================================================================

// Strategic Opportunity
func (s *FounderGameScreen) rebuildStrategicOpportunityMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{}

	if fg.PendingOpportunity != nil {
		opp := fg.PendingOpportunity
		items = append(items,
			components.MenuItem{ID: "info", Title: opp.Title, Description: opp.Description, Disabled: true, Icon: "💡"},
			components.MenuItem{ID: "accept", Title: "Accept Opportunity", Description: fmt.Sprintf("Benefit: %s", opp.Benefit), Icon: "✓"},
			components.MenuItem{ID: "decline", Title: "Decline Opportunity", Description: "Pass on this opportunity", Icon: "✗"},
		)
	}

	items = append(items, components.MenuItem{ID: "cancel", Title: "Back", Icon: "←"})

	s.strategicOpportunityMenu = components.NewMenu("STRATEGIC OPPORTUNITY", items)
	s.strategicOpportunityMenu.SetSize(60, 12)
	s.strategicOpportunityMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleStrategicOpportunitySelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	if fg.PendingOpportunity == nil {
		s.view = FounderViewActions
		return s, nil
	}

	opp := fg.PendingOpportunity

	switch id {
	case "accept":
		// Deduct cost
		if opp.Cost > 0 {
			if fg.Cash < opp.Cost {
				s.turnMessages = []string{fmt.Sprintf("❌ Not enough cash ($%s needed)", formatCompactMoney(opp.Cost))}
				s.view = FounderViewActions
				return s, nil
			}
			fg.Cash -= opp.Cost
		}

		msgs := []string{fmt.Sprintf("✓ Accepted: %s", opp.Title)}

		// Apply actual effects based on type
		switch opp.Type {
		case "press":
			// Press coverage brings customers and reduces CAC
			newCustomers := 5 + rand.Intn(15)
			newMRR := int64(newCustomers) * fg.AvgDealSize
			fg.Customers += newCustomers
			fg.DirectCustomers += newCustomers
			fg.DirectMRR += newMRR
			fg.BaseCAC = int64(float64(fg.BaseCAC) * 0.85) // 15% CAC reduction
			msgs = append(msgs,
				fmt.Sprintf("   +%d new customers (+$%s/mo MRR)", newCustomers, formatCompactMoney(newMRR)),
				"   -15% customer acquisition cost (brand awareness)")

		case "enterprise_pilot":
			// Big enterprise deal with 80% success chance
			if rand.Float64() < 0.80 {
				dealMRR := (50000 + rand.Int63n(150000)) / 12
				fg.Customers += 1
				fg.DirectCustomers += 1
				fg.DirectMRR += dealMRR
				msgs = append(msgs,
					fmt.Sprintf("   Enterprise pilot SUCCESS — +$%s/mo MRR", formatCompactMoney(dealMRR)),
					"   +1 enterprise customer (reference account)")
			} else {
				msgs = append(msgs, "   Enterprise pilot did not convert — they went with another vendor")
			}

		case "bridge_round":
			// Quick cash injection with equity dilution
			amount := 200000 + rand.Int63n(500000)
			equity := 3.0 + rand.Float64()*5.0
			fg.Cash += amount
			fg.EquityGivenAway += equity
			fg.CalculateRunway()
			msgs = append(msgs,
				fmt.Sprintf("   +$%s cash injection", formatCompactMoney(amount)),
				fmt.Sprintf("   -%.1f%% equity dilution", equity),
				fmt.Sprintf("   New runway: %d months", fg.CashRunwayMonths))

		case "conference":
			// Conference brings leads that convert to customers over time
			newCustomers := 3 + rand.Intn(8)
			newMRR := int64(newCustomers) * fg.AvgDealSize
			fg.Customers += newCustomers
			fg.DirectCustomers += newCustomers
			fg.DirectMRR += newMRR
			fg.BaseCAC = int64(float64(fg.BaseCAC) * 0.90) // 10% CAC reduction
			msgs = append(msgs,
				fmt.Sprintf("   +%d customers (+$%s/mo MRR)", newCustomers, formatCompactMoney(newMRR)),
				"   -10% CAC from industry credibility")

		case "talent":
			// Hire a star engineer with high impact
			eng := founder.Employee{
				Name:        "Star Engineer",
				Role:        founder.RoleEngineer,
				MonthlyCost: 200000 / 12, // $200k/yr
				Impact:      2.0,          // 2x normal engineer
				MonthHired:  fg.Turn,
			}
			fg.Team.Engineers = append(fg.Team.Engineers, eng)
			fg.CalculateTeamCost()
			fg.CalculateRunway()
			msgs = append(msgs,
				"   Hired Star Engineer (2x impact)",
				fmt.Sprintf("   Salary: $200k/yr | New runway: %d months", fg.CashRunwayMonths))

		case "competitor_distress":
			// Acquire competitor's customers
			newCustomers := 15 + rand.Intn(25)
			newMRR := int64(newCustomers) * fg.AvgDealSize
			fg.Customers += newCustomers
			fg.DirectCustomers += newCustomers
			fg.DirectMRR += newMRR
			// Deactivate a random competitor if any exist
			for i := range fg.Competitors {
				if fg.Competitors[i].Active {
					fg.Competitors[i].Active = false
					msgs = append(msgs, fmt.Sprintf("   Eliminated competitor: %s", fg.Competitors[i].Name))
					break
				}
			}
			msgs = append(msgs, fmt.Sprintf("   +%d customers acquired (+$%s/mo MRR)", newCustomers, formatCompactMoney(newMRR)))

		case "api_integration":
			// API partner brings recurring customers
			newCustomers := 10 + rand.Intn(20)
			newMRR := int64(newCustomers) * fg.AvgDealSize
			fg.Customers += newCustomers
			fg.DirectCustomers += newCustomers
			fg.DirectMRR += newMRR
			fg.MonthlyGrowthRate += 0.02 // Ongoing growth boost
			msgs = append(msgs,
				fmt.Sprintf("   +%d customers via API integration (+$%s/mo MRR)", newCustomers, formatCompactMoney(newMRR)),
				"   +2% ongoing monthly growth from partner channel")

		case "govt_contract":
			// Government contract — guaranteed revenue
			contractMRR := int64(20000 + rand.Intn(80000))
			fg.Customers += 1
			fg.DirectCustomers += 1
			fg.DirectMRR += contractMRR
			msgs = append(msgs,
				fmt.Sprintf("   Government contract: +$%s/mo MRR", formatCompactMoney(contractMRR)),
				"   3-year guaranteed revenue")

		case "influencer":
			newCustomers := 8 + rand.Intn(20)
			newMRR := int64(newCustomers) * fg.AvgDealSize
			fg.Customers += newCustomers
			fg.DirectCustomers += newCustomers
			fg.DirectMRR += newMRR
			msgs = append(msgs, fmt.Sprintf("   Influencer campaign: +%d customers (+$%s/mo MRR)", newCustomers, formatCompactMoney(newMRR)))

		case "patent":
			// Patent creates competitive moat
			for i := range fg.Competitors {
				if fg.Competitors[i].Active {
					fg.Competitors[i].MarketShare *= 0.8 // Reduce all competitor share
				}
			}
			msgs = append(msgs, "   Patent granted — competitors' market share reduced 20%")

		case "university_partnership":
			// Cheap talent pipeline + customers
			newCustomers := 5 + rand.Intn(10)
			newMRR := int64(newCustomers) * fg.AvgDealSize
			fg.Customers += newCustomers
			fg.DirectCustomers += newCustomers
			fg.DirectMRR += newMRR
			msgs = append(msgs,
				fmt.Sprintf("   +%d customers from university network (+$%s/mo MRR)", newCustomers, formatCompactMoney(newMRR)),
				"   Improved recruiting pipeline")
		}

		// Sync MRR to reflect new customers
		fg.MRR = fg.DirectMRR + fg.AffiliateMRR

		if opp.Cost > 0 {
			msgs = append(msgs, fmt.Sprintf("   Cost: $%s", formatCompactMoney(opp.Cost)))
		}

		fg.CalculateRunway()
		s.turnMessages = msgs
		fg.PendingOpportunity = nil

	case "decline":
		s.turnMessages = []string{
			fmt.Sprintf("✗ Declined: %s", opp.Title),
		}
		fg.PendingOpportunity = nil
	}

	s.view = FounderViewMain
	return s, nil
}

// Content Marketing
func (s *FounderGameScreen) rebuildContentMarketingMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{}

	if fg.ContentProgram == nil {
		items = append(items,
			components.MenuItem{ID: "launch_10k", Title: "Launch ($10k/mo)", Description: "Basic content program", Icon: "📝"},
			components.MenuItem{ID: "launch_25k", Title: "Launch ($25k/mo)", Description: "Medium content program", Icon: "📝"},
			components.MenuItem{ID: "launch_50k", Title: "Launch ($50k/mo)", Description: "Premium content program", Icon: "📝"},
		)
	} else {
		items = append(items,
			components.MenuItem{ID: "view", Title: "View Program Status", Description: fmt.Sprintf("SEO: %d, Traffic: %d", fg.ContentProgram.SEOScore, fg.ContentProgram.OrganicTraffic), Icon: "👁️"},
			components.MenuItem{ID: "end", Title: "End Program", Description: "Stop content marketing", Icon: "🛑"},
		)
	}

	items = append(items, components.MenuItem{ID: "cancel", Title: "Back", Icon: "←"})

	s.contentMarketingMenu = components.NewMenu("CONTENT MARKETING", items)
	s.contentMarketingMenu.SetSize(55, 12)
	s.contentMarketingMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleContentMarketingSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	switch id {
	case "launch_10k":
		err := fg.LaunchContentProgram(10000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched content marketing program ($10k/mo)"}
		}
	case "launch_25k":
		err := fg.LaunchContentProgram(25000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched content marketing program ($25k/mo)"}
		}
	case "launch_50k":
		err := fg.LaunchContentProgram(50000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched content marketing program ($50k/mo)"}
		}
	case "view":
		if fg.ContentProgram != nil {
			s.turnMessages = []string{
				"📝 Content Marketing Status:",
				fmt.Sprintf("   Monthly Budget: $%s", formatCompactMoney(fg.ContentProgram.MonthlyBudget)),
				fmt.Sprintf("   SEO Score: %d/100", fg.ContentProgram.SEOScore),
				fmt.Sprintf("   Organic Traffic: %d visitors/mo", fg.ContentProgram.OrganicTraffic),
				fmt.Sprintf("   Content Quality: %.0f/100", fg.ContentProgram.ContentQuality*100),
				fmt.Sprintf("   Inbound Leads: %d this month", fg.ContentProgram.InboundLeads),
			}
		}
	case "end":
		fg.EndContentProgram()
		s.turnMessages = []string{"✓ Ended content marketing program"}
	}

	s.view = FounderViewMain
	return s, nil
}

// CS Playbooks
func (s *FounderGameScreen) rebuildCSPlaybooksMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{
		{ID: "view", Title: "View Active Playbooks", Description: fmt.Sprintf("%d playbooks running", len(fg.CSPlaybooks)), Icon: "👁️"},
	}

	// Check which playbooks are not yet launched
	hasOnboarding := false
	hasHealth := false
	hasUpsell := false
	hasRenewal := false
	hasChurnPrevention := false

	for _, pb := range fg.CSPlaybooks {
		switch pb.Name {
		case "Onboarding":
			hasOnboarding = true
		case "Health Monitoring":
			hasHealth = true
		case "Upsell":
			hasUpsell = true
		case "Renewal":
			hasRenewal = true
		case "Churn Prevention":
			hasChurnPrevention = true
		}
	}

	if !hasOnboarding {
		items = append(items, components.MenuItem{ID: "launch_onboarding", Title: "Launch Onboarding", Description: "Help customers get started", Icon: "🚀"})
	}
	if !hasHealth {
		items = append(items, components.MenuItem{ID: "launch_health", Title: "Launch Health Monitoring", Description: "Track customer health scores", Icon: "❤️"})
	}
	if !hasUpsell {
		items = append(items, components.MenuItem{ID: "launch_upsell", Title: "Launch Upsell Program", Description: "Identify expansion opportunities", Icon: "📈"})
	}
	if !hasRenewal {
		items = append(items, components.MenuItem{ID: "launch_renewal", Title: "Launch Renewal Program", Description: "Proactive renewal management", Icon: "🔄"})
	}
	if !hasChurnPrevention {
		items = append(items, components.MenuItem{ID: "launch_churn", Title: "Launch Churn Prevention", Description: "Save at-risk customers", Icon: "🛡️"})
	}

	items = append(items, components.MenuItem{ID: "cancel", Title: "Back", Icon: "←"})

	s.csPlaybooksMenu = components.NewMenu("CS PLAYBOOKS", items)
	s.csPlaybooksMenu.SetSize(55, 15)
	s.csPlaybooksMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleCSPlaybooksSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	switch id {
	case "view":
		if len(fg.CSPlaybooks) == 0 {
			s.turnMessages = []string{"No CS playbooks active yet"}
		} else {
			msgs := []string{"📋 Active CS Playbooks:"}
			for _, pb := range fg.CSPlaybooks {
				msgs = append(msgs, fmt.Sprintf("   • %s (Budget: $%s/mo, NPS: %d)", pb.Name, formatCompactMoney(pb.MonthlyBudget), pb.NPSScore))
			}
			s.turnMessages = msgs
		}
	case "launch_onboarding":
		err := fg.LaunchCSPlaybook("Onboarding", 5000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched Onboarding playbook"}
		}
	case "launch_health":
		err := fg.LaunchCSPlaybook("Health Monitoring", 5000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched Health Monitoring playbook"}
		}
	case "launch_upsell":
		err := fg.LaunchCSPlaybook("Upsell", 5000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched Upsell playbook"}
		}
	case "launch_renewal":
		err := fg.LaunchCSPlaybook("Renewal", 5000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched Renewal playbook"}
		}
	case "launch_churn":
		err := fg.LaunchCSPlaybook("Churn Prevention", 5000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched Churn Prevention playbook"}
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// Competitive Intelligence
func (s *FounderGameScreen) rebuildCompetitiveIntelMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{}

	if fg.CompetitiveIntel == nil {
		items = append(items,
			components.MenuItem{ID: "hire_80k", Title: "Hire Analyst ($80k/yr)", Description: "Basic competitive intelligence", Icon: "🕵️"},
			components.MenuItem{ID: "hire_100k", Title: "Hire Analyst ($100k/yr)", Description: "Experienced analyst", Icon: "🕵️"},
			components.MenuItem{ID: "hire_120k", Title: "Hire Analyst ($120k/yr)", Description: "Senior analyst", Icon: "🕵️"},
		)
	} else {
		items = append(items,
			components.MenuItem{ID: "view", Title: "View Intel Reports", Description: fmt.Sprintf("%d reports available", len(fg.CompetitiveIntel.IntelReports)), Icon: "👁️"},
		)

		// Add commission options for each competitor
		for i, comp := range fg.Competitors {
			items = append(items, components.MenuItem{
				ID:          fmt.Sprintf("report_%d", i),
				Title:       fmt.Sprintf("Commission Report: %s", comp.Name),
				Description: fmt.Sprintf("$25-50k, threat: %s", comp.Threat),
				Icon:        "📊",
			})
		}
	}

	items = append(items, components.MenuItem{ID: "cancel", Title: "Back", Icon: "←"})

	s.competitiveIntelMenu = components.NewMenu("COMPETITIVE INTELLIGENCE", items)
	s.competitiveIntelMenu.SetSize(55, 15)
	s.competitiveIntelMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleCompetitiveIntelSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	switch id {
	case "hire_80k":
		err := fg.LaunchCompetitiveIntel(80000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Hired competitive intelligence analyst ($80k/yr)"}
		}
	case "hire_100k":
		err := fg.LaunchCompetitiveIntel(100000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Hired competitive intelligence analyst ($100k/yr)"}
		}
	case "hire_120k":
		err := fg.LaunchCompetitiveIntel(120000)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Hired competitive intelligence analyst ($120k/yr)"}
		}
	case "view":
		if fg.CompetitiveIntel != nil && len(fg.CompetitiveIntel.IntelReports) > 0 {
			msgs := []string{"🕵️ Intel Reports:"}
			for _, report := range fg.CompetitiveIntel.IntelReports {
				recentMove := "No recent moves"
				if len(report.RecentMoves) > 0 {
					recentMove = report.RecentMoves[0]
				}
				msgs = append(msgs, fmt.Sprintf("   • %s (Threat: %s): %s", report.CompetitorName, report.ThreatLevel, recentMove))
			}
			s.turnMessages = msgs
		} else {
			s.turnMessages = []string{"No intel reports yet. Commission reports on competitors."}
		}
	default:
		if strings.HasPrefix(id, "report_") {
			idxStr := strings.TrimPrefix(id, "report_")
			idx, _ := strconv.Atoi(idxStr)
			if idx >= 0 && idx < len(fg.Competitors) {
				comp := fg.Competitors[idx]
				err := fg.CommissionIntelReport(comp.Name, 35000) // Average cost
				if err != nil {
					s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
				} else {
					s.turnMessages = []string{fmt.Sprintf("✓ Commissioned intel report on %s ($35k)", comp.Name)}
				}
			}
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// Referral Program
func (s *FounderGameScreen) rebuildReferralProgramMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{}

	if fg.ReferralProgram == nil {
		items = append(items,
			components.MenuItem{ID: "launch_100", Title: "Launch ($100 reward)", Description: "Cash reward per referral", Icon: "🎁"},
			components.MenuItem{ID: "launch_250", Title: "Launch ($250 reward)", Description: "Cash reward per referral", Icon: "🎁"},
			components.MenuItem{ID: "launch_500", Title: "Launch ($500 reward)", Description: "Premium cash reward", Icon: "🎁"},
			components.MenuItem{ID: "launch_credit", Title: "Launch (Account Credit)", Description: "$200 account credit per referral", Icon: "💳"},
		)
	} else {
		items = append(items,
			components.MenuItem{ID: "view", Title: "View Program Status", Icon: "👁️"},
			components.MenuItem{ID: "end", Title: "End Referral Program", Icon: "🛑"},
		)
	}

	items = append(items, components.MenuItem{ID: "cancel", Title: "Back", Icon: "←"})

	s.referralProgramMenu = components.NewMenu("REFERRAL PROGRAM", items)
	s.referralProgramMenu.SetSize(55, 12)
	s.referralProgramMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleReferralProgramSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	switch id {
	case "launch_100":
		err := fg.LaunchReferralProgram(100, "cash")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched referral program ($100 cash reward)"}
		}
	case "launch_250":
		err := fg.LaunchReferralProgram(250, "cash")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched referral program ($250 cash reward)"}
		}
	case "launch_500":
		err := fg.LaunchReferralProgram(500, "cash")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched referral program ($500 cash reward)"}
		}
	case "launch_credit":
		err := fg.LaunchReferralProgram(200, "credit")
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Launched referral program ($200 account credit)"}
		}
	case "view":
		if fg.ReferralProgram != nil {
			s.turnMessages = []string{
				"🎁 Referral Program Status:",
				fmt.Sprintf("   Reward: $%s (%s)", formatCompactMoney(fg.ReferralProgram.RewardPerReferral), fg.ReferralProgram.RewardType),
				fmt.Sprintf("   Total Referrals: %d", fg.ReferralProgram.TotalReferrals),
				fmt.Sprintf("   Customers Acquired: %d", fg.ReferralProgram.CustomersAcquired),
				fmt.Sprintf("   Monthly Budget: $%s", formatCompactMoney(fg.ReferralProgram.MonthlyBudget)),
			}
		}
	case "end":
		err := fg.EndReferralProgram()
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{"✓ Ended referral program"}
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// Technical Debt
func (s *FounderGameScreen) rebuildTechDebtMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{}

	debtLevel := 0
	velocity := 1.0
	if fg.TechnicalDebt != nil {
		debtLevel = fg.TechnicalDebt.CurrentLevel
		velocity = fg.TechnicalDebt.VelocityImpact
	}

	items = append(items,
		components.MenuItem{ID: "view", Title: "View Tech Debt Status", Description: fmt.Sprintf("Level: %d/100, Velocity: %.0f%%", debtLevel, velocity*100), Icon: "👁️"},
	)

	if debtLevel > 20 {
		items = append(items,
			components.MenuItem{ID: "refactor_small", Title: "Small Refactor ($50k)", Description: "Reduce debt by ~10 points", Icon: "🔧"},
			components.MenuItem{ID: "refactor_medium", Title: "Medium Refactor ($100k)", Description: "Reduce debt by ~20 points", Icon: "🔧"},
			components.MenuItem{ID: "refactor_large", Title: "Large Refactor ($200k)", Description: "Reduce debt by ~35 points", Icon: "🔧"},
		)
	}

	items = append(items, components.MenuItem{ID: "cancel", Title: "Back", Icon: "←"})

	s.techDebtMenu = components.NewMenu("TECHNICAL DEBT", items)
	s.techDebtMenu.SetSize(55, 12)
	s.techDebtMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleTechDebtSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	switch id {
	case "view":
		if fg.TechnicalDebt != nil {
			td := fg.TechnicalDebt
			msgs := []string{
				"🔧 Technical Debt Status:",
				fmt.Sprintf("   Current Level: %d/100", td.CurrentLevel),
				fmt.Sprintf("   Velocity Impact: %.0f%%", td.VelocityImpact*100),
				fmt.Sprintf("   Bug Frequency: %.1f%%", td.BugFrequency*100),
				fmt.Sprintf("   Months Since Refactor: %d", td.MonthsSinceRefactor),
			}
			if td.ScalingProblems {
				msgs = append(msgs, "   ⚠️  SCALING PROBLEMS ACTIVE")
			}
			s.turnMessages = msgs
		} else {
			s.turnMessages = []string{"Technical debt not yet being tracked"}
		}
	case "refactor_small":
		engineers := len(fg.Team.Engineers)
		if engineers == 0 {
			engineers = 1
		}
		err := fg.RefactorTechDebt(50000, engineers)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				"✓ Small refactor complete ($50k)",
				fmt.Sprintf("   New debt level: %d/100", fg.TechnicalDebt.CurrentLevel),
			}
		}
	case "refactor_medium":
		engineers := len(fg.Team.Engineers)
		if engineers == 0 {
			engineers = 1
		}
		err := fg.RefactorTechDebt(100000, engineers)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				"✓ Medium refactor complete ($100k)",
				fmt.Sprintf("   New debt level: %d/100", fg.TechnicalDebt.CurrentLevel),
			}
		}
	case "refactor_large":
		engineers := len(fg.Team.Engineers)
		if engineers == 0 {
			engineers = 1
		}
		err := fg.RefactorTechDebt(200000, engineers)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				"✓ Large refactor complete ($200k)",
				fmt.Sprintf("   New debt level: %d/100", fg.TechnicalDebt.CurrentLevel),
			}
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// End Affiliate Program
func (s *FounderGameScreen) rebuildEndAffiliateMenu() {
	items := []components.MenuItem{
		{ID: "transition", Title: "End & Transition Customers", Description: "Convert affiliate customers to direct", Icon: "✓"},
		{ID: "end_only", Title: "End Program Only", Description: "Affiliate customers will churn", Icon: "⚠️"},
		{ID: "cancel", Title: "Back", Icon: "←"},
	}

	s.endAffiliateMenu = components.NewMenu("END AFFILIATE PROGRAM", items)
	s.endAffiliateMenu.SetSize(55, 10)
	s.endAffiliateMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleEndAffiliateSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewActions
		return s, nil
	}

	switch id {
	case "transition":
		err := fg.EndAffiliateProgram(true)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				"✓ Ended affiliate program",
				"  Affiliate customers converted to direct customers",
			}
		}
	case "end_only":
		err := fg.EndAffiliateProgram(false)
		if err != nil {
			s.turnMessages = []string{fmt.Sprintf("❌ Error: %v", err)}
		} else {
			s.turnMessages = []string{
				"✓ Ended affiliate program",
				"  ⚠️  Affiliate customers have churned",
			}
		}
	}

	s.view = FounderViewMain
	return s, nil
}

// View renders the founder game screen
func (s *FounderGameScreen) View() string {
	switch s.view {
	case FounderViewActions:
		return s.renderActions()
	case FounderViewHiring:
		return s.renderHiring()
	case FounderViewExecOffer:
		return s.renderExecOffer()
	case FounderViewHiringMarket:
		return s.renderHiringMarket()
	case FounderViewFiring:
		return s.renderFiring()
	case FounderViewMarketing:
		return s.renderMarketing()
	case FounderViewFunding:
		return s.renderFunding()
	case FounderViewFundingTerms:
		return s.renderFundingTerms()
	case FounderViewPartnership:
		return s.renderPartnership()
	case FounderViewAffiliate:
		return s.renderAffiliate()
	case FounderViewCompetitors:
		return s.renderCompetitors()
	case FounderViewExpansion:
		return s.renderExpansion()
	case FounderViewPivot:
		return s.renderPivot()
	case FounderViewBoard:
		return s.renderBoard()
	case FounderViewBoardAction:
		return s.renderExpandPool()
	case FounderViewBuyback:
		return s.renderBuyback()
	case FounderViewBuybackConfirm:
		return s.renderBuybackConfirm()
	case FounderViewTeamRoster:
		return s.renderTeamRoster()
	case FounderViewCustomers:
		return s.renderCustomers()
	case FounderViewFinancials:
		return s.renderFinancials()
	case FounderViewCapTable:
		return s.renderCapTable()
	case FounderViewExit:
		return s.renderExit()
	case FounderViewConfirmExit:
		return s.renderConfirmExit()
	// Phase 4: Advanced views
	case FounderViewRoadmap:
		return s.renderRoadmap()
	case FounderViewRoadmapStart:
		return s.renderRoadmapStart()
	case FounderViewSegments:
		return s.renderSegments()
	case FounderViewPricing:
		return s.renderPricing()
	case FounderViewAcquisitions:
		return s.renderAcquisitions()
	case FounderViewPlatform:
		return s.renderPlatform()
	case FounderViewSecurity:
		return s.renderSecurity()
	case FounderViewPRCrisis:
		return s.renderPRCrisis()
	case FounderViewEconomy:
		return s.renderEconomy()
	case FounderViewSuccession:
		return s.renderSuccession()
	case FounderViewSalesPipeline:
		return s.renderSalesPipeline()
	// New feature views
	case FounderViewStrategicOpportunity:
		return s.renderStrategicOpportunity()
	case FounderViewContentMarketing:
		return s.renderContentMarketing()
	case FounderViewCSPlaybooks:
		return s.renderCSPlaybooks()
	case FounderViewCompetitiveIntel:
		return s.renderCompetitiveIntel()
	case FounderViewReferralProgram:
		return s.renderReferralProgram()
	case FounderViewTechDebt:
		return s.renderTechDebt()
	case FounderViewAcquisitionOffer:
		return s.renderAcquisitionOffer()
	case FounderViewAdvisorExpertise:
		return s.renderAdvisorExpertise()
	case FounderViewAdvisorConfirm:
		return s.renderAdvisorConfirm()
	case FounderViewRemoveAdvisor:
		return s.renderRemoveAdvisor()
	case FounderViewBoardTable:
		return s.renderBoardTable()
	case FounderViewEndAffiliate:
		return s.renderEndAffiliate()
	case FounderViewConfirmQuit:
		return s.renderConfirmQuit()
	case FounderViewEngineerRealloc:
		return s.renderEngineerRealloc()
	default:
		return s.renderMain()
	}
}

func (s *FounderGameScreen) renderMain() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	header := fmt.Sprintf("🚀 %s - MONTH %d", fg.CompanyName, fg.Turn)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render(header)))
	b.WriteString("\n")

	// Low cash warning
	if fg.NeedsLowCashWarning() {
		warnBox := lipgloss.NewStyle().
			Foreground(styles.Red).
			Bold(true).
			Width(70).
			Align(lipgloss.Center)
		warnText := fmt.Sprintf("⚠️ LOW CASH WARNING: $%s remaining | %d months runway",
			formatCompactMoney(fg.Cash), fg.CashRunwayMonths)
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(warnBox.Render(warnText)))
	}
	b.WriteString("\n")

	// Main layout. Center with margin only so bordered boxes never reflow (avoids disjointed bottom border).
	leftPanel := s.renderCompanyPanel()
	rightPanel := s.renderMetricsPanel()
	panelRow := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, "  ", rightPanel)
	rowWidth := 35 + 2 + 35
	margin := (s.width - rowWidth) / 2
	if margin < 0 {
		margin = 0
	}
	b.WriteString(lipgloss.NewStyle().MarginLeft(margin).Render(panelRow))
	b.WriteString("\n")

	// SaaS Metrics + Board Sentiment row
	saasPanel := s.renderSaaSMetricsPanel()
	boardPanel := s.renderBoardSentimentPanel()
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, saasPanel, "  ", boardPanel)
	b.WriteString(lipgloss.NewStyle().MarginLeft(margin).Render(row2))
	b.WriteString("\n")

	// Monthly highlights
	highlights := fg.GenerateMonthlyHighlights()
	if len(highlights) > 0 {
		b.WriteString(s.renderHighlights(highlights))
		b.WriteString("\n")
	}

	// News
	if len(s.turnMessages) > 0 {
		b.WriteString(s.renderNews())
		b.WriteString("\n")
	}

	// Strategic opportunity notice
	if fg.PendingOpportunity != nil {
		oppStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Width(s.width).Align(lipgloss.Center)
		b.WriteString(oppStyle.Render(fmt.Sprintf("💡 Strategic Opportunity: %s (expires in %d mo)",
			fg.PendingOpportunity.Title, fg.PendingOpportunity.ExpiresIn)))
		b.WriteString("\n")
	}

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Foreground(styles.White).
		Background(styles.DarkGray).
		Width(70).
		Padding(0, 2)

	// Calculate valuation
	arr := fg.MRR * 12
	multiple := 10.0
	if fg.MonthlyGrowthRate > 0.05 {
		multiple = 15.0
	}
	valuation := int64(float64(arr) * multiple)

	runway := fmt.Sprintf("%d mo", fg.CashRunwayMonths)
	if fg.CashRunwayMonths < 0 {
		runway = "∞"
	}

	status := fmt.Sprintf("💰 Cash: $%s  |  📊 Valuation: $%s  |  ⏳ Runway: %s",
		formatCompactMoney(fg.Cash),
		formatCompactMoney(valuation),
		runway)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(statusStyle.Render(status)))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter actions • n next month • esc/q quit"))

	return b.String()
}

func (s *FounderGameScreen) renderCompanyPanel() string {
	fg := s.gameData.FounderState

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(0, 1).
		Width(35)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
	b.WriteString(titleStyle.Render("🏢 COMPANY"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)

	b.WriteString(labelStyle.Render("Founder: "))
	b.WriteString(fg.FounderName)
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Category: "))
	b.WriteString(fg.Category)
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Team Size: "))
	b.WriteString(fmt.Sprintf("%d", fg.Team.TotalEmployees))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("MRR: "))
	b.WriteString(fmt.Sprintf("$%s", formatCompactMoney(fg.MRR)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Your Equity: "))
	b.WriteString(fmt.Sprintf("%.1f%%", 100.0-fg.EquityGivenAway-fg.EquityPool))
	b.WriteString("\n")

	// Customer health summary
	healthy, atRisk, critical, _, criticalMRR := fg.GetCustomerHealthSegments()
	if fg.Customers > 0 {
		b.WriteString("\n")
		healthStyle := lipgloss.NewStyle().Foreground(styles.Green)
		riskStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
		critStyle := lipgloss.NewStyle().Foreground(styles.Red)
		b.WriteString(healthStyle.Render(fmt.Sprintf("✓%d", healthy)))
		b.WriteString(" ")
		b.WriteString(riskStyle.Render(fmt.Sprintf("⚠%d", atRisk)))
		b.WriteString(" ")
		b.WriteString(critStyle.Render(fmt.Sprintf("✗%d", critical)))
		if criticalMRR > 0 {
			b.WriteString(critStyle.Render(fmt.Sprintf(" ($%s risk)", formatCompactMoney(criticalMRR))))
		}
		b.WriteString("\n")
	}

	// Active partnerships
	activePartners := 0
	for _, p := range fg.Partnerships {
		if p.Status == "active" {
			activePartners++
		}
	}
	if activePartners > 0 {
		b.WriteString(labelStyle.Render("Partnerships: "))
		b.WriteString(fmt.Sprintf("%d active", activePartners))
		b.WriteString("\n")
	}

	return panelStyle.Render(b.String())
}

func (s *FounderGameScreen) renderMetricsPanel() string {
	fg := s.gameData.FounderState

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(0, 1).
		Width(35)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	b.WriteString(titleStyle.Render("📊 METRICS"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)

	b.WriteString(labelStyle.Render("Customers: "))
	b.WriteString(fmt.Sprintf("%d", fg.Customers))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Growth: "))
	growthStyle := lipgloss.NewStyle().Foreground(styles.Green)
	if fg.MonthlyGrowthRate < 0 {
		growthStyle = lipgloss.NewStyle().Foreground(styles.Red)
	}
	b.WriteString(growthStyle.Render(fmt.Sprintf("%.1f%%", fg.MonthlyGrowthRate*100)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Churn: "))
	b.WriteString(fmt.Sprintf("%.1f%%/mo", fg.CustomerChurnRate*100))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Product Maturity: "))
	b.WriteString(fmt.Sprintf("%.0f%%", fg.ProductMaturity*100))
	b.WriteString("\n")

	return panelStyle.Render(b.String())
}

func (s *FounderGameScreen) renderSaaSMetricsPanel() string {
	fg := s.gameData.FounderState

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(35)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	b.WriteString(titleStyle.Render("📈 SaaS METRICS"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)

	// LTV:CAC
	ltvCac := fg.CalculateLTVToCAC()
	b.WriteString(labelStyle.Render("LTV:CAC: "))
	ltvStyle := lipgloss.NewStyle().Foreground(styles.Red)
	if ltvCac > 3 {
		ltvStyle = lipgloss.NewStyle().Foreground(styles.Green)
	} else if ltvCac > 1 {
		ltvStyle = lipgloss.NewStyle().Foreground(styles.Yellow)
	}
	b.WriteString(ltvStyle.Render(fmt.Sprintf("%.1fx", ltvCac)))
	b.WriteString("\n")

	// CAC Payback
	cacPayback := fg.CalculateCACPayback()
	b.WriteString(labelStyle.Render("CAC Payback: "))
	paybackStyle := lipgloss.NewStyle().Foreground(styles.Red)
	if cacPayback > 0 && cacPayback < 12 {
		paybackStyle = lipgloss.NewStyle().Foreground(styles.Green)
	} else if cacPayback > 0 && cacPayback < 24 {
		paybackStyle = lipgloss.NewStyle().Foreground(styles.Yellow)
	}
	if cacPayback > 0 {
		b.WriteString(paybackStyle.Render(fmt.Sprintf("%.0f mo", cacPayback)))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(styles.Gray).Render("N/A"))
	}
	b.WriteString("\n")

	// Rule of 40
	ruleOf40 := fg.CalculateRuleOf40()
	b.WriteString(labelStyle.Render("Rule of 40: "))
	ro40Style := lipgloss.NewStyle().Foreground(styles.Red)
	if ruleOf40 > 40 {
		ro40Style = lipgloss.NewStyle().Foreground(styles.Green)
	} else if ruleOf40 > 20 {
		ro40Style = lipgloss.NewStyle().Foreground(styles.Yellow)
	}
	b.WriteString(ro40Style.Render(fmt.Sprintf("%.0f%%", ruleOf40)))
	b.WriteString("\n")

	// Burn Multiple
	burnMult := fg.CalculateBurnMultiple()
	b.WriteString(labelStyle.Render("Burn Multiple: "))
	burnStyle := lipgloss.NewStyle().Foreground(styles.Red)
	if burnMult > 0 && burnMult < 1 {
		burnStyle = lipgloss.NewStyle().Foreground(styles.Green)
	} else if burnMult > 0 && burnMult < 2 {
		burnStyle = lipgloss.NewStyle().Foreground(styles.Yellow)
	}
	if burnMult > 0 {
		b.WriteString(burnStyle.Render(fmt.Sprintf("%.1fx", burnMult)))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(styles.Green).Render("Profitable"))
	}
	b.WriteString("\n")

	// Magic Number
	magicNum := fg.CalculateMagicNumber()
	b.WriteString(labelStyle.Render("Magic Number: "))
	magicStyle := lipgloss.NewStyle().Foreground(styles.Red)
	if magicNum > 1 {
		magicStyle = lipgloss.NewStyle().Foreground(styles.Green)
	} else if magicNum > 0.5 {
		magicStyle = lipgloss.NewStyle().Foreground(styles.Yellow)
	}
	if magicNum > 0 {
		b.WriteString(magicStyle.Render(fmt.Sprintf("%.2f", magicNum)))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(styles.Gray).Render("N/A"))
	}
	b.WriteString("\n")

	return panelStyle.Render(b.String())
}

func (s *FounderGameScreen) renderBoardSentimentPanel() string {
	fg := s.gameData.FounderState

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(0, 1).
		Width(35)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(titleStyle.Render("🏛️ BOARD"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)

	// Board sentiment
	sentiment := fg.BoardSentiment
	if sentiment == "" {
		sentiment = "neutral"
	}

	sentimentIcon := "😐"
	sentimentStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	switch sentiment {
	case "happy":
		sentimentIcon = "😊"
		sentimentStyle = lipgloss.NewStyle().Foreground(styles.Green)
	case "pleased":
		sentimentIcon = "🙂"
		sentimentStyle = lipgloss.NewStyle().Foreground(styles.Green)
	case "concerned":
		sentimentIcon = "😟"
		sentimentStyle = lipgloss.NewStyle().Foreground(styles.Yellow)
	case "angry":
		sentimentIcon = "😡"
		sentimentStyle = lipgloss.NewStyle().Foreground(styles.Red)
	}

	b.WriteString(labelStyle.Render("Sentiment: "))
	b.WriteString(sentimentStyle.Render(fmt.Sprintf("%s %s", sentimentIcon, sentiment)))
	b.WriteString("\n")

	// Board pressure
	b.WriteString(labelStyle.Render("Pressure: "))
	pressureStyle := lipgloss.NewStyle().Foreground(styles.Green)
	if fg.BoardPressure > 70 {
		pressureStyle = lipgloss.NewStyle().Foreground(styles.Red)
	} else if fg.BoardPressure > 40 {
		pressureStyle = lipgloss.NewStyle().Foreground(styles.Yellow)
	}
	pressureBar := ""
	filled := fg.BoardPressure / 10
	for i := 0; i < 10; i++ {
		if i < filled {
			pressureBar += "█"
		} else {
			pressureBar += "░"
		}
	}
	b.WriteString(pressureStyle.Render(fmt.Sprintf("%s %d%%", pressureBar, fg.BoardPressure)))
	b.WriteString("\n")

	// Board members count
	activeMembers := 0
	var chairman string
	for _, m := range fg.BoardMembers {
		if m.IsActive {
			activeMembers++
			if m.IsChairman {
				chairman = m.Name
			}
		}
	}
	b.WriteString(labelStyle.Render("Members: "))
	b.WriteString(fmt.Sprintf("%d", activeMembers))
	if chairman != "" {
		b.WriteString(fmt.Sprintf(" (Chair: %s)", truncate(chairman, 12)))
	}
	b.WriteString("\n")

	// Board seats
	b.WriteString(labelStyle.Render("Board Seats: "))
	b.WriteString(fmt.Sprintf("%d", fg.BoardSeats))
	b.WriteString("\n")

	// Warning if board pressure high
	if fg.BoardPressure > 70 {
		warnStyle := lipgloss.NewStyle().Foreground(styles.Red)
		b.WriteString(warnStyle.Render("⚠ Board may force exit!"))
		b.WriteString("\n")
	}

	return panelStyle.Render(b.String())
}

func (s *FounderGameScreen) renderHighlights(highlights []founder.MonthlyHighlight) string {
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(72)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	b.WriteString(titleStyle.Render("📋 MONTHLY HIGHLIGHTS"))
	b.WriteString("\n")

	winStyle := lipgloss.NewStyle().Foreground(styles.Green)
	concernStyle := lipgloss.NewStyle().Foreground(styles.Yellow)

	shown := 0
	for _, h := range highlights {
		if shown >= 4 {
			break
		}
		msg := h.Message
		if len(msg) > 60 {
			msg = msg[:57] + "..."
		}
		if h.Type == "win" {
			b.WriteString(winStyle.Render(fmt.Sprintf("%s %s", h.Icon, msg)))
		} else {
			b.WriteString(concernStyle.Render(fmt.Sprintf("%s %s", h.Icon, msg)))
		}
		b.WriteString("\n")
		shown++
	}

	return lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(panelStyle.Render(b.String()))
}

func (s *FounderGameScreen) renderNews() string {
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(0, 1).
		Width(72)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(titleStyle.Render("📰 NEWS"))
	b.WriteString("\n")

	limit := 8
	if len(s.turnMessages) < limit {
		limit = len(s.turnMessages)
	}
	for i := 0; i < limit; i++ {
		msg := s.turnMessages[i]
		b.WriteString("• " + msg + "\n")
	}

	return lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(panelStyle.Render(b.String()))
}

func (s *FounderGameScreen) renderActions() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⚡ MONTHLY DECISIONS")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	b.WriteString(menuContainer.Render(menuBox.Render(s.actionsMenu.View())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderHiring() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("💼 HIRE TEAM MEMBER")))
	b.WriteString("\n\n")

	// Cash info
	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	b.WriteString(infoStyle.Render(fmt.Sprintf("Cash: $%s | Runway: %d months", formatCompactMoney(fg.Cash), fg.CashRunwayMonths)))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2)

	b.WriteString(menuContainer.Render(menuBox.Render(s.hiringMenu.View())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderHiringMarket() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📍 SELECT MARKET ASSIGNMENT")))
	b.WriteString("\n\n")

	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	b.WriteString(infoStyle.Render(fmt.Sprintf("Hiring: %s", s.selectedRole)))
	b.WriteString("\n\n")

	marketBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(50)

	var markets strings.Builder
	for i, m := range s.marketOptions {
		markets.WriteString(fmt.Sprintf("%d. %s\n", i+1, m))
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(marketBox.Render(markets.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("1-9 select market • esc back"))

	return b.String()
}

func (s *FounderGameScreen) renderFiring() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Red).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⚠️ LET GO TEAM MEMBER")))
	b.WriteString("\n\n")

	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	b.WriteString(infoStyle.Render(fmt.Sprintf("Team Size: %d | Monthly Cost: $%s", fg.Team.TotalEmployees, formatCompactMoney(fg.MonthlyTeamCost))))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Red).
		Padding(1, 2)

	b.WriteString(menuContainer.Render(menuBox.Render(s.firingMenu.View())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderMarketing() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Orange).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📣 MARKETING SPEND")))
	b.WriteString("\n\n")

	// Info
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Orange).
		Padding(1, 2).
		Width(55)

	var info strings.Builder
	info.WriteString(fmt.Sprintf("Current Cash: $%s\n", formatCompactMoney(fg.Cash)))
	info.WriteString(fmt.Sprintf("Base CAC: $%s\n", formatCompactMoney(fg.BaseCAC)))
	info.WriteString(fmt.Sprintf("Effective CAC: $%s\n", formatCompactMoney(fg.CustomerAcquisitionCost)))
	info.WriteString(fmt.Sprintf("Product Maturity: %.0f%% (reduces CAC up to 40%%)\n", fg.ProductMaturity*100))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(infoBox.Render(info.String())))
	b.WriteString("\n\n")

	// Input
	inputLabel := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Width(s.width).Align(lipgloss.Center)
	b.WriteString(inputLabel.Render("ENTER MARKETING BUDGET"))
	b.WriteString("\n")

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(30)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render("$ " + s.marketingInput.View())))
	b.WriteString("\n")

	// Error message
	if s.inputMessage != "" {
		msgStyle := lipgloss.NewStyle().Foreground(styles.Red).Width(s.width).Align(lipgloss.Center)
		b.WriteString(msgStyle.Render(s.inputMessage))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter confirm • esc cancel"))

	return b.String()
}

func (s *FounderGameScreen) renderFunding() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("💰 RAISE FUNDING")))
	b.WriteString("\n\n")

	// Current status
	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	currentEquity := 100.0 - fg.EquityGivenAway - fg.EquityPool
	b.WriteString(infoStyle.Render(fmt.Sprintf("Your Equity: %.1f%% | Rounds Raised: %d", currentEquity, len(fg.FundingRounds))))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2)

	b.WriteString(menuContainer.Render(menuBox.Render(s.fundingMenu.View())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderFundingTerms() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render(fmt.Sprintf("📄 %s TERM SHEETS", strings.ToUpper(s.selectedRoundName)))))
	b.WriteString("\n\n")

	currentEquity := 100.0 - fg.EquityGivenAway - fg.EquityPool
	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	b.WriteString(infoStyle.Render(fmt.Sprintf("Your Current Equity: %.2f%%", currentEquity)))
	if currentEquity >= 50.0 {
		b.WriteString(lipgloss.NewStyle().Foreground(styles.Green).Width(s.width).Align(lipgloss.Center).Render(" ✓ Majority control"))
	}
	b.WriteString("\n\n")

	// Term sheets
	termsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2).
		Width(65)

	var terms strings.Builder
	for i, sheet := range s.fundingTerms {
		newEquity := currentEquity - sheet.Equity
		terms.WriteString(fmt.Sprintf("%d. %s\n", i+1, sheet.Terms))
		terms.WriteString(fmt.Sprintf("   Amount: $%s | Valuation: $%s\n", formatCompactMoney(sheet.Amount), formatCompactMoney(sheet.PostValuation)))
		terms.WriteString(fmt.Sprintf("   Equity: %.1f%% | Your equity after: %.1f%%", sheet.Equity, newEquity))
		if newEquity < 50.0 && currentEquity >= 50.0 {
			terms.WriteString(" ⚠️ LOSES CONTROL")
		}
		terms.WriteString("\n\n")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(termsBox.Render(terms.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("1-4 select term sheet • esc back"))

	return b.String()
}

func (s *FounderGameScreen) renderPartnership() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🤝 STRATEGIC PARTNERSHIPS")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	b.WriteString(menuContainer.Render(menuBox.Render(s.partnershipMenu.View())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderAffiliate() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("💸 LAUNCH AFFILIATE PROGRAM")))
	b.WriteString("\n\n")

	// Info
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2).
		Width(55)

	var info strings.Builder
	info.WriteString("Affiliate programs let partners sell for commission.\n\n")
	info.WriteString(fmt.Sprintf("Setup Cost: $20-50k | Monthly Fees: $5-10k\n"))
	info.WriteString(fmt.Sprintf("Your Cash: $%s\n\n", formatCompactMoney(fg.Cash)))
	info.WriteString("Recommended rates:\n")
	info.WriteString("  • 10-15% for SaaS products\n")
	info.WriteString("  • 15-20% for marketplaces\n")
	info.WriteString("  • 20-30% for high-margin products\n")

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(infoBox.Render(info.String())))
	b.WriteString("\n\n")

	// Input
	inputLabel := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Width(s.width).Align(lipgloss.Center)
	b.WriteString(inputLabel.Render("ENTER COMMISSION RATE (5-30%)"))
	b.WriteString("\n")

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(20)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render(s.affiliateInput.View() + "%")))
	b.WriteString("\n")

	if s.inputMessage != "" {
		msgStyle := lipgloss.NewStyle().Foreground(styles.Red).Width(s.width).Align(lipgloss.Center)
		b.WriteString(msgStyle.Render(s.inputMessage))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter confirm • esc cancel"))

	return b.String()
}

func (s *FounderGameScreen) renderCompetitors() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Red).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⚔️ COMPETITOR MANAGEMENT")))
	b.WriteString("\n\n")

	// Show competitor list and/or action menu
	if s.competitorAction != nil && s.selectedCompetitorIdx >= 0 && s.selectedCompetitorIdx < len(fg.Competitors) {
		comp := fg.Competitors[s.selectedCompetitorIdx]
		infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
		b.WriteString(infoStyle.Render(fmt.Sprintf("Handling: %s (Threat: %s)", comp.Name, comp.Threat)))
		b.WriteString("\n\n")

		menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
		menuBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Red).
			Padding(1, 2)
		b.WriteString(menuContainer.Render(menuBox.Render(s.competitorAction.View())))
	} else {
		// Show competitors list
		compBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Red).
			Padding(1, 2).
			Width(60)

		var comps strings.Builder
		activeNum := 0
		for _, c := range fg.Competitors {
			if !c.Active {
				continue
			}
			activeNum++
			threatColor := styles.Green
			if c.Threat == "high" || c.Threat == "critical" {
				threatColor = styles.Red
			} else if c.Threat == "medium" {
				threatColor = styles.Yellow
			}
			threatStyle := lipgloss.NewStyle().Foreground(threatColor)
			comps.WriteString(fmt.Sprintf("%d. %s\n", activeNum, c.Name))
			comps.WriteString(fmt.Sprintf("   Threat: %s | Share: %.1f%%\n", threatStyle.Render(c.Threat), c.MarketShare*100))
			comps.WriteString(fmt.Sprintf("   Strategy: %s\n\n", c.Strategy))
		}

		if activeNum == 0 {
			comps.WriteString("No active competitors at this time")
		}

		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(compBox.Render(comps.String())))
	}

	b.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	if s.competitorAction != nil {
		b.WriteString(helpStyle.Render("enter select • esc back"))
	} else {
		b.WriteString(helpStyle.Render("1-9 select competitor • esc back"))
	}

	return b.String()
}

func (s *FounderGameScreen) renderExpansion() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🌍 GLOBAL EXPANSION")))
	b.WriteString("\n\n")

	// Show active markets info
	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	b.WriteString(infoStyle.Render(fmt.Sprintf("Markets: %d | Cash: $%s", len(fg.GlobalMarkets)+1, formatCompactMoney(fg.Cash))))
	b.WriteString("\n\n")

	// Show menu for expansion
	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	if s.expansionMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.expansionMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderTeamRoster() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("👥 TEAM ROSTER")))
	b.WriteString("\n\n")

	teamBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2).
		Width(72)

	var team strings.Builder
	team.WriteString(fmt.Sprintf("Total Employees: %d | Monthly Cost: $%s\n\n", fg.Team.TotalEmployees, formatCompactMoney(fg.MonthlyTeamCost)))

	renderEmployeeList := func(employees []founder.Employee, title string) {
		if len(employees) == 0 {
			return
		}
		team.WriteString(fmt.Sprintf("%s: %d\n", title, len(employees)))
		for _, e := range employees {
			// Vesting info
			vestInfo := ""
			if e.VestingMonths > 0 {
				if !e.HasCliff {
					vestInfo = fmt.Sprintf(" [cliff: %dmo]", e.CliffMonths-e.VestedMonths)
				} else {
					vestInfo = fmt.Sprintf(" [vested: %d/%dmo]", e.VestedMonths, e.VestingMonths)
				}
			}

			// Equity info
			eqInfo := ""
			if e.Equity > 0 {
				eqInfo = fmt.Sprintf(" %.2f%%eq", e.Equity)
			}

			// Market assignment
			marketInfo := ""
			if e.AssignedMarket != "" && e.AssignedMarket != "USA" {
				marketInfo = fmt.Sprintf(" [%s]", e.AssignedMarket)
			}

			// Salary
			salary := fmt.Sprintf("$%s/mo", formatCompactMoney(e.MonthlyCost))

			team.WriteString(fmt.Sprintf("  • %s %.1fx %s%s%s%s\n",
				truncate(e.Name, 15), e.Impact, salary, eqInfo, vestInfo, marketInfo))
		}
		team.WriteString("\n")
	}

	renderEmployeeList(fg.Team.Engineers, "Engineers")
	renderEmployeeList(fg.Team.Sales, "Sales")
	renderEmployeeList(fg.Team.CustomerSuccess, "Customer Success")
	renderEmployeeList(fg.Team.Marketing, "Marketing")

	if len(fg.Team.Executives) > 0 {
		team.WriteString(fmt.Sprintf("Executives: %d\n", len(fg.Team.Executives)))
		for _, e := range fg.Team.Executives {
			vestInfo := ""
			if e.VestingMonths > 0 {
				vestInfo = fmt.Sprintf(" [vested: %d/%dmo]", e.VestedMonths, e.VestingMonths)
			}
			eqInfo := ""
			if e.Equity > 0 {
				eqInfo = fmt.Sprintf(" %.2f%%eq", e.Equity)
			}
			team.WriteString(fmt.Sprintf("  • %s (%s) %.1fx $%s/mo%s%s\n",
				e.Name, e.Role, e.Impact, formatCompactMoney(e.MonthlyCost), eqInfo, vestInfo))
		}
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(teamBox.Render(team.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))

	return b.String()
}

func (s *FounderGameScreen) renderCustomers() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🏢 CUSTOMERS")))
	b.WriteString("\n\n")

	// Summary box
	custBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(72)

	var cust strings.Builder
	cust.WriteString(fmt.Sprintf("Active Customers: %d\n", fg.Customers))
	cust.WriteString(fmt.Sprintf("Total Ever: %d | Churned: %d\n", fg.TotalCustomersEver, fg.TotalChurned))
	cust.WriteString(fmt.Sprintf("MRR: $%s | Avg Deal: $%s/mo\n", formatCompactMoney(fg.MRR), formatCompactMoney(fg.AvgDealSize)))
	cust.WriteString(fmt.Sprintf("Deal Range: $%s - $%s\n\n", formatCompactMoney(fg.MinDealSize), formatCompactMoney(fg.MaxDealSize)))

	if fg.AffiliateProgram != nil {
		cust.WriteString(fmt.Sprintf("Direct: %d ($%s MRR) | Affiliate: %d ($%s MRR)\n\n",
			fg.DirectCustomers, formatCompactMoney(fg.DirectMRR),
			fg.AffiliateCustomers, formatCompactMoney(fg.AffiliateMRR)))
	}

	cust.WriteString(fmt.Sprintf("Churn Rate: %.1f%%/mo\n", fg.CustomerChurnRate*100))

	// Customer health
	healthy, atRisk, critical, atRiskMRR, criticalMRR := fg.GetCustomerHealthSegments()
	if healthy > 0 || atRisk > 0 || critical > 0 {
		cust.WriteString(fmt.Sprintf("\nHealth: 🟢 %d healthy | 🟡 %d at risk ($%s) | 🔴 %d critical ($%s)\n",
			healthy, atRisk, formatCompactMoney(atRiskMRR), critical, formatCompactMoney(criticalMRR)))
	}

	// Individual customer table (show up to 10 recent)
	activeCustomers := []founder.Customer{}
	for _, c := range fg.CustomerList {
		if c.IsActive {
			activeCustomers = append(activeCustomers, c)
		}
	}

	if len(activeCustomers) > 0 {
		cust.WriteString("\n── CUSTOMER DEALS ──\n")
		tableHeader := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
		cust.WriteString(tableHeader.Render(fmt.Sprintf("%-6s %-10s %-10s %-8s %-8s\n", "ID", "Source", "MRR", "Term", "Health")))

		shown := 0
		for i := len(activeCustomers) - 1; i >= 0 && shown < 10; i-- {
			c := activeCustomers[i]
			termStr := "∞"
			if c.TermMonths > 0 {
				termStr = fmt.Sprintf("%dmo", c.TermMonths)
			}
			healthStr := "🟢"
			if c.HealthScore < 0.3 {
				healthStr = "🔴"
			} else if c.HealthScore < 0.7 {
				healthStr = "🟡"
			}
			cust.WriteString(fmt.Sprintf("%-6d %-10s $%-9s %-8s %s %.0f%%\n",
				c.ID, truncate(c.Source, 10), formatCompactMoney(c.DealSize), termStr, healthStr, c.HealthScore*100))
			shown++
		}
		if len(activeCustomers) > 10 {
			cust.WriteString(fmt.Sprintf("  ... and %d more customers\n", len(activeCustomers)-10))
		}
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(custBox.Render(cust.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))

	return b.String()
}

func (s *FounderGameScreen) renderFinancials() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📊 FINANCIALS")))
	b.WriteString("\n\n")

	finBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2).
		Width(68)

	var fin strings.Builder
	fin.WriteString(fmt.Sprintf("Cash: $%s\n", formatCompactMoney(fg.Cash)))

	runway := fmt.Sprintf("%d months", fg.CashRunwayMonths)
	if fg.CashRunwayMonths < 0 {
		runway = "∞ (profitable!)"
	}
	fin.WriteString(fmt.Sprintf("Runway: %s\n\n", runway))

	// Revenue breakdown with deduction detail
	sectionStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	fin.WriteString(sectionStyle.Render("── MONTHLY REVENUE ──"))
	fin.WriteString("\n")
	fin.WriteString(fmt.Sprintf("  Gross MRR:           $%s\n", formatCompactMoney(fg.MRR)))

	// Revenue deductions detail
	taxes := int64(float64(fg.MRR) * 0.20)
	processing := int64(float64(fg.MRR) * 0.03)
	overhead := int64(float64(fg.MRR) * 0.05)
	savings := int64(float64(fg.MRR) * 0.05)
	totalDeductions := taxes + processing + overhead + savings
	netMRR := fg.MRR - totalDeductions

	deductStyle := lipgloss.NewStyle().Foreground(styles.Red)
	fin.WriteString(deductStyle.Render(fmt.Sprintf("  Taxes (20%%):         -$%s\n", formatCompactMoney(taxes))))
	fin.WriteString(deductStyle.Render(fmt.Sprintf("  Processing (3%%):     -$%s\n", formatCompactMoney(processing))))
	fin.WriteString(deductStyle.Render(fmt.Sprintf("  Overhead (5%%):       -$%s\n", formatCompactMoney(overhead))))
	fin.WriteString(deductStyle.Render(fmt.Sprintf("  Savings (5%%):        -$%s\n", formatCompactMoney(savings))))

	// Affiliate costs if applicable
	if fg.AffiliateProgram != nil {
		affCost := int64(float64(fg.AffiliateMRR) * fg.AffiliateProgram.Commission)
		fin.WriteString(deductStyle.Render(fmt.Sprintf("  Affiliate Comm:      -$%s\n", formatCompactMoney(affCost))))
		netMRR -= affCost
	}

	// Global market costs
	for _, market := range fg.GlobalMarkets {
		if market.MonthlyCost > 0 {
			fin.WriteString(deductStyle.Render(fmt.Sprintf("  %s ops:         -$%s\n", truncate(market.Region, 10), formatCompactMoney(market.MonthlyCost))))
			netMRR -= market.MonthlyCost
		}
	}

	netStyle := lipgloss.NewStyle().Foreground(styles.Green)
	fin.WriteString(netStyle.Render(fmt.Sprintf("  Net Revenue:         $%s\n", formatCompactMoney(netMRR))))

	// Gross margin
	if fg.MRR > 0 {
		grossMargin := float64(netMRR) / float64(fg.MRR) * 100.0
		fin.WriteString(fmt.Sprintf("  Gross Margin:        %.0f%%\n", grossMargin))
	}
	fin.WriteString("\n")

	// Expenses
	fin.WriteString(sectionStyle.Render("── MONTHLY EXPENSES ──"))
	fin.WriteString("\n")
	fin.WriteString(fmt.Sprintf("  Team Salaries:       $%s\n", formatCompactMoney(fg.MonthlyTeamCost)))
	fin.WriteString(fmt.Sprintf("  Compute/Cloud:       $%s\n", formatCompactMoney(fg.MonthlyComputeCost)))
	fin.WriteString(fmt.Sprintf("  Other Direct Costs:  $%s\n", formatCompactMoney(fg.MonthlyODCCost)))

	// Partnership costs
	partnerCost := int64(0)
	for _, p := range fg.Partnerships {
		if p.Status == "active" {
			partnerCost += p.Cost
		}
	}
	if partnerCost > 0 {
		fin.WriteString(fmt.Sprintf("  Partnerships:        $%s\n", formatCompactMoney(partnerCost)))
	}

	totalExpenses := fg.MonthlyTeamCost + fg.MonthlyComputeCost + fg.MonthlyODCCost + partnerCost
	fin.WriteString("  ─────────────────────────\n")
	fin.WriteString(fmt.Sprintf("  Total Expenses:      $%s\n\n", formatCompactMoney(totalExpenses)))

	// Net income
	netIncome := netMRR - totalExpenses
	if netIncome >= 0 {
		fin.WriteString(netStyle.Render(fmt.Sprintf("NET INCOME: +$%s/mo 🟢\n", formatCompactMoney(netIncome))))
	} else {
		fin.WriteString(deductStyle.Render(fmt.Sprintf("NET BURN: -$%s/mo 🔴\n", formatCompactMoney(-netIncome))))
	}

	// Net margin
	if fg.MRR > 0 {
		netMargin := float64(netIncome) / float64(fg.MRR) * 100.0
		marginColor := styles.Green
		if netMargin < 0 {
			marginColor = styles.Red
		}
		fin.WriteString(lipgloss.NewStyle().Foreground(marginColor).Render(fmt.Sprintf("Net Margin: %.0f%%\n", netMargin)))
	}

	// Funding history
	if len(fg.FundingRounds) > 0 {
		fin.WriteString("\n")
		fin.WriteString(sectionStyle.Render("── FUNDING HISTORY ──"))
		fin.WriteString("\n")
		totalRaised := int64(0)
		for _, r := range fg.FundingRounds {
			fin.WriteString(fmt.Sprintf("  %s: $%s (%.1f%% equity, %s terms)\n",
				r.RoundName, formatCompactMoney(r.Amount), r.EquityGiven, r.Terms))
			totalRaised += r.Amount
		}
		fin.WriteString(fmt.Sprintf("  Total Raised: $%s\n", formatCompactMoney(totalRaised)))
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(finBox.Render(fin.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))

	return b.String()
}

func (s *FounderGameScreen) renderCapTable() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📋 CAP TABLE")))
	b.WriteString("\n\n")

	capBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2).
		Width(70)

	var content strings.Builder

	labelStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	valStyle := lipgloss.NewStyle().Foreground(styles.White)
	dimStyle := lipgloss.NewStyle().Foreground(styles.Gray)

	// Summary row
	founderEquity := 100.0 - fg.EquityGivenAway - fg.EquityPool
	availablePool := fg.EquityPool - fg.EquityAllocated
	if availablePool < 0 {
		availablePool = 0
	}
	content.WriteString(labelStyle.Render("Total Shares: "))
	content.WriteString(valStyle.Render("100%"))
	content.WriteString("   ")
	content.WriteString(labelStyle.Render("Funding Rounds: "))
	content.WriteString(valStyle.Render(fmt.Sprintf("%d", len(fg.FundingRounds))))
	content.WriteString("\n\n")

	// Column header
	colHeader := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	content.WriteString(colHeader.Render(fmt.Sprintf("  %-32s %8s  %s", "HOLDER", "EQUITY", "DETAILS")))
	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  " + strings.Repeat("─", 64)))
	content.WriteString("\n")

	// 1. Founder
	founderColor := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	controlNote := ""
	if founderEquity >= 50.0 {
		controlNote = "majority control ✓"
	} else {
		controlNote = "minority ⚠️"
	}
	content.WriteString(founderColor.Render(fmt.Sprintf("  %-32s %7.1f%%  %s", "👤 "+fg.FounderName+" (Founder)", founderEquity, controlNote)))
	content.WriteString("\n")

	// 2. Investors (from funding rounds)
	totalInvestorEquity := 0.0
	investorStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	for _, round := range fg.FundingRounds {
		roundEquity := round.EquityGiven
		totalInvestorEquity += roundEquity
		investorLabel := round.RoundName
		if len(round.Investors) > 0 {
			investorLabel = round.Investors[0]
		}
		content.WriteString(investorStyle.Render(fmt.Sprintf("  %-32s %7.1f%%  %s @ $%s val",
			"💰 "+investorLabel,
			roundEquity,
			round.RoundName,
			formatCompactMoney(round.Valuation))))
		content.WriteString("\n")
	}

	// 3. Cap table entries (executives, employees, advisors)
	execStyle := lipgloss.NewStyle().Foreground(styles.Magenta)
	empStyle := lipgloss.NewStyle().Foreground(styles.White)
	advisorStyle := lipgloss.NewStyle().Foreground(styles.Cyan)

	for _, entry := range fg.CapTable {
		var style lipgloss.Style
		var icon string
		var detail string
		switch entry.Type {
		case "executive":
			style = execStyle
			icon = "⭐"
			detail = fmt.Sprintf("exec, month %d", entry.MonthGranted)
		case "employee":
			style = empStyle
			icon = "👷"
			detail = fmt.Sprintf("employee, month %d", entry.MonthGranted)
		case "advisor":
			style = advisorStyle
			icon = "🧠"
			detail = fmt.Sprintf("advisor, month %d", entry.MonthGranted)
		default:
			style = dimStyle
			icon = "📄"
			detail = entry.Type
		}
		content.WriteString(style.Render(fmt.Sprintf("  %-32s %7.2f%%  %s",
			icon+" "+entry.Name, entry.Equity, detail)))
		content.WriteString("\n")
	}

	// 4. Unallocated pool
	if availablePool > 0 {
		content.WriteString(dimStyle.Render(fmt.Sprintf("  %-32s %7.1f%%  reserved for future hires",
			"🏦 Unallocated Pool", availablePool)))
		content.WriteString("\n")
	}

	// Separator and totals
	content.WriteString(dimStyle.Render("  " + strings.Repeat("─", 64)))
	content.WriteString("\n")

	totalStyle := lipgloss.NewStyle().Foreground(styles.White).Bold(true)
	content.WriteString(totalStyle.Render(fmt.Sprintf("  %-32s %7.1f%%", "TOTAL", 100.0)))
	content.WriteString("\n\n")

	// Summary stats
	summaryLabel := lipgloss.NewStyle().Foreground(styles.Yellow)
	content.WriteString(summaryLabel.Render("  Founder Equity: "))
	if founderEquity >= 50 {
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Green).Render(fmt.Sprintf("%.1f%%", founderEquity)))
	} else {
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Red).Render(fmt.Sprintf("%.1f%%", founderEquity)))
	}
	content.WriteString(summaryLabel.Render("  |  Investor Equity: "))
	content.WriteString(valStyle.Render(fmt.Sprintf("%.1f%%", totalInvestorEquity)))
	content.WriteString(summaryLabel.Render("  |  Employee Pool: "))
	content.WriteString(valStyle.Render(fmt.Sprintf("%.1f%% (%.1f%% used)", fg.EquityPool, fg.EquityAllocated)))
	content.WriteString("\n")

	// At-exit note
	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  ℹ At exit, unallocated pool ("+fmt.Sprintf("%.1f%%", availablePool)+") cancels — founder effective equity: "+fmt.Sprintf("%.1f%%", 100.0-fg.EquityAllocated-fg.EquityGivenAway)+"%"))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(capBox.Render(content.String())))
	b.WriteString("\n\n")

	capHelpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(capHelpStyle.Render("esc back"))

	return b.String()
}

func (s *FounderGameScreen) renderExit() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🚪 EXIT OPTIONS")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2)

	b.WriteString(menuContainer.Render(menuBox.Render(s.exitMenu.View())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderConfirmExit() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Red).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⚠️ CONFIRM EXIT")))
	b.WriteString("\n\n")

	exits := fg.GetAvailableExits()
	var selectedExit *founder.ExitOption
	for _, e := range exits {
		if e.Type == s.selectedExitType {
			selectedExit = &e
			break
		}
	}

	if selectedExit == nil {
		return b.String()
	}

	confirmBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Red).
		Padding(1, 2).
		Width(60)

	// At exit, unallocated equity pool cancels — only allocated equity counts
	founderEquity := 100.0 - fg.EquityGivenAway - fg.EquityAllocated
	payout := int64(float64(selectedExit.Valuation) * founderEquity / 100.0)

	var confirm strings.Builder
	confirm.WriteString(fmt.Sprintf("Exit Type: %s\n", strings.ToUpper(s.selectedExitType)))
	confirm.WriteString(fmt.Sprintf("Valuation: $%s\n", formatCompactMoney(selectedExit.Valuation)))
	confirm.WriteString(fmt.Sprintf("Your Equity: %.1f%%\n", founderEquity))
	confirm.WriteString(fmt.Sprintf("Your Payout: $%s\n\n", formatCompactMoney(payout)))
	confirm.WriteString(selectedExit.Description)
	confirm.WriteString("\n\nThis action is PERMANENT. Are you sure?")

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(confirmBox.Render(confirm.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("y confirm • n cancel • esc back"))

	return b.String()
}

// ============================================================================
// PHASE 4: ADVANCED RENDER FUNCTIONS
// ============================================================================

func (s *FounderGameScreen) renderRoadmap() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Orange).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🔨 PRODUCT ROADMAP")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Orange).
		Padding(1, 2)

	if s.roadmapMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.roadmapMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderRoadmapStart() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Orange).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🚀 START NEW FEATURE")))
	b.WriteString("\n\n")

	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	allocated := fg.GetAllocatedEngineers()
	b.WriteString(infoStyle.Render(fmt.Sprintf("Engineers: %d total, %d allocated", len(fg.Team.Engineers), allocated)))
	b.WriteString("\n\n")

	featureBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Orange).
		Padding(1, 2).
		Width(65)

	var features strings.Builder
	for i, f := range s.roadmapFeatures {
		features.WriteString(fmt.Sprintf("%d. %s\n", i+1, f.Name))
		features.WriteString(fmt.Sprintf("   Category: %s | Cost: $%s | Engineer-months: %d\n\n", f.Category, formatCompactMoney(f.Cost), f.EngineerMonths))
	}

	if len(s.roadmapFeatures) == 0 {
		features.WriteString("No features available to start")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(featureBox.Render(features.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("1-9 select feature • esc back"))

	return b.String()
}

func (s *FounderGameScreen) renderSegments() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🎯 CUSTOMER SEGMENTS")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2)

	if s.segmentsMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.segmentsMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderPricing() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("💲 PRICING STRATEGY")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2)

	if s.pricingMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.pricingMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderAcquisitions() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🏢 ACQUISITIONS")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	if s.acquisitionsMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.acquisitionsMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderPlatform() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🌐 PLATFORM STRATEGY")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	if s.platformMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.platformMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderSecurity() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Red).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🔒 SECURITY & COMPLIANCE")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Red).
		Padding(1, 2)

	if s.securityMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.securityMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderPRCrisis() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Yellow).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📰 PR CRISIS MANAGEMENT")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2)

	if s.prCrisisMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.prCrisisMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderEconomy() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Orange).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📉 ECONOMIC STRATEGY")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Orange).
		Padding(1, 2)

	if s.economyMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.economyMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderSuccession() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("👤 SUCCESSION PLANNING")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2)

	if s.successionMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.successionMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderSalesPipeline() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📈 SALES PIPELINE")))
	b.WriteString("\n\n")

	pipeBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2).
		Width(65)

	var pipe strings.Builder

	if fg.SalesPipeline != nil {
		metrics := fg.GetPipelineMetrics()
		if totalDeals, ok := metrics["total_deals"].(int); ok {
			pipe.WriteString(fmt.Sprintf("Active Deals: %d\n", totalDeals))
		}
		if totalValue, ok := metrics["total_value"].(int64); ok {
			pipe.WriteString(fmt.Sprintf("Total Pipeline Value: $%s\n", formatCompactMoney(totalValue)))
		}
		if avgDeal, ok := metrics["avg_deal_size"].(int64); ok {
			pipe.WriteString(fmt.Sprintf("Avg Deal Size: $%s\n", formatCompactMoney(avgDeal)))
		}
		if closeRate, ok := metrics["close_rate"].(float64); ok {
			pipe.WriteString(fmt.Sprintf("Close Rate: %.0f%%\n", closeRate*100))
		}

		pipe.WriteString("\nDeals by Stage:\n")
		for _, deal := range fg.SalesPipeline.ActiveDeals {
			pipe.WriteString(fmt.Sprintf("  • %s: $%s (%s)\n", deal.CompanyName, formatCompactMoney(deal.DealSize), deal.Stage))
		}
	} else {
		pipe.WriteString("No active sales pipeline")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(pipeBox.Render(pipe.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))

	return b.String()
}

func (s *FounderGameScreen) renderPivot() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Orange).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🔄 EXECUTE PIVOT")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Orange).
		Padding(1, 2)

	if s.pivotMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.pivotMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderBoard() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("👔 BOARD & EQUITY")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	if s.boardMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.boardMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderBuyback() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📈 EQUITY BUYBACK")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2)

	if s.buybackMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.buybackMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderBuybackConfirm() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render(fmt.Sprintf("📈 BUYBACK FROM %s", s.selectedBuybackRound))))
	b.WriteString("\n\n")

	infoStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
	b.WriteString(infoStyle.Render(fmt.Sprintf("Cash: $%s", formatCompactMoney(fg.Cash))))
	b.WriteString("\n\n")

	inputLabel := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Width(s.width).Align(lipgloss.Center)
	b.WriteString(inputLabel.Render("ENTER EQUITY % TO BUY BACK"))
	b.WriteString("\n")

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1).
		Width(20)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render(s.buybackInput.View() + "%")))
	b.WriteString("\n")

	if s.inputMessage != "" {
		msgStyle := lipgloss.NewStyle().Foreground(styles.Red).Width(s.width).Align(lipgloss.Center)
		b.WriteString(msgStyle.Render(s.inputMessage))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter confirm • esc cancel"))

	return b.String()
}

// ============================================================================
// NEW FEATURE RENDER FUNCTIONS
// ============================================================================

func (s *FounderGameScreen) renderStrategicOpportunity() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Yellow).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("STRATEGIC OPPORTUNITY")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2)

	if s.strategicOpportunityMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.strategicOpportunityMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderContentMarketing() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("CONTENT MARKETING")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	if s.contentMarketingMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.contentMarketingMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderCSPlaybooks() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("CS PLAYBOOKS")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 2)

	if s.csPlaybooksMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.csPlaybooksMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderCompetitiveIntel() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("COMPETITIVE INTELLIGENCE")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2)

	if s.competitiveIntelMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.competitiveIntelMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderReferralProgram() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("REFERRAL PROGRAM")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2)

	if s.referralProgramMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.referralProgramMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderTechDebt() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Yellow).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("TECHNICAL DEBT")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2)

	if s.techDebtMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.techDebtMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

// ============================================================================
// Board Table View
// ============================================================================

func (s *FounderGameScreen) renderBoardTable() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("👔 BOARD COMPOSITION")))
	b.WriteString("\n\n")

	// Summary line
	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(65)

	var content strings.Builder

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	valStyle := lipgloss.NewStyle().Foreground(styles.White)

	content.WriteString(labelStyle.Render("Board Seats: "))
	content.WriteString(valStyle.Render(fmt.Sprintf("%d", fg.BoardSeats)))
	content.WriteString("    ")
	availablePool := fg.EquityPool - fg.EquityAllocated
	if availablePool < 0 {
		availablePool = 0
	}
	content.WriteString(labelStyle.Render("Equity Pool: "))
	content.WriteString(valStyle.Render(fmt.Sprintf("%.1f%% (%.1f%% available)", fg.EquityPool, availablePool)))
	content.WriteString("    ")
	content.WriteString(labelStyle.Render("Your Equity: "))
	content.WriteString(valStyle.Render(fmt.Sprintf("%.1f%%", 100.0-fg.EquityPool-fg.EquityGivenAway)))
	content.WriteString("\n\n")

	// Sentiment
	sentiment := fg.BoardSentiment
	if sentiment == "" {
		sentiment = "neutral"
	}
	sentimentIcon := "😐"
	sentimentColor := styles.Yellow
	switch sentiment {
	case "happy":
		sentimentIcon = "😊"
		sentimentColor = styles.Green
	case "pleased":
		sentimentIcon = "🙂"
		sentimentColor = styles.Green
	case "concerned":
		sentimentIcon = "😟"
		sentimentColor = styles.Yellow
	case "angry":
		sentimentIcon = "😡"
		sentimentColor = styles.Red
	}
	content.WriteString(labelStyle.Render("Sentiment: "))
	content.WriteString(lipgloss.NewStyle().Foreground(sentimentColor).Render(fmt.Sprintf("%s %s", sentimentIcon, sentiment)))
	content.WriteString("    ")
	content.WriteString(labelStyle.Render("Pressure: "))
	pressureColor := styles.Green
	if fg.BoardPressure > 70 {
		pressureColor = styles.Red
	} else if fg.BoardPressure > 40 {
		pressureColor = styles.Yellow
	}
	content.WriteString(lipgloss.NewStyle().Foreground(pressureColor).Render(fmt.Sprintf("%d%%", fg.BoardPressure)))
	content.WriteString("\n")

	// Divider
	divider := lipgloss.NewStyle().Foreground(styles.Gray)
	content.WriteString(divider.Render(strings.Repeat("─", 55)))
	content.WriteString("\n\n")

	// Chairman
	chairman := fg.GetChairman()
	if chairman != nil {
		chairStyle := lipgloss.NewStyle().Foreground(styles.Gold).Bold(true)
		content.WriteString(chairStyle.Render("👑 CHAIRMAN"))
		content.WriteString("\n")
		nameStyle := lipgloss.NewStyle().Foreground(styles.White).Bold(true)
		content.WriteString(nameStyle.Render(fmt.Sprintf("   %s", chairman.Name)))
		content.WriteString("\n")
		detailStyle := lipgloss.NewStyle().Foreground(styles.Gray)
		content.WriteString(detailStyle.Render(fmt.Sprintf("   Expertise: %s │ Impact: 2x │ Equity: %.2f%%",
			chairman.Expertise, chairman.EquityCost)))
		content.WriteString("\n")
		content.WriteString(detailStyle.Render(fmt.Sprintf("   Contribution: %.0f%% │ Added: Month %d",
			chairman.ContributionScore*100, chairman.MonthAdded)))
		content.WriteString("\n\n")
	}

	// Advisors
	advisorCount := 0
	for _, m := range fg.BoardMembers {
		if !m.IsActive || m.IsChairman || m.Type != "advisor" {
			continue
		}
		if advisorCount == 0 {
			sectionStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
			content.WriteString(sectionStyle.Render("🧠 ADVISORS"))
			content.WriteString("\n")
		}
		nameStyle := lipgloss.NewStyle().Foreground(styles.White)
		content.WriteString(nameStyle.Render(fmt.Sprintf("   %s", m.Name)))
		detailStyle := lipgloss.NewStyle().Foreground(styles.Gray)
		content.WriteString(detailStyle.Render(fmt.Sprintf(" (%s) %.2f%% eq, Score: %.0f%%",
			m.Expertise, m.EquityCost, m.ContributionScore*100)))
		content.WriteString("\n")
		advisorCount++
	}
	if advisorCount > 0 {
		content.WriteString("\n")
	}

	// Investor directors
	investorCount := 0
	for _, m := range fg.BoardMembers {
		if !m.IsActive || m.Type != "investor" {
			continue
		}
		if investorCount == 0 {
			sectionStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
			content.WriteString(sectionStyle.Render("💼 INVESTOR DIRECTORS"))
			content.WriteString("\n")
		}
		nameStyle := lipgloss.NewStyle().Foreground(styles.White)
		content.WriteString(nameStyle.Render(fmt.Sprintf("   %s", m.Name)))
		detailStyle := lipgloss.NewStyle().Foreground(styles.Gray)
		content.WriteString(detailStyle.Render(fmt.Sprintf(" (%s) %.2f%% equity", m.Expertise, m.EquityCost)))
		content.WriteString("\n")
		investorCount++
	}

	// Empty state
	if chairman == nil && advisorCount == 0 && investorCount == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(styles.Gray).Italic(true)
		content.WriteString(emptyStyle.Render("   No board members yet. Add advisors to get guidance!"))
		content.WriteString("\n")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(summaryBox.Render(content.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc/enter back to board menu"))

	return b.String()
}

// ============================================================================
// Expand Equity Pool
// ============================================================================

func (s *FounderGameScreen) renderExpandPool() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📊 EXPAND EQUITY POOL")))
	b.WriteString("\n\n")

	contentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(50)

	var content strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	content.WriteString(titleStyle.Render("Enter expansion percentage (1-10%):"))
	content.WriteString("\n\n")

	if s.inputMessage != "" {
		infoStyle := lipgloss.NewStyle().Foreground(styles.Gray)
		content.WriteString(infoStyle.Render(s.inputMessage))
		content.WriteString("\n\n")
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1)
	content.WriteString(inputStyle.Render(s.equityPoolInput.View()))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(contentBox.Render(content.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter confirm"))

	return b.String()
}

// ============================================================================
// Advisor Management
// ============================================================================

var svAdvisorNames = map[string][]string{
	"sales":       {"Marc Andreessen", "Ben Horowitz", "Jason Lemkin", "Tomas Tunguz", "Sam Altman"},
	"product":     {"Marissa Mayer", "Julie Zhuo", "Ken Norton", "Shreyas Doshi", "Lenny Rachitsky"},
	"fundraising": {"Reid Hoffman", "Peter Thiel", "Vinod Khosla", "Mary Meeker", "Bill Gurley"},
	"operations":  {"Sheryl Sandberg", "Keith Rabois", "Claire Hughes Johnson", "Frank Slootman", "Elad Gil"},
	"strategy":    {"Eric Schmidt", "Jeff Bezos", "Patrick Collison", "Stewart Butterfield", "Drew Houston"},
}

func (s *FounderGameScreen) rebuildAdvisorExpertiseMenu() {
	items := []components.MenuItem{
		{ID: "sales", Title: "Sales & Revenue", Description: "Close more deals, improve conversion", Icon: "💰"},
		{ID: "product", Title: "Product & Engineering", Description: "Build faster, improve maturity", Icon: "🔧"},
		{ID: "fundraising", Title: "Fundraising & Finance", Description: "Better terms, investor intros", Icon: "📈"},
		{ID: "operations", Title: "Operations & Scaling", Description: "Reduce costs, improve efficiency", Icon: "⚙️"},
		{ID: "strategy", Title: "Strategy & Growth", Description: "Market positioning, expansion", Icon: "🎯"},
		{ID: "cancel", Title: "Cancel", Icon: "←"},
	}

	s.advisorExpertiseMenu = components.NewMenu("SELECT ADVISOR EXPERTISE", items)
	s.advisorExpertiseMenu.SetSize(55, 15)
	s.advisorExpertiseMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleAdvisorExpertiseSelection(id string) (ScreenModel, tea.Cmd) {
	if id == "cancel" {
		s.rebuildBoardMenu()
		s.view = FounderViewBoard
		return s, nil
	}

	fg := s.gameData.FounderState
	s.selectedExpertise = id

	// Pick a name from the SV names pool
	names, ok := svAdvisorNames[id]
	if !ok {
		names = []string{"Alex Kim", "Jordan Lee", "Morgan Chen"}
	}
	nameIdx := len(fg.BoardMembers) % len(names)
	s.pendingAdvisorName = names[nameIdx]

	// Calculate costs
	equityCost := 0.25 + float64(len(fg.BoardMembers))*0.1
	if equityCost > 1.0 {
		equityCost = 1.0
	}
	s.pendingAdvisorCost = equityCost

	// Setup fee: $10k-50k based on expertise
	setupFees := map[string]int64{
		"sales": 25000, "product": 20000, "fundraising": 50000, "operations": 30000, "strategy": 40000,
	}
	s.pendingAdvisorSetup = setupFees[id]

	// Build confirmation menu
	items := []components.MenuItem{
		{
			ID:    "confirm",
			Title: fmt.Sprintf("Hire %s ($%sK setup + %.2f%% equity)", s.pendingAdvisorName, formatCompactMoney(s.pendingAdvisorSetup), equityCost),
			Icon:  "✓",
		},
		{ID: "cancel", Title: "Cancel", Icon: "←"},
	}
	s.advisorConfirmMenu = components.NewMenu("CONFIRM ADVISOR", items)
	s.advisorConfirmMenu.SetSize(60, 8)
	s.advisorConfirmMenu.SetHideHelp(true)
	s.view = FounderViewAdvisorConfirm
	return s, nil
}

func (s *FounderGameScreen) handleAdvisorConfirmSelection(id string) (ScreenModel, tea.Cmd) {
	if id == "cancel" {
		s.rebuildAdvisorExpertiseMenu()
		s.view = FounderViewAdvisorExpertise
		return s, nil
	}

	fg := s.gameData.FounderState

	// Deduct setup fee
	fg.Cash -= s.pendingAdvisorSetup

	// Create advisor
	advisor := founder.BoardMember{
		Name:       s.pendingAdvisorName,
		Type:       "advisor",
		Expertise:  s.selectedExpertise,
		MonthAdded: fg.Turn,
		EquityCost: s.pendingAdvisorCost,
		IsActive:   true,
	}
	fg.BoardMembers = append(fg.BoardMembers, advisor)
	fg.EquityAllocated += s.pendingAdvisorCost // Allocate from pool, don't shrink pool itself

	// Add to cap table
	fg.CapTable = append(fg.CapTable, founder.CapTableEntry{
		Name:         advisor.Name,
		Type:         "advisor",
		Equity:       s.pendingAdvisorCost,
		MonthGranted: fg.Turn,
	})

	s.turnMessages = []string{
		fmt.Sprintf("✓ Hired advisor: %s (%s)", advisor.Name, advisor.Expertise),
		fmt.Sprintf("   Setup fee: $%s", formatCompactMoney(s.pendingAdvisorSetup)),
		fmt.Sprintf("   Equity cost: %.2f%%", advisor.EquityCost),
		fmt.Sprintf("   Cash remaining: $%s", formatCompactMoney(fg.Cash)),
	}
	s.rebuildBoardMenu()
	s.view = FounderViewBoard
	return s, nil
}

func (s *FounderGameScreen) rebuildRemoveAdvisorMenu() {
	fg := s.gameData.FounderState

	items := []components.MenuItem{}
	for _, m := range fg.BoardMembers {
		if m.IsActive && m.Type == "advisor" {
			// Calculate buyback cost based on valuation
			arr := fg.MRR * 12
			valuation := int64(float64(arr) * 10.0)
			buybackCost := int64(float64(valuation) * m.EquityCost / 100.0)

			items = append(items, components.MenuItem{
				ID:          "buyback_" + m.Name,
				Title:       fmt.Sprintf("Remove %s (buy back equity: $%s)", m.Name, formatCompactMoney(buybackCost)),
				Description: fmt.Sprintf("%s advisor, %.2f%% equity", m.Expertise, m.EquityCost),
				Icon:        "💰",
			})
			items = append(items, components.MenuItem{
				ID:          "nobuyback_" + m.Name,
				Title:       fmt.Sprintf("Remove %s (no buyback)", m.Name),
				Description: "They keep equity, possible negative PR",
				Icon:        "❌",
			})
		}
	}

	items = append(items, components.MenuItem{ID: "cancel", Title: "Cancel", Icon: "←"})

	s.removeAdvisorMenu = components.NewMenu("REMOVE ADVISOR", items)
	s.removeAdvisorMenu.SetSize(60, 15)
	s.removeAdvisorMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleRemoveAdvisorSelection(id string) (ScreenModel, tea.Cmd) {
	if id == "cancel" {
		s.rebuildBoardMenu()
		s.view = FounderViewBoard
		return s, nil
	}

	fg := s.gameData.FounderState

	if strings.HasPrefix(id, "buyback_") {
		name := strings.TrimPrefix(id, "buyback_")
		// Calculate and process buyback
		arr := fg.MRR * 12
		valuation := int64(float64(arr) * 10.0)

		for i := range fg.BoardMembers {
			if fg.BoardMembers[i].Name == name && fg.BoardMembers[i].IsActive {
				buybackCost := int64(float64(valuation) * fg.BoardMembers[i].EquityCost / 100.0)
				if fg.Cash >= buybackCost {
					fg.Cash -= buybackCost
					fg.EquityAllocated -= fg.BoardMembers[i].EquityCost // Return to available pool
					if fg.EquityAllocated < 0 {
						fg.EquityAllocated = 0
					}
					// Remove from cap table
					for j := len(fg.CapTable) - 1; j >= 0; j-- {
						if fg.CapTable[j].Name == name && fg.CapTable[j].Type == "advisor" {
							fg.CapTable = append(fg.CapTable[:j], fg.CapTable[j+1:]...)
							break
						}
					}
					fg.BoardMembers[i].IsActive = false
					s.turnMessages = []string{
						fmt.Sprintf("✓ Removed advisor: %s", name),
						fmt.Sprintf("   Bought back %.2f%% equity for $%s", fg.BoardMembers[i].EquityCost, formatCompactMoney(buybackCost)),
						"   Equity returned to pool",
					}
				} else {
					s.turnMessages = []string{
						fmt.Sprintf("❌ Not enough cash for buyback ($%s needed)", formatCompactMoney(buybackCost)),
					}
				}
				break
			}
		}
	} else if strings.HasPrefix(id, "nobuyback_") {
		name := strings.TrimPrefix(id, "nobuyback_")
		fg.RemoveAdvisor(name, false)
		s.turnMessages = []string{
			fmt.Sprintf("✓ Removed advisor: %s (no buyback)", name),
			"   Advisor retains equity",
			"   ⚠️  Possible negative PR impact",
		}
	}

	s.rebuildBoardMenu()
	s.view = FounderViewBoard
	return s, nil
}

func (s *FounderGameScreen) renderAdvisorExpertise() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🧠 ADD ADVISOR")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	if s.advisorExpertiseMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.advisorExpertiseMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderAdvisorConfirm() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("CONFIRM ADVISOR HIRE")))
	b.WriteString("\n\n")

	// Details box
	detailBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(55)

	var details strings.Builder
	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	details.WriteString(labelStyle.Render("Advisor: "))
	details.WriteString(fmt.Sprintf("%s\n", s.pendingAdvisorName))
	details.WriteString(labelStyle.Render("Expertise: "))
	details.WriteString(fmt.Sprintf("%s\n", s.selectedExpertise))
	details.WriteString(labelStyle.Render("Setup Fee: "))
	details.WriteString(fmt.Sprintf("$%s\n", formatCompactMoney(s.pendingAdvisorSetup)))
	details.WriteString(labelStyle.Render("Equity Cost: "))
	details.WriteString(fmt.Sprintf("%.2f%%\n\n", s.pendingAdvisorCost))

	availPool := fg.EquityPool - fg.EquityAllocated
	if availPool < 0 {
		availPool = 0
	}
	details.WriteString(labelStyle.Render("After hire:\n"))
	details.WriteString(fmt.Sprintf("  Cash: $%s → $%s\n",
		formatCompactMoney(fg.Cash), formatCompactMoney(fg.Cash-s.pendingAdvisorSetup)))
	details.WriteString(fmt.Sprintf("  Pool Available: %.1f%% → %.1f%%\n",
		availPool, availPool-s.pendingAdvisorCost))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(detailBox.Render(details.String())))
	b.WriteString("\n\n")

	if s.advisorConfirmMenu != nil {
		menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
		b.WriteString(menuContainer.Render(s.advisorConfirmMenu.View()))
	}
	b.WriteString("\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderRemoveAdvisor() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Red).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("❌ REMOVE ADVISOR")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Red).
		Padding(1, 2)

	if s.removeAdvisorMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.removeAdvisorMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

// ============================================================================
// Acquisition Offer
// ============================================================================

func (s *FounderGameScreen) rebuildAcquisitionMenu() {
	fg := s.gameData.FounderState
	offer := s.pendingAcquisition
	if offer == nil {
		return
	}

	founderEquity := 100.0 - fg.EquityGivenAway - fg.EquityPool
	forcedAcceptance := founderEquity < 50.0

	items := []components.MenuItem{}
	if forcedAcceptance {
		items = append(items, components.MenuItem{
			ID:          "forced",
			Title:       "Board Forces Acquisition",
			Description: fmt.Sprintf("You have %.1f%% equity - board has control", founderEquity),
			Icon:        "⚠️",
		})
	} else {
		items = append(items, components.MenuItem{
			ID:          "accept",
			Title:       "Accept Offer",
			Description: fmt.Sprintf("Sell company for $%s", formatCompactMoney(offer.OfferAmount)),
			Icon:        "✓",
		})
		items = append(items, components.MenuItem{
			ID:          "decline",
			Title:       "Decline Offer",
			Description: "Continue building your company",
			Icon:        "✗",
		})
	}

	s.acquisitionMenu = components.NewMenu("YOUR DECISION", items)
	s.acquisitionMenu.SetSize(60, 10)
	s.acquisitionMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleAcquisitionSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState
	offer := s.pendingAcquisition
	if offer == nil {
		s.view = FounderViewMain
		return s, nil
	}

	// At exit, unallocated equity pool cancels — only allocated equity counts
	founderEquity := 100.0 - fg.EquityGivenAway - fg.EquityAllocated
	founderPayout := int64(float64(offer.OfferAmount) * founderEquity / 100.0)

	switch id {
	case "accept", "forced":
		fg.Cash = founderPayout
		fg.HasExited = true
		fg.ExitType = "acquisition"
		fg.ExitValuation = offer.OfferAmount
		fg.ExitMonth = fg.Turn
		fg.Turn = fg.MaxTurns + 1 // End game

		if offer.IsCompetitor {
			s.turnMessages = append(s.turnMessages, fmt.Sprintf("⚠️ %s acquired your company for $%s!", offer.Acquirer, formatCompactMoney(offer.OfferAmount)))
		} else {
			s.turnMessages = append(s.turnMessages, fmt.Sprintf("🎉 Acquisition complete! Sold to %s for $%s", offer.Acquirer, formatCompactMoney(offer.OfferAmount)))
		}
		s.turnMessages = append(s.turnMessages, fmt.Sprintf("💰 Your payout: $%s (%.1f%% equity)", formatCompactMoney(founderPayout), founderEquity))

	case "decline":
		s.turnMessages = append(s.turnMessages, fmt.Sprintf("✓ Declined acquisition offer from %s", offer.Acquirer))
	}

	s.pendingAcquisition = nil
	s.view = FounderViewMain
	return s, nil
}

func (s *FounderGameScreen) renderAcquisitionOffer() string {
	fg := s.gameData.FounderState
	offer := s.pendingAcquisition
	if offer == nil {
		return "No offer"
	}

	var b strings.Builder

	// Header
	var headerStyle lipgloss.Style
	if offer.IsCompetitor {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Black).
			Background(styles.Red).
			Bold(true).
			Width(70).
			Align(lipgloss.Center)
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⚠️ COMPETITOR ACQUISITION OFFER!")))
	} else {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Black).
			Background(styles.Green).
			Bold(true).
			Width(70).
			Align(lipgloss.Center)
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🎉 ACQUISITION OFFER!")))
	}
	b.WriteString("\n\n")

	// Offer details box
	offerBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 3).
		Width(65)

	var details strings.Builder

	acquirerStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	details.WriteString(acquirerStyle.Render(fmt.Sprintf("%s wants to acquire %s!", offer.Acquirer, fg.CompanyName)))
	details.WriteString("\n\n")

	amountStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	details.WriteString("Offer Amount: ")
	details.WriteString(amountStyle.Render(fmt.Sprintf("$%s", formatCompactMoney(offer.OfferAmount))))
	details.WriteString("\n")
	details.WriteString(fmt.Sprintf("Due Diligence: %s\n", offer.DueDiligence))
	details.WriteString(fmt.Sprintf("Terms Quality: %s\n", offer.TermsQuality))
	details.WriteString("\n")

	// Cap table payout breakdown
	headerRow := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	details.WriteString(headerRow.Render("── PAYOUT BREAKDOWN ──"))
	details.WriteString("\n\n")

	// At exit, unallocated equity pool cancels — only allocated equity counts
	founderEquity := 100.0 - fg.EquityGivenAway - fg.EquityAllocated
	founderPayout := int64(float64(offer.OfferAmount) * founderEquity / 100.0)

	founderLine := lipgloss.NewStyle().Foreground(styles.Green)
	details.WriteString(founderLine.Render(fmt.Sprintf("%-30s %6.1f%%  $%s", "You (Founder)", founderEquity, formatCompactMoney(founderPayout))))
	details.WriteString("\n")

	for _, entry := range fg.CapTable {
		payout := int64(float64(offer.OfferAmount) * entry.Equity / 100.0)
		label := entry.Name
		switch entry.Type {
		case "executive":
			label += " (Exec)"
		case "employee":
			label += " (Emp)"
		case "advisor":
			label += " (Adv)"
		}
		details.WriteString(fmt.Sprintf("%-30s %6.1f%%  $%s\n",
			truncate(label, 30), entry.Equity, formatCompactMoney(payout)))
	}

	unallocatedPool := fg.EquityPool - fg.EquityAllocated
	if unallocatedPool < 0 {
		unallocatedPool = 0
	}
	if unallocatedPool > 0 {
		poolStyle := lipgloss.NewStyle().Foreground(styles.Gray)
		details.WriteString(poolStyle.Render(fmt.Sprintf("%-30s %6.1f%%  (cancelled at exit)", "Unallocated Pool", unallocatedPool)))
		details.WriteString("\n")
	}

	// Forced acceptance warning
	if founderEquity < 50.0 {
		details.WriteString("\n")
		warnStyle := lipgloss.NewStyle().Foreground(styles.Red).Bold(true)
		details.WriteString(warnStyle.Render("⚠️ WARNING: You don't have majority ownership!"))
		details.WriteString("\n")
		warnDetail := lipgloss.NewStyle().Foreground(styles.Yellow)
		details.WriteString(warnDetail.Render(fmt.Sprintf("   Your equity: %.1f%% (need 50%%+ for control)", founderEquity)))
		details.WriteString("\n")
		details.WriteString(warnDetail.Render("   The board can force this acquisition."))
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(offerBox.Render(details.String())))
	b.WriteString("\n\n")

	// Menu
	if s.acquisitionMenu != nil {
		menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
		b.WriteString(menuContainer.Render(s.acquisitionMenu.View()))
	}
	b.WriteString("\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter select"))

	return b.String()
}

func (s *FounderGameScreen) renderEndAffiliate() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Red).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("END AFFILIATE PROGRAM")))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Red).
		Padding(1, 2)

	if s.endAffiliateMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.endAffiliateMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}

// ==================== Confirm Quit ====================

func (s *FounderGameScreen) rebuildConfirmQuitMenu() {
	items := []components.MenuItem{
		{ID: "resume", Title: "Resume Game", Description: "Continue playing", Icon: "▶️"},
		{ID: "quit", Title: "Quit to Main Menu", Description: "⚠️ Your progress will be lost!", Icon: "🚪"},
	}
	s.confirmQuitMenu = components.NewMenu("QUIT GAME?", items)
	s.confirmQuitMenu.SetSize(55, 8)
	s.confirmQuitMenu.SetHideHelp(true)
}

func (s *FounderGameScreen) handleConfirmQuitSelection(id string) (ScreenModel, tea.Cmd) {
	switch id {
	case "resume":
		s.view = FounderViewMain
		return s, nil
	case "quit":
		return s, SwitchTo(ScreenMainMenu)
	}
	return s, nil
}

func (s *FounderGameScreen) renderConfirmQuit() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Yellow).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⚠️ QUIT GAME?")))
	b.WriteString("\n\n")

	warnStyle := lipgloss.NewStyle().Foreground(styles.Red).Bold(true).Width(s.width).Align(lipgloss.Center)
	b.WriteString(warnStyle.Render("Your current game progress will be lost!"))
	b.WriteString("\n")
	b.WriteString(warnStyle.Render("This cannot be undone."))
	b.WriteString("\n\n")

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2)

	if s.confirmQuitMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.confirmQuitMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc resume • enter select"))

	return b.String()
}

// ==================== Engineer Reallocation ====================

func (s *FounderGameScreen) rebuildEngineerReallocMenu() {
	fg := s.gameData.FounderState
	inProgress := fg.GetInProgressFeatures()

	if len(inProgress) == 0 {
		s.reallocMenu = nil
		s.reallocFeatures = nil
		return
	}

	items := []components.MenuItem{}
	totalEngineers := len(fg.Team.Engineers)
	allocatedEngineers := fg.GetAllocatedEngineers()
	availableEngineers := totalEngineers - allocatedEngineers

	for _, f := range inProgress {
		items = append(items, components.MenuItem{
			ID:          "realloc_add_" + f.Name,
			Title:       fmt.Sprintf("Add engineer to %s", f.Name),
			Description: fmt.Sprintf("Currently: %d engineers, %d%% done", f.AllocatedEngineers, f.DevelopmentProgress),
			Icon:        "➕",
		})
		if f.AllocatedEngineers > 1 {
			items = append(items, components.MenuItem{
				ID:          "realloc_remove_" + f.Name,
				Title:       fmt.Sprintf("Remove engineer from %s", f.Name),
				Description: fmt.Sprintf("Currently: %d engineers, %d%% done", f.AllocatedEngineers, f.DevelopmentProgress),
				Icon:        "➖",
			})
		}
	}

	items = append(items, components.MenuItem{
		ID:       "info",
		Title:    fmt.Sprintf("Available: %d / %d engineers", availableEngineers, totalEngineers),
		Disabled: true,
		Icon:     "ℹ️",
	})

	items = append(items, components.MenuItem{ID: "cancel", Title: "Back", Icon: "←"})

	s.reallocMenu = components.NewMenu("ENGINEER REALLOCATION", items)
	s.reallocMenu.SetSize(60, 15)
	s.reallocMenu.SetHideHelp(true)
	s.reallocFeatures = inProgress
}

func (s *FounderGameScreen) handleEngineerReallocSelection(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	if id == "cancel" {
		s.view = FounderViewRoadmap
		return s, nil
	}

	if strings.HasPrefix(id, "realloc_add_") {
		featureName := strings.TrimPrefix(id, "realloc_add_")
		// Find current allocation
		for _, f := range fg.GetInProgressFeatures() {
			if f.Name == featureName {
				err := fg.ReallocateEngineers(featureName, f.AllocatedEngineers+1)
				if err != nil {
					s.turnMessages = []string{fmt.Sprintf("❌ %v", err)}
				} else {
					s.turnMessages = []string{fmt.Sprintf("✓ Added engineer to %s (now %d engineers)", featureName, f.AllocatedEngineers+1)}
				}
				break
			}
		}
		// Rebuild and stay
		s.rebuildEngineerReallocMenu()
		if s.reallocMenu == nil {
			s.view = FounderViewRoadmap
		}
		return s, nil
	}

	if strings.HasPrefix(id, "realloc_remove_") {
		featureName := strings.TrimPrefix(id, "realloc_remove_")
		for _, f := range fg.GetInProgressFeatures() {
			if f.Name == featureName {
				newCount := f.AllocatedEngineers - 1
				if newCount < 1 {
					newCount = 1
				}
				err := fg.ReallocateEngineers(featureName, newCount)
				if err != nil {
					s.turnMessages = []string{fmt.Sprintf("❌ %v", err)}
				} else {
					s.turnMessages = []string{fmt.Sprintf("✓ Removed engineer from %s (now %d engineers)", featureName, newCount)}
				}
				break
			}
		}
		s.rebuildEngineerReallocMenu()
		if s.reallocMenu == nil {
			s.view = FounderViewRoadmap
		}
		return s, nil
	}

	return s, nil
}

func (s *FounderGameScreen) renderEngineerRealloc() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🔧 ENGINEER REALLOCATION")))
	b.WriteString("\n\n")

	if len(s.turnMessages) > 0 {
		msgStyle := lipgloss.NewStyle().Foreground(styles.Green).Width(s.width).Align(lipgloss.Center)
		for _, msg := range s.turnMessages {
			b.WriteString(msgStyle.Render(msg))
			b.WriteString("\n")
		}
		b.WriteString("\n")
		s.turnMessages = nil
	}

	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	if s.reallocMenu != nil {
		b.WriteString(menuContainer.Render(menuBox.Render(s.reallocMenu.View())))
	}
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter select • esc back"))

	return b.String()
}

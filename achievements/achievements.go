package achievements

import (
	"time"
)

// Achievement represents an unlockable achievement
type Achievement struct {
	ID                   string
	Name                 string
	Description          string
	Icon                 string
	Category             string
	Points               int
	Rarity               string
	Hidden               bool
	RequiredAchievements []string // Must unlock these first
	ChainID              string   // Group related achievements
	ProgressTracking     bool     // Can show progress
	MaxProgress          int      // For progress tracking (e.g., "Win 10 games")
}

// PlayerAchievement tracks when a player unlocked an achievement
type PlayerAchievement struct {
	AchievementID string
	UnlockedAt    time.Time
}

// GameStats contains all stats needed for achievement checking
type GameStats struct {
	// Game mode
	GameMode string // "vc" or "founder"
	
	// Game results
	FinalNetWorth   int64
	ROI             float64
	SuccessfulExits int
	TurnsPlayed     int
	Difficulty      string
	
	// Portfolio details (VC mode)
	InvestmentCount int
	SectorsInvested []string
	TotalInvested   int64
	RiskScores      []float64 // Risk scores of all investments (for achievement tracking)
	
	// Performance (VC mode)
	PositiveInvestments int
	NegativeInvestments int
	BestROI             float64
	WorstROI            float64
	
	// Founder mode stats
	FinalMRR              int64
	FinalValuation        int64
	FinalEquity           float64
	Customers             int
	FundingRoundsRaised   int
	TotalFundingRaised    int64
	HasExited             bool
	ExitType              string // "ipo", "acquisition", "secondary"
	ExitValuation         int64
	MonthsToProfitability int
	RanOutOfCash          bool // True if founder ran out of cash (lost)

	// Phase 1 feature stats
	FeaturesCompleted          int
	InnovationLeader           bool
	EnterpriseFeatures         int
	CustomerLossDuringRoadmap  bool
	EnterpriseCustomers        int
	VerticalConcentration      float64
	PricingExperimentsCompleted int
	PremiumPricingSuccess      bool
	LowTouchCustomers          int
	DealsClosedWon             int
	HighProbabilityClose       bool
	MaxPipelineSize            int

	// Phase 2-3 feature stats
	ContentLeads             int
	SEOScore                 int
	MaxNPS                   int
	CustomerChurnRate        float64
	IntelReportsCommissioned int
	TechDebtKeptLow          bool
	MajorMediaMentions       int
	BoardPressureKeptLow     bool

	// Advanced growth mechanics stats
	AcquisitionsCompleted    int
	AcquisitionSynergy        float64
	FastIntegration           bool
	PlatformLaunched          bool
	NetworkEffectActivated    bool
	MarketplaceRevenue        int64
	DeepIntegrations          int
	NativeIntegrations        int
	PartnershipRevenuePercent float64

	// Crisis management stats
	SecurityScoreMaintained   int
	SecurityIncidentsResolved  int
	LowChurnDuringIncident    bool
	ComplianceCertsAchieved    int
	PRCrisesNavigated         int
	BrandScoreMaintained      int
	EconomicDownturnsSurvived  int
	FundingWinterRaised       bool
	MarketShareGained         bool
	SuccessionPlansCreated     int
	KeyPersonsRetained         bool
	KeyPersonReplaced         bool
	
	// Career stats
	TotalGames      int
	TotalWins       int
	WinStreak       int
	BestNetWorth    int64
	TotalExits      int
}

// Achievement categories
const (
	CategoryWealth      = "Wealth"
	CategoryPerformance = "Performance"
	CategoryStrategy    = "Strategy"
	CategoryCareer      = "Career"
	CategoryChallenge   = "Challenge"
	CategorySpecial     = "Special"
)

// Rarity levels
const (
	RarityCommon    = "Common"
	RarityRare      = "Rare"
	RarityEpic      = "Epic"
	RarityLegendary = "Legendary"
)

// All available achievements
var AllAchievements = map[string]Achievement{
	// Wealth Achievements
	"first_profit": {
		ID:          "first_profit",
		Name:        "First Profit",
		Description: "Make your first dollar of profit",
		Icon:        "$",
		Category:    CategoryWealth,
		Points:      5,
		Rarity:      RarityCommon,
	},
	"millionaire": {
		ID:          "millionaire",
		Name:        "Millionaire",
		Description: "Reach $1,000,000 net worth",
		Icon:        "💰",
		Category:    CategoryWealth,
		Points:      10,
		Rarity:      RarityCommon,
	},
	"multi_millionaire": {
		ID:          "multi_millionaire",
		Name:        "Multi-Millionaire",
		Description: "Reach $5,000,000 net worth",
		Icon:        "💵",
		Category:    CategoryWealth,
		Points:      25,
		Rarity:      RarityRare,
	},
	"deca_millionaire": {
		ID:          "deca_millionaire",
		Name:        "Deca-Millionaire",
		Description: "Reach $10,000,000 net worth",
		Icon:        "🏦",
		Category:    CategoryWealth,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"mega_rich": {
		ID:          "mega_rich",
		Name:        "Mega Rich",
		Description: "Reach $50,000,000 net worth",
		Icon:        "👑",
		Category:    CategoryWealth,
		Points:      100,
		Rarity:      RarityLegendary,
	},
	
	// Performance Achievements
	"break_even": {
		ID:          "break_even",
		Name:        "Break Even",
		Description: "Achieve 0% or better ROI",
		Icon:        "=",
		Category:    CategoryPerformance,
		Points:      5,
		Rarity:      RarityCommon,
	},
	"double_up": {
		ID:          "double_up",
		Name:        "Double Up",
		Description: "Achieve 100%+ ROI",
		Icon:        "📈",
		Category:    CategoryPerformance,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"great_investor": {
		ID:          "great_investor",
		Name:        "Great Investor",
		Description: "Achieve 200%+ ROI",
		Icon:        "⭐",
		Category:    CategoryPerformance,
		Points:      25,
		Rarity:      RarityRare,
	},
	"elite_vc": {
		ID:          "elite_vc",
		Name:        "Elite VC",
		Description: "Achieve 500%+ ROI",
		Icon:        "🏆",
		Category:    CategoryPerformance,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"unicorn_hunter": {
		ID:          "unicorn_hunter",
		Name:        "Unicorn Hunter",
		Description: "Achieve 1000%+ ROI",
		Icon:        "🦄",
		Category:    CategoryPerformance,
		Points:      100,
		Rarity:      RarityLegendary,
	},
	
	// Strategy Achievements
	"diversified": {
		ID:          "diversified",
		Name:        "Diversified",
		Description: "Invest in 5+ companies",
		Icon:        "📊",
		Category:    CategoryStrategy,
		Points:      10,
		Rarity:      RarityCommon,
	},
	"sector_master": {
		ID:          "sector_master",
		Name:        "Sector Master",
		Description: "Invest in 5+ different sectors",
		Icon:        "🏢",
		Category:    CategoryStrategy,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"all_in": {
		ID:          "all_in",
		Name:        "All In",
		Description: "Win with only 1 investment",
		Icon:        "🎲",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityEpic,
	},
	"sector_specialist": {
		ID:          "sector_specialist",
		Name:        "Sector Specialist",
		Description: "Win with all investments in same sector",
		Icon:        "🎯",
		Category:    CategoryStrategy,
		Points:      20,
		Rarity:      RarityRare,
	},
	"exit_master": {
		ID:          "exit_master",
		Name:        "Exit Master",
		Description: "3+ successful exits (5x) in one game",
		Icon:        "🚀",
		Category:    CategoryStrategy,
		Points:      25,
		Rarity:      RarityRare,
	},
	"perfect_portfolio": {
		ID:          "perfect_portfolio",
		Name:        "Perfect Portfolio",
		Description: "Win without any losing investments",
		Icon:        "✨",
		Category:    CategoryStrategy,
		Points:      50,
		Rarity:      RarityEpic,
	},
	
	// Career Achievements
	"first_game": {
		ID:          "first_game",
		Name:        "First Steps",
		Description: "Complete your first game",
		Icon:        "👣",
		Category:    CategoryCareer,
		Points:      5,
		Rarity:      RarityCommon,
	},
	"persistent": {
		ID:          "persistent",
		Name:        "Persistent",
		Description: "Play 10 games",
		Icon:        "💪",
		Category:    CategoryCareer,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"veteran": {
		ID:          "veteran",
		Name:        "Veteran",
		Description: "Play 25 games",
		Icon:        "🎖️",
		Category:    CategoryCareer,
		Points:      25,
		Rarity:      RarityRare,
	},
	"master_investor": {
		ID:          "master_investor",
		Name:        "Master Investor",
		Description: "Play 50 games",
		Icon:        "👑",
		Category:    CategoryCareer,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"win_streak_3": {
		ID:               "win_streak_3",
		Name:             "Hot Streak",
		Description:      "Win 3 games in a row",
		Icon:             "🔥",
		Category:         CategoryCareer,
		Points:           20,
		Rarity:           RarityRare,
		ChainID:          "win_streak",
		ProgressTracking: true,
		MaxProgress:      3,
	},
	"win_streak_5": {
		ID:                   "win_streak_5",
		Name:                 "On Fire",
		Description:          "Win 5 games in a row",
		Icon:                 "⚡",
		Category:             CategoryCareer,
		Points:               40,
		Rarity:               RarityEpic,
		RequiredAchievements: []string{"win_streak_3"},
		ChainID:              "win_streak",
		ProgressTracking:     true,
		MaxProgress:          5,
	},
	
	// Challenge Achievements
	"easy_win": {
		ID:          "easy_win",
		Name:        "Easy Money",
		Description: "Win on Easy difficulty",
		Icon:        "✅",
		Category:    CategoryChallenge,
		Points:      10,
		Rarity:      RarityCommon,
	},
	"medium_win": {
		ID:          "medium_win",
		Name:        "Rising Star",
		Description: "Win on Medium difficulty",
		Icon:        "⭐",
		Category:    CategoryChallenge,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"hard_win": {
		ID:          "hard_win",
		Name:        "Battle Tested",
		Description: "Win on Hard difficulty",
		Icon:        "🛡️",
		Category:    CategoryChallenge,
		Points:      25,
		Rarity:      RarityRare,
	},
	"expert_win": {
		ID:          "expert_win",
		Name:        "Expert Survivor",
		Description: "Win on Expert difficulty",
		Icon:        "💀",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"easy_master": {
		ID:          "easy_master",
		Name:        "Easy Domination",
		Description: "500%+ ROI on Easy",
		Icon:        "🥇",
		Category:    CategoryChallenge,
		Points:      30,
		Rarity:      RarityRare,
	},
	"expert_master": {
		ID:          "expert_master",
		Name:        "Expert Legend",
		Description: "500%+ ROI on Expert",
		Icon:        "🌟",
		Category:    CategoryChallenge,
		Points:      100,
		Rarity:      RarityLegendary,
	},
	"speed_runner": {
		ID:          "speed_runner",
		Name:        "Speed Runner",
		Description: "Win in under 60 turns",
		Icon:        "🏃",
		Category:    CategoryChallenge,
		Points:      30,
		Rarity:      RarityRare,
	},
	
	// Special Achievements
	"lucky_seven": {
		ID:          "lucky_seven",
		Name:        "Lucky Seven",
		Description: "Invest in exactly 7 companies and win",
		Icon:        "🍀",
		Category:    CategorySpecial,
		Points:      15,
		Rarity:      RarityRare,
	},
	"minimalist": {
		ID:          "minimalist",
		Name:        "Minimalist",
		Description: "Win with exactly 2 investments",
		Icon:        "🎯",
		Category:    CategorySpecial,
		Points:      20,
		Rarity:      RarityRare,
	},
	"tech_enthusiast": {
		ID:          "tech_enthusiast",
		Name:        "Tech Enthusiast",
		Description: "Only invest in tech sectors and win",
		Icon:        "💻",
		Category:    CategorySpecial,
		Points:      20,
		Rarity:      RarityRare,
	},
	"clean_investor": {
		ID:          "clean_investor",
		Name:        "Clean Investor",
		Description: "Only invest in CleanTech/AgriTech and win",
		Icon:        "🌱",
		Category:    CategorySpecial,
		Points:      20,
		Rarity:      RarityRare,
	},
	"risk_taker": {
		ID:          "risk_taker",
		Name:        "Risk Taker",
		Description: "Win with only high-risk companies",
		Icon:        "🎲",
		Category:    CategorySpecial,
		Points:      35,
		Rarity:      RarityEpic,
		Hidden:      true,
	},
	"cautious_investor": {
		ID:          "cautious_investor",
		Name:        "Cautious Investor",
		Description: "Win with only low-risk companies",
		Icon:        "🛡️",
		Category:    CategorySpecial,
		Points:      25,
		Rarity:      RarityRare,
	},
	
	// Founder Mode Achievements
	"first_revenue": {
		ID:          "first_revenue",
		Name:        "First Revenue",
		Description: "Generate your first $1,000 MRR",
		Icon:        "💵",
		Category:    CategoryWealth,
		Points:      10,
		Rarity:      RarityCommon,
	},
	"profitable": {
		ID:          "profitable",
		Name:        "Profitable",
		Description: "Reach profitability (positive cash flow)",
		Icon:        "📈",
		Category:    CategoryPerformance,
		Points:      25,
		Rarity:      RarityRare,
	},
	"100k_mrr": {
		ID:          "100k_mrr",
		Name:        "$100K MRR Club",
		Description: "Reach $100,000 monthly recurring revenue",
		Icon:        "🎯",
		Category:    CategoryWealth,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"1m_mrr": {
		ID:          "1m_mrr",
		Name:        "Unicorn MRR",
		Description: "Reach $1,000,000 monthly recurring revenue",
		Icon:        "🦄",
		Category:    CategoryWealth,
		Points:      100,
		Rarity:      RarityLegendary,
	},
	"seed_raised": {
		ID:          "seed_raised",
		Name:        "Seed Raiser",
		Description: "Raise your first funding round",
		Icon:        "🌱",
		Category:    CategoryStrategy,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"series_a": {
		ID:          "series_a",
		Name:        "Series A Graduate",
		Description: "Raise Series A funding",
		Icon:        "🚀",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityRare,
	},
	"ipo_exit": {
		ID:          "ipo_exit",
		Name:        "Public Debut",
		Description: "Take your company public via IPO",
		Icon:        "🏛️",
		Category:    CategorySpecial,
		Points:      75,
		Rarity:      RarityLegendary,
	},
	"acquired": {
		ID:          "acquired",
		Name:        "Acquired",
		Description: "Get acquired by another company",
		Icon:        "🤝",
		Category:    CategorySpecial,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"10000_customers": {
		ID:          "10000_customers",
		Name:        "10K Customers",
		Description: "Reach 10,000 customers",
		Icon:        "👥",
		Category:    CategoryPerformance,
		Points:      40,
		Rarity:      RarityEpic,
	},
	"bootstrapped": {
		ID:          "bootstrapped",
		Name:        "Bootstrapped",
		Description: "Reach $100K MRR without raising funding",
		Icon:        "💪",
		Category:    CategoryChallenge,
		Points:      60,
		Rarity:      RarityLegendary,
		Hidden:      true,
	},

	// Phase 1: Product Roadmap Achievements
	"feature_factory": {
		ID:          "feature_factory",
		Name:        "Feature Factory",
		Description: "Complete 10 product features",
		Icon:        "🔨",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityRare,
	},
	"innovation_leader": {
		ID:          "innovation_leader",
		Name:        "Innovation Leader",
		Description: "Complete a feature before any competitor",
		Icon:        "🚀",
		Category:    CategoryStrategy,
		Points:      40,
		Rarity:      RarityEpic,
	},
	"perfect_roadmap": {
		ID:          "perfect_roadmap",
		Name:        "Perfect Roadmap",
		Description: "Complete all enterprise features with no customer losses",
		Icon:        "🎯",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityLegendary,
		Hidden:      true,
	},

	// Phase 1: Customer Segmentation Achievements
	"enterprise_champion": {
		ID:          "enterprise_champion",
		Name:        "Enterprise Champion",
		Description: "Acquire 100+ enterprise customers",
		Icon:        "👔",
		Category:    CategoryWealth,
		Points:      35,
		Rarity:      RarityEpic,
	},
	"vertical_domination": {
		ID:          "vertical_domination",
		Name:        "Vertical Domination",
		Description: "Have 80% of customers in one vertical",
		Icon:        "🎯",
		Category:    CategoryStrategy,
		Points:      45,
		Rarity:      RarityEpic,
	},

	// Phase 1: Pricing Strategy Achievements
	"pricing_wizard": {
		ID:          "pricing_wizard",
		Name:        "Pricing Wizard",
		Description: "Run 3 successful pricing experiments",
		Icon:        "🧪",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityRare,
	},
	"premium_positioning": {
		ID:          "premium_positioning",
		Name:        "Premium Positioning",
		Description: "Charge 2x market rate and maintain growth",
		Icon:        "💎",
		Category:    CategoryWealth,
		Points:      40,
		Rarity:      RarityEpic,
	},
	"volume_play": {
		ID:          "volume_play",
		Name:        "Volume Play",
		Description: "Have 500+ customers on low-touch plan",
		Icon:        "📊",
		Category:    CategoryPerformance,
		Points:      35,
		Rarity:      RarityRare,
	},

	// Phase 1: Sales Pipeline Achievements
	"sales_machine": {
		ID:          "sales_machine",
		Name:        "Sales Machine",
		Description: "Close 50 deals in one game",
		Icon:        "🤝",
		Category:    CategoryPerformance,
		Points:      30,
		Rarity:      RarityRare,
	},
	"perfect_close": {
		ID:          "perfect_close",
		Name:        "Perfect Close",
		Description: "Close a deal with 90%+ probability",
		Icon:        "💯",
		Category:    CategoryPerformance,
		Points:      25,
		Rarity:      RarityRare,
	},
	"pipeline_master": {
		ID:          "pipeline_master",
		Name:        "Pipeline Master",
		Description: "Maintain 100+ deals in pipeline simultaneously",
		Icon:        "📈",
		Category:    CategoryStrategy,
		Points:      40,
		Rarity:      RarityEpic,
	},

	// Phase 2-3 Achievements
	"content_machine": {
		ID:          "content_machine",
		Name:        "Content Machine",
		Description: "Generate 1000+ inbound leads from content",
		Icon:        "📝",
		Category:    CategoryStrategy,
		Points:      35,
		Rarity:      RarityRare,
	},
	"seo_master": {
		ID:          "seo_master",
		Name:        "SEO Master",
		Description: "Achieve SEO score of 90+",
		Icon:        "🔍",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityRare,
	},
	"customer_champion": {
		ID:          "customer_champion",
		Name:        "Customer Champion",
		Description: "Achieve NPS of 70+",
		Icon:        "⭐",
		Category:    CategoryPerformance,
		Points:      35,
		Rarity:      RarityEpic,
	},
	"churn_slayer": {
		ID:          "churn_slayer",
		Name:        "Churn Slayer",
		Description: "Reduce churn below 2%",
		Icon:        "🛡️",
		Category:    CategoryPerformance,
		Points:      40,
		Rarity:      RarityEpic,
	},
	"know_thy_enemy": {
		ID:          "know_thy_enemy",
		Name:        "Know Thy Enemy",
		Description: "Commission 10 competitive intel reports",
		Icon:        "🔬",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityRare,
	},
	"technical_excellence": {
		ID:          "technical_excellence",
		Name:        "Technical Excellence",
		Description: "Keep tech debt below 20 for 12 months",
		Icon:        "⚙️",
		Category:    CategoryChallenge,
		Points:      45,
		Rarity:      RarityEpic,
	},
	"media_darling": {
		ID:          "media_darling",
		Name:        "Media Darling",
		Description: "Featured in 5+ major outlets",
		Icon:        "📰",
		Category:    CategoryStrategy,
		Points:      35,
		Rarity:      RarityEpic,
	},
	"board_whisperer": {
		ID:          "board_whisperer",
		Name:        "Board Whisperer",
		Description: "Maintain board pressure below 30 for 12 months",
		Icon:        "🤝",
		Category:    CategoryStrategy,
		Points:      40,
		Rarity:      RarityEpic,
	},

	// Advanced Growth Mechanics Achievements
	"serial_acquirer": {
		ID:          "serial_acquirer",
		Name:        "Serial Acquirer",
		Description: "Complete 3 acquisitions",
		Icon:        "🏢",
		Category:    CategoryStrategy,
		Points:      45,
		Rarity:      RarityEpic,
	},
	"synergy_master": {
		ID:          "synergy_master",
		Name:        "Synergy Master",
		Description: "Achieve 50%+ revenue boost from acquisition",
		Icon:        "⚡",
		Category:    CategoryPerformance,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"integration_expert": {
		ID:          "integration_expert",
		Name:        "Integration Expert",
		Description: "Complete acquisition integration in <4 months",
		Icon:        "🔧",
		Category:    CategoryPerformance,
		Points:      35,
		Rarity:      RarityRare,
	},
	"platform_builder": {
		ID:          "platform_builder",
		Name:        "Platform Builder",
		Description: "Reach 1000+ customers with platform model",
		Icon:        "🌐",
		Category:    CategoryStrategy,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"network_effect": {
		ID:          "network_effect",
		Name:        "Network Effect",
		Description: "Activate strong network effects",
		Icon:        "🔗",
		Category:    CategoryStrategy,
		Points:      45,
		Rarity:      RarityEpic,
	},
	"marketplace_master": {
		ID:          "marketplace_master",
		Name:        "Marketplace Master",
		Description: "Generate $500k+ from marketplace fees",
		Icon:        "💰",
		Category:    CategoryWealth,
		Points:      40,
		Rarity:      RarityEpic,
	},
	"integration_master": {
		ID:          "integration_master",
		Name:        "Integration Master",
		Description: "Create 10+ deep integrations",
		Icon:        "🔌",
		Category:    CategoryStrategy,
		Points:      35,
		Rarity:      RarityRare,
	},
	"ecosystem_builder": {
		ID:          "ecosystem_builder",
		Name:        "Ecosystem Builder",
		Description: "Generate 30%+ revenue from partnerships",
		Icon:        "🤝",
		Category:    CategoryPerformance,
		Points:      40,
		Rarity:      RarityEpic,
	},

	// Crisis Management Achievements
	"security_champion": {
		ID:          "security_champion",
		Name:        "Security Champion",
		Description: "Maintain 90+ security score for 12 months",
		Icon:        "🔒",
		Category:    CategoryChallenge,
		Points:      45,
		Rarity:      RarityEpic,
	},
	"incident_response": {
		ID:          "incident_response",
		Name:        "Incident Response",
		Description: "Resolve critical security incident with <5% churn",
		Icon:        "🛡️",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"compliance_master": {
		ID:          "compliance_master",
		Name:        "Compliance Master",
		Description: "Achieve all major certifications (SOC2, ISO27001, HIPAA, GDPR)",
		Icon:        "✅",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"crisis_manager": {
		ID:          "crisis_manager",
		Name:        "Crisis Manager",
		Description: "Successfully navigate 3+ PR crises",
		Icon:        "📰",
		Category:    CategoryChallenge,
		Points:      45,
		Rarity:      RarityEpic,
	},
	"brand_resilience": {
		ID:          "brand_resilience",
		Name:        "Brand Resilience",
		Description: "Maintain brand score >70 through crisis",
		Icon:        "💪",
		Category:    CategoryChallenge,
		Points:      40,
		Rarity:      RarityEpic,
	},
	"media_master": {
		ID:          "media_master",
		Name:        "Media Master",
		Description: "Turn PR crisis into positive coverage",
		Icon:        "📺",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"recession_survivor": {
		ID:          "recession_survivor",
		Name:        "Recession Survivor",
		Description: "Survive severe economic downturn",
		Icon:        "📉",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"funding_winter_warrior": {
		ID:          "funding_winter_warrior",
		Name:        "Funding Winter Warrior",
		Description: "Raise funding during funding winter",
		Icon:        "❄️",
		Category:    CategoryChallenge,
		Points:      55,
		Rarity:      RarityLegendary,
	},
	"market_leader": {
		ID:          "market_leader",
		Name:        "Market Leader",
		Description: "Gain market share during sector crash",
		Icon:        "👑",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"succession_ready": {
		ID:          "succession_ready",
		Name:        "Succession Ready",
		Description: "Have succession plans for all key persons",
		Icon:        "📋",
		Category:    CategoryStrategy,
		Points:      40,
		Rarity:      RarityEpic,
	},
	"retention_master": {
		ID:          "retention_master",
		Name:        "Retention Master",
		Description: "Keep all key persons for 60 months",
		Icon:        "👥",
		Category:    CategoryChallenge,
		Points:      45,
		Rarity:      RarityEpic,
	},
	"crisis_leader": {
		ID:          "crisis_leader",
		Name:        "Crisis Leader",
		Description: "Successfully replace key person with <10% impact",
		Icon:        "🎯",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
	},
	
	// Diversification Chain
	"diversified_starter": {
		ID:                   "diversified_starter",
		Name:                 "Diversification Starter",
		Description:          "Invest in 3 different sectors",
		Icon:                 "🌱",
		Category:             CategoryStrategy,
		Points:               10,
		Rarity:               RarityCommon,
		RequiredAchievements: []string{},
		ChainID:              "diversification",
		ProgressTracking:     false,
	},
	"portfolio_manager": {
		ID:                   "portfolio_manager",
		Name:                 "Portfolio Manager",
		Description:          "Invest in 5 different sectors",
		Icon:                 "📊",
		Category:             CategoryStrategy,
		Points:               25,
		Rarity:               RarityRare,
		RequiredAchievements: []string{"diversified_starter"},
		ChainID:              "diversification",
		ProgressTracking:     false,
	},
	"investment_conglomerate": {
		ID:                   "investment_conglomerate",
		Name:                 "Investment Conglomerate",
		Description:          "Invest in all available sectors",
		Icon:                 "🏢",
		Category:             CategoryStrategy,
		Points:               50,
		Rarity:               RarityEpic,
		RequiredAchievements: []string{"portfolio_manager"},
		ChainID:              "diversification",
		ProgressTracking:     false,
	},
	
	// Win Streak Chain - Extended with win_streak_10
	"win_streak_10": {
		ID:                   "win_streak_10",
		Name:                 "Unstoppable",
		Description:          "Win 10 games in a row",
		Icon:                 "🔥🔥🔥",
		Category:             CategoryCareer,
		Points:               75,
		Rarity:               RarityLegendary,
		RequiredAchievements: []string{"win_streak_5"},
		ChainID:              "win_streak",
		ProgressTracking:     true,
		MaxProgress:          10,
	},
	
	// Games Played Chain with Progress Tracking
	"games_10": {
		ID:               "games_10",
		Name:             "Getting Started",
		Description:      "Play 10 games",
		Icon:             "🎮",
		Category:         CategoryCareer,
		Points:           10,
		Rarity:           RarityCommon,
		ChainID:          "games_played",
		ProgressTracking: true,
		MaxProgress:      10,
	},
	"games_50": {
		ID:                   "games_50",
		Name:                 "Dedicated Player",
		Description:          "Play 50 games",
		Icon:                 "🎯",
		Category:             CategoryCareer,
		Points:               25,
		Rarity:               RarityRare,
		RequiredAchievements: []string{"games_10"},
		ChainID:              "games_played",
		ProgressTracking:     true,
		MaxProgress:          50,
	},
	"games_100": {
		ID:                   "games_100",
		Name:                 "Century Mark",
		Description:          "Play 100 games",
		Icon:                 "💯",
		Category:             CategoryCareer,
		Points:               50,
		Rarity:               RarityEpic,
		RequiredAchievements: []string{"games_50"},
		ChainID:              "games_played",
		ProgressTracking:     true,
		MaxProgress:          100,
	},
	"games_500": {
		ID:                   "games_500",
		Name:                 "Veteran",
		Description:          "Play 500 games",
		Icon:                 "🏅",
		Category:             CategoryCareer,
		Points:               100,
		Rarity:               RarityLegendary,
		RequiredAchievements: []string{"games_100"},
		ChainID:              "games_played",
		ProgressTracking:     true,
		MaxProgress:          500,
	},
	
	// Hidden Mystery Achievements
	"mystery_investor": {
		ID:               "mystery_investor",
		Name:             "Mystery Investor",
		Description:      "Invest in all 30+ startups across multiple games",
		Icon:             "🎭",
		Category:         CategorySpecial,
		Points:           100,
		Rarity:           RarityLegendary,
		Hidden:           true,
		ProgressTracking: true,
		MaxProgress:      30,
	},
	"perfect_month": {
		ID:          "perfect_month",
		Name:        "Perfect Month",
		Description: "Have all portfolio companies increase in value in one turn",
		Icon:        "✨",
		Category:    CategorySpecial,
		Points:      75,
		Rarity:      RarityLegendary,
		Hidden:      true,
	},
	"phoenix": {
		ID:          "phoenix",
		Name:        "Phoenix",
		Description: "Win a game after having negative net worth",
		Icon:        "🔥",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
		Hidden:      true,
	},
	"day_trader": {
		ID:          "day_trader",
		Name:        "Day Trader",
		Description: "Complete a game in under 5 minutes",
		Icon:        "⚡",
		Category:    CategoryChallenge,
		Points:      30,
		Rarity:      RarityRare,
		Hidden:      true,
	},
	
	// Investment Expertise Chain
	"investment_novice": {
		ID:               "investment_novice",
		Name:             "Investment Novice",
		Description:      "Make 10 successful investments",
		Icon:             "📈",
		Category:         CategoryPerformance,
		Points:           10,
		Rarity:           RarityCommon,
		ChainID:          "investment_count",
		ProgressTracking: true,
		MaxProgress:      10,
	},
	"investment_expert": {
		ID:                   "investment_expert",
		Name:                 "Investment Expert",
		Description:          "Make 50 successful investments",
		Icon:                 "📊",
		Category:             CategoryPerformance,
		Points:               30,
		Rarity:               RarityRare,
		RequiredAchievements: []string{"investment_novice"},
		ChainID:              "investment_count",
		ProgressTracking:     true,
		MaxProgress:          50,
	},
	"investment_master": {
		ID:                   "investment_master",
		Name:                 "Investment Master",
		Description:          "Make 100 successful investments",
		Icon:                 "💎",
		Category:             CategoryPerformance,
		Points:               75,
		Rarity:               RarityLegendary,
		RequiredAchievements: []string{"investment_expert"},
		ChainID:              "investment_count",
		ProgressTracking:     true,
		MaxProgress:          100,
	},
}

// CheckAchievements checks which achievements were earned this game
func CheckAchievements(stats GameStats, previouslyUnlocked []string) []Achievement {
	unlocked := make(map[string]bool)
	for _, id := range previouslyUnlocked {
		unlocked[id] = true
	}
	
	var newAchievements []Achievement
	
	for id, achievement := range AllAchievements {
		if unlocked[id] {
			continue
		}
		
		if checkAchievement(id, stats) {
			newAchievements = append(newAchievements, achievement)
		}
	}
	
	return newAchievements
}

func checkAchievement(id string, stats GameStats) bool {
	// Determine if player won
	var won bool
	if stats.GameMode == "founder" {
		// For founder mode: won = exited successfully OR reached max turns without running out of cash
		won = stats.HasExited || !stats.RanOutOfCash
	} else {
		// For VC mode: won = positive ROI
		won = stats.ROI > 0
	}
	
	switch id {
	// Wealth
	case "first_profit":
		return stats.FinalNetWorth > stats.TotalInvested
	case "millionaire":
		return stats.FinalNetWorth >= 1000000
	case "multi_millionaire":
		return stats.FinalNetWorth >= 5000000
	case "deca_millionaire":
		return stats.FinalNetWorth >= 10000000
	case "mega_rich":
		return stats.FinalNetWorth >= 50000000
		
	// Performance
	case "break_even":
		return stats.ROI >= 0
	case "double_up":
		return stats.ROI >= 100
	case "great_investor":
		return stats.ROI >= 200
	case "elite_vc":
		return stats.ROI >= 500
	case "unicorn_hunter":
		return stats.ROI >= 1000
		
	// Strategy
	case "diversified":
		return stats.InvestmentCount >= 5
	case "sector_master":
		return len(stats.SectorsInvested) >= 5
	case "all_in":
		return stats.InvestmentCount == 1 && won
	case "sector_specialist":
		return len(stats.SectorsInvested) == 1 && stats.InvestmentCount > 1 && won
	case "exit_master":
		return stats.SuccessfulExits >= 3
	case "perfect_portfolio":
		return stats.NegativeInvestments == 0 && stats.InvestmentCount > 0 && won
		
	// Career
	case "first_game":
		return stats.TotalGames >= 1
	case "persistent":
		return stats.TotalGames >= 10
	case "veteran":
		return stats.TotalGames >= 25
	case "master_investor":
		return stats.TotalGames >= 50
	case "win_streak_3":
		return stats.WinStreak >= 3
	case "win_streak_5":
		return stats.WinStreak >= 5
		
	// Challenge
	case "easy_win":
		return stats.Difficulty == "Easy" && won
	case "medium_win":
		return stats.Difficulty == "Medium" && won
	case "hard_win":
		return stats.Difficulty == "Hard" && won
	case "expert_win":
		return stats.Difficulty == "Expert" && won
	case "easy_master":
		return stats.Difficulty == "Easy" && stats.ROI >= 500
	case "expert_master":
		return stats.Difficulty == "Expert" && stats.ROI >= 500
	case "speed_runner":
		return stats.TurnsPlayed < 60 && won
		
	// Special
	case "lucky_seven":
		return stats.InvestmentCount == 7 && won
	case "minimalist":
		return stats.InvestmentCount == 2 && won
	case "tech_enthusiast":
		// Check if all sectors are tech-related
		// Must have at least one investment and won
		if stats.InvestmentCount == 0 || !won {
			return false
		}
		techSectors := map[string]bool{
			"CloudTech": true, "SaaS": true, "DeepTech": true,
			"FinTech": true, "HealthTech": true, "EdTech": true,
			"LegalTech": true, "Gaming": true, "Security": true,
		}
		// Must have at least one sector
		if len(stats.SectorsInvested) == 0 {
			return false
		}
		// All sectors must be tech sectors
		for _, sector := range stats.SectorsInvested {
			if !techSectors[sector] {
				return false // Found non-tech sector
			}
		}
		return true
	case "clean_investor":
		// Check if only CleanTech/AgriTech
		for _, sector := range stats.SectorsInvested {
			if sector != "CleanTech" && sector != "AgriTech" {
				return false
			}
		}
		return len(stats.SectorsInvested) > 0 && won
	case "risk_taker":
		// Win with only high-risk companies (risk score > 0.6)
		if stats.InvestmentCount == 0 || !won {
			return false
		}
		for _, risk := range stats.RiskScores {
			if risk <= 0.6 {
				return false // Not all high-risk
			}
		}
		return true
	case "cautious_investor":
		// Win with only low-risk companies (risk score < 0.3)
		if stats.InvestmentCount == 0 || !won {
			return false
		}
		for _, risk := range stats.RiskScores {
			if risk >= 0.3 {
				return false // Not all low-risk
			}
		}
		return true
		
	// Founder Mode Achievements
	case "first_revenue":
		return stats.GameMode == "founder" && stats.FinalMRR >= 1000
	case "profitable":
		return stats.GameMode == "founder" && stats.MonthsToProfitability > 0
	case "100k_mrr":
		return stats.GameMode == "founder" && stats.FinalMRR >= 100000
	case "1m_mrr":
		return stats.GameMode == "founder" && stats.FinalMRR >= 1000000
	case "seed_raised":
		return stats.GameMode == "founder" && stats.FundingRoundsRaised >= 1
	case "series_a":
		return stats.GameMode == "founder" && stats.FundingRoundsRaised >= 2 // Assuming Seed + Series A
	case "ipo_exit":
		return stats.GameMode == "founder" && stats.HasExited && stats.ExitType == "ipo"
	case "acquired":
		return stats.GameMode == "founder" && stats.HasExited && stats.ExitType == "acquisition"
	case "10000_customers":
		return stats.GameMode == "founder" && stats.Customers >= 10000
	case "bootstrapped":
		return stats.GameMode == "founder" && stats.FinalMRR >= 100000 && stats.FundingRoundsRaised == 0
	
	// Phase 1: Product Roadmap
	case "feature_factory":
		// Note: Need to add FeaturesCompleted to GameStats
		return stats.GameMode == "founder" && stats.FeaturesCompleted >= 10
	case "innovation_leader":
		// Note: Track if any feature completed before competitors launched same
		return stats.GameMode == "founder" && stats.InnovationLeader
	case "perfect_roadmap":
		// Note: Track enterprise features completed and customer losses
		return stats.GameMode == "founder" && stats.EnterpriseFeatures >= 3 && stats.CustomerLossDuringRoadmap == false

	// Phase 1: Customer Segmentation
	case "enterprise_champion":
		return stats.GameMode == "founder" && stats.EnterpriseCustomers >= 100
	case "vertical_domination":
		return stats.GameMode == "founder" && stats.VerticalConcentration >= 0.80

	// Phase 1: Pricing Strategy
	case "pricing_wizard":
		return stats.GameMode == "founder" && stats.PricingExperimentsCompleted >= 3
	case "premium_positioning":
		return stats.GameMode == "founder" && stats.PremiumPricingSuccess
	case "volume_play":
		return stats.GameMode == "founder" && stats.LowTouchCustomers >= 500

	// Phase 1: Sales Pipeline
	case "sales_machine":
		return stats.GameMode == "founder" && stats.DealsClosedWon >= 50
	case "perfect_close":
		return stats.GameMode == "founder" && stats.HighProbabilityClose
	case "pipeline_master":
		return stats.GameMode == "founder" && stats.MaxPipelineSize >= 100

	// Phase 2-3
	case "content_machine":
		return stats.GameMode == "founder" && stats.ContentLeads >= 1000
	case "seo_master":
		return stats.GameMode == "founder" && stats.SEOScore >= 90
	case "customer_champion":
		return stats.GameMode == "founder" && stats.MaxNPS >= 70
	case "churn_slayer":
		return stats.GameMode == "founder" && stats.CustomerChurnRate <= 0.02
	case "know_thy_enemy":
		return stats.GameMode == "founder" && stats.IntelReportsCommissioned >= 10
	case "technical_excellence":
		return stats.GameMode == "founder" && stats.TechDebtKeptLow
	case "media_darling":
		return stats.GameMode == "founder" && stats.MajorMediaMentions >= 5
	case "board_whisperer":
		return stats.GameMode == "founder" && stats.BoardPressureKeptLow

	// Advanced Growth Mechanics Achievements
	case "serial_acquirer":
		return stats.GameMode == "founder" && stats.AcquisitionsCompleted >= 3
	case "synergy_master":
		return stats.GameMode == "founder" && stats.AcquisitionSynergy >= 0.5
	case "integration_expert":
		return stats.GameMode == "founder" && stats.FastIntegration
	case "platform_builder":
		return stats.GameMode == "founder" && stats.PlatformLaunched && stats.Customers >= 1000
	case "network_effect":
		return stats.GameMode == "founder" && stats.NetworkEffectActivated
	case "marketplace_master":
		return stats.GameMode == "founder" && stats.MarketplaceRevenue >= 500000
	case "integration_master":
		return stats.GameMode == "founder" && stats.DeepIntegrations >= 10
	case "ecosystem_builder":
		return stats.GameMode == "founder" && stats.PartnershipRevenuePercent >= 0.30

	// Crisis Management Achievements
	case "security_champion":
		return stats.GameMode == "founder" && stats.SecurityScoreMaintained >= 90
	case "incident_response":
		return stats.GameMode == "founder" && stats.SecurityIncidentsResolved > 0 && stats.LowChurnDuringIncident
	case "compliance_master":
		return stats.GameMode == "founder" && stats.ComplianceCertsAchieved >= 4
	case "crisis_manager":
		return stats.GameMode == "founder" && stats.PRCrisesNavigated >= 3
	case "brand_resilience":
		return stats.GameMode == "founder" && stats.BrandScoreMaintained >= 70
	case "media_master":
		return stats.GameMode == "founder" && stats.PRCrisesNavigated > 0 && stats.BrandScoreMaintained >= 70
	case "recession_survivor":
		return stats.GameMode == "founder" && stats.EconomicDownturnsSurvived > 0
	case "funding_winter_warrior":
		return stats.GameMode == "founder" && stats.FundingWinterRaised
	case "market_leader":
		return stats.GameMode == "founder" && stats.MarketShareGained
	case "succession_ready":
		return stats.GameMode == "founder" && stats.SuccessionPlansCreated >= 3
	case "retention_master":
		return stats.GameMode == "founder" && stats.KeyPersonsRetained
	case "crisis_leader":
		return stats.GameMode == "founder" && stats.KeyPersonReplaced
	}
	
	return false
}

// CalculateCareerLevel calculates player level based on achievement points
func CalculateCareerLevel(totalPoints int) (level int, title string, nextLevelPoints int) {
	levels := []struct {
		points int
		level  int
		title  string
	}{
		{0, 0, "Intern"},
		{25, 1, "Analyst"},
		{75, 2, "Associate"},
		{150, 3, "Senior Associate"},
		{250, 4, "Principal"},
		{400, 5, "Partner"},
		{600, 6, "Senior Partner"},
		{850, 7, "Managing Partner"},
		{1150, 8, "Elite VC"},
		{1500, 9, "Master Investor"},
		{2000, 10, "Legendary Investor"},
	}
	
	for i := len(levels) - 1; i >= 0; i-- {
		if totalPoints >= levels[i].points {
			nextLevel := 2001 // Max
			if i < len(levels)-1 {
				nextLevel = levels[i+1].points
			}
			return levels[i].level, levels[i].title, nextLevel
		}
	}
	
	return 0, "Intern", 25
}

// GetAchievementsByCategory returns achievements for a category
func GetAchievementsByCategory(category string) []Achievement {
	var achievements []Achievement
	for _, ach := range AllAchievements {
		if ach.Category == category && !ach.Hidden {
			achievements = append(achievements, ach)
		}
	}
	return achievements
}

// GetRarityColor returns color code for rarity
func GetRarityColor(rarity string) int {
	switch rarity {
	case RarityCommon:
		return 37 // White
	case RarityRare:
		return 36 // Cyan
	case RarityEpic:
		return 35 // Magenta
	case RarityLegendary:
		return 33 // Yellow
	default:
		return 37
	}
}

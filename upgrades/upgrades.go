package upgrades

import (
	"time"
)

// Upgrade represents a purchasable upgrade
type Upgrade struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Category    string
	Cost        int
	Tier        int // 1-5 (higher = more expensive/powerful)
	Effect      string
}

// Upgrade categories
const (
	CategoryInvestmentTerms = "Investment Terms"
	CategoryFinancialPerks  = "Financial Perks"
	CategoryInformation     = "Information & Intel"
	CategoryGameModes       = "Game Modes"
	CategoryBoardPowers     = "Board Seat Powers"
	CategorySpecialAbilities = "Special Abilities"
	CategoryFounderPerks   = "Founder Perks"
)

// All available upgrades
var AllUpgrades = map[string]Upgrade{
	// Tier 1: Quick Wins (100-150 pts)
	"fund_booster": {
		ID:          "fund_booster",
		Name:        "Fund Booster",
		Description: "+10% starting cash on all difficulties",
		Icon:        "💰",
		Category:    CategoryFinancialPerks,
		Cost:        100,
		Tier:        1,
		Effect:      "cash_multiplier:1.1",
	},
	"due_diligence": {
		ID:          "due_diligence",
		Name:        "Due Diligence",
		Description: "See exact risk score numbers (not just labels)",
		Icon:        "🔍",
		Category:    CategoryInformation,
		Cost:        150,
		Tier:        1,
		Effect:      "show_risk_numbers:true",
	},
	"early_access": {
		ID:          "early_access",
		Name:        "Early Access",
		Description: "See 2 extra startups before investment phase starts",
		Icon:        "👁️",
		Category:    CategoryInformation,
		Cost:        100,
		Tier:        1,
		Effect:      "extra_startups:2",
	},
	
	// Tier 2: Medium Impact (150-200 pts)
	"enhanced_safe_discount": {
		ID:          "enhanced_safe_discount",
		Name:        "Enhanced SAFE Discount",
		Description: "SAFE discount increases from 20% → 25%",
		Icon:        "📈",
		Category:    CategoryInvestmentTerms,
		Cost:        150,
		Tier:        2,
		Effect:      "safe_discount:0.25",
	},
	"management_fee_reduction": {
		ID:          "management_fee_reduction",
		Name:        "Management Fee Reduction",
		Description: "Management fees reduced from 2% → 1.5%",
		Icon:        "💵",
		Category:    CategoryFinancialPerks,
		Cost:        150,
		Tier:        2,
		Effect:      "management_fee:0.015",
	},
	"revenue_tracker": {
		ID:          "revenue_tracker",
		Name:        "Revenue Tracker",
		Description: "See monthly revenue growth trends",
		Icon:        "📊",
		Category:    CategoryInformation,
		Cost:        100,
		Tier:        2,
		Effect:      "show_revenue_trends:true",
	},
	"fee_waiver": {
		ID:          "fee_waiver",
		Name:        "Fee Waiver",
		Description: "No management fees for first 12 months",
		Icon:        "🎫",
		Category:    CategoryFinancialPerks,
		Cost:        300,
		Tier:        2,
		Effect:      "fee_waiver_months:12",
	},
	"seed_accelerator": {
		ID:          "seed_accelerator",
		Name:        "Seed Accelerator",
		Description: "First investment gets 25% equity bonus",
		Icon:        "🚀",
		Category:    CategoryFinancialPerks,
		Cost:        400,
		Tier:        2,
		Effect:      "first_investment_bonus:0.25",
	},
	
	// Tier 3: Major Strategic Unlocks (200-350 pts)
	"double_board_seat": {
		ID:          "double_board_seat",
		Name:        "Double Board Seat",
		Description: "Get 2 board seats per $100k investment (double voting power)",
		Icon:        "🏛️",
		Category:    CategoryBoardPowers,
		Cost:        200,
		Tier:        3,
		Effect:      "board_seats_multiplier:2",
	},
	"super_pro_rata": {
		ID:          "super_pro_rata",
		Name:        "Super Pro-Rata",
		Description: "Can invest up to 50% of round (vs 20% max)",
		Icon:        "🎯",
		Category:    CategoryInvestmentTerms,
		Cost:        200,
		Tier:        3,
		Effect:      "max_investment_percent:0.50",
	},
	"follow_on_reserve_boost": {
		ID:          "follow_on_reserve_boost",
		Name:        "Follow-On Reserve Boost",
		Description: "+$200k to follow-on reserve",
		Icon:        "💎",
		Category:    CategoryFinancialPerks,
		Cost:        200,
		Tier:        3,
		Effect:      "follow_on_reserve_bonus:200000",
	},
	"founder_network": {
		ID:          "founder_network",
		Name:        "Founder Network",
		Description: "See one extra startup during initial investment phase",
		Icon:        "🤝",
		Category:    CategoryInformation,
		Cost:        200,
		Tier:        3,
		Effect:      "extra_startups_per_round:1",
	},
	"angel_investor": {
		ID:          "angel_investor",
		Name:        "Angel Investor",
		Description: "+$100k bonus starting cash (stacks with Fund Booster)",
		Icon:        "👼",
		Category:    CategoryFinancialPerks,
		Cost:        250,
		Tier:        3,
		Effect:      "bonus_cash:100000",
	},
	"liquidation_preference_2x": {
		ID:          "liquidation_preference_2x",
		Name:        "2x Liquidation Preference",
		Description: "Unlock 2x liquidation preference option (get paid 2x before others)",
		Icon:        "💎",
		Category:    CategoryInvestmentTerms,
		Cost:        250,
		Tier:        3,
		Effect:      "liquidation_pref_2x:true",
	},
	"strategic_advisor": {
		ID:          "strategic_advisor",
		Name:        "Strategic Advisor",
		Description: "Get preview of next board vote before it happens",
		Icon:        "🔮",
		Category:    CategoryBoardPowers,
		Cost:        250,
		Tier:        3,
		Effect:      "preview_board_votes:true",
	},
	"market_intelligence": {
		ID:          "market_intelligence",
		Name:        "Market Intelligence",
		Description: "See which categories are trending up/down (based on valuations)",
		Icon:        "📈",
		Category:    CategoryInformation,
		Cost:        250,
		Tier:        3,
		Effect:      "show_sector_trends:true",
	},
	
	// Tier 4: Game Mode Unlocks (200-400 pts)
	"speed_mode": {
		ID:          "speed_mode",
		Name:        "Speed Mode",
		Description: "30 turns instead of 60 (faster games)",
		Icon:        "⚡",
		Category:    CategoryGameModes,
		Cost:        200,
		Tier:        4,
		Effect:      "max_turns:30",
	},
	"endurance_mode": {
		ID:          "endurance_mode",
		Name:        "Endurance Mode",
		Description: "120 turns instead of 60 (longer games)",
		Icon:        "🏃",
		Category:    CategoryGameModes,
		Cost:        250,
		Tier:        4,
		Effect:      "max_turns:120",
	},
	
	// Tier 5: Special Abilities (400+ pts)
	"portfolio_insurance": {
		ID:          "portfolio_insurance",
		Name:        "Portfolio Insurance",
		Description: "Protect one investment from down rounds per game",
		Icon:        "🛡️",
		Category:    CategorySpecialAbilities,
		Cost:        500,
		Tier:        5,
		Effect:      "protect_from_down_round:1",
	},
	"time_machine": {
		ID:          "time_machine",
		Name:        "Time Machine",
		Description: "Rewind one investment decision per game",
		Icon:        "⏰",
		Category:    CategorySpecialAbilities,
		Cost:        600,
		Tier:        5,
		Effect:      "rewind_decision:1",
	},
	
	// Founder Mode Upgrades
	"fast_track": {
		ID:          "fast_track",
		Name:        "Fast Track",
		Description: "Start with 10% more product maturity",
		Icon:        "🚀",
		Category:    CategoryFounderPerks,
		Cost:        200,
		Tier:        3,
		Effect:      "product_maturity_boost:0.1",
	},
	"sales_boost": {
		ID:          "sales_boost",
		Name:        "Sales Boost",
		Description: "+15% to initial MRR",
		Icon:        "💰",
		Category:    CategoryFounderPerks,
		Cost:        250,
		Tier:        3,
		Effect:      "initial_mrr_boost:0.15",
	},
	"lower_burn": {
		ID:          "lower_burn",
		Name:        "Lower Burn",
		Description: "-10% monthly team costs",
		Icon:        "💸",
		Category:    CategoryFounderPerks,
		Cost:        300,
		Tier:        3,
		Effect:      "team_cost_reduction:0.1",
	},
	"better_terms": {
		ID:          "better_terms",
		Name:        "Better Terms",
		Description: "Raise funding with 5% less equity given away",
		Icon:        "📝",
		Category:    CategoryFounderPerks,
		Cost:        350,
		Tier:        4,
		Effect:      "equity_reduction:0.05",
	},
	"quick_hire": {
		ID:          "quick_hire",
		Name:        "Quick Hire",
		Description: "First 3 hires cost 50% less",
		Icon:        "👥",
		Category:    CategoryFounderPerks,
		Cost:        200,
		Tier:        3,
		Effect:      "first_hires_discount:0.5",
	},
	"market_insight": {
		ID:          "market_insight",
		Name:        "Market Insight",
		Description: "See competitor threat levels",
		Icon:        "🔍",
		Category:    CategoryFounderPerks,
		Cost:        250,
		Tier:        3,
		Effect:      "show_competitor_threats:true",
	},
	"churn_shield": {
		ID:          "churn_shield",
		Name:        "Churn Shield",
		Description: "Reduce churn by 10% permanently",
		Icon:        "🛡️",
		Category:    CategoryFounderPerks,
		Cost:        300,
		Tier:        3,
		Effect:      "churn_reduction:0.1",
	},
	"cloud_free_first_year": {
		ID:          "cloud_free_first_year",
		Name:        "Cloud Free First Year",
		Description: "No cloud compute costs for first 12 months",
		Icon:        "☁️",
		Category:    CategoryFounderPerks,
		Cost:        300,
		Tier:        3,
		Effect:      "free_cloud_months:12",
	},
}

// GetUpgradesByCategory returns upgrades for a category
func GetUpgradesByCategory(category string) []Upgrade {
	var upgrades []Upgrade
	for _, upgrade := range AllUpgrades {
		if upgrade.Category == category {
			upgrades = append(upgrades, upgrade)
		}
	}
	return upgrades
}

// GetUpgradesByTier returns upgrades for a tier
func GetUpgradesByTier(tier int) []Upgrade {
	var upgrades []Upgrade
	for _, upgrade := range AllUpgrades {
		if upgrade.Tier == tier {
			upgrades = append(upgrades, upgrade)
		}
	}
	return upgrades
}

// GetAllCategories returns all upgrade categories
func GetAllCategories() []string {
	categories := make(map[string]bool)
	for _, upgrade := range AllUpgrades {
		categories[upgrade.Category] = true
	}
	
	result := []string{}
	for cat := range categories {
		result = append(result, cat)
	}
	return result
}

// PlayerUpgrade tracks when a player purchased an upgrade
type PlayerUpgrade struct {
	UpgradeID   string
	PurchasedAt time.Time
}

// IsOwned checks if an upgrade is owned by checking the purchased upgrades list
func IsOwned(upgradeID string, ownedUpgrades []string) bool {
	for _, owned := range ownedUpgrades {
		if owned == upgradeID {
			return true
		}
	}
	return false
}


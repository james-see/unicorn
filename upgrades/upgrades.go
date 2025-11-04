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
		Description: "See one extra startup per round",
		Icon:        "🤝",
		Category:    CategoryInformation,
		Cost:        200,
		Tier:        3,
		Effect:      "extra_startups_per_round:1",
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


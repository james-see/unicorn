package game

import (
	"fmt"
	"math/rand"
)

// ValueAddAction represents operational support provided to portfolio companies
type ValueAddAction struct {
	ActionType     string // "recruiting", "sales", "technical", "board", "marketing"
	CompanyName    string
	Cost           int64
	Relationship   float64 // Relationship score increase
	ValuationBoost float64 // One-time or sustained valuation increase
	RiskReduction  float64 // Risk score reduction
	Duration       int     // How many turns the effect lasts
	AppliedTurn    int     // When the action was taken
	Description    string
}

// ValueAddType defines the characteristics of each value-add action
type ValueAddType struct {
	ID                string
	Name              string
	Description       string
	Cost              int64
	MinRelationship   float64
	MaxRelationship   float64
	MinValBoost       float64
	MaxValBoost       float64
	RiskReduction     float64
	Duration          int // Number of turns effect lasts
	RequiresBoardSeat bool
	MinEquityPct      float64 // Minimum equity required if no board seat
}

// GetAvailableValueAddTypes returns all value-add action types
func GetAvailableValueAddTypes() []ValueAddType {
	return []ValueAddType{
		{
			ID:                "recruiting",
			Name:              "Recruiting Support",
			Description:       "Help recruit key executives and engineering talent",
			Cost:              20000,
			MinRelationship:   5.0,
			MaxRelationship:   10.0,
			MinValBoost:       0.02, // 2% min
			MaxValBoost:       0.05, // 5% max
			RiskReduction:     0.0,
			Duration:          3, // Effect spread over 3 turns
			RequiresBoardSeat: false,
			MinEquityPct:      5.0,
		},
		{
			ID:                "sales",
			Name:              "Sales Introductions",
			Description:       "Introduce founders to potential customers and partners",
			Cost:              15000,
			MinRelationship:   3.0,
			MaxRelationship:   8.0,
			MinValBoost:       0.01, // 1% min
			MaxValBoost:       0.04, // 4% max
			RiskReduction:     0.0,
			Duration:          2,
			RequiresBoardSeat: false,
			MinEquityPct:      5.0,
		},
		{
			ID:                "technical",
			Name:              "Technical Advisory",
			Description:       "Provide technical architecture and scaling advice",
			Cost:              25000,
			MinRelationship:   5.0,
			MaxRelationship:   12.0,
			MinValBoost:       0.03, // 3% min
			MaxValBoost:       0.07, // 7% max
			RiskReduction:     0.05, // Reduces risk by 5%
			Duration:          4,
			RequiresBoardSeat: false,
			MinEquityPct:      5.0,
		},
		{
			ID:                "board_leadership",
			Name:              "Board Leadership",
			Description:       "Take active board leadership role and provide strategic guidance",
			Cost:              10000,
			MinRelationship:   4.0,
			MaxRelationship:   10.0,
			MinValBoost:       0.0,
			MaxValBoost:       0.0,
			RiskReduction:     0.0,
			Duration:          1,
			RequiresBoardSeat: true, // Requires board seat
			MinEquityPct:      0.0,
		},
		{
			ID:                "marketing",
			Name:              "Marketing Guidance",
			Description:       "Help with go-to-market strategy and brand positioning",
			Cost:              15000,
			MinRelationship:   5.0,
			MaxRelationship:   10.0,
			MinValBoost:       0.02, // 2% min
			MaxValBoost:       0.04, // 4% max
			RiskReduction:     0.0,
			Duration:          3,
			RequiresBoardSeat: false,
			MinEquityPct:      5.0,
		},
	}
}

// CanProvideValueAdd checks if player can provide value-add to a company
func (gs *GameState) CanProvideValueAdd(companyName string, actionType ValueAddType) (bool, string) {
	// Find the investment
	var inv *Investment
	for i := range gs.Portfolio.Investments {
		if gs.Portfolio.Investments[i].CompanyName == companyName {
			inv = &gs.Portfolio.Investments[i]
			break
		}
	}

	if inv == nil {
		return false, "You have not invested in this company"
	}

	// Check if company already received value-add this turn
	for _, action := range gs.ActiveValueAddActions {
		if action.CompanyName == companyName && action.AppliedTurn == gs.Portfolio.Turn {
			return false, "Company already received value-add this turn"
		}
	}

	// Check board seat requirement
	if actionType.RequiresBoardSeat && !inv.Terms.HasBoardSeat {
		return false, "Requires board seat"
	}

	// Check equity requirement
	if !actionType.RequiresBoardSeat && inv.EquityPercent < actionType.MinEquityPct {
		return false, fmt.Sprintf("Requires at least %.1f%% equity", actionType.MinEquityPct)
	}

	// Check cash
	if gs.Portfolio.Cash < actionType.Cost {
		return false, fmt.Sprintf("Insufficient funds (need $%d)", actionType.Cost)
	}

	// Check attention points (max 2 value-add actions per turn)
	actionsThisTurn := 0
	for _, action := range gs.ActiveValueAddActions {
		if action.AppliedTurn == gs.Portfolio.Turn {
			actionsThisTurn++
		}
	}
	if actionsThisTurn >= 2 {
		return false, "Maximum 2 value-add actions per turn"
	}

	return true, ""
}

// ProvideValueAdd executes a value-add action
func (gs *GameState) ProvideValueAdd(companyName string, actionTypeID string) error {
	// Find action type
	var actionType *ValueAddType
	for _, vat := range GetAvailableValueAddTypes() {
		if vat.ID == actionTypeID {
			actionType = &vat
			break
		}
	}

	if actionType == nil {
		return fmt.Errorf("invalid action type")
	}

	// Check if can provide
	can, reason := gs.CanProvideValueAdd(companyName, *actionType)
	if !can {
		return fmt.Errorf(reason)
	}

	// Find the investment
	var inv *Investment
	invIdx := -1
	for i := range gs.Portfolio.Investments {
		if gs.Portfolio.Investments[i].CompanyName == companyName {
			inv = &gs.Portfolio.Investments[i]
			invIdx = i
			break
		}
	}

	// Calculate actual values within ranges
	relationshipIncrease := actionType.MinRelationship +
		rand.Float64()*(actionType.MaxRelationship-actionType.MinRelationship)

	valuationBoost := 0.0
	if actionType.MaxValBoost > 0 {
		valuationBoost = actionType.MinValBoost +
			rand.Float64()*(actionType.MaxValBoost-actionType.MinValBoost)
	}

	// Create the action
	action := ValueAddAction{
		ActionType:     actionTypeID,
		CompanyName:    companyName,
		Cost:           actionType.Cost,
		Relationship:   relationshipIncrease,
		ValuationBoost: valuationBoost,
		RiskReduction:  actionType.RiskReduction,
		Duration:       actionType.Duration,
		AppliedTurn:    gs.Portfolio.Turn,
		Description:    actionType.Description,
	}

	// Apply immediate effects
	inv.RelationshipScore = ApplyRelationshipChange(inv.RelationshipScore, relationshipIncrease)
	inv.ValueAddProvided++
	inv.LastInteraction = gs.Portfolio.Turn

	// Deduct cost
	gs.Portfolio.Cash -= actionType.Cost

	// Add to active actions (for sustained effects)
	gs.ActiveValueAddActions = append(gs.ActiveValueAddActions, action)

	// Update the investment in the slice
	gs.Portfolio.Investments[invIdx] = *inv

	return nil
}

// ProcessActiveValueAddActions applies ongoing effects from active value-add actions
func (gs *GameState) ProcessActiveValueAddActions() []string {
	messages := []string{}

	for i, action := range gs.ActiveValueAddActions {
		turnsSinceAction := gs.Portfolio.Turn - action.AppliedTurn

		// Check if action is still active
		if turnsSinceAction >= action.Duration {
			continue // Expired
		}

		// Find the company
		for idx := range gs.AvailableStartups {
			if gs.AvailableStartups[idx].Name == action.CompanyName {
				startup := &gs.AvailableStartups[idx]

				// Apply valuation boost (spread over duration)
				if action.ValuationBoost > 0 {
					monthlyBoost := action.ValuationBoost / float64(action.Duration)
					boostAmount := int64(float64(startup.Valuation) * monthlyBoost)
					startup.Valuation += boostAmount

					// Update investment valuation
					for invIdx := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[invIdx].CompanyName == action.CompanyName {
							gs.Portfolio.Investments[invIdx].CurrentValuation = startup.Valuation
						}
					}

					if turnsSinceAction == 0 {
						messages = append(messages, fmt.Sprintf(
							"Your %s support for %s is showing results (+%.1f%% valuation)",
							action.ActionType, action.CompanyName, monthlyBoost*100))
					}
				}

				// Apply risk reduction (one-time on first turn)
				if action.RiskReduction > 0 && turnsSinceAction == 0 {
					startup.RiskScore -= action.RiskReduction
					if startup.RiskScore < 0.1 {
						startup.RiskScore = 0.1 // Floor
					}
				}

				break
			}
		}

		// Mark as processed for this turn
		gs.ActiveValueAddActions[i] = action
	}

	// Remove expired actions
	activeActions := []ValueAddAction{}
	for _, action := range gs.ActiveValueAddActions {
		if gs.Portfolio.Turn-action.AppliedTurn < action.Duration {
			activeActions = append(activeActions, action)
		}
	}
	gs.ActiveValueAddActions = activeActions

	return messages
}

// GetValueAddOpportunities returns companies that can receive value-add
func (gs *GameState) GetValueAddOpportunities() []string {
	companies := []string{}

	for _, inv := range gs.Portfolio.Investments {
		// Must have board seat OR 5%+ equity
		if inv.Terms.HasBoardSeat || inv.EquityPercent >= 5.0 {
			// Check if not already helped this turn
			alreadyHelped := false
			for _, action := range gs.ActiveValueAddActions {
				if action.CompanyName == inv.CompanyName && action.AppliedTurn == gs.Portfolio.Turn {
					alreadyHelped = true
					break
				}
			}

			if !alreadyHelped {
				companies = append(companies, inv.CompanyName)
			}
		}
	}

	return companies
}

// GetTotalValueAddInvestment returns total spent on value-add
func (gs *GameState) GetTotalValueAddInvestment() int64 {
	total := int64(0)
	for _, action := range gs.ActiveValueAddActions {
		total += action.Cost
	}
	return total
}

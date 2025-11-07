# VC Reputation System - Integration Guide

## Overview
This document describes how the VC Reputation System integrates into the existing game flow.

## System Components Created

### Core Game Logic Files
- `game/reputation.go` - Reputation calculations and scoring
- `game/founders.go` - Founder relationship system
- `game/value_add.go` - Operational value-add actions
- `game/due_diligence.go` - Due diligence system
- `game/secondary_market.go` - Secondary stake sales
- `game/deal_flow.go` - Reputation-based deal quality

### UI Files
- `ui/reputation_ui.go` - Reputation display
- `ui/value_add_ui.go` - Value-add menus
- `ui/dd_ui.go` - Due diligence flow
- `ui/secondary_market_ui.go` - Secondary market offers

### Database
- Added `vc_reputation` table in `database/database.go`
- Functions: `GetVCReputation()`, `SaveVCReputation()`

## Integration Points

### 1. Game Initialization (`ui/vc_ui.go` - `PlayVCMode()`)

```go
// Load player reputation
dbRep, err := database.GetVCReputation(username)
if err != nil {
    dbRep = &database.VCReputation{
        PlayerName:       username,
        PerformanceScore: 50.0,
        FounderScore:     50.0,
        MarketScore:      50.0,
    }
}

// Convert to game reputation
gs.PlayerReputation = &game.VCReputation{
    PlayerName:       dbRep.PlayerName,
    PerformanceScore: dbRep.PerformanceScore,
    FounderScore:     dbRep.FounderScore,
    MarketScore:      dbRep.MarketScore,
    TotalGamesPlayed: dbRep.TotalGamesPlayed,
    SuccessfulExits:  dbRep.SuccessfulExits,
    AvgROILast5:      dbRep.AvgROILast5,
}

// Generate startups based on reputation (replaces LoadStartups)
gs.AvailableStartups, err = game.GenerateStartupsWithReputation(
    gs.PlayerReputation, 
    30, // or 20 depending on difficulty
    "startups/")

// Show reputation summary
DisplayReputationSummary(gs.PlayerReputation)
```

### 2. Investment Phase (`ui/vc_ui.go` - `investmentPhase()`)

```go
// Before investment, optionally run DD (Manual Mode only)
ddLevel := ShowDueDiligenceMenu(gs, &startup, amount, autoMode)
if ddLevel == "cancelled" {
    continue // Investment cancelled after DD
}

// After investment is made, initialize founder relationship
investment.FounderName = game.GenerateFounderName()
investment.RelationshipScore = game.CalculateInitialRelationship(
    selectedTerms, 
    ddLevel != "none", 
    amount)
investment.HasDueDiligence = ddLevel != "none"
investment.DDLevel = ddLevel
investment.LastInteraction = gs.Portfolio.Turn
investment.ValueAddProvided = 0

// Apply reputation bonus to relationship
bonus := game.GetReputationBonus(gs.PlayerReputation)
game.ApplyReputationBonusToInvestment(&investment, bonus)
```

### 3. Turn Processing (`ui/vc_ui.go` - `PlayTurn()`)

```go
// At start of turn, process active value-add actions
messages := gs.ProcessActiveValueAddActions()
for _, msg := range messages {
    fmt.Println(msg)
}

// Generate relationship events for investments
for i := range gs.Portfolio.Investments {
    inv := &gs.Portfolio.Investments[i]
    event := game.GenerateRelationshipEvent(inv, gs.Portfolio.Turn)
    if event != nil {
        fmt.Printf("\n%s\n", event.Description)
        inv.RelationshipScore = game.ApplyRelationshipChange(
            inv.RelationshipScore, 
            event.ScoreChange)
        
        // Check if can be fired from board
        if inv.Terms.HasBoardSeat && game.CanBeFiredFromBoard(inv.RelationshipScore) {
            fmt.Printf("⚠️  Poor relationship with %s led to board removal!\n", 
                inv.FounderName)
            inv.Terms.HasBoardSeat = false
        }
    }
}

// Generate secondary market offers
newOffers := gs.GenerateSecondaryOffers()
gs.SecondaryMarketOffers = append(gs.SecondaryMarketOffers, newOffers...)

// Process offer expirations
expiredMsgs := gs.ProcessSecondaryOfferExpirations()
for _, msg := range expiredMsgs {
    fmt.Println(msg)
}
```

### 4. After Turn (Manual Mode Only) (`ui/vc_ui.go` - `PlayTurn()`)

```go
if !autoMode {
    // Show value-add opportunities
    ShowValueAddMenu(gs, autoMode)
    
    // Show secondary market offers
    ShowSecondaryMarketOffers(gs, autoMode)
}
```

### 5. Game End (`ui/vc_ui.go` - `PlayVCMode()`)

```go
// Calculate game results
netWorth, roi, successfulExits := gs.GetFinalScore()

// Calculate average founder relationship
totalRelationship := 0.0
count := 0
for _, inv := range gs.Portfolio.Investments {
    totalRelationship += inv.RelationshipScore
    count++
}
avgFounderRelationship := 50.0
if count > 0 {
    avgFounderRelationship = totalRelationship / float64(count)
}

// Get achievement points and win streak
achievementPoints := getTotalAchievementPoints(username)
winStreak := getWinStreak(username, roi > 0)

// Update reputation
hadSuccessfulExit := successfulExits > 0
updatedRep := game.UpdateReputationAfterGame(
    gs.PlayerReputation,
    roi,
    hadSuccessfulExit,
    achievementPoints,
    winStreak)

// Update founder score
updatedRep.UpdateFounderScore(avgFounderRelationship)

// Save to database
dbRep := &database.VCReputation{
    PlayerName:       updatedRep.PlayerName,
    PerformanceScore: updatedRep.PerformanceScore,
    FounderScore:     updatedRep.FounderScore,
    MarketScore:      updatedRep.MarketScore,
    TotalGamesPlayed: updatedRep.TotalGamesPlayed,
    SuccessfulExits:  updatedRep.SuccessfulExits,
    AvgROILast5:      updatedRep.AvgROILast5,
}
database.SaveVCReputation(dbRep)

// Show reputation changes
fmt.Println("\n" + strings.Repeat("=", 70))
color.Cyan("REPUTATION UPDATE")
fmt.Println(strings.Repeat("=", 70))
fmt.Printf("Performance: %.1f → %.1f\n", 
    gs.PlayerReputation.PerformanceScore, updatedRep.PerformanceScore)
fmt.Printf("Founder:     %.1f → %.1f\n", 
    gs.PlayerReputation.FounderScore, updatedRep.FounderScore)
fmt.Printf("Market:      %.1f → %.1f\n", 
    gs.PlayerReputation.MarketScore, updatedRep.MarketScore)
fmt.Printf("Overall:     %.1f → %.1f (%s)\n",
    gs.PlayerReputation.GetAggregateReputation(),
    updatedRep.GetAggregateReputation(),
    updatedRep.GetReputationLevel())
```

### 6. Main Menu Integration

Add option to view reputation:

```go
// In main menu
fmt.Println("9. View Your VC Reputation")

// Handler
case "9":
    rep, err := database.GetVCReputation(username)
    if err == nil {
        gameRep := convertToGameReputation(rep)
        DisplayReputation(gameRep)
    }
```

### 7. Portfolio Dashboard Updates

Show founder relationships in portfolio view:

```go
// For each investment
emoji := game.GetRelationshipEmoji(inv.RelationshipScore)
level := game.GetRelationshipLevel(inv.RelationshipScore)
fmt.Printf("  Founder: %s | Relationship: %s %s (%.0f/100)\n",
    inv.FounderName, emoji, level, inv.RelationshipScore)

if inv.ValueAddProvided > 0 {
    fmt.Printf("  Value-add actions provided: %d\n", inv.ValueAddProvided)
}
```

## Automated vs Manual Mode

### Manual Mode
- Full access to all features
- DD before investment
- Value-add menu after each turn
- Secondary market offer review
- Founder relationship actions

### Automated Mode
- Skip all interactive menus
- No DD (direct to investment)
- No value-add actions
- Auto-decline secondary offers
- Passive relationship evolution only
- Reputation still affects deal flow quality

## Testing Checklist

- [ ] Reputation loads correctly on game start
- [ ] Deal quality varies by reputation tier
- [ ] DD system works and applies findings
- [ ] Founder relationships initialize properly
- [ ] Value-add actions apply effects correctly
- [ ] Secondary market offers generate and process
- [ ] Relationship events trigger appropriately
- [ ] Reputation updates correctly at game end
- [ ] Database saves/loads reputation
- [ ] Automated mode skips interactive features
- [ ] Manual mode shows all features

## Balance Notes

- Reputation builds slowly (3-5 games to Tier 1)
- Value-add ROI: 6-12 turns to break even
- DD cost vs benefit: 15-30% chance of significant finding
- Secondary market friction: 20-30% discount
- Founder relationships: 10% impact on exit value
- Deal flow quality: 15-25% improvement at Tier 1

## Future Enhancements

1. Multiplayer reputation comparison
2. Reputation leaderboards
3. More value-add action types
4. Dynamic DD costs based on company size
5. Founder reputation tracking
6. Cross-game founder referrals
7. Reputation decay over time
8. Industry-specific reputation scores


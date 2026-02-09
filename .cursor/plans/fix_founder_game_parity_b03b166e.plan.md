---
name: Fix Founder Game Parity
overview: Bring the founder mode TUI to full feature parity with the old legacy UI and the VC TUI, restoring missing end-game flow, achievements/XP, board management, SaaS metrics, acquisition offers, and enriching several simplified views.
todos:
  - id: founder-results
    content: Create tui/founder_results.go with multi-phase end-game ceremony (results, XP, level-up, leaderboard, achievements) and wire into app.go
    status: completed
  - id: founder-setup
    content: Enhance tui/founder_setup.go with difficulty selection, returning player detection, player upgrades/profile loading
    status: completed
  - id: acquisition-offers
    content: Add FounderViewAcquisitionOffer to tui/founder_game.go with cap table payout, forced acceptance for <50% equity, accept/decline flow
    status: completed
  - id: saas-metrics
    content: Add SaaS metrics panel (LTV:CAC, CAC Payback, Rule of 40, Burn Multiple, Magic Number) to founder dashboard
    status: completed
  - id: board-advisors
    content: "Enrich board management: advisor expertise selection with real names, cost confirmation, chairman removal, advisor buyback, board table visual, sentiment display"
    status: completed
  - id: dashboard-enrichment
    content: Add monthly highlights, customer health, low cash warning, active partnerships, strategic opportunity to main dashboard
    status: completed
  - id: detail-views
    content: Enrich team roster (vesting), customer deals (individual table), financials (margins/deductions), affiliate stats views
    status: completed
  - id: equity-pool-input
    content: Add numeric input for equity pool expansion percentage (1-10%) instead of hardcoded 2%
    status: completed
isProject: false
---

# Fix Founder Mode TUI for Full Feature Parity

## Context

The v4.0.0 TUI refactor replaced the old `ui/founder_ui.go` (5,451 lines) with `tui/founder_game.go`. While much was ported, critical features are still missing -- most notably the entire end-game ceremony, XP/achievements, and several board/advisor management features. The VC side (`tui/vc_results.go`, `tui/vc_setup.go`, `tui/vc_turn.go`) already has all of these, so we have clear patterns to follow.

---

## 1. Create Founder Results Screen (CRITICAL)

Create new file `tui/founder_results.go` modeled on `[tui/vc_results.go](tui/vc_results.go)`.

- Add `ScreenFounderResults` to the Screen enum in `[tui/app.go](tui/app.go)`
- Multi-phase ceremony matching VC pattern: `PhaseResults -> PhaseXP -> PhaseLevelUp -> PhaseLeaderboard -> PhaseAchievements -> PhaseDone`
- **PhaseResults**: Show outcome (IPO/Acquisition/Secondary/Bankrupt/Survived), exit valuation, cap table payout breakdown, team summary, funding rounds, performance rating
- **PhaseXP**: Calculate and award XP (game completion, positive ROI, successful exit type bonuses like IPO +500, acquisition +300, profitability +100, achievement bonuses) -- port logic from `ui/founder_ui.go:3434-3494`
- **PhaseLevelUp**: Level-up celebration with old->new level, bonus points, new unlocks (same as VC)
- **PhaseLeaderboard**: Show local score saved, offer global leaderboard submission
- **PhaseAchievements**: Check and display newly unlocked achievements using `achievements.CheckAchievements()` with `GameMode: "founder"` and full founder stats -- port stat collection from `ui/founder_ui.go:3284-3396`
- **Score Saving**: Save via `database.SaveGameScore()` with `Difficulty: "Founder"` -- port from `ui/founder_ui.go:3266-3281`

**Wire up**: Change the game-over handler in `[tui/founder_game.go:968-970](tui/founder_game.go)` from `SwitchTo(ScreenMainMenu)` to `SwitchTo(ScreenFounderResults)`.

## 2. Enhance Founder Setup Screen

Update `[tui/founder_setup.go](tui/founder_setup.go)` to match VC setup depth:

- Add **difficulty selection** step (Easy/Medium/Hard/Expert) with level-gating like VC setup (`tui/vc_setup.go:186-225`)
- Add **returning player detection** using `database.GetPlayerStats()` with welcome-back message
- Load **player upgrades** from DB via `database.GetPlayerUpgrades()` 
- Load **player profile** for level checking
- Add progress steps: Name -> Company -> Category -> Difficulty -> Ready
- Pass difficulty to `FounderState` (will need a `Difficulty` field or apply modifiers to burn rate/competition/funding availability)

## 3. Acquisition Offer Handling (HIGH)

Add in-game acquisition offer display and interaction to `[tui/founder_game.go](tui/founder_game.go)`:

- Add `FounderViewAcquisitionOffer` view
- After `ProcessMonth()`, check `fs.CheckForAcquisition()` -- if offer returned, switch to acquisition offer view
- Render: acquirer name, offer amount, cap table payout breakdown, due diligence quality, terms quality
- Handle **forced acceptance** when founder equity < 50% (board forces exit)
- Accept/Decline menu with consequences display
- Port logic from `ui/founder_ui.go:1434-1563`

## 4. Key SaaS Metrics on Dashboard (HIGH)

Add SaaS metrics panel to the main founder dashboard in `[tui/founder_game.go](tui/founder_game.go)`:

- **LTV:CAC Ratio** via `fs.CalculateLTVToCAC()` with color-coded health (green >3, yellow 1-3, red <1)
- **CAC Payback** via `fs.CalculateCACPayback()` (green <12mo, yellow 12-24, red >24)
- **Rule of 40** via `fs.CalculateRuleOf40()` (green >40, yellow 20-40, red <20)
- **Burn Multiple** via `fs.CalculateBurnMultiple()` (green <1, yellow 1-2, red >2)
- **Magic Number** via `fs.CalculateMagicNumber()` (green >1, yellow 0.5-1, red <0.5)
- Port from `ui/founder_ui.go:357-457`

## 5. Enrich Board & Advisor Management (MEDIUM)

Enhance the existing `FounderViewBoard` / `FounderViewBoardAction` in `[tui/founder_game.go](tui/founder_game.go)`:

- **Advisor expertise selection**: When adding advisor, show 5 expertise options (sales, product, fundraising, operations, strategy) with real Silicon Valley names (Marc Andreessen, Reid Hoffman, etc.) -- port from `ui/founder_ui.go:2548-2580`
- **Advisor cost confirmation**: Show setup fee ($10-50k), optional monthly retainer, equity cost, before/after cash/equity -- port from `ui/founder_ui.go:2628-2714`
- **Remove chairman** as explicit menu option with consequences display
- **Advisor buyback option** on removal: calculate buyback cost, show buyback vs no-buyback consequences -- port from `ui/founder_ui.go:2906-2966`
- **Board table visual**: Render a styled board composition view showing chairman, advisors with expertise, investor board members, contribution scores -- port from `ui/founder_ui.go:2438-2537`
- **Board/Investor Sentiment display** on main dashboard: sentiment indicator + pressure gauge + warnings about founder replacement -- port from `ui/founder_ui.go:521-559`

## 6. Enhance Dashboard Information Density (MEDIUM)

Add to the main view in `[tui/founder_game.go](tui/founder_game.go)`:

- **Monthly Highlights** (wins & concerns) via `fs.GenerateMonthlyHighlights()` -- port from `ui/founder_ui.go:459-492`
- **Customer Health** summary (healthy/at-risk/critical with MRR at risk) via `fs.GetCustomerHealthSegments()` -- port from `ui/founder_ui.go:494-519`
- **Low Cash Warning** when `fs.NeedsLowCashWarning()` returns true -- port from `ui/founder_ui.go:224-246`
- **Active Partnerships** display with type, months remaining, benefits -- port from `ui/founder_ui.go:586-637`
- **Strategic Opportunity** proactive display on dashboard when pending -- port from `ui/founder_ui.go:562-583`

## 7. Enrich Detail Views (MEDIUM)

- **Team Roster**: Add equity vesting details, cliff status, vested months, market assignment, salary -- port from `ui/founder_ui.go:3940-3993`
- **Customer Deals**: Show individual customer table with ID, source, deal size, term, health score -- port from `ui/founder_ui.go:3571-3619`
- **Financials**: Add gross/net margins, revenue deduction breakdown (taxes 20%, processing 3%, overhead 5%, savings 5%), affiliate costs, global market costs -- port from `ui/founder_ui.go:3799-3903`
- **Affiliate Program**: Full stats view (launch month, affiliates, commission rate, setup cost, platform fee, customers acquired, revenue, avg per affiliate) -- port from `ui/founder_ui.go:3659-3713`

## 8. Equity Pool Expansion Input

Currently the TUI defaults to expanding equity pool by 2%. Add a numeric input to let the player specify 1-10% like the old UI.

---

## Key Files to Modify

- **New**: `tui/founder_results.go` (~500-700 lines, modeled on `tui/vc_results.go`)
- **Modify**: `[tui/app.go](tui/app.go)` - add `ScreenFounderResults` + wiring
- **Modify**: `[tui/founder_setup.go](tui/founder_setup.go)` - add difficulty + player detection
- **Modify**: `[tui/founder_game.go](tui/founder_game.go)` - acquisition offers, SaaS metrics, board enhancements, dashboard enrichment, detail views

## Reference Files (port from)

- `[ui/founder_ui.go](ui/founder_ui.go)` - all old features to port
- `[tui/vc_results.go](tui/vc_results.go)` - pattern for results screen
- `[tui/vc_setup.go](tui/vc_setup.go)` - pattern for setup enhancements


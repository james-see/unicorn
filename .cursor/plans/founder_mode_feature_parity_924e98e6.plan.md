---
name: Founder Mode Feature Parity
overview: Implement missing TUI menu options for features that exist in the backend but are not accessible through the new Bubble Tea UI, restoring feature parity with the old terminal UI.
todos:
  - id: strategic-opportunities
    content: Add Strategic Opportunities menu and handler (when PendingOpportunity != nil)
    status: pending
  - id: solicit-feedback
    content: Add Solicit Customer Feedback menu item (when customers > 0)
    status: pending
  - id: end-affiliate
    content: Add End Affiliate Program menu item (when AffiliateProgram != nil)
    status: pending
  - id: board-chairman
    content: Add Set Chairman and Fire Board Member to board menu
    status: pending
  - id: content-marketing
    content: "Add Content Marketing launch/manage/view menu (unlock: marketing hire OR 200k MRR)"
    status: pending
  - id: cs-playbooks
    content: "Add CS Playbooks launch/manage menu (unlock: CS hire OR 100 customers)"
    status: pending
  - id: competitive-intel
    content: "Add Competitive Intelligence menu (unlock: Series A OR 5+ competitors)"
    status: pending
  - id: referral-program
    content: "Add Referral Program launch/manage/end menu (unlock: 10+ customers)"
    status: pending
  - id: tech-debt
    content: "Add Technical Debt view and refactor menu (unlock: 5+ engineers OR 1M MRR)"
    status: pending
  - id: investor-updates
    content: "Add Investor Updates monthly composition (unlock: first funding)"
    status: pending
  - id: board-requests
    content: Add Board Requests menu (customer intros, recruiting, advice)
    status: pending
  - id: test-full-game
    content: End-to-end test all founder mode features work together
    status: pending
isProject: false
---

# Founder Mode Feature Parity Plan

## Problem Summary

The TUI refactor (Nov 2025) created a comprehensive Bubble Tea UI but left out several features that existed in the old `ui/founder_ui.go`. The game loop processes many features automatically, but players cannot actively manage them through the TUI.

## Missing Features Analysis

### Category 1: Features with Backend but NO TUI Menu


| Feature           | Backend File                           | Old UI Location | Status                             |
| ----------------- | -------------------------------------- | --------------- | ---------------------------------- |
| Content Marketing | `founder/founder_content_marketing.go` | Not in old UI   | Auto-runs if active, cannot launch |
| CS Playbooks      | `founder/founder_cs_playbooks.go`      | Not in old UI   | Auto-runs if active, cannot launch |
| Competitive Intel | `founder/founder_competitive_intel.go` | Not in old UI   | NOT integrated at all              |
| Technical Debt    | `founder/founder_technical_debt.go`    | Not in old UI   | Auto-accumulates, cannot refactor  |
| Referral Program  | `founder/founder_advisors.go`          | Not in old UI   | Full backend, no UI                |


### Category 2: Features in Old UI, Missing in TUI


| Feature                   | Backend                     | Old UI                             | TUI Status |
| ------------------------- | --------------------------- | ---------------------------------- | ---------- |
| Solicit Customer Feedback | `SolicitCustomerFeedback()` | Menu item when customers > 0       | Missing    |
| Strategic Opportunities   | `PendingOpportunity`        | Menu item when opportunity pending | Missing    |
| End Affiliate Program     | `EndAffiliateProgram()`     | "11j. End Affiliate Program"       | Missing    |
| Set Chairman              | `SetChairman()`             | In board management                | Missing    |
| Fire Board Member         | `FireBoardMember()`         | In board management                | Missing    |
| Investor Updates          | `InvestorUpdates`           | Built into board system            | Missing    |
| Board Requests            | `BoardRequests`             | Built into board system            | Missing    |


## Implementation Plan

### Phase 1: Quick Wins - Wire Existing Backend to TUI

These require only adding menu items and handlers in [tui/founder_game.go](tui/founder_game.go):

**1. Strategic Opportunities** (Priority: High)

- Add menu item when `fg.PendingOpportunity != nil`
- Reference: [ui/founder_ui.go](ui/founder_ui.go) lines 4006-4129 for `handleStrategicOpportunity()`
- Show opportunity details, accept/decline options

**2. Solicit Customer Feedback** (Priority: High)

- Add to VIEW DATA section when `fg.Customers > 0`
- Call `fg.SolicitCustomerFeedback()` 
- Reference: [ui/founder_ui.go](ui/founder_ui.go) lines 3621-3650

**3. End Affiliate Program** (Priority: Medium)

- Add menu item when `fg.AffiliateProgram != nil`
- Call `fg.EndAffiliateProgram(transitionCustomers bool)`
- Reference: [ui/founder_ui.go](ui/founder_ui.go) lines 891-897

**4. Board Management Enhancements** (Priority: Medium)

- Add "Set Chairman" option to board menu
- Add "Fire Board Member" option (investor type, requires 51%+ ownership)
- Reference: [founder/founder_advisors.go](founder/founder_advisors.go) lines 479-768

### Phase 2: Launch New Feature Programs

These need new menu sections/submenus:

**5. Launch Content Marketing** (Unlock: Marketing hire OR $200k MRR)

- Add menu item to launch with budget selection ($10-50k/month)
- Add view for current content program status
- Calls: `fg.LaunchContentProgram(budget)`, `fg.EndContentProgram()`

**6. Launch CS Playbooks** (Unlock: CS hire OR 100+ customers)  

- Add menu to launch playbooks (Onboarding, Health, Upsell, Renewal, Churn Prevention)
- Calls: `fg.LaunchCSPlaybook(name, budget)`

**7. Competitive Intelligence** (Unlock: Series A OR 5+ competitors)

- Add menu to hire analyst ($80-120k/year)
- Add menu to commission intel reports ($20-50k per competitor)
- Calls: `fg.LaunchCompetitiveIntel(salary)`, `fg.CommissionIntelReport(competitor, cost)`

**8. Referral Program** (Unlock: 10+ customers)

- Add menu to launch referral program
- Set reward per referral, reward type
- Calls: `fg.LaunchReferralProgram(reward, type)`, `fg.EndReferralProgram()`

**9. Technical Debt Refactoring** (Unlock: 5+ engineers OR $1M+ MRR)

- Add view for current tech debt level and impacts
- Add option to refactor ($50-100k cost, allocate engineers)
- Calls: `fg.RefactorTechDebt(cost, engineers)`

### Phase 3: Board/Investor Relations

**10. Investor Updates** (Unlock: First funding raised)

- Add monthly update composition (transparency level: full, optimistic, selective)

**11. Board Requests** (Unlock: First funding raised)

- Add ability to ask board for: customer intros, recruiting help, strategic advice

## File Changes Required

Primary file: [tui/founder_game.go](tui/founder_game.go)

Key patterns to follow:

- Add menu items to `getMainMenuItems()` with unlock conditions
- Add view states to `FounderViewState` enum
- Create `rebuild*Menu()` functions for submenus
- Create `handle*Selection()` functions for actions
- Add render functions in `View()` switch

## Testing Strategy

After each feature:

1. Build: `go build`
2. Run the game: `./unicorn`
3. Start founder mode
4. Progress to unlock conditions
5. Verify menu appears and actions work
6. Verify turn messages show results


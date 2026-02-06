---
name: Port Founder Mode to TUI
overview: Port all ~25 monthly turn options from the old `ui/founder_ui.go` to the new TUI in `tui/founder_game.go`, creating feature parity with full gameplay including hiring, funding, partnerships, strategic decisions, and data views.
todos:
  - id: phase1-hiring
    content: Port hiring system with all employee types and C-suite
    status: completed
  - id: phase1-firing
    content: Port firing system with department listing and confirmation
    status: completed
  - id: phase1-marketing
    content: Port marketing spend with amount selection
    status: completed
  - id: phase1-funding
    content: Port funding rounds with negotiation and investor selection
    status: completed
  - id: phase2-partnerships
    content: Port partnership creation system
    status: completed
  - id: phase2-affiliate
    content: Port affiliate program launch/view/end
    status: completed
  - id: phase2-competitors
    content: Port competitor management
    status: completed
  - id: phase2-expansion
    content: Port global expansion
    status: completed
  - id: phase2-pivot
    content: Port pivot/strategy change
    status: completed
  - id: phase3-board
    content: Port board and equity management
    status: completed
  - id: phase3-buyback
    content: Port equity buyback
    status: completed
  - id: phase4-advanced
    content: Port advanced features (roadmap, pricing, acquisitions, etc.)
    status: completed
  - id: phase5-views
    content: Port data view screens (team, customers, financials)
    status: completed
isProject: false
---

# Port Founder Mode Features to TUI

The new TUI founder mode at `[tui/founder_game.go](tui/founder_game.go)` currently has only 3 basic options:

- Continue to next month  
- View Team
- View Exit Options

The old version at `[ui/founder_ui.go](ui/founder_ui.go)` (5400+ lines) has 25+ monthly decision options. We need full parity.

---

## Architecture

The TUI uses a `ScreenModel` pattern with `FounderView` enum for sub-views. We'll expand this with:

- **Category-based menu system** (Team, Funding, Strategic, View)
- **Sub-screens for complex actions** (hiring wizard, funding negotiation, etc.)
- **Form components for input** (amounts, selections)

```
FounderViewMain (dashboard)
    └── FounderViewActions (category menu)
            ├── FounderViewHiring
            ├── FounderViewFiring
            ├── FounderViewMarketing
            ├── FounderViewFunding
            ├── FounderViewPartnership
            ├── FounderViewAffiliate
            ├── FounderViewCompetitors
            ├── FounderViewExpansion
            ├── FounderViewBoard
            ├── FounderViewRoadmap
            ├── FounderViewPricing
            ├── FounderViewFinancials
            └── FounderViewExit
```

---

## Implementation Plan

### Phase 1: Core Gameplay Actions (Priority)

Port essential features that make the game playable:

1. **Hiring System** - Port `handleHiring()` (lines 1004-1124)
  - Engineer, Sales, CS, Marketing ($100k each)
  - C-Suite: CTO, CGO, COO, CFO ($300k each, 3x impact)
  - Market assignment for global expansion
2. **Firing System** - Port `handleFiring()` (lines 1125-1199)
  - List employees by department
  - Select and confirm termination
3. **Marketing Spend** - Port `handleMarketing()` (lines 1200-1246)
  - Choose spend amount ($10k-$200k)
  - Apply growth boost
4. **Funding Rounds** - Port `handleFundraising()` (lines 1247-1802)
  - Round types: Seed, Series A/B/C
  - Valuation negotiation
  - Term sheet options (board seat, pro-rata, etc.)
  - Investor selection

### Phase 2: Strategic Options

1. **Partnerships** - Port `handlePartnership()` (lines 1803-1851)
  - Distribution, technology, co-marketing, data partnerships
2. **Affiliate Program** - Port `handleAffiliateLaunch()` (lines 1852-1886)
  - Launch/view/end affiliate program
3. **Competitor Management** - Port `handleCompetitorManagement()` (lines 1887-1977)
  - View competitors, respond to threats
4. **Global Expansion** - Port `handleGlobalExpansion()` (lines 1978-2107)
  - Expand to EU, APAC, LATAM, etc.
5. **Pivot/Strategy** - Port `handlePivot()` (lines 2108-2194)
  - Market pivot, product pivot, business model changes

### Phase 3: Board & Equity Management

1. **Board Management** - Port `handleBoardAndEquity()` (lines 2270-2437)
  - View cap table
    - Add/remove advisors
    - Set/remove chairman
    - Manage equity pool
2. **Buyback** - Port `handleBuyback()` (lines 2195-2269)
  - Buy back investor equity when profitable

### Phase 4: Advanced Strategic Features

1. **Product Roadmap** - Port `handleProductRoadmap()` (lines 4145-4417)
2. **Customer Segments** - Port `handleCustomerSegments()` (lines 4418-4642)
3. **Pricing Strategy** - Port `handlePricingStrategy()` (lines 4643-4975)
4. **Acquisitions** - Port `handleAcquisitions()` (lines 5134-5182)
5. **Platform Strategy** - Port `handlePlatformStrategy()` (lines 5183-5228)
6. **Security & Compliance** - Port `handleSecurityCompliance()` (lines 5229-5276)
7. **PR Crisis Management** - Port `handlePRCrisis()` (lines 5277-5323)
8. **Economic Strategy** - Port `handleEconomicStrategy()` (lines 5324-5368)
9. **Succession Planning** - Port `handleSuccessionPlanning()` (lines 5369-5452)

### Phase 5: View/Data Screens

1. **View Team Roster** - Port `handleViewTeamRoster()` (lines 3905-4005)
2. **View Customer Deals** - Port `handleViewCustomerDeals()` (lines 3571-3620)
3. **View Sales Pipeline** - Port `handleViewSalesPipeline()` (lines 4976-5133)
4. **View Financials** - Port `handleViewFinancials()` (lines 3799-3904)
5. **Solicit Feedback** - Port `handleSolicitFeedback()` (lines 3621-3658)
6. **Exit Options** - Port `handleExitOptions()` (lines 3060-3570)

---

## Testing Strategy

After each phase:

1. Run the TUI and navigate to founder mode
2. Verify each new action appears in the menu
3. Test the action executes correctly
4. Verify state changes are reflected in the dashboard

---

## Key Files to Modify

- `[tui/founder_game.go](tui/founder_game.go)` - Main game screen (expand from ~350 lines to ~2000+ lines)
- May split into multiple files:
  - `tui/founder_hiring.go` - Hiring/firing logic
  - `tui/founder_funding.go` - Funding rounds
  - `tui/founder_strategic.go` - Partnerships, expansion, etc.
  - `tui/founder_views.go` - Data view screens


# Changelog

## Version 3.11.0 - Founder Mode: Enhanced Fundraising, Exit Details & Market Growth (2025-11-02)

### Major Features Added

#### 💼 Executive Team Gets Real Names
- **Famous Silicon Valley Names** for C-suite hires drawn from the show and real tech world:
  - **CTO**: Gilfoyle, Steve Wozniak, Sergey Brin, Marc Andreessen, Brendan Eich
  - **CFO**: Jared Dunn, Ruth Porat, David Wehner, Ned Segal, Luca Maestri
  - **COO**: Sheryl Sandberg, Gwart, Tim Cook, Jeff Weiner, Stephanie McMahon
  - **CGO**: Richard Hendricks, Erlich Bachman, Andrew Chen, Alex Schultz, Sean Ellis
- **Impact**: "Hire a CTO" → "Hire Gilfoyle as CTO" - gives personality to your leadership team
- Names displayed in team roster, cap table, and acquisition payouts

#### 💰 Complete Acquisition Payout Breakdown
- **Full Cap Table Payout** displayed when receiving acquisition offers:
  - Your founder payout with equity %
  - Each investor's payout (split by name, not just "Series A Investors")
  - Executive team payouts (with their names)
  - Employee equity payouts
  - Unallocated employee pool shown
- **Example**:
  ```
  $50M Acquisition Offer:
  You (Founder)           45.2%    $22,600,000
  Sequoia Capital          8.5%     $4,250,000
  Y Combinator            3.2%     $1,600,000
  Gilfoyle (CTO)          7.1%     $3,550,000
  Employee Pool           5.0%     (unallocated)
  ```
- **Impact**: See exactly how the exit splits among all stakeholders

#### 🏦 Real Investor & Firm Names in Funding Rounds
- **Realistic Investors** based on round type and amount:
  - **Angel/Pre-Seed**: Naval Ravikant, Balaji Srinivasan, Jason Calacanis, David Sacks, etc.
  - **Seed**: Y Combinator, Sequoia Scout, First Round, SV Angel, Hustle Fund, etc.
  - **Series A**: Sequoia Capital, Andreessen Horowitz, Accel, Benchmark, Greylock, etc.
  - **Series B+**: Tiger Global, SoftBank Vision Fund, Coatue, DST Global, IVP, etc.
  - **Strategic**: Bezos Expeditions, Schmidt Futures, Gates Ventures, Cuban Companies, etc.
- **Multiple Investors** per round based on amount (larger rounds = more co-investors)
- **Displayed Everywhere**: Funding history, final score screen, cap table
- **Example**: "Series A: $12M from Sequoia Capital, Accel Partners, Bezos Expeditions"

#### 💡 Advisor Impact on Fundraising Terms
- **Advisors Help Fundraise**: If you have a fundraising or strategy advisor on your board:
  - Message appears: "💡 [Advisor Name] helped improve these terms!"
  - Shows their contribution when viewing term sheet options
- **Visual Recognition**: Advisors' value is now clearly demonstrated during critical decisions
- **Impact**: Reminds founders why they gave away equity for advisory support

#### 🌍 Smart Market Expansion Lists
- **Markets Disappear After Expansion**: Once you expand to Europe, it's removed from the available markets list
- **Active Markets Shown**: See all markets you're operating in with performance metrics
- **Clear End State**: Message when you've expanded to all 6 markets (Europe, Asia, LATAM, Middle East, Africa, Australia)
- **Impact**: No more confusion about which markets are available vs. active

#### 🚀 Dramatically Improved International Market Growth
- **Sales/Marketing Actually Works** in new markets:
  - **Sales Team**: Each rep adds 5% growth rate + 2 direct customers/month
  - **Marketing Team**: Each marketer adds 3% growth rate + 1 customer/month
  - **CGO Impact**: Chief Growth Officer adds 5% rate + 3 customers/month (with 3x multiplier)
- **Two-Pronged Growth Formula**:
  - **Percentage Growth**: Compounds over time (scales with existing base)
  - **Absolute Growth**: Helps new/small markets grow from scratch
- **Competition Rebalanced**:
  - Low competition: 1.2x multiplier (was 1.1x)
  - Very high competition: 0.6x multiplier (was 0.5x)
- **Example Impact**:
  - **Before**: Expand to Asia with 25 customers, gain 1-2 customers/month
  - **After**: With 3 sales reps + CGO, gain 15-20+ customers/month in same market
- **Result**: Markets actually become viable revenue centers with proper team investment

#### 💬 Customer Feedback Reduces Churn
- **Churn Reduction**: Soliciting customer feedback now reduces churn by **3-10%** (random)
- **Displayed Clearly**: UI shows both product maturity improvement AND churn reduction
- **Example**: "Churn rate reduced by 6.2% (now 4.1%)"
- **Strategic Value**: Makes customer feedback a powerful retention tool
- **Minimum Floor**: Churn can't drop below 1% (realistic minimum)

### Technical Changes

#### Modified Files
- `founder/founder.go`:
  - Added `GenerateInvestorNames()` function with 50+ real investors/firms categorized by stage
  - Updated `HireEmployee()` to assign executive names from curated lists
  - Enhanced `RaiseFundingWithTerms()` to generate investor names and split cap table entries by investor
  - Modified `SolicitCustomerFeedback()` to reduce churn by 3-10%
  - Updated `FundingRound` struct to include `Investors []string` field
  - Rewrote `UpdateGlobalMarkets()` growth formula with dual percentage + absolute growth
  - Added CGO-specific growth bonuses for international expansion
  - Improved competition multipliers and growth rates for markets

- `founder_ui.go`:
  - Updated `displayAcquisitionOffer()` to show complete cap table breakdown with all stakeholders
  - Modified `handleFundraising()` to detect and display fundraising advisor contributions
  - Enhanced `handleSolicitFeedback()` to display churn reduction amount
  - Rewrote `handleGlobalExpansion()` to filter out active markets and show dynamic menu
  - Added market availability checking and active market display
  - Updated funding history display to show investor names
  - Modified final score display to include investor details

---

## Version 3.10.1 - Critical Bug Fix: Follow-On Investment Equity (2025-11-02)

### Bug Fixes

#### 🐛 Fixed Follow-On Investment Equity Calculation
- **CRITICAL FIX**: Follow-on investments were calculating equity based on current valuation instead of post-money valuation
- **Impact**: This caused equity percentages to exceed 100% (impossible scenario)
- **Example of Bug**:
  - User invests $500k in follow-on round
  - Equity incorrectly jumped from 52% to 112% (based on old $830k valuation)
  - Should have calculated based on new $3.4M post-money valuation
- **Fix**: Now correctly uses post-money valuation from the funding round event
- **Result**: Realistic equity percentages that properly reflect your stake in the new round

### Technical Changes

#### Modified Files
- `game/game.go`:
  - Updated `MakeFollowOnInvestment()` to look up funding round event
  - Calculate additional equity using post-money valuation (pre-money + raise amount)
  - Added validation to ensure funding round exists for the turn
  - Formula: `additionalEquity = (investment / postMoneyValuation) * 100`

---

## Version 3.10.0 - VC Mode: Investment Terms & Dramatic Events (2025-11-02)

### Major Features Added

#### 💼 Investment Terms System
- **Three Term Sheet Options** for investments $50k+:
  - **Preferred Stock** (VC Standard): Pro-rata rights, information rights, board seat, 1x liquidation preference, anti-dilution protection
  - **SAFE** (Simple Agreement for Future Equity): 20% conversion discount, pro-rata rights, simpler structure
  - **Common Stock** (Founder-friendly): Basic equity ownership with no special protections
- **Automated Selection** for smaller investments (<$50k) defaults to common stock
- **Interactive Term Selection UI** with clear explanations of each term's benefits
- **SAFE Conversion Discount** properly applies 20% bonus equity

#### 🎭 Dramatic Events System
- **10 Event Types** inspired by Silicon Valley and real startup scandals:
  - 💔 Co-founder Splits: Falling outs, CEO resignations, board chaos
  - 🔥 Scandals: Harassment allegations, PR disasters, questionable practices
  - ⚖️ Lawsuits: Patent infringement, class-action suits, legal liabilities
  - 🚨 Fraud: Financial irregularities, CFO cooking books, SEC investigations
  - 🔓 Data Breaches: Customer data leaks, GDPR fines, security incidents
  - 👋 Key Hires Quit: CTO/VP departures, team exodus to competitors
  - 📋 Regulatory Issues: Compliance problems, business model threats
  - 🔄 Pivot Failures: Failed strategic changes, lost customers
  - ⚔️ Competitor Attacks: Predatory pricing, market share erosion
  - 💥 Product Failures: Launch flops, buggy releases, refund demands
- **Difficulty-Scaled Impact**:
  - Easy: 10% event frequency, 15-25% valuation drops
  - Medium: 20% frequency, 15-50% drops
  - Hard/Expert: 30-40% frequency, 15-60% drops
- **Severity Levels**: Minor, Moderate, Severe with appropriate messaging
- **Critical Event Pausing**: Game pauses on dramatic events even in automated mode

#### 📚 Comprehensive Investing FAQ
- **Investment Terms Guide**: Preferred vs Common stock, SAFEs, Pro-Rata rights explained
- **Valuation & Equity**: How ownership is calculated, dilution mechanics
- **Funding Stages**: Pre-seed through IPO timeline and amounts
- **Exit Strategies**: Acquisitions, IPOs, secondary sales
- **Risk Management**: Diversification strategies, what kills startups
- **Key Metrics**: MRR, burn rate, CAC, LTV, and ratios
- **Accessible from Main Menu**: Help → Startup Investing FAQ

#### 🎮 Enhanced Investment UX
- **Press Enter to Skip**: Can now press Enter (or type 0) to skip investments
- **Auto-Start at $0**: Game automatically starts when investment capital depleted
- **Follow-On Skip Improvements**: Enter key skips follow-on investments
- **Exit Event Notifications**: Game always pauses for acquisitions/IPOs, even in auto mode
- **Special Exit Alert**: "🎉 COMPANY EXIT EVENT! 🎉" notification

### Technical Changes

#### New Types & Structures
- `InvestmentTerms` struct with all term details
- `DramaticEvent` struct for scandal/crisis events
- `DramaticEventQueue` in GameState

#### Modified Files
- `game/game.go`:
  - Added `GenerateTermOptions()` for term sheet generation
  - Added `MakeInvestmentWithTerms()` for term-based investing
  - Added `ScheduleDramaticEvents()` and `ProcessDramaticEvents()`
  - Integrated dramatic event processing into game loop
  - Difficulty-based event frequency and severity
- `main.go`:
  - Added `selectInvestmentTerms()` interactive UI
  - Added `displayInvestingFAQ()` comprehensive FAQ
  - Enhanced `investmentPhase()` with term selection and auto-start
  - Updated `handleFollowOnOpportunities()` with Enter-to-skip
  - Added dramatic event detection to critical message handling
  - Enhanced `playTurn()` to always pause on exits
  - Updated `displayHelpGuide()` with submenu for FAQ

#### Investment Flow Changes
- Investments $50k+ show term selection UI
- Smaller investments auto-default to common stock
- SAFE discount properly applied to equity calculation
- Terms stored with each investment for future reference

---

## Version 3.7.1 - Welcome Back Player Stats (2025-11-02)

### Major Features Added

#### 🎉 Welcome Back Feature
- **Returning Player Recognition** - Game now checks if player has played before when entering their name
- **Player Stats Display** - Shows key stats for returning players:
  - Total games played
  - Best net worth achieved
  - Best ROI percentage
  - Total successful exits
  - Average net worth across all games
  - Win rate (% of games with positive ROI)
  - Total achievements unlocked
- **Seamless Experience** - Stats display automatically before game mode selection
- **Press to Continue** - Lets players review their stats before continuing

### Technical Changes

#### Modified Files
- `main.go` - Enhanced `initMenu()` to check player stats and display welcome back message
- Added `formatCurrency()` helper function for clean currency display

---

## Version 3.7.0 - Founder Mode: Enhanced Metrics & Equity Fixes (2025-01-XX)

### Major Features Added

#### 📊 Enhanced Customer & Revenue Tracking
- **Separate Affiliate Tracking** - Direct customers and MRR tracked separately from affiliate customers/MRR
- **Variable Deal Sizes** - Realistic pricing variation (50-200% of average deal size)
  - 70% of deals within ±30% of average
  - 30% are wider range (smaller or enterprise deals)
- **Deal Size Range Display** - Shows min/max deal sizes alongside average
- **Affiliate Breakdown** - Clear visibility into revenue sources (direct vs affiliate)

#### 💼 Equity Pool Fixes
- **Proper Dilution** - Expanding equity pool now correctly dilutes founder equity
- **Accurate Calculations** - Founder equity = 100% - EquityGivenAway - EquityPool
- **Consistent Display** - Equity calculations updated across all screens

### Technical Changes

#### Modified Files
- `founder/founder.go` - Added separate tracking fields, variable deal size generation, fixed equity calculations
- `founder_ui.go` - Enhanced metrics display with breakdowns and min/max ranges

#### New Fields
- `DirectMRR`, `AffiliateMRR` - Separate MRR tracking
- `DirectCustomers`, `AffiliateCustomers` - Separate customer tracking
- `MinDealSize`, `MaxDealSize` - Deal size range tracking

#### New Functions
- `generateDealSize()` - Creates realistic variable deal sizes
- `updateDealSizeRange()` - Tracks min/max deal sizes

### Gameplay Impact
- **More Realistic** - Deal sizes vary like real SaaS pricing (tiers, discounts, enterprise deals)
- **Better Visibility** - Clear breakdown of revenue sources and customer acquisition channels
- **Proper Dilution** - Equity pool expansion correctly affects founder ownership
- **Accurate Metrics** - Can now see exactly where revenue is coming from

---

## Version 3.6.0 - Founder Mode: Realistic Financials & Cap Table (2025-11-01)

### Major Features Added

#### 💰 MRR Cash Flow System
- **Realistic Revenue Conversion** - MRR now flows into company cash after realistic deductions
- **Tax Deductions** - 20% tax rate on revenue
- **Processing Fees** - 3% payment processing fees
- **Company Overhead** - 5% operational overhead costs
- **Savings Buffer** - 5% set aside for reserves
- **Net Cash Flow** - ~67% of MRR converts to cash (realistic for SaaS businesses)

#### 📊 Cap Table Management
- **Equity Tracking** - Complete cap table structure tracking all equity holders
- **Employee Equity Grants** - Initial employees receive 1-2% equity each at game start
- **Executive Equity Grants** - C-suite executives (CTO, CFO, COO, CGO) receive 3-10% equity when hired
- **Investor Tracking** - All funding rounds tracked on cap table with equity percentages
- **Cap Table Structure** - Detailed tracking of employee, executive, and investor equity

#### 🎲 Startup Randomization
- **Cash Variance** - Initial cash varies ±20% (0.8x to 1.2x multiplier)
- **Competition Randomization** - 30% chance to randomize competition level (low/medium/high/very_high)
- **More Variety** - Each game start has unique financial and market conditions

### Technical Changes

#### Modified Files
- `founder/founder.go` - Added MRR cash flow, cap table structures, equity grants, and randomization

#### New Structures
- `CapTableEntry` - Tracks individual equity ownership (name, type, equity %, month granted)
- `Employee.Equity` - Added equity percentage field to Employee struct
- `FounderState.CapTable` - Array tracking all equity holders

### Gameplay Impact
- **More Realistic** - MRR now properly affects cash runway and company financials
- **Strategic Equity** - Players must manage equity pool when hiring executives
- **Varied Starts** - Each game provides different initial conditions for replayability
- **Complete Ownership** - Full visibility into who owns what percentage of the company

---

## Version 2.0.0 - Global Leaderboard & Cloud Integration (2025-01-XX)

### Major Features Added

#### 🌐 Global Leaderboard System
- **Real-time Global Competition** - Players can submit scores and compete worldwide
- **One-Click Score Submission** - Submit scores directly from the game
- **Multiple Views** - Filter by difficulty, sort by Net Worth or ROI
- **Real-time Updates** - Leaderboard updates automatically
- **Beautiful UI** - Modern, responsive design on GitHub Pages
- **Privacy First** - Only game stats are submitted, no personal data
- **Free Hosting** - Powered by Vercel (free tier) and Datasette

#### ☁️ Cloud Infrastructure
- **Vercel Serverless API** - Go-based API endpoint for score submissions (`/api/submit-score`)
- **Datasette Integration** - Read-only JSON API for leaderboard data
- **GitHub Pages Display** - Live leaderboard on GitHub Pages
- **Automatic Deployment** - GitHub Actions workflow for auto-redeployment

#### 📊 New Components
- **API Endpoint** (`api/submit-score.go`) - Serverless function for score submission
- **Leaderboard Client** (`leaderboard/leaderboard.go`) - HTTP client for game integration
- **Database Schema** - SQLite database with UUID-based score tracking
- **Frontend Integration** - JavaScript fetches from Datasette API

### Technical Changes

#### New Files
- `api/submit-score.go` - Serverless function for score submission
- `api/go.mod` - Dependencies for API
- `leaderboard/leaderboard.go` - HTTP client for game
- `datasette-metadata.json` - Datasette configuration
- `vercel.json` - Vercel deployment config
- `leaderboard.db` - SQLite database with sample data
- `DATASETTE_SETUP.md` - Detailed setup guide
- `QUICKSTART_LEADERBOARD.md` - Quick start guide
- `README_LEADERBOARD.md` - Architecture documentation
- `DEPLOYMENT_STATUS.md` - Deployment status and next steps

#### Modified Files
- `main.go` - Added score submission prompt after game ends
- `docs/index.html` - Added global leaderboard display section
- `.gitignore` - Added Vercel and DB files

#### Database Schema
```sql
CREATE TABLE game_scores (
  id TEXT PRIMARY KEY,              -- UUID
  player_name TEXT NOT NULL,
  final_net_worth INTEGER NOT NULL,
  roi REAL NOT NULL,
  successful_exits INTEGER NOT NULL,
  turns_played INTEGER NOT NULL,
  difficulty TEXT NOT NULL,
  played_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Deployment
- **API Deployed:** Vercel serverless function
- **Datasette Deployed:** Vercel-hosted database API
- **Frontend:** GitHub Pages with live leaderboard
- **Status:** Ready for production use

### API Endpoints
- `POST /api/submit-score` - Submit game score
- `GET /leaderboard/game_scores.json` - Fetch leaderboard data

### Next Steps
- Disable Vercel authentication protection for public access
- Configure production domains
- Set up GitHub Actions for auto-deployment

---

## Phase 4 - Achievements & Career Progression (2025-10-31)

### Major Features Added

#### ?? Achievements System
- **35+ Achievements** across 6 categories:
  - Wealth (5): First Profit, Millionaire, Multi-Millionaire, Deca-Millionaire, Mega Rich
  - Performance (5): Break Even, Double Up, Great Investor, Elite VC, Unicorn Hunter
  - Strategy (6): Diversified, Sector Master, All In, Sector Specialist, Exit Master, Perfect Portfolio
  - Career (6): First Steps, Persistent, Veteran, Master Investor, Hot Streak, On Fire
  - Challenge (7): Easy Money, Rising Star, Battle Tested, Expert Survivor, Easy Domination, Expert Legend, Speed Runner
  - Special (6+): Lucky Seven, Minimalist, Tech Enthusiast, Clean Investor, Risk Taker, Cautious Investor
- **Rarity System:** Common, Rare, Epic, Legendary
- **Point Values:** 5-100 points per achievement
- **Hidden Achievements:** Discover special achievements through gameplay

#### ??? Career Progression
- **11 Career Levels:**
  - Level 0: Intern (0 pts)
  - Level 1: Analyst (25 pts)
  - Level 2: Associate (75 pts)
  - Level 3: Senior Associate (150 pts)
  - Level 4: Principal (250 pts)
  - Level 5: Partner (400 pts)
  - Level 6: Senior Partner (600 pts)
  - Level 7: Managing Partner (850 pts)
  - Level 8: Elite VC (1150 pts)
  - Level 9: Master Investor (1500 pts)
  - Level 10: Legendary Investor (2000 pts)
- **Point-Based Progression:** Accumulate points to level up
- **Progress Tracking:** See how many points needed for next level

#### ?? Database Integration
- **New Table:** `player_achievements`
  - Stores player name, achievement ID, unlock timestamp
  - Unique constraint prevents duplicates
  - Indexed for fast lookups
- **New Functions:**
  - `UnlockAchievement()` - Save achievement
  - `GetPlayerAchievements()` - Retrieve all unlocked
  - `GetPlayerAchievementCount()` - Count achievements
  - `GetWinStreak()` - Track consecutive wins

#### ?? Real-time Notifications
- **After-Game Display:**
  - Banner announcing new achievements
  - Icon, name, rarity for each
  - Description of accomplishment
  - Points earned
  - Updated career level display
  - Total achievement points
- **Color-Coded:**
  - White (Common), Cyan (Rare), Magenta (Epic), Yellow (Legendary)

#### ?? Achievements Menu
- **Main Menu Option 4:** Achievements
- **Sub-Menu:**
  - View My Achievements - Personal progress
  - Browse All Achievements - See what's possible
  - Leaderboard (Coming Soon) - Top achievers
- **Personal View:**
  - Progress percentage (X/35)
  - Total points earned
  - Current career level & title
  - Points to next level
  - Grouped by category

### Technical Changes

#### New Files
- `achievements/achievements.go` - Complete achievements system (564 lines)
- `PHASE4_SUMMARY.md` - Phase 4 documentation

#### Modified Files
- `database/database.go` - Achievement tracking functions
- `main.go` - Achievements menu, notifications, checking logic
- `README.md` - Updated with Phase 4 features
- `CHANGELOG.md` - This update

#### New Dependencies
None - Pure Go implementation

### Achievement Stats
- **Total Achievements:** 35+
- **Categories:** 6
- **Rarity Levels:** 4
- **Point Range:** 5-100
- **Career Levels:** 11
- **Max Career Points:** 2000+

### Game Flow Updates
```
Before:
Play ? Score ? Repeat

After:
Play ? Score ? ?? Achievements! ? Level Up ? Repeat
```

---

## Phase 3 - Content Expansion & Analytics (2025-10-31)

### Major Features Added

#### ?? Content Doubled
- **20 Startup Companies** (up from 10)
- 10 new companies across new sectors:
  - NanoMed (BioTech), CloudForge (CloudTech)
  - FoodLoop (CleanTech), GameStream (Gaming)
  - LegalAI (LegalTech), UrbanFarm (AgriTech)
  - DroneDeliver (Logistics), MusicGen (Creative)
  - SmartHome (IoT), QuantumSecure (Security)
- **60 Random Events** (up from 30)
- 30 new events covering:
  - Series B funding rounds
  - Product recalls and breakthroughs
  - Strategic alliances
  - Security vulnerabilities
  - Market disruptions
  - Regulatory changes

#### ?? Advanced Analytics System
- **Portfolio Analytics Package** (`analytics/analytics.go`)
- Comprehensive performance tracking:
  - Total invested vs. current value
  - Best/worst performer identification
  - Positive/negative investment ratio
- **Sector Performance Breakdown**
  - ROI by industry sector
  - Investment distribution
  - Sector rankings
  - Per-sector statistics
- Automatically displayed after each game

#### ?? Help & Information Menu
- Complete in-game guide
- Game overview and mechanics
- Company metrics explanation
- Scoring system details
- Difficulty breakdowns
- Strategy tips
- All companies listed
- Event categories explained

### Technical Changes

#### New Files
- `startups/11.json` through `startups/20.json` - 10 new companies
- `analytics/analytics.go` - Analytics engine (350+ lines)
- `PHASE3_SUMMARY.md` - Phase 3 documentation

#### Modified Files
- `game/game.go` - Load 20 companies (up from 10)
- `rounds/round-options.json` - 60 events (up from 30)
- `main.go` - Added help menu option and function
- `README.md` - Updated feature list
- `CHANGELOG.md` - This update

### Content Stats
- **Companies:** 10 ? 20 (100% increase)
- **Events:** 30 ? 60 (100% increase)
- **Sectors:** 8 ? 12+ (50% increase)
- **Valuation Range:** $8M - $45M
- **Event Impact Range:** 0.5x - 3.2x

---

## Phase 2 - Persistence & Competition (2025-10-31)

### Major Features Added

#### ??? Database Persistence
- SQLite database integration for score storage
- Automatic database creation on first run
- All games automatically saved with full details
- Database file: `unicorn_scores.db`

#### ?? Leaderboard System
- 7 different leaderboard views
- Top 10 by Net Worth (all difficulties)
- Top 10 by ROI (all difficulties)
- Separate boards for Easy/Medium/Hard/Expert
- Recent games view (last 10)
- Color-coded rankings (gold/silver/bronze)

#### ?? Player Statistics
- Career stats tracking per player
- Total games played
- Best net worth achieved
- Best ROI percentage
- Total successful exits across all games
- Average net worth
- Win rate (% positive ROI)

#### ?? Difficulty Levels
- **Easy Mode:** $500k start, 20% events, 3% volatility
- **Medium Mode:** $250k start, 30% events, 5% volatility
- **Hard Mode:** $150k start, 40% events, 7% volatility
- **Expert Mode:** $100k start, 50% events, 10% volatility, only 90 turns!

#### ?? Enhanced UI
- Main menu system (New Game, Leaderboards, Stats, Quit)
- Difficulty selection screen
- Score save confirmation
- Formatted leaderboard tables
- Player stats display
- Color-coded performance indicators

### Technical Changes

#### New Files
- `database/database.go` - Complete persistence layer
- `PHASE2_SUMMARY.md` - Phase 2 documentation
- `CHANGELOG.md` - This file

#### Modified Files
- `game/game.go` - Added Difficulty system
- `main.go` - Complete rewrite with menu system
- `README.md` - Updated features list
- `QUICKSTART.md` - Updated with new features
- `go.mod` - Added sqlite3 dependency

#### API Changes
- `game.NewGame()` now takes `Difficulty` parameter instead of `int64`
- Added difficulty-based event frequency
- Added difficulty-based volatility

### Dependencies
- Added: `github.com/mattn/go-sqlite3 v1.14.32`

---

## Phase 1 - Core Gameplay (2025-10-31)

### Initial Release

#### ? Core Features
- Investment mechanics with portfolio tracking
- 10 diverse startup companies
- 30 random game events
- Turn-based system (120 turns)
- Win/loss conditions
- Final scoring system
- ROI calculation
- Performance ratings (6 tiers)

#### ?? Content
- 10 startups across 8 sectors
- 30 events (positive and negative)
- Risk scoring system
- Growth potential indicators

#### ?? UI Features
- Color-coded profit/loss
- Money formatting
- Portfolio dashboard
- Event notifications
- Final results screen

### Initial Files
- `game/game.go` - Game engine
- `main.go` - UI and game loop
- `startups/1-10.json` - Company data
- `rounds/round-options.json` - Event data
- `README.md` - Documentation
- `QUICKSTART.md` - Quick start guide
- `PHASE1_SUMMARY.md` - Phase 1 docs

---

## Future Roadmap

### Phase 4 - Achievements & Polish (Planned)
- [ ] Achievement system
- [ ] Career progression
- [ ] Challenge modes
- [ ] Better analytics dashboard
- [ ] ASCII charts for trends
- [ ] Export game results

### Long-term Vision
- [ ] Multiplayer mode (hot-seat)
- [ ] Startup mode (you run a company)
- [ ] Secondary market (buy/sell shares)
- [ ] Follow-on investment rounds
- [ ] M&A activity
- [ ] Syndicate with other VCs

---

## Build Information

**Current Version:** Phase 4
**Go Version:** 1.x+
**Platform:** Linux/macOS/Windows
**Database:** SQLite 3
**Build Size:** ~7MB
**Companies:** 20
**Events:** 60
**Sectors:** 12+
**Achievements:** 35+
**Career Levels:** 11

## How to Build

```bash
go mod download
go build -o unicorn
```

## Dependencies

```
github.com/fatih/color v1.7.0
github.com/buger/goterm v0.0.0-20181115115552-c206103e1f37
github.com/mattn/go-sqlite3 v1.14.32
gopkg.in/yaml.v2 v2.2.8
```

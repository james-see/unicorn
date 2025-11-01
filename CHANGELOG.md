# Changelog

## Phase 5 - Datasette Global Leaderboard (2025-11-01)

### Major Features Added

- **Global Leaderboard Publishing:** Players can opt into submitting their final results to a Datasette-backed leaderboard after each game.
- **GitHub Pages Integration:** The project website now renders the live global leaderboard by consuming Datasette?s JSON API.
- **Configurable Deployment:** YAML/JSON configs make it easy to point both the game and website at your Datasette instance.

### Technical Changes

- Added `datasette/datasette.go` for authenticated inserts via the `datasette-write` API.
- Introduced `config/datasette.yaml` for runtime configuration with optional `UNICORN_DATASETTE_TOKEN` override.
- Updated `main.go` to load Datasette settings, prompt players, and submit global scores.
- Refreshed `docs/index.html` with a live leaderboard card, styles, and a Datasette-powered loader.
- Added `docs/leaderboard-config.json` to configure the Datasette endpoint used by GitHub Pages.
- Documented the workflow in `README.md`.

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

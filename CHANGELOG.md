# Changelog

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

### Phase 3 - Content Expansion (Planned)
- [ ] Expand to 20+ startup companies
- [ ] Add 50+ diverse events
- [ ] Company metrics visualization
- [ ] Sector performance tracking
- [ ] Economic cycle indicators
- [ ] Enhanced risk/growth calculations

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

**Current Version:** Phase 2
**Go Version:** 1.x+
**Platform:** Linux/macOS/Windows
**Database:** SQLite 3
**Build Size:** ~7MB

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

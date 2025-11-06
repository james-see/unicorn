# Changelog

## Version 3.20.1 - Bug Fix: XP Integration & Documentation Update (2025-11-06)

### Bug Fixes

#### 🐛 XP System Integration
- **Issue**: XP system was not being called after games completed - players didn't see XP earned
- **Fix**: 
  - Integrated XP calculation and display into both VC Mode and Founder Mode completion flows
  - Added `DisplayXPGained()` function calls to show detailed XP breakdown
  - Added level-up celebration screens with `DisplayLevelUp()` function
  - XP now properly awarded for: game completion (+100), positive ROI (+50), successful exits (+200), difficulty bonuses (0-200), achievements (+10 × points)
  - Founder mode bonuses: IPO (+500), acquisition (+300), profitability (+100)

#### 🎮 Upgrade Filtering by Game Mode
- **Issue**: All upgrades shown regardless of game mode (VC upgrades in Founder, Founder upgrades in VC)
- **Fix**:
  - Added `GetUpgradesForGameMode()` and `FilterUpgradeIDsForGameMode()` functions
  - VC Mode now shows only: Investment Terms, Financial Perks, Information, Board Powers, Special Abilities, Game Modes
  - Founder Mode now shows only: Founder Perks and Game Modes
  - Upgraded displays now labeled "ACTIVE UPGRADES FOR THIS GAME"

### Documentation Updates

#### 📚 README.md (GitHub)
- Added comprehensive Player Progression System section (v3.20.0)
- Updated "What's New" with v3.20.0 as latest release
- Updated feature list: 45+ achievements, 50 levels, 30 startups
- Updated main menu documentation to show all 9 menu options
- Added detailed XP sources and level unlock information

#### 🎮 In-Game Help (ui/help_ui.go)
- Updated game overview to highlight both VC and Founder modes
- Added Player Progression System section with XP breakdown
- Added Achievements & Upgrades section
- Added Analytics Dashboard section
- Updated strategy tips to include progression advice

#### 🌐 GitHub Pages (docs/index.html)
- **Fixed spacing issue**: "Two Game Modes" line now has proper margin (0.5rem top, 1.5rem bottom)
- **Updated header stats**: 50 Levels, 45+ Achievements, 30 Startups
- **Added new section**: Player Progression System (v3.20.0) with 6 feature cards
- **Updated VC Mode features**: Upgraded from 20 to 30 startups, mentioned upgrades system
- **Mobile-responsive leaderboards**: 
  - Horizontal scroll support
  - Responsive font sizing (0.85rem tablets, 0.75rem phones)
  - Hide less important columns on mobile (Exits, Difficulty, Date)
  - Optimized for 768px (tablet) and 480px (phone) breakpoints
  - Table wrapper with min-width for proper display

### Technical Changes

**Modified Files**:
- `ui/achievements_ui.go` - Added XP integration to VC mode completion
- `ui/founder_ui.go` - Added XP integration to Founder mode completion
- `ui/vc_ui.go` - Added upgrade filtering for VC mode welcome screen
- `upgrades/upgrades.go` - Added game mode filtering functions
- `README.md` - Comprehensive documentation update
- `ui/help_ui.go` - Enhanced in-game help with progression details
- `docs/index.html` - Landing page updates and mobile responsiveness

**New Functions**:
- `GetUpgradesForGameMode(gameMode string)` - Returns upgrades relevant to specific mode
- `FilterUpgradeIDsForGameMode(upgradeIDs []string, gameMode string)` - Filters owned upgrades by mode

### User Experience Improvements

**Example Impact**:
- Players now see exactly how much XP they earned after every game
- Level-up celebrations make progression feel rewarding
- Upgrade lists are no longer cluttered with irrelevant items
- Mobile players can view leaderboards comfortably
- Complete documentation helps new players understand all features

**XP Display Example**:
```
📊 EXPERIENCE EARNED:
   +100 XP - Game Completion
   +50 XP - Positive ROI
   +200 XP - Successful Exit
   +50 XP - Medium Difficulty
   +30 XP - New Achievements (1)
   
   Total XP Gained: +430 XP

   Level 3 Progress: [■■■■■■□□□□] 650/1000 XP
```

---

## Version 3.20.0 - Major Feature Update: Progression, Analytics & Enhanced Achievements (2025-11-06)

### Major Features Added

#### 💎 Player Progression System
- **Level-Based Progression**: 50 levels with exponential XP requirements (200 * level^1.5)
- **Experience Points**: Earn XP from games, achievements, difficulty bonuses, and successful exits
- **Level Unlocks**: Progressive unlocking of difficulties (Hard at L5, Expert at L10, Nightmare at L15)
- **Rank Titles**: From "Novice Investor" to "Titan of Industry"
- **Visual Progress**: XP bars, level-up celebrations, and progression tracking in main menu
- **Persistent Profiles**: Player profiles with level history and total XP earned

#### 🏆 Achievement System Enhancements
- **Achievement Chains**: Connected achievements with prerequisites (diversification, win streaks, investment count)
- **Progress Tracking**: 15+ progressive achievements with visual progress bars (e.g., "Win 10 games: 4/10")
- **Hidden Achievements**: Mystery achievements revealed only upon unlock (Phoenix, Perfect Month, Day Trader)
- **New Achievements**: Added 18+ new achievements including:
  - Diversification Chain (3 sectors → 5 sectors → all sectors)
  - Win Streak Chain (3 → 5 → 10 consecutive wins)
  - Games Played Milestones (10 → 50 → 100 → 500 games)
  - Investment Expertise Chain (10 → 50 → 100 successful investments)
- **Enhanced UI**: Filter by chains, view progress, see locked/available achievements

#### 📊 Analytics Dashboard
- **Performance Trends**: Automatic trend analysis for 7-day, 30-day, and all-time periods
- **Difficulty Breakdown**: Per-difficulty statistics (games, win rate, avg net worth)
- **Historical Performance**: Monthly reports showing games, wins, ROI, and best results
- **Top Games Tracking**: Personal leaderboard of your best 5 games
- **AI-Generated Insights**: Smart recommendations based on your performance patterns
- **Global Comparisons**: Compare your stats to global averages with percentile rankings
- **Visual Charts**: ASCII trend charts and performance heatmaps
- **Player Comparison**: Side-by-side comparison of two players' statistics

#### 🌍 Advanced Game Mechanics - Market Cycles
- **Dynamic Market Cycles**: 5 cycle types (Bull, Bear, Normal, Recession, Boom)
- **Economic Events**: 9+ event types including Interest Rate Hikes, Tech Booms, Credit Crunches
- **Valuation Effects**: Market cycles affect startup valuations (0.6x - 1.5x multipliers)
- **Funding Dynamics**: Market conditions impact funding availability (0.5x - 1.8x multipliers)
- **Sector-Specific Impact**: Targeted events affect specific sectors (e.g., AI Investment Frenzy)
- **Market Sentiment Indicators**: Real-time display of current market conditions

### Technical Changes

**New Files Created**:
- `progression/progression.go` - Core progression system with XP calculations
- `progression/unlocks.go` - Level-based unlock definitions
- `ui/progression_ui.go` - Progression UI displays and level-up screens
- `achievements/checker.go` - Achievement chain and progress tracking logic
- `analytics/analytics.go` - Trend analysis and performance comparison engine
- `ui/analytics_ui.go` - Comprehensive analytics dashboard UI
- `game/market.go` - Market cycle and economic event system

**Modified Files**:
- `database/database.go` - Added 5 new tables (player_profiles, player_level_history, achievement_progress, game_history_detailed)
- `achievements/achievements.go` - Extended Achievement struct with chains, progress tracking, and 18+ new achievements
- `ui/achievements_ui.go` - Enhanced with chain displays, progress bars, and hidden achievements
- `ui/main_menu.go` - Added Progression & Levels and Analytics Dashboard menu options
- `main.go` - Integrated new menu options and progression display

**Database Schema Changes**:
- `player_profiles` table - Stores player level, XP, and progression data
- `player_level_history` table - Tracks when players reached each level
- `achievement_progress` table - Progressive achievement tracking with current/max progress
- `game_history_detailed` table - Comprehensive game statistics for analytics

**New Package Functions**:
- Player progression: `GetPlayerProfile()`, `AddExperience()`, `GetLevelRequirement()`
- Achievement chains: `CheckAchievementChain()`, `GetNextInChain()`, `GetProgressiveAchievements()`
- Analytics: `GenerateTrendAnalysis()`, `CompareToGlobal()`, `GetMonthlyStats()`
- Market cycles: `InitializeMarketCycle()`, `AdvanceMarketCycle()`, `ApplyMarketEffects()`

### User Experience Improvements

**Example Impact**:
- Players now have clear progression goals with level-based unlocks
- Achievement chains provide longer-term objectives and satisfying progression
- Analytics dashboard offers deep insights into performance trends and improvement areas
- Market cycles add strategic depth and variety to each playthrough

---

## Version 3.19.2 - Code Quality Improvements (2025-11-06)

### Bug Fixes

#### 🧹 Code Formatting Cleanup
- **Issue**: Inconsistent whitespace in test files
- **Fix**: 
  - Cleaned up trailing whitespace in `founder/founder_test.go`
  - Standardized indentation and blank line usage
  - Improved code readability and consistency

### Technical Changes

#### Modified Files
- `founder/founder_test.go`:
  - Removed trailing whitespace from test functions
  - Standardized blank line formatting between test assertions
  - Improved code consistency across all test cases
- `REFACTORING_SUMMARY.md`:
  - Minor formatting adjustment

### Impact
- Better code maintainability
- Consistent code style across test suite
- No functional changes

---

## Version 3.19.1 - Animation Visibility & Achievement Display Fixes (2025-11-05)

### Bug Fixes

#### 🐛 Animations Not Visible to Users
- **Issue**: Animations were displaying but getting cleared or overwritten too quickly, making them invisible to players
- **Root Cause**: Screen clear operations and lack of pauses caused animations to flash by before users could see them
- **Fixes**:
  - Moved game over animation to display BEFORE screen clear (was clearing then showing)
  - Added manual pause with "Press Enter to see detailed results..." after game over animation
  - Added pause with "Press Enter to continue to main menu..." after splash screen animation
  - Users now have time to enjoy the animations before they're replaced

#### 🐛 Achievement Section Missing at Game End
- **Issue**: Players reported "no rewards at game end" - achievement section only showed when new achievements were unlocked
- **Root Cause**: Achievement check section was conditional and only displayed if `len(newAchievements) > 0`
- **Fix**: 
  - Always show "ACHIEVEMENT CHECK" section at game end
  - When no new achievements: Display helpful tips on how to unlock them
  - When achievements unlocked: Show animated celebration
  - Provides clear feedback either way

#### 🐛 Career Level Not Displaying
- **Issue**: Career level, lifetime points, and available balance weren't showing after games
- **Root Cause**: Career level display code was incorrectly nested inside the achievement unlock block
- **Fix**: Moved career level display outside the conditional to always show after every game

### Technical Changes

#### Modified Files
- `main.go`:
  - Reordered `displayFinalScore()` to show animation before clear
  - Added `bufio.NewReader(os.Stdin).ReadBytes('\n')` pauses after animations
  - Restructured `checkAndUnlockAchievements()` to always show achievement section
  - Added helpful tips display when no new achievements
  - Fixed indentation of career level display block (removed from inside achievement conditional)

### User Experience

#### 🎨 Animations Now Clearly Visible
- **Splash Screen**: Users see big UNICORN title, loading spinner, and success message with manual pause
- **Game Over**: Dramatic victory/defeat animation displays with pause before detailed results
- **Achievement Unlocks**: Flashy star effects and colored boxes clearly visible

#### 📊 Better Feedback & Guidance
- **Always Informed**: Achievement check section shows after every game
- **Helpful Tips**: When no achievements unlocked, game provides actionable guidance:
  - Wealth goals: Reach $1M, $5M, $10M, $50M net worth
  - Performance goals: Achieve positive ROI, 2x, 5x, 10x returns
  - Strategy goals: Diversify investments, master sectors, get successful exits
  - Career goals: Play more games, build win streaks
- **Progress Tracking**: Career level, total points, available balance always visible

#### 💡 Example Impact
**Before**: Animation flashed → Screen cleared → User saw nothing → "I didn't see any animations"
**After**: Animation displays → User sees it → "Press Enter" pause → User enjoys animation → Continues

**Before**: Game ends → No feedback if no new achievements → "No rewards at game end"
**After**: Game ends → "ACHIEVEMENT CHECK" section → Shows progress OR helpful tips → Always informed

---

## Version 3.19.0 - Animated CLI Experience (2025-11-05)

### Major Features Added

#### 🎨 Beautiful Terminal Animations
- **New Package**: Added pterm v0.12.82 for modern CLI animations and effects
- **Splash Screen**: Animated startup with big "UNICORN" title, styled info box, and loading spinner
- **Achievement Unlocks**: Flashy animations with stars (✨⭐🏆) and colored boxes when unlocking achievements
- **Round Milestones**: Animated header transitions every 5 turns showing progress
- **Game Over Screens**: Epic victory/defeat animations with big text and styled result boxes

#### 💼 New Animation Library Functions
- `ShowGameStartAnimation()` - Animated splash screen on game launch
- `ShowAchievementUnlock()` - Flashy achievement notification with effects
- `ShowRoundTransition()` - Milestone round transitions
- `ShowGameOverAnimation()` - Victory or defeat end-game animation
- `ShowLoadingSpinner()` - Customizable loading animations
- `ShowSuccessMessage()`, `ShowErrorMessage()`, `ShowWarningMessage()`, `ShowInfoMessage()` - Styled notification messages
- `ShowInvestmentAnimation()` - Investment processing animation
- `ShowExitAnimation()` - Successful exit celebration with fireworks effect
- `TypewriterEffect()` - Character-by-character text animation
- `ShowProgressBar()` - Animated progress bars

### Technical Changes

#### New Files
- `animations/animations.go` - Complete animation library with 15+ animation functions
- `ANIMATIONS.md` - Comprehensive documentation for animation system

#### Modified Files
- `main.go`:
  - Added animations package import
  - Integrated splash screen animation in `main()` function on startup
  - Added achievement unlock animations in `checkAndUnlockAchievements()`
  - Added round transition animations in `playTurn()` (every 5 turns)
  - Added game over animations in `displayFinalScore()`
  - Enhanced quit message with styled info notification
- `go.mod`:
  - Added pterm v0.12.82 as direct dependency
  - Added indirect dependencies: atomicgo.dev/cursor, atomicgo.dev/keyboard, atomicgo.dev/schedule, github.com/gookit/color, github.com/lithammer/fuzzysearch, github.com/mattn/go-runewidth, github.com/rivo/uniseg, github.com/xo/terminfo, golang.org/x/term, golang.org/x/text
  - Upgraded Go toolchain to 1.24.1

### User Experience

#### 🚀 Enhanced Visual Feedback
- **Immersive Startup**: Players are greeted with a beautiful animated splash screen showing the game loading
- **Celebration Moments**: Achievement unlocks now have satisfying visual effects with stars and animations
- **Progress Markers**: Every 5 turns shows a special milestone header to mark game progress
- **Epic Conclusions**: Win or lose, the game ends with a dramatic full-screen animation

#### 🎮 Improved Game Feel
- **Professional Polish**: The CLI now feels like a modern application with smooth animations
- **Visual Hierarchy**: Important moments (achievements, game over) are highlighted with animations
- **Pacing**: Animations are timed perfectly (0.5-2 seconds) to enhance without slowing gameplay
- **Feedback**: Loading spinners and status messages provide clear feedback during operations

#### 📊 Example Impact
Before: Plain text "Achievement Unlocked: High Roller"
After: Flashing stars → Animated box → "🏆 ACHIEVEMENT UNLOCKED 🏆" with colored borders and timed display

---

## Version 3.18.8 - Negative Equity Pool Bug Fix (2025-11-04)

### Bug Fixes

#### 🐛 Negative Equity Available Display
- **Issue**: After hiring executives, the equity pool display showed negative available equity (e.g., "9.1% used, -2.5% available")
- **Root Cause**: When executives were hired, the code reduced `EquityPool` directly but never incremented `EquityAllocated`, causing the calculation `EquityPool - EquityAllocated` to show negative values
- **Fix**: Updated executive hiring logic to increment `EquityAllocated` instead of reducing `EquityPool`, and changed the availability check to use `EquityPool - EquityAllocated`
- **Result**: Equity pool tracking now correctly shows available equity, matching the same logic used for advisors

### Technical Changes

#### Modified Files
- `founder/founder.go`:
  - Updated `HireEmployee()` function for executive roles
  - Changed availability check from `executiveEquity > fs.EquityPool` to `executiveEquity > availableEquity` where `availableEquity = fs.EquityPool - fs.EquityAllocated`
  - Changed from `fs.EquityPool -= executiveEquity` to `fs.EquityAllocated += executiveEquity`
  - Now consistent with advisor hiring logic which properly tracks allocated equity

### User Experience
- **Accurate Display**: Equity pool now always shows correct available percentage
- **Consistent Logic**: Executive and advisor equity tracking uses the same method
- **Better UX**: Users can now properly see how much equity remains available for hiring

---

## Version 3.18.7 - Copyright Year Update (2025-11-03)

### Maintenance

#### 📝 Copyright Update
- Updated copyright year from 2019 to 2025
- Updated in README.md and logo/logo.go
- Game logo and documentation now display correct copyright year

### Technical Changes

#### Modified Files
- `README.md`: Updated copyright notice in ASCII art header
- `logo/logo.go`: Updated copyright displayed in game welcome screen

---

## Version 3.18.6 - Founder Mode Menu Fixes (2025-11-03)

### Bug Fixes

#### 🐛 Menu Numbering Out of Order
- **Issue**: Menu options were numbered incorrectly (17, 12, 16, 15 appearing out of sequence)
- **Fix**: Implemented dynamic sequential numbering system
- **Result**: All menu options now numbered sequentially based on availability

#### 🐛 Invalid Choice Validation
- **Issue**: Valid options (like option 13) were being rejected as invalid
- **Fix**: Updated validation logic to properly calculate maxChoice and validate all options
- **Result**: All valid menu choices now work correctly

#### 🐛 Funding Round Numbering
- **Issue**: When raising Series A or Series B, menu showed "2. Series A" but required input "1"
- **Fix**: Implemented dynamic option numbering for funding rounds
- **Result**: If only Series A is available, it shows "1. Series A" and accepts "1" as input

### Technical Changes

#### Modified Files
- `founder_ui.go`:
  - Refactored `handleFounderDecisions()` to use dynamic sequential numbering
  - Added `nextOption` counter to track current option number
  - View options, Strategic Opportunity, and Exit options now numbered dynamically
  - Updated `handleFundraising()` to use `optionNum` counter instead of hardcoded 1, 2, 3
  - Fixed validation to properly check against calculated `maxChoice`
  - Improved option handling logic to support sequential numbering

### User Experience
- **Clear Menu**: Options are now numbered sequentially (1, 2, 3... instead of 1, 2, 3, 17, 12, 16, 15)
- **Correct Validation**: All valid choices are now accepted
- **Intuitive Funding**: Funding round selection matches displayed numbers

---

## Version 3.18.5 - Founder Mode Win Condition Fix (2025-11-03)

### Bug Fixes

#### 🐛 Win-Based Achievements Triggering on Losses
- **Issue**: Achievements like "Speed Runner" (Win in under 60 turns) were triggering when founder ran out of cash
- **Root Cause**: Win condition was based on ROI > 0, which could be true even when cash ran out (due to equity valuation)
- **Fix**: Updated founder mode win condition to properly check for successful exit OR reaching max turns without running out of cash
- **Win Conditions**:
  - **Won**: Exited successfully (IPO/Acquisition/Secondary) OR reached max turns with cash remaining
  - **Lost**: Ran out of cash before exiting or reaching max turns
- **Result**: Win-based achievements now only trigger on actual wins, not losses

### Technical Changes

#### Modified Files
- `achievements/achievements.go`:
  - Added `RanOutOfCash bool` field to `GameStats` struct
  - Updated `checkAchievement()` win logic to differentiate founder mode vs VC mode
  - Founder mode: `won = HasExited || !RanOutOfCash`
  - VC mode: `won = ROI > 0` (unchanged)

- `founder_ui.go`:
  - Set `RanOutOfCash` field when building `GameStats` for achievement checking
  - `RanOutOfCash = fs.Cash <= 0 && !fs.HasExited`

### User Experience
- **Accurate Achievements**: Win-based achievements now correctly validate actual game outcomes
- **Fair Progression**: Players can't earn win achievements by losing games
- **Clear Logic**: Win conditions are now explicit and match game outcomes

---

## Version 3.18.4 - Tech Enthusiast Fix & Points Display at Game End (2025-11-03)

### Bug Fixes

#### 🐛 Tech Enthusiast Achievement Logic
- **Issue**: Tech Enthusiast achievement was triggering incorrectly when non-tech sectors (like "Creative") were invested in
- **Fix**: Improved validation logic to ensure ALL sectors must be tech-related
- **Tech Sectors**: CloudTech, SaaS, DeepTech, FinTech, HealthTech, EdTech, LegalTech, Gaming, Security
- **Non-Tech Examples**: Creative, CleanTech, AgriTech, Social Media, Advertising, etc.
- **Result**: Achievement now correctly only triggers when exclusively investing in tech sectors

#### 🐛 Points Display at Game End
- **Issue**: Game end screen only showed total lifetime points, not available balance after new achievements
- **Fix**: Updated both VC and Founder mode end screens to show complete point breakdown
- **Display Now Shows**:
  - Available Balance: Spendable points (total - spent on upgrades)
  - Total Lifetime Points: All points ever earned
  - Spent Points: Amount spent on upgrades (if > 0)
- **Result**: Players can immediately see their new available balance after earning achievements

### Technical Changes

#### Modified Files
- `achievements/achievements.go`:
  - Improved `tech_enthusiast` achievement validation logic
  - Added explicit checks for investment count and sector validation
  - Ensures all sectors must be tech-related, not just some

- `main.go`:
  - Updated achievement display at VC mode game end
  - Added available balance calculation after new achievements
  - Shows complete point breakdown (available, total, spent)

- `founder_ui.go`:
  - Updated achievement display at Founder mode game end
  - Added available balance calculation after new achievements
  - Shows complete point breakdown (available, total, spent)

### User Experience
- **Accurate Achievements**: Tech Enthusiast now correctly validates investment portfolio
- **Clear Point Tracking**: Players see their new available balance immediately after earning achievements
- **Consistent Display**: Both VC and Founder modes show the same point breakdown format

---

## Version 3.18.3 - Achievement Fixes & Points Display (2025-11-03)

### Bug Fixes

#### 🐛 Achievement Logic Fixes
- **Issue**: "Cautious Investor" and "Risk Taker" achievements were triggering incorrectly
- **Fix**: Implemented proper risk score tracking and validation
- **Cautious Investor**: Now only triggers when ALL investments have risk < 0.3 (LOW risk)
- **Risk Taker**: Now only triggers when ALL investments have risk > 0.6 (High/VERY HIGH risk)
- **Result**: Achievements now correctly validate actual risk levels instead of always triggering

#### 🐛 Points Display Clarity
- **Issue**: Upgrade store showed total lifetime points instead of available balance
- **Fix**: Updated upgrade store to show both available balance and total lifetime points
- **Available Balance**: Shows spendable points (total - spent on upgrades)
- **Total Lifetime Points**: Shows all points ever earned (for career level calculation)
- **Spent Points**: Shows how much was spent on upgrades (if > 0)
- **Result**: Players can now clearly see how many points they have available vs total earned

### Technical Changes

#### Modified Files
- `achievements/achievements.go`:
  - Added `RiskScores []float64` field to `GameStats` struct
  - Updated `risk_taker` achievement check to validate all investments have risk > 0.6
  - Updated `cautious_investor` achievement check to validate all investments have risk < 0.3

- `main.go`:
  - Updated `displayUpgradeMenu()` to calculate and display available balance vs total lifetime points
  - Updated achievement stats collection to track risk scores for all investments
  - Updated `browseAllUpgrades()` to use available balance for affordability checks

### User Experience
- **Accurate Achievements**: Risk-based achievements now correctly reflect actual investment strategy
- **Clear Point Tracking**: Easy to see available balance vs total earned vs spent
- **Better Decisions**: Players can make informed upgrade purchase decisions

---

## Version 3.18.2 - Super Pro-Rata Fix & Build Fixes (2025-11-03)

### Bug Fixes

#### 🐛 Super Pro-Rata Upgrade Not Working
- **Issue**: Super Pro-Rata upgrade was active but UI still showed 20% max investment limit
- **Fix**: Updated startup display and investment prompt to check for Super Pro-Rata upgrade
- **Result**: Upgrade now correctly shows 50% max investment (instead of 20%) when active

#### 🐛 Build Compilation Errors
- **Issue**: Version 3.18.1 failed GitHub Actions builds due to missing method implementations
- **Fix**: Committed missing methods (`GetNextBoardVotePreview`, `GetSectorTrends`) and all related code
- **Result**: All builds now compile successfully across all platforms

### Technical Changes

#### Modified Files
- `main.go`:
  - Updated `displayStartup()` to check for Super Pro-Rata and show 50% limit when active
  - Updated `investmentPhase()` to validate against 50% limit when Super Pro-Rata is active
  - Fixed error messages to show correct percentage (50% vs 20%)

- `game/game.go`:
  - Verified `GetNextBoardVotePreview()` and `GetSectorTrends()` methods are properly committed

---

## Version 3.18.1 - Bug Fixes (2025-11-03)

### Bug Fixes

#### 🐛 Stats Display Separation
- **Issue**: Player stats were blended between VC and Founder modes
- **Fix**: Stats now display separately for each game mode
- **VC Mode Stats**: Shows stats for Easy/Medium/Hard/Expert difficulty games
- **Founder Mode Stats**: Shows stats for Founder difficulty games
- **Result**: Players can now see their performance in each game mode independently

#### 🐛 Upgrade Purchase Flow
- **Issue**: After purchasing an upgrade, system asked for player name again and showed stale point values
- **Fix**: Purchase flow now refreshes points and owned upgrades from database after each purchase
- **Result**: Points update immediately after purchase, allowing continuous purchases without re-entering name

### Technical Changes

#### Modified Files
- `database/database.go`:
  - Added `GetPlayerStatsByMode()` function to filter stats by game mode ("vc" or "founder")
  - VC mode filters by difficulty IN ('Easy', 'Medium', 'Hard', 'Expert')
  - Founder mode filters by difficulty = 'Founder'

- `main.go`:
  - Updated `displayPlayerStats()` to call `GetPlayerStatsByMode()` for both VC and Founder modes
  - Displays separate stat sections for each game mode
  - Modified `purchaseUpgrades()` to refresh points and owned upgrades from database after purchase
  - Purchase flow now loops back to purchase menu with updated values

### User Experience
- **Clearer Stats**: Players can now see how they perform in VC mode vs Founder mode separately
- **Smoother Purchases**: Upgrade purchases flow seamlessly without re-entering information
- **Real-time Updates**: Points update immediately after purchase

---

## Version 3.18.0 - Founder Mode Achievements & Upgrades (2025-11-03)

### Major Features Added

#### 🎯 Founder Mode Achievements & Points System
- **10 New Founder Achievements**: First Revenue, Profitable, $100K MRR Club, Unicorn MRR, Seed Raiser, Series A Graduate, IPO Exit, Acquired, 10K Customers, Bootstrapped
- **Point Accrual**: Founder mode now earns achievement points (10-100 points per achievement)
- **Career Progression**: Founder mode games contribute to career level progression
- **Score Saving**: Founder mode scores now saved to local database with ROI calculation
- **Profitability Tracking**: Automatically tracks when company reaches profitability for achievements

#### 🚀 Founder Mode Upgrades (8 New Upgrades!)
- **Fast Track** (200 pts): Start with 10% more product maturity
- **Sales Boost** (250 pts): +15% to initial MRR
- **Lower Burn** (300 pts): -10% monthly team costs permanently
- **Better Terms** (350 pts): Raise funding with 5% less equity given away
- **Quick Hire** (200 pts): First 3 hires cost 50% less
- **Market Insight** (250 pts): See competitor threat levels in competitor management
- **Churn Shield** (300 pts): Reduce churn by 10% permanently
- **Cloud Free First Year** (300 pts): No cloud compute costs for first 12 months

#### 💰 VC Mode Upgrades (2 New Upgrades!)
- **Seed Accelerator** (400 pts): First investment gets 25% equity bonus
- **Portfolio Insurance** (500 pts): Protect one investment from down round dilution per game

### Technical Changes

#### New Files
- None (upgrades added to existing `upgrades/upgrades.go`)

#### Database Changes
- Founder mode scores now saved to `game_scores` table with difficulty "Founder"
- Founder mode achievements tracked in `player_achievements` table

#### Modified Files
- `achievements/achievements.go`:
  - Extended `GameStats` struct with founder mode fields (FinalMRR, FinalValuation, Customers, FundingRoundsRaised, etc.)
  - Added 10 founder mode achievements
  - Updated achievement checking logic to support both VC and founder modes

- `founder/founder.go`:
  - Added `PlayerUpgrades []string` and `HiresCount int` fields to `FounderState`
  - Modified `NewFounderGame()` to accept `playerUpgrades` parameter
  - Added `MonthReachedProfitability` tracking
  - Implemented upgrade effects: Fast Track, Sales Boost, Lower Burn, Better Terms, Quick Hire, Churn Shield, Cloud Free First Year
  - Updated `CalculateTeamCost()` to apply Lower Burn upgrade
  - Updated `CalculateInfrastructureCosts()` to apply Cloud Free First Year upgrade
  - Updated `RaiseFundingWithTerms()` to apply Better Terms upgrade
  - Updated `HireEmployee()` to apply Quick Hire upgrade

- `founder_ui.go`:
  - Added `saveFounderScoreAndCheckAchievements()` function
  - Integrated achievement checking after founder mode games
  - Added upgrade display in welcome screen
  - Modified `playFounderMode()` to load and pass player upgrades

- `game/game.go`:
  - Added `InsuranceUsed` and `ProtectedCompany` fields to `GameState`
  - Implemented Seed Accelerator upgrade (25% equity bonus on first investment)
  - Implemented Portfolio Insurance upgrade (protect first investment from down rounds)
  - Updated `ProcessFundingRounds()` to check for Portfolio Insurance

- `upgrades/upgrades.go`:
  - Added new category: `CategoryFounderPerks`
  - Added 8 founder mode upgrades
  - Added 2 VC mode upgrades

### User Experience
- **Founder Mode Achievements**: Earn points and unlock achievements just like VC mode
- **Upgrade Visibility**: Active upgrades displayed when starting founder mode games
- **Portfolio Insurance Feedback**: Clear message when investment is protected from down round
- **Seed Accelerator Bonus**: First investment shows 25% equity bonus in investment summary

### Game Balance
- Founder mode upgrades priced competitively (200-350 pts)
- Seed Accelerator provides strong early game advantage (400 pts)
- Portfolio Insurance offers strategic protection (500 pts)
- All upgrades maintain game balance while providing meaningful progression

---

## Version 3.17.0 - Upgrade System & Meta-Progression (2025-11-03)

### Major Features Added

#### 🎁 Upgrade Store System
- **Meta-Progression**: Players can now spend achievement points to unlock permanent upgrades
- **Upgrade Menu**: New "5. Upgrades" option in main menu to browse, view, and purchase upgrades
- **Point-Based Economy**: Earn points from achievements (5-100 points per achievement), spend to unlock upgrades
- **Persistent Upgrades**: Purchased upgrades apply automatically to all future games

#### 💰 Financial Perks Upgrades
- **Fund Booster** (100 pts): +10% starting cash on all difficulties
- **Management Fee Reduction** (150 pts): Management fees reduced from 2% → 1.5%
- **Follow-On Reserve Boost** (200 pts): +$200k to follow-on reserve

#### 📈 Investment Terms Upgrades
- **Enhanced SAFE Discount** (150 pts): SAFE discount increases from 20% → 25%
- **Super Pro-Rata** (200 pts): Can invest up to 50% of round (vs 20% max)

#### 🎯 Information & Intel Upgrades
- **Early Access** (100 pts): See 2 extra startups before investment phase starts

#### ⚡ Game Mode Upgrades
- **Speed Mode** (200 pts): 30 turns instead of 60 (faster games)
- **Endurance Mode** (250 pts): 120 turns instead of 60 (longer games)

### Technical Changes

#### New Files
- `upgrades/upgrades.go`: Upgrade definitions, categories, and helper functions
- `UPGRADE_SYSTEM_PROPOSAL.md`: Comprehensive upgrade system design document

#### Database Changes
- Added `player_upgrades` table to track purchased upgrades
- Added `PurchaseUpgrade()`, `GetPlayerUpgrades()`, `HasUpgrade()` functions

#### Modified Files
- `game/game.go`:
  - Added `PlayerUpgrades []string` field to `GameState`
  - Modified `NewGame()` to accept and apply player upgrades
  - Updated `LoadStartups()` to show extra startups with Early Access upgrade
  - Updated `GenerateTermOptions()` to apply Enhanced SAFE Discount
  - Updated `MakeInvestmentWithTerms()` to apply Super Pro-Rata (50% max investment)
  - Applied upgrades for: fund booster, management fees, follow-on reserve, speed/endurance modes

- `main.go`:
  - Added upgrade menu to main menu (option 5)
  - Added `displayUpgradeMenu()` function with browse, view, and purchase options
  - Added `browseAllUpgrades()`, `viewPlayerUpgrades()`, `purchaseUpgrades()` functions
  - Modified `playNewGame()` to load and pass player upgrades to game initialization

- `database/database.go`:
  - Added `player_upgrades` table creation
  - Added upgrade-related database functions

### User Experience
- **Upgrade Store UI**: Clean, categorized display of all upgrades with ownership status
- **Point Tracking**: Shows current points, career level, and points needed for next level
- **Smart Filtering**: Only shows available upgrades you can afford
- **Visual Feedback**: Clear indicators for owned vs. available vs. locked upgrades

### Game Balance
- Upgrades are priced to provide meaningful progression without breaking game balance
- Early upgrades (100-150 pts) provide immediate value
- Advanced upgrades (200+ pts) unlock new strategies and game modes
- Meta-progression encourages replayability and long-term engagement

---

## Version 3.16.1 - Follow-On Investment Validation Fix (2025-11-03)

### Bug Fixes

#### 🐛 Fixed Follow-On Investment Validation Logic
- **Issue**: Follow-on investment validation incorrectly checked cumulative total investment instead of per-round limit
- **Problem**: When investing $100k in a follow-on round with $118k max, it rejected because it checked `($80k existing + $100k new) > $118k`
- **Fix**: 20% investment limit now applies per round separately, not cumulatively
- **Result**: Users can invest up to 20% of current pre-money valuation in each round, regardless of previous investments
- **Example**: Company valued at $592k pre-money → can invest up to $118k in THIS round, even if you invested $80k previously

### Technical Changes

#### Modified Files
- `game/game.go`:
  - Updated `MakeFollowOnInvestment()` validation to check only the new investment amount against 20% limit
  - Changed error message to clarify it's checking the follow-on amount, not cumulative total
  - The 20% limit now correctly applies to each funding round independently

---

## Version 3.16.0 - Investment Limit Enforcement & UI Improvements (2025-11-03)

### Major Features Added

#### 💰 20% Investment Limit Per Company
- **Standard VC Practice**: Maximum investment per company is now 20% of valuation (not 100%)
- **Realistic Investment Model**: Only 20% of company equity is available for investment in each round
- **Applies to All Investments**: Both initial and follow-on investments respect the 20% limit
- **Total Investment Cap**: Follow-on investments are capped so total investment doesn't exceed 20% of valuation
- **Example**: Company valued at $500k → Max investment is $100k (20% of $500k)

#### 🎯 Investment Validation
- **Minimum Investment**: $10,000 minimum for all investments (initial and follow-on)
- **Maximum Investment**: 20% of company valuation (enforced backend and UI)
- **Clear Error Messages**: Shows exact limits when validation fails
- **AI Players**: AI investments also respect the 20% limit

#### 📊 Enhanced UI Display
- **Company List Shows Max Investment**: Each startup displays "Max Investment: $X" in the list
- **Cash-Aware Display**: Shows if max is limited by available cash vs. valuation limit
- **Investment Selection Screen**: Shows company name, valuation, and max investment available when selecting
- **Clear Investment Range**: Prompt shows "$10,000 - $X" with clear max limit
- **No More Guessing**: Users can see investment limits before selecting companies

### Technical Changes

#### Modified Files
- `game/game.go`:
  - Added $10,000 minimum investment validation in `MakeInvestmentWithTerms()`
  - Added 20% maximum investment validation (20% of valuation)
  - Updated `MakeFollowOnInvestment()` to validate total investment doesn't exceed 20%
  - Updated `GetFollowOnOpportunities()` to calculate max as 20% of pre-money valuation
  - Updated AI investment logic to respect 20% limit
  - Equity calculation now assumes only 20% of company is available for investment

- `main.go`:
  - Updated `displayStartup()` to show max investment amount (considering cash limits)
  - Enhanced investment selection UI with company details and max investment display
  - Added pre-validation in UI before calling backend (better UX)
  - Shows cash-limited vs. valuation-limited max investment clearly

### Bug Fixes
- **Fixed Equity Calculation**: Equity now correctly calculated based on 20% available pool
- **SAFE Discount Limits**: SAFE discounts capped at 24% (20% × 1.2) maximum equity
- **Investment Validation**: Both UI and backend validate investment amounts

---

## Version 3.15.0 - Difficulty Balance & Risk Adjustments (2025-11-02)

### Major Features Added

#### 💰 Difficulty Level Balance Updates
- **More Capital for Harder Difficulties**: Harder difficulties now provide MORE capital (not less):
  - **Easy**: $1M fund (unchanged)
  - **Medium**: $750k → **$1.5M** (+$750k)
  - **Hard**: $500k → **$2M** (+$1.5M)
  - **Expert**: $500k → **$2.5M** (+$2M)
- **Rationale**: Higher difficulty = more capital to deploy, but increased risk events and volatility
- **Risk vs. Reward**: More money available means more opportunities, but adverse events increase

#### 🎯 Risk Score Minimum Adjustment
- **Minimum Risk is Medium**: All startups now have minimum risk score of 0.5 (medium)
- **Risk Range**: Changed from 0.0-1.0 to **0.5-1.0**
- **No Low-Risk Companies**: Even "safer" companies start at medium risk level
- **Difficulty Impact**: Higher difficulty = more adverse events occur, making risk management critical

#### 💵 Uninvested Cash Available for Follow-On Investments
- **Cash Available for Later Rounds**: Uninvested cash from initial fund is now available for follow-on investments
- **Combined Funds**: Follow-on investments can use `Cash + FollowOnReserve` (not just reserve)
- **Example**: Invest $600k initially, keep $400k cash → Available for follow-on = $400k + $1M reserve
- **Smart Deduction**: Uses cash first, then follow-on reserve
- **UI Updates**: Shows "Available Funds: $X (Cash: $Y + Reserve: $Z)" for clarity

### Technical Changes

#### Modified Files
- `game/game.go`:
  - Updated difficulty levels: Medium=$1.5M, Hard=$2M, Expert=$2.5M
  - Modified `calculateRiskScore()` to enforce minimum 0.5 (medium) risk
  - Updated `GetFollowOnOpportunities()` to use `Cash + FollowOnReserve` for max investment
  - Updated `MakeFollowOnInvestment()` to deduct from cash first, then reserve
  - Changed risk calculation range from 0.0-1.0 to 0.5-1.0

- `main.go`:
  - Updated follow-on investment UI to show combined available funds
  - Updated difficulty display messages to clarify cash availability
  - Updated FAQ to mention cash availability for follow-on investments

---

## Version 3.14.0 - Board Voting Mechanism & Equity Fixes (2025-11-02)

### Major Features Added

#### 🏛️ Board Voting System
- **Board Seat Participation**: Players with board seats ($100k+ Preferred Stock investments) now vote on critical company decisions
- **Voting Scenarios**:
  - **Acquisition Offers**: Vote to accept or reject acquisition offers
  - **Down Rounds**: Vote to approve or reject down rounds (valuations below current)
- **Voting Mechanics**:
  - Player vote counts as 1 vote
  - AI board members (2-3 simulated) vote based on strategy
  - Majority wins; outcomes are executed based on vote results
- **Interactive UI**: Clear voting interface with options, consequences, and results displayed

#### 🤖 AI Player Enhancements
- **Board Seats for AI**: AI players now get board seats when investing $100k+ (same as players)
- **Preferred Stock Terms**: AI investments now use Preferred Stock terms with all protections
- **Equal Starting Capital**: AI players start with same capital as player ($1M on Easy difficulty)

### Bug Fixes

#### 🐛 Equity Calculation Fixes
- **Fixed Equity Over 100% Bug**: Follow-on investments now correctly recalculate total equity instead of adding
- **SAFE Conversion Discount**: SAFE conversion discounts now properly apply during follow-on investments
- **Double Dilution Prevention**: Added flag to prevent double dilution when making follow-on investments
- **Equity Safety Cap**: Added 100% cap as safety check (should never be needed)

#### 📊 Investment Terms Improvements
- **SAFE Terms Execution**: Confirmed SAFE conversion discounts are correctly applied in math
- **Equity Recalculation**: Total equity now recalculated as `(totalAmountInvested / postMoneyValuation) * 100`
- **Follow-on Investment Accuracy**: Follow-on investments now use correct post-money valuation for equity calculation

### Technical Changes

#### Modified Files
- `game/game.go`:
  - Added `BoardVote` struct and `PendingBoardVotes` queue
  - Added `HasBoardSeat()`, `HasAnyBoardSeat()`, `GetPendingBoardVotes()` functions
  - Added `ProcessBoardVote()` and `ExecuteBoardVoteOutcome()` functions
  - Updated `MakeFollowOnInvestment()` to recalculate equity correctly and apply SAFE discounts
  - Updated `ProcessFundingRounds()` to create board votes for down rounds
  - Updated `ProcessAcquisitions()` to create board votes for acquisitions
  - Updated `AIPlayerMakeInvestments()` to include Preferred Stock terms
  - Added `FollowOnThisTurn` flag to prevent double dilution
  - Fixed equity calculation to use total invested amount

- `main.go`:
  - Added `handleBoardVotes()` function for interactive voting UI
  - Integrated board vote handling into turn processing flow

---

## Version 3.13.0 - 16 New Startup Options (2025-11-02)

### Major Features Added

#### 🚀 16 New Startup Options Across All Categories
- **Total Startups Increased**: 21 → 36 startups (+71% more choices!)
- **SaaS**: 9 → 14 startups (+5 new)
  - **GrowthEngine**: Marketing automation and attribution
  - **LegalStack**: Contract lifecycle management
  - **DealFlow Pro**: Sales engagement and pipeline intelligence
  - **TeamSync**: All-in-one workspace (Slack + Zoom + Asana replacement)
  - **DataUnify**: Customer data platform for personalization
  
- **DeepTech**: 6 → 9 startups (+3 new)
  - **OrbitLaunch**: Small satellite launch services and orbital logistics
  - **NanoMat Industries**: Advanced nanomaterials (graphene-based)
  - **FarmTech AI**: Precision agriculture with autonomous robots
  
- **FinTech**: 2 → 5 startups (+3 new)
  - **CreditFlow**: AI-powered lending for underbanked businesses
  - **ChainPay**: Cryptocurrency payment infrastructure
  - **InsureTech Pro**: Embedded insurance platform
  
- **HealthTech**: 2 → 4 startups (+2 new)
  - **MindWell AI**: Mental health platform for employers
  - **PharmaDirect**: Digital pharmacy and prescription management
  
- **GovTech**: 2 → 4 startups (+2 new)
  - **VoteSecure**: Blockchain-based voting and election management
  - **FirstRespond**: Emergency response coordination and dispatch

#### 💡 Diverse New Industries Represented
- **Space Tech**: Satellite launch services (OrbitLaunch)
- **Agriculture**: Autonomous farming robots (FarmTech AI)
- **Advanced Materials**: Nanomaterials for aerospace (NanoMat)
- **Mental Health**: Workplace wellness platforms (MindWell AI)
- **Pharmacy Tech**: Digital prescription management (PharmaDirect)
- **Legal Tech**: Contract management with AI (LegalStack)
- **Crypto Infrastructure**: Business payment gateways (ChainPay)
- **InsureTech**: Embedded insurance APIs (InsureTech Pro)
- **Election Tech**: Secure digital voting (VoteSecure)
- **Emergency Services**: First responder coordination (FirstRespond)

#### 📊 Realistic Parameters & Variety
- Deal sizes range from $300/mo (consumer fintech) to $500K (space launches)
- Initial cash ranges from $650K (education SaaS) to $5M (space tech)
- Competition levels vary from "low" (govtech) to "very high" (collaboration, marketing)
- Market sizes from 300 (space) to 200,000 (consumer fintech)
- Each startup has unique team compositions and growth dynamics

### Technical Changes

#### Modified Files
- `founder/startups.json`:
  - Added 16 new startup templates across all categories
  - Balanced parameters for realistic gameplay
  - Diverse initial team compositions
  - Industry-appropriate CAC and churn rates

---

## Version 3.12.0 - International Market Management & Regional Competitors (2025-11-02)

### Major Features Added

#### 📍 Market Assignment for Employees
- **Assign Employees to Specific Markets**: When hiring Sales, Marketing, or CS roles with active international markets:
  - Choose which market they focus on (USA, Europe, Asia, etc.)
  - Option to assign to "All Markets" (works globally)
  - Growth and churn mitigation now market-specific
- **Market-Based Impact**:
  - Sales/Marketing only boost growth in their assigned market
  - CS team only reduces churn in their assigned market
  - Executives (CGO, COO) work across all markets
- **Team Roster Shows Assignments**: View which employees are assigned to which markets
- **Example Impact**:
  - Before: All sales reps boosted all markets equally (unrealistic)
  - After: Assign 2 sales reps to Asia, 1 to Europe for targeted growth

#### 🌍 Regional Competitors Appear in New Markets
- **Competitors Generated on Market Entry**:
  - **Very High Competition**: 2-3 competitors appear (10-30% market share each)
  - **High Competition**: 1-2 competitors (8-23% market share)
  - **Medium Competition**: 0-1 competitors (5-15% market share)
  - **Low Competition**: No competitors
- **Real Regional Companies**:
  - **Europe**: Zalando, Klarna, N26, TransferWise, BlaBlaCar, Deliveroo EU
  - **Asia**: Grab, GoJek, Alibaba Local, Meituan, Tokopedia, Paytm
  - **LATAM**: Nubank, Mercado Libre, Rappi, Kavak, Creditas, QuintoAndar
  - **Middle East**: Careem, Souq, Fetchr, Talabat, Noon, Swvl
  - **Africa**: Jumia, Flutterwave, Andela, Paystack, Konga, M-Pesa
  - **Australia**: Afterpay, Canva AU, Atlassian, WiseTech, Xero, SEEK
- **Manage Like Other Competitors**: View in competitor menu, choose strategy (ignore/compete/partner)
- **Impact**: Competitive markets now have actual named competitors to deal with

#### 💬 Clearer Churn Mitigation in International Markets
- **Improved Messaging**: When expanding to new markets, clear guidance on:
  - "Assign CS team to this market to reduce churn!"
  - "Assign Sales/Marketing to grow this market faster!"
- **Churn Calculation Now Market-Specific**:
  - CS reps assigned to market reduce churn by 2% each
  - COO reduces churn by 6% across all markets
  - No CS in market = +30% churn penalty
  - Shows exactly which team members help which markets
- **Visual Feedback**: Competitor alerts show threat level and market share immediately

### Technical Changes

#### Modified Files
- `founder/founder.go`:
  - Added `AssignedMarket` field to `Employee` struct
  - Created `HireEmployeeWithMarket()` function for market-specific hiring
  - Updated `UpdateGlobalMarkets()` to only count employees assigned to each market
  - Rewrote churn calculation to check CS team assigned to specific markets
  - Added regional competitor generation in `ExpandToMarket()`
  - 40+ regional competitor names across 6 markets
  - Threat levels and market share based on competition level

- `founder_ui.go`:
  - Updated `handleHiring()` to ask for market assignment for Sales/Marketing/CS
  - Shows available markets with current performance metrics
  - Modified `handleViewTeamRoster()` to display market assignments
  - Enhanced `handleGlobalExpansion()` messaging with churn mitigation tips
  - Added competitor detection alerts when entering new markets
  - Shows new competitor names, threat levels, and market share

### Examples

**Before**: "Hired 3 sales reps" → all helped all markets equally

**After**: "Hired sales rep for Asia" → only boosts Asia growth; "Hired sales rep for All Markets" → helps globally

**Before**: Expand to Asia → no competitors mentioned

**After**: Expand to Asia → "2 new competitors detected! Grab (Asia) - High threat, 15% market share, Meituan (Asia) - High threat, 22% market share"

---

## Version 3.11.1 - Bug Fix: Acquisition Payout Breakdown on Final Screen (2025-11-02)

### Bug Fixes

#### 🐛 Fixed Missing Acquisition Payout Breakdown
- **Issue**: Acquisition payout breakdown was only shown when accepting the offer, not on the final results screen
- **Impact**: Players couldn't see how the acquisition split among investors, executives, and employees at game end
- **Fix**: Added complete cap table breakdown to final score screen for acquisition exits
- **Now Shows**:
  - Founder payout with equity %
  - Each investor's payout (by name)
  - Executive team payouts (with names like Gilfoyle, Jared Dunn, etc.)
  - Employee equity payouts
  - Unallocated employee pool
  - Total validation showing 100% of acquisition price

### Technical Changes

#### Modified Files
- `founder_ui.go`:
  - Updated `displayFounderFinalScore()` to show cap table breakdown for acquisition exits
  - Added same payout display logic as `displayAcquisitionOffer()`
  - Maintains consistent formatting between offer screen and final results

---

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

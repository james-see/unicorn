# Realistic VC Fund Simulation - New Features

This document describes all the realistic VC mechanics that have been integrated into the Unicorn game.

## ?? Fund Structure

### Larger, More Realistic Fund Sizes
- **Easy Mode**: $1,000,000 fund
- **Medium Mode**: $750,000 fund  
- **Hard Mode**: $500,000 fund
- **Expert Mode**: $500,000 fund (with extreme volatility)

### Management Fees (Carried Interest)
- **2% Annual Management Fee** charged regardless of performance
- Fees are charged monthly (1/12 of annual rate each month)
- Displayed in portfolio view and final score
- Applies to both player and AI competitors

Example:
- $1M fund = $20,000/year in management fees ($1,667/month)
- $500K fund = $10,000/year in management fees ($833/month)

## ?? Multiple Funding Rounds

### Scheduled Rounds for Each Company
Companies automatically raise multiple rounds throughout the game:

1. **Seed Round** (6-12 months)
   - Raises ~10% of current valuation
   
2. **Series A** (18-36 months)
   - Raises ~20% of current valuation
   
3. **Series B** (36-60 months)
   - Raises ~33% of current valuation
   
4. **Series C** (60-90 months, 50% of companies)
   - Raises ~50% of current valuation

### Post-Money Valuation
- **Pre-Money Valuation**: Company value before new investment
- **Post-Money Valuation**: Pre-money + new investment amount
- Formula: `Post-Money = Pre-Money + Investment`

## ?? Dilution Mechanics

### Equity Dilution
When a company raises a new round, existing shareholders get diluted:

**Dilution Formula**: 
```
New Equity % = Old Equity % ? (Pre-Money / Post-Money)
```

**Example**:
- You own 5% of a company worth $10M
- Company raises $5M Series A
- Pre-money: $10M, Post-money: $15M
- Your new equity: 5% ? ($10M / $15M) = 3.33%
- Your value still increases: 3.33% of $15M = $500K vs 5% of $10M = $500K

### Tracking
- Initial equity percentage is preserved
- Current equity shown with dilution info
- All funding rounds are tracked per investment
- Dilution notifications appear as news events

## ?? AI Competitor System

### Three AI VC Players

1. **CARL** - Sterling & Cooper
   - **Strategy**: Conservative
   - **Risk Tolerance**: 30% (low)
   - **Behavior**: Invests in low-risk, proven companies with good margins
   
2. **Sarah Chen** - Accel Partners
   - **Strategy**: Aggressive
   - **Risk Tolerance**: 80% (high)
   - **Behavior**: Chases high-risk, high-growth moonshots
   
3. **Marcus Williams** - Sequoia Capital
   - **Strategy**: Balanced
   - **Risk Tolerance**: 50% (medium)
   - **Behavior**: Mixes risk across portfolio

### AI Decision Making
- AI players make their own investment choices at game start
- Invest in 3-6 companies based on strategy
- Subject to same events, dilution, and fees as player
- Portfolios update every turn

## ?? Competitive Leaderboard

### During Game
- **Mini Leaderboard** shown every quarter (every 3 months)
- Shows current standings of all players
- Displays net worth and ROI for each competitor

### Final Results
- Complete leaderboard at game end
- Rankings by net worth
- Shows who won the competition
- Highlights player's position

**Leaderboard Format**:
```
RANK  INVESTOR              FIRM                  NET WORTH    ROI
1     Marcus Williams       Sequoia Capital       $2,145,000   114.5%
2     Your Name ? YOU       Your Fund             $1,987,500   98.8%
3     Sarah Chen            Accel Partners        $1,654,000   65.4%
4     CARL                  Sterling & Cooper     $1,321,000   32.1%
```

## ?? Enhanced Display Features

### Investment Phase
- Shows fund size prominently
- Displays management fees to be charged
- Lists all AI competitors and their strategies

### Turn Display
- Management fees shown (annually)
- Dilution events with before/after equity percentages
- Company funding round announcements
- Quarterly competitive leaderboard

### Portfolio View
Enhanced investment display showing:
```
CompanyName: $100,000 invested, 3.45% equity (was 5.00%, 2 rounds)
  Current Value: $500,000 (+$400,000)
```

### Final Score
- Total management fees paid
- Full competitive leaderboard
- Victory/defeat message based on ranking
- Detailed portfolio breakdown

## ?? Strategic Implications

### Early Investment Advantages
- Invest before dilution rounds
- Higher equity percentages initially
- More upside if company succeeds

### Management Fee Impact
- 2% annual fee reduces available capital
- Must generate returns above fees to profit
- Affects all players equally

### Dilution Strategy
- Early investments get diluted more
- But also benefit from lower valuations
- Balance timing vs. dilution risk

### AI Competition
- Different strategies create varied outcomes
- Conservative vs aggressive approaches
- Learn from AI behavior patterns

## ?? Technical Implementation

### New Data Structures
- `FundingRound`: Tracks each round's details
- `Investment`: Enhanced with dilution tracking
- `AIPlayer`: Computer opponent data
- `PlayerScore`: Leaderboard entries

### New Functions
- `ProcessFundingRounds()`: Handles dilution
- `ProcessManagementFees()`: Monthly fee deduction
- `AIPlayerMakeInvestments()`: AI decision logic
- `GetLeaderboard()`: Competitive rankings
- `ScheduleFundingRounds()`: Plan future rounds

### Game Flow
1. Player and AI make initial investments
2. Each turn:
   - Charge management fees (monthly)
   - Process scheduled funding rounds
   - Apply random events
   - Update valuations
   - Show quarterly leaderboards
3. Final leaderboard and winner announcement

## ?? Example Game Flow

**Month 1**: 
- Start with $1M fund
- Player invests $300K across 3 companies
- AI players make their investments

**Month 8**:
- Company A raises Seed round ($2M)
- Your 5% equity dilutes to 4.2%
- Management fee: $1,667 charged

**Month 20**:
- Company A raises Series A ($8M)
- Your 4.2% dilutes to 3.1%
- Quarterly leaderboard: You're #2

**Month 120** (End):
- Total management fees: $20,000
- Final net worth: $2.1M
- ROI: 110% (after fees)
- Final rank: #1 - Victory!

## ?? Future Enhancements

Possible additions:
- Pro-rata rights to participate in follow-on rounds
- Carry interest (20% of profits above hurdle)
- LP commitments and capital calls
- Term sheets and valuation caps
- Employee option pool dilution
- Liquidation preferences
- More AI personalities and strategies
- Historical performance tracking across games

---

**Note**: All these features are now live in the game. Start a new game to experience realistic VC fund simulation with competition, dilution, multiple rounds, and management fees!

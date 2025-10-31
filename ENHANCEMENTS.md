# Unicorn Game Enhancement Suggestions

## Core Gameplay Enhancements

### 1. **Investment System**
- Allow players to invest specific amounts in startups (currently just displays)
- Track portfolio holdings with ownership percentages
- Implement exit events (acquisitions, IPOs, failures)
- Add partial investment options (not all-or-nothing)

### 2. **Turn-Based Game Loop**
- Implement the 120-month (10-year) gameplay loop
- Monthly events affecting portfolio companies
- Random market events (recessions, tech booms, etc.)
- Company performance updates (growth, churn, milestones)

### 3. **Strategic Decision Making**
- **Due Diligence System**: Spend money to research companies before investing
- **Diversification Requirements**: Penalties/bonuses for portfolio concentration
- **Follow-on Rounds**: Option to invest more in portfolio companies
- **Exit Timing**: Choose when to sell holdings (immediate vs. wait for better valuation)

### 4. **Company Dynamics**
- **Growth Trajectories**: Companies grow/decline based on metrics
- **Competition**: Rival companies can impact valuations
- **Team Quality**: Hidden stat affecting company success
- **Market Fit**: How well company fits current market conditions

### 5. **Events & Randomness**
- **Market Events**: Economic downturns, tech bubbles, regulatory changes
- **Company-Specific Events**: Use the round-options.json for monthly events
- **Portfolio Synergy**: Companies can help/hurt each other
- **Liquidity Events**: Acquisitions, IPOs, shutdowns happen randomly

## Strategy & Depth

### 6. **Risk Management**
- **Portfolio Risk Score**: Track diversification and risk exposure
- **Emergency Fund**: Maintain cash reserves for opportunities/crises
- **Risk/Reward Profiles**: Different investment strategies (aggressive, conservative, balanced)

### 7. **Investment Types**
- **Seed/Series A/B/C**: Different stages with different risk/reward
- **Angel Investments**: Smaller checks, higher risk
- **Strategic Investments**: Corporate partnerships
- **Secondary Market**: Buy/sell existing shares

### 8. **Metrics & Analytics**
- **Portfolio Dashboard**: Real-time portfolio value, ROI, performance metrics
- **Company Health Scores**: Track burn rate vs runway, growth metrics
- **Market Trends**: Show sector performance, hot categories
- **Prediction Models**: Historical data to inform decisions

### 9. **Player Actions**
- **Research Actions**: Pay for market research, competitor analysis
- **Board Meetings**: Influence company decisions (costs money/time)
- **Networking**: Unlock better deals through connections
- **Fundraising**: Raise additional capital (with trade-offs)

## Scoring & Progression

### 10. **Local High Scores**
- JSON file storing top scores with player name, final portfolio value, date
- Display top 10 leaderboard
- Score based on final portfolio value + multipliers for:
  - Time to first exit
  - ROI percentage
  - Number of successful exits
  - Portfolio diversity score

### 11. **Achievements/Badges**
- "First Unicorn" - Invest in a company that reaches $1B valuation
- "Diversifier" - Invest in 10+ different companies
- "Patient Investor" - Hold an investment for 5+ years
- "Exit Master" - Successfully exit 5+ investments
- "Survivor" - Survive a market crash
- "Early Bird" - Invest in seed rounds only

### 12. **Difficulty Levels**
- **Easy**: More cash, fewer bad events, slower burn rates
- **Normal**: Current settings
- **Hard**: Less cash, more events, faster burn rates, stricter rules
- **Expert**: Hidden company stats, limited information, aggressive markets

## Fun & Engagement

### 13. **Narrative Elements**
- **Company Stories**: Rich descriptions of what each startup does
- **News Headlines**: Monthly news affecting portfolio
- **Founder Personalities**: Different founder types affect company outcomes
- **Industry Trends**: Follow trends (AI, crypto, biotech waves)

### 14. **Visual Enhancements**
- **ASCII Charts**: Show portfolio value over time
- **Company Status Indicators**: Visual health meters
- **Progress Bars**: Show time remaining, milestones
- **Color Coding**: Green/red for gains/losses

### 15. **Replayability**
- **Random Seed**: Different startup sets each game
- **Dynamic Events**: Events scale with portfolio size
- **Unlock System**: Unlock new startup types/events with achievements
- **Save/Load**: Save game state to resume later

### 16. **Multi-Company Management**
- Show all portfolio companies at once
- Compare performance side-by-side
- Make decisions for multiple companies per turn
- Portfolio rebalancing options

## Technical Enhancements

### 17. **Data Persistence**
- Save game state (JSON)
- Load saved games
- High scores persistence
- Configurable game settings

### 18. **Improved UX**
- Clear menu system with navigation
- Input validation and error handling
- Help/instructions accessible in-game
- Keyboard shortcuts for common actions

### 19. **Testing & Balance**
- Game balance testing (ensure winnable but challenging)
- Event probability tuning
- Startup stat ranges validation
- Performance optimization

## Implementation Priority

### Phase 1 (Core Gameplay)
- Investment system with portfolio tracking
- Turn-based loop (120 months)
- Basic events system
- Local high scores

### Phase 2 (Strategy)
- Due diligence system
- Company growth/decline mechanics
- Exit events
- Portfolio dashboard

### Phase 3 (Polish)
- Achievements
- Save/load
- Visual improvements
- Narrative elements

### Phase 4 (Advanced)
- Multiple difficulty levels
- Advanced investment types
- Complex events and synergies
- Analytics and prediction systems

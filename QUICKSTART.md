# Quick Start Guide

## Building and Running

```bash
# Build the game
go build -o unicorn

# Run it
./unicorn
```

## Main Menu

When you start the game, you'll see:
```
?? UNICORN - MAIN MENU ??

1. New Game
2. Leaderboards
3. Player Statistics
4. Quit
```

## How to Play a New Game

1. **Select "1. New Game"** from main menu

2. **Enter your name** when prompted

3. **Choose difficulty** (1-4):
   - Easy: $500k starting cash, gentler volatility
   - Medium: $250k starting cash, balanced
   - Hard: $150k starting cash, more volatile
   - Expert: $100k starting cash, extreme volatility, only 7.5 years!

4. **Review the 10 available startups** - each shows:
   - Company name and description
   - Valuation (in millions)
   - Risk level (Low/Medium/High)
   - Growth potential (Low/Medium/High)
   - Key metrics (sales, margins, visitors)

3. **Make your investments:**
   - Type the company number (1-10)
   - Enter investment amount (e.g., 50000)
   - Repeat for multiple companies
   - Type 'done' when finished

4. **Watch the game unfold:**
   - Each turn = 1 month
   - Random events affect your companies
   - Portfolio value updates automatically
   - Green = profit, Red = loss

5. **Final score after 90-120 turns (depends on difficulty):**
   - Final net worth
   - ROI percentage
   - Performance rating
   - Detailed portfolio breakdown
   - **Score automatically saved to leaderboard!**

6. **Return to main menu** to:
   - Play again with different difficulty
   - View leaderboards
   - Check your statistics

## Leaderboards (Option 2)

View top scores in multiple ways:
- **By Net Worth (All)** - Highest total wealth
- **By ROI (All)** - Best return percentage
- **Easy/Medium/Hard/Expert** - Difficulty-specific boards
- **Recent Games** - Last 10 games played

Features:
- ?????? Color-coded rankings
- Shows: Player, Net Worth, ROI, Exits, Difficulty
- Top 10 for each category

## Player Statistics (Option 3)

Check career performance:
- ?? Total games played
- ?? Best net worth ever
- ?? Best ROI percentage
- ?? Total successful exits
- ?? Average net worth
- ?? Win rate (% positive ROI)

## Difficulty Comparison

| Level | Cash | Event % | Volatility | Turns | Best For |
|-------|------|---------|-----------|-------|----------|
| Easy | $500k | 20% | 3% | 120 | Learning |
| Medium | $250k | 30% | 5% | 120 | Standard |
| Hard | $150k | 40% | 7% | 120 | Challenge |
| Expert | $100k | 50% | 10% | 90 | Masters |

## Pro Tips

- **Diversify:** Don't put all your money in one company
- **Risk/Reward:** High risk companies can have huge payoffs
- **Read carefully:** Company metrics give hints about potential
- **Budget wisely:** You start with $250,000 - make it count!
- **Long game:** Companies take time to grow (or fail)

## Example Strategy

**Conservative:**
- Invest in 4-5 low-risk companies
- Keep $50k as reserve
- Aim for steady 100-200% ROI

**Aggressive:**
- Go all-in on 1-2 high-growth companies
- Higher risk, potential 500%+ ROI
- Could lose everything!

**Balanced:**
- 3-4 medium-risk companies
- Spread across different sectors
- Target 200-300% ROI

## Performance Ratings

?? **1000%+ ROI** = UNICORN HUNTER - Legendary!
?? **500%+ ROI** = Elite VC - Outstanding!
? **200%+ ROI** = Great Investor - Excellent!
?? **50%+ ROI** = Solid Performance - Good!
?? **0%+ ROI** = Break Even - Not Bad
?? **Negative ROI** = Lost Money - Try Again

## Competitive Play

**Beat Your Own Records:**
- Try to top your personal best
- Improve your win rate
- Master all difficulty levels

**Compete on Leaderboards:**
- Get #1 net worth in your difficulty
- Highest ROI across all players
- Most successful exits

**Challenge Yourself:**
- Start on Easy, work up to Expert
- Try to win with minimal investments
- Diversify vs. focused strategy

## Data Persistence

- All scores saved to `unicorn_scores.db`
- Database created automatically on first run
- Portable - can backup or share
- Never expires - build your legacy!

Good luck! ??

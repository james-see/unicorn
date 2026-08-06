# Quick Start Guide

## Install

### Easiest: Homebrew

```bash
brew install james-see/tap/unicorn
```

> Don't have Homebrew? Install it from [brew.sh](https://brew.sh).

### Download Binary

Grab the latest binary from the [Releases page](https://github.com/james-see/unicorn/releases/latest). Pick the file matching your platform, make it executable, and run:

```bash
chmod +x unicorn
./unicorn
```

### Build from Source

Requires [Go 1.24+](https://go.dev/dl/):

```bash
go build -o unicorn
./unicorn
```

## Main Menu

When you start the game, you'll see:

```
🦄 UNICORN — MAIN MENU

1. New Game
2. Leaderboards
3. Player Statistics
4. Achievements
5. Progression & Levels
6. Analytics Dashboard
7. Upgrades
8. Help & Info
9. Quit
```

## VC Investor Mode

1. Select **New Game** → choose **VC Investor Mode**
2. Enter your name and firm name
3. Choose difficulty:

| Level | Fund Size | Reserve | Event Chance | Volatility | Turns |
|---|---|---|---|---|---|
| Easy | $1M | $2.5M | 20% | 3% | 60 |
| Medium | $1.5M | $3M | 30% | 5% | 60 |
| Hard 🔒 | $2M | $3.5M | 40% | 7% | 60 |
| Expert 🔒 | $2.5M | $4M | 50% | 10% | 60 |

> 🔒 Hard unlocks at player level 5, Expert at level 10.

4. Browse 45 available startups — each shows valuation, risk level, and growth potential
5. Select a startup and enter your investment amount
6. Press `d` when done investing — the game begins!
7. Each turn = 1 month. Watch for:
   - **Follow-on opportunities** — invest more when portfolio companies raise new rounds
   - **Board votes** — vote on acquisitions, down rounds, strategic pivots
   - **Random events** — market shifts, scandals, breakthroughs
   - **Opportunity Fund** — when a company hits 3x growth, extra capital unlocks for follow-ons
8. After 60 turns, see your final score: net worth, ROI, and performance rating

### Performance Ratings

| ROI | Rating |
|---|---|
| 1000%+ | 👑 UNICORN HUNTER — Legendary! |
| 500%+ | 🏆 Elite VC — Outstanding! |
| 200%+ | ⭐ Great Investor — Excellent! |
| 50%+ | ✓ Solid Performance — Good! |
| 0%+ | = Break Even — Not Bad |
| < 0% | ⚠ Lost Money — Try Again |

## Startup Founder Mode

1. Select **New Game** → choose **Startup Founder Mode**
2. Pick a startup template (SaaS, DeepTech, GovTech, Hardware)
3. Each month, make strategic decisions:
   - Hire team (engineers, sales, CS, marketing, executives)
   - Spend on marketing to acquire customers
   - Launch partnerships and affiliate programs
   - Raise funding rounds (Seed → Series A → Series B)
   - Manage board, equity pool, and advisors
   - Expand to new markets
4. Watch your MRR grow, manage churn, track SaaS metrics
5. Exit when ready:
   - 🏛️ IPO — $20M ARR, 40%+ growth
   - 🤝 Acquisition — $5M ARR, 50+ customers
   - 💼 Secondary Sale — $10M ARR, profitable

## Difficulty Comparison

| Level | Fund | Reserve | Opp Fund | Max New Bets | Events | Volatility | Best For |
|---|---|---|---|---|---|---|---|
| Easy | $1M | $2.5M | $1.5M | 12 | 20% | 3% | Learning |
| Medium | $1.5M | $3M | $2.25M | 10 | 30% | 5% | Standard |
| Hard | $2M | $3.5M | $2.5M | 8 | 40% | 7% | Challenge |
| Expert | $2.5M | $4M | $2.5M | 8 | 50% | 10% | Masters |

## Pro Tips

- **Diversify:** Spread investments across sectors to reduce risk
- **Reserve discipline:** Don't deploy all capital early — later rounds are expensive
- **Power law:** Concentrate follow-ons on your top 2-3 performers, not every company
- **Opportunity Fund:** When a company hits 3x growth, you unlock bonus capital for follow-ons
- **Read the metrics:** Risk score and growth potential are hints, not guarantees
- **Management fees:** 2% annual fee chips away at your cash — don't sit on it too long

## Data Persistence

- All scores saved to `~/.config/unicorn/unicorn_scores.db` (or `~/Library/Application Support/unicorn/` on macOS)
- Database created automatically on first run
- Achievements and progression persist across games

Good luck! 🦄
# unicorn

text based startup adventure — play as VC and make bets on startups **OR** build your own startup and try to become the unicorn! 🦄

<img width="674" height="489" alt="Unicorn game splash screen" src="https://github.com/user-attachments/assets/99a8342a-98bf-44e6-b0db-399c2758a196" />

## Install

### Homebrew (macOS/Linux)

```bash
brew install james-see/tap/unicorn
```

> **Don't have Homebrew?** Install it first from [brew.sh](https://brew.sh), then run the command above in Terminal.

### Download Binary (No Tools Required)

Grab the latest pre-built binary from the [Releases page](https://github.com/james-see/unicorn/releases/latest) — look for the file matching your platform (e.g., `unicorn-darwin-arm64` for Apple Silicon Macs, `unicorn-darwin-amd64` for Intel Macs, `unicorn-linux-amd64` for Linux).

```bash
# After downloading, make it executable and run:
chmod +x unicorn
./unicorn
```

### From Source

Requires [Go 1.24+](https://go.dev/dl/):

```bash
go install github.com/jamesacampbell/unicorn@latest
```

## Nutshell

**TWO GAME MODES:**

### 🎩 VC Investor Mode (Classic)

You manage a realistic VC fund competing against AI investors like CARL from Sterling & Cooper! Deploy capital across startups, manage dilution through multiple funding rounds, pay 2% annual management fees, and use pro-rata rights to defend your ownership. With the new **Opportunity Fund** mechanic, breakout companies (3x+ growth) unlock additional capital for follow-on investments — just like real VCs raising SPVs. Over 5 years (60 turns), will you beat the AI VCs and become the top investor?

### 🚀 Startup Founder Mode

**Build your own startup from the ground up!** Choose from 10 different startup templates (SaaS, DeepTech, GovTech, Hardware), hire your team, acquire customers, raise funding rounds, and navigate the challenges of being a founder. Make strategic decisions about product development, marketing, hiring, partnerships, and more. Can you grow to $20M ARR and IPO? Or will you get acquired along the way?

## Features at a Glance

| | VC Mode | Founder Mode |
|---|---|---|
| **Goal** | Maximize fund ROI | Build & exit your startup |
| **Turns** | 60 (5 years) | 60 (5 years) |
| **Startups** | 45 to invest in | 10 templates to build |
| **AI Competition** | 7 AI VCs with strategies | Silicon Valley TV show competitors |
| **Funding** | Follow-ons, syndicates, opportunity fund | Seed → Series A → Series B |
| **Exits** | N/A (fund returns) | IPO, Acquisition, Secondary Sale |
| **Upgrades** | 10 VC upgrades | 8 founder upgrades |
| **Achievements** | 35 VC achievements | 10 founder achievements |
| **Progression** | 50 levels, cross-mode XP | Shared with VC mode |

## Difficulty Levels

| Level | Fund Size | Reserve | Opp Fund | Max New Bets | Event Chance | Volatility |
|---|---|---|---|---|---|---|
| **Easy** | $1M | $2.5M | $1.5M | 12 | 20% | 3% |
| **Medium** | $1.5M | $3M | $2.25M | 10 | 30% | 5% |
| **Hard** 🔒 | $2M | $3.5M | $2.5M | 8 | 40% | 7% |
| **Expert** 🔒 | $2.5M | $4M | $2.5M | 8 | 50% | 10% |

> 🔒 = unlocked at player level 5 (Hard) and level 10 (Expert)

## How to Play

```bash
# After installing via Homebrew or binary:
unicorn

# Or if building from source:
go build -o unicorn && ./unicorn
```

**Main Menu:**
1. New Game — Choose VC Mode or Founder Mode
2. Leaderboards — View top scores (local + global)
3. Player Statistics — Career stats
4. Achievements — 45+ achievements, 50 levels
5. Progression & Levels — XP, unlocks, rank titles
6. Analytics Dashboard — Performance trends
7. Upgrades — Permanent upgrades with achievement points
8. Help & Info — Complete in-game guide
9. Quit

**VC Mode Quick Start:** Select difficulty → browse 45 startups → invest → press `d` when done → watch your portfolio evolve over 60 turns. Follow-on opportunities appear when portfolio companies raise new rounds.

**Founder Mode Quick Start:** Choose a startup template → hire team → acquire customers → raise funding → build MRR → exit via IPO ($20M ARR), acquisition ($5M ARR), or secondary sale ($10M ARR).

## What's New

### v3.36.0 — Opportunity Fund (Latest)
- **Opportunity Fund**: separate capital pool unlocks when a portfolio company hits 3x growth, mirroring how real VCs raise SPVs for breakout winners
- Scales per difficulty: $1.5M (Easy) to $2.5M (Expert)
- Follow-on screen shows "⭐ Breakout qualified!" badge and available opp fund capital

### v3.35.0 — Spring Physics Animations
- Splash screen now uses Harmonica spring physics for smooth title fade-in, info box slide-in, and pulsing prompt
- New `AnimatedCounter` component: net worth and ROI spring-animate from 0 to final values on results screens

### v3.35.1 — Auto-Mode Fix
- Fixed: auto-mode tick no longer preempts follow-on investment input

### v3.34.0 — Board Seats & Capital Rebalance
- Founder mode: explicit board seat tracking for advisors and funding rounds
- VC mode: per-difficulty follow-on reserves, LP commitment multipliers, and first-check investment caps

### v3.33.0 — Founder/VC Mode Parity
- Wired broken menus, removed dead code, fixed emoji alignment issues

<details>
<summary>Older releases</summary>

### v3.26.0 — VC Firm Name Customization
- Set your own firm name (defaults to last name + "Capital")

### v3.25.0 — Employee Option Pool Dilution & ROI Predictor
- 15-20% dilution per funding round for employee options
- ROI predictor with confidence levels
- Customer referral program

### v3.20.0 — Player Progression System
- 50 levels with exponential XP and rank titles
- Level unlocks: Hard at L5, Expert at L10, Nightmare at L15
- Analytics dashboard with trend analysis

### v3.8.0 — Exit Options & Advisory Board
- IPO, Strategic Acquisition, or Secondary Sale
- Hire real Silicon Valley legends as advisors

### v3.7.0 — Founder Mode
- Build your own startup from the ground up
- Customer acquisition, MRR tracking, churn management
- Raise Seed, Series A, Series B

</details>

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md) for the complete feature roadmap.

**Quick Summary:**
- ✅ Opportunity fund, board seats, follow-ons, syndicates, global leaderboards
- 🚧 Planned: participating liquidation prefs, tournament mode, founder-to-VC crossover
- 📋 Full details and priorities in [docs/ROADMAP.md](docs/ROADMAP.md)

## Tech

- **Language:** Go 1.24
- **TUI Framework:** Charmbracelet (Bubble Tea, Lipgloss, Bubbles, Harmonica)
- **Database:** SQLite (local scores), Datasette/Vercel (global leaderboard)
- **License:** AGPL-3.0
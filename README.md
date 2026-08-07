<p align="center">
  <img src="./assets/readme/hero.svg" width="100%" alt="UNICORN — text-based VC and founder simulation">
</p>

<p align="center">
  <a href="https://james-see.github.io/unicorn/"><strong>Play site &amp; leaderboards</strong></a>
  ·
  <a href="https://github.com/james-see/unicorn/releases/latest">Releases</a>
  ·
  <a href="https://james-see.github.io/unicorn/#global-leaderboard">Live scores</a>
</p>

Text-based startup adventure — manage a VC fund **or** build the company and chase the exit.

<img width="674" height="489" alt="Unicorn TUI splash screen" src="https://github.com/user-attachments/assets/99a8342a-98bf-44e6-b0db-399c2758a196" />

## Install

### Homebrew (macOS / Linux)

```bash
brew install james-see/tap/unicorn
```

Then run `unicorn`.

### Binary

Grab the latest build from [Releases](https://github.com/james-see/unicorn/releases/latest) for your platform (`unicorn-darwin-arm64`, `unicorn-darwin-amd64`, `unicorn-linux-amd64`, or Windows).

```bash
chmod +x unicorn
./unicorn
```

### From source

Requires [Go 1.24+](https://go.dev/dl/):

```bash
go install github.com/jamesacampbell/unicorn@latest
```

## Nutshell

### VC Investor Mode

You manage a realistic VC fund against AI investors. Deploy capital across startups, manage dilution through funding rounds, pay 2% annual management fees, and use pro-rata rights to defend ownership. When a portfolio company hits 3× growth, the **Opportunity Fund** unlocks follow-on capital — like raising an SPV for a breakout. Over 60 turns (5 years), beat the AI VCs and climb the global board.

### Startup Founder Mode

Build from a template (SaaS, DeepTech, GovTech, Hardware). Hire the team, acquire customers, raise Seed → Series C, ship the roadmap, and navigate board pressure. Grow toward IPO, acquisition, or secondary sale.

## Modes at a glance

| | VC Mode | Founder Mode |
|---|---|---|
| **Goal** | Maximize fund ROI | Build & exit your startup |
| **Turns** | 60 (5 years) | 60 (5 years) |
| **Startups** | 45 to invest in | 10 templates to build |
| **AI competition** | 7 AI VCs with strategies | Silicon Valley–style rivals |
| **Funding** | Follow-ons, syndicates, opportunity fund | Seed → Series A → Series B → C |
| **Exits** | Fund returns | IPO, acquisition, secondary |
| **Upgrades** | 10 VC upgrades | 8 founder upgrades |
| **Achievements** | 35 VC | 10 founder |
| **Progression** | 50 levels, shared XP | Shared with VC mode |

## Difficulty (VC)

| Level | Fund Size | Reserve | Opp Fund | Max New Bets | Event Chance | Volatility |
|---|---|---|---|---|---|---|
| **Easy** | $1M | $2.5M | $1.5M | 12 | 20% | 3% |
| **Medium** | $1.5M | $3M | $2.25M | 10 | 30% | 5% |
| **Hard** (L5) | $2M | $3.5M | $2.5M | 8 | 40% | 7% |
| **Expert** (L10) | $2.5M | $4M | $2.5M | 8 | 50% | 10% |

Hard unlocks at player level 5; Expert at level 10.

## How to play

```bash
unicorn
```

**Main menu:** New Game (VC or Founder) · Leaderboards · Stats · Achievements · Progression · Analytics · Upgrades · Help · Quit

**VC quick start:** Select difficulty → browse 45 startups → invest → press `d` when done → watch the portfolio over 60 turns. Follow-ons appear when portfolio companies raise.

**Founder quick start:** Choose a template → hire → acquire customers → raise → build MRR → exit via IPO ($20M ARR), acquisition ($5M ARR), or secondary ($10M ARR).

## What's new

### v3.36.0 — Opportunity Fund

- Separate capital pool unlocks when a portfolio company hits 3× growth (SPV-style follow-ons)
- Scales per difficulty: $1.5M (Easy) to $2.5M (Expert)
- Follow-on screen shows breakout badge and available opp-fund capital

<details>
<summary>Older releases</summary>

### v3.35.0 — Spring physics animations
- Splash screen Harmonica springs; `AnimatedCounter` on results screens

### v3.35.1 — Auto-mode fix
- Auto-mode tick no longer preempts follow-on investment input

### v3.34.0 — Board seats & capital rebalance
- Founder: explicit board seat tracking; VC: per-difficulty reserves and first-check caps

### v3.33.0 — Founder/VC mode parity
- Wired broken menus, removed dead code, fixed emoji alignment

### v3.26.0 — VC firm name customization
### v3.25.0 — Option pool dilution & ROI predictor
### v3.20.0 — Player progression (50 levels)
### v3.8.0 — Exit options & advisory board
### v3.7.0 — Founder Mode

Full history: [CHANGELOG.md](CHANGELOG.md)

</details>

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md).

- Done: opportunity fund, board seats, follow-ons, syndicates, global leaderboards
- Planned: participating prefs, tournament mode, founder-to-VC crossover

## Tech

- **Language:** Go 1.24
- **TUI:** Charmbracelet (Bubble Tea, Lipgloss, Bubbles, Harmonica)
- **Data:** SQLite (local), Datasette/Vercel (global leaderboard)
- **License:** [AGPL-3.0](LICENSE)
- **Site:** [james-see.github.io/unicorn](https://james-see.github.io/unicorn/)

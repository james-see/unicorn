# Upgrade & Unlock System Proposal

## Overview
Players earn achievement points (currently 5-100 per achievement) that unlock permanent upgrades, new game modes, enhanced features, and special abilities. This creates meta-progression similar to games like **Cookie Clicker**, **Idle Games**, and **Roguelikes**.

---

## 🎮 Unlock Categories

### 1. 💼 Investment Terms Upgrades
**Spend points to unlock better investment terms permanently**

| Upgrade | Cost | Effect |
|---------|------|--------|
| **Enhanced SAFE Discount** | 150 pts | SAFE discount increases from 20% → 25% |
| **Double Board Seat** | 200 pts | Get 2 board seats per $100k investment (double voting power) |
| **2x Liquidation Preference** | 250 pts | Unlock 2x liquidation preference option (get paid 2x before others) |
| **Weighted Voting Rights** | 300 pts | Your board votes count as 1.5x on acquisitions |
| **Super Pro-Rata** | 200 pts | Can invest up to 50% of round (vs 20% max) |
| **Early Access** | 100 pts | See 2 extra startups before investment phase starts |
| **Founder-Friendly Terms** | 150 pts | Startups more likely to accept your investments (better deals) |

**Why It's Fun**: Gives players strategic choices and permanent advantages that change how they play.

---

### 2. 🏛️ Board Seat Powers
**Enhanced voting and board influence**

| Upgrade | Cost | Effect |
|---------|------|--------|
| **Board Majority** | 400 pts | Your vote counts as 2 votes (need 2 board seats) |
| **Veto Power** | 500 pts | Can veto one down round per game |
| **Founder Alignment** | 300 pts | Founders vote with you 60% of the time |
| **Strategic Advisor** | 250 pts | Get preview of next board vote before it happens |
| **Board Network** | 350 pts | Learn about acquisitions 1 turn earlier |

**Why It's Fun**: Makes board seats more powerful and strategic, not just cosmetic.

---

### 3. 💰 Financial Perks
**Starting advantages and economic bonuses**

| Upgrade | Cost | Effect |
|---------|------|--------|
| **Fund Booster** | 100 pts | +10% starting cash on all difficulties |
| **Management Fee Reduction** | 150 pts | Management fees reduced from 2% → 1.5% |
| **Follow-On Reserve Boost** | 200 pts | +$200k to follow-on reserve |
| **Angel Investor** | 250 pts | Start with $100k bonus cash (stacks) |
| **Fee Waiver** | 300 pts | No management fees for first 12 months |
| **Seed Accelerator** | 400 pts | First investment gets 25% equity bonus |

**Why It's Fun**: Helps players progress faster and customize their starting conditions.

---

### 4. 🎯 Startup Access & Information
**Better intel and access to deals**

| Upgrade | Cost | Effect |
|---------|------|--------|
| **Due Diligence** | 150 pts | See risk score numbers (not just labels) |
| **Revenue Tracker** | 100 pts | See monthly revenue growth trends |
| **Founder Network** | 200 pts | See one extra startup per round |
| **Market Intelligence** | 250 pts | See which sectors are trending up/down |
| **Founder Secrets** | 300 pts | Know if acquisition is coming 3 turns early |
| **Portfolio Synergy** | 350 pts | Investments in same sector get +5% valuation boost |

**Why It's Fun**: Information is power - helps players make better decisions.

---

### 5. 🎲 Game Mode Unlocks
**New ways to play the game**

| Mode | Cost | Description |
|------|------|-------------|
| **Speed Mode** | 200 pts | 30 turns instead of 60 (faster games) |
| **Endurance Mode** | 250 pts | 120 turns instead of 60 (longer games) |
| **Hardcore Mode** | 300 pts | Only 1 investment allowed, Expert difficulty |
| **Sandbox Mode** | 400 pts | Custom difficulty settings (adjust volatility, events) |
| **Career Mode** | 500 pts | Play multiple games, carry over reputation |
| **Co-Investment Mode** | 350 pts | Partner with AI funds on deals (split equity) |
| **Founder Selection** | 400 pts | Choose which startups appear (curated selection) |

**Why It's Fun**: Adds replayability and different challenges for experienced players.

---

### 6. 🤖 AI Opponent Unlocks
**New AI personalities and strategies**

| Opponent | Cost | Strategy |
|----------|------|----------|
| **Tiger Global** | 150 pts | Aggressive growth investor (high valuations) |
| **Y Combinator** | 200 pts | Early-stage specialist (low valuations, high equity) |
| **SoftBank** | 250 pts | Mega-fund (massive investments) |
| **Angel Investor** | 100 pts | Solo investor (smaller deals, more portfolio) |
| **Corporate VC** | 200 pts | Strategic investor (sector-focused) |
| **Impact Investor** | 150 pts | Only invests in CleanTech/AgriTech |

**Why It's Fun**: Different AI strategies create different competitive dynamics each game.

---

### 7. 📊 Analysis Tools
**Better analytics and insights**

| Tool | Cost | Effect |
|------|------|--------|
| **Portfolio Dashboard** | 150 pts | See detailed portfolio metrics each turn |
| **ROI Predictor** | 200 pts | See projected ROI for each investment |
| **Risk Calculator** | 150 pts | See probability of different outcomes |
| **Market Trends** | 250 pts | Visualize sector trends over time |
| **Exit Simulator** | 300 pts | Predict when companies might exit |

**Why It's Fun**: Appeals to players who love data and optimization.

---

### 8. 🎨 Cosmetic & Quality of Life
**Visual and UX improvements**

| Upgrade | Cost | Effect |
|---------|------|--------|
| **Color Themes** | 50 pts | Unlock new color schemes |
| **Achievement Badges** | 25 pts | Show badges next to name in leaderboard |
| **Advanced Stats** | 100 pts | Detailed post-game statistics |
| **Export Portfolio** | 150 pts | Export game results to CSV |
| **Replay Mode** | 200 pts | Watch replay of your game decisions |

**Why It's Fun**: Personalization and convenience features.

---

### 9. 🔥 Special Abilities
**Unique powers that change gameplay**

| Ability | Cost | Effect |
|---------|------|--------|
| **Time Machine** | 600 pts | Rewind one investment decision per game |
| **Portfolio Insurance** | 500 pts | Protect one investment from down rounds |
| **Valuation Boost** | 400 pts | +10% valuation boost to one company per game |
| **Exit Accelerator** | 450 pts | Force one company to exit early (2x return) |
| **Network Effect** | 350 pts | Portfolio companies get +10% growth from synergy |
| **Lucky Streak** | 300 pts | One random positive event guaranteed per game |

**Why It's Fun**: Powerful abilities that create memorable moments and strategic depth.

---

## 🎯 Implementation Strategy

### Phase 1: Core System
1. Add `player_upgrades` table to database
2. Create upgrade menu in main menu
3. Implement 3-5 most impactful upgrades (Financial Perks, Investment Terms)
4. Track points spending

### Phase 2: Game Mode Expansions
1. Add new game modes
2. Unlockable AI opponents
3. Sandbox/custom difficulty

### Phase 3: Advanced Features
1. Board seat powers
2. Special abilities
3. Analysis tools

---

## 💡 Similar Games Inspiration

### **Cookie Clicker** (Idle Games)
- Unlock upgrades that persist across games
- Meta-progression feels rewarding
- Players can focus on different upgrade paths

### **Slay the Spire** (Roguelikes)
- Unlock new cards/characters
- Each unlock opens new strategies
- Meta-progression doesn't make game easier, just different

### **Football Manager** (Simulation)
- Unlock new tactics and abilities
- Career progression feels meaningful
- Depth without complexity

### **Poker** (Strategy)
- Skill progression feels rewarding
- Long-term improvements matter
- Meta-game around optimization

---

## 🎮 Example Player Journey

**New Player** (0 points):
- Plays first game, earns 25 points
- Unlocks "Fund Booster" (+10% cash)
- Next game feels easier, earns 50 points
- Unlocks "Due Diligence" (see risk scores)
- Makes better decisions, earns 100 points
- Unlocks "Enhanced SAFE Discount"
- Strategy evolves - now uses SAFE more
- Earns 200 points, unlocks "Speed Mode"
- Plays faster games, earns points faster
- Unlocks "Double Board Seat"
- Gameplay completely changes - focuses on board control

**Veteran Player** (2000+ points):
- Has unlocked most upgrades
- Plays Hardcore Mode for challenge
- Uses Time Machine for perfect runs
- Focuses on leaderboard rankings
- Creates custom strategies with all tools

---

## 🏆 Suggested Point Costs

**Tier System**:
- **Tier 1 (50-100 pts)**: Small QoL improvements
- **Tier 2 (100-200 pts)**: Medium gameplay changes
- **Tier 3 (200-350 pts)**: Major strategic unlocks
- **Tier 4 (350-500 pts)**: Game mode unlocks
- **Tier 5 (500+ pts)**: Special abilities

**Balance**: Players should unlock 1-2 upgrades per game initially, then slower as costs increase.

---

## 📝 User Experience

### Upgrade Menu (Main Menu → "5. Upgrades")
```
╔════════════════════════════════════════════════╗
║           🎁 UPGRADE STORE 🎁                  ║
╚════════════════════════════════════════════════╝

Your Points: 425
Career Level: 4 - Principal

╔════════════════════════════════════════════════╗
║  💼 INVESTMENT TERMS                           ║
╚════════════════════════════════════════════════╝

[✓] Enhanced SAFE Discount (150 pts) - OWNED
[✓] Early Access (100 pts) - OWNED
[ ] Double Board Seat (200 pts) - AVAILABLE
[ ] Super Pro-Rata (200 pts) - AVAILABLE
[ ] 2x Liquidation Preference (250 pts) - AVAILABLE

╔════════════════════════════════════════════════╗
║  💰 FINANCIAL PERKS                            ║
╚════════════════════════════════════════════════╝

[✓] Fund Booster (100 pts) - OWNED
[ ] Management Fee Reduction (150 pts) - AVAILABLE
[ ] Follow-On Reserve Boost (200 pts) - AVAILABLE

Press number to purchase, or 'q' to quit...
```

---

## 🎯 Why This Works

1. **Meta-Progression**: Players feel like they're getting better over time
2. **Strategic Depth**: Different upgrade paths = different playstyles
3. **Replayability**: Each unlock opens new strategies to try
4. **Reward System**: Points feel meaningful, not just cosmetic
5. **Flexibility**: Players can focus on what they enjoy
6. **Longevity**: Hundreds of hours of content through unlocks

---

## 🚀 Quick Wins (Easy to Implement)

1. **Fund Booster** (+10% cash) - Simple multiplier
2. **Enhanced SAFE Discount** - Change one number
3. **Early Access** - Show 2 extra startups
4. **Due Diligence** - Show risk score numbers
5. **Speed Mode** - Change max turns

These 5 upgrades alone would add significant value and engagement!


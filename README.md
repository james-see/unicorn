# unicorn
text based startup adventure - play as vc and make bets on startups or try to be the unicorn

<pre>
               \
                \\
                 \%,     ,'     , ,.
                  \%\,';/J,";";";;,,.
     ~.------------\%;((`);)));`;;,.,-----------,~
    ~~:           ,`;@)((;`,`((;(;;);;,`         :~~
   ~~ :           ;`(@```))`~ ``; );(;));;,      : ~~
  ~~  :            `X `(( `),    (;;);;;;`       :  ~~
 ~~~~ :            / `) `` /;~   `;;;;;;;);,     :  ~~~~
~~~~  :           / , ` ,/` /     (`;;(;;;;,     : ~~~~
  ~~~ :          (o  /]_/` /     ,);;;`;;;;;`,,  : ~~~
   ~~ :           `~` `~`  `      ``;,  ``;" ';, : ~~
    ~~:                             `'   `'  `'  :~~
     ~`-----------------------------------------`~
       ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐
                  WELCOME TO UNICORN
       │            COPYRIGHT 2019             │
        ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
</pre>

## Install

`go get github.com/jamesacampbell/unicorn`

## Nutshell

You start out as a VC with $250,000 USD. Your goal is to invest and make as much money as you can by the end of 10 years (120 turns). Each turn represents 1 month - random events affect your portfolio companies, and you'll see your investments grow or decline. Make strategic decisions about which startups to back based on their risk profiles and growth potential.

## Features

### 🎮 Core Gameplay
- **Investment Mechanics:** Invest any amount in 20 different startups (NEW!)
- **Portfolio Tracking:** Real-time valuation of all your investments
- **Turn-Based System:** 90-120 turns depending on difficulty
- **Random Events:** 60+ events that impact company valuations (NEW!)
- **Performance Ratings:** From "Unicorn Hunter" (1000%+ ROI) to "Lost Money"
- **Advanced Analytics:** Sector breakdown, best/worst performers (NEW!)

### 🏆 Difficulty Levels (NEW!)
- **Easy:** $500k starting cash, 20% event chance, 3% volatility
- **Medium:** $250k starting cash, 30% event chance, 5% volatility  
- **Hard:** $150k starting cash, 40% event chance, 7% volatility
- **Expert:** $100k starting cash, 50% event chance, 10% volatility, only 90 turns!

### 📊 Persistence & Competition (NEW!)
- **Leaderboards:** Track top 10 scores by net worth or ROI
- **Statistics:** View career stats for any player
- **Recent Games:** See the last 10 games played
- **SQLite Database:** All scores saved locally
- **Difficulty Filters:** Separate leaderboards for each difficulty

### 🏢 20 Diverse Startups (NEW!)
Choose from companies across 12+ sectors:
- **FinTech:** AI-powered trading platforms
- **BioTech:** Nanotechnology drug delivery  
- **CleanTech:** Sustainable packaging, food waste conversion
- **HealthTech:** VR meditation apps
- **EdTech:** TikTok-style education
- **Robotics:** Automated kitchen systems
- **Security:** Blockchain & quantum-resistant encryption
- **Gaming:** Cloud gaming platforms
- **LegalTech:** AI legal document automation
- **AgriTech:** Vertical farming kits
- **Logistics:** Last-mile delivery drones
- **IoT:** Smart home control hubs
- **Creative:** AI-generated music
- **CloudTech:** Infrastructure optimization
- **Social Media:** Pet social networks
- **Advertising:** Unicycle billboards
- **Consumer Goods:** IoT finger puppets, pet umbrellas

### 📈 Strategic Depth
- **Risk Indicators:** See which companies are high/medium/low risk
- **Growth Potential:** Evaluate each startup's growth prospects
- **Sector Diversity:** Spread your bets across different industries
- **Capital Allocation:** Balance investment vs. cash reserves
- **Difficulty Selection:** Choose your challenge level

### 🎯 Scoring System
- Net worth calculation (cash + portfolio value)
- ROI (Return on Investment) percentage
- Successful exits counter (5x+ returns)
- Performance tier ratings (6 levels)
- Persistent leaderboards
- Career statistics

## How to Play

```bash
# Build the game
go build -o unicorn

# Run it
./unicorn

# Menu options:
# 1. New Game - Start a new investment game
# 2. Leaderboards - View top scores
# 3. Player Statistics - Check your career stats
# 4. Help & Info - Complete game guide (NEW!)
# 5. Quit
```

## What's New in Phase 3

✅ **20 Startup Companies** (doubled!)
✅ **60 Random Events** (doubled!)
✅ **Advanced Analytics System** with sector breakdown
✅ **Help & Information Menu** with complete guide
✅ **12+ Industry Sectors** for strategic diversification

## Future Plans

**Phase 4:** Achievements system, challenge modes, career progression
**Long-term:** Game mode where you run a startup trying to become a Unicorn

## Demo

![unicorn-demo 2019-01-02 16_21_43](https://user-images.githubusercontent.com/616585/50613136-95163300-0eaa-11e9-9e0b-a4ed7c57bc71.gif)

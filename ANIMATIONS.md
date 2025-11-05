# 🎨 Unicorn Game Animations

## Overview
The Unicorn game now features beautiful CLI animations powered by the [pterm](https://github.com/pterm/pterm) library!

## Features Added

### 1. **Splash Screen Animation** 🚀
- Big animated UNICORN title on startup
- Loading spinner with "Loading your portfolio..." message
- Styled info box with game description

### 2. **Achievement Unlock Animations** 🏆
- Flashy achievement notification with stars and sparkles
- Animated box display showing achievement name and description
- Timed display with visual effects

### 3. **Round Transition Animations** 🎯
- Header animation every 5 turns showing milestone progress
- Color-coded round numbers
- Full-width header display

### 4. **Game Over Animations** 💎
- Victory animation with big "VICTORY" text (when you win)
- Game Over animation with big "GAME OVER" text (when you lose)
- Styled boxes showing final net worth
- Color-coded based on win/loss

### 5. **Helper Functions** ⚡
Additional animation utilities available:
- `ShowLoadingSpinner()` - Customizable loading animations
- `ShowSuccessMessage()` - Green checkmark success messages
- `ShowErrorMessage()` - Red error notifications
- `ShowWarningMessage()` - Yellow warning alerts
- `ShowInfoMessage()` - Blue info messages
- `ShowInvestmentAnimation()` - Investment processing animation
- `ShowExitAnimation()` - Successful exit celebration with fireworks
- `TypewriterEffect()` - Character-by-character text typing
- `ShowProgressBar()` - Animated progress bars

## Technical Details

### Package Structure
```
animations/
  └── animations.go  - All animation functions
```

### Dependencies
- **pterm v0.12.82** - Modern terminal animation library
- Provides spinners, progress bars, big text, and styled boxes

### Integration Points
1. **main()** - Splash screen on game start
2. **playTurn()** - Round transitions every 5 turns
3. **displayFinalScore()** - Animated game over screen
4. **checkAndUnlockAchievements()** - Achievement unlock animations

## Usage Examples

### Basic Spinner
```go
animations.ShowLoadingSpinner("Processing...", 2*time.Second)
```

### Success Message
```go
animations.ShowSuccessMessage("Investment successful!")
```

### Achievement Unlock
```go
animations.ShowAchievementUnlock("🏆 High Roller", "Invested over $1M in a single round")
```

## Animation Timing
- Splash screen: ~2 seconds
- Achievement unlocks: ~1.5 seconds each
- Round transitions: ~0.5 seconds
- Game over: ~2 seconds

## Customization
All animations can be customized by modifying `animations/animations.go`. The pterm library offers extensive customization options including:
- Custom colors and styles
- Different spinner sequences
- Box border styles
- Text alignment
- Animation speeds

## Future Enhancement Ideas
- [ ] Animated leaderboard reveals
- [ ] Investment decision countdown timer
- [ ] Portfolio value change animations
- [ ] Startup funding round progress bars
- [ ] Real-time market event notifications with effects
- [ ] Multiplayer waiting room animations
- [ ] Achievement progress tracking bars

---
*Built with ❤️ using pterm*


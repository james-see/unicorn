package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// VCGameScreen is a wrapper that delegates to VCTurnScreen
// This exists for backward compatibility with the screen enum
type VCGameScreen struct {
	turnScreen *VCTurnScreen
}

// NewVCGameScreen creates a new VC game screen
func NewVCGameScreen(width, height int, gameData *GameData) *VCGameScreen {
	return &VCGameScreen{
		turnScreen: NewVCTurnScreen(width, height, gameData),
	}
}

// Init initializes the game screen
func (g *VCGameScreen) Init() tea.Cmd {
	return g.turnScreen.Init()
}

// Update handles game screen input
func (g *VCGameScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	return g.turnScreen.Update(msg)
}

// View renders the game screen
func (g *VCGameScreen) View() string {
	return g.turnScreen.View()
}

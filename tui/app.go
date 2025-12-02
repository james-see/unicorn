package tui

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/founder"
	"github.com/jamesacampbell/unicorn/game"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// Screen represents different screens in the application
type Screen int

const (
	ScreenSplash Screen = iota
	ScreenMainMenu
	ScreenVCSetup
	ScreenVCGame
	ScreenVCInvest
	ScreenVCTurn
	ScreenVCResults
	ScreenFounderSetup
	ScreenFounderGame
	ScreenLeaderboard
	ScreenAchievements
	ScreenUpgrades
	ScreenStats
	ScreenProgression
	ScreenAnalytics
	ScreenHelp
)

// Global key bindings
type appKeyMap struct {
	Quit key.Binding
	Back key.Binding
	Help key.Binding
}

var appKeys = appKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// Messages for screen transitions
type (
	// SwitchScreenMsg requests a screen change
	SwitchScreenMsg struct {
		Screen Screen
		Data   interface{}
	}

	// PushScreenMsg pushes a screen onto the stack
	PushScreenMsg struct {
		Screen Screen
		Data   interface{}
	}

	// PopScreenMsg pops the current screen and goes back
	PopScreenMsg struct{}

	// QuitMsg signals the app to quit
	QuitMsg struct{}

	// ErrorMsg carries error information
	ErrorMsg struct {
		Err error
	}

	// WindowSizeMsg carries terminal size
	WindowSizeMsg struct {
		Width  int
		Height int
	}
)

// GameData holds shared game state between screens
type GameData struct {
	PlayerName     string
	FirmName       string
	Difficulty     game.Difficulty
	GameState      *game.GameState
	FounderState   *founder.FounderState
	PlayerUpgrades []string
	AutoMode       bool
}

// ScreenModel interface for all screen models
type ScreenModel interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (ScreenModel, tea.Cmd)
	View() string
}

// App is the root model for the entire application
type App struct {
	width        int
	height       int
	currentScreen Screen
	screenStack  []Screen
	gameData     *GameData
	
	// Screen models
	splash      ScreenModel
	mainMenu    ScreenModel
	vcSetup     ScreenModel
	vcGame      ScreenModel
	vcInvest    ScreenModel
	vcTurn      ScreenModel
	vcResults   ScreenModel
	founderSetup ScreenModel
	founderGame  ScreenModel
	leaderboard  ScreenModel
	achievements ScreenModel
	upgrades     ScreenModel
	stats        ScreenModel
	progression  ScreenModel
	analytics    ScreenModel
	help         ScreenModel
	
	quitting    bool
	showHelp    bool
}

// NewApp creates a new application instance
func NewApp() *App {
	app := &App{
		width:         80,
		height:        24,
		currentScreen: ScreenSplash,
		screenStack:   []Screen{},
		gameData:      &GameData{},
	}
	
	// Initialize screens - they will be created lazily
	return app
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	// Get user config directory (~/.config/unicorn)
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("HOME")
	}
	unicornDir := filepath.Join(configDir, "unicorn")
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(unicornDir, 0755); err == nil {
		dbPath := filepath.Join(unicornDir, "unicorn_scores.db")
		database.InitDB(dbPath)
	}
	
	// Start with splash screen
	a.splash = NewSplashScreen(a.width, a.height)
	return a.splash.Init()
}

// Update handles messages
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit handler
		if key.Matches(msg, appKeys.Quit) && a.currentScreen != ScreenVCGame && a.currentScreen != ScreenFounderGame {
			a.quitting = true
			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		// Propagate to current screen
		return a, nil

	case SwitchScreenMsg:
		return a.switchScreen(msg.Screen, msg.Data)

	case PushScreenMsg:
		a.screenStack = append(a.screenStack, a.currentScreen)
		return a.switchScreen(msg.Screen, msg.Data)

	case PopScreenMsg:
		if len(a.screenStack) > 0 {
			prevScreen := a.screenStack[len(a.screenStack)-1]
			a.screenStack = a.screenStack[:len(a.screenStack)-1]
			return a.switchScreen(prevScreen, nil)
		}
		// If no stack, go to main menu
		return a.switchScreen(ScreenMainMenu, nil)

	case QuitMsg:
		a.quitting = true
		return a, tea.Quit

	case ErrorMsg:
		// Handle errors - could show a dialog
		return a, nil
	}

	// Update current screen
	var cmd tea.Cmd
	switch a.currentScreen {
	case ScreenSplash:
		if a.splash != nil {
			a.splash, cmd = a.splash.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenMainMenu:
		if a.mainMenu != nil {
			a.mainMenu, cmd = a.mainMenu.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenVCSetup:
		if a.vcSetup != nil {
			a.vcSetup, cmd = a.vcSetup.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenVCGame:
		if a.vcGame != nil {
			a.vcGame, cmd = a.vcGame.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenVCInvest:
		if a.vcInvest != nil {
			a.vcInvest, cmd = a.vcInvest.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenVCTurn:
		if a.vcTurn != nil {
			a.vcTurn, cmd = a.vcTurn.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenVCResults:
		if a.vcResults != nil {
			a.vcResults, cmd = a.vcResults.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenFounderSetup:
		if a.founderSetup != nil {
			a.founderSetup, cmd = a.founderSetup.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenFounderGame:
		if a.founderGame != nil {
			a.founderGame, cmd = a.founderGame.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenLeaderboard:
		if a.leaderboard != nil {
			a.leaderboard, cmd = a.leaderboard.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenAchievements:
		if a.achievements != nil {
			a.achievements, cmd = a.achievements.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenUpgrades:
		if a.upgrades != nil {
			a.upgrades, cmd = a.upgrades.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenStats:
		if a.stats != nil {
			a.stats, cmd = a.stats.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenProgression:
		if a.progression != nil {
			a.progression, cmd = a.progression.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenAnalytics:
		if a.analytics != nil {
			a.analytics, cmd = a.analytics.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenHelp:
		if a.help != nil {
			a.help, cmd = a.help.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return a, tea.Batch(cmds...)
}

// View renders the current screen
func (a *App) View() string {
	if a.quitting {
		return styles.Center(a.width, "Thanks for playing! 🦄")
	}

	var content string
	switch a.currentScreen {
	case ScreenSplash:
		if a.splash != nil {
			content = a.splash.View()
		}
	case ScreenMainMenu:
		if a.mainMenu != nil {
			content = a.mainMenu.View()
		}
	case ScreenVCSetup:
		if a.vcSetup != nil {
			content = a.vcSetup.View()
		}
	case ScreenVCGame:
		if a.vcGame != nil {
			content = a.vcGame.View()
		}
	case ScreenVCInvest:
		if a.vcInvest != nil {
			content = a.vcInvest.View()
		}
	case ScreenVCTurn:
		if a.vcTurn != nil {
			content = a.vcTurn.View()
		}
	case ScreenVCResults:
		if a.vcResults != nil {
			content = a.vcResults.View()
		}
	case ScreenFounderSetup:
		if a.founderSetup != nil {
			content = a.founderSetup.View()
		}
	case ScreenFounderGame:
		if a.founderGame != nil {
			content = a.founderGame.View()
		}
	case ScreenLeaderboard:
		if a.leaderboard != nil {
			content = a.leaderboard.View()
		}
	case ScreenAchievements:
		if a.achievements != nil {
			content = a.achievements.View()
		}
	case ScreenUpgrades:
		if a.upgrades != nil {
			content = a.upgrades.View()
		}
	case ScreenStats:
		if a.stats != nil {
			content = a.stats.View()
		}
	case ScreenProgression:
		if a.progression != nil {
			content = a.progression.View()
		}
	case ScreenAnalytics:
		if a.analytics != nil {
			content = a.analytics.View()
		}
	case ScreenHelp:
		if a.help != nil {
			content = a.help.View()
		}
	default:
		content = "Loading..."
	}

	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// switchScreen handles screen transitions
func (a *App) switchScreen(screen Screen, data interface{}) (tea.Model, tea.Cmd) {
	a.currentScreen = screen
	var cmd tea.Cmd

	switch screen {
	case ScreenSplash:
		a.splash = NewSplashScreen(a.width, a.height)
		cmd = a.splash.Init()

	case ScreenMainMenu:
		a.mainMenu = NewMainMenuScreen(a.width, a.height, a.gameData)
		cmd = a.mainMenu.Init()

	case ScreenVCSetup:
		a.vcSetup = NewVCSetupScreen(a.width, a.height, a.gameData)
		cmd = a.vcSetup.Init()

	case ScreenVCGame:
		a.vcGame = NewVCGameScreen(a.width, a.height, a.gameData)
		cmd = a.vcGame.Init()

	case ScreenVCInvest:
		a.vcInvest = NewVCInvestScreen(a.width, a.height, a.gameData)
		cmd = a.vcInvest.Init()

	case ScreenVCTurn:
		a.vcTurn = NewVCTurnScreen(a.width, a.height, a.gameData)
		cmd = a.vcTurn.Init()

	case ScreenVCResults:
		a.vcResults = NewVCResultsScreen(a.width, a.height, a.gameData)
		cmd = a.vcResults.Init()

	case ScreenFounderSetup:
		a.founderSetup = NewFounderSetupScreen(a.width, a.height, a.gameData)
		cmd = a.founderSetup.Init()

	case ScreenFounderGame:
		a.founderGame = NewFounderGameScreen(a.width, a.height, a.gameData)
		cmd = a.founderGame.Init()

	case ScreenLeaderboard:
		a.leaderboard = NewLeaderboardScreen(a.width, a.height)
		cmd = a.leaderboard.Init()

	case ScreenAchievements:
		a.achievements = NewAchievementsScreen(a.width, a.height, a.gameData.PlayerName)
		cmd = a.achievements.Init()

	case ScreenUpgrades:
		a.upgrades = NewUpgradesScreen(a.width, a.height, a.gameData.PlayerName)
		cmd = a.upgrades.Init()

	case ScreenStats:
		a.stats = NewStatsScreen(a.width, a.height)
		cmd = a.stats.Init()

	case ScreenProgression:
		a.progression = NewProgressionScreen(a.width, a.height, a.gameData.PlayerName)
		cmd = a.progression.Init()

	case ScreenAnalytics:
		a.analytics = NewAnalyticsScreen(a.width, a.height)
		cmd = a.analytics.Init()

	case ScreenHelp:
		a.help = NewHelpScreen(a.width, a.height)
		cmd = a.help.Init()
	}

	return a, cmd
}

// Run starts the Bubble Tea program
func Run() error {
	p := tea.NewProgram(
		NewApp(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	
	_, err := p.Run()
	
	// Clean up database on exit
	database.CloseDB()
	
	return err
}

// Helper functions for creating messages

// SwitchTo creates a command to switch to a screen
func SwitchTo(screen Screen) tea.Cmd {
	return func() tea.Msg {
		return SwitchScreenMsg{Screen: screen}
	}
}

// SwitchToWithData creates a command to switch to a screen with data
func SwitchToWithData(screen Screen, data interface{}) tea.Cmd {
	return func() tea.Msg {
		return SwitchScreenMsg{Screen: screen, Data: data}
	}
}

// PushTo creates a command to push a screen onto the stack
func PushTo(screen Screen) tea.Cmd {
	return func() tea.Msg {
		return PushScreenMsg{Screen: screen}
	}
}

// PopScreen creates a command to pop the screen stack
func PopScreen() tea.Cmd {
	return func() tea.Msg {
		return PopScreenMsg{}
	}
}

// Quit creates a command to quit the app
func Quit() tea.Cmd {
	return func() tea.Msg {
		return QuitMsg{}
	}
}

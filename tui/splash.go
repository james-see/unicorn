package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Inline styles for splash screen
var (
	splashMagenta = lipgloss.Color("#FF00FF")
	splashCyan    = lipgloss.Color("#00FFFF")
	splashYellow  = lipgloss.Color("#FFFF00")
	splashGray    = lipgloss.Color("#808080")
)

// ASCII art unicorn logo
const unicornLogo = `
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣤⣄⡀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⣿⣿⣿⣿⣷⡀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆⠀
    ⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀
    ⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⠁⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⡟⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⡿⠋⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠁⠀⠀⠀⠀⠀
`

const unicornTitle = `
██╗   ██╗███╗   ██╗██╗ ██████╗ ██████╗ ██████╗ ███╗   ██╗
██║   ██║████╗  ██║██║██╔════╝██╔═══██╗██╔══██╗████╗  ██║
██║   ██║██╔██╗ ██║██║██║     ██║   ██║██████╔╝██╔██╗ ██║
██║   ██║██║╚██╗██║██║██║     ██║   ██║██╔══██╗██║╚██╗██║
╚██████╔╝██║ ╚████║██║╚██████╗╚██████╔╝██║  ██║██║ ╚████║
 ╚═════╝ ╚═╝  ╚═══╝╚═╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝
`

// SplashScreen is the animated splash screen
type SplashScreen struct {
	width      int
	height     int
	spinner    spinner.Model
	phase      int
	startTime  time.Time
	ready      bool
}

// NewSplashScreen creates a new splash screen
func NewSplashScreen(width, height int) *SplashScreen {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"🦄", "✨🦄", "✨✨🦄", "✨🦄", "🦄"},
		FPS:    time.Second / 6,
	}
	s.Style = lipgloss.NewStyle()

	return &SplashScreen{
		width:     width,
		height:    height,
		spinner:   s,
		phase:     0,
		startTime: time.Now(),
		ready:     false,
	}
}

// splashTickMsg is sent to advance the splash animation
type splashTickMsg struct{}

func splashTick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return splashTickMsg{}
	})
}

// Init initializes the splash screen
func (s *SplashScreen) Init() tea.Cmd {
	return tea.Batch(s.spinner.Tick, splashTick())
}

// Update handles splash screen updates
func (s *SplashScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Any key press advances to main menu (intentional first-key-skip)
		return s, SwitchTo(ScreenMainMenu)

	case splashTickMsg:
		s.phase++
		if s.phase >= 4 {
			s.ready = true
		}
		return s, splashTick()

	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}

	return s, nil
}

// View renders the splash screen
func (s *SplashScreen) View() string {
	var b strings.Builder

	// Title with gradient effect
	titleStyle := lipgloss.NewStyle().
		Foreground(splashMagenta).
		Bold(true)
	
	b.WriteString(titleStyle.Render(unicornTitle))
	b.WriteString("\n")

	// Spinner
	spinnerStyle := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(spinnerStyle.Render(s.spinner.View()))
	b.WriteString("\n\n")

	// Tagline with fade-in effect based on phase
	taglines := []string{
		"The Ultimate Startup Adventure",
		"Invest as a VC, or Build as a Founder",
		"Hunt for Unicorns, or Become One",
		"Where Legends Are Made",
	}

	taglineStyle := lipgloss.NewStyle().
		Foreground(splashCyan).
		Italic(true).
		Width(s.width).
		Align(lipgloss.Center)

	if s.phase < len(taglines) {
		b.WriteString(taglineStyle.Render(taglines[s.phase]))
	} else {
		b.WriteString(taglineStyle.Render("Where Dreams Become Unicorns... Or Don't"))
	}
	b.WriteString("\n\n")

	// Box with game info
	if s.phase >= 2 {
		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(splashCyan).
			Padding(1, 2).
			Width(50).
			Align(lipgloss.Center)

		info := "🚀 Invest in startups\n💰 Compete against AI VCs\n🏆 Build your reputation\n🦄 Find the next unicorn"
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(infoBox.Render(info)))
		b.WriteString("\n\n")
	}

	// Press any key prompt
	if s.ready {
		promptStyle := lipgloss.NewStyle().
			Foreground(splashYellow).
			Bold(true).
			Blink(true).
			Width(s.width).
			Align(lipgloss.Center)
		b.WriteString(promptStyle.Render("Press any key to continue..."))
	}

	// Copyright
	b.WriteString("\n\n")
	copyrightStyle := lipgloss.NewStyle().
		Foreground(splashGray).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(copyrightStyle.Render("© 2025 Unicorn Game"))

	return b.String()
}

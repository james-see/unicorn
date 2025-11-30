package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// LoadingSpinner is a spinner with a message
type LoadingSpinner struct {
	spinner spinner.Model
	message string
	width   int
}

// NewLoadingSpinner creates a new loading spinner
func NewLoadingSpinner(message string) *LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle
	
	return &LoadingSpinner{
		spinner: s,
		message: message,
		width:   40,
	}
}

// SetMessage updates the spinner message
func (l *LoadingSpinner) SetMessage(message string) {
	l.message = message
}

// SetWidth sets the spinner width
func (l *LoadingSpinner) SetWidth(width int) {
	l.width = width
}

// Init initializes the spinner
func (l *LoadingSpinner) Init() tea.Cmd {
	return l.spinner.Tick
}

// Update handles spinner updates
func (l *LoadingSpinner) Update(msg tea.Msg) (*LoadingSpinner, tea.Cmd) {
	var cmd tea.Cmd
	l.spinner, cmd = l.spinner.Update(msg)
	return l, cmd
}

// View renders the spinner
func (l *LoadingSpinner) View() string {
	spinnerView := l.spinner.View()
	
	messageStyle := lipgloss.NewStyle().
		Foreground(styles.White).
		MarginLeft(1)
	
	content := spinnerView + messageStyle.Render(l.message)
	
	containerStyle := lipgloss.NewStyle().
		Width(l.width).
		Align(lipgloss.Center)
	
	return containerStyle.Render(content)
}

// SpinnerTypes provides different spinner styles
var SpinnerTypes = map[string]spinner.Spinner{
	"dot":      spinner.Dot,
	"line":     spinner.Line,
	"minidot":  spinner.MiniDot,
	"jump":     spinner.Jump,
	"pulse":    spinner.Pulse,
	"points":   spinner.Points,
	"globe":    spinner.Globe,
	"moon":     spinner.Moon,
	"monkey":   spinner.Monkey,
	"meter":    spinner.Meter,
	"hamburger": spinner.Hamburger,
}

// SetSpinnerType changes the spinner animation
func (l *LoadingSpinner) SetSpinnerType(name string) {
	if s, ok := SpinnerTypes[name]; ok {
		l.spinner.Spinner = s
	}
}

// MoneySpinner creates a money-themed spinner
func MoneySpinner(message string) *LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"💰", "💸", "💵", "💴", "💶", "💷"},
		FPS:    10,
	}
	s.Style = lipgloss.NewStyle()
	
	return &LoadingSpinner{
		spinner: s,
		message: message,
		width:   40,
	}
}

// RocketSpinner creates a rocket-themed spinner
func RocketSpinner(message string) *LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"🚀", "🚀 ", "🚀  ", "🚀   ", " 🚀  ", "  🚀 ", "   🚀"},
		FPS:    8,
	}
	s.Style = lipgloss.NewStyle()
	
	return &LoadingSpinner{
		spinner: s,
		message: message,
		width:   40,
	}
}

// UnicornSpinner creates a unicorn-themed spinner
func UnicornSpinner(message string) *LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"🦄", "✨🦄", "✨✨🦄", "✨🦄", "🦄"},
		FPS:    6,
	}
	s.Style = lipgloss.NewStyle()
	
	return &LoadingSpinner{
		spinner: s,
		message: message,
		width:   40,
	}
}

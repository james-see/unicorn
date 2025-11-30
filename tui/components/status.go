package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// StatusBar displays game status information at the bottom of the screen
type StatusBar struct {
	width       int
	items       []StatusItem
	helpText    string
	showHelp    bool
}

// StatusItem represents a single status item
type StatusItem struct {
	Label string
	Value string
	Style lipgloss.Style
}

// NewStatusBar creates a new status bar
func NewStatusBar(width int) *StatusBar {
	return &StatusBar{
		width:    width,
		items:    []StatusItem{},
		helpText: "",
		showHelp: true,
	}
}

// SetWidth updates the status bar width
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// SetItems sets the status items
func (s *StatusBar) SetItems(items []StatusItem) {
	s.items = items
}

// AddItem adds a status item
func (s *StatusBar) AddItem(label, value string, style lipgloss.Style) {
	s.items = append(s.items, StatusItem{
		Label: label,
		Value: value,
		Style: style,
	})
}

// ClearItems removes all status items
func (s *StatusBar) ClearItems() {
	s.items = []StatusItem{}
}

// SetHelp sets the help text
func (s *StatusBar) SetHelp(help string) {
	s.helpText = help
}

// SetShowHelp toggles help visibility
func (s *StatusBar) SetShowHelp(show bool) {
	s.showHelp = show
}

// View renders the status bar
func (s *StatusBar) View() string {
	var parts []string

	// Status items
	for _, item := range s.items {
		var itemStr string
		if item.Label != "" {
			labelStyle := lipgloss.NewStyle().Foreground(styles.Gray)
			itemStr = labelStyle.Render(item.Label+": ") + item.Style.Render(item.Value)
		} else {
			itemStr = item.Style.Render(item.Value)
		}
		parts = append(parts, itemStr)
	}

	statusContent := strings.Join(parts, " │ ")

	// Build the bar
	var result strings.Builder

	// Main status line
	if len(s.items) > 0 {
		statusStyle := styles.StatusBarStyle.Width(s.width)
		result.WriteString(statusStyle.Render(statusContent))
		result.WriteString("\n")
	}

	// Help line
	if s.showHelp && s.helpText != "" {
		helpStyle := styles.HelpStyle.Width(s.width)
		result.WriteString(helpStyle.Render(s.helpText))
	}

	return result.String()
}

// GameStatusBar creates a status bar for the VC game
func GameStatusBar(width int, turn, maxTurns int, cash, netWorth int64) *StatusBar {
	bar := NewStatusBar(width)

	bar.SetItems([]StatusItem{
		{
			Label: "Turn",
			Value: fmt.Sprintf("%d/%d", turn, maxTurns),
			Style: lipgloss.NewStyle().Foreground(styles.Yellow),
		},
		{
			Label: "Cash",
			Value: formatMoney(cash),
			Style: lipgloss.NewStyle().Foreground(styles.Green),
		},
		{
			Label: "Net Worth",
			Value: formatMoney(netWorth),
			Style: lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true),
		},
	})

	bar.SetHelp("d dashboard • v value-add • s secondary • enter next turn • q quit")

	return bar
}

// FounderStatusBar creates a status bar for the Founder game
func FounderStatusBar(width int, month int, cash, valuation int64, runway int) *StatusBar {
	bar := NewStatusBar(width)

	runwayStyle := lipgloss.NewStyle().Foreground(styles.Green)
	if runway < 6 {
		runwayStyle = lipgloss.NewStyle().Foreground(styles.Red).Bold(true)
	} else if runway < 12 {
		runwayStyle = lipgloss.NewStyle().Foreground(styles.Yellow)
	}

	bar.SetItems([]StatusItem{
		{
			Label: "Month",
			Value: fmt.Sprintf("%d", month),
			Style: lipgloss.NewStyle().Foreground(styles.Yellow),
		},
		{
			Label: "Cash",
			Value: formatMoney(cash),
			Style: lipgloss.NewStyle().Foreground(styles.Green),
		},
		{
			Label: "Valuation",
			Value: formatMoney(valuation),
			Style: lipgloss.NewStyle().Foreground(styles.Cyan),
		},
		{
			Label: "Runway",
			Value: fmt.Sprintf("%d mo", runway),
			Style: runwayStyle,
		},
	})

	bar.SetHelp("enter next month • q quit")

	return bar
}

// formatMoney formats an int64 as currency
func formatMoney(amount int64) string {
	if amount < 0 {
		return "-$" + formatPositiveMoney(-amount)
	}
	return "$" + formatPositiveMoney(amount)
}

func formatPositiveMoney(amount int64) string {
	if amount >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(amount)/1000000000)
	}
	if amount >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(amount)/1000000)
	}
	if amount >= 1000 {
		return fmt.Sprintf("%.0fK", float64(amount)/1000)
	}
	return fmt.Sprintf("%d", amount)
}

// Header creates a header bar with title
type Header struct {
	title    string
	subtitle string
	width    int
}

// NewHeader creates a new header
func NewHeader(title, subtitle string, width int) *Header {
	return &Header{
		title:    title,
		subtitle: subtitle,
		width:    width,
	}
}

// SetWidth updates the header width
func (h *Header) SetWidth(width int) {
	h.width = width
}

// View renders the header
func (h *Header) View() string {
	var b strings.Builder

	// Title
	titleStyle := styles.HeaderStyle.Width(h.width)
	b.WriteString(titleStyle.Render(h.title))

	// Subtitle
	if h.subtitle != "" {
		b.WriteString("\n")
		subStyle := styles.SubtitleStyle.Width(h.width).Align(lipgloss.Center)
		b.WriteString(subStyle.Render(h.subtitle))
	}

	return b.String()
}

// GameHeader creates a header for the game screen
func GameHeader(width int, playerName, firmName string, difficulty string) *Header {
	title := fmt.Sprintf("🦄 UNICORN VC - %s", firmName)
	subtitle := fmt.Sprintf("%s | %s Mode", playerName, difficulty)
	return NewHeader(title, subtitle, width)
}

// FounderHeader creates a header for the founder screen
func FounderHeader(width int, companyName, playerName string) *Header {
	title := fmt.Sprintf("🚀 %s", companyName)
	subtitle := fmt.Sprintf("Founded by %s", playerName)
	return NewHeader(title, subtitle, width)
}

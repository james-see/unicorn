package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// ProgressBar wraps bubbles/progress with game styling
type ProgressBar struct {
	progress progress.Model
	label    string
	value    float64
	width    int
	showPct  bool
}

// NewProgressBar creates a new progress bar
func NewProgressBar(label string, width int) *ProgressBar {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(width-10),
		progress.WithoutPercentage(),
	)

	return &ProgressBar{
		progress: p,
		label:    label,
		value:    0,
		width:    width,
		showPct:  true,
	}
}

// SetValue sets the progress value (0.0 to 1.0)
func (p *ProgressBar) SetValue(value float64) {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	p.value = value
}

// SetLabel updates the label
func (p *ProgressBar) SetLabel(label string) {
	p.label = label
}

// SetShowPercent toggles percentage display
func (p *ProgressBar) SetShowPercent(show bool) {
	p.showPct = show
}

// SetWidth updates the width
func (p *ProgressBar) SetWidth(width int) {
	p.width = width
	p.progress.Width = width - 10
}

// SetGradient sets custom gradient colors
func (p *ProgressBar) SetGradient(colorA, colorB string) {
	p.progress = progress.New(
		progress.WithGradient(colorA, colorB),
		progress.WithWidth(p.width-10),
		progress.WithoutPercentage(),
	)
}

// Init initializes the progress bar
func (p *ProgressBar) Init() tea.Cmd {
	return nil
}

// Update handles progress bar updates
func (p *ProgressBar) Update(msg tea.Msg) (*ProgressBar, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := p.progress.Update(msg)
	p.progress = model.(progress.Model)
	return p, cmd
}

// View renders the progress bar
func (p *ProgressBar) View() string {
	var b strings.Builder

	// Label
	if p.label != "" {
		labelStyle := lipgloss.NewStyle().
			Foreground(styles.Yellow).
			Width(p.width).
			Align(lipgloss.Left)
		b.WriteString(labelStyle.Render(p.label))
		b.WriteString("\n")
	}

	// Progress bar
	b.WriteString(p.progress.ViewAs(p.value))

	// Percentage
	if p.showPct {
		pctStyle := lipgloss.NewStyle().
			Foreground(styles.Cyan).
			MarginLeft(2)
		b.WriteString(pctStyle.Render(fmt.Sprintf("%.0f%%", p.value*100)))
	}

	return b.String()
}

// TurnProgress creates a progress bar for game turns
func TurnProgress(currentTurn, maxTurns, width int) *ProgressBar {
	p := NewProgressBar(
		fmt.Sprintf("Turn %d of %d", currentTurn, maxTurns),
		width,
	)
	p.SetValue(float64(currentTurn) / float64(maxTurns))
	p.SetGradient("#00FFFF", "#FF00FF")
	return p
}

// XPProgress creates a progress bar for XP/leveling
func XPProgress(currentXP, nextLevelXP, width int) *ProgressBar {
	p := NewProgressBar(
		fmt.Sprintf("XP: %d / %d", currentXP, nextLevelXP),
		width,
	)
	p.SetValue(float64(currentXP) / float64(nextLevelXP))
	p.SetGradient("#00FF00", "#FFFF00")
	return p
}

// RunwayProgress creates a progress bar for startup runway
func RunwayProgress(months, width int) *ProgressBar {
	label := fmt.Sprintf("Runway: %d months", months)
	p := NewProgressBar(label, width)

	// Color based on runway health
	if months < 6 {
		p.SetGradient("#FF0000", "#FF4444")
	} else if months < 12 {
		p.SetGradient("#FFAA00", "#FFFF00")
	} else {
		p.SetGradient("#00FF00", "#00FFAA")
	}

	// Cap at 24 months for display
	maxMonths := 24.0
	value := float64(months) / maxMonths
	if value > 1 {
		value = 1
	}
	p.SetValue(value)
	p.SetShowPercent(false)

	return p
}

// HealthBar creates a health/status bar
type HealthBar struct {
	label      string
	current    int
	max        int
	width      int
	filledChar string
	emptyChar  string
	fillColor  lipgloss.Color
	emptyColor lipgloss.Color
}

// NewHealthBar creates a new health bar
func NewHealthBar(label string, current, max, width int) *HealthBar {
	return &HealthBar{
		label:      label,
		current:    current,
		max:        max,
		width:      width,
		filledChar: "█",
		emptyChar:  "░",
		fillColor:  styles.Green,
		emptyColor: styles.DarkGray,
	}
}

// SetValue updates the current value
func (h *HealthBar) SetValue(current int) {
	h.current = current
	if h.current < 0 {
		h.current = 0
	}
	if h.current > h.max {
		h.current = h.max
	}
}

// SetColors sets the fill and empty colors
func (h *HealthBar) SetColors(fill, empty lipgloss.Color) {
	h.fillColor = fill
	h.emptyColor = empty
}

// View renders the health bar
func (h *HealthBar) View() string {
	var b strings.Builder

	// Label
	if h.label != "" {
		labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
		b.WriteString(labelStyle.Render(h.label + " "))
	}

	// Calculate bar sections
	barWidth := h.width - len(h.label) - 10
	if barWidth < 10 {
		barWidth = 10
	}

	filledWidth := int(float64(h.current) / float64(h.max) * float64(barWidth))
	emptyWidth := barWidth - filledWidth

	// Build bar
	filledStyle := lipgloss.NewStyle().Foreground(h.fillColor)
	emptyStyle := lipgloss.NewStyle().Foreground(h.emptyColor)

	b.WriteString(filledStyle.Render(strings.Repeat(h.filledChar, filledWidth)))
	b.WriteString(emptyStyle.Render(strings.Repeat(h.emptyChar, emptyWidth)))

	// Value
	valueStyle := lipgloss.NewStyle().Foreground(styles.White).MarginLeft(1)
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d/%d", h.current, h.max)))

	return b.String()
}
